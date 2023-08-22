package helper

import (
	"errors"

	"github.com/meinside/infisical-go"
)

// Value returns the secret value for given parameters.
//
// Works only on E2EE-disabled workspaces.
func Value(workspaceID, token, environment string, secretType infisical.SecretType, secretKeyPath string) (string, error) {
	client := infisical.NewClientWithoutAPIKey(map[string]infisical.WorkspaceToken{
		workspaceID: {
			Token: token,
			E2EE:  false,
		},
	})

	return client.RetrieveSecretValue(workspaceID, environment, secretType, secretKeyPath)
}

// Values returns the secret values for given parameters.
//
// Works only on E2EE-disabled workspaces.
func Values(workspaceID, token, environment string, secretType infisical.SecretType, secretKeyPaths []string) (map[string]string, error) {
	client := infisical.NewClientWithoutAPIKey(map[string]infisical.WorkspaceToken{
		workspaceID: {
			Token: token,
			E2EE:  false,
		},
	})

	var value string
	var err error
	values := map[string]string{}
	errs := []error{}
	for _, secretKeyPath := range secretKeyPaths {
		value, err = client.RetrieveSecretValue(workspaceID, environment, secretType, secretKeyPath)
		if err == nil {
			values[secretKeyPath] = value
		} else {
			errs = append(errs, err)
		}
	}

	return values, errors.Join(errs...)
}

// E2EEValue returns the decrypted secret value for given parameters.
//
// Works only on E2EE-enabled workspaces.
func E2EEValue(apiKey, workspaceID, token, environment string, secretType infisical.SecretType, secretKeyPath string) (string, error) {
	client := infisical.NewClient(apiKey, map[string]infisical.WorkspaceToken{
		workspaceID: {
			Token: token,
			E2EE:  true,
		},
	})

	return client.RetrieveSecretValue(workspaceID, environment, secretType, secretKeyPath)
}

// E2EEValue returns the decrypted secret values for given parameters.
//
// Works only on E2EE-enabled workspaces.
func E2EEValues(apiKey, workspaceID, token, environment string, secretType infisical.SecretType, secretKeyPaths []string) (map[string]string, error) {
	client := infisical.NewClient(apiKey, map[string]infisical.WorkspaceToken{
		workspaceID: {
			Token: token,
			E2EE:  true,
		},
	})

	var value string
	var err error
	values := map[string]string{}
	errs := []error{}
	for _, secretKeyPath := range secretKeyPaths {
		value, err = client.RetrieveSecretValue(workspaceID, environment, secretType, secretKeyPath)
		if err == nil {
			values[secretKeyPath] = value
		} else {
			errs = append(errs, err)
		}
	}

	return values, errors.Join(errs...)
}
