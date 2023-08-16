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
	var isVerbose = (verbose == "true")

	////////////////////////////////
	// initialize client
	var client *Client
	if apiKey == "" {
		t.Fatalf("no environment variable: `INFISICAL_API_KEY` was found")
	} else {
		client = NewClient()

		client.SetAPIKey(apiKey)
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
