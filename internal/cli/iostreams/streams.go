// Package iostreams owns the process streams used by the CLI presentation layer.
package iostreams

import (
	"io"
	"os"
)

// IOStreams makes terminal input and output explicit and replaceable in tests.
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// System returns the streams attached to the current process.
func System() IOStreams {
	return IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
}
