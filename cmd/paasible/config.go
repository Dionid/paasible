package main

import (
	"fmt"
	"io/fs"
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
) (*paasible.ConfigFile, error) {
	// # Read os env
	viper.AutomaticEnv()

	// # Yaml
	yamlConfig, err := paasible.ParseConfig(yamlConfigPath)
	if err != nil {
		return nil, err
	}

	// # Env
	envConfigViper := viper.New()

	if yamlConfig.Paasible == nil || yamlConfig.Paasible.CliEnvPath == "" {
		envConfigViper.SetConfigFile("./paasible.env")
	} else {
		// ## Tell viper the path/location of your env file. If it is root just add "."
		envConfigViper.SetConfigFile(yamlConfig.Paasible.CliEnvPath)
	}

	// ## Tell viper the type of your file
	envConfigViper.SetConfigType("env")

	envConfigViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ## Viper reads all the variables from env file and log error if any found
	if err := envConfigViper.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			_, ok = err.(*fs.PathError)
		}

		if !ok {
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

	yamlConfig.Auth = &paasible.AuthEntity{
		User:    envConfig.User,
		Machine: envConfig.Machine,
	}

	yamlConfig.UI = &paasible.UIEntity{
		Port: envConfig.Port,
	}

	return yamlConfig, nil
}
