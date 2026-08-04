package plugins

import (
	"os"
	"path/filepath"
	"testing"
)

// withFakeHome points $HOME at a fresh temp dir so PluginDir() (and
// everything built on top of it) is fully isolated per test.
func withFakeHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	return home
}

func TestPluginDir_UsesHomeGoPortPlugins(t *testing.T) {
	home := withFakeHome(t)

	want := filepath.Join(home, ".goPort", "plugins")
	if got := PluginDir(); got != want {
		t.Errorf("PluginDir() = %q, want %q", got, want)
	}
}

func TestInstall_CopiesFileIntoPluginDir(t *testing.T) {
	withFakeHome(t)

	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "myplugin.lua")
	content := []byte("function scan(ip, port, banner) end\n")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("writing source plugin: %v", err)
	}

	if err := Install(srcPath); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}

	dstPath := filepath.Join(PluginDir(), "myplugin.lua")
	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("plugin was not copied to %q: %v", dstPath, err)
	}
	if string(got) != string(content) {
		t.Errorf("copied content = %q, want %q", got, content)
	}
}

func TestInstall_MissingSourceFileReturnsError(t *testing.T) {
	withFakeHome(t)

	err := Install(filepath.Join(t.TempDir(), "does-not-exist.lua"))
	if err == nil {
		t.Fatal("expected an error for a missing source file, got none")
	}
}

func TestList_EmptyWhenPluginDirMissing(t *testing.T) {
	// Deliberately do not create PluginDir(): List must degrade
	// gracefully to an empty slice rather than propagating the error.
	withFakeHome(t)

	got, err := List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("List() = %v, want empty", got)
	}
}

func TestList_OnlyReturnsLuaFiles(t *testing.T) {
	withFakeHome(t)

	if err := os.MkdirAll(PluginDir(), 0755); err != nil {
		t.Fatalf("creating plugin dir: %v", err)
	}
	files := map[string]string{
		"a.lua":     "-- lua plugin",
		"b.lua":     "-- another lua plugin",
		"notes.txt": "not a plugin",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(PluginDir(), name), []byte(content), 0644); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
	}

	got, err := List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("List() = %v, want 2 .lua files", got)
	}
	for _, name := range got {
		if filepath.Ext(name) != ".lua" {
			t.Errorf("List() returned non-.lua file: %s", name)
		}
	}
}

func TestUninstall_RemovesExistingPlugin(t *testing.T) {
	withFakeHome(t)
	if err := os.MkdirAll(PluginDir(), 0755); err != nil {
		t.Fatalf("creating plugin dir: %v", err)
	}
	target := filepath.Join(PluginDir(), "todelete.lua")
	if err := os.WriteFile(target, []byte("-- plugin"), 0644); err != nil {
		t.Fatalf("writing plugin: %v", err)
	}

	if err := Uninstall("todelete.lua"); err != nil {
		t.Fatalf("Uninstall returned error: %v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Error("plugin file still exists after Uninstall")
	}
}

func TestUninstall_MissingPluginReturnsError(t *testing.T) {
	withFakeHome(t)

	err := Uninstall("nope.lua")
	if err == nil {
		t.Fatal("expected an error for a missing plugin, got none")
	}
}

// Regression test: Uninstall(name) must not allow name to escape
// PluginDir() via "../" path traversal.
func TestUninstall_RejectsPathTraversal(t *testing.T) {
	home := withFakeHome(t)

	// A file outside PluginDir() that a traversal attempt might target.
	victim := filepath.Join(home, "victim.txt")
	if err := os.WriteFile(victim, []byte("do not delete me"), 0644); err != nil {
		t.Fatalf("writing victim file: %v", err)
	}

	// This resolves to PluginDir()/../victim.txt if not sanitized.
	_ = Uninstall(filepath.Join("..", "victim.txt"))

	if _, err := os.Stat(victim); err != nil {
		t.Fatalf("path traversal deleted a file outside PluginDir(): %v", err)
	}
}
