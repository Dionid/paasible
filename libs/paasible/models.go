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

	TargetId      string `json:"target_id" db:"target_id"`           // the target where the command was run
	ApplicationId string `json:"application_id" db:"application_id"` // the application where the command was run
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

var _ core.Model = (*ApplicationTemplate)(nil)

type ApplicationTemplate struct {
	core.BaseModel

	Name       string `json:"name" db:"name"`
	Repository string `json:"repository" db:"repository"`
}

func (r ApplicationTemplate) TableName() string {
	return "application_template"
}

func ApplicationTemplateQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(ApplicationTemplate{}.TableName())
}

func ApplicationTemplateModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*Application)(nil)

type Application struct {
	core.BaseModel

	Name       string `json:"name" db:"name"`
	TemplateId string `json:"template_id" db:"template_id"`
	Path       string `json:"path" db:"path"`
}

func (r Application) TableName() string {
	return "application"
}

func ApplicationQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(Application{}.TableName())
}

func ApplicationModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}

// ---

var _ core.Model = (*ApplicationTarget)(nil)

type ApplicationTarget struct {
	core.BaseModel

	Created       types.DateTime `json:"created" db:"created"`
	Updated       types.DateTime `json:"updated" db:"updated"`
	ApplicationId string         `json:"application_id" db:"application_id"`
	TargetId      string         `json:"target_id" db:"target_id"`
}

func (r ApplicationTarget) TableName() string {
	return "application_target"
}

func ApplicationTargetQuery(app *pocketbase.PocketBase) *dbx.SelectQuery {
	return app.RecordQuery(ApplicationTarget{}.TableName())
}

func ApplicationTargetModel(app *pocketbase.PocketBase, entity *SshKey) *dbx.ModelQuery {
	return app.DB().Model(entity)
}
