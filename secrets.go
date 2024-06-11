package infisical

import (
	"fmt"
	"net/http"
	"strings"
)

// SecretType type and constants
type SecretType string

const (
	SecretTypeShared   SecretType = "shared"
	SecretTypePersonal SecretType = "personal"
)

type SecretImport struct {
	Environment string   `json:"environment"`
	FolderID    *string  `json:"folderId,omitempty"`
	SecretPath  string   `json:"secretPath"`
	Secrets     []Secret `json:"secrets"`
}

// Secret struct for one secret
type Secret struct {
	ID_           string     `json:"_id"`
	Environment   string     `json:"environment"`
	ID            string     `json:"id"`
	SecretComment *string    `json:"secretComment,omitempty"`
	SecretKey     string     `json:"secretKey"`
	SecretValue   string     `json:"secretValue"`
	Type          SecretType `json:"type"`
	Version       int        `json:"version"`
	Workspace     string     `json:"workspace"`
}

type ParamsListSecrets map[string]any

func NewParamsListSecrets() ParamsListSecrets {
	return ParamsListSecrets{}
}

func (p ParamsListSecrets) SetWorkspaceID(workspaceID string) ParamsListSecrets {
	p["workspaceId"] = workspaceID
	return p
}

func (p ParamsListSecrets) SetEnvironment(environment string) ParamsListSecrets {
	p["environment"] = environment
	return p
}

func (p ParamsListSecrets) SetSecretPath(secretPath string) ParamsListSecrets {
	p["secretPath"] = secretPath
	return p
}

func (p ParamsListSecrets) SetIncludeImports(includeImports bool) ParamsListSecrets {
	p["include_imports"] = includeImports
	return p
}

// SecretsData struct for secrets response
type SecretsData struct {
	Imports []SecretImport `json:"imports"`
	Secrets []Secret       `json:"secrets"`
}

// ListSecrets lists all secrets for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/list
func (c *Client) ListSecrets(params ParamsListSecrets) (result SecretsData, err error) {
	if params == nil {
		params = NewParamsListSecrets()
	}

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", "/v3/secrets/raw", AuthMethodNormal, params)
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

	return SecretsData{}, fmt.Errorf("failed to list secrets: %s", err)
}

type ParamsCreateSecret map[string]any

func NewParamsCreateSecret() ParamsCreateSecret {
	return ParamsCreateSecret{
		"secretPath": "/",
		"type":       SecretTypeShared,
	}
}

func (p ParamsCreateSecret) SetSecretComment(secretComment string) ParamsCreateSecret {
	p["secretComment"] = secretComment
	return p
}

func (p ParamsCreateSecret) SetSecretPath(secretPath string) ParamsCreateSecret {
	if secretPath != "/" {
		secretPath = strings.TrimSuffix(secretPath, "/")
	}
	p["secretPath"] = secretPath
	return p
}

func (p ParamsCreateSecret) SetType(typ SecretType) ParamsCreateSecret {
	p["type"] = typ
	return p
}

// CreateSecret creates a secret with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/create
func (c *Client) CreateSecret(workspaceID, environment, secretKey, secretValue string, params ParamsCreateSecret) (err error) {
	if params == nil {
		params = NewParamsCreateSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["secretValue"] = secretValue

	var req *http.Request
	req, err = c.newRequestWithJSONBody("POST", fmt.Sprintf("/v3/secrets/raw/%s", secretKey), AuthMethodNormal, params)
	if err != nil {
		return err
	}

	c.dumpRequest(req)

	var res *http.Response
	if res, err = c.httpClient.Do(req); err == nil {
		c.dumpResponse(res)

		return c.parseResponse(res, nil)
	}

	return err
}

type ParamsRetrieveSecret map[string]any

func NewParamsRetrieveSecret() ParamsRetrieveSecret {
	return ParamsRetrieveSecret{
		"secretPath": "/",
		"type":       SecretTypePersonal,
	}
}

func (p ParamsRetrieveSecret) SetSecretPath(secretPath string) ParamsRetrieveSecret {
	if secretPath != "/" {
		secretPath = strings.TrimSuffix(secretPath, "/")
	}
	p["secretPath"] = secretPath
	return p
}

func (p ParamsRetrieveSecret) SetType(typ SecretType) ParamsRetrieveSecret {
	p["type"] = typ
	return p
}

// SecretData struct for secret response
type SecretData struct {
	Secret Secret `json:"secret"`
}

// RetrieveSecret retrieves a secret for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/read
func (c *Client) RetrieveSecret(workspaceID, environment, secretKey string, params ParamsRetrieveSecret) (result SecretData, err error) {
	if params == nil {
		params = NewParamsRetrieveSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", fmt.Sprintf("/v3/secrets/raw/%s", secretKey), AuthMethodNormal, params)
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

	return SecretData{}, fmt.Errorf("failed to retrieve secret: %s", err)
}

// RetrieveSecretValue retrieves a secret value for given path + key.
//
// Just a helper function for `RetrieveSecret`.
//
// `secretKeyWithPath` is in form of: "/folder1/folder2/.../secret_key_name"
func (c *Client) RetrieveSecretValue(workspaceID, environment string, secretType SecretType, secretKeyWithPath string) (value string, err error) {
	// secretKeyWithPath => secretKey + secretPath
	splitted := strings.Split(secretKeyWithPath, "/")
	secretKey := splitted[len(splitted)-1]
	secretPath := strings.TrimSuffix(secretKeyWithPath, secretKey)

	params := NewParamsRetrieveSecret().
		SetSecretPath(secretPath).
		SetType(secretType)

	var retrieved SecretData
	if retrieved, err = c.RetrieveSecret(workspaceID, environment, secretKey, params); err == nil {
		return retrieved.Secret.SecretValue, nil
	}

	return "", fmt.Errorf("failed to retrieve secret value for key path '%s': %s", secretKeyWithPath, err)
}

type ParamsUpdateSecret map[string]any

func NewParamsUpdateSecret() ParamsUpdateSecret {
	return ParamsUpdateSecret{
		"secretPath": "/",
		"type":       SecretTypeShared,
	}
}

func (p ParamsUpdateSecret) SetSecretPath(secretPath string) ParamsUpdateSecret {
	if secretPath != "/" {
		secretPath = strings.TrimSuffix(secretPath, "/")
	}
	p["secretPath"] = secretPath
	return p
}

func (p ParamsUpdateSecret) SetType(typ SecretType) ParamsUpdateSecret {
	p["type"] = typ
	return p
}

func (p ParamsUpdateSecret) SetSecretComment(comment string) ParamsUpdateSecret {
	p["secretComment"] = comment
	return p
}

// UpdateSecret updates a secret with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/update
func (c *Client) UpdateSecret(workspaceID, environment, secretKey, secretValue string, params ParamsUpdateSecret) (err error) {
	if params == nil {
		params = NewParamsUpdateSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["secretValue"] = secretValue

	var req *http.Request
	req, err = c.newRequestWithJSONBody("PATCH", fmt.Sprintf("/v3/secrets/raw/%s", secretKey), AuthMethodNormal, params)
	if err != nil {
		return err
	}

	c.dumpRequest(req)

	var res *http.Response
	if res, err = c.httpClient.Do(req); err == nil {
		c.dumpResponse(res)

		return c.parseResponse(res, nil)
	}

	return err
}

type ParamsDeleteSecret map[string]any

func NewParamsDeleteSecret() ParamsDeleteSecret {
	return ParamsDeleteSecret{
		"secretPath": "/",
		"type":       SecretTypePersonal,
	}
}

func (p ParamsDeleteSecret) SetSecretPath(secretPath string) ParamsDeleteSecret {
	if secretPath != "/" {
		secretPath = strings.TrimSuffix(secretPath, "/")
	}
	p["secretPath"] = secretPath
	return p
}

func (p ParamsDeleteSecret) SetType(typ SecretType) ParamsDeleteSecret {
	p["type"] = typ
	return p
}

// DeleteSecret deletes a secret for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/delete
func (c *Client) DeleteSecret(workspaceID, environment, secretKey string, params ParamsDeleteSecret) (err error) {
	if params == nil {
		params = NewParamsDeleteSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var req *http.Request
	req, err = c.newRequestWithJSONBody("DELETE", fmt.Sprintf("/v3/secrets/raw/%s", secretKey), AuthMethodNormal, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			return c.parseResponse(res, nil)
		}
	}

	return err
}
