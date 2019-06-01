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
)

// ScriptCommand is a Command implementation that applies stuff
type ScriptCommand struct {
	Context       *Context
	Log           *logrus.Entry
	Command       config.Command
	GlobalOptions *cmd.OptionsSet
	Script        shell.Script
	Function      string
}

// Run builds the source and executes it
func (c *ScriptCommand) Run(args []string) int {
	c.Log.Debugf("Executing command \"%v\"", c.Function)

	var source string
	switch c.Script.Language() {
	case "bash":
		source = generateBashSource(c, args)
		c.Log.Debugf("Generated bash source:\n%s\n", source)
	default:
		c.Log.Errorf("Unsupported script language %s", c.Script.Language())
		return 1
	}

	err := c.Script.Executable().Run(c.Context.io, []string{"-c", source})
	if err != nil {
		exitCode := 1

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}

		c.Log.Errorf("Command %v exited with error! %v", c.Function, err)
		return exitCode
	}

	c.Log.Debugf("Finished executing command %v...", c.Function)
	return 0
}

// Help returns the help text of the ScriptCommand
func (c *ScriptCommand) Help() string {
	return c.Command.Help
}

// Synopsis returns the synopsis of the ScriptCommand
func (c *ScriptCommand) Synopsis() string {
	return c.Command.Description
}

func generateBashSource(c *ScriptCommand, args []string) string {
	source := []string{}
	source = append(source, "#!/usr/bin/env bash")

	source = append(source, "")
	source = append(source, "# Set working directory")
	source = append(source, fmt.Sprintf("cd %s || exit 1", c.Context.manifest.BasePath))

	source = append(source, "")
	source = append(source, "# Set exports from flags")

	for _, v := range optionsSetToEnvVars(c.GlobalOptions) {
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
	for _, s := range c.Context.manifest.Scripts {
		source = append(source, fmt.Sprintf("source %s", s))
	}

	source = append(source, "")
	source = append(source, "# Sourcing command")
	source = append(source, fmt.Sprintf("source %s", c.Script.FullPath()))

	source = append(source, "")
	source = append(source, "# Executing command")
	source = append(source, fmt.Sprintf("%s %s", c.Function, strings.Join(args, " ")))

	return strings.Join(source, "\n")
}
