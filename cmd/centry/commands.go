package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/shell"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func registerInternalCommands(runtime *Runtime) {
	context := runtime.context

	if context.executor == CLI {
		serveCmd := &ServeCommand{
			Manifest: context.manifest,
			Log: context.log.GetLogger().WithFields(logrus.Fields{
				"command": "serve",
			}),
		}
		internalCmd := &cli.Command{
			Name:      "internal",
			Usage:     "Internal centry commands",
			UsageText: "",
			Hidden:    context.manifest.Config.HideInternalCommands,
			Subcommands: []*cli.Command{
				serveCmd.ToCLICommand(),
			},
		}
		runtime.cli.Commands = append(runtime.cli.Commands, internalCmd)
	}
}

func registerManifestCommands(runtime *Runtime, options *cmd.OptionsSet) {
	context := runtime.context

	for _, cmd := range context.manifest.Commands {
		cmd := cmd

		if context.commandEnabledFunc != nil && context.commandEnabledFunc(cmd) == false {
			continue
		}

		script := createScript(cmd, context)

		funcs, err := script.Functions()
		if err != nil {
			context.log.GetLogger().WithFields(logrus.Fields{
				"command": cmd.Name,
			}).Errorf("Failed to parse script functions. %v", err)
		} else {
			for _, fn := range funcs {
				fn := fn
				cmd := cmd
				namespace := script.FunctionNamespace(cmd.Name)

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
					Log:           context.log.GetLogger().WithFields(logrus.Fields{}),
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
						if getCommand(runtime.cli.Commands, cmdKeyPart) == nil {
							if depth == len(cmdKeyParts)-1 {
								// add destination command
								runtime.cli.Commands = append(runtime.cli.Commands, cliCmd)
							} else {
								// add placeholder
								runtime.cli.Commands = append(runtime.cli.Commands, &cli.Command{
									Name:            cmdKeyPart,
									Usage:           cmdDescription,
									UsageText:       cmdHelp,
									HideHelpCommand: true,
									Action:          nil,
								})
							}
						}
						root = getCommand(runtime.cli.Commands, cmdKeyPart)
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

func createScript(cmd config.Command, context *Context) shell.Script {
	return &shell.BashScript{
		BasePath: context.manifest.BasePath,
		Path:     cmd.Path,
		Log: context.log.GetLogger().WithFields(logrus.Fields{
			"script": cmd.Path,
		}),
	}
}
