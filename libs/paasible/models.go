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

// ---

var _ core.Model = (*Playbook)(nil)

type Playbook struct {
	core.BaseModel

	Name             string `json:"name" db:"name"`
	Path             string `json:"path" db:"path"`
	PlaybookPath     string `json:"playbook_path" db:"playbook_path"`
	OriginRepository string `json:"origin_repository" db:"origin_repository"`
}

func (r Playbook) TableName() string {
	return "playbook"
}

func PlaybookQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(Playbook{}.TableName())
}

func PlaybookModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*PlaybookTarget)(nil)

type PlaybookTarget struct {
	core.BaseModel

	Created    types.DateTime `json:"created" db:"created"`
	Updated    types.DateTime `json:"updated" db:"updated"`
	Name       string         `json:"name" db:"name"`
	User       string         `json:"user" db:"user"`
	PlaybookId string         `json:"playbook_id" db:"playbook_id"`
	TargetId   string         `json:"target_id" db:"target_id"`
}

func (r PlaybookTarget) TableName() string {
	return "playbook_target"
}

func PlaybookTargetQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(PlaybookTarget{}.TableName())
}

func PlaybookTargetModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}
