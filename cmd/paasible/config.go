package main

import (
	"fmt"
	"strings"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	AppVersion string `mapstructure:"APP_VERSION"`

	Port    int    `mapstructure:"PAASIBLE_UI_PORT"`
	User    string `mapstructure:"PAASIBLE_USER"`
	Machine string `mapstructure:"PAASIBLE_MACHINE"`
}

type YamlConfigPaasible struct {
	CliVersion             string `mapstructure:"cli_version"`
	ProjectName            string `mapstructure:"project_name"`
	DataFolderRelativePath string `mapstructure:"data_folder"`
}

type YamlConfig struct {
	Paasible YamlConfigPaasible `mapstructure:"paasible"`
}

// Call to load the variables from env
func initConfig(
	envConfigName string,
	envConfigPath string,
	yamlConfigName string,
	yamlConfigPath string,
) (*EnvConfig, *YamlConfig, error) {
	// # Read os env
	viper.AutomaticEnv()

	// # Env
	envConfigViper := viper.New()

	// ## Tell viper the path/location of your env file. If it is root just add "."
	envConfigViper.AddConfigPath(envConfigPath)

	// ## Tell viper the name of your file
	envConfigViper.SetConfigName(envConfigName + ".env")

	// ## Tell viper the type of your file
	envConfigViper.SetConfigType("env")

	// envConfigViper.SetDefault("PORT", 8080)
	envConfigViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ## Viper reads all the variables from env file and log error if any found
	if err := envConfigViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, nil, fmt.Errorf("Error while reading env config: %w", err)
		}
	}

	envConfig := &EnvConfig{}

	// ## Viper unmarshals the loaded env varialbes into the struct
	if err := envConfigViper.Unmarshal(envConfig); err != nil {
		return nil, nil, err
	}

	if envConfig.Port == 0 {
		envConfig.Port = 8080
	}

	// # Yaml
	yamlConfigViper := viper.New()

	// ## Tell viper the path/location of your env file. If it is root just add "."
	yamlConfigViper.AddConfigPath(yamlConfigPath)

	// ## Tell viper the name of your file
	yamlConfigViper.SetConfigName(yamlConfigName)

	// ## Tell viper the type of your file
	yamlConfigViper.SetConfigType("yaml")

	// ## Viper reads all the variables from env file and log error if any found
	if err := yamlConfigViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, nil, fmt.Errorf("Error reading yaml config: %w", err)
		}
	}

	yamlConfig := &YamlConfig{}

	// ## Viper unmarshals the loaded env varialbes into the struct
	if err := yamlConfigViper.Unmarshal(yamlConfig); err != nil {
		return nil, nil, err
	}

	if yamlConfig.Paasible.DataFolderRelativePath == "" {
		yamlConfig.Paasible.DataFolderRelativePath = paasible.DATA_FOLDER_NAME
	}

	return envConfig, yamlConfig, nil
}
