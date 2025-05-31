package features

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func InitRunPlaybookCmd(
	app *pocketbase.PocketBase,
	config *paasible.ConfigFile,
	paasibleRootFolderPath string,
) {
	runPlaybookCmd := &cobra.Command{
		Use:     "run",
		Short:   "Run playbook",
		Long:    "Run playbook with the given arguments",
		Example: "run --id <playbook_id>",
		Run: func(cmd *cobra.Command, args []string) {
			// # Current folder path
			currentFolderPath, err := os.Getwd()
			if err != nil {
				log.Fatal("Error getting current folder path:", err)
			}

			// # Get Machine
			machineId := config.Auth.Machine

			if machineId == "" {
				log.Fatal("Can't find machine id in config file!")
			}

			userId := config.Auth.User

			if userId == "" {
				log.Fatal("Can't find user id in config file!")
			}

			// # Get the playbook ID from the command line arguments
			playbookId, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatalf("Failed to get playbook ID: %v", err)
			}

			// # Get the playbook from the database
			playbook := paasible.Playbook{}

			err = paasible.PlaybookQuery(app).Where(dbx.HashExp{
				"id": playbookId,
			}).One(&playbook)
			if err != nil {
				log.Fatalf("Failed to find playbook with ID %s: %v", args[0], err)
			}

			// # Get playbook project
			project := paasible.Project{}

			err = paasible.ProjectQuery(app).Where(dbx.HashExp{
				"id": playbook.ProjectId,
			}).One(&project)
			if err != nil {
				log.Fatalf("Failed to find project with ID %s: %v", playbook.ProjectId, err)
			}

			// # Create inventory
			inventoriesPaths := make([]string, 0)

			// ## From Playbook
			if playbook.InventoriesPaths != "" {
				splitedPaths := strings.Split(playbook.InventoriesPaths, ",")

				for _, inventoryPath := range splitedPaths {
					inventoryPath = strings.TrimSpace(inventoryPath)

					if inventoryPath == "" {
						continue
					}

					// # Check if the inventory file exists
					fullInventoryPath := path.Join(
						paasibleRootFolderPath,
						project.Path,
						inventoryPath,
					)

					if _, err := os.Stat(fullInventoryPath); os.IsNotExist(err) {
						log.Fatalf("Inventory file %s does not exist", fullInventoryPath)
					}

					inventoriesPaths = append(inventoriesPaths, fullInventoryPath)
				}
			}

			// ## From Targets
			playbookTargets := []paasible.PlaybookTarget{}

			err = paasible.PlaybookTargetQuery(app).Where(dbx.HashExp{
				"playbook_id": playbookId,
			}).All(&playbookTargets)
			if err != nil {
				log.Fatalf("Failed to find targets for playbook with ID %s: %v", playbookId, err)
			}

			if len(playbookTargets) != 0 {
				inventoryByGroup := make(map[string]string)

				for _, playbookTarget := range playbookTargets {
					target := paasible.Target{}

					err = paasible.TargetQuery(app).Where(dbx.HashExp{
						"id": playbookTarget.TargetId,
					}).One(&target)
					if err != nil {
						log.Fatalf("Failed to find target with ID %s: %v", playbookTarget.TargetId, err)
					}

					targetSshKey := paasible.TargetSshKey{}

					err = paasible.TargetSshKeyQuery(app).Where(dbx.HashExp{
						"target_id": target.Id,
					}).One(&targetSshKey)
					if err != nil {
						log.Fatalf("Failed to find SSH key for target with ID %s: %v", target.Id, err)
					}

					sshKey := paasible.SshKey{}

					err = paasible.SshKeyQuery(app).Where(dbx.HashExp{
						"id": targetSshKey.SshKeyId,
					}).One(&sshKey)
					if err != nil {
						log.Fatalf("Failed to find SSH key with ID %s: %v", targetSshKey.SshKeyId, err)
					}

					// # Create ssh file
					pathToSshKey := path.Join(
						paasibleRootFolderPath,
						project.Path,
						fmt.Sprintf("%s_%s.ssh_key", target.Id, sshKey.Name),
					)

					// ## Create the SSH key file
					err = os.WriteFile(pathToSshKey, []byte(sshKey.Private), 0600)
					if err != nil {
						log.Fatalf("Failed to write SSH key file: %v", err)
					}

					inventoryByGroup[playbookTarget.Group] = fmt.Sprintf(
						`%s ansible_host=%s ansible_ssh_user=%s ansible_ssh_private_key_file=%s`,
						playbook.Name,
						target.Address,
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
						paasibleRootFolderPath,
						project.Path,
						"generated_inventory.ini",
					)
					err = os.WriteFile(inventoryFilePath, []byte(inventoryContent), 0644)
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
					paasibleRootFolderPath,
					project.Path,
					playbook.Path,
				),
			)

			err = RunAndSave(
				currentFolderPath,
				paasibleRootFolderPath,
				ansiblePlaybookArgs,
				machineId,
				userId,
				playbookId,
			)
			if err != nil {
				log.Fatalf("Error running ansible-playbook: %v", err)
			}

			// # Print the result
			log.Println("Success")
		},
	}

	runPlaybookCmd.Flags().StringP("id", "i", "", "ID of the playbook to run")
	runPlaybookCmd.MarkFlagRequired("id")

	app.RootCmd.AddCommand(runPlaybookCmd)
}
