package hpilo4

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Api struct {
	ServerAddr  string
	client      *http.Client
	insecureSsl bool
	username    string
	password    string
}

func (api *Api) Name() string {
	return api.ServerAddr
}

func (api *Api) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	u := "https://" + api.ServerAddr + path
	return http.NewRequest(method, u, body)
}

func (api *Api) Do(req *http.Request) (*http.Response, error) {
	return api.client.Do(req)
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
	if contentType != "application/json" && contentType != "application/x-javascript" {
		return fmt.Errorf("unexpected content type for %s: %s", path, contentType)
	}

	err = json.NewDecoder(
		io.LimitReader(res.Body, 1024*1024*10),
	).Decode(obj)
	return err
}

func newHttpClientInternal(insecureSsl bool) *http.Client {
	transport := http.Transport{}
	if insecureSsl {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			CipherSuites: []uint16{
				// This is the best iLO4 can negotiate with modern
				// systems but it isn't included by default in Go
				// anymore unless GODEBUG=tlsrsakex=1
				// https://github.com/golang/go/issues/63413
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &http.Client{
		Timeout:   time.Duration(5 * time.Second),
		Transport: &transport,
		Jar:       jar,
	}
}

func NewApi(serverAddr string, insecureSsl bool) (*Api, error) {
	return &Api{
		ServerAddr:  serverAddr,
		insecureSsl: insecureSsl,
		client:      newHttpClientInternal(insecureSsl),
	}, nil
}
