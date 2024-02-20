package infisical

import (
	"os"
	"testing"
)

func TestUsersAndOrganizations(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	var apiKey, clientID, clientSecret, verbose string
	apiKey = os.Getenv("INFISICAL_API_KEY")
	clientID = os.Getenv("INFISICAL_CLIENT_ID")
	clientSecret = os.Getenv("INFISICAL_CLIENT_SECRET")
	verbose = os.Getenv("VERBOSE") // NOTE: "true" or not

	////////////////////////////////
	// initialize client
	if apiKey == "" || clientID == "" || clientSecret == "" {
		t.Fatalf("no environment variables: `INFISICAL_API_KEY`, `INFISICAL_CLIENT_ID`, or `INFISICAL_CLIENT_SECRET` were found.")
	}
	client := NewClient(apiKey, clientID, clientSecret)
	client.Verbose = (verbose == "true")

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

			// (retrieve workspaces)
			if projects, err := client.RetrieveProjects(organizationID); err != nil {
				t.Errorf("failed to retrieve workspaces: %s", err)
			} else {
				if len(projects.Workspaces) <= 0 {
					t.Errorf("there were no workspaces in organization with id: %s", organizationID)
				}
			}
		}
	}
}
