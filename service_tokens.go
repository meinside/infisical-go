package infisical

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RetrieveServiceToken retrieves service token data.
//
// https://infisical.com/docs/api-reference/endpoints/service-tokens/get
func (c *Client) RetrieveServiceToken(token WorkspaceToken) (result ServiceToken, err error) {
	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", "/v2/service-token/", AuthMethodNormal, &token, nil)
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

	return ServiceToken{}, fmt.Errorf("failed to retrieve service token: %s", err)
}
