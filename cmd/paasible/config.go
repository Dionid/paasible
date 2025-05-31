package main

import (
	"fmt"
	"strings"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	AppVersion string `mapstructure:"APP_VERSION"`
	Port       int    `mapstructure:"PAASIBLE_UI_PORT"`
	User       string `mapstructure:"PAASIBLE_USER"`
	Machine    string `mapstructure:"PAASIBLE_MACHINE"`
}

// Call to load the variables from env
func initConfig(
	yamlConfigPath string,
) (*paasible.Config, error) {
	// # Read os env
	viper.AutomaticEnv()

	// # Yaml
	yamlConfig, err := paasible.ParseConfig(yamlConfigPath)
	if err != nil {
		return nil, err
	}

	// # Env
	envConfigViper := viper.New()

	if yamlConfig.Paasible.CliEnvPath == "" {
		envConfigViper.AddConfigPath(".")

		// ## Tell viper the name of your file
		envConfigViper.SetConfigName("paasible.env")
	} else {
		// ## Tell viper the path/location of your env file. If it is root just add "."
		envConfigViper.AddConfigPath(yamlConfig.Paasible.CliEnvPath)
	}

	// ## Tell viper the type of your file
	envConfigViper.SetConfigType("env")

	// envConfigViper.SetDefault("PORT", 8080)
	envConfigViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ## Viper reads all the variables from env file and log error if any found
	if err := envConfigViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Error while reading env config: %w", err)
		}
	}

	envConfig := &EnvConfig{}

	// ## Viper unmarshals the loaded env varialbes into the struct
	if err := envConfigViper.Unmarshal(envConfig); err != nil {
		return nil, err
	}

	if envConfig.Port == 0 {
		envConfig.Port = 8080
	}

	yamlConfig.Auth.User = envConfig.User
	yamlConfig.Auth.Machine = envConfig.Machine

	yamlConfig.UI.Port = envConfig.Port

	return yamlConfig, nil
}
