package UDP

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Slashas632/goPort/internal/output"
	"github.com/Slashas632/goPort/internal/ratelimit"
)

func init() {
	// tryProbe/Udp call ratelimit.Limiter.Wait() unconditionally; the
	// package-level Limiter is nil until ratelimit.Init runs, which
	// would otherwise panic every test in this file.
	ratelimit.Init(500)
}

// --- isPrintable ---------------------------------------------------------

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"plain ascii", "hello world", true},
		{"printable punctuation", "SSH-2.0-OpenSSH_9.2!", true},
		{"contains null byte", "abc\x00def", false},
		{"contains control char", "abc\x01def", false},
		{"contains high byte", "abc\xffdef", false},
		{"contains newline", "abc\ndef", false},
		{"contains tab", "abc\tdef", false},
		{"boundary space (32) is printable", " ", true},
		{"boundary tilde (126) is printable", "~", true},
		{"boundary del (127) is not printable", "\x7f", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrintable(tt.input); got != tt.want {
				t.Errorf("isPrintable(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// --- probe data integrity -------------------------------------------------

func TestPortProbes_AllHaveNonEmptyPayloads(t *testing.T) {
	if len(PortProbes) == 0 {
		t.Fatal("PortProbes is empty")
	}
	for port, probes := range PortProbes {
		if len(probes) == 0 {
			t.Errorf("port %d has no probes defined", port)
		}
		for _, p := range probes {
			if len(p.Payload) == 0 {
				t.Errorf("port %d probe %q has an empty payload", port, p.Name)
			}
			if p.Name == "" {
				t.Errorf("port %d has a probe with an empty Name", port)
			}
		}
	}
}

func TestGenericProbes_NonEmpty(t *testing.T) {
	if len(GenericProbes) == 0 {
		t.Fatal("GenericProbes is empty")
	}
	for _, p := range GenericProbes {
		if len(p.Payload) == 0 {
			t.Errorf("generic probe %q has an empty payload", p.Name)
		}
	}
}

func TestKnownServices_MatchDocumentedPorts(t *testing.T) {
	// Every port advertised in the README's UDP probe table should
	// resolve to a friendly service name for non-printable banners.
	documented := []int{53, 111, 123, 137, 161, 500, 514, 1900, 5353, 11211, 51820}
	for _, port := range documented {
		if _, ok := knownServices[port]; !ok {
			t.Errorf("port %d is documented but missing from knownServices", port)
		}
	}
}

// --- tryProbe / Udp against a local fake UDP service ----------------------

// fakeUDPServer starts a UDP listener on 127.0.0.1 that replies to every
// datagram it receives with the given response, and returns its address
// plus a stop function.
func fakeUDPServer(t *testing.T, response []byte) (addr string, stop func()) {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("starting fake UDP server: %v", err)
	}

	done := make(chan struct{})
	go func() {
		buf := make([]byte, 4096)
		for {
			n, raddr, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}
			_ = n
			select {
			case <-done:
				return
			default:
			}
			if response != nil {
				conn.WriteTo(response, raddr)
			}
		}
	}()

	return conn.LocalAddr().String(), func() {
		close(done)
		conn.Close()
	}
}

func TestTryProbe_ReturnsBannerOnResponse(t *testing.T) {
	addr, stop := fakeUDPServer(t, []byte("PONG-1.0"))
	defer stop()

	banner, ok := tryProbe(addr, GenericProbes[0])
	if !ok {
		t.Fatal("tryProbe returned ok=false, want true")
	}
	if banner != "PONG-1.0" {
		t.Errorf("banner = %q, want PONG-1.0", banner)
	}
}

func TestTryProbe_NoResponseTimesOut(t *testing.T) {
	// A server that never replies should make tryProbe give up (after its
	// retry) and report ok=false rather than hang or panic.
	addr, stop := fakeUDPServer(t, nil)
	defer stop()

	_, ok := tryProbe(addr, GenericProbes[0])
	if ok {
		t.Error("tryProbe returned ok=true for a server that never responds")
	}
}

func TestUdp_UnreachableHostDoesNotPanic(t *testing.T) {
	output.JSONPath = ""
	// Port 1 on localhost should have nothing listening; Udp must return
	// quietly rather than panicking or blocking indefinitely.
	Udp(1, "127.0.0.1")
}

func TestUdp_RecordsResultViaOutputPackage(t *testing.T) {
	addr, stop := fakeUDPServer(t, []byte("PONG-1.0"))
	defer stop()

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("splitting fake server address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parsing fake server port: %v", err)
	}

	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "results.json")
	output.JSONPath = jsonPath
	defer func() { output.JSONPath = "" }()

	Udp(port, host)

	if err := output.SaveJson(jsonPath); err != nil {
		t.Fatalf("SaveJson: %v", err)
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("reading results file: %v", err)
	}
	var got []output.Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("results file is not valid JSON: %v", err)
	}

	found := false
	for _, r := range got {
		if r.IP == host && r.Port == port && r.Banner == "PONG-1.0" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a result for %s:%d with banner PONG-1.0, got %+v", host, port, got)
	}
}
