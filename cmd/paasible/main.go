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
	currentFolder, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// # Take inital env
	envFilename, ok := os.LookupEnv("PAASIBLE_ENV_FILENAME")
	if ok == false {
		envFilename = "paasible"
	}

	envPath, ok := os.LookupEnv("PAASIBLE_ENV_PATH")
	if ok == false {
		envPath = "."
	}

	yamlFilename, ok := os.LookupEnv("PAASIBLE_YAML_FILENAME")
	if ok == false {
		yamlFilename = "paasible"
	}

	yamlPath, ok := os.LookupEnv("PAASIBLE_YAML_PATH")
	if ok == false {
		yamlPath = "."
	}

	envConfig, yamlConfig, err := initConfig(
		envFilename,
		envPath,
		yamlFilename,
		yamlPath,
	)
	if err != nil {
		log.Fatal(err)
	}

	// # Paasible specific Config
	paasibleConfig := paasible.CliConfig{
		Machine:                envConfig.Machine,
		User:                   envConfig.User,
		CliVersion:             yamlConfig.Paasible.CliVersion,
		ProjectName:            yamlConfig.Paasible.ProjectName,
		DataFolderRelativePath: yamlConfig.Paasible.DataFolderRelativePath,
	}

	// # Paasible data folder path
	paasibleDataFolderPath := path.Join(
		currentFolder,
		paasibleConfig.DataFolderRelativePath,
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
		DefaultDataDir: path.Join(paasibleDataFolderPath, "/pb_data"),
	})

	// # Migrations
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
		Dir:         path.Join(currentFolder, "pb_migrations"),
	})

	// # Commands
	features.InitAnsiblePlaybookCmd(
		app,
		&paasibleConfig,
		paasibleDataFolderPath,
	)

	features.InitRunAppilcationCmd(
		app,
		&paasibleConfig,
		paasibleDataFolderPath,
	)

	features.InitInitCmd(
		app,
		envFilename,
		envPath,
		yamlFilename,
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
			path.Join(currentFolder, yamlConfig.Paasible.DataFolderRelativePath),
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

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Can't create watcher", err)
	}
	defer watcher.Close()

	//
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
		path.Join(currentFolder, paasibleConfig.DataFolderRelativePath, paasible.DATA_RUN_RESULT_FOLDER_NAME),
	); err != nil {
		log.Fatal("Can't add file watcher", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal("App start: ", err)
	}
}
