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
	"time"

	"github.com/andrewjjenkins/powerlab/pkg/responsecache"
)

type Session struct {
	SessionId string
	CsrfToken string
}

type Api struct {
	ServerAddr string
	session    *Session
	client     *http.Client
	cache      *responsecache.Cache
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
		cache: responsecache.New(),
	}
	return &a, nil
}
