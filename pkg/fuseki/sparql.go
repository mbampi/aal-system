package fuseki

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is a Apache Jena Fuseki client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	dataset    string
}

// NewClient creates a new Apache Jena Fuseki client.
func NewClient(dataset string) *Client {
	defaultBaseURL := url.URL{
		Scheme: "http",
		Host:   "localhost:3030",
	}

	return &Client{
		baseURL: &defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		dataset: dataset,
	}
}

// SetBaseURL sets the base URL of the Apache Jena Fuseki server.
func (c *Client) SetBaseURL(baseURL string) error {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("error parsing base URL: %w", err)
	}
	c.baseURL = parsedBaseURL
	return nil
}

// IsConnected returns true if the client is connected to Apache Jena Fuseki.
func (c *Client) IsConnected() bool {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/$/ping"}),
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
