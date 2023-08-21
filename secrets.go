package infisical

import (
	"fmt"
	"net/http"
	"strings"
)

// RetrieveSecrets retrieves secrets for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/read
func (c *Client) RetrieveSecrets(workspaceID, environment string, params ParamsRetrieveSecrets) (result SecretsData, err error) {
	var token WorkspaceToken
	var exists bool
	if token, exists = c.workspaceTokens[workspaceID]; !exists {
		return SecretsData{}, fmt.Errorf("`token` for given workspace id was not found: %s", workspaceID)
	}

	if params == nil {
		params = NewParamsRetrieveSecrets()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var path string
	if token.E2EE {
		path = "/v3/secrets/"
	} else {
		path = "/v3/secrets/raw/"
	}

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", path, AuthMethodNormal, &token, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			if err = c.parseResponse(res, &result); err == nil {
				if token.E2EE {
					// decrypt it
					var secrets []Secret
					if secrets, err = c.decryptSecrets(token, result.Secrets); err != nil {
						return SecretsData{}, fmt.Errorf("failed to decrypt retrieved secrets: %s", err)
					}
					return SecretsData{Secrets: secrets}, nil
				} else {
					// return as it is
					return result, nil
				}
			}
		}
	}

	return SecretsData{}, fmt.Errorf("failed to retrieve secrets: %s", err)
}

// RetrieveSecretsAtPath retrieves all secrets at given path.
//
// Just a helper function for `RetrieveSecrets`.
//
// `secretPath` is in form of: "/folder1/folder2/..."
func (c *Client) RetrieveSecretsAtPath(workspaceID, environment, secretPath string) (secrets []Secret, err error) {
	params := NewParamsRetrieveSecrets().
		SetSecretPath(secretPath)

	var retrieved SecretsData
	if retrieved, err = c.RetrieveSecrets(workspaceID, environment, params); err == nil {
		return retrieved.Secrets, nil
	}

	return nil, fmt.Errorf("failed to retrieve secrets at secret path '%s': %s", secretPath, err)
}

// CreateSecret creates a secret with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/create
func (c *Client) CreateSecret(workspaceID, environment, secretKey, secretValue string, params ParamsCreateSecret) (err error) {
	var token WorkspaceToken
	var exists bool
	if token, exists = c.workspaceTokens[workspaceID]; !exists {
		return fmt.Errorf("`token` for given workspace id was not found: %s", workspaceID)
	}

	if params == nil {
		params = NewParamsCreateSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["secretValue"] = secretValue

	var path string
	if token.E2EE {
		path = "/v3/secrets/%s"

		var projectKey []byte
		if projectKey, err = c.projectKey(token); err != nil {
			return err
		}

		// encrypt things
		var encrypted, nonce, authTag []byte
		// (key)
		if encrypted, nonce, authTag, err = encrypt(projectKey, []byte(secretKey)); err != nil {
			return err
		}
		params["secretKeyCiphertext"] = encodeBase64(encrypted)
		params["secretKeyIV"] = encodeBase64(nonce)
		params["secretKeyTag"] = encodeBase64(authTag)
		delete(params, "secretKey")
		// (value)
		if encrypted, nonce, authTag, err = encrypt(projectKey, []byte(secretValue)); err != nil {
			return err
		}
		params["secretValueCiphertext"] = encodeBase64(encrypted)
		params["secretValueIV"] = encodeBase64(nonce)
		params["secretValueTag"] = encodeBase64(authTag)
		delete(params, "secretValue")
		// (comment)
		if comment, exists := params["secretComment"]; exists {
			if comment, ok := comment.(string); ok {
				if encrypted, nonce, authTag, err = encrypt(projectKey, []byte(comment)); err != nil {
					return err
				}
				params["secretCommentCiphertext"] = encodeBase64(encrypted)
				params["secretCommentIV"] = encodeBase64(nonce)
				params["secretCommentTag"] = encodeBase64(authTag)
			}
		}
		delete(params, "secretComment")
	} else {
		path = "/v3/secrets/raw/%s"
	}

	var req *http.Request
	req, err = c.newRequestWithJSONBody("POST", fmt.Sprintf(path, secretKey), AuthMethodNormal, &token, params)
	if err != nil {
		return err
	}

	c.dumpRequest(req)

	var res *http.Response
	res, err = c.httpClient.Do(req)

	c.dumpResponse(res)

	return err
}

// RetrieveSecret retrieves a secret for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/read-one
func (c *Client) RetrieveSecret(workspaceID, environment, secretKey string, params ParamsRetrieveSecret) (result SecretData, err error) {
	var token WorkspaceToken
	var exists bool
	if token, exists = c.workspaceTokens[workspaceID]; !exists {
		return SecretData{}, fmt.Errorf("`token` for given workspace id was not found: %s", workspaceID)
	}

	if params == nil {
		params = NewParamsRetrieveSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var path string
	if token.E2EE {
		path = "/v3/secrets/%s"
	} else {
		path = "/v3/secrets/raw/%s"
	}

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", fmt.Sprintf(path, secretKey), AuthMethodAPIKeyOnly, &token, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			if err = c.parseResponse(res, &result); err == nil {
				if token.E2EE {
					// decrypt it
					var secret Secret
					if secret, err = c.decryptSecret(token, result.Secret); err != nil {
						return SecretData{}, fmt.Errorf("failed to decrypt retrieved secret: %s", err)
					}
					return SecretData{Secret: secret}, nil
				} else {
					// return as it is
					return result, nil
				}
			}
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

	if !emptyString(c.apiKey) {
		params := NewParamsRetrieveSecret().
			SetSecretPath(secretPath).
			SetType(secretType)

		var retrieved SecretData
		if retrieved, err = c.RetrieveSecret(workspaceID, environment, secretKey, params); err == nil {
			return retrieved.Secret.SecretValue, nil
		}

		return "", fmt.Errorf("failed to retrieve secret value for key path '%s': %s", secretKeyWithPath, err)
	} else {
		// FIXME: when E2EE is disabled, `RetrieveSecret` may fail due to missing `api_key`, so use `RetrieveSecrets` instead

		params := NewParamsRetrieveSecrets().SetSecretPath(secretPath)

		var retrieved SecretsData
		if retrieved, err = c.RetrieveSecrets(workspaceID, environment, params); err == nil {
			for _, secret := range retrieved.Secrets {
				if secret.SecretKey == secretKey && secret.Environment == environment && secret.Type == secretType {
					return secret.SecretValue, nil
				}
			}
		}

		return "", fmt.Errorf("failed to retrieve secret value for key path '%s': no matching secret in the result of `RetrieveSecrets`", secretKeyWithPath)
	}
}

// UpdateSecret updates a secret with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/update
func (c *Client) UpdateSecret(workspaceID, environment, secretKey, secretValue string, params ParamsUpdateSecret) (err error) {
	var token WorkspaceToken
	var exists bool
	if token, exists = c.workspaceTokens[workspaceID]; !exists {
		return fmt.Errorf("`token` for given workspace id was not found: %s", workspaceID)
	}

	if params == nil {
		params = NewParamsUpdateSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["secretValue"] = secretValue

	var path string
	if token.E2EE {
		path = "/v3/secrets/%s"
	} else {
		path = "/v3/secrets/raw/%s"
	}

	if token.E2EE {
		var projectKey []byte
		if projectKey, err = c.projectKey(token); err != nil {
			return err
		}

		// set encrypted values
		var encrypted, nonce, authTag []byte
		// (value)
		if encrypted, nonce, authTag, err = encrypt(projectKey, []byte(secretValue)); err != nil {
			return err
		}
		params["secretValueCiphertext"] = encodeBase64(encrypted)
		params["secretValueIV"] = encodeBase64(nonce)
		params["secretValueTag"] = encodeBase64(authTag)
		delete(params, "secretValue")
		// (comment)
		if comment, exists := params["secretComment"]; exists {
			if comment, ok := comment.(string); ok {
				if encrypted, nonce, authTag, err = encrypt(projectKey, []byte(comment)); err != nil {
					return err
				}
				params["secretCommentCiphertext"] = encodeBase64(encrypted)
				params["secretCommentIV"] = encodeBase64(nonce)
				params["secretCommentTag"] = encodeBase64(authTag)
			}
		}
		delete(params, "secretComment")
	}

	var req *http.Request
	req, err = c.newRequestWithJSONBody("PATCH", fmt.Sprintf(path, secretKey), AuthMethodNormal, &token, params)
	if err != nil {
		return err
	}

	c.dumpRequest(req)

	var res *http.Response
	res, err = c.httpClient.Do(req)

	c.dumpResponse(res)

	return err
}

// DeleteSecret deletes a secret for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/delete
func (c *Client) DeleteSecret(workspaceID, environment, secretKey string, params ParamsDeleteSecret) (err error) {
	var token WorkspaceToken
	var exists bool
	if token, exists = c.workspaceTokens[workspaceID]; !exists {
		return fmt.Errorf("`token` for given workspace id was not found: %s", workspaceID)
	}

	if params == nil {
		params = NewParamsDeleteSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var path string
	if token.E2EE {
		path = "/v3/secrets/%s"
	} else {
		path = "/v3/secrets/raw/%s"
	}

	var req *http.Request
	req, err = c.newRequestWithJSONBody("DELETE", fmt.Sprintf(path, secretKey), AuthMethodNormal, &token, params)
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
