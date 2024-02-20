package infisical

import (
	"fmt"
	"net/http"
)

// RetrieveServiceToken retrieves service token data.
//
// (DEPRECATED)
//
// https://infisical.com/docs/api-reference/endpoints/service-tokens/get
func (c *Client) RetrieveServiceToken(token WorkspaceToken) (result ServiceToken, err error) {
	var req *http.Request
	//req, err = c.newRequestWithQueryParams("GET", "/v2/service-token", AuthMethodNormal, &token, nil)
	req, err = c.newRequestWithQueryParams("GET", "/v2/service-token", AuthMethodTokenOnly, &token, nil)
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

	return ServiceToken{}, fmt.Errorf("failed to retrieve service token: %s", err)
}
