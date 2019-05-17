package shell

import "github.com/kristofferahl/go-centry/pkg/io"

// Executable defines the interface of an executable program
type Executable interface {
	Run(io io.InputOutput, args []string) error
}

// Script defines the interface of a script file
type Script interface {
	Language() string
	Executable() Executable
	FullPath() string
	RelativePath() string
	Functions() (funcs []string, err error)
	CreateFunctionNamespace(name string) string
	FunctionNameSplitChar() string
}
