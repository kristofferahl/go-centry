package shell

import "github.com/kristofferahl/go-centry/internal/pkg/io"

// Executable defines the interface of an executable program
type Executable interface {
	Run(io io.InputOutput, args []string) error
}

// Function defines a function
type Function struct {
	Name        string
	Description string
	Help        string
}

// Script defines the interface of a script file
type Script interface {
	Language() string
	Executable() Executable
	FullPath() string
	RelativePath() string
	Functions() (funcs []*Function, err error)
	CreateFunctionNamespace(name string) string
	FunctionNameSplitChar() string
}
