package infisical

import "strings"

// WorkspaceToken struct
type WorkspaceToken struct {
	Token string `json:"token"`
	E2EE  bool   `json:"e2ee"`
}

// type aliases

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
	ID           string `json:"_id"`
	CreatedAt    string `json:"createdAt"`
	EncryptedKey string `json:"encryptedKey"`
	Environment  string `json:"environment"`
	ExpiresAt    string `json:"expiresAt"`
	IV           string `json:"iv"`
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	UpdatedAt    string `json:"updatedAt"`
	User         struct {
		ID        string `json:"_id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
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
	ID                      string     `json:"_id"`
	Version                 int        `json:"version"`
	Workspace               string     `json:"workspace"`
	Type                    SecretType `json:"type"`
	Environment             string     `json:"environment"`
	SecretKey               string     `json:"secretKey,omitempty"`
	SecretKeyCiphertext     string     `json:"secretKeyCiphertext,omitempty"`
	SecretKeyIV             string     `json:"secretKeyIV,omitempty"`
	SecretKeyTag            string     `json:"secretKeyTag,omitempty"`
	SecretValue             string     `json:"secretValue,omitempty"`
	SecretValueCiphertext   string     `json:"secretValueCiphertext,omitempty"`
	SecretValueIV           string     `json:"secretValueIV,omitempty"`
	SecretValueTag          string     `json:"secretValueTag,omitempty"`
	SecretComment           string     `json:"secretComment,omitempty"`
	SecretCommentCiphertext string     `json:"secretCommentCiphertext,omitempty"`
	SecretCommentIV         string     `json:"secretCommentIV,omitempty"`
	SecretCommentTag        string     `json:"secretCommentTag,omitempty"`
}

// OrganizationsData struct for organizations response
type OrganizationsData struct {
	Organizations []Organization `json:"organizations"`
}

// Organization struct for one organization
type Organization struct {
	ID         string `json:"_id"`
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
}

// ProjectsData struct for projects response
type ProjectsData struct {
	Workspaces []Workspace `json:"workspaces"`
}

// Workspace struct for project
type Workspace struct {
	ID                 string                 `json:"_id"`
	Name               string                 `json:"name"`
	Organization       string                 `json:"organization"`
	Environments       []WorkspaceEnvironment `json:"environments"`
	AutoCapitalization bool                   `json:"autoCapitalization"`
}

// WorkspaceEnvironment struct for environments
type WorkspaceEnvironment struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
