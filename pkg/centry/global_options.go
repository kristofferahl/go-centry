package centry

import "github.com/kristofferahl/go-centry/pkg/config"

func createGlobalOptions(manifest *config.Manifest) *OptionsSet {
	// Add global options
	options := NewOptionsSet(optionSetGlobal)
	options.Add(&Option{
		Name:        "config.log.level",
		Description: "Overrides the log level",
		Default:     manifest.Config.Log.Level,
	})
	options.Add(&Option{
		Name:        "quiet",
		Short:       "q",
		Description: "Disables logging",
	})
	options.Add(&Option{
		Name:        "help",
		Short:       "h",
		Description: "Displays help",
	})
	options.Add(&Option{
		Name:        "version",
		Short:       "v",
		Description: "Displays the version fo the cli",
	})

	// Adding global options specified by the manifest
	for _, o := range manifest.Options {
		o := o
		options.Add(&Option{
			Name:        o.Name,
			Description: o.Description,
			EnvName:     o.EnvName,
			Default:     o.Default,
		})
	}

	return options
}
