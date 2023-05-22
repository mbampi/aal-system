package sparql

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is a SPARQL client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	dataset    string
}

// NewClient creates a new SPARQL client.
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

// SetBaseURL sets the base URL of the SPARQL server.
func (c *Client) SetBaseURL(baseURL string) error {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("error parsing base URL: %w", err)
	}
	c.baseURL = parsedBaseURL
	return nil
}

// IsConnected returns true if the client is connected to SPARQL.
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

func (c *Client) Query(query string) (map[string]any, error) {
	sparqlURL := c.baseURL.ResolveReference(
		&url.URL{
			Path: fmt.Sprintf("/%s/sparql", c.dataset),
		})
	req := &http.Request{
		Method: http.MethodGet,
		URL:    sparqlURL,
		Header: map[string][]string{
			"Accept": {"*/*"},
		},
	}
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.RawQuery)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", res.Status)
	}
	defer res.Body.Close()
	bytes, _ := io.ReadAll(res.Body)
	fmt.Printf("%s", bytes)

	var response map[string]any
	// err = json.NewDecoder(res.Body).Decode(&response)
	// if err != nil {
	// 	return nil, fmt.Errorf("error decoding response: %w", err)
	// }

	return response, nil
}
