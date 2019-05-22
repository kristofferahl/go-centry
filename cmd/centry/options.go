package main

import (
	"github.com/kristofferahl/go-centry/pkg/cmd"
)

func createGlobalOptions(context *Context) *cmd.OptionsSet {
	manifest := context.manifest

	// Add global options
	options := cmd.NewOptionsSet(cmd.OptionSetGlobal)
	options.Add(&cmd.Option{
		Name:        "config.log.level",
		Description: "Overrides the log level",
		Default:     manifest.Config.Log.Level,
	})
	options.Add(&cmd.Option{
		Name:        "quiet",
		Short:       "q",
		Description: "Disables logging",
	})
	options.Add(&cmd.Option{
		Name:        "help",
		Short:       "h",
		Description: "Displays help",
	})
	options.Add(&cmd.Option{
		Name:        "version",
		Short:       "v",
		Description: "Displays the version of the cli",
	})

	// Adding global options specified by the manifest
	for _, o := range manifest.Options {
		o := o

		if context.optionEnabledFunc != nil && context.optionEnabledFunc(o) == false {
			continue
		}

		options.Add(&cmd.Option{
			Name:        o.Name,
			Description: o.Description,
			EnvName:     o.EnvName,
			Default:     o.Default,
		})
	}

	return options
}
