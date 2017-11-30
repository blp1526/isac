package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

type Client struct {
	AccessToken       string
	AccessTokenSecret string
}

func NewClient(accessToken string, accessTokenSecret string) *Client {
	client := &Client{
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
	}
	return client
}

func (client *Client) url(zone string, resource string, id string) (url string) {
	scheme := "https://"
	domain := "secure.sakura.ad.jp"
	path := "/" + path.Join("cloud", "zone", zone, "api", "cloud", "1.1", resource, id)
	url = scheme + domain + path
	return url
}

func (client *Client) Request(method string, zone string, resource string, id string, params []byte) (statusCode int, respBody []byte, err error) {
	url := client.url(zone, resource, id)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(params))
	if err != nil {
		return statusCode, respBody, err
	}

	req.SetBasicAuth(client.AccessToken, client.AccessTokenSecret)
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return statusCode, respBody, err
	}

	statusCode = resp.StatusCode

	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return statusCode, respBody, err
	}

	return statusCode, respBody, nil
}
