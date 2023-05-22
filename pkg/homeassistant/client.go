package homeassistant

import (
	"net/http"
	"net/url"
	"time"
)

// Client is a Home Assistant client.
type Client struct {
	httpClient           *http.Client
	baseURL              *url.URL
	longLivedAccessToken string
}

// NewClient creates a new Home Assistant client.
func NewClient() *Client {
	defaultBaseURL := url.URL{
		Scheme: "http",
		Host:   "homeassistant.local:8123",
	}

	return &Client{
		longLivedAccessToken: "",
		baseURL:              &defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetBaseURL sets the base URL for the HTTP client.
func (c *Client) SetBaseURL(baseURL url.URL) {
	c.baseURL = &baseURL
}

// SetLongLivedAccessToken sets the long-lived access token for the HTTP client.
func (c *Client) SetLongLivedAccessToken(longLivedAccessToken string) {
	c.longLivedAccessToken = longLivedAccessToken
}

// SetTimeout sets the timeout for the HTTP client.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// getHeaders returns the headers for the HTTP client,
// which includes the authorization header.
func (c *Client) getHeaders() http.Header {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+c.longLivedAccessToken)
	header.Add("Content-Type", "application/json")
	return header
}

// IsConnected returns true if the client is connected to Home Assistant.
func (c *Client) IsConnected() bool {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/api/"}),
		Header: c.getHeaders(),
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	if res.StatusCode != http.StatusOK {
		return false
	}
	return true
}
