package infisical

import (
	"fmt"
	"net/http"
)

// RetrieveOrganizations retrieves all my organizations.
//
// (DEPRECATED)
//
// https://infisical.com/docs/api-reference/endpoints/users/my-organizations
func (c *Client) RetrieveOrganizations() (result OrganizationsData, err error) {
	path := "/v2/users/me/organizations"

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", path, AuthMethodAPIKeyOnly, nil)
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

	return OrganizationsData{}, fmt.Errorf("failed to retrieve organizations: %s", err)
}
