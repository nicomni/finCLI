// Package iostreams provides a simple abstraction for standard input, output, and error streams.
package iostreams

import "io"

// IOStreams groups standard input, output, and error streams.
// It can be used to inject custom readers and writers for testing or redirection.
type IOStreams struct {
	In  io.Reader // In is the input stream, typically os.Stdin.
	Out io.Writer // Out is the output stream, typically os.Stdout.
	Err io.Writer // Err is the error output stream, typically os.Stderr.
}
