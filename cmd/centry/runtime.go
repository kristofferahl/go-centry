package main

import (
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/log"
	"github.com/kristofferahl/go-centry/internal/pkg/shell"
	"github.com/mitchellh/cli"
	"github.com/sirupsen/logrus"
)

// Runtime defines the runtime
type Runtime struct {
	context *Context
	cli     *cli.CLI
}

// NewRuntime builds a runtime based on the given arguments
func NewRuntime(inputArgs []string, context *Context) (*Runtime, error) {
	// Create the runtime
	runtime := &Runtime{}

	// Args
	file := ""
	args := []string{}
	if len(inputArgs) >= 1 {
		file = inputArgs[0]
		args = inputArgs[1:]
	}

	// Load manifest
	manifest, err := config.LoadManifest(file)
	if err != nil {
		return nil, err
	}

	context.manifest = manifest

	// Create the log manager
	context.log = log.CreateManager(context.manifest.Config.Log.Level, context.manifest.Config.Log.Prefix, context.io)

	// Create global options
	options := createGlobalOptions(context)

	// Parse global options to get cli args
	args, err = options.Parse(args, context.io)
	if err != nil {
		return nil, err
	}

	// Initialize cli
	c := &cli.CLI{
		Name:    context.manifest.Config.Name,
		Version: context.manifest.Config.Version,

		Commands:   map[string]cli.CommandFactory{},
		Args:       args,
		HelpFunc:   cliHelpFunc(context.manifest, options),
		HelpWriter: context.io.Stderr,

		// Autocomplete:          true,
		// AutocompleteInstall:   "install-autocomplete",
		// AutocompleteUninstall: "uninstall-autocomplete",
	}

	// Override the current log level from options
	logLevel := options.GetString("config.log.level")
	if options.GetBool("quiet") {
		logLevel = "panic"
	}
	context.log.TrySetLogLevel(logLevel)

	logger := context.log.GetLogger()

	// Register builtin commands
	if context.executor == CLI {
		c.Commands["serve"] = func() (cli.Command, error) {
			return &ServeCommand{
				Manifest: context.manifest,
				Log: logger.WithFields(logrus.Fields{
					"command": "serve",
				}),
			}, nil
		}
	}

	// Build commands
	for _, cmd := range context.manifest.Commands {
		cmd := cmd

		if context.commandEnabledFunc != nil && context.commandEnabledFunc(cmd) == false {
			continue
		}

		script := createScript(cmd, context)

		funcs, err := script.Functions()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"command": cmd.Name,
			}).Errorf("Failed to parse script functions. %v", err)

		} else {
			for _, fn := range funcs {
				fn := fn
				namespace := script.CreateFunctionNamespace(cmd.Name)

				if fn != cmd.Name && strings.HasPrefix(fn, namespace) == false {
					continue
				}

				cmdKey := strings.Replace(fn, script.FunctionNameSplitChar(), " ", -1)
				c.Commands[cmdKey] = func() (cli.Command, error) {
					return &ScriptCommand{
						Context:       context,
						Log:           logger.WithFields(logrus.Fields{}),
						GlobalOptions: options,
						Command:       cmd,
						Script:        script,
						Function:      fn,
					}, nil
				}

				logger.Debugf("Registered command \"%s\"", cmdKey)
			}
		}
	}

	runtime.context = context
	runtime.cli = c

	return runtime, nil
}

// Execute runs the CLI and exits with a code
func (runtime *Runtime) Execute() int {
	// Run cli
	exitCode, err := runtime.cli.Run()
	if err != nil {
		logger := runtime.context.log.GetLogger()
		logger.Error(err)
	}

	return exitCode
}

func createScript(cmd config.Command, context *Context) shell.Script {
	return &shell.BashScript{
		BasePath: context.manifest.BasePath,
		Path:     cmd.Path,
	}
}
