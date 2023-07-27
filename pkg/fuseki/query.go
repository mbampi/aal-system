package fuseki

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type QueryResult struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results struct {
		Bindings []struct {
			Label struct {
				Type    string `json:"type"`
				Value   string `json:"value"`
				XMLLang string `json:"xml:lang"`
			} `json:"label"`
			X struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"x"`
		} `json:"bindings"`
	} `json:"results"`
}

// Query sends a SPARQL query to the server.
func (c *Client) Query(query string) (*QueryResult, error) {
	sparqlURL := c.baseURL.ResolveReference(
		&url.URL{
			Path: fmt.Sprintf("/%s/sparql", c.dataset),
		})
	req := &http.Request{
		Method: http.MethodPost,
		URL:    sparqlURL,
		Header: http.Header{"Accept": {"application/json"}},
	}

	encoded := url.QueryEscape(query)
	encoded = strings.ReplaceAll(encoded, "%25", "%")
	req.URL.RawQuery = fmt.Sprintf("query=%s", encoded)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", res.Status)
	}
	defer res.Body.Close()

	var response QueryResult
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// Update sends a SPARQL update query to the server.
func (c *Client) Update(query string) error {
	sparqlURL := c.baseURL.ResolveReference(
		&url.URL{
			Path: fmt.Sprintf("/%s", c.dataset),
		})
	req := &http.Request{
		Method: http.MethodPost,
		URL:    sparqlURL,
		Header: http.Header{
			"Accept":       {"text/plain"},
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
	}

	encoded := url.QueryEscape(query)
	encoded = strings.ReplaceAll(encoded, "%25", "%")
	encoded = strings.ReplaceAll(encoded, "%5Cn", "%0A")
	req.URL.RawQuery = fmt.Sprintf("update=%s", encoded)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		bodyRes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error response from server: %s (%s)", res.Status, bodyRes)
	}
	defer res.Body.Close()

	return nil
}
