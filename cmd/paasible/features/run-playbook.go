package features

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func InitRunPlaybookCmd(
	app *pocketbase.PocketBase,
	config *paasible.CliConfig,
	paasibleDataFolderPath string,
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
			machineId := config.Machine

			if machineId == "" {
				log.Fatal("Can't find machine id in config file!")
			}

			userId := config.User

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

			// # Create inventory
			playbookTargets := []paasible.PlaybookTarget{}

			err = paasible.PlaybookTargetQuery(app).Where(dbx.HashExp{
				"playbook_id": playbookId,
			}).All(&playbookTargets)
			if err != nil {
				log.Fatalf("Failed to find targets for playbook with ID %s: %v", playbookId, err)
			}

			if len(playbookTargets) == 0 {
				log.Fatalf("No targets found for playbook with ID %s", playbookId)
			}

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
					paasibleDataFolderPath,
					paasible.DATA_PLAYBOOKS_FOLDER_NAME,
					playbook.CodePath,
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

			inventory := ""

			for group, hosts := range inventoryByGroup {
				inventory += fmt.Sprintf("[%s]\n%s\n\n", group, hosts)
			}

			// # Create inventory file
			inventoryFilePath := path.Join(
				paasibleDataFolderPath,
				paasible.DATA_PLAYBOOKS_FOLDER_NAME,
				playbook.CodePath,
				"inventory.ini",
			)
			err = os.WriteFile(inventoryFilePath, []byte(inventory), 0644)
			if err != nil {
				log.Fatalf("Failed to write inventory file: %v", err)
			}

			// # Create ansible-playbook command
			ansiblePlaybookArgs := []string{
				"-i", inventoryFilePath,
				path.Join(
					paasibleDataFolderPath,
					paasible.DATA_PLAYBOOKS_FOLDER_NAME,
					playbook.CodePath,
					playbook.Path,
				),
			}

			err = RunAndSave(
				currentFolderPath,
				paasibleDataFolderPath,
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
