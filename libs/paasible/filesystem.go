package paasible

import (
	"fmt"
	"os"
	"path"
)

var DATA_RUN_RESULT_FOLDER_NAME = "run_result"
var DATA_APPLICATIONS_FOLDER_NAME = "applications"

func CreateDataFolder(
	paasibleDataFolder string,
) error {
	// ## Create paasible_data folder if not exists
	if _, err := os.Stat(paasibleDataFolder); os.IsNotExist(err) {
		err = os.Mkdir(paasibleDataFolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating paasible_data folder: %w", err)
		}
	}

	// ## Create paasible_data/run_result folder if not exists
	playbookRunResultFolder := path.Join(paasibleDataFolder, DATA_RUN_RESULT_FOLDER_NAME)
	if _, err := os.Stat(playbookRunResultFolder); os.IsNotExist(err) {
		err = os.Mkdir(playbookRunResultFolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating paasible_data/run_result folder: %w", err)
		}
	}

	return nil
}
