package paasible

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Path         string              `yaml:"-"`
	UI           UIConfig            `yaml:"web"`
	Auth         AuthConfig          `yaml:"auth"`
	Paasible     PaasibleConfig      `yaml:"paasible"`
	SSHKeys      []SSHKeyConfig      `yaml:"ssh_keys"`
	Hosts        []HostConfig        `yaml:"hosts"`
	Inventories  []InventoryConfig   `yaml:"inventories"`
	Projects     []ProjectConfig     `yaml:"projects"`
	Performances []PerformanceConfig `yaml:"performances"`
}

func ParseConfig(
	yamlConfigPath string,
) (*Config, error) {
	// # Yaml
	yamlConfigViper := viper.New()

	// ## Tell viper the path/location of your env file. If it is root just add "."
	yamlConfigViper.AddConfigPath(yamlConfigPath)

	// ## Tell viper the type of your file
	yamlConfigViper.SetConfigType("yaml")

	// ## Viper reads all the variables from env file and log error if any found
	if err := yamlConfigViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Error reading yaml config: %w", err)
		}
	}

	yamlConfig := &Config{}

	yamlConfig.Path = yamlConfigPath

	// ## Viper unmarshals the loaded env varialbes into the struct
	if err := yamlConfigViper.Unmarshal(yamlConfig); err != nil {
		return nil, err
	}

	return yamlConfig, nil
}

type UIConfig struct {
	Port int `yaml:"port"`
}

type AuthConfig struct {
	User    string
	Machine string
}

type PaasibleConfig struct {
	CliVersion string `yaml:"cli_version"`
	CliEnvPath string `yaml:"cli_env_path"`

	DataFolderPath string `yaml:"data_folder_path"`
}

type SSHKeyConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	PrivatePath string `yaml:"private_path"`
	PublicPath  string `yaml:"public_path"`
}

type HostConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Address     string   `yaml:"address"`
	SSHKeys     []string `yaml:"ssh_keys"`
}

type InventoryConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}

type ProjectConfig struct {
	ID               string     `yaml:"id"`
	Name             string     `yaml:"name"`
	Description      string     `yaml:"description"`
	Version          string     `yaml:"version"`
	CliVersion       string     `yaml:"cli_version"`
	Repository       string     `yaml:"repository"`
	RepositoryBranch string     `yaml:"repository_branch"`
	RepositoryPath   string     `yaml:"repository_path"`
	LocalPath        string     `yaml:"local_path"`
	Playbooks        []Playbook `yaml:"playbooks"`
}

type PlaybookConfig struct {
	ID              string               `yaml:"id"`
	Name            string               `yaml:"name"`
	Description     string               `yaml:"description"`
	Path            string               `yaml:"path"`
	VariablesSchema VariableSchemaConfig `yaml:"variables_schema"`
}

type VariableSchemaConfig struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Schema      map[string]string `yaml:"schema"`
}

type PerformanceConfig struct {
	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Playbooks   []string         `yaml:"playbooks"`
	Inventories []string         `yaml:"inventories"`
	Variables   []VariableConfig `yaml:"variables"`
	Targets     []Target         `yaml:"targets"`
}

type VariableConfig struct {
	Path string `yaml:"path"`
}

type TargetConfig struct {
	Host   string `yaml:"host"`
	SSHKey string `yaml:"ssh_key"`
	User   string `yaml:"user"`
	Port   int    `yaml:"port"`
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
