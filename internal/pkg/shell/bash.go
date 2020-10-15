package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
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
	callArgs := []string{"-c", fmt.Sprintf("set -e; source %s; declare -F", s.FullPath())}

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

// FunctionAnnotations returns function annotations declared in the script file
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

	// TODO: Make this part shared for all shell types
	for _, fname := range fnames {
		s.Log.WithFields(logrus.Fields{
			"func": fname,
		}).Debugf("building function")

		f := &Function{
			Name:    fname,
			Options: cmd.NewOptionsSet(fname),
		}

		options := make(map[string]*cmd.Option, 0)

		for _, a := range annotations {
			cmdName := a.NamespaceValues["cmd"]
			if cmdName == "" || cmdName != f.Name {
				continue
			}

			s.Log.WithFields(logrus.Fields{
				"func":      a.NamespaceValues["cmd"],
				"namespace": a.Namespace,
				"key":       a.Key,
			}).Debugf("handling annotation")

			switch a.Namespace {
			case config.CommandAnnotationCmdOptionNamespace:
				name := a.NamespaceValues["option"]
				if name == "" {
					continue
				}
				if options[name] == nil {
					options[name] = &cmd.Option{Type: cmd.StringOption, Name: name}
				}
				switch a.Key {
				case "type":
					options[name].Type = cmd.StringToOptionType(a.Value)
				case "short":
					options[name].Short = a.Value
				case "envName":
					options[name].EnvName = a.Value
				case "description":
					options[name].Description = a.Value
				case "default":
					options[name].Default = a.Value
				}
			case config.CommandAnnotationCmdNamespace:
				switch a.Key {
				case "description":
					f.Description = a.Value
				case "help":
					f.Help = a.Value
				}
			}
		}

		for _, v := range options {
			if err := v.Validate(); err != nil {
				s.Log.WithFields(logrus.Fields{
					"option": v.Name,
					"type":   v.Type,
				}).Warn(err.Error())
			} else {
				f.Options.Add(v)
			}
		}

		funcs = append(funcs, f)
	}

	return funcs, nil
}

// FunctionNamespace returns a namespaced function name
func (s *BashScript) FunctionNamespace(name string) string {
	return fmt.Sprintf("%s%s", name, s.FunctionNamespaceSplitChar())
}

// FunctionNamespaceSplitChar returns the separator used for function namespaces
func (s *BashScript) FunctionNamespaceSplitChar() string {
	return ":"
}
