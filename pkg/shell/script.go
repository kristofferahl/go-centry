package shell

// ScriptFile defines the interface of a script file
type ScriptFile interface {
	Language() string
	FullPath() string
	Functions() (funcs []string, err error)
	CreateFunctionNamespace(name string) string
	FunctionNameSplitChar() string
}
