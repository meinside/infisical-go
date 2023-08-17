# infisicli

Infisicli, a simple [Infisical](https://infisical.com/) CLI.

## Features

- [X] List and view organizations, workspaces, and secrets.
- [ ] Create/Update/Delete organizations, workspaces, and secrets.

## Install

```bash
$ go install github.com/meinside/infisical-go/cmd/infisicli@latest
```

## Configuration

Put a `config.json` file in `$XDG_CONFIG_HOME/infisicli/` directory with following content:

```json
{
  "api_key": "ak.1234567890.abcdefghijk",
  "token": "st.xyzwabcd.0987654321.abcdefghijklmnop",
  "e2ee_enabled": true
}
```
where `api_key` is your API key of Infisical account,

`token` is your token of workspace,

and `e2ee_enabled` is whether you enabled your E2EE setting or not.

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
org: your-org-name   | id: <your-org-id> (customer id: your-customer-id)
```

### List Workspaces

List workspaces with <your-org-id> obtained from above:

```bash
$ infisicli -lw -o=<your-org-id>
workspace: <workspace1-id>     | name: workspace1-name
 - env: dev     | name: Development, id: workspace1-dev-env-id
 - env: staging | name: Staging, id: workspace1-staging-env-id
 - env: prod    | name: Production, id: workspace1-prod-env-id
workspace: <workspace2-id> | name: workspace2-name
 - env: dev     | name: Development, id: workspace2-dev-env-id
 - env: staging | name: Staging, id: workspace2-staging-env-id
 - env: prod    | name: Production, id: workspace2-prod-env-id
...
```

### List Secrets

Now list secrets at a folder:

```bash
$ infisicli -las -w=<workspace1-id> -e=dev -f=/folder1/folder2
workspace: <workspace1-id>     | env: dev      | type: <type1>  | <key1> = <value2>
workspace: <workspace1-id>     | env: dev      | type: <type2>  | <key2> = <value2>
```

### Print secret value

Following will print the value of given key-path (folder + key) without a trailing newline:

```bash
$ infisicli -p -w=<workspace1-id> -e=dev -t=<type2> -k=/folder1/folder2/<value2>
<value2>
```

It can be used in shell scripts like:

```bash
VALUE=$(infisicli -p -w=<workspace1-id> -e=dev -t=<type1> -k=/folder1/folder2/<value1>)

echo "value for key: <value1> = $VALUE"
```

## License

MIT

