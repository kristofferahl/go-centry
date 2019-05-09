package testing

import (
	"bytes"
	"io"
	"os"
)

// Output represents config for capturing stdout and or stderr.
type Output struct {
	captureStdout bool
	captureStderr bool
}

// OutputCapture contains the result of the capture opreation.
type OutputCapture struct {
	Stdout string
	Stderr string
}

// CaptureOutput captures stdout and stderr.
func CaptureOutput(f func()) *OutputCapture {
	output := &Output{captureStdout: true, captureStderr: true}
	return output.capture(f)
}

func (output *Output) capture(f func()) *OutputCapture {
	rOut, wOut, errOut := os.Pipe()
	if errOut != nil {
		panic(errOut)
	}

	rErr, wErr, errErr := os.Pipe()
	if errErr != nil {
		panic(errErr)
	}

	if output.captureStdout {
		stdout := os.Stdout
		os.Stdout = wOut
		defer func() {
			os.Stdout = stdout
		}()
	}

	if output.captureStderr {
		stderr := os.Stderr
		os.Stderr = wErr
		defer func() {
			os.Stderr = stderr
		}()
	}

	f()

	wOut.Close()
	wErr.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	io.Copy(&stdoutBuf, rOut)
	io.Copy(&stderrBuf, rErr)

	return &OutputCapture{
		Stdout: stdoutBuf.String(),
		Stderr: stderrBuf.String(),
	}
}
