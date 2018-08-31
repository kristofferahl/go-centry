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

	// Args
	file := ""
	args := []string{}
	if len(os.Args) >= 2 {
		file = os.Args[1]
		args = os.Args[2:]
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
	globalFlags := flag.NewFlagSet("global", flag.ExitOnError)
	logLevel := globalFlags.String("config.log.level", "", "Overrides the manifest log level")
	loggingOff := globalFlags.Bool("quiet", false, "Disables logging")
	globalFlags.BoolVar(loggingOff, "q", false, "Disables logging")
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

	// Parse global flags
	parseGlobalFlags := true
	for _, arg := range args {
		if arg == "-v" || arg == "--v" || arg == "-version" || arg == "--version" {
			parseGlobalFlags = false
			break
		}
		if arg == "-h" || arg == "--h" || arg == "-help" || arg == "--help" {
			parseGlobalFlags = false
			break
		}
	}
	if parseGlobalFlags {
		globalFlags.Parse(args)
		args = globalFlags.Args()
	}

	// Initialize cli
	c := &cli.CLI{
		Name:    manifest.Config.Name,
		Version: manifest.Config.Version,

		Commands: map[string]cli.CommandFactory{},
		Args:     args,
		HelpFunc: centryHelpFunc(manifest.Config.Name, globalFlags), // TODO: Pass manifest instead of globalFlags to get correct flags order and short flags grouped with long option

		// Autocomplete:          true,
		// AutocompleteInstall:   "install-autocomplete",
		// AutocompleteUninstall: "uninstall-autocomplete",
	}

	// Override the manifest log level
	if *loggingOff == true {
		*logLevel = "panic"
	}
	if *logLevel != "" {
		log.Debugf("Overriding manifest loglevel %s", *logLevel, l)
		l, _ := logrus.ParseLevel(*logLevel)
		log.SetLevel(l)
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
					GlobalFlags:  globalFlags,
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
