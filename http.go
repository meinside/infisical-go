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
	AuthMethodAPIKey      AuthMethod = 1
	AuthMethodToken       AuthMethod = 2
	AuthMethodPreferToken AuthMethod = 9
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
func (c *Client) newRequestWithQueryParams(method, path string, auth AuthMethod, params map[string]any) (req *http.Request, err error) {
	url := fmt.Sprintf("%s%s", APIBaseURL, path)

	if req, err = http.NewRequest(method, url, nil); err == nil {
		// query parameters
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()

		// add headers for authorization
		token := c.token
		apiKey := c.apiKey
		if token == nil && apiKey == nil {
			return nil, fmt.Errorf("api key and token are missing, cannot generate a request")
		}
		if auth&AuthMethodAPIKey > 0 {
			if apiKey != nil {
				req.Header.Set("X-API-KEY", *apiKey)
			} else {
				return nil, fmt.Errorf("api key is missing, cannot generate a request")
			}
		}
		if auth&AuthMethodToken > 0 {
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
			} else {
				return nil, fmt.Errorf("token is missing, cannot generate a request")
			}
		}
	}

	return req, nil
}

// newRequestWithJSONBody creates a new http request with JSON body.
func (c *Client) newRequestWithJSONBody(method, path string, auth AuthMethod, params map[string]any) (req *http.Request, err error) {
	var encoded []byte
	encoded, err = json.Marshal(params)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", APIBaseURL, path)

	if req, err = http.NewRequest(method, url, bytes.NewReader(encoded)); err == nil {
		req.Header.Set("Content-Type", "application/json")

		// add headers for authorization
		token := c.token
		apiKey := c.apiKey
		if token == nil && apiKey == nil {
			return nil, fmt.Errorf("api key and token are missing, cannot generate a request")
		}
		if auth&AuthMethodPreferToken > 0 {
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
			} else if apiKey != nil {
				req.Header.Set("X-API-KEY", *apiKey)
			}
		} else {
			if auth&AuthMethodAPIKey > 0 {
				if apiKey != nil {
					req.Header.Set("X-API-KEY", *apiKey)
				} else {
					return nil, fmt.Errorf("api key is missing, cannot generate a request")
				}
			}
			if auth&AuthMethodToken > 0 {
				if token != nil {
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
				} else {
					return nil, fmt.Errorf("token is missing, cannot generate a request")
				}
			}
		}
	}

	return req, err
}
