package main

import (
	"strings"

	"github.com/kristofferahl/cli"
	"github.com/sirupsen/logrus"
)

func centry(osArgs []string) int {
	// Args
	file := ""
	args := []string{}
	if len(osArgs) >= 2 {
		file = osArgs[1]
		args = osArgs[2:]
	}

	// Load manifest
	manifest := loadManifest(file)

	// Configure and create logger
	lf := &loggerFactory{
		config: &loggerConfig{
			level:  manifest.Config.Log.Level,
			prefix: manifest.Config.Log.Prefix,
		},
	}
	log := lf.createLogger()

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
	lf.trySetLogLevel(logLevel)

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

		for _, bf := range command.GeBashFunctions() {
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
