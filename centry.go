package main

import (
	"os"
	"strings"

	"github.com/kristofferahl/cli"
	"github.com/sirupsen/logrus"
)

func centry(osArgs []string) int {
	log := logrus.New()

	// Args
	file := ""
	args := []string{}
	if len(osArgs) >= 2 {
		file = osArgs[1]
		args = osArgs[2:]
	}

	// Load manifest
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Error("The first argument of centry must be a path to a valid manfest file")
		return 1
	}
	manifest := loadManifest(file)

	// Configure logger
	l, _ := logrus.ParseLevel(manifest.Config.Log.Level)
	log.SetLevel(l)
	log.Formatter = &PrefixedTextFormatter{
		Prefix: manifest.Config.Log.Prefix,
	}

	// Add global option flags
	// TODO: Allow anything under config to be overridden using flags
	options := NewOptionsSet(optionSetGlobal)
	options.Add(&Option{
		Name:        "config.log.level",
		Description: "Overrides the log level",
		Default:     log.Level.String(),
	})
	options.Add(&Option{
		Name:        "quiet",
		Short:       "q",
		Description: "Disables logging",
	})
	options.Add(&Option{
		Name:        "help",
		Short:       "h",
		Description: "Displays help",
	})
	options.Add(&Option{
		Name:        "version",
		Short:       "v",
		Description: "Displays the version fo the cli",
	})

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

	// Override the current log level
	logLevel := options.GeString("config.log.level")
	log.Debugf("Current loglevel is (%s)..", l)

	loggingOff := options.GetBool("quiet")
	if loggingOff == true {
		logLevel = "panic"
	}

	if logLevel != "" {
		log.Debugf("Changing loglevel to value from option (%s)..", logLevel)
		l, _ := logrus.ParseLevel(logLevel)
		log.SetLevel(l)
		log.Debugf("Changed loglevel to (%s)..", l)
	}

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
