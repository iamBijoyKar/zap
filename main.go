package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/iamBijoyKar/zap/internal/out"
	"github.com/iamBijoyKar/zap/internal/tasks"
	"github.com/iamBijoyKar/zap/internal/utils"
	"gopkg.in/yaml.v3"
)

func run_command(command []string) error {
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

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
		return
	}
	f, err := os.ReadFile(filepath.Join(wd, "zap.yaml"))
	if err != nil {
		log.Fatal("Unable to open zap.yaml file!")
		return
	}
	var yml tasks.Yml
	err = yaml.Unmarshal(f, &yml)
	if err != nil {
		log.Fatal("Error parsing zap.yaml:", err)
		return
	}
	if len(yml.Tasks) == 0 {
		log.Fatal("No tasks found in zap.yaml")
		return
	}

	var completed_tasks []tasks.Task = make([]tasks.Task, len(yml.Tasks))
	var failed_tasks []tasks.Task = make([]tasks.Task, len(yml.Tasks))

	for idx, val := range yml.Tasks {
		// check depends
		if !utils.CheckDeps(completed_tasks, val.Depends_On) {
			out.PrintDefault(fmt.Sprintf("%d. Skipping Task ... %s\n > %s\n  🏗️", idx+1, val.Name, strings.Join(val.Command, " ")))
			out.PrintInfo("Due to dependencies does not matches...\n")
			continue
		}
		out.PrintDefault(fmt.Sprintf("%d. Running Task ... 🔨 %s\n > %s\n", idx+1, color.CyanString(val.Name), strings.Join(val.Command, " ")))
		err := run_command(val.Command)
		if err != nil {
			out.PrintError(fmt.Sprintf("Failed to complete the task! Unable to run the provided command! \n%s\n", strings.Join(val.Command, " ")))
			failed_tasks = append(failed_tasks, val)
			continue
		}
		completed_tasks = append(completed_tasks, val)
		out.PrintDefault(fmt.Sprintf("Task completed ✅\n\n"))
	}
}
