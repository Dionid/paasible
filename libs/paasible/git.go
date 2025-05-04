package paasible

import (
	"bytes"
	"log"
	"os/exec"
	"strings"

	"github.com/Dionid/paasible/libs/uuidv7"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ExecGitCommandResult struct {
	Stdout string
	Stderr string
	Error  error
}

func ExecGitCommand(
	cmd *exec.Cmd,
	currentFolderPath string,
) ExecGitCommandResult {
	cmd.Dir = currentFolderPath
	var cmdStderr bytes.Buffer
	cmd.Stderr = &cmdStderr

	revParseCmdOutput, err := cmd.Output()
	revParseCmdOutputStr := strings.TrimSpace(string(revParseCmdOutput))
	stdErrString := cmdStderr.String()

	return ExecGitCommandResult{
		Stdout: revParseCmdOutputStr,
		Stderr: stdErrString,
		Error:  err,
	}
}

// ALL CODE TILL THE END OF THE FILE IS FOR FUTURE VERSIONS
func GetAndSyncOrCreateRepository(
	app *pocketbase.PocketBase,
	currentFolderPath string,
) (*Repository, error) {
	// # Find or create Repository
	repository := Repository{}

	// ## Check if repository exists by remote
	getUrlResult := ExecGitCommand(
		exec.Command("git", "remote", "get-url", "origin"),
		currentFolderPath,
	)
	if getUrlResult.Error != nil {
		if !strings.Contains(getUrlResult.Stderr, "No such remote") {
			log.Println("Error running git command remote:", getUrlResult.Error.Error(), " ", getUrlResult.Stderr)
			return nil, getUrlResult.Error
		}
	}

	// ## Check if repository exists by top_level
	revParseResult := ExecGitCommand(
		exec.Command("git", "rev-parse", "--show-toplevel"),
		currentFolderPath,
	)
	if revParseResult.Error != nil {
		if !strings.Contains(revParseResult.Stderr, "No such remote") {
			log.Println("Error running git command rev-parse:", revParseResult.Error)
			return nil, revParseResult.Error
		}
	}

	// ## Check if repository exists by top_level or remote
	err := RepositoryQuery(app).Where(dbx.HashExp{
		"top_level": revParseResult.Stdout,
	}).OrWhere(
		dbx.HashExp{
			"remote": getUrlResult.Stdout,
		},
	).One(&repository)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return nil, err
		}

		// # Create new Repository
		repository = Repository{
			BaseModel: core.BaseModel{
				Id: uuidv7.NewE().String(),
			},
			Created:  types.NowDateTime(),
			Updated:  types.NowDateTime(),
			Remote:   getUrlResult.Stdout,
			TopLevel: revParseResult.Stdout,
		}

		err = RepositoryModel(app, &repository).Insert()
		if err != nil {
			log.Println("Error inserting Repository:", err)
			return nil, err
		}
	}

	// # Update Repository
	repository.Remote = getUrlResult.Stdout
	repository.TopLevel = revParseResult.Stdout
	repository.Updated = types.NowDateTime()

	err = RepositoryModel(app, &repository).Update()
	if err != nil {
		log.Println("Error updating Repository:", err)
		return nil, err
	}

	return &repository, nil
}
