package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/kristofferahl/cli"
)

func centryHelpFunc(manifest *manifest, globalFlags *flag.FlagSet) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("Usage: %s [--version] [--help] <command> [<args>]\n\n", manifest.Config.Name))

		writeCommands(&buf, commands, manifest)
		writeOptions(&buf, "Global", globalFlags)

		return buf.String()
	}
}

func writeCommands(buf *bytes.Buffer, commands map[string]cli.CommandFactory, manifest *manifest) {
	buf.WriteString("Available commands are:\n")

	// Get the list of keys so we can sort them, and also get the maximum
	// key length so they can be aligned properly.
	keys := make([]string, 0, len(commands))
	maxKeyLen := 0
	for key := range commands {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}

		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		commandFunc, ok := commands[key]
		if !ok {
			// This should never happen since we JUST built the list of
			// keys.
			panic("command not found: " + key)
		}

		command, err := commandFunc()
		if err != nil {
			log.Printf("[ERR] cli: Command '%s' failed to load: %s", key, err)
			continue
		}

		synopsis := command.Synopsis()
		if synopsis == "" {
			for _, mc := range manifest.Commands {
				if mc.Name == key {
					synopsis = mc.Description
				}
			}
		}
		key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
		buf.WriteString(fmt.Sprintf("    %s    %s\n", key, synopsis))
	}
}

func writeOptions(buf *bytes.Buffer, group string, flags *flag.FlagSet) {
	buf.WriteString(fmt.Sprintf("\n%s options are:\n", group))

	options := make(map[string]string, 0)
	keys := make([]string, 0)
	maxKeyLen := 0

	flags.VisitAll(func(f *flag.Flag) {
		key := fmt.Sprintf("--%s", f.Name)
		if len(f.Name) == 1 {
			key = fmt.Sprintf("-%s", f.Name)
		}

		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}

		options[key] = f.Usage
		keys = append(keys, key)
	})
	sort.Strings(keys)

	for _, key := range keys {
		o := fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
		d := options[key]
		buf.WriteString(fmt.Sprintf("    %s    %s\n", o, d))
	}
}

type helpOption struct {
	Name        string
	Description string
}
