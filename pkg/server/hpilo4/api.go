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

/*
func (api *Api) NewRequest(method, path string, body io.Reader) (*http.Request, error) {

}
*/

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
