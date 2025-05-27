package paasible

type CliConfig struct {
	User    string
	Machine string

	CliVersion             string
	DataFolderRelativePath string
}

func CliConfigYaml() string {
	return `paasible:
  cli_version: 0.0.1
`
}

func CliConfigEnv() string {
	return `PAASIBLE_UI_PORT=8080
PAASIBLE_USER=user
PAASIBLE_MACHINE=local
`
}
