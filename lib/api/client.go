package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

// Client shows SAKURA Cloud API client.
type Client struct {
	AccessToken       string
	AccessTokenSecret string
}

// NewClient initializes *Client.
func NewClient(accessToken string, accessTokenSecret string) *Client {
	client := &Client{
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
	}
	return client
}

// URL creates an API endpoint.
func (client *Client) URL(zone string, paths []string) (url string) {
	scheme := "https://"
	domain := "secure.sakura.ad.jp"
	last := strings.Join(paths, "/")
	path := "/" + path.Join("cloud", "zone", zone, "api", "cloud", "1.1", last)
	url = scheme + domain + path
	return url
}

// Request creates an API request.
func (client *Client) Request(method string, url string, params []byte) (statusCode int, respBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(params))
	if err != nil {
		return statusCode, respBody, err
	}

	req.SetBasicAuth(client.AccessToken, client.AccessTokenSecret)
	c := &http.Client{Timeout: 30 * time.Second}
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
