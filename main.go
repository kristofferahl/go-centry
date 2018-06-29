package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kristofferahl/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log := logrus.New()

	c := &cli.CLI{
		Name:     "go-centry",
		Version:  "1.0.0",
		Args:     os.Args,
		Commands: map[string]cli.CommandFactory{},

		Autocomplete:          true,
		AutocompleteInstall:   "install-autocomplete",
		AutocompleteUninstall: "uninstall-autocomplete",
	}

	if len(os.Args) >= 2 {
		c.Args = os.Args[2:]

		manifest := loadManifest(os.Args[1])

		l, _ := logrus.ParseLevel(manifest.Config.Log.Level)
		log.SetLevel(l)
		log.Formatter = &PrefixedTextFormatter{
			Prefix: manifest.Config.Log.Prefix,
		}

		for _, cmd := range manifest.Commands {
			cmd := cmd
			command := &DynamicCommand{
				Log: log.WithFields(logrus.Fields{
					"command": cmd.Name,
				}),
				Command:  cmd,
				Manifest: manifest,
			}

			bashCommands := command.GeBashCommands()

			for _, bc := range bashCommands {
				bc := bc
				cmdKey := strings.Replace(bc, ":", " ", -1)
				log.Debug(fmt.Sprintf("Adding command %s", cmdKey))

				// TODO: Keep this for a while
				// if i == 1 && bc != cmd.Name {
				// 	c.Commands[cmd.Name] = func() (cli.Command, error) {
				// 		return &BashCommand{
				// 			Manifest: manifest,
				// 			Log: log.WithFields(logrus.Fields{
				// 				"command": cmd.Name,
				// 			}),
				// 			Name:         cmd.Name,
				// 			Path:         cmd.Path,
				// 			HelpText:     cmd.Help,
				// 			SynopsisText: cmd.Synopsis,
				// 		}, nil
				// 	}
				// }

				help := cmd.Help
				synopsis := cmd.Synopsis

				if bc != cmd.Name {
					help = fmt.Sprintf("TODO: Help: %s", cmdKey)
					synopsis = fmt.Sprintf("TODO: Synopsis: %s", cmdKey)
				}

				c.Commands[cmdKey] = func() (cli.Command, error) {
					return &BashCommand{
						Manifest: manifest,
						Log: log.WithFields(logrus.Fields{
							"command": cmdKey,
						}),
						Name:         bc,
						Path:         cmd.Path,
						HelpText:     help,
						SynopsisText: synopsis,
					}, nil
				}
			}
		}
	}

	log.Debug(fmt.Sprintf("Args passed on from centry: %v", c.Args))

	exitStatus, err := c.Run()
	if err != nil {
		log.Error(err)
	}

	return exitStatus
}
