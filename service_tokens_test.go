package infisical

import (
	"os"
	"testing"
)

func TestServiceTokens(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	var token, apiKey, e2ee, workspaceID, verbose string
	token = os.Getenv("INFISICAL_TOKEN")
	apiKey = os.Getenv("INFISICAL_API_KEY")
	e2ee = os.Getenv("INFISICAL_E2EE") // NOTE: "enabled" or not
	workspaceID = os.Getenv("INFISICAL_WORKSPACE_ID")
	verbose = os.Getenv("VERBOSE") // NOTE: "true" or not
	var e2eeEnabled = (e2ee == "enabled")
	var isVerbose = (verbose == "true")

	////////////////////////////////
	// initialize client
	var client *Client
	if token == "" || apiKey == "" {
		t.Fatalf("No environment variables: `INFISICAL_TOKEN` or `INFISICAL_API_KEY` were found.")
	} else {
		if e2eeEnabled {
			client = NewE2EEEnabledClient(apiKey, token)
		} else {
			client = NewE2EEDisabledClient(token).SetAPIKey(apiKey)
		}
	}
	if workspaceID == "" {
		t.Fatalf("No environment variable: `INFISICAL_WORKSPACE_ID` was found.")
	}
	client.Verbose = isVerbose

	////////////////////////////////
	// test api functions

	// (retrieve service token)
	if serviceToken, err := client.RetrieveServiceToken(); err != nil {
		t.Errorf("failed to retrieve service token: %s", err)
	} else {
		if serviceToken.Workspace != workspaceID {
			t.Errorf("workspace id of retrieved service token data is not equal to requested one: '%s' vs '%s'", serviceToken.Workspace, workspaceID)
		}
	}
}
