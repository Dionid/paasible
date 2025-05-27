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

// ---

var _ core.Model = (*SshKey)(nil)

type SshKey struct {
	core.BaseModel

	Created types.DateTime `json:"created" db:"created"`
	Updated types.DateTime `json:"updated" db:"updated"`

	Name    string `json:"name" db:"name"`
	Public  string `json:"public" db:"public"`
	Private string `json:"private" db:"private"`
}

func (r SshKey) TableName() string {
	return "ssh_key"
}

func SshKeyQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(SshKey{}.TableName())
}

func SshKeyModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*Target)(nil)

type Target struct {
	core.BaseModel

	Created types.DateTime `json:"created" db:"created"`
	Updated types.DateTime `json:"updated" db:"updated"`

	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	Type    string `json:"type" db:"type"` // e.g. "ssh", "docker", etc.
}

func (r Target) TableName() string {
	return "target"
}

func TargetQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(Target{}.TableName())
}

func TargetModel(app *pocketbase.PocketBase, entity *Target) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*TargetSshKey)(nil)

type TargetSshKey struct {
	core.BaseModel

	Created types.DateTime `json:"created" db:"created"`
	Updated types.DateTime `json:"updated" db:"updated"`

	SshKeyId string `json:"ssh_key_id" db:"ssh_key_id"`
	TargetId string `json:"target_id" db:"target_id"`
}

func (r TargetSshKey) TableName() string {
	return "target_ssh_key"
}

func TargetSshKeyQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(TargetSshKey{}.TableName())
}

func TargetSshKeyModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}
