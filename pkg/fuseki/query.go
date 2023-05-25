package fuseki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Query sends a SPARQL query to the server.
func (c *Client) Query(query string) (map[string]any, error) {
	sparqlURL := c.baseURL.ResolveReference(
		&url.URL{
			Path: fmt.Sprintf("/%s/sparql", c.dataset),
		})
	req := &http.Request{
		Method: http.MethodGet,
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

	var response map[string]any
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}
