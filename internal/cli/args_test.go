package cli

import (
	"flag"
	"os"
	"testing"
)

// --- portCheck --------------------------------------------------------

func TestPortCheck(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{"empty string means no port given", "", 0, 0, false},
		{"single valid port", "80", 80, 80, false},
		{"single port zero", "0", 0, 0, false},
		{"single port max valid", "65535", 65535, 65535, false},
		{"single port out of range", "70000", 0, 0, true},
		{"single port negative", "-5", 0, 0, true},
		{"single port not a number", "abc", 0, 0, true},
		{"valid range", "0-1024", 0, 1024, false},
		{"valid range full", "0-65535", 0, 65535, false},
		{"range equal bounds", "443-443", 443, 443, false},
		{"range start greater than end", "100-50", 0, 0, true},
		{"range end out of bounds", "60000-70000", 0, 0, true},
		{"range start invalid number", "x-100", 0, 0, true},
		{"range end invalid number", "100-x", 0, 0, true},
		{"range missing right side", "100-", 0, 0, true},
		{"range missing left side", "-100", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := portCheck(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("portCheck(%q) expected an error, got none (start=%d end=%d)", tt.input, start, end)
				}
				return
			}

			if err != nil {
				t.Fatalf("portCheck(%q) unexpected error: %v", tt.input, err)
			}
			if start != tt.wantStart || end != tt.wantEnd {
				t.Errorf("portCheck(%q) = (%d, %d), want (%d, %d)", tt.input, start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

// "-5" as a range's left side is a special case: strings.Contains(port, "-")
// is true so it goes down the range path, splits into ["", "5"], and Atoi("")
// fails independently of the bounds check. Verified separately for clarity.
func TestPortCheck_NegativeSingleGoesThroughEmptyLeftSplit(t *testing.T) {
	_, _, err := portCheck("-5")
	if err == nil {
		t.Fatal("expected an error for \"-5\", got none")
	}
}

// --- ParseArgs ----------------------------------------------------------

// resetFlags gives each subtest a clean flag.CommandLine and os.Args, since
// the flag package's default FlagSet is a shared global and flag.Parse can
// only be called once per FlagSet without panicking on redefinition.
func resetFlags(args []string) {
	flag.CommandLine = flag.NewFlagSet(args[0], flag.ContinueOnError)
	os.Args = args
}

func TestParseArgs_ValidTCPScan(t *testing.T) {
	resetFlags([]string{"goPort", "-tcp", "-ip", "10.0.0.1", "-p", "0-1024"})

	opts, err := ParseArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.TCP {
		t.Error("expected TCP=true")
	}
	if opts.UDP {
		t.Error("expected UDP=false")
	}
	if opts.IP != "10.0.0.1" {
		t.Errorf("IP = %q, want 10.0.0.1", opts.IP)
	}
	if opts.StartPort != 0 || opts.EndPort != 1024 {
		t.Errorf("port range = (%d, %d), want (0, 1024)", opts.StartPort, opts.EndPort)
	}
	if opts.Workers != 500 {
		t.Errorf("Workers = %d, want default 500", opts.Workers)
	}
}

func TestParseArgs_CustomWorkersAndJSON(t *testing.T) {
	resetFlags([]string{"goPort", "-udp", "-ip", "192.168.1.1", "-p", "53", "-w", "20", "-json", "out.json"})

	opts, err := ParseArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Workers != 20 {
		t.Errorf("Workers = %d, want 20", opts.Workers)
	}
	if opts.Json != "out.json" {
		t.Errorf("Json = %q, want out.json", opts.Json)
	}
	if opts.StartPort != 53 || opts.EndPort != 53 {
		t.Errorf("port range = (%d, %d), want (53, 53)", opts.StartPort, opts.EndPort)
	}
}

func TestParseArgs_InvalidPortReturnsError(t *testing.T) {
	resetFlags([]string{"goPort", "-tcp", "-ip", "10.0.0.1", "-p", "99999"})

	_, err := ParseArgs()
	if err == nil {
		t.Fatal("expected an error for an out-of-range port, got none")
	}
}

func TestParseArgs_InvalidPortIgnoredWhenInstalling(t *testing.T) {
	// Install/uninstall runs are allowed to skip port validation entirely,
	// since they never reach the scanning stage.
	resetFlags([]string{"goPort", "-install", "myplugin.lua", "-p", "notaport"})

	opts, err := ParseArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Install != "myplugin.lua" {
		t.Errorf("Install = %q, want myplugin.lua", opts.Install)
	}
}

func TestParseArgs_NoPortGiven(t *testing.T) {
	resetFlags([]string{"goPort", "-tcp", "-ip", "10.0.0.1"})

	opts, err := ParseArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.StartPort != 0 || opts.EndPort != 0 {
		t.Errorf("port range = (%d, %d), want (0, 0)", opts.StartPort, opts.EndPort)
	}
}
