package cmd

import (
	"fmt"
	"sort"
	"strconv"
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

// Validate returns true if the option is concidered valid
func (o *Option) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("missing option name")
	}

	if o.Type == "" {
		return fmt.Errorf("missing option type")
	}

	return nil
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

	err := option.Validate()
	if err != nil {
		return err
	}

	if _, ok := s.items[option.Name]; ok {
		return fmt.Errorf("an option with the name \"%s\" has already been added", option.Name)
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
