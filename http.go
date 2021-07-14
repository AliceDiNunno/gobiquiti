package gobiquiti

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func buildClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return client
}

func buildRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		request = nil
		return nil, err
	}

	request.Header.Add("content-type", "application/json")
	request.Header.Add("Accept", "*/*")

	return request, nil
}

func httpPOST(url string, body interface{}, cookie *http.Cookie) (*http.Response, error) {
	bodyJson, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	client := buildClient()

	bodyBuf := bytes.NewBuffer(bodyJson)
	request, err := buildRequest("POST", url, bodyBuf)
	bodyBuf = nil
	bodyJson = nil

	if err != nil {
		client = nil
		request = nil
		return nil, err
	}

	if cookie != nil {
		request.AddCookie(cookie)
	}

	response, err := client.Do(request)
	request = nil
	client = nil

	return response, err
}

func httpGET(url string, cookie *http.Cookie) (*http.Response, error) {
	client := buildClient()
	request, err := buildRequest("GET", url, nil)

	if err != nil {
		client = nil
		request = nil
		return nil, err
	}

	if cookie != nil {
		request.AddCookie(cookie)
	}

	response, err := client.Do(request)
	request = nil
	client = nil

	return response, err
}
