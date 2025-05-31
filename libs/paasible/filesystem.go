package paasible

import (
	"fmt"
	"os"
	"path"
)

var RUN_RESULTS_FOLDER_NAME = "run_results"

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

	// ## Create paasible_data/run_results folder if not exists
	playbookRunResultFolder := path.Join(paasibleDataFolder, RUN_RESULTS_FOLDER_NAME)
	if _, err := os.Stat(playbookRunResultFolder); os.IsNotExist(err) {
		err = os.Mkdir(playbookRunResultFolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating paasible_data/run_results folder: %w", err)
		}
	}

	return nil
}
