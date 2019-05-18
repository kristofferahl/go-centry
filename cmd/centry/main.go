package main

import (
	"os"

	"github.com/kristofferahl/go-centry/pkg/io"
)

func main() {
	args := os.Args[1:]

	// Create the context
	context := NewContext(CLI, io.Standard())

	// Create the runtime
	runtime := NewRuntime(args, context)

	// Run and exit
	os.Exit(runtime.Execute())
}
