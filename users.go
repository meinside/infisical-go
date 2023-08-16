package infisical

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RetrieveOrganizations retrieves all my organizations.
//
// https://infisical.com/docs/api-reference/endpoints/users/my-organizations
func (c *Client) RetrieveOrganizations() (result OrganizationsData, err error) {
	path := fmt.Sprintf("/v2/users/me/organizations")

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", path, AuthMethodAPIKey, nil)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			var body []byte
			if res.StatusCode == 200 {
				if body, err = io.ReadAll(res.Body); err == nil {
					if err = json.Unmarshal(body, &result); err == nil {
						return result, nil
					}
				}
			} else {
				if body, err = io.ReadAll(res.Body); err == nil {
					err = fmt.Errorf("HTTP %d error: `%s`", res.StatusCode, string(body))
				} else {
					err = fmt.Errorf("HTTP %d error", res.StatusCode)
				}
			}
		}
	}

	return OrganizationsData{}, fmt.Errorf("failed to retrieve organizations: %s", err)
}
