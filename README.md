# infisical-go

[Infisical](https://infisical.com/) client library for golang.

## Usage

### Sample Code

```go
package main

import (
	"log"

	"github.com/meinside/infisical-go"
)

// NOTE: put yours here
const (
	// authentication
	apiKey       = "ak.1234567890.abcdefghijk"
	clientID     = "abcdefgh-0987-6543-xyzw-0123abcd4567"
	clientSecret = "abcdefghijklmnopqrstuvwxyz0123456789"

	workspaceID = "012345abcdefg"
	environment = "dev"
	keyPath     = "/folder1/folder2"

	//verbose = true // => for dumping HTTP requests & responses
	verbose = false
)

func main() {
	// create a client,
	client := infisical.NewClient(apiKey, clientID, clientSecret)
	//client.SetAPIBaseURL("https://app.infisical.com") // change API base URL (eg. for self-hosted infisical servers)
	client.Verbose = verbose

	// fetch all secrets at a path,
	if res, err := client.ListSecrets(infisical.NewParamsListSecrets().
		SetWorkspaceID(workspaceID).
		SetEnvironment(environment).
		SetSecretPath(keyPath),
	); err == nil {
		log.Printf("retrieved %d secret(s) at path '%s'", len(res.Secrets), keyPath)

		for _, secret := range res.Secrets {
			// fetch a value directly with path + key
			key := keyPath + "/" + secret.SecretKey

			if value, err := client.RetrieveSecretValue(secret.Workspace, secret.Environment, secret.Type, key); err == nil {
				log.Printf("retrieved value for secret keypath '%s' = '%s'", key, value)
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

### Helper Functions

Use `helper.Value()` for retrieving values:

```go
package main

import (
	"log"

	"github.com/meinside/infisical-go"
	"github.com/meinside/infisical-go/helper"
)

// NOTE: put yours here
const (
	clientID     = "abcdefgh-0987-6543-xyzw-0123abcd4567"
	clientSecret = "abcdefghijklmnopqrstuvwxyz0123456789"

	workspaceID   = "012345abcdefg"
	environment   = "dev"
	secretType    = infisical.SecretTypeShared
	secretKeyPath = "/folder1/folder2/KEY_A"
)

func main() {
	value, err := helper.Value(clientID, clientSecret, workspaceID, environment, secretType, secretKeyPath)
	if err != nil {
		panic(err)
	}

	log.Printf("retrieved value for key: %s = %s", secretKeyPath, value)
}
```

## Implemented APIs

* Users (./users.go)
- [ ] ~~[Get My User](https://infisical.com/docs/api-reference/endpoints/users/me)~~ // DEPRECATED
- [X] ~~[Get My Organization](https://infisical.com/docs/api-reference/endpoints/users/my-organizations)~~ // DEPRECATED

* Identities (./identities.go)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/identities/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/identities/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/identities/delete)

* Universal Auth (./universal_auth.go)
- [X] [Login](https://infisical.com/docs/api-reference/endpoints/universal-auth/login)
- [ ] [Attach](https://infisical.com/docs/api-reference/endpoints/universal-auth/attach)
- [ ] [Retrieve](https://infisical.com/docs/api-reference/endpoints/universal-auth/retrieve)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/universal-auth/update)
- [ ] [Create Client Secret](https://infisical.com/docs/api-reference/endpoints/universal-auth/create-client-secret)
- [ ] [List Client Secrets](https://infisical.com/docs/api-reference/endpoints/universal-auth/list-client-secrets)
- [ ] [Revoke Client Secret](https://infisical.com/docs/api-reference/endpoints/universal-auth/revoke-client-secret)
- [X] [Renew Access Token](https://infisical.com/docs/api-reference/endpoints/universal-auth/renew-access-token)

* Organizations (./organizations.go)
- [ ] [Get User Memberships](https://infisical.com/docs/api-reference/endpoints/organizations/memberships)
- [ ] [Update User Membership](https://infisical.com/docs/api-reference/endpoints/organizations/update-membership)
- [ ] [Delete Membership](https://infisical.com/docs/api-reference/endpoints/organizations/delete-membership)
- [ ] [List Identity Memberships](https://infisical.com/docs/api-reference/endpoints/organizations/list-identity-memberships)
- [X] [Get Projects](https://infisical.com/docs/api-reference/endpoints/organizations/workspaces)

* Projects (./projects.go)
- [ ] [Get User Memberships](https://infisical.com/docs/api-reference/endpoints/workspaces/memberships)
- [ ] [Update User Membership](https://infisical.com/docs/api-reference/endpoints/workspaces/update-membership)
- [ ] [Delete User Membership](https://infisical.com/docs/api-reference/endpoints/workspaces/delete-membership)
- [ ] [List Identity Memberships](https://infisical.com/docs/api-reference/endpoints/workspaces/list-identity-memberships)
- [ ] [Update Identity Membership](https://infisical.com/docs/api-reference/endpoints/workspaces/update-identity-membership)
- [ ] [Delete Identity Membership](https://infisical.com/docs/api-reference/endpoints/workspaces/delete-identity-membership)
- [ ] ~~[Get Key](https://infisical.com/docs/api-reference/endpoints/workspaces/workspace-key)~~ // DEPRECATED
- [ ] [Get Snapshots](https://infisical.com/docs/api-reference/endpoints/workspaces/secret-snapshots)
- [ ] [Roll Back to Snapshot](https://infisical.com/docs/api-reference/endpoints/workspaces/rollback-snapshot)

* Environments (./environments.go)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/environments/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/environments/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/environments/delete)

* Folders (./folders.go)
- [X] [List](https://infisical.com/docs/api-reference/endpoints/folders/list)
- [X] [Create](https://infisical.com/docs/api-reference/endpoints/folders/create)
- [X] [Update](https://infisical.com/docs/api-reference/endpoints/folders/update)
- [X] [Delete](https://infisical.com/docs/api-reference/endpoints/folders/delete)

* Secrets (./secrets.go)
- [X] [List](https://infisical.com/docs/api-reference/endpoints/secrets/list)
- [X] [Create](https://infisical.com/docs/api-reference/endpoints/secrets/create)
- [X] [Retrieve](https://infisical.com/docs/api-reference/endpoints/secrets/read)
- [X] [Update](https://infisical.com/docs/api-reference/endpoints/secrets/update)
- [X] [Delete](https://infisical.com/docs/api-reference/endpoints/secrets/delete)

* Secret imports (./secret_imports.go)
- [ ] [List](https://infisical.com/docs/api-reference/endpoints/secret-imports/list)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/secret-imports/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/secret-imports/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/secret-imports/delete)

* Audit Logs (./audit_logs.go)
- [ ] [Export](https://infisical.com/docs/api-reference/endpoints/audit-logs/export-audit-log)

## Error Codes

There is no detailed description in error responses from API (for now),

so it's sometimes quite hard to find out what is going wrong.

In my case, the reasons for common HTTP errors were:

* HTTP `400`: there were some missing parameters, or some of them were wrong/misformatted.
* HTTP `401`: was trying to access something with expired or wrong API key and/or token.
* HTTP `403`: was trying to access things that were not accessible with current API key and/or token.
* HTTP `404`: was trying to access something that doesn't exist; wrong key-path or etc.

## Test

With some environment variables:

```bash
export INFISICAL_API_KEY=ak.1234567890.abcdefghijk
export INFISICAL_WORKSPACE_ID=01234567-abcd-efgh-0987-ijklmnopqrst
export INFISICAL_CLIENT_ID=abcdefgh-0987-6543-xyzw-0123abcd4567
export INFISICAL_CLIENT_SECRET=abcdefghijklmnopqrstuvwxyz0123456789
export INFISICAL_ENVIRONMENT=dev
#export VERBOSE=true
```

run test:

```bash
$ go test
```

## CLI

I built a [CLI](https://github.com/meinside/infisical-go/tree/master/cmd/infisicli) for testing and personal use.

## Known Issues / Todos

### E2EE

E2EE features were removed due to [the deprecation of related endpoints](https://infisical.com/docs/api-reference/endpoints/service-tokens/get).

So projects with E2EE setting enabled may not work.

Version [v0.2.0](https://github.com/meinside/infisical-go/releases/tag/v0.2.0) will be the last version with E2EE support.

