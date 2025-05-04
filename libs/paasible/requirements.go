package paasible

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func CheckRequirements() error {
	command := exec.Command("ansible-playbook", "--version")

	// Run the command
	err := command.Run()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("You must install ansible-playbook executable or add it to $PATH")
		}

		return fmt.Errorf("Error running ansible-playbook exec: %w", err)
	}

	return nil
}
