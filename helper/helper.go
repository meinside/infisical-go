package helper

import (
	"errors"

	"github.com/meinside/infisical-go"
)

// Value returns the secret value for given parameters.
func Value(clientID, clientSecret, workspaceID, environment string, secretType infisical.SecretType, secretKeyPath string) (string, error) {
	client := infisical.NewClientWithoutAPIKey(clientID, clientSecret)

	return client.RetrieveSecretValue(workspaceID, environment, secretType, secretKeyPath)
}

// Values returns multiple secret values for given parameters.
func Values(clientID, clientSecret, workspaceID, environment string, secretType infisical.SecretType, secretKeyPaths []string) (map[string]string, error) {
	client := infisical.NewClientWithoutAPIKey(clientID, clientSecret)

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
