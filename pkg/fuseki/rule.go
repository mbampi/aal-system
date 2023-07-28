package fuseki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// UploadRule uploads a SWRL rule to the Fuseki server.
// TODO review this code
func (c *Client) UploadRule(swrlRule string) error {
	bodyJSON, err := json.Marshal(map[string]string{"update": swrlRule})
	if err != nil {
		return fmt.Errorf("error marshalling SWRL rule: %w", err)
	}

	req := &http.Request{
		Method: http.MethodPut,
		URL:    c.baseURL.ResolveReference(&url.URL{Path: "/$/server"}),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(bodyJSON)),
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending SWRL rule: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error uploading SWRL rule: %s", res.Status)
	}
	return nil
}
