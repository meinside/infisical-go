package infisical

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// constants for E2EE
const tagSize = 16 // default(?) gcmTagSize = 16

// encrypt function for E2EE
//
// `aesKey` should be 32 bytes long (AES-256-GCM)
func encrypt(aesKey, text []byte) (encrypted, nonce, authTag []byte, err error) {
	if len(aesKey) != 32 {
		return nil, nil, nil, fmt.Errorf("the length of AES key is not 32 but %d, cannot encrypt with given things", len(aesKey))
	}

	var block cipher.Block
	if block, err = aes.NewCipher(aesKey); err != nil {
		return nil, nil, nil, err
	}

	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCM(block); err != nil {
		return nil, nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, err
	}

	encryptedWhole := gcm.Seal(nil, nonce, text, nil)

	encrypted = encryptedWhole[:len(encryptedWhole)-tagSize]
	authTag = encryptedWhole[len(encryptedWhole)-tagSize:]

	return encrypted, nonce, authTag, nil
}

// decrypt function for E2EE
//
// `aesKey` should be 32 bytes long (AES-256-GCM)
func decrypt(aesKey, encrypted, nonce, authTag []byte) (decrypted []byte, err error) {
	if len(aesKey) != 32 {
		return nil, fmt.Errorf("the length of AES key is not 32 but %d, cannot decrypt with given things", len(aesKey))
	}

	var block cipher.Block
	if block, err = aes.NewCipher(aesKey); err != nil {
		return nil, err
	}

	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCMWithNonceSize(block, len(nonce)); err != nil {
		return nil, err
	}

	encrypted = append(encrypted, authTag...)
	if decrypted, err = gcm.Open(nil, nonce, encrypted, nil); err != nil {
		return nil, err
	}

	return decrypted, nil
}

// encode base64 string from bytes array
func encodeBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// decode base64 string to bytes array
func decodeBase64(encoded string) (decoded []byte, err error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// get service token and decrypt project key from it
func (c *Client) projectKey(token WorkspaceToken) (projectKey []byte, err error) {
	var serviceToken ServiceToken
	if serviceToken, err = c.RetrieveServiceToken(token); err != nil {
		return nil, err
	}

	if serviceToken.EncryptedKey == nil {
		return nil, fmt.Errorf("returned service token's encrypted key is nil")
	}
	if serviceToken.IV == nil {
		return nil, fmt.Errorf("returned service token's IV is nil")
	}
	if serviceToken.Tag == nil {
		return nil, fmt.Errorf("returned service token's tag is nil")
	}

	if token.Token == "" {
		return nil, fmt.Errorf("`token` is missing, cannot decrypt project key")
	}

	splittedToken := strings.Split(token.Token, ".")
	serviceTokenSecret := splittedToken[len(splittedToken)-1]

	// decrypt things
	var decodedEncryptedKey, decodedIV, decodedTag []byte
	if decodedEncryptedKey, err = decodeBase64(*serviceToken.EncryptedKey); err != nil {
		return nil, err
	}
	if decodedIV, err = decodeBase64(*serviceToken.IV); err != nil {
		return nil, err
	}
	if decodedTag, err = decodeBase64(*serviceToken.Tag); err != nil {
		return nil, err
	}
	if projectKey, err = decrypt([]byte(serviceTokenSecret), decodedEncryptedKey, decodedIV, decodedTag); err != nil {
		return nil, err
	}

	return projectKey, nil
}

// decrypt encrypted secrets (when E2EE is enabled)
func (c *Client) decryptSecrets(token WorkspaceToken, secrets []Secret) (decrypted []Secret, err error) {
	var projectKey []byte
	if projectKey, err = c.projectKey(token); err != nil {
		return nil, err
	}

	decrypted = []Secret{}
	var encrypted, iv, tag []byte
	var key, value, comment []byte
	for _, secret := range secrets {
		// (key)
		if encrypted, err = decodeBase64(secret.SecretKeyCiphertext); err != nil {
			return nil, err
		}
		if iv, err = decodeBase64(secret.SecretKeyIV); err != nil {
			return nil, err
		}
		if tag, err = decodeBase64(secret.SecretKeyTag); err != nil {
			return nil, err
		}
		if key, err = decrypt(projectKey, encrypted, iv, tag); err != nil {
			return nil, err
		}
		secret.SecretKey = string(key)

		// (value)
		if encrypted, err = decodeBase64(secret.SecretValueCiphertext); err != nil {
			return nil, err
		}
		if iv, err = decodeBase64(secret.SecretValueIV); err != nil {
			return nil, err
		}
		if tag, err = decodeBase64(secret.SecretValueTag); err != nil {
			return nil, err
		}
		if value, err = decrypt(projectKey, encrypted, iv, tag); err != nil {
			return nil, err
		}
		secret.SecretValue = string(value)

		// (comment)
		if secret.SecretCommentCiphertext != "" && secret.SecretCommentIV != "" && secret.SecretCommentTag != "" {
			if encrypted, err = decodeBase64(secret.SecretCommentCiphertext); err != nil {
				return nil, err
			}
			if iv, err = decodeBase64(secret.SecretCommentIV); err != nil {
				return nil, err
			}
			if tag, err = decodeBase64(secret.SecretCommentTag); err != nil {
				return nil, err
			}
			if comment, err = decrypt(projectKey, encrypted, iv, tag); err != nil {
				return nil, err
			}
			secret.SecretComment = string(comment)
		}

		decrypted = append(decrypted, secret)
	}

	return decrypted, nil
}

// decrypt an encrypted secret (when E2EE is enabled)
func (c *Client) decryptSecret(token WorkspaceToken, secret Secret) (decrypted Secret, err error) {
	var projectKey []byte
	if projectKey, err = c.projectKey(token); err != nil {
		return Secret{}, err
	}

	// decrypt key, value, and comment
	var text, iv, tag []byte
	var key, value, comment []byte

	// (key)
	if text, err = decodeBase64(secret.SecretKeyCiphertext); err != nil {
		return Secret{}, err
	}
	if iv, err = decodeBase64(secret.SecretKeyIV); err != nil {
		return Secret{}, err
	}
	if tag, err = decodeBase64(secret.SecretKeyTag); err != nil {
		return Secret{}, err
	}
	if key, err = decrypt(projectKey, text, iv, tag); err != nil {
		return Secret{}, err
	}
	secret.SecretKey = string(key)

	// (value)
	if text, err = decodeBase64(secret.SecretValueCiphertext); err != nil {
		return Secret{}, err
	}
	if iv, err = decodeBase64(secret.SecretValueIV); err != nil {
		return Secret{}, err
	}
	if tag, err = decodeBase64(secret.SecretValueTag); err != nil {
		return Secret{}, err
	}
	if value, err = decrypt(projectKey, text, iv, tag); err != nil {
		return Secret{}, err
	}
	secret.SecretValue = string(value)

	// (comment)
	if secret.SecretCommentCiphertext != "" && secret.SecretCommentIV != "" && secret.SecretCommentTag != "" {
		if text, err = decodeBase64(secret.SecretCommentCiphertext); err != nil {
			return Secret{}, err
		}
		if iv, err = decodeBase64(secret.SecretCommentIV); err != nil {
			return Secret{}, err
		}
		if tag, err = decodeBase64(secret.SecretCommentTag); err != nil {
			return Secret{}, err
		}
		if comment, err = decrypt(projectKey, text, iv, tag); err != nil {
			return Secret{}, err
		}
		secret.SecretComment = string(comment)
	}

	return secret, nil
}
