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

	Verbose bool // NOTE: set for dumping http requests & responses
}

// NewClient creates a new client and return it.
func NewClient(apiKey string, workspaceTokens map[string]WorkspaceToken) *Client {
	return &Client{
		apiKey:          &apiKey,
		workspaceTokens: workspaceTokens,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},
	}
}

// NewClientWithoutAPIKey creates a new client only with tokens and return it.
func NewClientWithoutAPIKey(workspaceTokens map[string]WorkspaceToken) *Client {
	return &Client{
		workspaceTokens: workspaceTokens,

		httpClient: &http.Client{
			Timeout: TimeoutSeconds * time.Second,
		},
	}
}
