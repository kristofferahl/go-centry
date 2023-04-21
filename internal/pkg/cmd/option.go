package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// OptionsSet represents a set of flags that can be passed to the cli
type OptionsSet struct {
	Name  string
	items map[string]*Option
}

// OptionType defines the type of an option
type OptionType string

// StringOption defines a string value option
const StringOption OptionType = "string"

// BoolOption defines a boolean value option
const BoolOption OptionType = "bool"

// IntegerOption defines an interger value option
const IntegerOption OptionType = "integer"

// SelectOption defines a select value option
const SelectOption OptionType = "select"

// SelectOptionV2 defines a select value option
const SelectOptionV2 OptionType = "select/v2"

// StringToOptionType returns the OptionType matching the provided string
func StringToOptionType(s string) OptionType {
	s = strings.ToLower(s)
	switch s {
	case "string":
		return StringOption
	case "bool":
		return BoolOption
	case "integer":
		return IntegerOption
	case "select":
		return SelectOption
	case "select/v2":
		return SelectOptionV2
	default:
		return StringOption
	}
}

// Option represents a flag that can be passed to the cli
type Option struct {
	Type        OptionType
	Name        string
	Short       string
	EnvName     string
	Description string
	Required    bool
	Hidden      bool
	Internal    bool
	Values      []OptionValue
	Default     interface{}
}

type OptionValue struct {
	Name  string `json:"name,omitempty"`
	Short string `json:"short,omitempty"`
	Value string `json:"value,omitempty"`
}

// Validate returns true if the option is considered valid
func (o *Option) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("missing option name")
	}

	if o.Type == "" {
		return fmt.Errorf("missing option type")
	}

	for _, ov := range o.Values {
		if ov.Name == "" {
			return fmt.Errorf("missing option value name")
		}
	}

	return nil
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

	err := option.Validate()
	if err != nil {
		return err
	}

	err = convertDefaultValueToCorrectType(option)
	if err != nil {
		return err
	}

	names := make([]string, 0)
	for k, o := range s.items {
		names = append(names, k)
		for _, ov := range o.Values {
			names = append(names, ov.Name)
		}
	}

	shortNames := make([]string, 0)
	for _, o := range s.items {
		if len(o.Short) > 0 {
			shortNames = append(shortNames, o.Short)
		}
		for _, ov := range o.Values {
			if len(ov.Short) > 0 {
				shortNames = append(shortNames, ov.Short)
			}
		}
	}

	if contains(names, option.Name) {
		return fmt.Errorf("an option with the name \"%s\" has already been added", option.Name)
	}

	if len(option.Short) > 0 && contains(shortNames, option.Short) {
		return fmt.Errorf("an option with the short name \"%s\" has already been added", option.Short)
	}

	for _, ov := range option.Values {
		if contains(names, ov.Name) {
			return fmt.Errorf("an option value with the name \"%s\" has already been added", ov.Name)
		}

		if len(ov.Short) > 0 && contains(shortNames, ov.Short) {
			return fmt.Errorf("an option with the short name \"%s\" has already been added", ov.Short)
		}
	}

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

func convertDefaultValueToCorrectType(option *Option) error {
	var def interface{}

	switch option.Type {
	case SelectOption:
		def = false
	case SelectOptionV2:
		def = false
	case IntegerOption:
		def = 0
		switch option.Default.(type) {
		case string:
			if option.Default != "" {
				val, err := strconv.Atoi(option.Default.(string))
				if err != nil {
					return err
				}
				def = val
			}
		}
	case BoolOption:
		def = false
	case StringOption:
		def = option.Default
	default:
		return fmt.Errorf("default value conversion not registered for type \"%s\"", option.Type)
	}

	option.Default = def

	return nil
}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
