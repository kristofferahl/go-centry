package shell

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// BashScript encapsulates operations on the script file containing commands
type BashScript struct {
	BasePath string
	Path     string
	Log      *logrus.Entry
}

// Language returns the name of the script language
func (s *BashScript) Language() string {
	return "bash"
}

// FullPath returns the absolute path of the script file
func (s *BashScript) FullPath() string {
	return path.Join(s.BasePath, s.Path)
}

// Functions returns the command functions matching the command name
func (s *BashScript) Functions() ([]string, error) {
	callArgs := []string{"-c", fmt.Sprintf("source %s; declare -F", s.FullPath())}
	out, err := exec.Command("/bin/bash", callArgs...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(string(out), err)
	}

	functions := []string{}

	for _, fun := range strings.Split(string(out), "\n") {
		if fun != "" {
			name := strings.Replace(fun, "declare -f ", "", -1)
			functions = append(functions, name)
		}
	}

	return functions, nil
}

// CreateFunctionNamespace returns a namespaced function name
func (s *BashScript) CreateFunctionNamespace(name string) string {
	return fmt.Sprintf("%s%s", name, s.FunctionNameSplitChar())
}

// FunctionNameSplitChar returns the separator used for function namespaces
func (s *BashScript) FunctionNameSplitChar() string {
	return ":"
}
