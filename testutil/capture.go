//nolint:reassign // os.Stdout and os.Stderr are reassigned to capture output
package testutil

import (
	"bytes"
	"os"
)

// Capture captures the output written to os.Stdout during the execution of the provided function.
func Capture(run func()) string {
	origStdout := os.Stdout
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		out <- buf.String()
	}()

	run()

	_ = w.Close()
	os.Stdout = origStdout
	os.Stderr = origStderr

	return <-out
}
