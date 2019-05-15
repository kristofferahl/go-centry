package io

import (
	"io"
)

// InputOutput holds the reader and writers used during execution
type InputOutput struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
