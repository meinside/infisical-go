package infisical

import (
	"fmt"
	"net/http"
)

type ParamsListFolders map[string]any

func NewParamsListFolders() ParamsListFolders {
	return ParamsListFolders{}
}

func (p ParamsListFolders) SetPath(path string) ParamsListFolders {
	p["path"] = path
	return p
}

func (p ParamsListFolders) SetDirectory(directory string) ParamsListFolders {
	p["directory"] = directory
	return p
}

// FoldersData struct for folders response
type FoldersData struct {
	Folders []Folder `json:"folders"`
}

// Folder struct
type Folder struct {
	CreatedAt     string  `json:"createdAt"`
	EnvironmentID string  `json:"envId"`
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	ParentID      *string `json:"parentId,omitempty"`
	UpdatedAt     string  `json:"updatedAt"`
	Version       int     `json:"version,omitempty"`
}

// ListFolders lists folders for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/folders/list
func (c *Client) ListFolders(workspaceID, environment string, params ParamsListFolders) (result FoldersData, err error) {
	if params == nil {
		params = NewParamsListFolders()
	}

	// essential parameters
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", "/v1/folders", AuthMethodNormal, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			if err = c.parseResponse(res, &result); err == nil {
				return result, nil
			}
		}
	}

	return FoldersData{}, fmt.Errorf("failed to list folders: %s", err)
}

type ParamsCreateFolder map[string]any

func NewParamsCreateFolder() ParamsCreateFolder {
	return ParamsCreateFolder{
		"directory": "/",
		"path":      "/",
	}
}

func (p ParamsCreateFolder) SetDirectory(directory string) ParamsCreateFolder {
	p["directory"] = directory
	return p
}

func (p ParamsCreateFolder) SetPath(path string) ParamsCreateFolder {
	p["path"] = path
	return p
}

type FolderData struct {
	Folder Folder `json:"folder"`
}

// CreateFolder creates a new folder with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/folders/create
func (c *Client) CreateFolder(workspaceID, environment, name string, params ParamsCreateFolder) (result FolderData, err error) {
	if params == nil {
		params = NewParamsCreateFolder()
	}

	// essential parameters
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["name"] = name

	var req *http.Request
	req, err = c.newRequestWithJSONBody("POST", "/v1/folders", AuthMethodNormal, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			if err = c.parseResponse(res, &result); err == nil {
				return result, nil
			}
		}
	}

	return FolderData{}, fmt.Errorf("failed to create a folder: %s", err)
}

type ParamsUpdateFolder map[string]any

func NewParamsUpdateFolder() ParamsUpdateFolder {
	return ParamsUpdateFolder{}
}

func (p ParamsUpdateFolder) SetDirectory(directory string) ParamsUpdateFolder {
	p["directory"] = directory
	return p
}

func (p ParamsUpdateFolder) SetPath(path string) ParamsUpdateFolder {
	p["path"] = path
	return p
}

// UpdateFolder updates a folder with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/folders/update
func (c *Client) UpdateFolder(workspaceID, environment, folderID, name string, params ParamsUpdateFolder) (result FolderData, err error) {
	if params == nil {
		params = NewParamsUpdateFolder()
	}

	// essential parameters
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["name"] = name

	var req *http.Request
	req, err = c.newRequestWithJSONBody("PATCH", fmt.Sprintf("/v1/folders/%s", folderID), AuthMethodNormal, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			if err = c.parseResponse(res, &result); err == nil {
				return result, nil
			}
		}
	}

	return FolderData{}, fmt.Errorf("failed to update a folder: %s", err)
}

type ParamsDeleteFolder map[string]any

func NewParamsDeleteFolder() ParamsDeleteFolder {
	return ParamsDeleteFolder{}
}

func (p ParamsDeleteFolder) SetDirectory(directory string) ParamsDeleteFolder {
	p["directory"] = directory
	return p
}

func (p ParamsDeleteFolder) SetPath(path string) ParamsDeleteFolder {
	p["path"] = path
	return p
}

// DeleteFolder deletes a folder with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/folders/delete
func (c *Client) DeleteFolder(workspaceID, environment, folderID string, params ParamsDeleteFolder) (result FolderData, err error) {
	if params == nil {
		params = NewParamsDeleteFolder()
	}

	// essential parameters
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var req *http.Request
	req, err = c.newRequestWithJSONBody("DELETE", fmt.Sprintf("/v1/folders/%s", folderID), AuthMethodNormal, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			if err = c.parseResponse(res, &result); err == nil {
				return result, nil
			}
		}
	}

	return FolderData{}, fmt.Errorf("failed to delete a folder: %s", err)
}
