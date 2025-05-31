package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/Dionid/paasible/cmd/paasible/features"
	_ "github.com/Dionid/paasible/cmd/paasible/pb_migrations"
	"github.com/Dionid/paasible/libs/paasible"
	"github.com/fsnotify/fsnotify"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	paasibleCliPwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	yamlPath, ok := os.LookupEnv("PAASIBLE_CONFIG_PATH")
	if ok == false {
		yamlPath = path.Join(
			paasibleCliPwd,
			"./paasible.yaml",
		)
	}

	yamlConfig, err := initConfig(
		yamlPath,
	)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := paasible.EntityStorageFromOrigin(
		yamlConfig,
	)
	if err != nil {
		log.Fatal(err)
	}

	// # Paasible config folder path
	paasibleRootConfigFolderPath := path.Dir(
		yamlPath,
	)

	// # Paasible data folder path
	paasibleDataFolderPath := path.Join(
		paasibleRootConfigFolderPath,
		storage.Paasible.DataFolderPath,
	)

	// # Create data folder
	err = paasible.CreateDataFolder(
		paasibleDataFolderPath,
	)
	if err != nil {
		log.Fatal("Can't create paasible directory: ", err)
	}

	// # Pocketbase
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: path.Join(paasibleDataFolderPath, "db"),
	})

	// # Migrations
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
		Dir:         path.Join(paasibleCliPwd, "pb_migrations"),
	})

	// # Commands
	features.InitAnsiblePlaybookCmd(
		app,
		storage,
		paasibleCliPwd,
		paasibleDataFolderPath,
	)

	features.InitRunPlaybookCmd(
		app,
		storage,
		paasibleRootConfigFolderPath,
		paasibleCliPwd,
		paasibleDataFolderPath,
	)

	features.InitInitCmd(
		app,
		yamlPath,
	)

	// # Send verification email on sign-up
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		return mails.SendRecordVerification(app, e.Record)
	})

	// # Parse paasible data folder
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		err := features.ParsePaasibleData(
			app,
			path.Join(paasibleCliPwd, storage.Paasible.DataFolderPath),
		)

		if err != nil {
			return err
		}

		return se.Next()
	})

	// # Check requiremenets
	err = paasible.CheckRequirements()
	if err != nil {
		log.Fatal("Check requiremenets: ", err)
	}

	// # When new log file added sabe it to DB
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Can't create watcher", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				{
					err := features.UpsertPlaybookRunResult(
						app,
						event.Name,
					)

					if err != nil {
						app.Logger().Error("Upsert paasible run result", err)
					}

					app.Logger().Debug("Upserted new json file")
				}

			// watch for errors
			case err := <-watcher.Errors:
				{
					app.Logger().Error("Error listening files", err)
				}
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(
		path.Join(
			paasibleDataFolderPath,
			paasible.RUN_RESULT_FOLDER_NAME,
		),
	); err != nil {
		log.Fatal("Can't add file watcher", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal("App start: ", err)
	}
}
