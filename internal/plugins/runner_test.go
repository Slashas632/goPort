package plugins

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunAll_NoPluginsInstalled_DoesNotPanic(t *testing.T) {
	withFakeHome(t)
	// Nothing installed, PluginDir() doesn't even exist yet.
	RunAll("10.0.0.1", 80, "HTTP/1.1 200 OK")
}

func TestRunAll_ExecutesInstalledLuaPlugin(t *testing.T) {
	withFakeHome(t)
	if err := os.MkdirAll(PluginDir(), 0755); err != nil {
		t.Fatalf("creating plugin dir: %v", err)
	}

	// The plugin writes its received arguments to a marker file so the
	// test can prove the Lua scan() function actually ran with the
	// expected ip/port/banner.
	marker := filepath.Join(PluginDir(), "marker.txt")
	script := `
function scan(ip, port, banner)
    local f = io.open("` + escapeLua(marker) + `", "w")
    f:write(ip .. "|" .. port .. "|" .. banner)
    f:close()
end
`
	if err := os.WriteFile(filepath.Join(PluginDir(), "test.lua"), []byte(script), 0644); err != nil {
		t.Fatalf("writing plugin: %v", err)
	}

	RunAll("10.0.0.1", 22, "SSH-2.0-OpenSSH")

	// Lua execution here is synchronous, but give the filesystem a brief
	// moment just in case, then assert the marker's exact contents.
	deadline := time.Now().Add(2 * time.Second)
	var got []byte
	var err error
	for time.Now().Before(deadline) {
		got, err = os.ReadFile(marker)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("plugin did not run (marker file missing): %v", err)
	}

	want := "10.0.0.1|22|SSH-2.0-OpenSSH"
	if string(got) != want {
		t.Errorf("marker content = %q, want %q", got, want)
	}
}

func TestRunAll_BrokenPluginDoesNotPanic(t *testing.T) {
	withFakeHome(t)
	if err := os.MkdirAll(PluginDir(), 0755); err != nil {
		t.Fatalf("creating plugin dir: %v", err)
	}
	// Invalid Lua syntax: runPlugin must swallow the DoFile error, not
	// crash the whole scan over one bad plugin.
	if err := os.WriteFile(filepath.Join(PluginDir(), "broken.lua"), []byte("this is not valid lua {{{"), 0644); err != nil {
		t.Fatalf("writing plugin: %v", err)
	}

	RunAll("10.0.0.1", 80, "banner")
}

func TestRunAll_PluginMissingScanFunctionDoesNotPanic(t *testing.T) {
	withFakeHome(t)
	if err := os.MkdirAll(PluginDir(), 0755); err != nil {
		t.Fatalf("creating plugin dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(PluginDir(), "noscan.lua"), []byte("x = 1\n"), 0644); err != nil {
		t.Fatalf("writing plugin: %v", err)
	}

	RunAll("10.0.0.1", 80, "banner")
}

// escapeLua escapes backslashes for embedding a filesystem path inside a
// Lua double-quoted string literal (relevant on Windows-style paths).
func escapeLua(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			out = append(out, '\\', '\\')
		} else {
			out = append(out, s[i])
		}
	}
	return string(out)
}
