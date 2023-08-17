package infisical

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RetrieveSecrets retrieves secrets for given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/read
func (c *Client) RetrieveSecrets(workspaceID, environment string, params ParamsRetrieveSecrets) (result SecretsData, err error) {
	if params == nil {
		params = NewParamsRetrieveSecrets()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var path string
	if c.e2eeEnabled {
		path = "/v3/secrets/"
	} else {
		path = "/v3/secrets/raw/"
	}

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", path, AuthMethodPreferToken, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			var body []byte
			if res.StatusCode == 200 {
				if body, err = io.ReadAll(res.Body); err == nil {
					if err = json.Unmarshal(body, &result); err == nil {
						if c.e2eeEnabled {
							// decrypt it
							var secrets []Secret
							if secrets, err = c.decryptSecrets(result.Secrets); err != nil {
								return SecretsData{}, fmt.Errorf("failed to decrypt retrieved secrets: %s", err)
							}
							return SecretsData{Secrets: secrets}, nil
						} else {
							// return as it is
							return result, nil
						}
					}
				}
			} else {
				if body, err = io.ReadAll(res.Body); err == nil {
					err = fmt.Errorf("HTTP %d error: `%s`", res.StatusCode, string(body))
				} else {
					err = fmt.Errorf("HTTP %d error", res.StatusCode)
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
func (c *Client) RetrieveSecretsAtPath(secretPath, workspaceID, environment string) (secrets []Secret, err error) {
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
func (c *Client) CreateSecret(secretKey, workspaceID, environment, secretValue string, params ParamsCreateSecret) (err error) {
	if params == nil {
		params = NewParamsCreateSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["secretValue"] = secretValue

	var path string
	if c.e2eeEnabled {
		path = "/v3/secrets/%s"

		var projectKey []byte
		if projectKey, err = c.projectKey(); err != nil {
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
	req, err = c.newRequestWithJSONBody("POST", fmt.Sprintf(path, secretKey), AuthMethodPreferToken, params)
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
func (c *Client) RetrieveSecret(secretKey, workspaceID, environment string, params ParamsRetrieveSecret) (result SecretData, err error) {
	if params == nil {
		params = NewParamsRetrieveSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var path string
	if c.e2eeEnabled {
		path = "/v3/secrets/%s"
	} else {
		path = "/v3/secrets/raw/%s"
	}

	var req *http.Request
	req, err = c.newRequestWithQueryParams("GET", fmt.Sprintf(path, secretKey), AuthMethodPreferToken, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			var body []byte
			if res.StatusCode == 200 {
				if body, err = io.ReadAll(res.Body); err == nil {
					if err = json.Unmarshal(body, &result); err == nil {
						if c.e2eeEnabled {
							// decrypt it
							var secret Secret
							if secret, err = c.decryptSecret(result.Secret); err != nil {
								return SecretData{}, fmt.Errorf("failed to decrypt retrieved secret: %s", err)
							}
							return SecretData{Secret: secret}, nil
						} else {
							// return as it is
							return result, nil
						}
					}
				}
			} else {
				if body, err = io.ReadAll(res.Body); err == nil {
					err = fmt.Errorf("HTTP %d error: `%s`", res.StatusCode, string(body))
				} else {
					err = fmt.Errorf("HTTP %d error", res.StatusCode)
				}
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
func (c *Client) RetrieveSecretValue(secretKeyWithPath, workspaceID, environment string, secretType SecretType) (value string, err error) {
	// secretKeyWithPath => secretKey + secretPath
	splitted := strings.Split(secretKeyWithPath, "/")
	secretKey := splitted[len(splitted)-1]
	secretPath := strings.TrimSuffix(secretKeyWithPath, secretKey)

	params := NewParamsRetrieveSecret().
		SetSecretPath(secretPath).
		SetType(secretType)

	var retrieved SecretData
	if retrieved, err = c.RetrieveSecret(secretKey, workspaceID, environment, params); err == nil {
		return retrieved.Secret.SecretValue, nil
	}

	return "", fmt.Errorf("failed to retrieve secret value for key path '%s': %s", secretKeyWithPath, err)
}

// UpdateSecret updates a secret with given parameters.
//
// https://infisical.com/docs/api-reference/endpoints/secrets/update
func (c *Client) UpdateSecret(secretKey, workspaceID, environment, secretValue string, params ParamsUpdateSecret) (err error) {
	if params == nil {
		params = NewParamsUpdateSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment
	params["secretValue"] = secretValue

	var path string
	if c.e2eeEnabled {
		path = "/v3/secrets/%s"
	} else {
		path = "/v3/secrets/raw/%s"
	}

	if c.e2eeEnabled {
		var projectKey []byte
		if projectKey, err = c.projectKey(); err != nil {
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
	req, err = c.newRequestWithJSONBody("PATCH", fmt.Sprintf(path, secretKey), AuthMethodPreferToken, params)
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
func (c *Client) DeleteSecret(secretKey, workspaceID, environment string, params ParamsDeleteSecret) (err error) {
	if params == nil {
		params = NewParamsDeleteSecret()
	}

	// set essential params
	params["workspaceId"] = workspaceID
	params["environment"] = environment

	var path string
	if c.e2eeEnabled {
		path = "/v3/secrets/%s"
	} else {
		path = "/v3/secrets/raw/%s"
	}

	var req *http.Request
	req, err = c.newRequestWithJSONBody("DELETE", fmt.Sprintf(path, secretKey), AuthMethodPreferToken, params)
	if err == nil {
		c.dumpRequest(req)

		var res *http.Response
		if res, err = c.httpClient.Do(req); err == nil {
			c.dumpResponse(res)

			var body []byte
			if res.StatusCode == 200 {
				return nil
			} else {
				if body, err = io.ReadAll(res.Body); err == nil {
					err = fmt.Errorf("HTTP %d error: `%s`", res.StatusCode, string(body))
				} else {
					err = fmt.Errorf("HTTP %d error", res.StatusCode)
				}
			}
		}
	}

	return err
}
