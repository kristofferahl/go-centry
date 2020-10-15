package main

import (
	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const metadataExitCode string = "exitcode"

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
		Usage:     context.manifest.Config.Description,
		UsageText: "",
		Version:   context.manifest.Config.Version,

		Commands: make([]*cli.Command, 0),
		Flags:    optionsSetToFlags(options),

		HideHelpCommand:       true,
		CustomAppHelpTemplate: cliHelpTemplate,

		Writer:    context.io.Stdout,
		ErrWriter: context.io.Stderr,

		Before: func(c *cli.Context) error {
			return handleBefore(runtime, c)
		},
		CommandNotFound: func(c *cli.Context, command string) {
			handleCommandNotFound(runtime, c, command)
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			handleExitErr(runtime, c, err)
		},
	}

	// Register internal commands
	registerInternalCommands(runtime)

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
		runtime.context.log.GetLogger().Error(err)
	}

	// Return exitcode defined in metadata
	if runtime.cli.Metadata[metadataExitCode] != nil {
		switch runtime.cli.Metadata[metadataExitCode].(type) {
		case int:
			return runtime.cli.Metadata[metadataExitCode].(int)
		}
		return 128
	}

	return 0
}

func handleBefore(runtime *Runtime, c *cli.Context) error {
	// Override the current log level from options
	logLevel := c.String("config.log.level")
	if c.Bool("quiet") {
		logLevel = "panic"
	}
	runtime.context.log.TrySetLogLevel(logLevel)

	// Print runtime events
	logger := runtime.context.log.GetLogger()
	for _, e := range runtime.events {
		logger.Debugf(e)
	}

	return nil
}

func handleCommandNotFound(runtime *Runtime, c *cli.Context, command string) {
	logger := runtime.context.log.GetLogger()
	logger.WithFields(logrus.Fields{
		"command": command,
	}).Warnf("Command not found!")
	c.App.Metadata[metadataExitCode] = 127
}

// Handles errors implementing ExitCoder by printing their
// message and calling OsExiter with the given exit code.
// If the given error instead implements MultiError, each error will be checked
// for the ExitCoder interface, and OsExiter will be called with the last exit
// code found, or exit code 1 if no ExitCoder is found.
func handleExitErr(runtime *Runtime, context *cli.Context, err error) {
	if err == nil {
		return
	}

	logger := runtime.context.log.GetLogger()

	if exitErr, ok := err.(cli.ExitCoder); ok {
		if err.Error() != "" {
			if _, ok := exitErr.(cli.ErrorFormatter); ok {
				logger.WithFields(logrus.Fields{
					"command": context.Command.Name,
					"code":    exitErr.ExitCode(),
				}).Errorf("%+v\n", err)
			} else {
				logger.WithFields(logrus.Fields{
					"command": context.Command.Name,
					"code":    exitErr.ExitCode(),
				}).Error(err)
			}
		}
		cli.OsExiter(exitErr.ExitCode())
		return
	}

	if multiErr, ok := err.(cli.MultiError); ok {
		code := handleMultiError(runtime, context, multiErr)
		cli.OsExiter(code)
		return
	}
}

func handleMultiError(runtime *Runtime, context *cli.Context, multiErr cli.MultiError) int {
	code := 1
	for _, merr := range multiErr.Errors() {
		if multiErr2, ok := merr.(cli.MultiError); ok {
			code = handleMultiError(runtime, context, multiErr2)
		} else if merr != nil {
			if exitErr, ok := merr.(cli.ExitCoder); ok {
				code = exitErr.ExitCode()
				runtime.context.log.GetLogger().WithFields(logrus.Fields{
					"command": context.Command.Name,
					"code":    code,
				}).Error(merr)
			} else {
				runtime.context.log.GetLogger().WithFields(logrus.Fields{
					"command": context.Command.Name,
				}).Error(merr)
			}
		}
	}
	return code
}
