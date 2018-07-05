package main

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// BashCommand is a Command implementation that applies stuff
type BashCommand struct {
	Manifest     *manifest
	Log          *logrus.Entry
	Name         string
	Path         string
	HelpText     string
	SynopsisText string
}

func (bc *BashCommand) GetCommandPath() string {
	absPath := path.Join(bc.Manifest.BasePath, bc.Path)
	return absPath
}

func (bc *BashCommand) Run(args []string) int {
	callArgs := []string{}
	callArgs = append(callArgs, "-c", fmt.Sprintf("source %s; %s %s", bc.GetCommandPath(), bc.Name, strings.Join(args, " ")))

	bc.Log.Infof("Executing command %v", bc.Name)

	out, err := exec.Command("/bin/bash", callArgs...).CombinedOutput()

	fmt.Printf(string(out))

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

	bc.Log.Infof("Finished executing command %v...", bc.Name)
	return 0
}

func (bc *BashCommand) Help() string {
	return bc.HelpText
}

func (bc *BashCommand) Synopsis() string {
	return bc.SynopsisText
}
