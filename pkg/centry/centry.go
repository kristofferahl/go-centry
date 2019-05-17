package centry

import (
	"strings"

	"github.com/kristofferahl/cli"
	"github.com/kristofferahl/go-centry/pkg/config"
	"github.com/kristofferahl/go-centry/pkg/io"
	"github.com/kristofferahl/go-centry/pkg/logging"
	"github.com/sirupsen/logrus"
)

// Runtime defines the runtime
type Runtime struct {
	context *Context
	cli     *cli.CLI
}

// Executor is the name of the executor
type Executor string

// CLI executor
var CLI Executor = "CLI"
// Context defines the current context
type Context struct {
	executor Executor
	io       io.InputOutput
	log      *logging.LogManager
	manifest *config.Manifest
}

// NewContext creates a new context
func NewContext(executor Executor, io io.InputOutput) *Context {
	return &Context{
		executor: executor,
		io:       io,
	}
}

// Create builds a runtime based on the given arguments
func Create(inputArgs []string, context *Context) *Runtime {
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
	context.manifest = config.LoadManifest(file)

	// Create the log manager
	context.log = logging.CreateManager(context.manifest.Config.Log.Level, context.manifest.Config.Log.Prefix, context.io)

	// Create global options
	options := createGlobalOptions(context.manifest)

	// Parse global options to get cli args
	args = options.Parse(args)

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
	logLevel := options.GeString("config.log.level")
	if options.GetBool("quiet") {
		logLevel = "panic"
	}
	context.log.TrySetLogLevel(logLevel)

	logger := context.log.GetLogger()

	// Build commands
	for _, cmd := range context.manifest.Commands {
		cmd := cmd
		command := &DynamicCommand{
			Log: logger.WithFields(logrus.Fields{
				"command": cmd.Name,
			}),
			Command:  cmd,
			Manifest: context.manifest,
		}

		for _, bf := range command.GetFunctions() {
			cmdName := bf
			cmdKey := strings.Replace(cmdName, ":", " ", -1)
			logger.Debugf("Adding command %s", cmdKey)

			c.Commands[cmdKey] = func() (cli.Command, error) {
				return &BashCommand{
					Manifest: context.manifest,
					Log: logger.WithFields(logrus.Fields{
						"command": cmdKey,
					}),
					GlobalOptions: options,
					Name:          cmdName,
					Path:          cmd.Path,
					HelpText:      cmd.Help,
					SynopsisText:  cmd.Description,
					IO:            context.io,
				}, nil
			}
		}
	}

	runtime.context = context
	runtime.cli = c

	return runtime
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
