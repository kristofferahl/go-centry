package main

import (
	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/log"
	"github.com/urfave/cli/v2"
)

// Runtime defines the runtime
type Runtime struct {
	cli     *cli.App
	context *Context
	file    string
	args    []string
	events  []string
}

// NewRuntime builds a runtime based on the given arguments
func NewRuntime(inputArgs []string, context *Context) (*Runtime, error) {
	// Create the runtime
	runtime := &Runtime{
		cli:     nil,
		context: context,
		file:    "",
		args:    []string{},
		events:  []string{},
	}

	// Args
	if len(inputArgs) >= 1 {
		runtime.file = inputArgs[0]
		runtime.args = inputArgs[1:]
	}

	// Load manifest
	manifest, err := config.LoadManifest(runtime.file)
	if err != nil {
		return nil, err
	}
	context.manifest = manifest

	// Create the log manager
	context.log = log.CreateManager(context.manifest.Config.Log.Level, context.manifest.Config.Log.Prefix, context.io)

	// Create global options
	options := createGlobalOptions(context)

	// Configure default options
	configureDefaultOptions()

	// Initialize cli
	runtime.cli = &cli.App{
		Name:      context.manifest.Config.Name,
		HelpName:  context.manifest.Config.Name,
		Usage:     "A tool for building declarative CLI's over bash scripts, written in go.", // TODO: Set from manifest config
		UsageText: "",
		Version:   context.manifest.Config.Version,

		Commands: make([]*cli.Command, 0),
		Flags:    optionsSetToFlags(options),

		HideHelpCommand:       true,
		CustomAppHelpTemplate: cliHelpTemplate,

		Writer:    context.io.Stdout,
		ErrWriter: context.io.Stderr,

		Before: func(c *cli.Context) error {
			// Override the current log level from options
			logLevel := c.String("config.log.level")
			if c.Bool("quiet") {
				logLevel = "panic"
			}
			context.log.TrySetLogLevel(logLevel)

			// Print runtime events
			logger := context.log.GetLogger()
			for _, e := range runtime.events {
				logger.Debugf(e)
			}

			return nil
		},
	}

	// Register builtin commands
	registerBuiltinCommands(runtime)

	// Register manifest commands
	registerManifestCommands(runtime, options)

	// Sort commands
	sortCommands(runtime.cli.Commands)

	return runtime, nil
}

// Execute runs the CLI and exits with a code
func (runtime *Runtime) Execute() int {
	args := append([]string{""}, runtime.args...)

	// Run cli
	err := runtime.cli.Run(args)
	if err != nil {
		logger := runtime.context.log.GetLogger()
		logger.Error(err)
	}

	return 0
}
