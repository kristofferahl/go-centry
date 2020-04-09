package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/franela/goblin"
	"github.com/kristofferahl/go-centry/internal/pkg/io"
	test "github.com/kristofferahl/go-centry/internal/pkg/test"
)

func TestMain(t *testing.T) {
	g := Goblin(t)

	// Esuring the workdir is the root of the repo
	os.Chdir("../../")

	g.Describe("runtime", func() {
		g.It("returns error when manifest fails to load", func() {
			context := NewContext(CLI, io.Headless())
			runtime, err := NewRuntime([]string{}, context)
			g.Assert(runtime == nil).IsTrue("expected runtime to be nil")
			g.Assert(err != nil).IsTrue("expected error")
			g.Assert(strings.HasPrefix(err.Error(), "Failed to read manifest file.")).IsTrue("expected error message")
		})
	})

	g.Describe("scripts", func() {
		g.It("loads script in the expected order", func() {
			expected := "Loading init.sh\nLoading helpers.sh"
			os.Setenv("OUTPUT_DEBUG", "true")
			out := execQuiet("scripttest")
			test.AssertStringContains(g, out.Stdout, expected)
			os.Unsetenv("OUTPUT_DEBUG")
		})
	})

	g.Describe("commands", func() {
		g.Describe("invoking command", func() {
			g.Describe("with arguments", func() {
				g.It("should have arguments passed", func() {
					expected := "command args (foo bar)"
					out := execQuiet("commandtest foo bar")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})

			g.Describe("without arguments", func() {
				g.It("should have no arguments passed", func() {
					expected := "command args ()"
					out := execQuiet("commandtest")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})
		})

		g.Describe("invoking sub command", func() {
			g.Describe("with arguments", func() {
				g.It("should have arguments passed", func() {
					expected := "subcommand args (foo bar)"
					out := execQuiet("commandtest subcommand foo bar")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})

			g.Describe("without arguments", func() {
				g.It("should have no arguments passed", func() {
					expected := "subcommand args ()"
					out := execQuiet("commandtest subcommand")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})
		})

		g.Describe("command options", func() {
			g.Describe("invoking command with options", func() {
				g.It("should have arguments passed", func() {
					expected := "command args (foo bar baz)"
					out := execQuiet("commandtest options args --cmdstringopt=hello --cmdboolopt --cmdsel1 --cmdsel2 foo bar baz")
					test.AssertStringContains(g, out.Stdout, expected)
				})

				g.It("should have environment variables set", func() {
					out := execQuiet("commandtest options printenv --cmdstringopt=world --cmdboolopt --cmdsel1 --cmdsel2")
					test.AssertStringHasKeyValue(g, out.Stdout, "CMDSTRINGOPT", "world")
					test.AssertStringHasKeyValue(g, out.Stdout, "CMDBOOLOPT", "true")
					test.AssertStringHasKeyValue(g, out.Stdout, "CMDSELECTOPT", "cmdsel2")
				})
			})
		})
	})

	g.Describe("options", func() {
		g.Describe("invoke without option", func() {
			g.It("should pass arguments", func() {
				expected := "args (foo bar)"
				out := execQuiet("optiontest args foo bar")
				test.AssertStringContains(g, out.Stdout, expected)
			})

			// TODO: Add assertions for all default values?
			g.It("should have default value for environment variable set", func() {
				out := execQuiet("optiontest printenv")
				test.AssertStringHasKeyValue(g, out.Stdout, "STRINGOPT", "foobar")
			})
		})

		g.Describe("invoke with single option", func() {
			g.It("should have arguments passed", func() {
				expected := "args (foo bar)"
				out := execQuiet("--boolopt optiontest args foo bar")
				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should have environment set for select option", func() {
				out := execQuiet("--selectopt1 optiontest printenv")
				test.AssertStringHasKeyValue(g, out.Stdout, "SELECTOPT", "selectopt1")
			})

			g.It("should have environment set to last select option with same env_name (selectopt2)", func() {
				out := execQuiet("--selectopt1 --selectopt2 optiontest printenv")
				test.AssertStringHasKeyValue(g, out.Stdout, "SELECTOPT", "selectopt2")
			})

			// TODO: Do we really need =false??
			g.It("should have environment set to last select option with same env_name (selectopt1)", func() {
				out := execQuiet("--selectopt2=false --selectopt1 optiontest printenv")
				test.AssertStringHasKeyValue(g, out.Stdout, "SELECTOPT", "selectopt1")
			})
		})

		g.Describe("invoke with multiple options", func() {
			g.It("should have arguments passed", func() {
				expected := "args (bar foo)"
				out := execQuiet("--boolopt --stringopt=foo optiontest args bar foo")
				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should have multipe environment variables set", func() {
				out := execQuiet("--selectopt2 --stringopt=blazer --boolopt optiontest printenv")

				test.AssertStringHasKeyValue(g, out.Stdout, "STRINGOPT", "blazer")
				test.AssertStringHasKeyValue(g, out.Stdout, "BOOLOPT", "true")
				test.AssertStringHasKeyValue(g, out.Stdout, "SELECTOPT", "selectopt2")
			})
		})

		g.Describe("invoke with invalid option", func() {
			g.It("should print error message", func() {
				out := execQuiet("--invalidoption optiontest args")
				test.AssertStringContains(g, out.Stdout, "Incorrect Usage. flag provided but not defined: -invalidoption")
				test.AssertStringContains(g, out.Stderr, "flag provided but not defined: -invalidoption")
				test.AssertNoError(g, out.Error)
			})
		})
	})

	g.Describe("global options", func() {
		g.Describe("version", func() {
			g.Describe("--version", func() {
				g.It("should display version", func() {
					expected := `1.0.0`
					out := execQuiet("--version")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})

			g.Describe("-v", func() {
				g.It("should display version", func() {
					expected := `1.0.0`
					out := execQuiet("-v")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})
		})

		g.Describe("quiet", func() {
			g.Describe("--quiet", func() {
				g.It("should disable logging", func() {
					expected := `Changing loglevel to panic (from debug)`
					out := execWithLogging("--quiet")
					test.AssertStringContains(g, out.Stderr, expected)
				})
			})

			g.Describe("-q", func() {
				g.It("should disable logging", func() {
					expected := `Changing loglevel to panic (from debug)`
					out := execWithLogging("-q")
					test.AssertStringContains(g, out.Stderr, expected)
				})
			})
		})

		g.Describe("--config.log.level=info", func() {
			g.It("should change log level to info", func() {
				expected := `Changing loglevel to info (from debug)`
				out := execWithLogging("--config.log.level=info")
				test.AssertStringContains(g, out.Stderr, expected)
			})
		})

		g.Describe("--config.log.level=error", func() {
			g.It("should change log level to error", func() {
				expected := `Changing loglevel to error (from debug)`
				out := execWithLogging("--config.log.level=error")
				test.AssertStringContains(g, out.Stderr, expected)
			})
		})
	})

	g.Describe("help", func() {
		g.Describe("call with no arguments", func() {
			g.It("should display help", func() {
				expected := "USAGE:"
				out := execQuiet("")
				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("call with -h", func() {
			g.It("should display help", func() {
				expected := "USAGE:"
				out := execQuiet("-h")
				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("call with --help", func() {
			g.It("should display help", func() {
				expected := "USAGE:"
				out := execQuiet("--help")
				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("output", func() {
			out := execQuiet("")
			g.It("should display available commands", func() {
				expected := `COMMANDS:
   commandtest  Command tests
   helptest     Help tests
   optiontest   Option tests
   scripttest   Script tests`

				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should display global options", func() {
				expected := `OPTIONS:
   --boolopt, -B                A custom option (default: false)
   --config.log.level value     Overrides the log level (default: "debug")
   --quiet, -q                  Disables logging (default: false)
   --selectopt1                 Sets the selection to option 1 (default: false)
   --selectopt2                 Sets the selection to option 2 (default: false)
   --stringopt value, -S value  A custom option (default: "foobar")
   --help, -h                   Show help (default: false)
   --version, -v                Print the version (default: false)`

				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("command help output", func() {
			g.It("should display full help", func() {
				expected := `NAME:
   centry helptest - Help tests

USAGE:
   centry helptest command [command options] [arguments...]

COMMANDS:
   placeholder  ...
   subcommand   Description for subcommand

OPTIONS:
   --help, -h     Show help (default: false)
   --version, -v  Print the version (default: false)`

				out := execQuiet("helptest --help")
				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("subcommand help output", func() {
			g.It("should display full help", func() {
				expected := `NAME:
   centry helptest subcommand - Description for subcommand

USAGE:
   Help text for sub command

OPTIONS:
   --opt1 value, -o value  Help text for opt1 (default: "footothebar")
   --help, -h              Show help (default: false)`

				out := execQuiet("helptest subcommand --help")
				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("placeholder help output", func() {
			g.It("should display full help", func() {
				expected := `NAME:
   centry helptest placeholder - ...

USAGE:
   centry helptest placeholder command [command options] [arguments...]

COMMANDS:
   subcommand1  Description for placeholder subcommand1
   subcommand2  Description for placeholder subcommand2

OPTIONS:
   --help, -h     Show help (default: false)
   --version, -v  Print the version (default: false)`

				out := execQuiet("helptest placeholder --help")
				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.Describe("placeholder subcommand help", func() {
				g.It("should display full help", func() {
					expected := `NAME:
   centry helptest placeholder subcommand1 - Description for placeholder subcommand1

USAGE:
   Help text for placeholder subcommand1

OPTIONS:
   --opt1 value  Help text for opt1
   --help, -h    Show help (default: false)`

					out := execQuiet("helptest placeholder subcommand1 --help")
					test.AssertStringContains(g, out.Stdout, expected)
				})
			})
		})
	})
}

type execResult struct {
	Source   string
	ExitCode int
	Error    error
	Stdout   string
	Stderr   string
}

func execQuiet(source string) *execResult {
	return execCentry(source, true)
}

func execWithLogging(source string) *execResult {
	return execCentry(source, false)
}

func execCentry(source string, quiet bool) *execResult {
	var exitCode int
	var runtimeErr error

	out := test.CaptureOutput(func() {
		if quiet {
			source = fmt.Sprintf("--quiet %s", source)
		}
		context := NewContext(CLI, io.Headless())
		runtime, err := NewRuntime(strings.Split(fmt.Sprintf("test/data/runtime_test.yaml %s", source), " "), context)
		if err != nil {
			exitCode = 1
			runtimeErr = err
		} else {
			exitCode = runtime.Execute()
		}
	})

	return &execResult{
		Source:   source,
		ExitCode: exitCode,
		Error:    runtimeErr,
		Stdout:   out.Stdout,
		Stderr:   out.Stderr,
	}
}
