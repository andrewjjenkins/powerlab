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
