package main

import (
	"fmt"
	"os"

	"github.com/iamBijoyKar/zap/internal/tasks"
	"github.com/urfave/cli/v2"
)

var Version string = "v1.0.0"

var Yellow string = "\033[33m"
var Cyan string = "\033[0;36m"
var BoldWhite string = "\033[1;37m"
var Reset string = "\033[0m"

// Helper function to repeat a string
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func main() {
	app := &cli.App{
		Name:      "Zap",
		UsageText: "Cli app to run tasks",
		Version:   Version,
		Action: func(ctx *cli.Context) error {
			text := fmt.Sprintf(" %s#%s Welcome to ⚡ Zap Toolkit", Cyan, Reset)
			border := "   ┌" + repeat("─", len(text)-9) + "┐\n"
			border += "   │ " + text + "  │\n"
			border += "   └" + repeat("─", len(text)-9) + "┘\n"
			fmt.Print(border)
			fmt.Printf("\n\t%s███████╗ █████╗ ██████╗\n\t╚══███╔╝██╔══██╗██╔══██╗\n\t  ███╔╝ ███████║██████╔╝\n\t ███╔╝  ██╔══██║██╔═══╝\n\t███████╗██║  ██║██║\n\t╚══════╝╚═╝  ╚═╝╚═╝%s\n\n", Yellow, Reset)
			fmt.Printf("⚡Zap - Run your tasks sequentialy or parallely\n\n\tUse %sZap --help%s for more info\n", BoldWhite, Reset)
			return nil
		},
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
