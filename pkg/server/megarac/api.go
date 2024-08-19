package megarac

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/andrewjjenkins/powerlab/pkg/responsecache"
)

type Session struct {
	SessionId string
	CsrfToken string
}

type Api struct {
	ServerAddr  string
	session     *Session
	client      *http.Client
	relogin     sync.Mutex
	cache       *responsecache.Cache
	insecureSsl bool
	username    string
	password    string
}

func (api *Api) Name() string {
	return api.ServerAddr
}

func (api *Api) SessionId() string {
	if api.session == nil {
		return ""
	}
	return api.session.SessionId
}

func (api *Api) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	if api.session == nil {
		return nil, errors.New("no session")
	}

	u := "https://" + api.ServerAddr + path
	r, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("X-CSRFTOKEN", api.session.CsrfToken)
	return r, nil
}

func (api *Api) Do(req *http.Request) (*http.Response, error) {
	// FIXME: Maybe a body size limit and don't retry huge requests?
	var bodyBytes []byte
	var err error
	var hijackedBody io.ReadCloser
	if req.Body != nil && req.Body != http.NoBody {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		hijackedBody = req.Body
		defer hijackedBody.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	res, err := api.client.Do(req)

	shouldRetry := api.checkRetry(req, res, err)
	if !shouldRetry {
		slog.Debug("not retrying request", "status", res.StatusCode, "error", err)
		return res, err
	}

	slog.Warn("retrying request", "status", res.StatusCode, "error", err)

	// Reset request body, CSRFToken and cookies
	if bodyBytes != nil {
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	req.Header.Set("X-CSRFTOKEN", api.session.CsrfToken)
	req.Header.Del("cookie")
	for _, cookie := range api.client.Jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	// Replay request
	res, err = api.client.Do(req)
	slog.Warn("retried request", "status", res.StatusCode, "error", err)
	return res, err
}

func (api *Api) Get(path string) (*http.Response, error) {
	r, err := api.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return api.Do(r)
}

func (api *Api) GetJson(path string, obj interface{}) error {
	res, err := api.Get(path)
	if err != nil {
		return fmt.Errorf("failed getting %s: %v", path, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed getting %s: %v", path, res.StatusCode)
	}
	contentType := res.Header.Get("content-type")
	if contentType != "application/json" {
		return fmt.Errorf("unexpected content type for %s: %s", path, contentType)
	}

	err = json.NewDecoder(
		io.LimitReader(res.Body, 1024*1024*10),
	).Decode(obj)
	return err
}

func (api *Api) Post(path string, data interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(body)

	r, err := api.NewRequest("POST", path, bodyReader)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-type", "application/json")
	return api.Do(r)
}

func (api *Api) Delete(path string) (*http.Response, error) {
	req, err := api.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	return api.Do(req)
}

func (api *Api) checkRetry(
	req *http.Request, res *http.Response, err error) bool {
	if res.StatusCode != 401 {
		slog.Debug("relogin: error code not fixable", "error code", res.StatusCode)
		return false
	}
	if req.URL.Path == "/api/session" {
		slog.Debug(
			"relogin: cannot retry a login request",
			"error code", res.StatusCode,
			"method", req.Method,
			"path", req.URL.Path,
		)
	}

	if !api.relogin.TryLock() {
		slog.Warn("relogin: already in progress")
		return false
	}
	defer api.relogin.Unlock()

	jar, session, err := api.loginInternal()
	if err != nil {
		slog.Warn("relogin failed", "error", err)
		return false
	}
	api.client.Jar = *jar
	api.session = session

	slog.Info("Relogin successful, signaling retry")

	// Megarac is not ready to use the session immediately after login.
	time.Sleep(200 * time.Millisecond)

	return true
}

func newHttpClientInternal(insecureSsl bool) *http.Client {
	transport := http.Transport{}
	if insecureSsl {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &http.Client{
		Timeout:   time.Duration(15 * time.Second),
		Transport: &transport,
		Jar:       jar,
	}
}

func NewApi(serverAddr string, insecureSsl bool) (*Api, error) {
	a := Api{
		ServerAddr:  serverAddr,
		session:     nil,
		cache:       responsecache.New(),
		insecureSsl: insecureSsl,
		client:      newHttpClientInternal(insecureSsl),
	}

	return &a, nil
}
