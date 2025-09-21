package tasks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

func RunCommand(command []string, verbose bool) error {
	if len(command) == 0 {
		out.PrintError("No command specified!\n")
		return errors.New("no command found in the task")
	}

	if verbose {
		out.PrintInfo(fmt.Sprintf("Executing command: %s\n", strings.Join(command, " ")))
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if verbose {
			out.PrintError(fmt.Sprintf("Command failed with error: %v\n", err))
		}
		return errors.New("unable to run the command")
	}
	return nil
}

func RunTaskWithRetries(task Task, verbose bool) (error, time.Duration) {
	startTime := time.Now()
	maxRetries := task.Retries
	if maxRetries <= 0 {
		maxRetries = 1
	} else {
		maxRetries = maxRetries + 1 // Add 1 because retries is the number of retries, not total attempts
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if verbose && attempt > 1 {
			out.PrintInfo(fmt.Sprintf("Retry attempt %d/%d for task '%s'\n", attempt, maxRetries, task.Name))
		}

		err := RunCommand(task.Command, verbose)
		if err == nil {
			if verbose && attempt > 1 {
				out.PrintInfo(fmt.Sprintf("Task '%s' succeeded on attempt %d\n", task.Name, attempt))
			}
			duration := time.Since(startTime)
			return nil, duration
		}

		lastErr = err
		if attempt < maxRetries {
			if verbose {
				out.PrintInfo(fmt.Sprintf("Task '%s' failed on attempt %d, retrying...\n", task.Name, attempt))
			}
			time.Sleep(time.Second * 2) // Wait 2 seconds before retry
		}
	}

	duration := time.Since(startTime)
	return fmt.Errorf("task '%s' failed after %d attempts: %v", task.Name, maxRetries, lastErr), duration
}

func RunTask(taskName string, verbose bool) error {
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
		return runSpecificTask(yml.Tasks, taskName, verbose)
	}

	out.PrintDefault(fmt.Sprintf("\t - Total Tasks: %d\n\n", len(yml.Tasks)))

	var completed_tasks []Task
	var failed_tasks []Task
	var mutex sync.Mutex

	// Separate parallel and sequential tasks
	var parallelTasks []Task
	var sequentialTasks []Task

	for _, task := range yml.Tasks {
		if task.Parallel {
			parallelTasks = append(parallelTasks, task)
		} else {
			sequentialTasks = append(sequentialTasks, task)
		}
	}

	// Run parallel tasks concurrently
	if len(parallelTasks) > 0 {
		if verbose {
			out.PrintInfo(fmt.Sprintf("Running %d parallel tasks concurrently...\n", len(parallelTasks)))
		}

		var wg sync.WaitGroup
		for _, task := range parallelTasks {
			wg.Add(1)
			go func(t Task) {
				defer wg.Done()

				if verbose {
					out.PrintDefault(fmt.Sprintf("Running parallel task: 🔨 %s\n > %s\n", color.CyanString(t.Name), strings.Join(t.Command, " ")))
				}

				err, duration := RunTaskWithRetries(t, verbose)

				mutex.Lock()
				if err != nil {
					out.PrintError(fmt.Sprintf("Parallel task '%s' failed: %v (took %v)\n", t.Name, err, duration.Round(time.Millisecond)))
					failed_tasks = append(failed_tasks, t)
				} else {
					completed_tasks = append(completed_tasks, t)
					out.PrintDefault(fmt.Sprintf("Parallel task '%s' completed ✅ (took %v)\n", t.Name, duration.Round(time.Millisecond)))
				}
				mutex.Unlock()
			}(task)
		}
		wg.Wait()
	}

	// Run sequential tasks
	for idx, task := range sequentialTasks {
		// Check dependencies
		if !CheckDeps(completed_tasks, task.Depends_On) {
			out.PrintDefault(fmt.Sprintf("%d. Skipping Task ... %s\n > %s\n  🏗️", idx+1, task.Name, strings.Join(task.Command, " ")))
			out.PrintInfo("Due to dependencies does not match...\n")
			continue
		}

		out.PrintDefault(fmt.Sprintf("%d. Running Task ... 🔨 %s\n > %s\n", idx+1, color.CyanString(task.Name), strings.Join(task.Command, " ")))

		err, duration := RunTaskWithRetries(task, verbose)
		if err != nil {
			out.PrintError(fmt.Sprintf("Failed to complete the task! %v (took %v)\n", err, duration.Round(time.Millisecond)))
			failed_tasks = append(failed_tasks, task)
			continue
		}

		completed_tasks = append(completed_tasks, task)
		out.PrintDefault(fmt.Sprintf("Task completed ✅ (took %v)\n\n", duration.Round(time.Millisecond)))
	}

	out.PrintDefault(fmt.Sprintf("Total Completed Tasks: %d\n", len(completed_tasks)))
	out.PrintDefault(fmt.Sprintf("Total Failed Tasks: %d\n", len(failed_tasks)))
	return nil
}

func runSpecificTask(tasks []Task, taskName string, verbose bool) error {
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

	err, duration := RunTaskWithRetries(*targetTask, verbose)
	if err != nil {
		out.PrintError(fmt.Sprintf("Failed to complete the task! %v (took %v)\n", err, duration.Round(time.Millisecond)))
		return fmt.Errorf("task '%s' failed: %v", taskName, err)
	}

	out.PrintDefault(fmt.Sprintf("Task completed ✅ (took %v)\n", duration.Round(time.Millisecond)))
	return nil
}
