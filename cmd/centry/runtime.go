package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/log"
	"github.com/kristofferahl/go-centry/internal/pkg/shell"
	"github.com/sirupsen/logrus"
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

	// Initialize cli
	app := &cli.App{
		Name:      context.manifest.Config.Name,
		HelpName:  context.manifest.Config.Name,
		Usage:     "A tool for building declarative CLI's over bash scripts, written in go.", // TODO: Set from manifest config
		UsageText: "",
		Version:   context.manifest.Config.Version,

		Commands: make([]*cli.Command, 0),
		Flags:    toCliFlags(options),

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

	logger := context.log.GetLogger()

	// Register builtin commands
	if context.executor == CLI {
		serveCmd := &ServeCommand{
			Manifest: context.manifest,
			Log: logger.WithFields(logrus.Fields{
				"command": "serve",
			}),
		}
		app.Commands = append(app.Commands, serveCmd.ToCLICommand())
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
				cmd := cmd
				namespace := script.CreateFunctionNamespace(cmd.Name)

				if fn.Name != cmd.Name && strings.HasPrefix(fn.Name, namespace) == false {
					continue
				}

				cmdDescription := cmd.Description
				if fn.Description != "" {
					cmd.Description = fn.Description
				}

				cmdHelp := cmd.Help
				if fn.Help != "" {
					cmd.Help = fn.Help
				}

				scriptCmd := &ScriptCommand{
					Context:       context,
					Log:           logger.WithFields(logrus.Fields{}),
					GlobalOptions: options,
					Command:       cmd,
					Script:        script,
					Function:      *fn,
				}
				cliCmd := scriptCmd.ToCLICommand()

				cmdKeyParts := scriptCmd.GetCommandInvocationPath()

				var root *cli.Command
				for depth, cmdKeyPart := range cmdKeyParts {
					if depth == 0 {
						if getCommand(app.Commands, cmdKeyPart) == nil {
							if depth == len(cmdKeyParts)-1 {
								// add destination command
								app.Commands = append(app.Commands, cliCmd)
							} else {
								// add placeholder
								app.Commands = append(app.Commands, &cli.Command{
									Name:            cmdKeyPart,
									Usage:           cmdDescription,
									UsageText:       cmdHelp,
									HideHelpCommand: true,
									Action:          nil,
								})
							}
						}
						root = getCommand(app.Commands, cmdKeyPart)
					} else {
						if getCommand(root.Subcommands, cmdKeyPart) == nil {
							if depth == len(cmdKeyParts)-1 {
								// add destination command
								root.Subcommands = append(root.Subcommands, cliCmd)
							} else {
								// add placeholder
								root.Subcommands = append(root.Subcommands, &cli.Command{
									Name:            cmdKeyPart,
									Usage:           "...",
									UsageText:       "",
									HideHelpCommand: true,
									Action:          nil,
								})
							}
						}
						root = getCommand(root.Subcommands, cmdKeyPart)
					}
				}

				runtime.events = append(runtime.events, fmt.Sprintf("Registered command \"%s\"", scriptCmd.GetCommandInvocation()))
			}
		}
	}

	sortCommands(app.Commands)

	runtime.cli = app

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

func createScript(cmd config.Command, context *Context) shell.Script {
	return &shell.BashScript{
		BasePath: context.manifest.BasePath,
		Path:     cmd.Path,
		Log: context.log.GetLogger().WithFields(logrus.Fields{
			"script": cmd.Path,
		}),
	}
}

func getCommand(commands []*cli.Command, name string) *cli.Command {
	for _, c := range commands {
		if c.HasName(name) {
			return c
		}
	}

	return nil
}

func sortCommands(commands []*cli.Command) {
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
}
