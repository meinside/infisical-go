package infisical

import (
	"os"
	"testing"
)

func TestFolders(t *testing.T) {
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

	// list folders
	if folders, err := client.ListFolders(workspaceID, environment, NewParamsListFolders()); err != nil {
		t.Errorf("failed to list folders: %s", err)
	} else {
		numFolders := len(folders.Folders)

		const (
			newFolderName     = "NewFolderForTest"
			updatedFolderName = "UpdatedFolderForTest"
		)

		// create a folder
		if created, err := client.CreateFolder(workspaceID, environment, newFolderName, NewParamsCreateFolder()); err != nil {
			t.Errorf("failed to create a folder: %s", err)
		} else {

			// update a folder
			if updated, err := client.UpdateFolder(workspaceID, environment, created.Folder.ID, updatedFolderName, NewParamsUpdateFolder()); err != nil {
				t.Errorf("failed to update folder: %s", err)
			} else {
				if updated.Folder.Name != updatedFolderName {
					t.Errorf("folder name was not properly updated: %s != %s", updated.Folder.Name, updatedFolderName)
				}
			}

			// delete a folder
			if deleted, err := client.DeleteFolder(workspaceID, environment, created.Folder.ID, NewParamsDeleteFolder()); err != nil {
				t.Errorf("failed to delete folder: %s", err)
			} else {
				if deleted.Folder.ID != created.Folder.ID {
					t.Errorf("deleted folder id differs from the created one: %s != %s", deleted.Folder.ID, created.Folder.ID)
				}
			}

			if folders, err := client.ListFolders(workspaceID, environment, NewParamsListFolders()); err != nil {
				t.Errorf("failed to list folders: %s", err)
			} else {
				if numFolders != len(folders.Folders) {
					t.Errorf("the number of folders differ from the old one")
				}
			}
		}
	}
}
