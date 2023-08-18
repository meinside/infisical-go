package helper

import "github.com/meinside/infisical-go"

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
