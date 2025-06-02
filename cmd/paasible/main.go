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
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/spf13/pflag"
)

func main() {
	paasibleCliPwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// # Paasible root config path
	// ## Take from env
	paasibleRootConfigPathFromEnv, ok := os.LookupEnv("PAASIBLE_CONFIG_PATH")
	if ok == false {
		paasibleRootConfigPathFromEnv = "./paasible.yaml"
	}

	// ## Take from command line flags
	paasibleRootConfigPathP := pflag.StringP("config", "c", paasibleRootConfigPathFromEnv, "Path to paasible.yaml config file")
	pflag.Parse()
	paasibleRootConfigPath := *paasibleRootConfigPathP

	// ## Parse paasible config file
	paasibleRootConfig, err := newConfigFromPath(
		paasibleRootConfigPath,
	)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: decide which is better
	// # Paasible config folder path
	// paasibleRootConfigFolderPath := path.Join(
	// 	paasibleCliPwd,
	// 	path.Dir(
	// 		paasibleRootConfigPath,
	// 	),
	// )
	paasibleRootConfigFolderPath := path.Dir(
		paasibleRootConfigPath,
	)

	// # Init storage from paasible config file
	storage, err := paasible.EntityStorageFromOrigin(
		paasibleRootConfig,
	)
	if err != nil {
		log.Fatal(err)
	}

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
		paasibleRootConfigPath,
	)

	// # Send verification email on sign-up
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		return mails.SendRecordVerification(app, e.Record)
	})

	// # Parse paasible data folder
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// # Parse paasible data folder
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
						app.Logger().Error("Upsert paasible run result", "error", err)
					}

					app.Logger().Debug("Upserted new json file")
				}

			// watch for errors
			case err := <-watcher.Errors:
				{
					app.Logger().Error("Error listening files", "error", err)
				}
			}
		}
	}()

	// ## out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(
		path.Join(
			paasibleDataFolderPath,
			paasible.RUN_RESULTS_FOLDER_NAME,
		),
	); err != nil {
		log.Fatal("Can't add file watcher", err)
	}

	// # Map ui command
	uiCmd := cmd.NewServeCommand(
		app,
		true,
	)
	uiCmd.Use = "ui [domain(s)]"
	app.RootCmd.AddCommand(uiCmd)

	// # Add superuser cmd
	app.RootCmd.AddCommand(cmd.NewSuperuserCommand(app))

	// # Start app
	if err := app.Execute(); err != nil {
		log.Fatal("App start error: ", err)
	}
}
