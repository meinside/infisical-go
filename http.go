package infisical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	DefaultAPIBaseURL = "https://app.infisical.com"
	TimeoutSeconds    = 10
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

// requestURL returns a URL string for HTTP request with given path.
func (c *Client) requestURL(path string) string {
	return fmt.Sprintf("%s/api%s", c.baseURL, path)
}

// newRequestWithQueryParams creates a new http request with query strings.
func (c *Client) newRequestWithQueryParams(method, path string, authMethod AuthMethod, params map[string]any) (req *http.Request, err error) {
	apiKey := c.apiKey
	var token *UniversalAuthToken
	if token, err = c.getToken(); err != nil {
		return nil, fmt.Errorf("failed to fetch token, cannot generated a request: %s", err)
	}

	if authMethod&AuthMethodAPIKeyOnly != 0 && emptyString(c.apiKey) {
		return nil, fmt.Errorf("%s %s requires `api_key` that is missing, cannot generate a request", method, path)
	}
	if authMethod&AuthMethodTokenOnly != 0 && token == nil {
		return nil, fmt.Errorf("%s %s requires `token` that is missing, cannot generate a request", method, path)
	}

	if req, err = http.NewRequest(method, c.requestURL(path), nil); err == nil {
		// query parameters
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()

		// add headers for authorization
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
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			} else {
				return nil, fmt.Errorf("`token` is missing, cannot generate a request")
			}
		} else {
			if !emptyString(apiKey) {
				req.Header.Set("X-API-KEY", *apiKey)
			}
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			}
		}
	}

	return req, nil
}

// newRequestWithJSONBody creates a new http request with JSON body.
func (c *Client) newRequestWithJSONBody(method, path string, authMethod AuthMethod, params map[string]any) (req *http.Request, err error) {
	apiKey := c.apiKey
	var token *UniversalAuthToken
	if token, err = c.getToken(); err != nil {
		return nil, fmt.Errorf("failed to fetch token, cannot generated a request: %s", err)
	}

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

	if req, err = http.NewRequest(method, c.requestURL(path), bytes.NewReader(encoded)); err == nil {
		req.Header.Set("Content-Type", "application/json")

		// add headers for authorization
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
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			} else {
				return nil, fmt.Errorf("`token` is missing, cannot generate a request")
			}
		} else {
			if !emptyString(apiKey) {
				req.Header.Set("X-API-KEY", *apiKey)
			}
			if token != nil {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			}
		}
	}

	return req, err
}

// parse response into given interface
func (c *Client) parseResponse(res *http.Response, into any) (err error) {
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == 200 {
		if into != nil {
			return json.NewDecoder(res.Body).Decode(&into)
		} else {
			return nil
		}
	} else {
		var body []byte
		if body, err = io.ReadAll(res.Body); err == nil {
			err = fmt.Errorf("%s: `%s`", httpStatusToErr(res.StatusCode), string(body))
		} else {
			err = httpStatusToErr(res.StatusCode)
		}
	}

	return err
}

// convert HTTP status code to a meaningful error
func httpStatusToErr(status int) error {
	httpError := fmt.Sprintf("HTTP %d", status)

	switch status {
	case 400:
		return fmt.Errorf("%s; bad request", httpError)
	case 401:
		return fmt.Errorf("%s; unauthorized", httpError)
	case 403:
		return fmt.Errorf("%s; forbidden", httpError)
	case 404:
		return fmt.Errorf("%s; not found", httpError)
	case 500:
		return fmt.Errorf("%s; internal server error", httpError)
	case 503:
		return fmt.Errorf("%s; service unavailable", httpError)
	}

	// fallback
	return fmt.Errorf(httpError)
}

// checks if given string pointer is an empty string
func emptyString(str *string) bool {
	if str == nil || len(*str) == 0 {
		return true
	}
	return false
}
