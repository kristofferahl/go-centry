package centry

import (
	"strings"

	"github.com/kristofferahl/cli"
	"github.com/kristofferahl/go-centry/pkg/config"
	"github.com/kristofferahl/go-centry/pkg/logger"
	"github.com/sirupsen/logrus"
)

// RunOnce executes centry with the given arguments and exits with a code
func RunOnce(inputArgs []string) int {
	// Args
	file := ""
	args := []string{}
	if len(inputArgs) >= 2 {
		file = inputArgs[1]
		args = inputArgs[2:]
	}

	// Load manifest
	manifest := config.LoadManifest(file)

	// Configure and create logger
	lf := logger.CreateFactory(manifest.Config.Log.Level, manifest.Config.Log.Prefix)
	log := lf.CreateLogger()

	// Create global options
	options := createGlobalOptions(manifest)

	// Parse global options to get cli args
	args = options.Parse(args)

	// Initialize cli
	c := &cli.CLI{
		Name:    manifest.Config.Name,
		Version: manifest.Config.Version,

		Commands: map[string]cli.CommandFactory{},
		Args:     args,
		HelpFunc: centryHelpFunc(manifest, options),

		// Autocomplete:          true,
		// AutocompleteInstall:   "install-autocomplete",
		// AutocompleteUninstall: "uninstall-autocomplete",
	}

	// Override the current log level from options
	logLevel := options.GeString("config.log.level")
	if options.GetBool("quiet") {
		logLevel = "panic"
	}
	lf.TrySetLogLevel(logLevel)

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

		for _, bf := range command.GetFunctions() {
			cmdName := bf
			cmdKey := strings.Replace(cmdName, ":", " ", -1)
			log.Debugf("Adding command %s", cmdKey)

			c.Commands[cmdKey] = func() (cli.Command, error) {
				return &BashCommand{
					Manifest: manifest,
					Log: log.WithFields(logrus.Fields{
						"command": cmdKey,
					}),
					GlobalOptions: options,
					Name:          cmdName,
					Path:          cmd.Path,
					HelpText:      cmd.Help,
					SynopsisText:  cmd.Description,
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
