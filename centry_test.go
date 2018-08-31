package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/franela/goblin"
)

func TestMain(t *testing.T) {
	g := Goblin(t)

	g.Describe("scripts", func() {
		g.It("loads script in the expected order", func() {
			os.Setenv("OUTPUT_DEBUG", "true")
			g.Assert(strings.HasPrefix(ExecCentry("get").StdOut, "Loading init.sh\nLoading helpers.sh\n")).IsTrue()
			os.Unsetenv("OUTPUT_DEBUG")
		})
	})

	g.Describe("commands", func() {
		g.Describe("call without argument", func() {
			g.It("should have no arguments passed", func() {
				g.Assert(ExecCentry("get").StdOut).Equal("get ()\n")
			})
		})

		g.Describe("call with argument", func() {
			g.It("should have arguments passed", func() {
				g.Assert(ExecCentry("get foobar").StdOut).Equal("get (foobar)\n")
			})
		})
	})

	g.Describe("help", func() {
		g.Describe("call with --help", func() {
			result := ExecCentry("")

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

		g.Describe("call without command", func() {
			result := ExecCentry("")

			g.It("should display help text", func() {
				g.Assert(strings.HasPrefix(result.StdErr, "Usage: centry")).IsTrue()
			})
		})
	})
}

type ExecResult struct {
	Source string
	StdOut string
	StdErr string
	StdIn  string
	Error  error
}

func ExecCentry(source string) *ExecResult {
	return RunBash(fmt.Sprintf("./centry ./test/data/main_test.yaml %s", source))
}

func RunBash(source string) *ExecResult {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var stdin bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", source)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = &stdin
	err := cmd.Run()

	return &ExecResult{
		Source: source,
		StdOut: stdout.String(),
		StdErr: stderr.String(),
		StdIn:  stdin.String(),
		Error:  err,
	}
}
