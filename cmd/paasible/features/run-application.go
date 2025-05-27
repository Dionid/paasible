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

func InitRunAppilcationCmd(
	app *pocketbase.PocketBase,
	config *paasible.CliConfig,
	paasibleDataFolderPath string,
) {
	runApplicationCmd := &cobra.Command{
		Use:     "run",
		Short:   "Run application",
		Long:    "Run application with the given arguments",
		Example: "run --id <application_id>",
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

			// # Get the application ID from the command line arguments
			applicationId, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatalf("Failed to get application ID: %v", err)
			}

			// # Get the application from the database
			application := paasible.Application{}

			err = paasible.ApplicationQuery(app).Where(dbx.HashExp{
				"id": applicationId,
			}).One(&application)
			if err != nil {
				log.Fatalf("Failed to find application with ID %s: %v", args[0], err)
			}

			// # Create inventory
			applicationTargets := []paasible.ApplicationTarget{}

			err = paasible.ApplicationTargetQuery(app).Where(dbx.HashExp{
				"application_id": applicationId,
			}).All(&applicationTargets)
			if err != nil {
				log.Fatalf("Failed to find targets for application with ID %s: %v", applicationId, err)
			}

			if len(applicationTargets) == 0 {
				log.Fatalf("No targets found for application with ID %s", applicationId)
			}

			inventory := ""

			for _, applicationTarget := range applicationTargets {
				target := paasible.Target{}

				err = paasible.TargetQuery(app).Where(dbx.HashExp{
					"id": applicationTarget.TargetId,
				}).One(&target)
				if err != nil {
					log.Fatalf("Failed to find target with ID %s: %v", applicationTarget.TargetId, err)
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
					paasible.DATA_APPLICATIONS_FOLDER_NAME,
					application.Path,
					fmt.Sprintf("%s_%s.ssh_key", target.Id, sshKey.Name),
				)

				// ## Create the SSH key file
				err = os.WriteFile(pathToSshKey, []byte(sshKey.Private), 0600)
				if err != nil {
					log.Fatalf("Failed to write SSH key file: %v", err)
				}

				inventory += fmt.Sprintf(`[%s]
%s ansible_host=%s ansible_ssh_user=%s ansible_ssh_private_key_file=%s
`, application.Name, application.Name, target.Address, applicationTarget.User, pathToSshKey)
			}

			// # Create inventory file
			inventoryFilePath := path.Join(
				paasibleDataFolderPath,
				paasible.DATA_APPLICATIONS_FOLDER_NAME,
				application.Path,
				"inventory.ini",
			)
			err = os.WriteFile(inventoryFilePath, []byte(inventory), 0644)
			if err != nil {
				log.Fatalf("Failed to write inventory file: %v", err)
			}
			log.Printf("Inventory file created at %s", inventoryFilePath)

			// # Create ansible-playbook command
			ansiblePlaybookArgs := []string{
				"-i", inventoryFilePath,
				path.Join(
					paasibleDataFolderPath,
					paasible.DATA_APPLICATIONS_FOLDER_NAME,
					application.Path,
					application.PlaybookPath,
				),
			}

			err = RunAndSave(
				currentFolderPath,
				paasibleDataFolderPath,
				ansiblePlaybookArgs,
				machineId,
				userId,
				applicationId,
			)
			if err != nil {
				log.Fatalf("Error running ansible-playbook: %v", err)
			}

			// # Print the result
			log.Println("Success")
		},
	}

	runApplicationCmd.Flags().StringP("id", "i", "", "ID of the application to run")
	runApplicationCmd.MarkFlagRequired("id")

	app.RootCmd.AddCommand(runApplicationCmd)
}
