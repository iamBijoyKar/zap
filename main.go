package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/iamBijoyKar/zap/internal/out"
	"github.com/iamBijoyKar/zap/internal/tasks"
	"github.com/iamBijoyKar/zap/internal/utils"
	"gopkg.in/yaml.v3"
)

var Version string = "1.0.0"

func main() {
	out.PrintDefault((fmt.Sprintf("\n\t⚡ %s %s (golang)\n", color.YellowString("Zap"), color.YellowString(Version))))
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

	fmt.Printf("\t - Total Tasks: %d\n\n", len(yml.Tasks))

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
		err := tasks.RunCommand(val.Command)
		if err != nil {
			out.PrintError(fmt.Sprintf("Failed to complete the task! Unable to run the provided command! \n%s\n", strings.Join(val.Command, " ")))
			failed_tasks = append(failed_tasks, val)
			continue
		}
		completed_tasks = append(completed_tasks, val)
		out.PrintDefault(fmt.Sprintf("Task completed ✅\n\n"))
	}
	out.PrintDefault(fmt.Sprintf("Total Completed Tasks: %d\n", len(completed_tasks)))
	out.PrintDefault(fmt.Sprintf("Total Failed Tasks: %d\n", len(failed_tasks)))
}
