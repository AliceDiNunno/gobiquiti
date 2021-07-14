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

var cookie *http.Cookie = nil
var cookieLastFetch time.Time

func cookieRefreshRequired() bool {
	if cookie == nil {
		return true
	}

	return time.Since(cookieLastFetch).Minutes() > 50 //Ubiquiti's JWT elapses an hour after generation
}

func Login(server Config) (*http.Cookie, error) {
	if !cookieRefreshRequired() {
		return cookie, nil
	}

	cookie = nil

	loginEndpoint := "/api/auth/login"

	url := "https://" + server.Hostname + loginEndpoint

	login := loginRequest{
		Username: server.Username,
		Password: server.Password,
	}

	request, err := httpPOST(url, login, nil)

	if err != nil {
		request = nil
		return nil, err
	}

	if len(request.Cookies()) > 0 {
		cookieLastFetch = time.Now()
		cookie := request.Cookies()[0]
		request = nil
		return cookie, nil
	}
	request = nil

	return nil, errors.New("server did not respond with a valid token")
}
