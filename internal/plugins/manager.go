package plugins

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func PluginDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".goPort", "plugins")
}

func Install(src string) error {
	os.MkdirAll(PluginDir(), 0755)

	filename := filepath.Base(src)
	dst := filepath.Join(PluginDir(), filename)

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening plugin file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("installing plugin file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copying plugin: %w", err)
	}

	fmt.Println("Plugin installed:", filename)
	return nil
}

func List() ([]string, error) {
	entries, err := os.ReadDir(PluginDir())
	if err != nil {
		return nil, nil
	}

	var plugins []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".lua" {
			plugins = append(plugins, e.Name())
		}
	}
	return plugins, nil
}

func Uninstall(name string) error {
	path := filepath.Join(PluginDir(), name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s not found", name)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("uninstalling plugin %s: %w", name, err)
	}

	fmt.Println("Plugin uninstalled:", name)
	return nil
}
