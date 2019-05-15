package main

import (
	"os"

	"github.com/kristofferahl/go-centry/pkg/centry"
	"github.com/kristofferahl/go-centry/pkg/io"
)

func main() {
	args := os.Args[1:]

	// Create the context
	context := centry.NewContext(io.InputOutput{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	// Create the runtime
	runtime := centry.Create(args, context)

	// Run and exit
	os.Exit(runtime.Execute())
}
