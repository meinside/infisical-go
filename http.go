package infisical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	APIBaseURL     = "https://app.infisical.com/api"
	TimeoutSeconds = 10
)

type AuthMethod int

const (
	AuthMethodNormal     AuthMethod = 0
	AuthMethodAPIKeyOnly AuthMethod = 1 << iota
	AuthMethodTokenOnly  AuthMethod = 1 << iota
)

// dump http request
func (c *Client) dumpRequest(req *http.Request) {
	if c.Verbose {
		if bytes, err := httputil.DumpRequest(req, true); err == nil {
			log.Printf("**** dumping HTTP request:\n\n%s\n", string(bytes))
		} else {
			log.Printf("**** failed to dump HTTP request:\n\n%s\n", err)
		}
	}
}

// dump http response
func (c *Client) dumpResponse(req *http.Response) {
	if c.Verbose {
		if bytes, err := httputil.DumpResponse(req, true); err == nil {
			log.Printf("**** dumping HTTP response:\n\n%s\n", string(bytes))
		} else {
			log.Printf("**** failed to dump HTTP response:\n\n%s\n", err)
		}
	}
}

// newRequestWithQueryParams creates a new http request with query strings.
func (c *Client) newRequestWithQueryParams(method, path string, authMethod AuthMethod, token *WorkspaceToken, params map[string]any) (req *http.Request, err error) {
	if authMethod&AuthMethodAPIKeyOnly != 0 && emptyString(c.apiKey) {
		return nil, fmt.Errorf("%s %s requires `api_key` that is missing, cannot generate a request", method, path)
	}
	if authMethod&AuthMethodTokenOnly != 0 && token == nil {
		return nil, fmt.Errorf("%s %s requires `token` that is missing, cannot generate a request", method, path)
	}

	url := fmt.Sprintf("%s%s", APIBaseURL, path)

	if req, err = http.NewRequest(method, url, nil); err == nil {
		// query parameters
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()

		// add headers for authorization
		apiKey := c.apiKey
		if token == nil && emptyString(apiKey) {
			return nil, fmt.Errorf("both `api_key` and `token` are missing, cannot generate a request")
		}
		if authMethod&AuthMethodAPIKeyOnly != 0 {
			if !emptyString(apiKey) {
				req.Header.Set("X-API-KEY", *apiKey)
			} else {
				return nil, fmt.Errorf("`api_key` is missing, cannot generate a request")
			}
		} else if authMethod&AuthMethodTokenOnly != 0 {
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
			} else {
				return nil, fmt.Errorf("`token` is missing, cannot generate a request")
			}
		} else {
			if !emptyString(apiKey) {
				req.Header.Set("X-API-KEY", *apiKey)
			}
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
			}
		}
	}

	return req, nil
}

// newRequestWithJSONBody creates a new http request with JSON body.
func (c *Client) newRequestWithJSONBody(method, path string, authMethod AuthMethod, token *WorkspaceToken, params map[string]any) (req *http.Request, err error) {
	if authMethod&AuthMethodAPIKeyOnly != 0 && emptyString(c.apiKey) {
		return nil, fmt.Errorf("%s %s requires `api_key` that is missing, cannot generate a request", method, path)
	}
	if authMethod&AuthMethodTokenOnly != 0 && token == nil {
		return nil, fmt.Errorf("%s %s requires `token` that is missing, cannot generate a request", method, path)
	}

	var encoded []byte
	encoded, err = json.Marshal(params)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", APIBaseURL, path)

	if req, err = http.NewRequest(method, url, bytes.NewReader(encoded)); err == nil {
		req.Header.Set("Content-Type", "application/json")

		// add headers for authorization
		apiKey := c.apiKey
		if token == nil && emptyString(apiKey) {
			return nil, fmt.Errorf("both `api_key` and `token` are missing, cannot generate a request")
		}
		if authMethod&AuthMethodAPIKeyOnly != 0 {
			if !emptyString(apiKey) {
				req.Header.Set("X-API-KEY", *apiKey)
			} else {
				return nil, fmt.Errorf("`api_key` is missing, cannot generate a request")
			}
		} else if authMethod&AuthMethodTokenOnly != 0 {
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
			} else {
				return nil, fmt.Errorf("`token` is missing, cannot generate a request")
			}
		} else {
			if !emptyString(apiKey) {
				req.Header.Set("X-API-KEY", *apiKey)
			}
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
			}
		}
	}

	return req, err
}

// checks if given string pointer is an empty string
func emptyString(str *string) bool {
	if str == nil || len(*str) == 0 {
		return true
	}
	return false
}
