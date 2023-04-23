package main

import (
	"log"
	"os"

	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/io"
)

func main() {
	args := os.Args[1:]

	// Create the context
	context := NewContext(CLI, io.Standard())

	// Create the runtime
	runtime, err := NewRuntime(args, context)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(runtime.args) == 0 && context.manifest.Config.HelpMode == config.HelpModeInteractive {
		interactive(runtime)
	}

	// Run and exit
	os.Exit(runtime.Execute())
}
