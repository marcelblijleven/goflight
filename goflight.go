package goflight

import (
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://opensky-network.org"
)

// Client represents the API client used to communicate with the Opensky API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	username   string
	password   string

	States  *statesService
	Flights *flightService
}

// NewClient creates a new client with the provided credentials
func NewClient(username, password string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Second * 30}
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    u,
		username:   username,
		password:   password,
	}

	c.States = &statesService{client: c}
	c.Flights = &flightService{client: c}

	return c, nil
}

// GetBaseURL returns the base url of the client
func (c *Client) GetBaseURL() *url.URL {
	return c.baseURL
}

// GetCredentials returns the username and password of the client
func (c *Client) GetCredentials() (username string, password string) {
	return c.username, c.password
}

func (c *Client) setBaseURL(u *url.URL) {
	c.baseURL = u
}
