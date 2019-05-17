package centry

import (
	"github.com/kristofferahl/go-centry/pkg/config"
	"github.com/kristofferahl/go-centry/pkg/io"
	"github.com/kristofferahl/go-centry/pkg/log"
)

// Executor is the name of the executor
type Executor string

// CLI executor
var CLI Executor = "CLI"

// API Executor
var API Executor = "API"

// Context defines the current context
type Context struct {
	executor       Executor
	io             io.InputOutput
	log            *log.Manager
	manifest       *config.Manifest
	commandEnabled func(config.Command) bool
}

// NewContext creates a new context
func NewContext(executor Executor, io io.InputOutput) *Context {
	return &Context{
		executor: executor,
		io:       io,
	}
}
