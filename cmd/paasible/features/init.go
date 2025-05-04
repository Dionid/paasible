package features

import (
	"log"
	"os"
	"path"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func InitInitCmd(
	app *pocketbase.PocketBase,
	envConfigName string,
	envConfigPath string,
	yamlConfigName string,
	yamlConfigPath string,
) {
	initCmd := &cobra.Command{
		Use:     "init",
		Short:   "Run init",
		Long:    "Create initial files and folders",
		Example: "init",
		Run: func(cmd *cobra.Command, args []string) {
			yamlPath := path.Join(yamlConfigPath, yamlConfigName+".yaml")

			// Check and create yaml config file
			if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
				content := paasible.CliConfigYaml()
				if err := os.WriteFile(yamlPath, []byte(content), 0644); err != nil {
					log.Fatalf("Failed to write yaml config file: %v", err)
				}
			}

			envPath := path.Join(envConfigPath, envConfigName+".env")
			// Check and create yaml config file
			if _, err := os.Stat(envPath); os.IsNotExist(err) {
				content := paasible.CliConfigEnv()
				if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
					log.Fatalf("Failed to write env config file: %v", err)
				}
			}

			app.Logger().Info("Initial files created")
		},
	}

	app.RootCmd.AddCommand(initCmd)
}
