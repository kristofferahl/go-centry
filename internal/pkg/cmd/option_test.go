package cmd

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/franela/goblin"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain(t *testing.T) {
	g := Goblin(t)

	g.Describe("OptionsSet", func() {
		g.Describe("NewOptionsSet", func() {
			g.It("should create a named, empty OptionsSet", func() {
				os := NewOptionsSet("Name")
				g.Assert(os.Name).Equal("Name")
				g.Assert(len(os.Sorted())).Equal(0)
			})
		})

		g.Describe("Add", func() {
			g.It("should add option", func() {
				os := NewOptionsSet("Name")
				os.Add(&Option{Name: "Option", Type: StringOption})
				g.Assert(len(os.Sorted())).Equal(1)
			})

			g.It("should return error when option is nil", func() {
				os := NewOptionsSet("Name")
				err := os.Add(nil)
				g.Assert(err != nil).IsTrue("expected an error")
				g.Assert(err.Error()).Equal("an option is required")
				g.Assert(len(os.Sorted())).Equal(0)
			})

			g.It("should return error when option name is unset", func() {
				os := NewOptionsSet("Name")
				err := os.Add(&Option{})
				g.Assert(len(os.Sorted())).Equal(0)
				g.Assert(err.Error()).Equal("missing option name")
			})

			g.It("should return error when option type is unset", func() {
				os := NewOptionsSet("Name")
				err := os.Add(&Option{Name: "foo"})
				g.Assert(len(os.Sorted())).Equal(0)
				g.Assert(err.Error()).Equal("missing option type")
			})

			g.It("should return error when option name already exists", func() {
				os := NewOptionsSet("Name")
				err1 := os.Add(&Option{Name: "Option", Type: StringOption})
				err2 := os.Add(&Option{Name: "Option", Type: StringOption})
				g.Assert(len(os.Sorted())).Equal(1)
				g.Assert(err1).Equal(nil)
				g.Assert(err2 != nil).IsTrue("expected an error")
				g.Assert(err2.Error()).Equal("an option with the name \"Option\" has already been added")
			})

			g.It("should return error when select option value name already exists as option name", func() {
				os := NewOptionsSet("Name")
				err1 := os.Add(&Option{Name: "Option", Type: StringOption})
				err2 := os.Add(&Option{Name: "Foo", Type: SelectOptionV2, Values: []OptionValue{{Name: "Option"}}})
				g.Assert(len(os.Sorted())).Equal(1)
				g.Assert(err1).Equal(nil)
				g.Assert(err2 != nil).IsTrue("expected an error")
				g.Assert(err2.Error()).Equal("an option value with the name \"Option\" has already been added")
			})

			g.It("should return error when select option value name already exists as option value name", func() {
				os := NewOptionsSet("Name")
				err1 := os.Add(&Option{Name: "Foo", Type: SelectOptionV2, Values: []OptionValue{{Name: "Opt1"}, {Name: "Opt2"}}})
				err2 := os.Add(&Option{Name: "Bar", Type: SelectOptionV2, Values: []OptionValue{{Name: "Opt2"}, {Name: "Opt3"}}})
				g.Assert(len(os.Sorted())).Equal(1)
				g.Assert(err1).Equal(nil)
				g.Assert(err2 != nil).IsTrue("expected an error")
				g.Assert(err2.Error()).Equal("an option value with the name \"Opt2\" has already been added")
			})
		})
	})

	g.Describe("StringToOptionType", func() {
		g.It("should default to StringOption", func() {
			g.Assert(StringToOptionType(randomString(10))).Equal(StringOption)
		})

		g.It("should return StringOption", func() {
			g.Assert(StringToOptionType("string")).Equal(StringOption)
			g.Assert(StringToOptionType("String")).Equal(StringOption)
			g.Assert(StringToOptionType("STRING")).Equal(StringOption)
		})

		g.It("should return BoolOption", func() {
			g.Assert(StringToOptionType("bool")).Equal(BoolOption)
			g.Assert(StringToOptionType("Bool")).Equal(BoolOption)
			g.Assert(StringToOptionType("BOOL")).Equal(BoolOption)
		})

		g.It("should return IntegerOption", func() {
			g.Assert(StringToOptionType("integer")).Equal(IntegerOption)
			g.Assert(StringToOptionType("Integer")).Equal(IntegerOption)
			g.Assert(StringToOptionType("INTEGER")).Equal(IntegerOption)
		})

		g.It("should return SelectOption", func() {
			g.Assert(StringToOptionType("select")).Equal(SelectOption)
			g.Assert(StringToOptionType("Select")).Equal(SelectOption)
			g.Assert(StringToOptionType("SELECT")).Equal(SelectOption)
		})

		g.It("should return SelectOption v2", func() {
			g.Assert(StringToOptionType("select/v2")).Equal(SelectOptionV2)
			g.Assert(StringToOptionType("Select/V2")).Equal(SelectOptionV2)
			g.Assert(StringToOptionType("SELECT/V2")).Equal(SelectOptionV2)
		})
	})
}
