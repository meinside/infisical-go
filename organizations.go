package infisical

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RetrieveProjects retrieves all projects for given organization id.
//
// https://infisical.com/docs/api-reference/endpoints/organizations/workspaces
func (c *Client) RetrieveProjects(organizationID string) (result ProjectsData, err error) {
	path := fmt.Sprintf("/v2/organizations/%s/workspaces", organizationID)

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", path, AuthMethodAPIKeyOnly, nil, nil)
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

	return ProjectsData{}, fmt.Errorf("failed to retrieve projects: %s", err)
}
