package shell

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/kristofferahl/go-centry/pkg/io"
	"github.com/sirupsen/logrus"
)

// Bash is thin wrapper around the bash executable
type Bash struct {
	Path string
}

// NewBash creates a new bash instance
func NewBash() *Bash {
	return &Bash{
		Path: "/bin/bash",
	}
}

// Run executes the bash with the given arguments
func (bash *Bash) Run(io io.InputOutput, args []string) error {
	cmd := exec.Command(bash.Path, args...)
	cmd.Stdin = io.Stdin
	cmd.Stdout = io.Stdout
	cmd.Stderr = io.Stderr
	return cmd.Run()
}

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

// Executable returns an executable
func (s *BashScript) Executable() Executable {
	return NewBash()
}

// RelativePath returns the relative path of the script file
func (s *BashScript) RelativePath() string {
	return s.Path
}

// FullPath returns the absolute path of the script file
func (s *BashScript) FullPath() string {
	return path.Join(s.BasePath, s.Path)
}

// Functions returns the command functions matching the command name
func (s *BashScript) Functions() ([]string, error) {
	callArgs := []string{"-c", fmt.Sprintf("source %s; declare -F", s.FullPath())}

	var buf bytes.Buffer
	io := io.InputOutput{
		Stdin:  nil,
		Stdout: &buf,
		Stderr: &buf,
	}

	err := NewBash().Run(io, callArgs)
	if err != nil {
		return nil, err
	}

	functions := []string{}

	out := buf.String()
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
