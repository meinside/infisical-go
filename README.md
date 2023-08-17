# infisical-go

[Infisical](https://infisical.com/) client library for golang.

## Usage

```go
package main

import (
	"log"

	"github.com/meinside/infisical-go"
)

// NOTE: put yours here
const (
	// authentication
	apiKey      = "ak.1234567890.abcdefghijk"
	token       = "st.xyzwabcd.0987654321.abcdefghijklmnop"

	// workspace & environment
	workspaceID = "012345abcdefg"
	environment = "dev"

	keyPath    = "/folder1/folder2"
	secretType = infisical.SecretTypeShared

	//verbose = true // => for dumping HTTP requests & responses
	verbose = false
)

func main() {
	// create a client, (NOTE: enable E2EE setting in your project)
	client := infisical.NewE2EEEnabledClient(apiKey, token)
	client.Verbose = verbose

	// fetch all secrets at a path,
	if secrets, err := client.RetrieveSecretsAtPath(keyPath, workspaceID, environment); err == nil {
		log.Printf("retrieved %d secret(s) at path '%s'", len(secrets), keyPath)

		for _, secret := range secrets {
			// fetch a value directly with path + key
			key := keyPath + "/" + secret.SecretKey

			if value, err := client.RetrieveSecretValue(key, workspaceID, environment, secret.Type); err == nil {
				log.Printf("retrieved value for secret key '%s' = '%s'", key, value)
			} else {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}
```

Output:

```
2023/08/16 14:30:33 retrieved 2 secret(s) at path '/folder1/folder2'
2023/08/16 14:30:34 retrieved value for secret key '/folder1/folder2/KEY_A' = 'value A'
2023/08/16 14:30:36 retrieved value for secret key '/folder1/folder2/KEY_B' = 'value B'
```

## Implemented

* Users
- [ ] [Get My User](https://infisical.com/docs/api-reference/endpoints/users/me)
- [X] [Get My Organizations](https://infisical.com/docs/api-reference/endpoints/users/my-organizations)

* Organizations
- [ ] [Get Memberships](https://infisical.com/docs/api-reference/endpoints/organizations/memberships)
- [ ] [Update Membership](https://infisical.com/docs/api-reference/endpoints/organizations/update-membership)
- [ ] [Delete Membership](https://infisical.com/docs/api-reference/endpoints/organizations/delete-membership)
- [X] [Get Projects](https://infisical.com/docs/api-reference/endpoints/organizations/workspaces)

* Projects
- [ ] [Get Memberships](https://infisical.com/docs/api-reference/endpoints/workspaces/memberships)
- [ ] [Update Membership](https://infisical.com/docs/api-reference/endpoints/workspaces/update-membership)
- [ ] [Delete Membership](https://infisical.com/docs/api-reference/endpoints/workspaces/delete-membership)
- [ ] [Get Key](https://infisical.com/docs/api-reference/endpoints/workspaces/workspace-key)
- [ ] [Get Logs](https://infisical.com/docs/api-reference/endpoints/workspaces/logs)
- [ ] [Get Snapshots](https://infisical.com/docs/api-reference/endpoints/workspaces/secret-snapshots)
- [ ] [Roll Back to Snapshot](https://infisical.com/docs/api-reference/endpoints/workspaces/rollback-snapshot)

* Secrets
- [X] [Retrieve All](https://infisical.com/docs/api-reference/endpoints/secrets/read)
- [X] [Create](https://infisical.com/docs/api-reference/endpoints/secrets/create)
- [X] [Retrieve](https://infisical.com/docs/api-reference/endpoints/secrets/read-one)
- [X] [Update](https://infisical.com/docs/api-reference/endpoints/secrets/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/secrets/delete) // FIXME: fails with 404 error
- [ ] [Get Versions](https://infisical.com/docs/api-reference/endpoints/secrets/versions)
- [ ] [Roll Back to Version](https://infisical.com/docs/api-reference/endpoints/secrets/rollback-version)

* Service Tokens
- [X] [Get](https://infisical.com/docs/api-reference/endpoints/service-tokens/get)

## Test

With some environment variables:

```bash
export INFISICAL_API_KEY=ak.1234567890.abcdefghijk
export INFISICAL_TOKEN=st.xyzwabcd.0987654321.abcdefghijklmnop
export INFISICAL_WORKSPACE_ID=012345abcdefg
#export INFISICAL_ENVIRONMENT=dev
#export INFISICAL_E2EE=enabled
#export VERBOSE=true
```

run test:

```bash
$ go test
```

## CLI

I built a [CLI](https://github.com/meinside/infisical-go/tree/master/cmd/infisicli) for testing and personal use.

## License

MIT

