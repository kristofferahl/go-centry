package main

import (
	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
)

// OptionSetGlobal is the name of the global OptionsSet
const OptionSetGlobal = "Global"

func createGlobalOptions(context *Context) *cmd.OptionsSet {
	manifest := context.manifest

	// Add global options
	options := cmd.NewOptionsSet(OptionSetGlobal)
	options.ShortCircuitParseFunc = func(arg string) bool {
		return arg == "-v" || arg == "--v" || arg == "-version" || arg == "--version" || arg == "-h" || arg == "--h" || arg == "-help" || arg == "--help"
	}

	options.Add(&cmd.Option{
		Type:        cmd.StringOption,
		Name:        "config.log.level",
		Description: "Overrides the log level",
		Default:     manifest.Config.Log.Level,
	})
	options.Add(&cmd.Option{
		Type:        cmd.BoolOption,
		Name:        "quiet",
		Short:       "q",
		Description: "Disables logging",
		Default:     false,
	})
	options.Add(&cmd.Option{
		Type:        cmd.BoolOption,
		Name:        "help",
		Short:       "h",
		Description: "Displays help",
		Default:     false,
	})
	options.Add(&cmd.Option{
		Type:        cmd.BoolOption,
		Name:        "version",
		Short:       "v",
		Description: "Displays the version of the cli",
		Default:     false,
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
		})
	}

	return options
}
