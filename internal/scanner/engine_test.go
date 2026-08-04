package scanner

import (
	"bytes"
	"io"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Slashas632/goPort/internal/cli"
)

// captureStdout redirects os.Stdout for the duration of fn.
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

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name    string
		opts    cli.Options
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "install takes priority and always returns false",
			opts:    cli.Options{Install: "/nonexistent/plugin.lua"},
			wantOK:  false,
			wantMsg: "",
		},
		{
			name:    "uninstall takes priority and always returns false",
			opts:    cli.Options{Uninstall: "nope.lua"},
			wantOK:  false,
			wantMsg: "",
		},
		{
			name:    "neither tcp nor udp",
			opts:    cli.Options{IP: "10.0.0.1", StartPort: 1, EndPort: 1, Workers: 500},
			wantOK:  false,
			wantMsg: "specify -tcp and/or -udp",
		},
		{
			name:    "missing ip",
			opts:    cli.Options{TCP: true, StartPort: 1, EndPort: 1, Workers: 500},
			wantOK:  false,
			wantMsg: "specify -ip",
		},
		{
			name:    "missing port range",
			opts:    cli.Options{TCP: true, IP: "10.0.0.1", Workers: 500},
			wantOK:  false,
			wantMsg: "specify -p",
		},
		{
			name:    "zero workers is rejected",
			opts:    cli.Options{TCP: true, IP: "10.0.0.1", StartPort: 1, EndPort: 1, Workers: 0},
			wantOK:  false,
			wantMsg: "-w must be greater than 0",
		},
		{
			name:    "negative workers is rejected",
			opts:    cli.Options{TCP: true, IP: "10.0.0.1", StartPort: 1, EndPort: 1, Workers: -5},
			wantOK:  false,
			wantMsg: "-w must be greater than 0",
		},
		{
			name:   "fully valid tcp options",
			opts:   cli.Options{TCP: true, IP: "10.0.0.1", StartPort: 1, EndPort: 1024, Workers: 500},
			wantOK: true,
		},
		{
			name:   "fully valid udp options",
			opts:   cli.Options{UDP: true, IP: "10.0.0.1", StartPort: 53, EndPort: 53, Workers: 500},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			out := captureStdout(t, func() {
				got = ErrorHandler(tt.opts)
			})
			if got != tt.wantOK {
				t.Errorf("ErrorHandler() = %v, want %v (stdout: %q)", got, tt.wantOK, out)
			}
			if tt.wantMsg != "" && !bytesContains(out, tt.wantMsg) {
				t.Errorf("stdout = %q, want it to contain %q", out, tt.wantMsg)
			}
		})
	}
}

func bytesContains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}

// TestRun_EndToEndTCPScan is an integration test: it spins up a local TCP
// listener with a known banner, points Run() at just that one port, and
// checks the process completes and writes the expected JSON results
// instead of hanging (e.g. due to the Workers<=0 goroutine-leak bug).
func TestRun_EndToEndTCPScan(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting fake server: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				c.Write([]byte("BANNER-OK\r\n"))
			}(conn)
		}
	}()

	host, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("splitting listener addr: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parsing port: %v", err)
	}

	jsonPath := t.TempDir() + "/results.json"

	opts := cli.Options{
		TCP:       true,
		IP:        host,
		StartPort: port,
		EndPort:   port,
		Workers:   4,
		Json:      jsonPath,
	}

	done := make(chan struct{})
	go func() {
		captureStdout(t, func() { Run(opts) })
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Run() did not finish in time (possible deadlock)")
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("reading results file: %v", err)
	}
	if !bytesContains(string(data), "BANNER-OK") {
		t.Errorf("results file %q does not contain the expected banner", data)
	}
}

// TestRun_EndToEndUDPScan mirrors the TCP test above but exercises the UDP
// branch of Run() (UDPworkers was previously untested).
func TestRun_EndToEndUDPScan(t *testing.T) {
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting fake UDP server: %v", err)
	}
	defer conn.Close()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, raddr, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}
			_ = n
			conn.WriteTo([]byte("UDP-BANNER-OK"), raddr)
		}
	}()

	host, portStr, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		t.Fatalf("splitting listener addr: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parsing port: %v", err)
	}

	jsonPath := t.TempDir() + "/results.json"

	opts := cli.Options{
		UDP:       true,
		IP:        host,
		StartPort: port,
		EndPort:   port,
		Workers:   4,
		Json:      jsonPath,
	}

	done := make(chan struct{})
	go func() {
		captureStdout(t, func() { Run(opts) })
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Run() did not finish in time (possible deadlock)")
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("reading results file: %v", err)
	}
	if !bytesContains(string(data), "UDP-BANNER-OK") {
		t.Errorf("results file %q does not contain the expected banner", data)
	}
}
