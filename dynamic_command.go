package main

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// DynamicCommand is a Command implementation that applies stuff
type DynamicCommand struct {
	Manifest *manifest
	Log      *logrus.Entry
	Command  command
}

func (dc *DynamicCommand) GetCommandPath() string {
	absPath := path.Join(dc.Manifest.BasePath, dc.Command.Path)
	return absPath
}

func (dc *DynamicCommand) GeBashCommands() []string {
	callArgs := []string{"-c", fmt.Sprintf("source %s; declare -F", dc.GetCommandPath())}
	out, err := exec.Command("/bin/bash", callArgs...).CombinedOutput()
	if err != nil {
		dc.Log.Fatal(err, string(out))
	}

	commands := []string{}
	for _, fun := range strings.Split(string(out), "\n") {
		if fun != "" {
			name := strings.Replace(fun, "declare -f ", "", -1)
			if strings.HasPrefix(name, dc.Command.Name) {
				commands = append(commands, name)
			}
		}
	}
	return commands
}
