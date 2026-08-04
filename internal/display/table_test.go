package display

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// everything written to it, since PrintHeader/PrintResult write directly
// to fmt's default output rather than accepting a io.Writer.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading captured stdout: %v", err)
	}
	return string(out)
}

func TestPrintHeader(t *testing.T) {
	out := captureStdout(t, PrintHeader)

	for _, want := range []string{"STATUS", "IP", "PORT", "BANNER"} {
		if !strings.Contains(out, want) {
			t.Errorf("header output %q missing column %q", out, want)
		}
	}
	if !strings.Contains(out, "─") {
		t.Error("header output missing separator line")
	}
}

func TestPrintResult(t *testing.T) {
	out := captureStdout(t, func() {
		PrintResult("10.0.0.1", 22, "SSH-2.0-OpenSSH_9.2")
	})

	for _, want := range []string{"[OPEN]", "10.0.0.1", "22", "SSH-2.0-OpenSSH_9.2"} {
		if !strings.Contains(out, want) {
			t.Errorf("result output %q missing %q", out, want)
		}
	}
	// Sanity-check the ANSI color codes are actually emitted around the
	// values so the exported color constants stay wired up correctly.
	if !strings.Contains(out, Green) || !strings.Contains(out, Cyan) || !strings.Contains(out, Yellow) || !strings.Contains(out, Reset) {
		t.Error("result output missing expected ANSI color codes")
	}
}
