package gobiquiti

import (
	"encoding/json"
	"fmt"
)

type Client struct {
	Id         string `json:"_id"`
	Mac        string `json:"mac"`
	Site       string `json:"site_id"`
	Guest      bool   `json:"is_guest"`
	FirstSeen  int    `json:"first_seen"`
	LastSeen   int    `json:"last_seen"`
	Wired      bool   `json:"is_wired"`
	HostName   string `json:"hostname"`
	DeviceName string `json:"device_name"`
	CustomName string `json:"name"`
}

type ClientStats struct {
	Id              string `json:"_id"`
	UserId          string `json:"user_id"`
	AssocTime       int    `json:"assoc_time"`
	LatestAssocTime int    `json:"latest_assoc_time"`
	IsWired         bool   `json:"is_wired"`
	Rssi            int    `json:"rssi"`
	Ccq             int    `json:"ccq"`
	Noise           int    `json:"noise"`
	Signal          int    `json:"signal"`
	TxPower         int    `json:"tx_power"`
	TxRetries       int    `json:"tx_retries"`
	BytesSent       int    `json:"tx_bytes-r"`
	BytesReceived   int    `json:"rx_bytes-r"`
}

type ClientsStatsResponse struct {
	Response
	Data []ClientStats `json:"data"`
}

type ClientsResponse struct {
	Response
	Data []Client `json:"data"`
}

func (c Client) GetDeviceName() string {
	names := []string{c.CustomName, c.HostName, c.DeviceName}

	for _, name := range names {
		if ValidateHostName(name) {
			names = nil
			return name
		}
	}

	names = nil
	return ""
}

func (c *CloudKeyInstance) GetClients() (*ClientsResponse, error) {
	url := fmt.Sprintf("https://%s/proxy/network/api/s/%s/rest/user", c.Config.Hostname, c.Config.Site)

	serverRequest, err := httpGET(url, c.cookie)

	if err != nil {
		serverRequest = nil
		return nil, err
	}

	var inter ClientsResponse
	decoder := json.NewDecoder(serverRequest.Body)

	err = decoder.Decode(&inter)
	defer serverRequest.Body.Close()

	decoder = nil
	serverRequest = nil
	if err != nil {
		return nil, err
	}
	return &inter, nil
}

func (c *CloudKeyInstance) GetClientsStats() (*ClientsStatsResponse, error) {
	url := fmt.Sprintf("https://%s/proxy/network/api/s/%s/stat/sta", c.Config.Hostname, c.Config.Site)

	serverRequest, err := httpGET(url, c.cookie)

	if err != nil {
		serverRequest = nil
		return nil, err
	}

	var inter ClientsStatsResponse
	decoder := json.NewDecoder(serverRequest.Body)
	err = decoder.Decode(&inter)
	defer serverRequest.Body.Close()

	serverRequest = nil
	decoder = nil
	if err != nil {
		return nil, err
	}
	return &inter, nil
}
