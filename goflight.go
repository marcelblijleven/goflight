package goflight

import (
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://opensky-network.org/api/"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	username   string
	password   string

	States *statesService
}

func NewClient(username, password string, httpClient *http.Client, host *url.URL) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Second * 30}
	}

	if host == nil {
		u, err := url.Parse(BaseURL)
		if err != nil {
			return nil, err
		}

		host = u
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    host,
		username:   username,
		password:   password,
	}

	c.States = &statesService{client: c}

	return c, nil
}

func (c *Client) GetBaseURL() *url.URL {
	return c.baseURL
}

func (c *Client) GetCredentials() (username string, password string) {
	return c.username, c.password
}
