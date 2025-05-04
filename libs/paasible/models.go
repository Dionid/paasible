package paasible

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ core.Model = (*PlaybookRunResult)(nil)

type PlaybookRunResult struct {
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
}

func (t PlaybookRunResult) TableName() string {
	return "playbook_run_result"
}

func PlaybookRunResultQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(PlaybookRunResult{}.TableName())
}

func PlaybookRunResultModel(app *pocketbase.PocketBase, entity *PlaybookRunResult) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ALL CODE TILL THE END OF THE FILE IS FOR FUTURE VERSIONS

// ---

var _ core.Model = (*Repository)(nil)

type Repository struct {
	core.BaseModel

	Created  types.DateTime `json:"created" db:"created"`
	Updated  types.DateTime `json:"updated" db:"updated"`
	Remote   string         `json:"remote" db:"remote"`
	TopLevel string         `json:"top_level" db:"top_level"`
}

func (r Repository) TableName() string {
	return "repository"
}

func RepositoryQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(Repository{}.TableName())
}

func RepositoryModel(app *pocketbase.PocketBase, entity *Repository) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*Playbook)(nil)

type Playbook struct {
	core.BaseModel

	Created      types.DateTime `json:"created" db:"created"`
	Updated      types.DateTime `json:"updated" db:"updated"`
	Name         string         `json:"name" db:"name"`
	Path         string         `json:"path" db:"path"`
	Content      string         `json:"content" db:"content"`
	RepositoryId string         `json:"repository_id" db:"repository_id"`
}

func (r Playbook) TableName() string {
	return "playbook"
}

func PlaybookQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(Playbook{}.TableName())
}

func PlaybookModel(app *pocketbase.PocketBase, entity *Playbook) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*Machine)(nil)

type Machine struct {
	core.BaseModel

	Created types.DateTime `json:"created" db:"created"`
	Updated types.DateTime `json:"updated" db:"updated"`
	Name    string         `json:"name" db:"name"`
	Current bool           `json:"current" db:"current"`
}

func (r Machine) TableName() string {
	return "machine"
}

func MachineQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(Machine{}.TableName())
}

func MachineModel(app *pocketbase.PocketBase, entity *Machine) *dbx.ModelQuery {
	return app.DB().Model(entity)
}
