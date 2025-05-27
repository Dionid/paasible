package paasible

import (
	"bytes"
	"os/exec"
	"strings"
)

type ExecGitCommandResult struct {
	Stdout string
	Stderr string
	Error  error
}

func ExecGitCommand(
	cmd *exec.Cmd,
	currentFolderPath string,
) ExecGitCommandResult {
	cmd.Dir = currentFolderPath
	var cmdStderr bytes.Buffer
	cmd.Stderr = &cmdStderr

	revParseCmdOutput, err := cmd.Output()
	revParseCmdOutputStr := strings.TrimSpace(string(revParseCmdOutput))
	stdErrString := cmdStderr.String()

	return ExecGitCommandResult{
		Stdout: revParseCmdOutputStr,
		Stderr: stdErrString,
		Error:  err,
	}
}
