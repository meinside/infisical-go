package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/meinside/infisical-go"
	"github.com/meinside/version-go"
	"github.com/tailscale/hujson"
)

const (
	applicationName = "infisicli"
	configFilename  = "config.json"
)

const (
	// arguments for commands
	cmdHelpShort           = "-h"
	cmdHelpLong            = "--help"
	cmdListWorkspacesShort = "-lw"
	cmdListWorkspacesLong  = "--list-workspaces"
	cmdListAllSecretsShort = "-las"
	cmdListAllSecretsLong  = "--list-all-secrets"
	cmdPrintValueShort     = "-p"
	cmdPrintValueLong      = "--print"
	cmdNewValueShort       = "-n"
	cmdNewValueLong        = "--new-value"
	cmdUpdateValueShort    = "-u"
	cmdUpdateValueLong     = "--update-value"
	cmdDeleteValueShort    = "-d"
	cmdDeleteValueLong     = "--delete-value"
	cmdVersionShort        = "-v"
	cmdVersionLong         = "--version"

	// arguments for command parameters
	argOrganizationShort = "-o"
	argOrganizationLong  = "--organization"
	argWorkspaceShort    = "-w"
	argWorkspaceLong     = "--workspace"
	argEnvironmentShort  = "-e"
	argEnvironmentLong   = "--environment"
	argTypeShort         = "-t"
	argTypeLong          = "--type"
	argKeyShort          = "-k"
	argKeyLong           = "--key"
	argSecretShort       = "-s"
	argSecretLong        = "--secret"
	argCommentShort      = "-c"
	argCommentLong       = "--comment"
	argFolderShort       = "-f"
	argFolderLong        = "--folder"

	// arguments for something else
	argVerboseShort = "-V"
	argVerboseLong  = "--verbose"

	// regex for parsing command param arguments (in 'key=value' format)
	regexKeyValue = `(.*?)=['"]?(.*?)['"]?$`
)

// config struct
type config struct {
	// Infisical Account's API Key
	APIKey *string `json:"api_key,omitempty"`

	// Universal Auth values
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// standardize given JSON (JWCC) bytes
func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()

	return ast.Pack(), nil
}

// load config file
func loadConfig() (conf config, err error) {
	// https://xdgbasedirectoryspecification.com
	configDir := os.Getenv("XDG_CONFIG_HOME")

	// If the value of the environment variable is unset, empty, or not an absolute path, use the default
	if configDir == "" || configDir[0:1] != "/" {
		var homeDir string
		if homeDir, err = os.UserHomeDir(); err == nil {
			configDir = filepath.Join(homeDir, ".config", applicationName)
		}
	} else {
		configDir = filepath.Join(configDir, applicationName)
	}

	if err == nil {
		configFilepath := filepath.Join(configDir, configFilename)

		var bytes []byte
		if bytes, err = os.ReadFile(configFilepath); err == nil {
			if bytes, err = standardizeJSON(bytes); err == nil {
				if err = json.Unmarshal(bytes, &conf); err == nil {
					return conf, nil
				}
			}
		}
	}

	return conf, fmt.Errorf("failed to load config: %s", err)
}

// check if given short/long argument is included
func hasArg(args []string, short, long string) bool {
	for _, arg := range args {
		if arg == short || arg == long {
			return true
		}
	}
	return false
}

// show version string
func showVersion() {
	fmt.Printf("%s\n", version.Minimum())

	os.Exit(0)
}

// show help message
func showHelp(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf(`
Usage:

  %[1]s %[2]s
  %[1]s %[3]s
  : Print this help message to stdout.

  %[1]s %[4]s
  %[1]s %[5]s
  : Print the version string to stdout.

  %[1]s %[6]s
  %[1]s %[7]s
  : List all your workspaces with given org id to stdout.
  > NOTE: Your org id can be retrieved from infisical console URL.
    eg. %[1]s %[6]s %[18]s=0a1b2c3d4e5f
    eg. %[1]s %[7]s %[19]s=0a1b2c3d4e5f

  %[1]s %[8]s
  %[1]s %[9]s
  : List all secret values to stdout.
    eg. %[1]s %[8]s %[22]s=012345abcdefg %[24]s=dev
    eg. %[1]s %[9]s %[23]s=012345abcdefg %[25]s=dev
    eg. %[1]s %[8]s %[22]s=012345abcdefg %[24]s=dev %[20]s=/folder1/folder2
    eg. %[1]s %[9]s %[23]s=012345abcdefg %[25]s=dev %[21]s=/folder1/folder2

  %[1]s %[10]s
  %[1]s %[11]s
  : Print the secret value to stdout without a trailing newline.
    eg. %[1]s %[10]s %[22]s=012345abcdefg %[24]s=dev %[26]s=shared %[28]s=/folder/SECRET_KEY_1
    eg. %[1]s %[11]s %[23]s=012345abcdefg %[25]s=dev %[27]s=shared %[29]s=/folder/SECRET_KEY_1

  %[1]s %[12]s
  %[1]s %[13]s
  : Create a new secret with given parameters.
    eg. %[1]s %[12]s %[22]s=012345abcdefg %[24]s=dev %[26]s=shared %[28]s=/folder/NEW_KEY %[30]s=NEW_VALUE
    eg. %[1]s %[13]s %[23]s=012345abcdefg %[25]s=dev %[27]s=shared %[29]s=/folder/NEW_KEY %[31]s=NEW_VALUE
    eg. %[1]s %[12]s %[22]s=012345abcdefg %[24]s=dev %[26]s=shared %[28]s=/folder/NEW_KEY %[30]s=NEW_VALUE %[32]s=COMMENT
    eg. %[1]s %[13]s %[23]s=012345abcdefg %[25]s=dev %[27]s=shared %[29]s=/folder/NEW_KEY %[31]s=NEW_VALUE %[33]s=COMMENT

  %[1]s %[14]s
  %[1]s %[15]s
  : Update the value of a secret at the given key-path.
    eg. %[1]s %[14]s %[22]s=012345abcdefg %[24]s=dev %[26]s=shared %[28]s=/folder/SOME_KEY %[30]s=UPDATED_VALUE
    eg. %[1]s %[15]s %[23]s=012345abcdefg %[25]s=dev %[27]s=shared %[29]s=/folder/SOME_KEY %[31]s=UPDATED_VALUE

  %[1]s %[16]s
  %[1]s %[17]s
  : Delete a secret at the given key-path.
    eg. %[1]s %[16]s %[22]s=012345abcdefg %[24]s=dev %[26]s=shared %[28]s=/folder/SOME_KEY
    eg. %[1]s %[17]s %[23]s=012345abcdefg %[25]s=dev %[27]s=shared %[29]s=/folder/SOME_KEY

Other optional arguments:

  %[34]s / %[35]s
  : Dump http requests/responses for debugging.
`,
		// executable name
		applicationName,

		// commands
		cmdHelpShort, cmdHelpLong,
		cmdVersionShort, cmdVersionLong,
		cmdListWorkspacesShort, cmdListWorkspacesLong,
		cmdListAllSecretsShort, cmdListAllSecretsLong,
		cmdPrintValueShort, cmdPrintValueLong,
		cmdNewValueShort, cmdNewValueLong,
		cmdUpdateValueShort, cmdUpdateValueLong,
		cmdDeleteValueShort, cmdDeleteValueLong,

		// parameters
		argOrganizationShort, argOrganizationLong,
		argFolderShort, argFolderLong,
		argWorkspaceShort, argWorkspaceLong,
		argEnvironmentShort, argEnvironmentLong,
		argTypeShort, argTypeLong,
		argKeyShort, argKeyLong,
		argSecretShort, argSecretLong,
		argCommentShort, argCommentLong,

		// others
		argVerboseShort, argVerboseLong,
	)

	if err != nil {
		os.Exit(1)
	}
}

// convert an array of strings like "key1=value1" into a string-string map
func convertKeyValueParams(params []string) (result map[string]string) {
	result = map[string]string{}

	regexKV := regexp.MustCompile(regexKeyValue)

	for _, param := range params {
		matches := regexKV.FindStringSubmatch(param)

		if len(matches) == 3 {
			k := matches[1]
			v := matches[2]

			result[k] = v
		}
	}

	return result
}

// search value for short/long key from given kvs
func valueFromKVs(short, long string, kvs map[string]string) (string, error) {
	if s, exists := kvs[short]; exists && len(s) > 0 {
		return s, nil
	}
	if l, exists := kvs[long]; exists && len(l) > 0 {
		return l, nil
	}

	return "", fmt.Errorf("no value for argument '%s' or '%s' was provided", short, long)
}

// run things with given arguments from main()
func run(args []string) {
	verbose := hasArg(args, argVerboseShort, argVerboseLong)

	var err error
	if hasArg(args, cmdVersionShort, cmdVersionLong) {
		// show version
		showVersion()
	} else if hasArg(args, cmdHelpShort, cmdHelpLong) {
		// do nothing
	} else {
		//handle commands here
		if hasArg(args, cmdListWorkspacesShort, cmdListWorkspacesLong) {
			err = doListWorkspaces(args, verbose)
		} else if hasArg(args, cmdListAllSecretsShort, cmdListAllSecretsLong) {
			err = doListAllSecrets(args, verbose)
		} else if hasArg(args, cmdPrintValueShort, cmdPrintValueLong) {
			err = doPrintValue(args, verbose)
		} else if hasArg(args, cmdNewValueShort, cmdNewValueLong) {
			err = doCreateValue(args, verbose)
		} else if hasArg(args, cmdUpdateValueShort, cmdUpdateValueLong) {
			err = doUpdateValue(args, verbose)
		} else if hasArg(args, cmdDeleteValueShort, cmdDeleteValueLong) {
			err = doDeleteValue(args, verbose)
		}
	}

	showHelp(err)
}

// get max length of given items
func maxLength[T any](items []T, lenFunc func(item T) int) (max int) {
	max = math.MinInt

	var current int
	for _, item := range items {
		current = lenFunc(item)
		if current > max {
			max = current
		}
	}

	return max
}

// do something with the client
func do(fn func(c *infisical.Client) error, verbose bool) error {
	cfg, err := loadConfig()

	if err == nil {
		var client *infisical.Client
		if cfg.APIKey != nil {
			client = infisical.NewClient(*cfg.APIKey, cfg.ClientID, cfg.ClientSecret)
		} else {
			client = infisical.NewClientWithoutAPIKey(cfg.ClientID, cfg.ClientSecret)
		}
		client.Verbose = verbose

		return fn(client)
	}

	return err
}

// list workspaces, will os.Exit(0) on success
func doListWorkspaces(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		orgs := []string{}

		var err error
		params := convertKeyValueParams(args)
		if org, _ := valueFromKVs(argOrganizationShort, argOrganizationLong, params); org != "" {
			orgs = append(orgs, org)
		} else {
			err = fmt.Errorf("organization id (`%s` or `%s`) is not provided", argOrganizationShort, argOrganizationLong)
		}

		if err == nil {
			allWorkspaces := []infisical.Workspace{}

			// fetch things
			for _, org := range orgs {
				var workspaces infisical.ProjectsData
				if workspaces, err = c.RetrieveProjects(org); err == nil {
					allWorkspaces = append(allWorkspaces, workspaces.Workspaces...)
				} else {
					break
				}
			}

			// print result
			if err == nil {
				if len(allWorkspaces) <= 0 {
					fmt.Printf("* There was no workspace for given parameters.\n")
				} else {
					// calculate max lengths for formatting
					maxLenOrgID := maxLength(allWorkspaces, func(workspace infisical.Workspace) int {
						return len(workspace.Organization)
					})
					maxLenWorkspaceID := maxLength(allWorkspaces, func(workspace infisical.Workspace) int {
						return len(workspace.ID)
					})
					maxLenWorkspaceName := maxLength(allWorkspaces, func(workspace infisical.Workspace) int {
						return len(workspace.Name)
					})
					workspaceFormat := fmt.Sprintf("%%%ds | %%%ds | %%-%ds\n", maxLenOrgID, maxLenWorkspaceID, maxLenWorkspaceName)

					// print headers
					fmt.Printf(workspaceFormat, "org id", "workspace id", "workspace name")

					for _, workspace := range allWorkspaces {
						maxLenSlug := maxLength(workspace.Environments, func(env infisical.WorkspaceEnvironment) int {
							return len(env.Slug)
						})
						maxLenName := maxLength(workspace.Environments, func(env infisical.WorkspaceEnvironment) int {
							return len(env.Name)
						})
						envFormat := fmt.Sprintf("  %%%ds | %%%ds\n", maxLenSlug, maxLenName)

						// print workspace
						fmt.Printf("----\n")
						fmt.Printf(workspaceFormat, workspace.Organization, workspace.ID, workspace.Name)

						// print environments
						for _, env := range workspace.Environments {
							fmt.Printf(envFormat, env.Slug, env.Name)
						}
					}
				}

				os.Exit(0)
			}
		}

		return err
	}, verbose)
}

// fetch folder paths in `folderPath` recursively
func fetchFolderPaths(c *infisical.Client, workspace, environment string, folderPath *string) (folderPaths []string, err error) {
	var dir string
	if folderPath == nil {
		dir = "/"
	} else {
		dir = *folderPath
	}

	params := infisical.NewParamsListFolders().
		SetPath(dir)

	var subdir string
	var result infisical.FoldersData
	if result, err = c.ListFolders(workspace, environment, params); err == nil {
		for _, folder := range result.Folders {
			subdir = path.Join(dir, folder.Name)

			folderPaths = append(folderPaths, subdir)
		}

		// recurse subfolder paths
		for _, folder := range result.Folders {
			subdir = path.Join(dir, folder.Name)

			var subFolderPaths []string
			if subFolderPaths, err = fetchFolderPaths(c, workspace, environment, &subdir); err == nil {
				folderPaths = append(folderPaths, subFolderPaths...)
			} else {
				break
			}
		}
	}

	return folderPaths, err
}

// print secrets to stdout
func printSecrets(all []infisical.Secret, secrets []infisical.Secret, imports []infisical.SecretImport, foldersMap map[string]string) {
	maxLenWorkspace := maxLength(all, func(secret infisical.Secret) int {
		return len(secret.Workspace)
	})
	maxLenEnv := maxLength(all, func(secret infisical.Secret) int {
		return len(secret.Environment)
	})
	maxLenType := maxLength(all, func(secret infisical.Secret) int {
		return len(secret.Type)
	})
	format := fmt.Sprintf("%%%ds | %%%ds | %%%ds | %%s\n", maxLenWorkspace, maxLenEnv, maxLenType)

	// print headers
	fmt.Printf(format, "workspace", "env", "type", "path/key=value")
	fmt.Printf("----\n")

	var folder string
	var exists bool

	// print key-values
	for _, secret := range secrets {
		if folder, exists = foldersMap[secret.ID]; exists {
			fmt.Printf(format,
				secret.Workspace,
				secret.Environment,
				secret.Type,
				path.Join(folder, secret.SecretKey)+"="+secret.SecretValue,
			)
		}
	}
	if len(imports) > 0 {
		fmt.Printf("<imported>\n")

		for _, imp := range imports {
			for _, secret := range imp.Secrets {
				if folder, exists = foldersMap[secret.ID]; exists {
					fmt.Printf(format,
						secret.Workspace,
						secret.Environment,
						secret.Type,
						path.Join(folder, secret.SecretKey)+"="+secret.SecretValue,
					)
				}
			}
		}
	}
}

// list all secrets, will os.Exit(0) on success
func doListAllSecrets(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		var err error
		params := convertKeyValueParams(args)

		var workspace, environment string
		if workspace, err = valueFromKVs(argWorkspaceShort, argWorkspaceLong, params); err != nil {
			return err
		}
		if environment, err = valueFromKVs(argEnvironmentShort, argEnvironmentLong, params); err != nil {
			return err
		}

		var result infisical.SecretsData
		listParams := infisical.NewParamsListSecrets().
			SetIncludeImports(true).
			SetWorkspaceID(workspace).
			SetEnvironment(environment)

		// folder
		folder, _ := valueFromKVs(argFolderShort, argFolderLong, params)
		var allFolderPaths []string
		if folder != "" {
			allFolderPaths = []string{folder}
		} else { // if folder is not given, iterate all folders
			folderPaths, _ := fetchFolderPaths(c, workspace, environment, nil)

			allFolderPaths = append([]string{"/"}, folderPaths...)
		}

		// and fetch all secrets from them
		all := []infisical.Secret{}
		news := []infisical.Secret{}
		secrets := []infisical.Secret{}
		imports := []infisical.SecretImport{}
		foldersMap := map[string]string{}
		for _, folderPath := range allFolderPaths {
			listParams = listParams.
				SetSecretPath(folderPath)

			// fetch secrets
			if result, err = c.ListSecrets(listParams); err == nil {
				news = result.Secrets
				for _, imp := range result.Imports {
					news = append(news, imp.Secrets...)
				}

				secrets = append(secrets, result.Secrets...)
				imports = append(imports, result.Imports...)
				all = append(all, news...)

				for _, secret := range news {
					foldersMap[secret.ID] = folderPath
				}
			}
		}

		if len(all) > 0 {
			printSecrets(all, secrets, imports, foldersMap)
		} else {
			fmt.Printf("* There was no secret for given parameters.\n")
		}

		if err == nil {
			os.Exit(0)
		}

		return err
	}, verbose)
}

// print a secret value, will os.Exit(0) on success
func doPrintValue(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		var err error
		params := convertKeyValueParams(args)

		var key, workspace, environment, typ string
		if key, err = valueFromKVs(argKeyShort, argKeyLong, params); err != nil {
			return err
		}
		if workspace, err = valueFromKVs(argWorkspaceShort, argWorkspaceLong, params); err != nil {
			return err
		}
		if environment, err = valueFromKVs(argEnvironmentShort, argEnvironmentLong, params); err != nil {
			return err
		}
		if typ, err = valueFromKVs(argTypeShort, argTypeLong, params); err != nil {
			return err
		}

		var result string
		result, err = c.RetrieveSecretValue(workspace, environment, infisical.SecretType(typ), key)
		if err == nil {
			fmt.Printf("%s", result)

			os.Exit(0)
		}

		return err
	}, verbose)
}

// create a new secret value, will os.Exit(0) on success
func doCreateValue(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		var err error
		params := convertKeyValueParams(args)

		var key, value, workspace, environment, typ, comment string
		if key, err = valueFromKVs(argKeyShort, argKeyLong, params); err != nil {
			return err
		}
		if value, err = valueFromKVs(argSecretShort, argSecretLong, params); err != nil {
			return err
		}
		if workspace, err = valueFromKVs(argWorkspaceShort, argWorkspaceLong, params); err != nil {
			return err
		}
		if environment, err = valueFromKVs(argEnvironmentShort, argEnvironmentLong, params); err != nil {
			return err
		}
		if typ, err = valueFromKVs(argTypeShort, argTypeLong, params); err != nil {
			return err
		}
		comment, _ = valueFromKVs(argCommentShort, argCommentLong, params)

		// key => folder + secretKey
		splitted := strings.Split(key, "/")
		secretKey := splitted[len(splitted)-1]
		folder := strings.TrimSuffix(key, secretKey)

		// request params
		createParams := infisical.NewParamsCreateSecret().
			SetType(infisical.SecretType(typ))
		if folder != "" {
			createParams.SetSecretPath(folder)
		}
		if comment != "" {
			createParams.SetSecretComment(comment)
		}

		// create
		err = c.CreateSecret(workspace, environment, secretKey, value, createParams)
		if err == nil {
			fmt.Printf("> Successfully created a new secret value at: %s\n", key)

			os.Exit(0)
		}

		return err
	}, verbose)
}

// update a secret value, will os.Exit(0) on success
func doUpdateValue(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		var err error
		params := convertKeyValueParams(args)

		var key, value, workspace, environment, typ string
		if key, err = valueFromKVs(argKeyShort, argKeyLong, params); err != nil {
			return err
		}
		if value, err = valueFromKVs(argSecretShort, argSecretLong, params); err != nil {
			return err
		}
		if workspace, err = valueFromKVs(argWorkspaceShort, argWorkspaceLong, params); err != nil {
			return err
		}
		if environment, err = valueFromKVs(argEnvironmentShort, argEnvironmentLong, params); err != nil {
			return err
		}
		if typ, err = valueFromKVs(argTypeShort, argTypeLong, params); err != nil {
			return err
		}

		// key => folder + secretKey
		splitted := strings.Split(key, "/")
		secretKey := splitted[len(splitted)-1]
		folder := strings.TrimSuffix(key, secretKey)

		// request params
		updateParams := infisical.NewParamsUpdateSecret().
			SetType(infisical.SecretType(typ))
		if folder != "" {
			updateParams.SetSecretPath(folder)
		}

		// update
		err = c.UpdateSecret(workspace, environment, secretKey, value, updateParams)
		if err == nil {
			fmt.Printf("> Successfully updated a secret value at: %s\n", key)

			os.Exit(0)
		}

		return err
	}, verbose)
}

// delete a secret value, will os.Exit(0) on success
func doDeleteValue(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		var err error
		params := convertKeyValueParams(args)

		var key, workspace, environment, typ string
		if key, err = valueFromKVs(argKeyShort, argKeyLong, params); err != nil {
			return err
		}
		if workspace, err = valueFromKVs(argWorkspaceShort, argWorkspaceLong, params); err != nil {
			return err
		}
		if environment, err = valueFromKVs(argEnvironmentShort, argEnvironmentLong, params); err != nil {
			return err
		}
		if typ, err = valueFromKVs(argTypeShort, argTypeLong, params); err != nil {
			return err
		}

		// key => folder + secretKey
		splitted := strings.Split(key, "/")
		secretKey := splitted[len(splitted)-1]
		folder := strings.TrimSuffix(key, secretKey)

		// request params
		deleteParams := infisical.NewParamsDeleteSecret().
			SetType(infisical.SecretType(typ))
		if folder != "" {
			deleteParams.SetSecretPath(folder)
		}

		// delete
		err = c.DeleteSecret(workspace, environment, secretKey, deleteParams)
		if err == nil {
			fmt.Printf("> Successfully deleted a secret value at: %s\n", key)

			os.Exit(0)
		}

		return err
	}, verbose)
}
