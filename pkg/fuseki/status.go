package fuseki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Status struct {
	Version       string    `json:"version"`
	Built         time.Time `json:"built"`
	StartDateTime time.Time `json:"startDateTime"`
	Uptime        int       `json:"uptime"`
	Datasets      []struct {
		DsName     string `json:"ds.name"`
		DsState    bool   `json:"ds.state"`
		DsServices []struct {
			SrvType        string   `json:"srv.type"`
			SrvDescription string   `json:"srv.description"`
			SrvEndpoints   []string `json:"srv.endpoints"`
		} `json:"ds.services"`
	} `json:"datasets"`
}

func (c *Client) GetStatus() (*Status, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/$/status"}),
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", res.Status)
	}
	defer res.Body.Close()

	var response Status
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}
