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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Enable verbose output",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Run all tasks",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "verbose",
								Usage: "Enable verbose output",
							},
						},
						Action: func(c *cli.Context) error {
							verbose := c.Bool("verbose")
							return tasks.RunTask("", verbose)
						},
					},
					{
						Name:      "task",
						Usage:     "Run a specific task",
						ArgsUsage: "<task-name>",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "verbose",
								Usage: "Enable verbose output",
							},
						},
						Action: func(c *cli.Context) error {
							if c.Args().Len() == 0 {
								return fmt.Errorf("task name is required")
							}
							taskName := c.Args().First()
							verbose := c.Bool("verbose")
							return tasks.RunTask(taskName, verbose)
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
