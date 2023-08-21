package infisical

import (
	"fmt"
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

			if err = c.parseResponse(res, &result); err == nil {
				return result, nil
			}
		}
	}

	return ProjectsData{}, fmt.Errorf("failed to retrieve projects: %s", err)
}
