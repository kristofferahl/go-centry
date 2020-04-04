package main

import (
	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	"github.com/urfave/cli/v2"
)

// OptionSetGlobal is the name of the global OptionsSet
const OptionSetGlobal = "Global"

func createGlobalOptions(context *Context) *cmd.OptionsSet {
	manifest := context.manifest

	// Add global options
	options := cmd.NewOptionsSet(OptionSetGlobal)

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

	// TODO: Override default version and help flags to get unified descriptions?

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

func toCliFlags(options *cmd.OptionsSet) []cli.Flag {
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
			})
		}
	}

	return flags
}
