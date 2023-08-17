package infisical

import (
	"net/http"
	"time"
)

// Client struct
type Client struct {
	apiKey      *string
	token       *string
	e2eeEnabled bool

	httpClient *http.Client

	Verbose bool // NOTE: set for dumping http requests & responses
}

// NewE2EEEnabledClient returns a client struct which is E2EE enabled.
func NewE2EEEnabledClient(apiKey, token string) *Client {
	return &Client{
		apiKey:      &apiKey,
		token:       &token,
		e2eeEnabled: true,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},
	}
}

// NewE2EEDisabledClient returns a client struct which is E2EE disabled.
func NewE2EEDisabledClient(token string) *Client {
	return &Client{
		token:       &token,
		e2eeEnabled: false,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},
	}
}

// SetAPIKey sets the `api_key` value of the client and returns it.
func (c *Client) SetAPIKey(apiKey string) *Client {
	c.apiKey = &apiKey
	return c
}
