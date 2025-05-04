package features

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// ALL CODE TILL THE END OF THE FILE IS FOR FUTURE VERSIONS

type LocalFile struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

type LocalFolder struct {
	FolderName string        `json:"folder_name"`
	SubFolders []LocalFolder `json:"sub_folders"`
	Files      []LocalFile   `json:"files"`
}

func getLocalFolderStructure(folderPath string) (LocalFolder, error) {
	folder, err := os.ReadDir(folderPath)
	if err != nil {
		return LocalFolder{}, err
	}

	localFolder := LocalFolder{
		FolderName: folderPath,
		SubFolders: []LocalFolder{},
		Files:      []LocalFile{},
	}

	for _, f := range folder {
		if f.IsDir() {
			subFolder, err := getLocalFolderStructure(folderPath + "/" + f.Name())
			if err != nil {
				return LocalFolder{}, err
			}
			localFolder.SubFolders = append(localFolder.SubFolders, subFolder)
		} else {
			localFile := LocalFile{
				FileName: f.Name(),
				FilePath: folderPath + "/" + f.Name(),
			}
			localFolder.Files = append(localFolder.Files, localFile)
		}
	}

	return localFolder, nil
}

type LocalFolderStructure struct {
	RootFolder LocalFolder `json:"root_folder"`
}

func InitLocalFiles(
	se *core.ServeEvent,
) {
	se.Router.GET("/get-local-folders-structure", func(e *core.RequestEvent) error {
		// read all folders in the public dir
		folderPath := "../../.."

		localFolders, err := getLocalFolderStructure(folderPath)
		if err != nil {
			return e.BadRequestError(err.Error(), nil)
		}

		// create a LocalFolderStructure struct to hold the root folder
		localFolderStructure := LocalFolderStructure{
			RootFolder: localFolders,
		}

		// return the LocalFolderStructure struct as a JSON response
		return e.JSON(200, localFolderStructure)
	})

	// serves static files from the provided public dir (if exists)
	se.Router.GET("/files/{path...}", func(e *core.RequestEvent) error {
		// read folder content
		path := e.Request.PathValue("path")
		fullPath := "../../../" + path

		files, err := os.ReadDir(fullPath)

		// # if it is file than return file content
		if strings.Contains(err.Error(), "not a directory") {
			// read file content
			// check if it's a file instead
			content, fileErr := os.ReadFile(fullPath)
			if fileErr == nil {
				return e.String(200, string(content))
			}

			// if it's not a file, return an error
			return e.BadRequestError("File not found", nil)
		}

		if err != nil {
			return e.BadRequestError(err.Error(), nil)
		}

		// create a string slice to hold the file names
		var fileNames []string
		// iterate over the files and append their names to the slice
		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}

		// join the file names into a single string
		fileList := ""
		for _, file := range fileNames {
			fileList += file + "\n"
		}

		// return the file names as a response
		return e.String(200, fileList)
	})
}
