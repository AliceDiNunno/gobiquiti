package gobiquiti

import (
	"errors"
	"net/http"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CloudKeyInstance struct {
	cookie          *http.Cookie
	cookieLastFetch time.Time
	Config          Config
}

func (c *CloudKeyInstance) cookieRefreshRequired() bool {
	if c.cookie == nil {
		return true
	}

	return time.Since(c.cookieLastFetch).Minutes() > 50 //Ubiquiti's JWT elapses an hour after generation
}

func (c *CloudKeyInstance) Login() error {
	if !c.cookieRefreshRequired() {
		return nil
	}

	c.cookie = nil

	loginEndpoint := "/api/auth/login"

	url := "https://" + c.Config.Hostname + loginEndpoint

	login := loginRequest{
		Username: c.Config.Username,
		Password: c.Config.Password,
	}

	request, err := httpPOST(url, login, nil)

	if err != nil {
		request = nil
		return err
	}

	if len(request.Cookies()) > 0 {
		c.cookieLastFetch = time.Now()
		c.cookie = request.Cookies()[0]
		request = nil
		return nil
	}
	request = nil

	return errors.New("server did not respond with a valid token")
}
