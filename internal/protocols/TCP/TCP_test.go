package TCP

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Slashas632/goPort/internal/output"
	"github.com/Slashas632/goPort/internal/ratelimit"
)

func init() {
	// Tcp() calls ratelimit.Limiter.Wait() unconditionally; without this
	// the package-level Limiter is nil and every test below panics.
	ratelimit.Init(500)
}

// listenAndBanner starts a TCP listener on 127.0.0.1 that, for each
// connection, immediately writes banner (if non-empty) and then closes.
// Returns the host and numeric port to pass into Tcp(), plus a stop func.
func listenAndBanner(t *testing.T, banner string) (host string, port int, stop func()) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting fake TCP server: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				if banner != "" {
					c.Write([]byte(banner))
				}
			}(conn)
		}
	}()

	h, p, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("splitting listener address: %v", err)
	}
	portNum, err := strconv.Atoi(p)
	if err != nil {
		t.Fatalf("parsing listener port: %v", err)
	}

	return h, portNum, func() { ln.Close() }
}

func TestTcp_RecordsBannerOnImmediateGreeting(t *testing.T) {
	host, port, stop := listenAndBanner(t, "SSH-2.0-OpenSSH_9.2p1 Debian\r\nextra ignored line")
	defer stop()

	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "results.json")
	output.JSONPath = jsonPath
	defer func() { output.JSONPath = "" }()

	Tcp(port, host)

	got := readResults(t, jsonPath)
	found := false
	for _, r := range got {
		if r.IP == host && r.Port == port && r.Banner == "SSH-2.0-OpenSSH_9.2p1 Debian" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a result with banner up to the first CRLF, got %+v", got)
	}
}

func TestTcp_FallsBackToHEADProbeWhenServerIsSilent(t *testing.T) {
	// A server that sends nothing until it receives data (e.g. a bare
	// HTTP server) should be handled via Tcp's HEAD-request fallback path.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting fake HTTP-like server: %v", err)
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
				reader := bufio.NewReader(c)
				// Wait for the HEAD request Tcp() sends after its first
				// read times out, then reply.
				line, _ := reader.ReadString('\n')
				if line != "" {
					c.Write([]byte("HTTP/1.0 200 OK\r\nServer: test\r\n\r\n"))
				}
			}(conn)
		}
	}()

	host, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("splitting listener address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parsing port: %v", err)
	}

	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "results.json")
	output.JSONPath = jsonPath
	defer func() { output.JSONPath = "" }()

	Tcp(port, host)

	got := readResults(t, jsonPath)
	found := false
	for _, r := range got {
		if r.IP == host && r.Port == port && r.Banner == "HTTP/1.0 200 OK" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected HEAD-probe fallback to record HTTP/1.0 200 OK, got %+v", got)
	}
}

func TestTcp_ClosedPortDoesNotPanicOrRecordResult(t *testing.T) {
	// Bind a listener just to grab a free port, then close it immediately
	// so the port is (almost certainly) refused when Tcp dials it.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("finding a free port: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	ln.Close()

	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "results.json")
	output.JSONPath = jsonPath
	defer func() { output.JSONPath = "" }()

	Tcp(port, "127.0.0.1")

	// output.results is a package-level slice shared by every test in this
	// binary, so assert on this specific port rather than on an empty
	// total length.
	got := readResults(t, jsonPath)
	for _, r := range got {
		if r.IP == "127.0.0.1" && r.Port == port {
			t.Errorf("did not expect a result for closed port %d, got %+v", port, r)
		}
	}
}

func readResults(t *testing.T, jsonPath string) []output.Result {
	t.Helper()

	// Give the test's own goroutine time to observe the write; Tcp() runs
	// synchronously so this is just a safety margin, not a real race.
	deadline := time.Now().Add(2 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := output.SaveJson(jsonPath); err == nil {
			break
		} else {
			lastErr = err
		}
		time.Sleep(10 * time.Millisecond)
	}
	if lastErr != nil {
		t.Fatalf("SaveJson never succeeded: %v", lastErr)
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("reading results file: %v", err)
	}
	var got []output.Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("results file is not valid JSON: %v", err)
	}
	return got
}
