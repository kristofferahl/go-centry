package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/shell"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ScriptCommand is a Command implementation that applies stuff
type ScriptCommand struct {
	Context       *Context
	Log           *logrus.Entry
	Command       config.Command
	GlobalOptions *cmd.OptionsSet
	Script        shell.Script
	Function      shell.Function
}

// GetCommandInvocation returns the command invocation string
func (sc *ScriptCommand) GetCommandInvocation() string {
	return strings.Replace(sc.Function.Name, sc.Script.FunctionNamespaceSplitChar(), " ", -1)
}

// GetCommandInvocationPath returns the command invocation path
func (sc *ScriptCommand) GetCommandInvocationPath() []string {
	return strings.Split(sc.GetCommandInvocation(), " ")
}

// ToCLICommand returns a CLI command
func (sc *ScriptCommand) ToCLICommand() *cli.Command {
	cmdKeys := sc.GetCommandInvocationPath()
	cmdName := cmdKeys[len(cmdKeys)-1]
	return &cli.Command{
		Name:            cmdName,
		Usage:           sc.Command.Description,
		UsageText:       sc.Command.Help,
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			ec := sc.Run(c, c.Args().Slice())
			if ec > 0 {
				return cli.Exit("command exited with non zero exit code", ec)
			}
			return nil
		},
		Flags: optionsSetToFlags(sc.Function.Options),
	}
}

// Run builds the source and executes it
func (sc *ScriptCommand) Run(c *cli.Context, args []string) int {
	sc.Log.Debugf("Executing command \"%v\"", sc.Function.Name)

	var source string
	switch sc.Script.Language() {
	case "bash":
		source = generateBashSource(c, sc, args)
		sc.Log.Debugf("Generated bash source:\n%s\n", source)
	default:
		sc.Log.Errorf("Unsupported script language %s", sc.Script.Language())
		return 1
	}

	err := sc.Script.Executable().Run(sc.Context.io, []string{"-c", source})
	if err != nil {
		exitCode := 1

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}

		sc.Log.Errorf("Command %v exited with error! %v", sc.Function.Name, err)
		return exitCode
	}

	sc.Log.Debugf("Finished executing command %v...", sc.Function.Name)
	return 0
}

func generateBashSource(c *cli.Context, sc *ScriptCommand, args []string) string {
	source := []string{}
	source = append(source, "#!/usr/bin/env bash")

	source = append(source, "")
	source = append(source, "# Set working directory")
	source = append(source, fmt.Sprintf("cd %s || exit 1", sc.Context.manifest.BasePath))

	source = append(source, "")
	source = append(source, "# Set exports from global options")

	for _, v := range optionsSetToEnvVars(c, sc.GlobalOptions) {
		if v.Value != "" {
			value := v.Value
			if v.IsString() {
				value = fmt.Sprintf("'%s'", v.Value)
			}
			source = append(source, fmt.Sprintf("export %s=%s", v.Name, value))
		}
	}

	source = append(source, "")
	source = append(source, "# Set exports from local options")

	for _, v := range optionsSetToEnvVars(c, sc.Function.Options) {
		if v.Value != "" {
			value := v.Value
			if v.IsString() {
				value = fmt.Sprintf("'%s'", v.Value)
			}
			source = append(source, fmt.Sprintf("export %s=%s", v.Name, value))
		}
	}

	source = append(source, "")
	source = append(source, "# Sourcing scripts")
	for _, s := range sc.Context.manifest.Scripts {
		source = append(source, fmt.Sprintf("source %s", s))
	}

	source = append(source, "")
	source = append(source, "# Sourcing command")
	source = append(source, fmt.Sprintf("source %s", sc.Script.FullPath()))

	source = append(source, "")
	source = append(source, "# Executing command")
	source = append(source, fmt.Sprintf("%s %s", sc.Function.Name, strings.Join(args, " ")))

	return strings.Join(source, "\n")
}
