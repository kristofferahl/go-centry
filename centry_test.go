package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/franela/goblin"
)

func TestMain(t *testing.T) {
	g := Goblin(t)

	g.Describe("scripts", func() {
		g.It("loads script in the expected order", func() {
			os.Setenv("OUTPUT_DEBUG", "true")
			g.Assert(strings.HasPrefix(execCentry("get").StdOut, "Loading init.sh\nLoading helpers.sh\n")).IsTrue()
			os.Unsetenv("OUTPUT_DEBUG")
		})
	})

	g.Describe("commands", func() {
		g.Describe("call without argument", func() {
			g.It("should have no arguments passed", func() {
				g.Assert(execCentry("get").StdOut).Equal("get ()\n")
			})
		})

		g.Describe("call with argument", func() {
			g.It("should have arguments passed", func() {
				g.Assert(execCentry("get foobar").StdOut).Equal("get (foobar)\n")
			})
		})
	})

	g.Describe("help", func() {
		g.Describe("call with --help", func() {
			result := execCentry("")

			g.It("should display available commands", func() {
				expected := `Available commands are:
    delete    Deletes stuff
    get       Gets stuff
    put       Creates stuff`

				g.Assert(strings.Contains(result.StdErr, expected)).IsTrue("\n\nEXPECTED:\n\n", expected, "\n\nTO BE FOUND IN:\n\n", result.StdErr)
			})

			g.It("should display global options", func() {
				expected := `Global options are:
    --config.log.level
        Overrides the manifest log level
    -q
        Disables logging
    --quiet
        Disables logging`

				g.Assert(strings.Contains(result.StdErr, expected)).IsTrue("\n\nEXPECTED:\n\n", expected, "\n\nTO BE FOUND IN:\n\n", result.StdErr)
			})
		})

		g.Describe("call without arguments", func() {
			result := execCentry("")

			g.It("should display help text", func() {
				g.Assert(strings.HasPrefix(result.StdErr, "Usage: centry")).IsTrue()
			})
		})
	})
}

type execResult struct {
	Source string
	StdOut string
	StdErr string
}

func execCentry(source string) *execResult {
	out := CaptureOutput(func() {
		centry(strings.Split(fmt.Sprintf("./centry ./test/data/main_test.yaml %s", source), " "))
	})

	return &execResult{
		Source: source,
		StdOut: out.Stdout,
		StdErr: out.Stderr,
	}
}
