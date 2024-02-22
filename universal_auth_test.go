package infisical

import (
	"os"
	"testing"
)

func TestUniversalAuth(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	apiKey := os.Getenv("INFISICAL_API_KEY")
	workspaceID := os.Getenv("INFISICAL_WORKSPACE_ID")
	clientID := os.Getenv("INFISICAL_CLIENT_ID")
	clientSecret := os.Getenv("INFISICAL_CLIENT_SECRET")
	environment := os.Getenv("INFISICAL_ENVIRONMENT")
	verbose := os.Getenv("VERBOSE") // NOTE: "true" or not

	////////////////////////////////
	// initialize client
	if apiKey == "" || clientID == "" || clientSecret == "" || workspaceID == "" || environment == "" {
		t.Fatalf("no environment variables: `INFISICAL_API_KEY`, `INFISICAL_CLIENT_ID`, `INFISICAL_CLIENT_SECRET`, `INFISICAL_WORKSPACE_ID`, or `INFISICAL_ENVIRONMENT` were found.")
	}
	client := NewClient(apiKey, clientID, clientSecret)
	client.Verbose = (verbose == "true")

	////////////////////////////////
	// test api functions

	// login
	if _, err := client.login(); err != nil {
		t.Errorf("failed to login with universal auth: %s", err)
	}

	// refresh
	if _, err := client.refresh(); err != nil {
		t.Errorf("failed to refresh token with universal auth: %s", err)
	}
}
