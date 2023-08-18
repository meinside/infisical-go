package infisical

import (
	"os"
	"testing"
)

func TestServiceTokens(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	apiKey := os.Getenv("INFISICAL_API_KEY")
	workspaceID := os.Getenv("INFISICAL_WORKSPACE_ID")
	token := os.Getenv("INFISICAL_TOKEN")
	e2ee := os.Getenv("INFISICAL_E2EE") // NOTE: "enabled" or not
	verbose := os.Getenv("VERBOSE")     // NOTE: "true" or not

	////////////////////////////////
	// initialize client
	if apiKey == "" || token == "" {
		t.Fatalf("No environment variables: `INFISICAL_API_KEY` or `INFISICAL_TOKEN` were found.")
	}
	if workspaceID == "" {
		t.Fatalf("No environment variable: `INFISICAL_WORKSPACE_ID` was found.")
	}
	workspaceToken := WorkspaceToken{
		Token: token,
		E2EE:  (e2ee == "enabled"),
	}
	client := NewClient(apiKey, map[string]WorkspaceToken{
		workspaceID: workspaceToken,
	})
	client.Verbose = (verbose == "true")

	////////////////////////////////
	// test api functions

	// (retrieve service token)
	if serviceToken, err := client.RetrieveServiceToken(workspaceToken); err != nil {
		t.Errorf("failed to retrieve service token: %s", err)
	} else {
		if serviceToken.Workspace != workspaceID {
			t.Errorf("workspace id of retrieved service token data is not equal to requested one: '%s' vs '%s'", serviceToken.Workspace, workspaceID)
		}
	}
}
