# infisicli

Infisicli, a simple [Infisical](https://infisical.com/) CLI.

## Features

- [X] List organizations, workspaces, environments, and secrets.
- [X] Create/Update/Delete secrets.
- [ ] Create/Update/Delete organizations, workspaces, and environments.

## Install

```bash
$ go install github.com/meinside/infisical-go/cmd/infisicli@latest
```

## Configuration

Put a `config.json` file in `$XDG_CONFIG_HOME/infisicli/` directory with following content:

```json
{
  "api_key": "ak.1234567890.abcdefghijk",
  "workspaces": {
    "012345abcdefg": {
      "token": "st.xyzwabcd.0987654321.abcdefghijklmnop",
      "e2ee": true
    }
  }
}
```

where `api_key` is the API key of your Infisical account,

keys in `workspaces` are organization IDs,

`token` is the token of each organization,

and `e2ee` is whether the workspace's E2EE setting is enabled or not.


**NOTE**: When your E2EE setting is off, you can omit the `api_key` value:

```json
{
  "workspaces": {
    "012345abcdefg": {
      "token": "st.xyzwabcd.0987654321.abcdefghijklmnop"
    }
  }
}
```

but in this case some features that require `api_key` (eg. listing organizations, listing workspaces, …) will not function.

## Usage

You can see detailed help messages with:

```bash
$ infisicli -h
# or
$ infisicli --help
```

With the valid configuration, you can do following tasks:

### List Organizations

List organizations info with:

```bash
$ infisicli -lo

           id | name
----
<your-org-id> | your-org-name
```

### List Workspaces

List workspaces with <your-org-id> obtained from above:

```bash
$ infisicli -lw -o=<your-org-id>

   workspace id | name
----
<workspace1-id> | workspace1-name
      dev | Development (workspace1-dev-env-id)
  staging |     Staging (workspace1-staging-env-id)
     prod |  Production (workspace1-prod-env-id)
----
<workspace2-id> | workspace2-name
      dev | Development (workspace2-dev-env-id)
  staging |     Staging (workspace2-staging-env-id)
     prod |  Production (workspace2-prod-env-id)
...
```

### List Secrets

Now list secrets at a folder:

```bash
$ infisicli -las -w=<workspace1-id> -e=dev -f=/folder1/folder2

      workspace | env |    type | path/key=value
----
<workspace1-id> | dev | <type1> | /folder1/folder2/<key1>=<value1>
<workspace1-id> | dev | <type2> | /folder1/folder2/<key2>=<value2>
...
```

### Print a Secret Value

Following will print the value of given key-path (folder + key) without a trailing newline:

```bash
$ infisicli -p -w=<workspace1-id> -e=dev -t=<type1> -k=/folder1/folder2/<key1>

<value1>
```

It can also be used in shell scripts like:

```bash
VALUE=$(infisicli -p -w=<workspace1-id> -e=dev -t=<type1> -k=/folder1/folder2/<key1>)

echo "value for key: <key1> = $VALUE"
```

### Create/Update/Delete a Secret

Create a new secret:

```bash
$ infisicli -n -w=<workspace1-id> -e=dev -t=shared -k=<path/key> -s=<new-value>

> Successfully created a new secret value at: <path/key>
```

Update a secret value:

```bash
$ infisicli -u -w=<workspace1-id> -e=dev -t=shared -k=<path/key> -s=<updated-value>

> Successfully updated a secret value at: <path/key>
```

Delete a secret:

```bash
$ infisicli -d -w=<workspace1-id> -e=dev -t=shared -k=<path/key>

> Successfully deleted a secret value at: <path/key>
```

## License

MIT

