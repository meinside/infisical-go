package infisical

import "strings"

// WorkspaceToken struct
type WorkspaceToken struct {
	Token string `json:"token"`
	E2EE  bool   `json:"e2ee"`
}

// type aliases and helper functions

type ParamsRetrieveSecrets map[string]any

func NewParamsRetrieveSecrets() ParamsRetrieveSecrets {
	return ParamsRetrieveSecrets{
		"secretPath": "/",
	}
}

func (p ParamsRetrieveSecrets) SetSecretPath(secretPath string) ParamsRetrieveSecrets {
	if secretPath != "/" {
		secretPath = strings.TrimSuffix(secretPath, "/")
	}
	p["secretPath"] = secretPath
	return p
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

// ServiceToken struct
type ServiceToken struct {
	V__          int      `json:"__v,omitempty"`
	ID_          string   `json:"_id"`
	CreatedAt    string   `json:"createdAt"`
	CreatedBy    string   `json:"createdBy"`
	EncryptedKey *string  `json:"encryptedKey,omitempty"`
	ExpiresAt    *string  `json:"expiresAt,omitempty"`
	ID           string   `json:"id"`
	IV           *string  `json:"iv,omitempty"`
	LastUsed     *string  `json:"lastUsed,omitempty"`
	Name         string   `json:"name"`
	Permissions  []string `json:"permissions"`
	ProjectID    string   `json:"projectId"`
	Scopes       any      `json:"scopes,omitempty"`
	SecretHash   string   `json:"secretHash"`
	Tag          *string  `json:"tag,omitempty"`
	UpdatedAt    string   `json:"updatedAt"`
	User         struct {
		V__         int      `json:"__v,omitempty"`
		ID_         string   `json:"_id"`
		AuthMethods []string `json:"authMethods"`
		CreatedAt   string   `json:"createdAt"`
		Devices     any      `json:"devices,omitempty"`
		Email       string   `json:"email"`
		FirstName   *string  `json:"firstName,omitempty"`
		ID          string   `json:"id"`
		LastName    *string  `json:"lastName,omitempty"`
		MFAMethods  []string `json:"mfaMethods,omitempty"`
		UpdatedAt   string   `json:"updatedAt"`
	} `json:"user"`
	Workspace string `json:"workspace"`
}

// SecretsData struct for secrets response
type SecretsData struct {
	Secrets []Secret `json:"secrets"`
}

// SecretType type and constants
type SecretType string

const (
	SecretTypeShared   SecretType = "shared"
	SecretTypePersonal SecretType = "personal"
)

// SecretData struct for secret response
type SecretData struct {
	Secret Secret `json:"secret"`
}

// Secret struct for one secret
type Secret struct {
	ID_                     string     `json:"_id"`
	Environment             string     `json:"environment"`
	ID                      string     `json:"id"`
	SecretComment           string     `json:"secretComment,omitempty"`
	SecretCommentCiphertext string     `json:"secretCommentCiphertext,omitempty"`
	SecretCommentIV         string     `json:"secretCommentIV,omitempty"`
	SecretCommentTag        string     `json:"secretCommentTag,omitempty"`
	SecretKey               string     `json:"secretKey,omitempty"`
	SecretKeyCiphertext     string     `json:"secretKeyCiphertext,omitempty"`
	SecretKeyIV             string     `json:"secretKeyIV,omitempty"`
	SecretKeyTag            string     `json:"secretKeyTag,omitempty"`
	SecretValue             string     `json:"secretValue,omitempty"`
	SecretValueCiphertext   string     `json:"secretValueCiphertext,omitempty"`
	SecretValueIV           string     `json:"secretValueIV,omitempty"`
	SecretValueTag          string     `json:"secretValueTag,omitempty"`
	Type                    SecretType `json:"type"`
	Version                 int        `json:"version"`
	Workspace               string     `json:"workspace"`
}

// OrganizationsData struct for organizations response
type OrganizationsData struct {
	Organizations []Organization `json:"organizations"`
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
