package infisical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UniversalAuthToken is a struct of universal-auth token
type UniversalAuthToken struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenMaxTTL int64  `json:"accessTokenMaxTTL"`
	ExpiresIn         int64  `json:"expiresIn"`
	TokenType         string `json:"tokenType"`
}

// login without token and load the result into the client
//
// https://infisical.com/docs/api-reference/endpoints/universal-auth/login
func (c *Client) login() (result UniversalAuthToken, err error) {
	var encoded []byte
	if encoded, err = json.Marshal(map[string]any{
		"clientId":     c.clientID,
		"clientSecret": c.clientSecret,
	}); err == nil {
		var req *http.Request
		if req, err = http.NewRequest("POST", c.requestURL("/v1/auth/universal-auth/login"), bytes.NewReader(encoded)); err == nil {
			req.Header.Set("Content-Type", "application/json")

			c.dumpRequest(req)

			var res *http.Response
			if res, err = c.httpClient.Do(req); err == nil {
				c.dumpResponse(res)

				if err = c.parseResponse(res, &result); err == nil {
					c.token = &UniversalAuthToken{
						AccessToken:       result.AccessToken,
						AccessTokenMaxTTL: result.AccessTokenMaxTTL,
						ExpiresIn:         result.ExpiresIn,
						TokenType:         result.TokenType,
					}
					c.tokenExpiresOn = time.Now().Add(time.Second * time.Duration(result.ExpiresIn))

					return result, nil
				}
			}
		}
	}

	return UniversalAuthToken{}, fmt.Errorf("failed to login: %s", err)
}

// refresh access token and load the result into the client
//
// https://infisical.com/docs/api-reference/endpoints/universal-auth/renew-access-token
func (c *Client) refresh() (result UniversalAuthToken, err error) {
	var encoded []byte
	if encoded, err = json.Marshal(map[string]any{
		"accessToken": c.token.AccessToken,
	}); err == nil {
		var req *http.Request
		if req, err = http.NewRequest("POST", c.requestURL("/v1/auth/token/renew"), bytes.NewReader(encoded)); err == nil {
			req.Header.Set("Content-Type", "application/json")

			c.dumpRequest(req)

			var res *http.Response
			if res, err = c.httpClient.Do(req); err == nil {
				c.dumpResponse(res)

				if err = c.parseResponse(res, &result); err == nil {
					c.token = &UniversalAuthToken{
						AccessToken:       result.AccessToken,
						AccessTokenMaxTTL: result.AccessTokenMaxTTL,
						ExpiresIn:         result.ExpiresIn,
						TokenType:         result.TokenType,
					}
					c.tokenExpiresOn = time.Now().Add(time.Second * time.Duration(result.ExpiresIn))

					return result, nil
				}
			}
		}
	}

	return UniversalAuthToken{}, fmt.Errorf("failed to refresh access token: %s", err)
}

/*
TODO:
* Universal Auth
- [ ] [Attach](https://infisical.com/docs/api-reference/endpoints/universal-auth/attach)
- [ ] [Retrieve](https://infisical.com/docs/api-reference/endpoints/universal-auth/retrieve)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/universal-auth/update)
- [ ] [Create Client Secret](https://infisical.com/docs/api-reference/endpoints/universal-auth/create-client-secret)
- [ ] [List Client Secrets](https://infisical.com/docs/api-reference/endpoints/universal-auth/list-client-secrets)
- [ ] [Revoke Client Secret](https://infisical.com/docs/api-reference/endpoints/universal-auth/revoke-client-secret)
- [ ] [Revoke Access Token](https://infisical.com/docs/api-reference/endpoints/universal-auth/revoke-access-token)
*/
