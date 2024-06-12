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

* (DEPRECATED) Users (./users.go)
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
- [ ] [Revoke Access Token](https://infisical.com/docs/api-reference/endpoints/universal-auth/revoke-access-token)

* Organizations (./organizations.go)
- [ ] [Get User Memberships](https://infisical.com/docs/api-reference/endpoints/organizations/memberships)
- [ ] [Update User Membership](https://infisical.com/docs/api-reference/endpoints/organizations/update-membership)
- [ ] [Delete User Membership](https://infisical.com/docs/api-reference/endpoints/organizations/delete-membership)
- [ ] [List Identity Memberships](https://infisical.com/docs/api-reference/endpoints/organizations/list-identity-memberships)
- [X] [Get Projects](https://infisical.com/docs/api-reference/endpoints/organizations/workspaces)

* Projects (./projects.go)
- [ ] [Create Project](https://infisical.com/docs/api-reference/endpoints/workspaces/create-workspace)
- [ ] [Delete Project](https://infisical.com/docs/api-reference/endpoints/workspaces/delete-workspace)
- [ ] [Get Project](https://infisical.com/docs/api-reference/endpoints/workspaces/get-workspace)
- [ ] [Update Project](https://infisical.com/docs/api-reference/endpoints/workspaces/update-workspace)
- [ ] [Get Snapshots](https://infisical.com/docs/api-reference/endpoints/workspaces/secret-snapshots)
- [ ] [Roll Back to Snapshot](https://infisical.com/docs/api-reference/endpoints/workspaces/rollback-snapshot)

* Project Users (./project_users.go)
- [ ] [Invite Member](https://infisical.com/docs/api-reference/endpoints/project-users/invite-member-to-workspace)
- [ ] [Remove Member](https://infisical.com/docs/api-reference/endpoints/project-users/remove-member-from-workspace)
- [ ] [Get User Memberships](https://infisical.com/docs/api-reference/endpoints/project-users/memberships)
- [ ] [Get By Username](https://infisical.com/docs/api-reference/endpoints/project-users/get-by-username)
- [ ] [Update User Membership](https://infisical.com/docs/api-reference/endpoints/project-users/update-membership)

* Project Identities (./project_identities.go)
- [ ] [Create Identity Membership](https://infisical.com/docs/api-reference/endpoints/project-identities/add-identity-membership)
- [ ] [List Identity Memberships](https://infisical.com/docs/api-reference/endpoints/project-identities/list-identity-memberships)
- [ ] [Get Identity By ID](https://infisical.com/docs/api-reference/endpoints/project-identities/get-by-id)
- [ ] [Update Identity Membership](https://infisical.com/docs/api-reference/endpoints/project-identities/update-identity-membership)
- [ ] [Delete Identity Membership](https://infisical.com/docs/api-reference/endpoints/project-identities/delete-identity-membership)

* Project Roles (./project_roles.go)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/project-roles/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/project-roles/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/project-roles/delete)
- [ ] [Get By Slug](https://infisical.com/docs/api-reference/endpoints/project-roles/get-by-slug)
- [ ] [List](https://infisical.com/docs/api-reference/endpoints/project-roles/list)

* Environments (./environments.go)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/environments/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/environments/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/environments/delete)

* Folders (./folders.go)
- [X] [List](https://infisical.com/docs/api-reference/endpoints/folders/list)
- [X] [Create](https://infisical.com/docs/api-reference/endpoints/folders/create)
- [X] [Update](https://infisical.com/docs/api-reference/endpoints/folders/update)
- [X] [Delete](https://infisical.com/docs/api-reference/endpoints/folders/delete)

* Secret Tags (./secret_tags.go)
- [ ] [List](https://infisical.com/docs/api-reference/endpoints/secret-tags/list)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/secret-tags/create)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/secret-tags/delete)

* Secrets (./secrets.go)
- [X] [List](https://infisical.com/docs/api-reference/endpoints/secrets/list)
- [X] [Create](https://infisical.com/docs/api-reference/endpoints/secrets/create)
- [X] [Retrieve](https://infisical.com/docs/api-reference/endpoints/secrets/read)
- [X] [Update](https://infisical.com/docs/api-reference/endpoints/secrets/update)
- [X] [Delete](https://infisical.com/docs/api-reference/endpoints/secrets/delete)
- [ ] [Bulk Create](https://infisical.com/docs/api-reference/endpoints/secrets/create-many)
- [ ] [Bulk Update](https://infisical.com/docs/api-reference/endpoints/secrets/update-many)
- [ ] [Bulk Delete](https://infisical.com/docs/api-reference/endpoints/secrets/delete-many)
- [ ] [Attach Tags](https://infisical.com/docs/api-reference/endpoints/secrets/attach-tags)
- [ ] [Detach Tags](https://infisical.com/docs/api-reference/endpoints/secrets/detach-tags)

* Secret imports (./secret_imports.go)
- [ ] [List](https://infisical.com/docs/api-reference/endpoints/secret-imports/list)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/secret-imports/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/secret-imports/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/secret-imports/delete)

* Identity Specific Privilege (./identity_privileges.go)
- [ ] [Create Permanent](https://infisical.com/docs/api-reference/endpoints/identity-specific-privilege/create-permanent)
- [ ] [Create Temporary](https://infisical.com/docs/api-reference/endpoints/identity-specific-privilege/create-temporary)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/identity-specific-privilege/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/identity-specific-privilege/delete)
- [ ] [Find By Privilege Slug](https://infisical.com/docs/api-reference/endpoints/identity-specific-privilege/find-by-slug)
- [ ] [List](https://infisical.com/docs/api-reference/endpoints/identity-specific-privilege/list)

* Integrations (./integrations.go)
- [ ] [Create Auth](https://infisical.com/docs/api-reference/endpoints/integrations/create-auth)
- [ ] [List Auth](https://infisical.com/docs/api-reference/endpoints/integrations/list-auth)
- [ ] [Get Auth By ID](https://infisical.com/docs/api-reference/endpoints/integrations/find-auth)
- [ ] [Delete Auth](https://infisical.com/docs/api-reference/endpoints/integrations/delete-auth)
- [ ] [Delete Auth By ID](https://infisical.com/docs/api-reference/endpoints/integrations/delete-auth-by-id)
- [ ] [Create](https://infisical.com/docs/api-reference/endpoints/integrations/create)
- [ ] [Update](https://infisical.com/docs/api-reference/endpoints/integrations/update)
- [ ] [Delete](https://infisical.com/docs/api-reference/endpoints/integrations/delete)
- [ ] [List Project Integrations](https://infisical.com/docs/api-reference/endpoints/integrations/list-project-integrations)

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

### Organization ID

Due to the deprecation of [Get My Organizations](https://infisical.com/docs/api-reference/endpoints/users/my-organizations) API,

CLI will not iterate over all organizations automatically; you need to provide your desired origanization ID manually.

NOTE: organization ID can be retrieved from the Infisical console URL(eg. `https://app.infisical.com/org/<your-organization-id-here>/overview`).

CLI version [v0.3.1](https://github.com/meinside/infisical-go/releases/tag/v0.3.1) will be the last version which iterates over all organizations automatically.

