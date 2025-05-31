package paasible

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ core.Model = (*RunResult)(nil)

type RunResult struct {
	core.BaseModel

	Created          types.DateTime `json:"created" db:"created"`
	Updated          types.DateTime `json:"updated" db:"updated"`
	Command          string         `json:"command" db:"command"`
	Stdout           string         `json:"stdout" db:"stdout"`
	Stderr           string         `json:"stderr" db:"stderr"`
	Error            string         `json:"error" db:"error"`
	Pwd              string         `json:"pwd" db:"pwd"` // runned from this directory
	RepositoryBranch string         `json:"repository_branch" db:"repository_branch"`
	MachineId        string         `json:"machine_id" db:"machine_id"`
	UserId           string         `json:"user_id" db:"user_id"`

	PlaybookId string `json:"playbook_id" db:"playbook_id"` // the playbook where the command was run
}

func (t RunResult) TableName() string {
	return "run_result"
}

func PlaybookRunResultQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(RunResult{}.TableName())
}

func PlaybookRunResultModel(app *pocketbase.PocketBase, entity *RunResult) *dbx.ModelQuery {
	return app.DB().Model(entity)
}
