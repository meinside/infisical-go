package infisical

import (
	"os"
	"testing"
)

const (
	secretType           = SecretTypeShared
	secretKey            = "new_secret_key"
	secretValue          = "newly_created_secret_value"
	secretComment        = "newly_created_secret_comment"
	secretValueUpdated   = "updated_secret_value"
	secretCommentUpdated = "updated_secret_comment"
)

func TestSecrets(t *testing.T) {
	////////////////////////////////
	// read values from environment variables
	var token, apiKey, e2ee, workspaceID, env, verbose string
	token = os.Getenv("INFISICAL_TOKEN")
	apiKey = os.Getenv("INFISICAL_API_KEY")
	e2ee = os.Getenv("INFISICAL_E2EE") // NOTE: "enabled" or not
	workspaceID = os.Getenv("INFISICAL_WORKSPACE_ID")
	env = os.Getenv("INFISICAL_ENVIRONMENT") // NOTE: default: "dev"
	verbose = os.Getenv("VERBOSE")           // NOTE: "true" or not
	var e2eeEnabled = (e2ee == "enabled")
	var environment = "dev"
	if env != "" {
		environment = env
	}
	var isVerbose = (verbose == "true")

	////////////////////////////////
	// initialize client
	var client *Client
	if token == "" || apiKey == "" {
		t.Fatalf("no environment variables: `INFISICAL_TOKEN` or `INFISICAL_API_KEY` were found.")
	} else {
		if e2eeEnabled {
			client = NewE2EEEnabledClient(apiKey, token)
		} else {
			client = NewE2EEDisabledClient(token).SetAPIKey(apiKey)
		}
	}
	if workspaceID == "" {
		t.Fatalf("no environment variable: `INFISICAL_WORKSPACE_ID` was found.")
	}
	client.Verbose = isVerbose

	////////////////////////////////
	// test api functions

	// (retrieve all secrets)
	var initialNumSecrets int
	if secrets, err := client.RetrieveSecrets(
		workspaceID,
		environment,
		NewParamsRetrieveSecrets(),
	); err != nil {
		t.Errorf("failed to retrieve secrets: %s", err)
	} else {
		initialNumSecrets = len(secrets.Secrets)

		// (create a secret)
		if err := client.CreateSecret(
			secretKey,
			workspaceID,
			environment,
			secretValue,
			NewParamsCreateSecret().
				SetType(secretType).
				SetSecretComment(secretComment),
		); err != nil {
			t.Errorf("failed to create a secret: %s", err)
		}

		// (retrieve a secret)
		if secret, err := client.RetrieveSecret(
			secretKey,
			workspaceID,
			environment,
			NewParamsRetrieveSecret().
				SetType(secretType),
		); err != nil {
			t.Errorf("failed to retrieve a secret: %s", err)
		} else {
			if secret.Secret.SecretValue != secretValue {
				t.Errorf("newly-created secret value is not equal to the requested one: '%s' vs '%s'", secret.Secret.SecretValue, secretValue)
			} else {
				// (update a secret)
				if err := client.UpdateSecret(
					secretKey,
					workspaceID,
					environment,
					secretValueUpdated,
					NewParamsUpdateSecret().
						SetType(secretType).
						SetSecretComment(secretCommentUpdated),
				); err != nil {
					t.Errorf("failed to update a secret: %s", err)
				} else {
					// (retrieve a secret)
					if secret, err := client.RetrieveSecret(
						secretKey,
						workspaceID,
						environment,
						NewParamsRetrieveSecret().
							SetType(secretType),
					); err != nil {
						t.Errorf("failed to retrieve an updated secret: %s", err)
					} else {
						if secret.Secret.SecretValue != secretValueUpdated {
							t.Errorf("retrieved secret value is not equal to the updated one: '%s' vs '%s'", secret.Secret.SecretValue, secretValueUpdated)
						} else {
							// (delete the newly-created & updated secret)
							if err := client.DeleteSecret(
								secret.Secret.SecretKey,
								secret.Secret.Workspace,
								secret.Secret.Environment,
								NewParamsDeleteSecret().
									SetType(secret.Secret.Type),
							); err != nil {
								t.Errorf("failed to delete secret: %s", err)
							}
						}
					}
				}
			}
		}

		// (retrieve all secrets)
		if secrets, err := client.RetrieveSecrets(
			workspaceID,
			environment,
			NewParamsRetrieveSecrets(),
		); err != nil {
			t.Errorf("failed to retrieve secrets: %s", err)
		} else {
			if len(secrets.Secrets) != initialNumSecrets {
				t.Errorf("the number of remaining secrets: %d is not equal to the initial count: %d", len(secrets.Secrets), initialNumSecrets)
			}
		}
	}
}
