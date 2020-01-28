package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/io"
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

// FunctionNames returns functions in declared in the script
func (s *BashScript) FunctionNames() ([]string, error) {
	callArgs := []string{"-c", fmt.Sprintf("source %s; declare -F", s.FullPath())}

	io, buf := io.BufferedCombined()

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

// FunctionAnnotations returns function annotations in declared in the script
func (s *BashScript) FunctionAnnotations() ([]*config.Annotation, error) {
	annotations := make([]*config.Annotation, 0)

	file, err := os.Open(s.FullPath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "#") {
			a, err := config.ParseAnnotation(strings.TrimLeft(t, "#"))
			if err != nil {
				s.Log.Debugf("%s", err.Error())
			} else if a != nil {
				annotations = append(annotations, a)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return annotations, nil
}

// Functions returns the command functions
func (s *BashScript) Functions() ([]*Function, error) {
	funcs := make([]*Function, 0)

	fnames, err := s.FunctionNames()
	if err != nil {
		return nil, err
	}

	annotations, err := s.FunctionAnnotations()
	if err != nil {
		return nil, err
	}

	for _, fname := range fnames {
		f := &Function{Name: fname}
		for _, a := range annotations {
			if a.Key == fname {
				switch a.Namespace {
				case config.CommandAnnotationDescriptionNamespace:
					f.Description = a.Value
				case config.CommandAnnotationHelpNamespace:
					f.Help = a.Value
				}
			}
		}
		funcs = append(funcs, f)
	}

	return funcs, nil
}

// CreateFunctionNamespace returns a namespaced function name
func (s *BashScript) CreateFunctionNamespace(name string) string {
	return fmt.Sprintf("%s%s", name, s.FunctionNameSplitChar())
}

// FunctionNameSplitChar returns the separator used for function namespaces
func (s *BashScript) FunctionNameSplitChar() string {
	return ":"
}
