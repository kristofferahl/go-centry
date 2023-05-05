package main

import (
	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	"github.com/urfave/cli/v2"
)

type SelectOptionFlag struct {
	cli.BoolFlag

	GroupName     string
	GroupRequired bool
	Values        []cmd.OptionValue
}
