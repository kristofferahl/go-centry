package main

import (
	"sort"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
)

type envType string

const (
	envTypeString envType = "string"
	envTypeBool   envType = "bool"
)

type envVar struct {
	Name  string
	Value string
	Type  envType
}

func (v envVar) IsString() bool {
	return v.Type == envTypeString
}

func (v envVar) IsBool() bool {
	return v.Type == envTypeBool
}

func optionsSetToEnvVars(set *cmd.OptionsSet) []envVar {
	envVars := make([]envVar, 0)
	for _, o := range set.Sorted() {
		o := o

		envName := o.EnvName
		if envName == "" {
			envName = o.Name
		}
		envName = strings.Replace(strings.ToUpper(envName), ".", "_", -1)
		envName = strings.Replace(strings.ToUpper(envName), "-", "_", -1)

		value := set.GetValueString(o.Name)

		switch o.Type {
		case cmd.BoolOption:
			envVars = append(envVars, envVar{
				Name:  envName,
				Value: value,
				Type:  envTypeBool,
			})
		case cmd.StringOption:
			envVars = append(envVars, envVar{
				Name:  envName,
				Value: value,
				Type:  envTypeString,
			})
		case cmd.SelectOption:
			if value == "true" {
				envVars = append(envVars, envVar{
					Name:  envName,
					Value: o.Name,
					Type:  envTypeString,
				})
			}
		}
	}

	return sortEnv(envVars)
}

func sortEnv(vars []envVar) []envVar {
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}
