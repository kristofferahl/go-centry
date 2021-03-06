package main

import (
	"fmt"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	"github.com/kristofferahl/go-centry/internal/pkg/shell"
	"github.com/urfave/cli/v2"
)

func configureDefaultOptions() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show help",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the version",
	}
}

func createGlobalOptions(runtime *Runtime) *cmd.OptionsSet {
	context := runtime.context
	manifest := context.manifest

	// Add global options
	options := cmd.NewOptionsSet("Global")

	options.Add(&cmd.Option{
		Type:        cmd.StringOption,
		Name:        "centry-config-log-level",
		Description: "Overrides the log level",
		Default:     manifest.Config.Log.Level,
		Hidden:      manifest.Config.HideInternalOptions,
		Internal:    true,
	})
	options.Add(&cmd.Option{
		Type:        cmd.BoolOption,
		Name:        "centry-quiet",
		Description: "Disables logging",
		Default:     false,
		Hidden:      manifest.Config.HideInternalOptions,
		Internal:    true,
	})

	// Adding global options specified by the manifest
	for _, o := range manifest.Options {
		o := o

		if context.optionEnabledFunc != nil && context.optionEnabledFunc(o) == false {
			continue
		}

		var def interface{}

		switch o.Type {
		case cmd.SelectOption:
			def = false
		case cmd.BoolOption:
			def = false
		default:
			def = o.Default
		}

		options.Add(&cmd.Option{
			Type:        o.Type,
			Name:        o.Name,
			Short:       o.Short,
			Description: o.Description,
			EnvName:     o.EnvName,
			Default:     def,
			Hidden:      o.Hidden,
		})

		runtime.events = append(runtime.events, fmt.Sprintf("registered global option \"%s\"", o.Name))
	}

	return options
}

func optionsSetToFlags(options *cmd.OptionsSet) []cli.Flag {
	flags := make([]cli.Flag, 0)

	for _, o := range options.Sorted() {
		short := []string{o.Short}
		if o.Short == "" {
			short = nil
		}

		switch o.Type {
		case cmd.SelectOption:
			def := false
			if o.Default != nil {
				def = o.Default.(bool)
			}
			flags = append(flags, &cli.BoolFlag{
				Name:    o.Name,
				Aliases: short,
				Usage:   o.Description,
				Value:   def,
				Hidden:  o.Hidden,
			})
		case cmd.BoolOption:
			def := false
			if o.Default != nil {
				def = o.Default.(bool)
			}
			flags = append(flags, &cli.BoolFlag{
				Name:    o.Name,
				Aliases: short,
				Usage:   o.Description,
				Value:   def,
				Hidden:  o.Hidden,
			})
		case cmd.StringOption:
			def := ""
			if o.Default != nil {
				def = o.Default.(string)
			}
			flags = append(flags, &cli.StringFlag{
				Name:    o.Name,
				Aliases: short,
				Usage:   o.Description,
				Value:   def,
				Hidden:  o.Hidden,
			})
		}
	}

	return flags
}

func optionsSetToEnvVars(c *cli.Context, set *cmd.OptionsSet, prefix string) []shell.EnvironmentVariable {
	envVars := make([]shell.EnvironmentVariable, 0)
	for _, o := range set.Sorted() {
		o := o

		envName := o.EnvName
		if envName == "" {
			envName = o.Name
		}
		envName = strings.Replace(strings.ToUpper(envName), ".", "_", -1)
		envName = strings.Replace(strings.ToUpper(envName), "-", "_", -1)

		if prefix != "" && o.Internal == false {
			envName = prefix + envName
		}

		value := c.String(o.Name)

		switch o.Type {
		case cmd.BoolOption:
			envVars = append(envVars, shell.EnvironmentVariable{
				Name:  envName,
				Value: value,
				Type:  shell.EnvironmentVariableTypeBool,
			})
		case cmd.StringOption:
			envVars = append(envVars, shell.EnvironmentVariable{
				Name:  envName,
				Value: value,
				Type:  shell.EnvironmentVariableTypeString,
			})
		case cmd.SelectOption:
			if value == "true" {
				envVars = append(envVars, shell.EnvironmentVariable{
					Name:  envName,
					Value: o.Name,
					Type:  shell.EnvironmentVariableTypeString,
				})
			}
		}
	}

	return shell.SortEnvironmentVariables(envVars)
}
