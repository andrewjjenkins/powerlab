package megarac

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/glog"
)

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

func (api *Api) loginInternal() (*http.CookieJar, *Session, error) {
	u := "https://" + api.ServerAddr + "/api/session"
	data := url.Values{}
	data.Set("username", api.username)
	data.Set("password", api.password)

	req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	// We do this with a separate client so we don't retry-our-retry
	client := newHttpClientInternal(api.insecureSsl)
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Login failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		return nil, nil, fmt.Errorf("Login authorization failed")
	} else if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("Login failed: %v", res.StatusCode)
	}

	l := loginResponse{}
	err = json.NewDecoder(res.Body).Decode(&l)
	if err != nil {
		return nil, nil, err
	}
	if l.Ok != 0 {
		// "Ok" is 0 to indicate ok. Shrug.
		return nil, nil, fmt.Errorf("Login not ok (%d)", l.Ok)
	}
	if l.CsrfToken == "" {
		return nil, nil, fmt.Errorf("Login response did not contain CSRFToken")
	}
	session := &Session{
		CsrfToken: l.CsrfToken,
		SessionId: fmt.Sprintf("%d", l.RacSessionId),
	}
	return &client.Jar, session, nil
}

func (api *Api) Login(username, password string) error {
	if api.session != nil {
		return errors.New("already logged in")
	}

	api.username = username
	api.password = password

	jar, session, err := api.loginInternal()
	if err != nil {
		return err
	}

	api.client.Jar = *jar
	api.session = session
	glog.Info("Logged in")
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
