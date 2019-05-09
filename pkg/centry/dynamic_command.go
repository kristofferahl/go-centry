package centry

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/kristofferahl/go-centry/pkg/config"
	"github.com/sirupsen/logrus"
)

// DynamicCommand encapsulates operations on the script file containing commands
type DynamicCommand struct {
	Manifest *config.Manifest
	Log      *logrus.Entry
	Command  config.Command
}

// GetFullPath returns the absolute path of the script file
func (dc *DynamicCommand) GetFullPath() string {
	absPath := path.Join(dc.Manifest.BasePath, dc.Command.Path)
	return absPath
}

// GetFunctions returns the command functions matching the command name
func (dc *DynamicCommand) GetFunctions() []string {
	callArgs := []string{"-c", fmt.Sprintf("source %s; declare -F", dc.GetFullPath())}
	out, err := exec.Command("/bin/bash", callArgs...).CombinedOutput()
	if err != nil {
		dc.Log.Fatal(err, string(out))
	}

	commands := []string{}
	for _, fun := range strings.Split(string(out), "\n") {
		if fun != "" {
			name := strings.Replace(fun, "declare -f ", "", -1)
			prefixName := fmt.Sprintf("%s:", dc.Command.Name)

			if name == dc.Command.Name || strings.HasPrefix(name, prefixName) {
				commands = append(commands, name)
			}
		}
	}
	return commands
}
