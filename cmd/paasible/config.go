package main

import (
	"github.com/Dionid/paasible/libs/paasible"
	"github.com/spf13/viper"
)

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

	return yamlConfig, nil
}
