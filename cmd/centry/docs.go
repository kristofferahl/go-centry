package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// GenerateMarkdownCommand is a Command implementation that generates markdown documentation
type GenerateMarkdownCommand struct {
	CLI      *cli.App
	Manifest *config.Manifest
	Log      *logrus.Entry
}

// ToCLICommand returns a CLI command
func (sc *GenerateMarkdownCommand) ToCLICommand() *cli.Command {
	return withCommandDefaults(&cli.Command{
		Name:      "generate-markdown",
		Usage:     "Generate markdown documentation",
		UsageText: "",
		Hidden:    false,
		Action: func(c *cli.Context) error {
			ec := sc.Run(c.Path("file"))
			if ec > 0 {
				return cli.Exit("failed to generate markdown documentation", ec)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Outputs the generated markedown to the specified file",
			},
		},
	})
}

// Run generates markdown documentation
func (sc *GenerateMarkdownCommand) Run(path string) int {
	sc.Log.Debugf("generating markdown documenation")

	md, err := sc.CLI.ToMarkdown()
	if err != nil {
		sc.Log.Error(err)
		return 1
	}

	if path == "" {
		fmt.Print(md)
		sc.Log.Debugf("generated markdown documenation to stdout")
	} else {
		file, err := os.Create(path)
		if err != nil {
			sc.Log.Error(err)
			return 1
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		bc, err := w.WriteString(md)
		if err != nil {
			sc.Log.Error(err)
			return 1
		}
		sc.Log.Tracef("wrote %d bytes", bc)

		err = w.Flush()
		if err != nil {
			sc.Log.Error(err)
			return 1
		}

		sc.Log.Infof("generated markdown documenation to file (path=%s)", path)
	}

	return 0
}
