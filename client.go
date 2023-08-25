package infisical

import (
	"net/http"
	"time"
)

// Client struct
type Client struct {
	apiKey          *string
	workspaceTokens map[string]WorkspaceToken

	httpClient *http.Client

	baseURL string

	Verbose bool // NOTE: set `true` for dumping http requests & responses
}

// NewClient creates a new client and return it.
func NewClient(apiKey string, workspaceTokens map[string]WorkspaceToken) *Client {
	return &Client{
		apiKey:          &apiKey,
		workspaceTokens: workspaceTokens,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},

		baseURL: DefaultAPIBaseURL,
	}
}

// NewClientWithoutAPIKey creates a new client only with tokens and return it.
func NewClientWithoutAPIKey(workspaceTokens map[string]WorkspaceToken) *Client {
	return &Client{
		workspaceTokens: workspaceTokens,

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
