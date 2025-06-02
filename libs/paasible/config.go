package paasible

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/spf13/viper"
)

type ConfigFile struct {
	FilePath string `mapstructure:"-"` // where file stored

	IncludePaths []string      `mapstructure:"include"`
	Includes     []*ConfigFile `mapstructure:"-"` // included config files

	Paasible *PaasibleEntity
	UI       *UIEntity   `mapstructure:"ui"`
	Auth     *AuthEntity `mapstructure:"auth"`

	SshKeys      []SSHKeyEntity      `mapstructure:"ssh_keys"`
	Hosts        []HostEntity        `mapstructure:"hosts"`
	Inventories  []InventoryEntity   `mapstructure:"inventories"`
	Projects     []ProjectEntity     `mapstructure:"projects"`
	Playbooks    []PlaybookEntity    `mapstructure:"playbooks"`
	Performances []PerformanceEntity `mapstructure:"performances"`
}

func ParseConfig(
	yamlConfigPath string,
) (*ConfigFile, error) {
	// # Yaml
	yamlConfigViper := viper.New()

	// ## Tell viper the path/location of your env file. If it is root just add "."
	yamlConfigViper.SetConfigFile(
		yamlConfigPath,
	)

	// ## Viper reads all the variables from env file and log error if any found
	if err := yamlConfigViper.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			_, ok = err.(*fs.PathError)
		}

		if !ok {
			return nil, fmt.Errorf("Error reading yaml config: %w", err)
		}
	}

	yamlConfig := &ConfigFile{}

	yamlConfig.FilePath = yamlConfigPath

	// ## Viper unmarshals the loaded env varialbes into the struct
	if err := yamlConfigViper.Unmarshal(yamlConfig); err != nil {
		return nil, err
	}

	// ## Load includes
	for _, includePathRaw := range yamlConfig.IncludePaths {
		includePath := path.Join(
			path.Dir(yamlConfig.FilePath),
			includePathRaw,
		)

		includedConfig, err := ParseConfig(
			includePath,
		)
		if err != nil {
			return nil, fmt.Errorf("Error parsing included config '%s': %w", includePath, err)
		}

		yamlConfig.Includes = append(yamlConfig.Includes, includedConfig)
	}

	return yamlConfig, nil
}

func PaasibleDefaultConfigYaml() string {
	return `paasible:
  cli_version: 0.0.2

include:
  - ./paasible.hidden.yaml

# Custom inventories
# inventories:
#  - name: default_inventory
#    description: Default inventory for the project
#    path: ./inventory.ini

# Describe the project
# projects:
#   - id: blank_project
#     name: Simple project
#     description: A simple example of a Paasible project
#     local_path: .
#     # Describe project playbooks
#     playbooks:
#       - id: blank_playbook
#         name: Blank playbook
#         description: A simple example of a Paasible project playbook
#         path: ./playbook.yml

# # Describe the performances (playbook executions)
# performances:
#   - id: simple_performance
#     name: Simple performance
#     description: A simple performance that runs the first playbook.
#     playbooks:
#       - project: blank_project
#         playbook: blank_playbook
#	  inventories:
#		- default_inventory
#     targets:
#       - host: localhost
#         ssh_key: local_ssh_key
#         user: root
#         group: all
#         port: 22

# # Desctibe hosts
# hosts:
#   - name: localhost
#     address: 127.0.0.1
#     ssh_keys: # ssh keys to access the host
#       - local_ssh_key

# # Describe ssh keys
# ssh_keys:
#   - name: local_ssh_key
#     description: Local SSH key for the project
#     private_path: ~/.ssh/id_rsa
#     public_path: ~/.ssh/id_rsa.pub
`
}

func PaasibleDefaultHiddenConfigYaml() string {
	return `
ui:
  port: 8080 # port for paasible UI

auth:
  user: user # user name that will be stored in run results
  machine: local # machine name that will be stored in run results
`
}
