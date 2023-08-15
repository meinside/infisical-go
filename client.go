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

// NewClient returns an empty client struct.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},
	}
}

// SetAPIKey sets the api key value of the client.
func (c *Client) SetAPIKey(apiKey string) *Client {
	c.apiKey = &apiKey
	return c
}

// SetToken sets the token value of the client.
func (c *Client) SetToken(token string) *Client {
	c.token = &token
	return c
}

// SetE2EEEnabled sets the e2eeEnabled value of the client.
func (c *Client) SetE2EEEnabled(enabled bool) *Client {
	c.e2eeEnabled = enabled
	return c
}
