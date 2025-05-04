package paasible

import (
	"os"
	"strings"

	"github.com/Dionid/paasible/libs/ansible"
	"github.com/Dionid/paasible/libs/uuidv7"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"gopkg.in/yaml.v3"
)

// ALL CODE TILL THE END OF THE FILE IS FOR FUTURE VERSIONS
func GetAndSyncOrCreatePlaybook(
	app *pocketbase.PocketBase,
	currentFolderPath string,
	playbookRelativePath string,
	repository *Repository,
) (*Playbook, error) {
	// # Parse file to paasible.AnsiblePlaybook
	// ## Check if file exists
	if _, err := os.Stat(playbookRelativePath); os.IsNotExist(err) {
		return nil, err
	}

	// ## Parse file
	ansiblePlaybook := ansible.Playbook{}
	ansiblePlaybookFile, err := os.ReadFile(playbookRelativePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(ansiblePlaybookFile, &ansiblePlaybook); err != nil {
		return nil, err
	}

	// ## Get file playbook path
	playbookAbsolutePath := currentFolderPath + "/" + playbookRelativePath
	playbookPath := playbookAbsolutePath
	if repository.Id != "" {
		// ## substring repository path from absolute path
		playbookPath = strings.Replace(playbookAbsolutePath, repository.TopLevel, "", 1)
	}

	playbookName := ""

	for _, item := range ansiblePlaybook {
		playbookName += item.Name + "; "
	}

	// ## Find or create playbook
	playbook := Playbook{}

	err = PlaybookQuery(app).Where(dbx.HashExp{
		"path": playbookPath,
	}).OrWhere(
		dbx.NewExp(
			"UPPER(name) LIKE '%{:name}%'",
			dbx.Params{
				"name": playbookName,
			},
		),
	).One(&playbook)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return nil, err
		}

		// ## Create new playbook
		playbook = Playbook{
			BaseModel: core.BaseModel{
				Id: uuidv7.NewE().String(),
			},
			Name:         playbookName,
			Created:      types.NowDateTime(),
			Updated:      types.NowDateTime(),
			Content:      string(ansiblePlaybookFile),
			Path:         playbookPath,
			RepositoryId: repository.Id,
		}

		// ## Insert playbook
		err = PlaybookModel(app, &playbook).Insert()
		if err != nil {
			return nil, err
		}
	}

	return &playbook, nil
}
