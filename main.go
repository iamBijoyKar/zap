package main

import (
	"fmt"
	"os"

	"github.com/iamBijoyKar/zap/internal/tasks"
	"github.com/urfave/cli/v2"
)

var Version string = "1.0.0"

func main() {
	app := &cli.App{
		Name:    "zap",
		Usage:   "A task runner for Go projects",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run tasks from zap.yaml",
				Subcommands: []*cli.Command{
					{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Run all tasks",
						Action: func(c *cli.Context) error {
							return tasks.RunTask("")
						},
					},
					{
						Name:      "task",
						Usage:     "Run a specific task",
						ArgsUsage: "<task-name>",
						Action: func(c *cli.Context) error {
							if c.Args().Len() == 0 {
								return fmt.Errorf("task name is required")
							}
							taskName := c.Args().First()
							return tasks.RunTask(taskName)
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
