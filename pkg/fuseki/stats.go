package fuseki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// EndpointStats is the stats for a single endpoint.
type EndpointStats struct {
	RequestsBad  int    `json:"RequestsBad"`
	Requests     int    `json:"Requests"`
	RequestsGood int    `json:"RequestsGood"`
	Operation    string `json:"operation"`
	Description  string `json:"description"`
}

// Stats is the response from the stats endpoint.
type Stats struct {
	Datasets map[string]struct {
		Requests     int                      `json:"Requests"`
		RequestsGood int                      `json:"RequestsGood"`
		RequestsBad  int                      `json:"RequestsBad"`
		Endpoints    map[string]EndpointStats `json:"endpoints"`
	} `json:"datasets"`
}

// GetStats gets the stats from the Apache Jena Fuseki server.
func (c *Client) GetStats() (*Stats, error) {
	endpoint := c.baseURL.ResolveReference(
		&url.URL{
			Path: fmt.Sprintf("/$/stats/%s", c.dataset),
		})
	req := &http.Request{
		Method: http.MethodGet,
		URL:    endpoint,
		Header: map[string][]string{
			"Accept": {"*/*"},
		},
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending stats request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error stats response from server: %s", res.Status)
	}
	defer res.Body.Close()

	var response *Stats
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}
