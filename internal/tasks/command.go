package tasks

import (
	"errors"
	"os/exec"

	"github.com/iamBijoyKar/zap/internal/out"
)

func RunCommand(command []string) error {
	if len(command) == 0 {
		out.PrintError("No command specified!\n")
		return errors.New("no command found in the task")
	}
	cmd := exec.Command(command[0], command[1:]...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("unable to the command")
	}
	return nil
}
