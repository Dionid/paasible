package features

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/Dionid/paasible/libs/paasible"
	"github.com/Dionid/paasible/libs/sqlify"
	"github.com/pocketbase/pocketbase"
)

func UpsertPlaybookRunResult(
	app *pocketbase.PocketBase,
	filePath string,
) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// # Parse file content to paasible.PlaybookRunResult
	var playbookRunResult paasible.PlaybookRunResult
	err = json.Unmarshal(fileContent, &playbookRunResult)
	if err != nil {
		return err
	}

	// # Insert new TaskRunResult in DB
	err = sqlify.Upsert(
		paasible.PlaybookRunResultModel(app, &playbookRunResult),
	)
	if err != nil {
		return err
	}

	return nil
}

func ParsePaasibleData(
	app *pocketbase.PocketBase,
	rootFolderPath string,
) error {
	// # Read all json files from rootFolderPath/plaubook_run_result and
	// parse them to paasible.PlaybookRunResult and insert into DB

	// # Get all files in the folder
	files, err := os.ReadDir(path.Join(rootFolderPath, paasible.DATA_PLAYBOOK_RUN_RESULT_FOLDER_NAME))
	if err != nil {
		return err
	}

	// # Loop through all files
	for _, file := range files {
		// # Check if file is json
		if path.Ext(file.Name()) != ".json" {
			continue
		}

		// # Read file
		filePath := path.Join(rootFolderPath, paasible.DATA_PLAYBOOK_RUN_RESULT_FOLDER_NAME, file.Name())
		err = UpsertPlaybookRunResult(app, filePath)
		if err != nil {
			return fmt.Errorf("ParsePaasibleData: %w", err)
		}
	}

	return nil
}
