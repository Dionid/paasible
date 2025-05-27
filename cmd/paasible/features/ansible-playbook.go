package features

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/Dionid/paasible/libs/uuidv7"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cobra"
)

func InitAnsiblePlaybookCmd(
	app *pocketbase.PocketBase,
	config *paasible.CliConfig,
) {
	ansiblePlaybookCmd := &cobra.Command{
		Use:     "ansible-playbook",
		Short:   "Run ansible-playbook",
		Long:    "Run ansible-playbook with the given arguments",
		Args:    cobra.ArbitraryArgs,
		Example: "ansible-playbook playbook.yml",
		Run: func(cmd *cobra.Command, args []string) {
			// # Current folder path
			currentFolderPath, err := os.Getwd()
			if err != nil {
				log.Fatal("Error getting current folder path:", err)
			}

			// # Take playbook relative path
			playbookRelativePath := ""
			for _, ar := range args {
				if strings.Contains(ar, ".yml") || strings.Contains(ar, ".yaml") {
					playbookRelativePath = ar
					break
				}
			}

			if playbookRelativePath == "" {
				log.Fatal("Error getting relative path")
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

			// # Create new TaskRunResult
			playbookRunResult := paasible.RunResult{
				BaseModel: core.BaseModel{
					Id: uuidv7.NewE().String(),
				},
				Created:   types.NowDateTime(),
				Updated:   types.NowDateTime(),
				Pwd:       currentFolderPath,
				MachineId: machineId,
				UserId:    userId,
			}

			// ## Check branch
			gitBranchResult := paasible.ExecGitCommand(
				exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD"),
				currentFolderPath,
			)
			if gitBranchResult.Error == nil {
				playbookRunResult.RepositoryBranch = gitBranchResult.Stdout
			}

			// # Run ansible-playbook
			command := exec.Command("ansible-playbook", args...)

			playbookRunResult.Command = command.String()

			// ## Add to std out
			var stdoutBuf, stderrBuf bytes.Buffer

			command.Stdout = io.MultiWriter(&stdoutBuf, os.Stdout)
			command.Stderr = io.MultiWriter(&stderrBuf, os.Stderr)

			// Run the command
			err = command.Run()
			if err != nil {
				playbookRunResult.Error = err.Error()
				log.Println("Error running ansible-playbook:", err)
			}

			// # Add to std out
			playbookRunResult.Stdout = stdoutBuf.String()
			playbookRunResult.Stderr = stderrBuf.String()

			// # Create new TaskRunResult json

			// ## Marshal
			playbookRunResultJson, err := json.MarshalIndent(playbookRunResult, "", "  ")
			if err != nil {
				log.Fatal("Error marshaling playbookRunResult to json:", err)
			}

			// ## Create paasible_data folder if not exists
			err = paasible.CreateDataFolder(currentFolderPath)
			if err != nil {
				log.Fatal("Error creating data folder: ", err)
			}

			// ## Write playbookRunResult to file
			unixTimestampString := strconv.FormatInt(playbookRunResult.Created.Time().Unix(), 10)

			playbookRunResultFileName := path.Join(
				paasible.DATA_RUN_RESULT_FOLDER_PATH,
				unixTimestampString+"__"+config.User+"__"+config.Machine+"__"+playbookRunResult.Id+".json",
			)
			playbookRunResultFile, err := os.Create(playbookRunResultFileName)
			if err != nil {
				log.Fatal("Error creating playbookRunResult file:", err)
			}
			defer playbookRunResultFile.Close()

			// # Write playbookRunResult to file
			_, err = playbookRunResultFile.Write(playbookRunResultJson)
			if err != nil {
				log.Fatal("Error writing playbookRunResult to file:", err)
			}

			// // # Insert new TaskRunResult in DB
			// err = paasible.TaskRunResultModel(app, &playbookRunResult).Insert()
			// if err != nil {
			// 	log.Println("Error inserting TaskRunResult:", err)
			// 	return
			// }

			if playbookRunResult.Error != "" {
				log.Fatal(playbookRunResult.Error)
			}

			// # Print the result
			log.Println("Success")
		},
	}

	app.RootCmd.AddCommand(ansiblePlaybookCmd)
}
