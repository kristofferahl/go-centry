package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"strconv"

	"github.com/kristofferahl/go-centry/pkg/io"
)

// OptionsSet represents a set of flags that can be passed to the cli
type OptionsSet struct {
	Name                  string
	items                 map[string]*Option
	ShortCircuitParseFunc func(arg string) bool
	flags                 *flag.FlagSet
}

// OptionType defines the type of an option
type OptionType string

// StringOption defines a string value option
const StringOption OptionType = "string"

// BoolOption defines a boolean value option
const BoolOption OptionType = "bool"

// SelectOption defines a boolean select value option
const SelectOption OptionType = "select"

// Option represents a flag that can be passed to the cli
type Option struct {
	Type        OptionType
	Name        string
	Short       string
	EnvName     string
	Description string
	Default     interface{}
	value       valuePointer
}

type boolValue bool

func (b *boolValue) string() string { return strconv.FormatBool(bool(*b)) }

type stringValue string

func (s *stringValue) string() string { return string(*s) }

type valuePointer interface {
	string() string
}

// NewOptionsSet creates a new set of options
func NewOptionsSet(name string) *OptionsSet {
	return &OptionsSet{
		Name:  name,
		items: make(map[string]*Option, 0),
	}
}

// Add adds options to the set
func (s *OptionsSet) Add(option *Option) error {
	if option == nil {
		return fmt.Errorf("an option is required")
	}

	if option.Name == "" {
		return fmt.Errorf("missing option name")
	}

	if _, ok := s.items[option.Name]; ok {
		return fmt.Errorf("an option with the name \"%s\" has already been added", option.Name)
	}

	// TODO: Validate option type

	s.items[option.Name] = option

	return nil
}

// Sorted returns the options sorted by it's key
func (s *OptionsSet) Sorted() []*Option {
	keys := make([]string, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	options := make([]*Option, 0)
	for _, key := range keys {
		options = append(options, s.items[key])
	}

	return options
}

// AsFlagSet returns the set of options as a FlagSet
func (s *OptionsSet) AsFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(s.Name, flag.ContinueOnError)

	for _, o := range s.Sorted() {
		o := o
		switch o.Type {
		case StringOption:
			val := o.Default
			if val == nil {
				val = ""
			}
			newStringFlag(fs, o, val.(string))
		case BoolOption:
			newBoolFlag(fs, o, o.Default.(bool))
		case SelectOption:
			newBoolFlag(fs, o, o.Default.(bool))
		default:
			// TODO: Handle unsupported type
		}
	}

	return fs
}

func newStringFlag(fs *flag.FlagSet, o *Option, def string) {
	f := fs.String(o.Name, def, o.Description)
	if o.Short != "" {
		fs.StringVar(f, o.Short, def, o.Description)
	}
	o.value = (*stringValue)(f)
}

func newBoolFlag(fs *flag.FlagSet, o *Option, def bool) {
	f := fs.Bool(o.Name, def, o.Description)
	if o.Short != "" {
		fs.BoolVar(f, o.Short, def, o.Description)
	}
	o.value = (*boolValue)(f)
}

// Parse pareses the args using a flagset and returns the remaining arguments
func (s *OptionsSet) Parse(args []string, io io.InputOutput) ([]string, error) {
	parse := true

	if s.ShortCircuitParseFunc != nil {
		for _, arg := range args {
			if s.ShortCircuitParseFunc(arg) {
				parse = false
				break
			}
		}
	}

	if parse {
		s.flags = s.AsFlagSet()
		s.flags.SetOutput(bytes.NewBufferString(""))

		err := s.flags.Parse(args)
		if err != nil {
			return nil, err
		}

		soc := make(map[string][]string)

		for _, o := range s.Sorted() {
			if o.Type == SelectOption && s.GetBool(o.Name) != o.Default.(bool) {
				key := o.EnvName
				if key == "" {
					key = o.Name
				}

				if _, ok := soc[key]; !ok {
					soc[key] = make([]string, 0)
				}

				soc[key] = append(soc[key], o.Name)

				if len(soc[key]) > 1 {
					return nil, fmt.Errorf("ambiguous flag usage %v", soc[key])
				}
			}
		}

		args = s.flags.Args()
	}

	return args, nil
}

// GetValueString returns the value of a given option
func (s *OptionsSet) GetValueString(key string) string {
	option := s.items[key]
	result := ""

	if option.value != nil {
		result = option.value.string()
	}

	return result
}

// GetBool returns the parsed value of a given option
func (s *OptionsSet) GetBool(key string) bool {
	option := s.items[key]
	result := false

	if option.value != nil {
		switch v := s.items[key].value.(type) {
		case *boolValue:
			result = bool(*v)
		}
	}

	return result
}

// GetString returns the parsed value of a given option
func (s *OptionsSet) GetString(key string) string {
	option := s.items[key]
	result := ""

	if option.value != nil {
		switch v := s.items[key].value.(type) {
		case *stringValue:
			result = string(*v)
		}
	}

	return result
}
