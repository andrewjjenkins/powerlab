package megarac

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	r.Header.Add("Cookie", "QSESSIONID="+api.session.SessionId)
	return r, nil
}

func (api *Api) Get(path string) (*http.Response, error) {
	r, err := api.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return api.client.Do(r)
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

func NewApi(serverAddr string, insecureSsl bool) (*Api, error) {
	transport := http.Transport{}
	if insecureSsl {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	a := Api{
		ServerAddr: serverAddr,
		session:    nil,
		client: &http.Client{
			Timeout:   time.Duration(15 * time.Second),
			Transport: &transport,
		},
	}
	return &a, nil
}

type loginResponse struct {
	Ok int `json:"ok"`
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
	glog.Info("Logged in: %v", l.Ok)
	return nil
}
