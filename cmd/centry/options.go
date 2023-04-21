package main

import (
	"fmt"
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/shell"
	"github.com/sirupsen/logrus"
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

		err := options.Add(&cmd.Option{
			Type:        o.Type,
			Name:        o.Name,
			Short:       o.Short,
			Description: o.Description,
			EnvName:     o.EnvName,
			Values:      mapOptionValuesToCmdOptionValues(o),
			Default:     o.Default,
			Required:    o.Required,
			Hidden:      o.Hidden,
		})

		if err != nil {
			runtime.events = append(runtime.events, fmt.Sprintf("failed to register global option \"%s\", error: %v", o.Name, err))
			continue
		}

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
				Name:     o.Name,
				Aliases:  short,
				Usage:    o.Description,
				Value:    def,
				Required: false,
				Hidden:   o.Hidden,
			})
		case cmd.SelectOptionV2:
			def := false
			if o.Default != nil {
				def = o.Default.(bool)
			}
			for _, v := range o.Values {
				short := []string{v.Short}
				if v.Short == "" {
					short = nil
				}
				value := v.Value
				if value == "" {
					value = v.Name
				}
				flags = append(flags, &SelectOptionFlag{
					BoolFlag: cli.BoolFlag{
						Name:     v.Name,
						Aliases:  short,
						Usage:    fmt.Sprintf("%s (%s=%s)", o.Description, o.Name, value),
						Value:    def,
						Required: false,
						Hidden:   o.Hidden,
					},
					GroupName:     o.Name,
					GroupRequired: o.Required,
					Values:        o.Values,
				})
			}
		case cmd.IntegerOption:
			def := 0
			if o.Default != nil {
				def = o.Default.(int)
			}
			flags = append(flags, &cli.IntFlag{
				Name:     o.Name,
				Aliases:  short,
				Usage:    o.Description,
				Value:    def,
				Required: o.Required,
				Hidden:   o.Hidden,
			})
		case cmd.BoolOption:
			def := false
			if o.Default != nil {
				def = o.Default.(bool)
			}
			flags = append(flags, &cli.BoolFlag{
				Name:     o.Name,
				Aliases:  short,
				Usage:    o.Description,
				Value:    def,
				Required: o.Required,
				Hidden:   o.Hidden,
			})
		case cmd.StringOption:
			def := ""
			if o.Default != nil {
				def = o.Default.(string)
			}
			flags = append(flags, &cli.StringFlag{
				Name:     o.Name,
				Aliases:  short,
				Usage:    o.Description,
				Value:    def,
				Required: o.Required,
				Hidden:   o.Hidden,
			})
		default:
			panic(fmt.Sprintf("option type \"%s\" not implemented", o.Type))
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
		case cmd.StringOption:
			envVars = append(envVars, shell.EnvironmentVariable{
				Name:  envName,
				Value: value,
				Type:  shell.EnvironmentVariableTypeString,
			})
		case cmd.BoolOption:
			envVars = append(envVars, shell.EnvironmentVariable{
				Name:  envName,
				Value: value,
				Type:  shell.EnvironmentVariableTypeBool,
			})
		case cmd.IntegerOption:
			envVars = append(envVars, shell.EnvironmentVariable{
				Name:  envName,
				Value: value,
				Type:  shell.EnvironmentVariableTypeInteger,
			})
		case cmd.SelectOption:
			if value == "true" {
				envVars = append(envVars, shell.EnvironmentVariable{
					Name:  envName,
					Value: o.Name,
					Type:  shell.EnvironmentVariableTypeString,
				})
			}
		case cmd.SelectOptionV2:
			value := ""
			for _, v := range o.Values {
				ov := c.String(v.Name)
				if ov == "true" {
					value = v.Value
					if value == "" {
						value = v.Name
					}
					break
				}
			}

			if value != "" {
				envVars = append(envVars, shell.EnvironmentVariable{
					Name:  envName,
					Value: value,
					Type:  shell.EnvironmentVariableTypeString,
				})
			}
		default:
			panic(fmt.Sprintf("option type \"%s\" not implemented", o.Type))
		}
	}

	return shell.SortEnvironmentVariables(envVars)
}

func mapOptionValuesToCmdOptionValues(o config.Option) []cmd.OptionValue {
	values := []cmd.OptionValue{}
	for _, v := range o.Values {
		values = append(values, cmd.OptionValue{
			Name:  v.Name,
			Short: v.Short,
			Value: v.Value,
		})
	}
	return values
}

func validateOptionsSet(c *cli.Context, set *cmd.OptionsSet, cmdName string, level string, log *logrus.Entry) error {
	selectOptions := make(map[string][]string)
	selectOptionRequired := make(map[string]bool)
	selectOptionSelectedValues := make(map[string][]string)

	for _, o := range set.Sorted() {
		o := o

		switch o.Type {
		case cmd.SelectOption:
			group := o.EnvName
			selectOptions[group] = append(selectOptions[group], o.Name)
			if o.Required {
				selectOptionRequired[group] = true
			}
			v := c.String(o.Name)
			log.Debugf("found select option %s (group=%s value=%v required=%v)\n", o.Name, group, v, o.Required)
			if v == "true" {
				selectOptionSelectedValues[group] = append(selectOptionSelectedValues[group], o.Name)
			}
		case cmd.SelectOptionV2:
			group := o.Name
			if o.Required {
				selectOptionRequired[group] = true
			}
			for _, ov := range o.Values {
				selectOptions[group] = append(selectOptions[group], ov.Name)
				v := c.String(ov.Name)
				log.Debugf("found select option %s (group=%s value=%v required=%v)\n", ov.Name, group, v, o.Required)
				if v == "true" {
					selectOptionSelectedValues[group] = append(selectOptionSelectedValues[group], ov.Name)
				}
			}
		}
	}

	for group := range selectOptions {
		if selectOptionRequired[group] {
			optionValues, ok := selectOptionSelectedValues[group]
			if ok && optionValues[0] != "" {
				log.Debugf("select option group %s was set by option %s", group, optionValues[0])
			} else {
				cli.ShowCommandHelp(c, cmdName)
				return fmt.Errorf("Required %s flag missing for select option group \"%s\" (one of \" %s \" must be provided)\n", level, group, strings.Join(selectOptions[group], " | "))
			}
		} else {
			log.Debugf("select option group %s does not require a value", group)
		}

		if optionValues, ok := selectOptionSelectedValues[group]; ok && len(optionValues) > 1 {
			return fmt.Errorf("%s flag specified multiple times for select option group \"%s\" (one of \" %s \" must be provided)\n", level, group, strings.Join(selectOptions[group], " | "))
		}
	}
	return nil
}
