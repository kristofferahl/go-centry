package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// BashCommand is a Command implementation that applies stuff
type BashCommand struct {
	Manifest      *manifest
	Log           *logrus.Entry
	GlobalOptions *OptionsSet
	Name          string
	Path          string
	HelpText      string
	SynopsisText  string
}

// Run builds the source and executes it
func (bc *BashCommand) Run(args []string) int {
	bc.Log.Debugf("Executing command \"%v\"", bc.Name)

	source := []string{}
	source = append(source, "#!/usr/bin/env bash")

	source = append(source, "")
	source = append(source, "# Set working directory")
	source = append(source, fmt.Sprintf("cd %s || exit 1", bc.Manifest.BasePath))

	source = append(source, "")
	source = append(source, "# Set exports from flags")
	exports := map[string]string{}
	for _, o := range bc.GlobalOptions.Items {
		envName := o.EnvName
		var envValue string
		value := bc.GlobalOptions.GetValue(o.Name)

		if envName == "" {
			envName = strings.Replace(strings.ToUpper(o.Name), ".", "_", -1)
		}

		switch value {
		case "true":
			envValue = o.Name
		case "false":
			envValue = ""
		default:
			envValue = value
		}

		if envValue != "" {
			exports[envName] = envValue
		}
	}
	for envName, envValue := range exports {
		source = append(source, fmt.Sprintf("export %s='%s'", envName, envValue))
	}

	source = append(source, "")
	source = append(source, "# Sourcing scripts")
	for _, s := range bc.Manifest.Scripts {
		source = append(source, fmt.Sprintf("source %s", s))
	}

	source = append(source, "")
	source = append(source, "# Sourcing command")
	source = append(source, fmt.Sprintf("source %s", bc.Path))

	source = append(source, "")
	source = append(source, "# Executing command")
	source = append(source, fmt.Sprintf("%s %s", bc.Name, strings.Join(args, " ")))

	bc.Log.Debugf("Command source code:\n%s\n", strings.Join(source, "\n"))

	callArgs := []string{}
	callArgs = append(callArgs, "-c", strings.Join(source, "\n"))

	err := bc.ExecBash(strings.Join(source, "\n"))
	if err != nil {
		exitCode := 1

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}

		bc.Log.Fatalf("Command %v exited with error! %v", bc.Name, err)
		return exitCode
	}

	bc.Log.Debugf("Finished executing command %v...", bc.Name)
	return 0
}

// ExecBash executes source code using /bin/bash
func (bc *BashCommand) ExecBash(source string) error {
	cmd := exec.Command("/bin/bash", "-c", source)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

// Help returns the help text of the BashCommand
func (bc *BashCommand) Help() string {
	return bc.HelpText
}

// Synopsis returns the synopsis of the BashCommand
func (bc *BashCommand) Synopsis() string {
	return bc.SynopsisText
}
