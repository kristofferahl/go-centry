package cmd

import (
	"flag"
	"strconv"
	"strings"
)

// OptionSetGlobal is the name of the global OptionsSet
const OptionSetGlobal = "Global"

// OptionsSet represents a set of flags that can be passed to the cli
type OptionsSet struct {
	Name  string
	Items map[string]*Option
	Flags *flag.FlagSet
}

// Option represents a flag that can be passed to the cli
type Option struct {
	Name        string
	Short       string
	EnvName     string
	Default     string
	Description string
	Value       value
}

type boolValue bool

func (b *boolValue) string() string { return strconv.FormatBool(bool(*b)) }

type stringValue string

func (s *stringValue) string() string { return string(*s) }

type value interface {
	string() string
}

// NewOptionsSet creates a new set of options
func NewOptionsSet(name string) *OptionsSet {
	return &OptionsSet{
		Name:  name,
		Items: make(map[string]*Option, 0),
	}
}

// Add adds options to the set
func (s *OptionsSet) Add(option *Option) {
	s.Items[option.Name] = option
}

// CreateFlagSet returns the set of options as a FlagSet
func (s *OptionsSet) CreateFlagSet() {
	s.Flags = flag.NewFlagSet(s.Name, flag.ContinueOnError)

	for _, o := range s.Items {
		d := strings.ToLower(o.Default)
		if d == "" {
			f := s.Flags.Bool(o.Name, false, o.Description)
			if o.Short != "" {
				s.Flags.BoolVar(f, o.Short, false, o.Description)
			}
			o.Value = (*boolValue)(f)
			continue
		}

		f := s.Flags.String(o.Name, o.Default, o.Description)
		if o.Short != "" {
			s.Flags.StringVar(f, o.Short, o.Default, o.Description)
		}
		o.Value = (*stringValue)(f)
	}
}

// Parse pareses the args using a flagset and returns the remaining arguments
func (s *OptionsSet) Parse(args []string) ([]string, error) {
	parse := true
	if s.Name == OptionSetGlobal {
		for _, arg := range args {
			if arg == "-v" || arg == "--v" || arg == "-version" || arg == "--version" {
				parse = false
				break
			}
			if arg == "-h" || arg == "--h" || arg == "-help" || arg == "--help" {
				parse = false
				break
			}
		}
	}

	if parse {
		s.CreateFlagSet()
		err := s.Flags.Parse(args)
		if err != nil {
			return nil, err
		}
		args = s.Flags.Args()
	}

	return args, nil
}

// GetValue returns the parsed value of a given option
func (s *OptionsSet) GetValue(key string) string {
	option := s.Items[key]
	result := ""

	if option.Value != nil {
		switch v := s.Items[key].Value.(type) {
		case *boolValue:
			result = (*v).string()
		case *stringValue:
			result = (*v).string()
		}
	}

	return result
}

// GetBool returns the parsed value of a given option
func (s *OptionsSet) GetBool(key string) bool {
	option := s.Items[key]
	result := false

	if option.Value != nil {
		switch v := s.Items[key].Value.(type) {
		case *boolValue:
			result = bool(*v)
		}
	}

	return result
}

// GeString returns the parsed value of a given option
func (s *OptionsSet) GeString(key string) string {
	option := s.Items[key]
	result := ""

	if option.Value != nil {
		switch v := s.Items[key].Value.(type) {
		case *stringValue:
			result = string(*v)
		}
	}

	return result
}
