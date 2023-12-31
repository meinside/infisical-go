package infisical

import (
	"os"
	"testing"
)

func TestUsersAndOrganizations(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	var apiKey, verbose string
	apiKey = os.Getenv("INFISICAL_API_KEY")
	verbose = os.Getenv("VERBOSE") // NOTE: "true" or not

	////////////////////////////////
	// initialize client
	if apiKey == "" {
		t.Fatalf("no environment variable: `INFISICAL_API_KEY` was found.")
	}
	client := NewClient(apiKey, nil)
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
