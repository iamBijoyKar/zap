package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iamBijoyKar/zap/internal/out"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Name       string   `yaml:"name"`
	Command    []string `yaml:"command"`
	Retries    int      `yaml:"retries"`
	Parallel   bool     `yaml:"parallel"`
	Depends_On []string `yaml:"depends_on"`
}
type Yml struct {
	Tasks []Task `yaml:"tasks"`
}

func run_command(command []string) error {
	if len(command) == 0 {
		out.PrintError("No command specified!\n")
		return errors.New("no command found in the task")
	}
	cmd := exec.Command("./", command...)
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
	var yml Yml
	err = yaml.Unmarshal(f, &yml)
	if err != nil {
		log.Fatal("Error parsing zap.yaml:", err)
		return
	}
	if len(yml.Tasks) == 0 {
		log.Fatal("No tasks found in zap.yaml")
		return
	}

	failed_flag := false

	for idx, val := range yml.Tasks {
		if failed_flag {
			out.PrintDefault(fmt.Sprintf("%d. Skipping Task %s\n", idx+1, val.Name))
			continue
		}
		out.PrintDefault(fmt.Sprintf("%d. Running Task %s\n", idx+1, val.Name))
		err := run_command(val.Command)
		if err != nil {
			out.PrintError(fmt.Sprintf("Failed to complete the task! Unable to run the provided command! \n%s\n", strings.Join(val.Command, " ")))
			failed_flag = true
		}

	}
}
