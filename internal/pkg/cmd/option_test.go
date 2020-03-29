package cmd

import (
	"testing"

	. "github.com/franela/goblin"
)

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
				os.Add(&Option{Name: "Option"})
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

			g.It("should return error when option name already exists", func() {
				os := NewOptionsSet("Name")
				err1 := os.Add(&Option{Name: "Option"})
				err2 := os.Add(&Option{Name: "Option"})
				g.Assert(len(os.Sorted())).Equal(1)
				g.Assert(err1).Equal(nil)
				g.Assert(err2 != nil).IsTrue("expected an error")
				g.Assert(err2.Error()).Equal("an option with the name \"Option\" has already been added")
			})
		})
	})
}
