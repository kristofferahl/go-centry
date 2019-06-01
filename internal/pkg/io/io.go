package io

import (
	"bytes"
	"io"
	"os"
)

// InputOutput holds the reader and writers used during execution
type InputOutput struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Standard creates InputOutput for use from a terminal
func Standard() InputOutput {
	return InputOutput{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Headless creates InputOutput for use from a terminal that can't accept input
func Headless() InputOutput {
	return InputOutput{
		Stdin:  nil,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Buffered creates a buffered InputOutput object
func Buffered() (io InputOutput, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	return InputOutput{
		Stdin:  nil,
		Stdout: &bufOut,
		Stderr: &bufErr,
	}, &bufOut, &bufErr
}

// BufferedCombined creates a buffered InputOutput object
func BufferedCombined() (InputOutput, *bytes.Buffer) {
	var buf bytes.Buffer
	return InputOutput{
		Stdin:  nil,
		Stdout: &buf,
		Stderr: &buf,
	}, &buf
}
