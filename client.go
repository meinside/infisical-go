package infisical

import (
	"net/http"
	"time"
)

// Client struct
type Client struct {
	// API Key (optional)
	apiKey *string

	// universal-auth token
	clientID       string
	clientSecret   string
	token          *UniversalAuthToken
	tokenExpiresOn time.Time

	httpClient *http.Client

	baseURL string

	Verbose bool // NOTE: set `true` for dumping http requests & responses
}

// NewClient creates a new client and return it.
func NewClient(apiKey, clientID, clientSecret string) *Client {
	return &Client{
		apiKey: &apiKey,

		clientID:     clientID,
		clientSecret: clientSecret,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},

		baseURL: DefaultAPIBaseURL,
	}
}

// NewClientWithoutAPIKey creates and returns a new client only with tokens.
func NewClientWithoutAPIKey(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},

		baseURL: DefaultAPIBaseURL,
	}
}

// SetAPIBaseURL changes the `baseURL`.
//
// (eg. for using in self-hosted infisical servers)
func (c *Client) SetAPIBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// get token, retrieve/refresh it if needed
func (c *Client) getToken() (token *UniversalAuthToken, err error) {
	if c.token == nil {
		_, err = c.login()
	} else {
		if time.Now().After(c.tokenExpiresOn) {
			_, err = c.refresh()
		}
	}

	return c.token, err
}
