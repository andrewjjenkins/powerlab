package hpilo4

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type loginRequestBody struct {
	UserLogin string `json:"user_login"`
	Password  string `json:"password"`
	Method    string `json:"method"`
}

type loginResponseBody struct {
	SessionKey       string `json:"session_key"`
	UserName         string `json:"user_name"`
	UserAccount      string `json:"user_account"`
	UserDn           string `json:"user_dn"`
	UserIp           string `json:"user_ip"`
	UserExpires      string `json:"user_expires"`
	LoginPriv        int    `json:"login_priv"`
	RemoteConsPriv   int    `json:"remote_cons_priv"`
	VirtualMediaPriv int    `json:"virtual_media_priv"`
	ResetPriv        int    `json:"reset_priv"`
	ConfigPriv       int    `json:"config_priv"`
	UserPriv         int    `json:"user_priv"`
}

func (api *Api) Login(username, password string) error {
	u := "https://" + api.ServerAddr + "/json/login_session"
	bodyStruct := loginRequestBody{
		UserLogin: username,
		Password:  password,
		Method:    "login",
	}
	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("login_session returned %d", res.StatusCode)
	}
	contentType := res.Header.Get("content-type")
	if contentType == "" {
		return fmt.Errorf("no content-type")
	}
	if contentType != "application/x-javascript" &&
		contentType != "application/json" {
		return fmt.Errorf("unexpected content-type %s", contentType)
	}
	var resBody loginResponseBody
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return fmt.Errorf("unexpected body: %v", err)
	}

	// Session key is stored in the cookie jar, we don't need to
	// handle it here. (FIXME maybe for logout?)
	return nil
}

func (api *Api) Logout() error {
	return nil
}
