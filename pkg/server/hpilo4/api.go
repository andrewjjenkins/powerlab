package hpilo4

import (
	"crypto/tls"
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
