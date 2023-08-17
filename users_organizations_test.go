package infisical

import (
	"os"
	"testing"
)

func TestUsersAndOrganizations(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	var token, apiKey, e2ee, verbose string
	token = os.Getenv("INFISICAL_TOKEN")
	apiKey = os.Getenv("INFISICAL_API_KEY")
	e2ee = os.Getenv("INFISICAL_E2EE") // NOTE: "enabled" or not
	verbose = os.Getenv("VERBOSE")     // NOTE: "true" or not
	var e2eeEnabled = (e2ee == "enabled")
	var isVerbose = (verbose == "true")

	////////////////////////////////
	// initialize client
	var client *Client
	if token == "" || apiKey == "" {
		t.Fatalf("no environment variables: `INFISICAL_TOKEN` or `INFISICAL_API_KEY` were found.")
	} else {
		if e2eeEnabled {
			client = NewE2EEEnabledClient(apiKey, token)
		} else {
			client = NewE2EEDisabledClient(token).SetAPIKey(apiKey)
		}
	}
	client.Verbose = isVerbose

	////////////////////////////////
	// test api functions

	// (retrieve organizations)
	if organizations, err := client.RetrieveOrganizations(); err != nil {
		t.Errorf("failed to retrieve organizations: %s", err)
	} else {
		if len(organizations.Organizations) <= 0 {
			t.Errorf("there were no organizations")
		} else {
			organizationID := organizations.Organizations[0].ID

			// (retrieve projects)
			if projects, err := client.RetrieveProjects(organizationID); err != nil {
				t.Errorf("failed to retrieve projects: %s", err)
			} else {
				if len(projects.Workspaces) <= 0 {
					t.Errorf("there were no projects in organization with id: %s", organizationID)
				}
			}
		}
	}
}
