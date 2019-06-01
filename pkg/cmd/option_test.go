package cmd

import (
	"flag"
	"testing"

	. "github.com/franela/goblin"
	"github.com/kristofferahl/go-centry/pkg/io"
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

		g.Describe("AsFlagSet", func() {
			g.Describe("with no options", func() {
				os := NewOptionsSet("Name")
				fs := os.AsFlagSet()

				g.It("should have error handling turned off", func() {
					g.Assert(fs.ErrorHandling()).Equal(flag.ContinueOnError)
				})

				g.It("should have no flags", func() {
					c := 0
					fs.VisitAll(func(flag *flag.Flag) {
						c++
					})
					g.Assert(c).Equal(0)
				})
			})

			g.Describe("with options", func() {
				os := NewOptionsSet("Name")
				os.Add(&Option{Type: StringOption, Name: "Option"})
				fs := os.AsFlagSet()

				g.It("should have error handling turned off", func() {
					g.Assert(fs.ErrorHandling()).Equal(flag.ContinueOnError)
				})

				g.It("should have flags", func() {
					c := 0
					var f *flag.Flag
					fs.VisitAll(func(flag *flag.Flag) {
						c++
						f = flag
					})
					g.Assert(c).Equal(1)
					g.Assert(f.Name).Equal("Option")
				})
			})
		})

		g.Describe("Parse", func() {
			g.Describe("with no options", func() {
				os := NewOptionsSet("Name")

				g.Describe("passing nil as args", func() {
					g.It("should return 0 args", func() {
						rest, err := os.Parse(nil, io.Headless())
						g.Assert(err).Equal(nil)
						g.Assert(len(rest)).Equal(0)
					})
				})

				g.Describe("passing 0 args", func() {
					g.It("should return 0 args", func() {
						rest, err := os.Parse([]string{}, io.Headless())
						g.Assert(err).Equal(nil)
						g.Assert(len(rest)).Equal(0)
					})
				})

				g.Describe("passing args", func() {
					g.It("should return all args", func() {
						rest, err := os.Parse([]string{"a1", "a2", "a3"}, io.Headless())
						g.Assert(err).Equal(nil)
						g.Assert(len(rest)).Equal(3)
					})
				})
			})

			g.Describe("boolean options", func() {
				withBooleanOption := func(defaultValue bool) *OptionsSet {
					os := NewOptionsSet("Name")
					os.Add(&Option{
						Type:    BoolOption,
						Name:    "Option",
						Default: defaultValue,
					})
					return os
				}

				g.Describe("with default value false", func() {
					g.Describe("passing nil as args", func() {
						os := withBooleanOption(false)
						rest, err := os.Parse(nil, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should not have value for otpions", func() {
							g.Assert(os.GetBool("Option")).Equal(false)
						})
					})

					g.Describe("passing flag with value", func() {
						os := withBooleanOption(false)
						rest, err := os.Parse([]string{"--Option=true", "arg"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(1)
							g.Assert(rest[0]).Equal("arg")
						})

						g.It("should not have value for option", func() {
							g.Assert(os.GetBool("Option")).Equal(true)
						})
					})

					g.Describe("passing flag without value", func() {
						os := withBooleanOption(false)
						rest, err := os.Parse([]string{"--Option", "arg"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(1)
							g.Assert(rest[0]).Equal("arg")
						})

						g.It("should not have value for option", func() {
							g.Assert(os.GetBool("Option")).Equal(true)
						})
					})
				})

				g.Describe("with default value true", func() {
					g.Describe("passing nil as args", func() {
						os := withBooleanOption(true)
						rest, err := os.Parse(nil, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should not have value for otpions", func() {
							g.Assert(os.GetBool("Option")).Equal(true)
						})
					})

					g.Describe("passing flag with value", func() {
						os := withBooleanOption(true)
						rest, err := os.Parse([]string{"--Option=false", "arg"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(1)
							g.Assert(rest[0]).Equal("arg")
						})

						g.It("should not have value for option", func() {
							g.Assert(os.GetBool("Option")).Equal(false)
						})
					})

					g.Describe("passing flag without value", func() {
						os := withBooleanOption(true)
						rest, err := os.Parse([]string{"--Option", "arg"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(1)
							g.Assert(rest[0]).Equal("arg")
						})

						g.It("should not have value for option", func() {
							g.Assert(os.GetBool("Option")).Equal(true)
						})
					})
				})
			})

			g.Describe("string options", func() {
				withStringOption := func(defaultValue string) *OptionsSet {
					os := NewOptionsSet("Name")
					os.Add(&Option{
						Type:    StringOption,
						Name:    "Option",
						Default: defaultValue,
					})
					return os
				}

				g.Describe("without default value", func() {
					g.Describe("passing nil as args", func() {
						os := withStringOption("")
						rest, err := os.Parse(nil, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should not have value for otpions", func() {
							g.Assert(os.GetString("Option")).Equal("")
						})
					})

					g.Describe("passing flag with value", func() {
						os := withStringOption("")
						rest, err := os.Parse([]string{"--Option", "value"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should not have value for option", func() {
							g.Assert(os.GetString("Option")).Equal("value")
						})
					})
				})

				g.Describe("string option with default value", func() {
					g.Describe("passing nil", func() {
						os := withStringOption("DefaultValue")
						rest, err := os.Parse(nil, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should not have value for otpions", func() {
							g.Assert(os.GetString("Option")).Equal("DefaultValue")
						})
					})

					g.Describe("passing flag with value", func() {
						os := withStringOption("DefaultValue")
						rest, err := os.Parse([]string{"--Option", "value"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should override default value", func() {
							g.Assert(os.GetString("Option")).Equal("value")
						})
					})
				})
			})

			g.Describe("select options", func() {
				withSelectOptions := func(defaultValue bool) *OptionsSet {
					os := NewOptionsSet("Name")
					os.Add(&Option{
						Type:    SelectOption,
						Name:    "Option1",
						EnvName: "OneOf",
						Default: defaultValue,
					})
					os.Add(&Option{
						Type:    SelectOption,
						Name:    "Option2",
						EnvName: "OneOf",
						Default: defaultValue,
					})
					os.Add(&Option{
						Type:    SelectOption,
						Name:    "Option3",
						EnvName: "OneOf",
						Default: defaultValue,
					})
					return os
				}

				g.Describe("passing nil", func() {
					os := withSelectOptions(false)
					rest, err := os.Parse(nil, io.Headless())

					g.It("should return 0 args and no error", func() {
						g.Assert(err).Equal(nil)
						g.Assert(len(rest)).Equal(0)
					})

					g.It("should not have default value for otpions", func() {
						g.Assert(os.GetBool("Option1")).Equal(false)
						g.Assert(os.GetBool("Option2")).Equal(false)
					})
				})

				g.Describe("passing flag should", func() {
					g.Describe("with no value", func() {
						os := withSelectOptions(false)
						rest, err := os.Parse([]string{"--Option1"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should set value", func() {
							g.Assert(os.GetBool("Option1")).Equal(true)
							g.Assert(os.GetBool("Option2")).Equal(false)
						})
					})

					g.Describe("with value false", func() {
						os := withSelectOptions(false)
						rest, err := os.Parse([]string{"--Option1=false"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should set value", func() {
							g.Assert(os.GetBool("Option1")).Equal(false)
							g.Assert(os.GetBool("Option2")).Equal(false)
						})
					})

					g.Describe("with value true", func() {
						os := withSelectOptions(false)
						rest, err := os.Parse([]string{"--Option1=true"}, io.Headless())

						g.It("should return 0 args and no error", func() {
							g.Assert(err).Equal(nil)
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should set value", func() {
							g.Assert(os.GetBool("Option1")).Equal(true)
							g.Assert(os.GetBool("Option2")).Equal(false)
						})
					})
				})

				g.Describe("passing multiple flags", func() {
					g.Describe("when default is false", func() {
						os := withSelectOptions(false)
						rest, err := os.Parse([]string{"--Option3=false", "--Option2=true", "--Option1"}, io.Headless())

						g.It("should return 0 args and error", func() {
							g.Assert(err != nil).IsTrue("expected error")
							g.Assert(err.Error()).Equal("ambiguous flag usage [Option1 Option2]")
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should set value", func() {
							g.Assert(os.GetBool("Option1")).Equal(true)
							g.Assert(os.GetBool("Option2")).Equal(true)
							g.Assert(os.GetBool("Option3")).Equal(false)
						})
					})

					g.Describe("when default is true", func() {
						os := withSelectOptions(true)
						rest, err := os.Parse([]string{"--Option3=false", "--Option2", "--Option1=false"}, io.Headless())

						g.It("should return 0 args and error", func() {
							g.Assert(err != nil).IsTrue("expected error")
							g.Assert(err.Error()).Equal("ambiguous flag usage [Option1 Option3]")
							g.Assert(len(rest)).Equal(0)
						})

						g.It("should set value", func() {
							g.Assert(os.GetBool("Option1")).Equal(false)
							g.Assert(os.GetBool("Option2")).Equal(true)
							g.Assert(os.GetBool("Option3")).Equal(false)
						})
					})
				})
			})
		})
	})
}
