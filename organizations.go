package infisical

import (
	"fmt"
	"net/http"
)

// OrganizationsData struct for organizations response
type OrganizationsData struct {
	Organizations []Organization `json:"organizations"`
}

// ProjectsData struct for projects response
type ProjectsData struct {
	Workspaces []Workspace `json:"workspaces"`
}

// Workspace struct for project
type Workspace struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Organization string                 `json:"organization"`
	Environments []WorkspaceEnvironment `json:"environments"`
}

// WorkspaceEnvironment struct for environments
type WorkspaceEnvironment struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Organization struct for one organization
//
// (DEPRECATED)
//
// https://infisical.com/docs/api-reference/endpoints/users/my-organizations
type Organization struct {
	AuthEnforced bool   `json:"authEnforced"`
	CreatedAt    string `json:"createdAt"`
	CustomerID   string `json:"customerId"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	UpdatedAt    string `json:"updatedAt"`
}

// RetrieveProjects retrieves all workspaces for given organization id.
//
// https://infisical.com/docs/api-reference/endpoints/organizations/workspaces
func (c *Client) RetrieveProjects(organizationID string) (result ProjectsData, err error) {
	path := fmt.Sprintf("/v2/organizations/%s/workspaces", organizationID)

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", path, AuthMethodNormal, nil)
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

	return ProjectsData{}, fmt.Errorf("failed to retrieve workspaces: %s", err)
}
