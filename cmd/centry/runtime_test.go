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
		g.Describe("manifest file", func() {
			g.It("tries to use ./centry.yaml as the default file", func() {
				context := NewContext(CLI, io.Headless())
				runtime, err := NewRuntime([]string{}, context)
				g.Assert(runtime == nil).IsTrue("expected runtime to be nil, %v", runtime)
				g.Assert(err != nil).IsTrue("expected error, %v", err)
				g.Assert(err.Error()).Eql("manifest file not found (path=./centry.yaml)")
			})

			g.It("tries to use file specified CENTRY_FILE environment variable", func() {
				os.Setenv("CENTRY_FILE", "./centry-environment.yaml")
				context := NewContext(CLI, io.Headless())
				runtime, err := NewRuntime([]string{}, context)
				g.Assert(runtime == nil).IsTrue("expected runtime to be nil, %v", runtime)
				g.Assert(err != nil).IsTrue("expected error, %v", err)
				g.Assert(err.Error()).Eql("manifest file not found (path=./centry-environment.yaml)")
				os.Setenv("CENTRY_FILE", "")
			})

			g.It("tries to use file specified by --centry-file flag", func() {
				context := NewContext(CLI, io.Headless())
				runtime, err := NewRuntime([]string{"--centry-file", "./centry-flag.yaml"}, context)
				g.Assert(runtime == nil).IsTrue("expected runtime to be nil, %v", runtime)
				g.Assert(err != nil).IsTrue("expected error, %v", err)
				g.Assert(err.Error()).Eql("manifest file not found (path=./centry-flag.yaml)")
			})

			g.It("tries to use file specified by --centry-file= flag", func() {
				context := NewContext(CLI, io.Headless())
				runtime, err := NewRuntime([]string{"--centry-file=./centry-flag-equals.yaml"}, context)
				g.Assert(runtime == nil).IsTrue("expected runtime to be nil, %v", runtime)
				g.Assert(err != nil).IsTrue("expected error, %v", err)
				g.Assert(err.Error()).Eql("manifest file not found (path=./centry-flag-equals.yaml)")
			})

			g.It("tries to use file specified by --centry-file= flag even when it contains equal signs", func() {
				context := NewContext(CLI, io.Headless())
				runtime, err := NewRuntime([]string{"--centry-file=./foo=bar.yaml"}, context)
				g.Assert(runtime == nil).IsTrue("expected runtime to be nil, %v", runtime)
				g.Assert(err != nil).IsTrue("expected error, %v", err)
				g.Assert(err.Error()).Eql("manifest file not found (path=./foo=bar.yaml)")
			})

			g.It("gives an error when --centry-file flag is missing it's value", func() {
				context := NewContext(CLI, io.Headless())
				runtime, err := NewRuntime([]string{"--centry-file"}, context)
				g.Assert(runtime == nil).IsTrue("expected runtime to be nil, %v", runtime)
				g.Assert(err != nil).IsTrue("expected error, %v", err)
				g.Assert(err.Error()).Eql("a value must be specified for --centry-file")
			})
		})
	})

	g.Describe("scripts", func() {
		g.It("loads script in the expected order", func() {
			expected := "Loading init.sh\nLoading helpers.sh"
			os.Setenv("OUTPUT_DEBUG", "true")
			out := execQuiet("scripttest")
			test.AssertNoError(g, out.Error)
			test.AssertStringContains(g, out.Stdout, expected)
			os.Unsetenv("OUTPUT_DEBUG")
		})
	})

	g.Describe("commands", func() {
		g.Describe("invoking invalid command", func() {
			g.It("should exit with status code 127", func() {
				out := execQuiet("commandnotdefined")
				g.Assert(out.ExitCode).Equal(127)
			})
		})

		g.Describe("invoking command that exits with a status code", func() {
			g.It("should exit with exit code from command", func() {
				out := execQuiet("commandtest exitcode")
				g.Assert(out.ExitCode).Equal(111)
			})
		})

		g.Describe("invoking command with undefined option", func() {
			g.It("should exit with exit code", func() {
				out := execWithLogging("commandtest --undef")
				g.Assert(out.ExitCode).Equal(127)
			})
		})

		g.Describe("invoking command", func() {
			g.Describe("with arguments", func() {
				g.It("should have arguments passed", func() {
					expected := "command args (foo bar)"
					out := execQuiet("commandtest foo bar")
					test.AssertStringContains(g, out.Stdout, expected)
				})

				g.It("should pass any flags followed by -- as arguments", func() {
					expected := "command args (--foo bar)"
					out := execQuiet("commandtest -- --foo bar")
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

				g.It("should pass any flags followed by -- as arguments", func() {
					expected := "subcommand args (--foo bar)"
					out := execQuiet("commandtest subcommand -- --foo bar")
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
					out := execQuiet("commandtest options printenv --cmdstringopt=world --cmdboolopt --cmdsel1 --cmdsel2 --dashed-opt dashed-val")
					test.AssertStringHasKeyValue(g, out.Stdout, "CMDSTRINGOPT", "world")
					test.AssertStringHasKeyValue(g, out.Stdout, "CMDBOOLOPT", "true")
					test.AssertStringHasKeyValue(g, out.Stdout, "CMDSELECTOPT", "cmdsel2")
					test.AssertStringHasKeyValue(g, out.Stdout, "DASHED_OPT", "dashed-val")
				})

				g.It("should hav prefixed environment variables set", func() {
					out := execCentry("commandtest options printenv --cmdstringopt=world --cmdboolopt --cmdsel1 --cmdsel2", true, "test/data/runtime_test_environment_prefix.yaml")
					test.AssertStringHasKeyValue(g, out.Stdout, "ENV_PREFIX_CMDSTRINGOPT", "world")
					test.AssertStringHasKeyValue(g, out.Stdout, "ENV_PREFIX_CMDBOOLOPT", "true")
					test.AssertStringHasKeyValue(g, out.Stdout, "ENV_PREFIX_CMDSELECTOPT", "cmdsel2")
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

			g.It("should have default value for environment variable set", func() {
				out := execQuiet("optiontest printenv")
				test.AssertStringHasKeyValue(g, out.Stdout, "BOOLOPT", "false")
				test.AssertStringHasKeyValue(g, out.Stdout, "STRINGOPT", "foobar")
			})

			g.It("should have prefixed environment variables set", func() {
				out := execCentry("optiontest printenv", true, "test/data/runtime_test_environment_prefix.yaml")
				test.AssertStringHasKeyValue(g, out.Stdout, "ENV_PREFIX_BOOLOPT", "false")
				test.AssertStringHasKeyValue(g, out.Stdout, "ENV_PREFIX_STRINGOPT", "foobar")
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

		g.Describe("invoke without required option", func() {
			g.Describe("of type string", func() {
				g.It("should fail with error message", func() {
					out := execCentry("optiontest required --boolopt --selectopt1", false, "test/data/runtime_test.yaml")
					test.AssertStringContains(g, out.Stderr, "level=error msg=\"Required flag \\\"stringopt\\\" not set\"")
				})
			})
			g.Describe("of type bool", func() {
				g.It("should fail with error message", func() {
					out := execCentry("optiontest required --stringopt=foo --selectopt1", false, "test/data/runtime_test.yaml")
					test.AssertStringContains(g, out.Stderr, "level=error msg=\"Required flag \\\"boolopt\\\" not set\"")
				})
			})
			g.Describe("of type select", func() {
				g.It("should fail with error message", func() {
					out := execCentry("optiontest required --stringopt=foo --boolopt", false, "test/data/runtime_test.yaml")
					test.AssertStringContains(g, out.Stderr, "level=error msg=\"Required command flag missing for select option group SELECT (one of \\\" selectopt1 | selectopt2 \\\" must be provided)")
				})
			})
		})

		g.Describe("invoke with required option", func() {
			g.It("should pass", func() {
				out := execCentry("optiontest required --stringopt=foo --boolopt --selectopt1", false, "test/data/runtime_test.yaml")
				test.AssertStringHasKeyValue(g, out.Stdout, "STRINGOPT", "foo")
				test.AssertStringHasKeyValue(g, out.Stdout, "BOOLOPT", "true")
				test.AssertStringHasKeyValue(g, out.Stdout, "SELECTOPT", "selectopt1")
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

		g.Describe("--centry-quiet", func() {
			g.It("should disable logging", func() {
				expected := `changing loglevel to panic (from debug)`
				out := execWithLogging("--centry-quiet")
				test.AssertStringContains(g, out.Stderr, expected)
			})

			g.It("should have environment variable set", func() {
				out := execCentry("optiontest printenv", true, defaultManifestPath)
				test.AssertStringHasKeyValue(g, out.Stdout, "CENTRY_QUIET", "true")
			})

			g.It("should not have prefixed environment variable set", func() {
				out := execCentry("optiontest printenv", true, "test/data/runtime_test_environment_prefix.yaml")
				test.AssertStringHasKeyValue(g, out.Stdout, "CENTRY_QUIET", "true") // Make sure we don't prefix internal options
			})
		})

		g.Describe("--centry-config-log-level=info", func() {
			g.It("should change log level to info", func() {
				expected := `changing loglevel to info (from debug)`
				out := execWithLogging("--centry-config-log-level=info")
				test.AssertStringContains(g, out.Stderr, expected)
			})
		})

		g.Describe("--centry-config-log-level=error", func() {
			g.It("should change log level to error", func() {
				expected := `changing loglevel to error (from debug)`
				out := execWithLogging("--centry-config-log-level=error")
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

			g.It("should display the program name", func() {
				expected := `NAME:
   centry`
				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should display the program description", func() {
				expected := "A manifest file used for testing purposes"
				test.AssertStringContains(g, out.Stdout, expected)
			})

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
   --selectopt1                 Sets the selection to option 1 (default: false)
   --selectopt2                 Sets the selection to option 2 (default: false)
   --stringopt value, -S value  A custom option (default: "foobar")
   --help, -h                   Show help (default: false)
   --version, -v                Print the version (default: false)`

				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.Describe("default config output", func() {
				g.It("should display the default program description", func() {
					expected := `NAME:
   name - A declarative cli built using centry`
					out := execQuiet("", "test/data/runtime_test_default_config.yaml")
					test.AssertNoError(g, out.Error)
					test.AssertStringContains(g, out.Stdout, expected)
				})
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
   --help, -h  Show help (default: false)`

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
   --help, -h  Show help (default: false)`

				out := execQuiet("helptest placeholder --help")
				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("placeholder subcommand help output", func() {
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

		g.Describe("hidden commands", func() {
			g.It("should not display internal or hidden commands", func() {
				out := execQuiet("", "test/data/runtime_test_hidden_commands.yaml")
				expected := `COMMANDS:
   helptest  Help tests`

				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should not display hidden subcommands", func() {
				out := execQuiet("helptest --help", "test/data/runtime_test_hidden_commands.yaml")
				expected := `COMMANDS:
   placeholder  ...
   subcommand   Description for subcommand`

				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should display internal commands when hide is set to false", func() {
				out := execQuiet("", "test/data/runtime_test_display_internal_commands.yaml")
				expected := `COMMANDS:
   internal  Internal centry commands`

				test.AssertStringContains(g, out.Stdout, expected)
			})
		})

		g.Describe("hidden options", func() {
			g.It("should not display internal or hidden options", func() {
				out := execQuiet("", "test/data/runtime_test_hidden_options.yaml")
				expected := `OPTIONS:
   --visible value  A visible option
   --help, -h       Show help (default: false)
   --version, -v    Print the version (default: false)`

				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should not display hidden subcommand options", func() {
				out := execQuiet("helptest subcommand --help", "test/data/runtime_test_hidden_options.yaml")
				expected := `OPTIONS:
   --opt1 value, -o value  Help text for opt1 (default: "footothebar")
   --help, -h              Show help (default: false)`

				test.AssertStringContains(g, out.Stdout, expected)
			})

			g.It("should display internal options when hide is set to false", func() {
				out := execQuiet("", "test/data/runtime_test_display_internal_options.yaml")
				expected := `OPTIONS:
   --centry-config-log-level value  Overrides the log level (default: "info")
   --centry-quiet                   Disables logging (default: false)`

				test.AssertStringContains(g, out.Stdout, expected)
			})
		})
	})
}

const defaultManifestPath string = "test/data/runtime_test.yaml"

type execResult struct {
	Source   string
	ExitCode int
	Error    error
	Stdout   string
	Stderr   string
}

func execQuiet(source string, params ...string) *execResult {
	manifestPath := defaultManifestPath
	if len(params) > 0 {
		if params[0] != "" {
			manifestPath = params[0]
		}
	}
	return execCentry(source, true, manifestPath)
}

func execWithLogging(source string, params ...string) *execResult {
	manifestPath := defaultManifestPath
	if len(params) > 0 {
		if params[0] != "" {
			manifestPath = params[0]
		}
	}
	return execCentry(source, false, manifestPath)
}

func execCentry(source string, quiet bool, manifestPath string) *execResult {
	var exitCode int
	var runtimeErr error

	out := test.CaptureOutput(func() {
		if quiet {
			source = fmt.Sprintf("--centry-quiet %s", source)
		}
		if manifestPath != "" {
			source = fmt.Sprintf("--centry-file %s %s", manifestPath, source)
		}
		context := NewContext(CLI, io.Headless())
		os.Args = strings.Split(fmt.Sprintf("program %s", source), " ")
		runtime, err := NewRuntime(os.Args[1:], context)
		if err != nil {
			exitCode = 1
			runtimeErr = err
		} else {
			exitCode = runtime.Execute()
		}
	})

	if out.ExitCode > 0 {
		exitCode = out.ExitCode
	}

	return &execResult{
		Source:   source,
		ExitCode: exitCode,
		Error:    runtimeErr,
		Stdout:   out.Stdout,
		Stderr:   out.Stderr,
	}
}
