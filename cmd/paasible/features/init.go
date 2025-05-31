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
	yamlConfigPath string,
) {
	initCmd := &cobra.Command{
		Use:     "init",
		Short:   "Run init",
		Long:    "Create initial files and folders",
		Example: "init",
		Run: func(cmd *cobra.Command, args []string) {
			// # Check and create yaml config file
			if _, err := os.Stat(yamlConfigPath); os.IsNotExist(err) {
				content := paasible.PaasibleDefaultConfigYaml()
				if err := os.WriteFile(yamlConfigPath, []byte(content), 0644); err != nil {
					log.Fatalf("Failed to write yaml config file: %v", err)
				}
			}

			// # Check and create hidden yaml config file
			hiddenYamlConfigPath := path.Join(
				path.Dir(yamlConfigPath),
				"paasible.hidden.yaml",
			)
			if _, err := os.Stat(hiddenYamlConfigPath); os.IsNotExist(err) {
				content := paasible.PaasibleDefaultHiddenConfigYaml()
				if err := os.WriteFile(hiddenYamlConfigPath, []byte(content), 0644); err != nil {
					log.Fatalf("Failed to write yaml config file: %v", err)
				}
			}

			app.Logger().Info("Initial files created")
		},
	}

	app.RootCmd.AddCommand(initCmd)
}
