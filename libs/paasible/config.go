package paasible

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigFile struct {
	FilePath string `mapstructure:"-"` // where file stored

	IncludePaths []string      `mapstructure:"include"`
	Includes     []*ConfigFile `mapstructure:"-"` // included config files

	Paasible     *PaasibleEntity
	UI           *UIEntity           `mapstructure:"web"`
	Auth         *AuthEntity         `mapstructure:"auth"`
	SshKeys      []SSHKeyEntity      `mapstructure:"ssh_keys"`
	Hosts        []HostEntity        `mapstructure:"hosts"`
	Inventories  []InventoryEntity   `mapstructure:"inventories"`
	Projects     []ProjectEntity     `mapstructure:"projects"`
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
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
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
	for _, includePath := range yamlConfig.IncludePaths {
		includedConfig, err := ParseConfig(includePath)
		if err != nil {
			return nil, fmt.Errorf("Error parsing included config '%s': %w", includePath, err)
		}

		yamlConfig.Includes = append(yamlConfig.Includes, includedConfig)
	}

	return yamlConfig, nil
}

func CliConfigYaml() string {
	return `paasible:
  cli_version: 0.0.2
  cli_env_path: ./paasible.env
`
}

func CliConfigEnv() string {
	return `PAASIBLE_UI_PORT=8080
PAASIBLE_USER=user
PAASIBLE_MACHINE=local
`
}
