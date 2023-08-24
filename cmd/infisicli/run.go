package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/meinside/infisical-go"
	"github.com/meinside/version-go"
)

const (
	applicationName = "infisicli"
	configFilename  = "config.json"
)

const (
	// arguments for commands
	argHelpShort              = "-h"
	argHelpLong               = "--help"
	argListOrganizationsShort = "-lo"
	argListOrganizationsLong  = "--list-organizations"
	argListWorkspacesShort    = "-lw"
	argListWorkspacesLong     = "--list-workspaces"
	argListAllSecretsShort    = "-las"
	argListAllSecretsLong     = "--list-all-secrets"
	argPrintValueShort        = "-p"
	argPrintValueLong         = "--print"
	argVersionShort           = "-v"
	argVersionLong            = "--version"
	argVerboseShort           = "-V"
	argVerboseLong            = "--verbose"

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
	argFolderShort       = "-f"
	argFolderLong        = "--folder"

	// regex for parsing command param arguments (in 'key=value' format)
	regexKeyValue = `(.*?)=['"]?(.*?)['"]?$`
)

// config struct
type config struct {
	// Infisical Account's API Key
	APIKey string `json:"api_key,omitempty"`

	// key = worksace ID
	// value = workspace token
	Workspaces map[string]infisical.WorkspaceToken `json:"workspaces"`
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
			if err = json.Unmarshal(bytes, &conf); err == nil {
				return conf, nil
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
	: print this help message to stdout.

	%[1]s %[4]s
	%[1]s %[5]s
	: print the version string to stdout.

	%[1]s %[6]s
	%[1]s %[7]s
	: list all your organizations to stdout. (needed: 'api_key' and 'token')

	%[1]s %[8]s
	%[1]s %[9]s
	: list all your workspaces to stdout. (needed: 'api_key' and 'token')
		eg. %[1]s %[8]s %[14]s=0a1b2c3d4e5f
		eg. %[1]s %[9]s %[15]s=0a1b2c3d4e5f

	%[1]s %[10]s
	%[1]s %[11]s
	: list all secret values in a given folder (default: /) to stdout. (only 'token' is needed when E2EE is disabled)
		eg. %[1]s %[10]s %[18]s=012345abcdefg %[20]s=dev
		eg. %[1]s %[11]s %[19]s=012345abcdefg %[21]s=dev
		eg. %[1]s %[10]s %[18]s=012345abcdefg %[20]s=dev %[16]s=/folder1/folder2
		eg. %[1]s %[11]s %[19]s=012345abcdefg %[21]s=dev %[17]s=/folder1/folder2

	%[1]s %[12]s
	%[1]s %[13]s
	: print the secret value to stdout without a trailing newline. (only 'token' is needed when E2EE is disabled)
		eg. %[1]s %[12]s %[18]s=012345abcdefg %[20]s=dev %[22]s=shared %[24]s=/folder1/folder2/SECRET_KEY_1
		eg. %[1]s %[13]s %[19]s=012345abcdefg %[21]s=dev %[23]s=shared %[25]s=/folder1/folder2/SECRET_KEY_1

Other optional arguments:

	%[26]s / %[27]s
	: dump http requests/responses for debugging.
`,
		// executable name
		applicationName,

		// commands
		argHelpShort, argHelpLong,
		argVersionShort, argVersionLong,
		argListOrganizationsShort, argListOrganizationsLong,
		argListWorkspacesShort, argListWorkspacesLong,
		argListAllSecretsShort, argListAllSecretsLong,
		argPrintValueShort, argPrintValueLong,

		// parameters
		argOrganizationShort, argOrganizationLong,
		argFolderShort, argFolderLong,
		argWorkspaceShort, argWorkspaceLong,
		argEnvironmentShort, argEnvironmentLong,
		argTypeShort, argTypeLong,
		argKeyShort, argKeyLong,

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
	if hasArg(args, argVersionShort, argVersionLong) {
		// show version
		showVersion()
	} else if hasArg(args, argHelpShort, argHelpLong) {
		// do nothing
	} else {
		//handle commands here
		if hasArg(args, argListOrganizationsShort, argListOrganizationsLong) {
			err = doListOrganizations(verbose)
		} else if hasArg(args, argListWorkspacesShort, argListWorkspacesLong) {
			err = doListWorkspaces(args, verbose)
		} else if hasArg(args, argListAllSecretsShort, argListAllSecretsLong) {
			err = doListAllSecrets(args, verbose)
		} else if hasArg(args, argPrintValueShort, argPrintValueLong) {
			err = doPrintValue(args, verbose)
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
		if cfg.APIKey != "" {
			client = infisical.NewClient(cfg.APIKey, cfg.Workspaces)
		} else {
			client = infisical.NewClientWithoutAPIKey(cfg.Workspaces)
		}
		client.Verbose = verbose

		return fn(client)
	}

	return err
}

// list organizations, will os.Exit(0) on success
func doListOrganizations(verbose bool) error {
	return do(func(c *infisical.Client) error {
		result, err := c.RetrieveOrganizations()
		if err == nil {
			// calculate max lengths for formatting
			maxLenOrg := maxLength(result.Organizations, func(org infisical.Organization) int {
				return len(org.Name)
			})
			maxLenID := maxLength(result.Organizations, func(org infisical.Organization) int {
				return len(org.ID)
			})
			format := fmt.Sprintf("%%%ds | %%-%ds\n", maxLenID, maxLenOrg)

			// print headers
			fmt.Printf(format, "id", "name")
			fmt.Printf("----\n")

			// print organizations
			for _, org := range result.Organizations {
				fmt.Printf(format, org.ID, org.Name)
			}

			os.Exit(0)
		}

		return err
	}, verbose)
}

// list workspaces, will os.Exit(0) on success
func doListWorkspaces(args []string, verbose bool) error {
	return do(func(c *infisical.Client) error {
		var err error
		params := convertKeyValueParams(args)

		var org string
		org, err = valueFromKVs(argOrganizationShort, argOrganizationLong, params)

		if err == nil {
			var workspaces infisical.ProjectsData
			workspaces, err = c.RetrieveProjects(org)

			if err == nil {
				// calculate max lengths for formatting
				maxLenWorkspaceID := maxLength(workspaces.Workspaces, func(workspace infisical.Workspace) int {
					return len(workspace.ID)
				})
				maxLenWorkspaceName := maxLength(workspaces.Workspaces, func(workspace infisical.Workspace) int {
					return len(workspace.Name)
				})
				workspaceFormat := fmt.Sprintf("%%%ds | %%-%ds\n", maxLenWorkspaceID, maxLenWorkspaceName)

				// print headers
				fmt.Printf(workspaceFormat, "workspace id", "name")

				for _, workspace := range workspaces.Workspaces {
					maxLenSlug := maxLength(workspace.Environments, func(env infisical.WorkspaceEnvironment) int {
						return len(env.Slug)
					})
					maxLenName := maxLength(workspace.Environments, func(env infisical.WorkspaceEnvironment) int {
						return len(env.Name)
					})
					envFormat := fmt.Sprintf("  %%%ds | %%%ds (%%s)\n", maxLenSlug, maxLenName)

					// print workspace
					fmt.Printf("----\n")
					fmt.Printf(workspaceFormat, workspace.ID, workspace.Name)

					// print environments
					for _, env := range workspace.Environments {
						fmt.Printf(envFormat, env.Slug, env.Name, env.ID)
					}
				}

				os.Exit(0)
			}
		}

		return err
	}, verbose)
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
		folder, _ := valueFromKVs(argFolderShort, argFolderLong, params)

		secretsParam := infisical.NewParamsRetrieveSecrets()
		if folder != "" {
			secretsParam.SetSecretPath(folder)
		}

		var result infisical.SecretsData
		result, err = c.RetrieveSecrets(workspace, environment, secretsParam)

		if err == nil {
			maxLenWorkspace := maxLength(result.Secrets, func(secret infisical.Secret) int {
				return len(secret.Workspace)
			})
			maxLenEnv := maxLength(result.Secrets, func(secret infisical.Secret) int {
				return len(secret.Environment)
			})
			maxLenType := maxLength(result.Secrets, func(secret infisical.Secret) int {
				return len(secret.Type)
			})
			format := fmt.Sprintf("%%%ds | %%%ds | %%%ds | %%s\n", maxLenWorkspace, maxLenEnv, maxLenType)

			// print headers
			fmt.Printf(format, "workspace", "env", "type", "path/key=value")
			fmt.Printf("----\n")

			// print key-values
			for _, secret := range result.Secrets {
				fmt.Printf(format, secret.Workspace, secret.Environment, secret.Type, path.Join(folder, secret.SecretKey)+"="+secret.SecretValue)
			}

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
