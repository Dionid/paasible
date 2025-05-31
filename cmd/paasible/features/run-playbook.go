package features

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func InitRunPlaybookCmd(
	app *pocketbase.PocketBase,
	storage *paasible.EntityStorage,
	paasibleRootConfigFolderPath string,
	paasibleCliPwd string,
	paasibleDataFolderPath string,
) {
	runPlaybookCmd := &cobra.Command{
		Use:     "run",
		Short:   "Run playbook",
		Long:    "Run playbook with the given arguments",
		Example: "run <performance_id>",
		Run: func(cmd *cobra.Command, args []string) {
			// # Get Machine
			machineId := storage.Auth.Machine

			if machineId == "" {
				log.Fatal("Can't find machine id in config file!")
			}

			userId := storage.Auth.User

			if userId == "" {
				log.Fatal("Can't find user id in config file!")
			}

			performanceId := args[0]
			if performanceId == "" {
				log.Fatalf("Performance ID required: %s", performanceId)
			}

			// # Get the playbook from the database
			performance, ok := storage.Performances[performanceId]
			if !ok {
				log.Fatalf("Failed to find playbook with ID %s", args[0])
			}

			// # Get playbook
			for _, playbookId := range performance.Playbooks {
				playbook, ok := storage.PlaybookEntity[playbookId]
				if !ok {
					log.Fatalf("Failed to find playbook with ID %s", args[0])
				}

				// # Get playbook project
				project, ok := storage.Projects[playbook.ProjectId]
				if !ok {
					log.Fatalf("Failed to find project with ID %s", playbook.ProjectId)
				}

				// # Create inventory
				inventoriesPaths := make([]string, 0)

				// ## Add Invetories
				for _, inventoryId := range performance.Inventories {
					inventory, ok := storage.Inventories[inventoryId]
					if !ok {
						log.Fatalf("Failed to find inventory with ID %s", inventoryId)
					}

					// # Check if the inventory file exists
					fullInventoryPath := path.Join(
						paasibleRootConfigFolderPath,
						inventory.Path,
					)

					if _, err := os.Stat(fullInventoryPath); os.IsNotExist(err) {
						log.Fatalf("Inventory file %s does not exist", fullInventoryPath)
					}

					inventoriesPaths = append(inventoriesPaths, fullInventoryPath)
				}

				// ## Generate inventories based on Targets
				if len(performance.Targets) != 0 {
					inventoryByGroup := make(map[string]string)

					for _, playbookTarget := range performance.Targets {
						targetSshKey, ok := storage.SSHKeys[playbookTarget.SSHKey]
						if !ok {
							log.Fatalf("Failed to find SSH key with ID %s", playbookTarget.SSHKey)
						}

						host, ok := storage.Hosts[playbookTarget.Host]
						if !ok {
							log.Fatalf("Failed to find host with ID %s", playbookTarget.Host)
						}

						// # Create ssh file
						pathToSshKey := path.Join(
							paasibleRootConfigFolderPath,
							targetSshKey.PrivatePath,
						)

						group := "all"
						if playbookTarget.Group != "" {
							group = playbookTarget.Group
						}

						inventoryByGroup[group] = fmt.Sprintf(
							`%s ansible_host=%s ansible_ssh_user=%s ansible_ssh_private_key_file=%s`,
							host.Name,
							host.Address,
							playbookTarget.User,
							pathToSshKey,
						)
					}

					inventoryContent := ""

					for group, hosts := range inventoryByGroup {
						inventoryContent += fmt.Sprintf("[%s]\n%s\n\n", group, hosts)
					}

					if inventoryContent != "" {
						// # Create inventory file
						inventoryFilePath := path.Join(
							paasibleRootConfigFolderPath,
							project.LocalPath,
							"inventory.paasible.ini",
						)
						err := os.WriteFile(inventoryFilePath, []byte(inventoryContent), 0644)
						if err != nil {
							log.Fatalf("Failed to write inventory file: %v", err)
						}

						inventoriesPaths = append(inventoriesPaths, inventoryFilePath)
					}
				}

				// # Create ansible-playbook command
				ansiblePlaybookArgs := []string{}

				// # Inventories
				for _, inventoryPath := range inventoriesPaths {
					ansiblePlaybookArgs = append(
						ansiblePlaybookArgs,
						"-i", inventoryPath,
					)
				}

				// # Playbook path
				ansiblePlaybookArgs = append(
					ansiblePlaybookArgs,
					path.Join(
						paasibleRootConfigFolderPath,
						project.LocalPath,
						playbook.Path,
					),
				)

				err := RunAndSave(
					paasibleCliPwd,
					paasibleDataFolderPath,
					ansiblePlaybookArgs,
					machineId,
					userId,
					playbookId,
				)
				if err != nil {
					log.Fatalf("Error running ansible-playbook: %v", err)
				}
			}

			// # Print the result
			log.Println("Success")
		},
	}

	// runPlaybookCmd.Flags().StringP("id", "i", "", "ID of the performance to run")
	// runPlaybookCmd.MarkFlagRequired("id")

	app.RootCmd.AddCommand(runPlaybookCmd)
}
