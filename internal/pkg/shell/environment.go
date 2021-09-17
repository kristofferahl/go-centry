package shell

import (
	"sort"
)

// EnvironmentVariableType represents a type of an environment variable
type EnvironmentVariableType string

const (
	// EnvironmentVariableTypeString represents a string environment variable
	EnvironmentVariableTypeString EnvironmentVariableType = "string"

	// EnvironmentVariableTypeBool represents a boolean environment variable
	EnvironmentVariableTypeBool EnvironmentVariableType = "bool"

	// EnvironmentVariableTypeInteger represents an integer environment variable
	EnvironmentVariableTypeInteger EnvironmentVariableType = "integer"
)

// EnvironmentVariable represents an environment variable
type EnvironmentVariable struct {
	Name  string
	Value string
	Type  EnvironmentVariableType
}

// IsString returns true if the environment variable is of type string
func (v EnvironmentVariable) IsString() bool {
	return v.Type == EnvironmentVariableTypeString
}

// IsBool returns true if the environment variable is of type boolean
func (v EnvironmentVariable) IsBool() bool {
	return v.Type == EnvironmentVariableTypeBool
}

// IsInteger returns true if the environment variable is of type integer
func (v EnvironmentVariable) IsInteger() bool {
	return v.Type == EnvironmentVariableTypeInteger
}

// SortEnvironmentVariables sorts environment variables by name
func SortEnvironmentVariables(vars []EnvironmentVariable) []EnvironmentVariable {
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}
