package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

func interactive(runtime *Runtime) {
	replArgs := []string{}

	cmd, cmdArgs := promptForCommands(nil, runtime.cli.VisibleCommands(), []string{})
	if cmd != nil {
		replArgs = append(replArgs, promptForOptions(runtime.cli.VisibleFlags(), []string{})...)
		replArgs = append(replArgs, cmdArgs...)
		replArgs = append(replArgs, promptForOptions(cmd.VisibleFlags(), []string{})...)
		replArgs = trimEmpty(replArgs)
	}

	fmt.Println()

	if cmd != nil {
		exec := false
		confirm := &survey.Confirm{
			Message: fmt.Sprintf("%s %s\n  would you like to run the command above:", runtime.cli.Name, strings.Join(replArgs, " ")),
		}
		survey.AskOne(confirm, &exec)

		if exec {
			runtime.args = replArgs
		}
	}
}

func promptForCommands(parent *cli.Command, commands []*cli.Command, in []string) (cmd *cli.Command, args []string) {
	descriptions := make(map[string]string)
	values := make([]string, 0)
	for _, c := range commands {
		if !c.Hidden {
			values = append(values, c.Name)
			descriptions[c.Name] = c.Usage
		}
	}

	reply := ""
	msg := "select command:"
	if parent != nil {
		msg = "select subcommand:"
	}
	prompt := &survey.Select{
		Message: msg,
		Options: values,
		Description: func(value string, index int) string {
			return descriptions[value]
		},
	}
	survey.AskOne(prompt, &reply)

	for _, c := range commands {
		if c.Name == reply {
			return promptForCommands(c, c.Subcommands, append(in, reply))
		}
	}

	return parent, append(in, reply)
}

func promptForOptions(flags []cli.Flag, in []string) []string {
	handled := make(map[string]string)
	values := make([]string, 0)
	descriptions := make(map[string]string)
	optional := make(map[string]cli.Flag)

	for _, f := range flags {
		name := f.Names()[0]

		if rf, ok := f.(cli.RequiredFlag); ok && rf.IsRequired() {
			in = appendFlagValue(name, f, in)
		} else if sf, ok := f.(*SelectOptionFlag); ok {
			if _, ok := handled[sf.GroupName]; ok {
				continue
			}
			handled[sf.GroupName] = ""
			if sf.GroupRequired {
				in = appendFlagValue(name, f, in)
			} else {
				values = append(values, sf.GroupName)
				descriptions[sf.GroupName] = sf.GetUsage()
				optional[sf.GroupName] = f
			}
		} else {
			if df, ok := f.(cli.DocGenerationFlag); ok {
				values = append(values, name)
				descriptions[name] = df.GetUsage()
				optional[name] = f
			}
		}
	}

	selected := []string{}
	prompt := &survey.MultiSelect{
		Message: "select options to set:",
		Options: values,
		Description: func(value string, index int) string {
			return descriptions[value]
		},
		PageSize: 20,
	}
	survey.AskOne(prompt, &selected)

	for _, name := range selected {
		if f, ok := optional[name]; ok {
			in = appendFlagValue(name, f, in)
		}
	}

	return in
}

func appendFlagValue(name string, f cli.Flag, args []string) []string {
	required := "[optional] "
	if rf, ok := f.(cli.RequiredFlag); ok && rf.IsRequired() {
		required = "[required] "
	}

	switch v := f.(type) {
	case *SelectOptionFlag:
		if v.GroupRequired {
			required = "[required] "
		}
		values := make(map[string]string)
		for _, val := range v.Values {
			values[val.Name] = fmt.Sprintf("%s (%s=%s)", v.GetUsage(), v.GroupName, val.ResolveValue())
		}
		prompt := fmt.Sprintf("%soption \"%s\"", required, v.GroupName)
		val := promptSelectValue(prompt, values)
		args = append(args, fmt.Sprintf("--%s", val))
	case *cli.BoolFlag:
		args = append(args, fmt.Sprintf("--%s", name))
	case *cli.StringFlag:
		prompt := fmt.Sprintf("%soption \"%s\" (%s)", required, name, v.GetUsage())
		val := promptValue(prompt, v.GetValue())
		args = append(args, fmt.Sprintf("--%s=%s", name, val))
	case *cli.IntFlag:
		prompt := fmt.Sprintf("%soption \"%s\" (%s)", required, name, v.GetUsage())
		val := promptValue(prompt, v.GetValue())
		args = append(args, fmt.Sprintf("--%s=%s", name, val))
	default:
		panic(fmt.Errorf("unhnadled flag type, %v", f))
	}
	return args
}

func promptSelectValue(text string, values map[string]string) string {
	options := make([]string, 0)
	for k, _ := range values {
		options = append(options, k)
	}

	for {
		selectedIndex := -1
		prompt := &survey.Select{
			Message: fmt.Sprintf("select value for %s:", text),
			Options: options,
			Default: options[0],
			Description: func(value string, index int) string {
				return values[value]
			},
			PageSize: 20,
		}
		survey.AskOne(prompt, &selectedIndex)

		if selectedIndex >= 0 {
			return options[selectedIndex]
		}
	}
}

func promptValue(text string, def string) string {
	v := ""
	prompt := &survey.Input{
		Message: fmt.Sprintf("enter a value for %s:", text),
		Default: def,
	}
	if def == "" {
		survey.AskOne(prompt, &v, survey.WithValidator(survey.Required))
	} else {
		survey.AskOne(prompt, &v)
	}
	return v
}

func trimEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
