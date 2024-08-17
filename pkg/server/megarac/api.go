package megarac

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
)

type Session struct {
	SessionId string
	CsrfToken string
}

type Api struct {
	ServerAddr string
	session    *Session
	client     *http.Client
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
	r.Header.Add("X-CSRFTOKEN", api.session.CsrfToken)
	return r, nil
}

func (api *Api) Get(path string) (*http.Response, error) {
	r, err := api.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return api.client.Do(r)
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
	return api.client.Do(r)
}

func (api *Api) Delete(path string) (*http.Response, error) {
	req, err := api.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	return api.client.Do(req)
}

func NewApi(serverAddr string, insecureSsl bool) (*Api, error) {
	transport := http.Transport{}
	if insecureSsl {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	a := Api{
		ServerAddr: serverAddr,
		session:    nil,
		client: &http.Client{
			Timeout:   time.Duration(15 * time.Second),
			Transport: &transport,
			Jar:       jar,
		},
	}
	return &a, nil
}

type loginResponse struct {
	Ok             int    `json:"ok"`
	Privilege      int    `json:"privilege"`
	ExtendedPriv   int    `json:"extendedpriv"`
	RacSessionId   int    `json:"racsession_id"`
	RemoteAddr     string `json:"remote_addr"`
	ServerName     string `json:"server_name"`
	ServerAddr     string `json:"server_addr"`
	HttpsEnabled   int    `json:"HTTPSEnabled"`
	CsrfToken      string `json:"CSRFToken"`
	Channel        int    `json:"channel"`
	PasswordStatus int    `json:"passwordStatus"`
}

func (api *Api) Login(username, password string) error {
	if api.session != nil {
		return errors.New("already logged in")
	}

	u := "https://" + api.ServerAddr + "/api/session"
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	res, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("Login failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		return fmt.Errorf("Login authorization failed")
	} else if res.StatusCode != 200 {
		return fmt.Errorf("Login failed: %v", res.StatusCode)
	}

	l := loginResponse{}
	err = json.NewDecoder(res.Body).Decode(&l)
	if err != nil {
		return err
	}
	if l.Ok != 0 {
		// "Ok" is 0 to indicate ok. Shrug.
		return fmt.Errorf("Login not ok (%d)", l.Ok)
	}
	if l.CsrfToken == "" {
		return fmt.Errorf("Login response did not contain CSRFToken")
	}
	api.session = &Session{
		CsrfToken: l.CsrfToken,
		SessionId: fmt.Sprintf("%d", l.RacSessionId),
	}
	glog.Info("Logged in: %v", l.Ok)
	return nil
}

func (api *Api) Logout() error {
	if api.session == nil {
		return nil
	}

	res, err := api.Delete("/api/session")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("failed logout: %s", res.StatusCode)
	}
	return nil
}
