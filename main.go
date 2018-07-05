package main

import (
	"flag"
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
		Commands: map[string]cli.CommandFactory{},

		Autocomplete:          true,
		AutocompleteInstall:   "install-autocomplete",
		AutocompleteUninstall: "uninstall-autocomplete",
	}

	// Parse global flags
	centryFlags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	file := centryFlags.String("file", "./centry.yaml", "The path to the manifest file")
	centryFlags.StringVar(file, "f", "./centry.yaml", "The path to the manifest file")
	centryFlags.Parse(os.Args[1:2])

	// Load manifest
	manifest := loadManifest(*file)

	// Configure logger
	l, _ := logrus.ParseLevel(manifest.Config.Log.Level)
	log.SetLevel(l)
	log.Formatter = &PrefixedTextFormatter{
		Prefix: manifest.Config.Log.Prefix,
	}

	// Add global option flags
	globalFlags := flag.NewFlagSet("global", flag.ExitOnError)
	for _, opt := range manifest.Options {
		opt := opt
		log.Debugf("Adding global option %s (default value: %s)", opt.Name, opt.Default)
		switch opt.Default {
		case "":
			globalFlags.Bool(opt.Name, false, opt.Description)
		default:
			globalFlags.String(opt.Name, opt.Default, opt.Description)
		}
	}

	// Re-parse global flags
	globalFlags.Parse(os.Args[2:])

	// Set args for cli
	c.Args = globalFlags.Args()

	// Build commands
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
			cmdName := bc
			cmdKey := strings.Replace(cmdName, ":", " ", -1)
			log.Debugf("Adding command %s", cmdKey)

			// TODO: Keep this for a while
			// if i == 1 && cmdName != cmd.Name {
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

			c.Commands[cmdKey] = func() (cli.Command, error) {
				return &BashCommand{
					Manifest: manifest,
					Log: log.WithFields(logrus.Fields{
						"command": cmdKey,
					}),
					Name:         cmdName,
					Path:         cmd.Path,
					HelpText:     cmd.Help,
					SynopsisText: cmd.Description,
				}, nil
			}
		}
	}

	// Run cli
	exitStatus, err := c.Run()
	if err != nil {
		log.Error(err)
	}

	return exitStatus
}
