package tasks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/iamBijoyKar/zap/internal/out"
	"gopkg.in/yaml.v3"
)

func CheckDeps(completed_tasks []Task, deps_on []string) bool {
	dep_checks := make([]bool, len(deps_on))
	for idx, dep := range deps_on {
		for _, comp := range completed_tasks {
			if dep == comp.Name {
				dep_checks[idx] = true
				break
			}
		}
	}
	result := true
	for _, val := range dep_checks {
		result = result && val
	}
	return result
}

func RunCommand(command []string) error {
	if len(command) == 0 {
		out.PrintError("No command specified!\n")
		return errors.New("no command found in the task")
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("unable to the command")
	}
	return nil
}

func RunTask(taskName string) error {
	out.PrintDefault((fmt.Sprintf("\n\t⚡ %s %s (golang)\n", color.YellowString("Zap"), color.YellowString("1.0.0"))))

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	f, err := os.ReadFile(filepath.Join(wd, "zap.yaml"))
	if err != nil {
		return fmt.Errorf("unable to open zap.yaml file: %v", err)
	}

	var yml Yml
	err = yaml.Unmarshal(f, &yml)
	if err != nil {
		return fmt.Errorf("error parsing zap.yaml: %v", err)
	}

	if len(yml.Tasks) == 0 {
		return fmt.Errorf("no tasks found in zap.yaml")
	}

	// If a specific task name is provided, find and run only that task
	if taskName != "" {
		return runSpecificTask(yml.Tasks, taskName)
	}

	out.PrintDefault(fmt.Sprintf("\t - Total Tasks: %d\n\n", len(yml.Tasks)))

	var completed_tasks []Task = make([]Task, len(yml.Tasks))
	var failed_tasks []Task = make([]Task, len(yml.Tasks))

	for idx, val := range yml.Tasks {
		// check depends
		if !CheckDeps(completed_tasks, val.Depends_On) {
			out.PrintDefault(fmt.Sprintf("%d. Skipping Task ... %s\n > %s\n  🏗️", idx+1, val.Name, strings.Join(val.Command, " ")))
			out.PrintInfo("Due to dependencies does not matches...\n")
			continue
		}
		out.PrintDefault(fmt.Sprintf("%d. Running Task ... 🔨 %s\n > %s\n", idx+1, color.CyanString(val.Name), strings.Join(val.Command, " ")))
		err := RunCommand(val.Command)
		if err != nil {
			out.PrintError(fmt.Sprintf("Failed to complete the task! Unable to run the provided command! \n%s\n", strings.Join(val.Command, " ")))
			failed_tasks = append(failed_tasks, val)
			continue
		}
		completed_tasks = append(completed_tasks, val)
		out.PrintDefault("Task completed ✅\n\n")
	}
	out.PrintDefault(fmt.Sprintf("Total Completed Tasks: %d\n", len(completed_tasks)))
	out.PrintDefault(fmt.Sprintf("Total Failed Tasks: %d\n", len(failed_tasks)))
	return nil
}

func runSpecificTask(tasks []Task, taskName string) error {
	// Find the specific task
	var targetTask *Task
	for i, task := range tasks {
		if task.Name == taskName {
			targetTask = &tasks[i]
			break
		}
	}

	if targetTask == nil {
		return fmt.Errorf("task '%s' not found in zap.yaml", taskName)
	}

	out.PrintDefault(fmt.Sprintf("Running specific task: %s\n", color.CyanString(taskName)))

	// Check if the task has dependencies and if they are satisfied
	// For now, we'll run the task directly without dependency checking for specific tasks
	// This could be enhanced later to check dependencies
	out.PrintDefault(fmt.Sprintf("Running Task ... 🔨 %s\n > %s\n", color.CyanString(targetTask.Name), strings.Join(targetTask.Command, " ")))

	err := RunCommand(targetTask.Command)
	if err != nil {
		out.PrintError(fmt.Sprintf("Failed to complete the task! Unable to run the provided command! \n%s\n", strings.Join(targetTask.Command, " ")))
		return fmt.Errorf("task '%s' failed: %v", taskName, err)
	}

	out.PrintDefault("Task completed ✅\n")
	return nil
}
