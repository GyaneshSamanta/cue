package plugin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/macro"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Plugin represents a TOML plugin file.
type Plugin struct {
	Meta   PluginMeta    `toml:"meta"`
	Macros []PluginMacro `toml:"macro"`
}

// PluginMeta holds plugin metadata.
type PluginMeta struct {
	Name              string `toml:"name"`
	Version           string `toml:"version"`
	Description       string `toml:"description"`
	Author            string `toml:"author"`
	MinGyaneshVersion string `toml:"min_cue_version"`
}

// PluginMacro is a macro defined in a plugin.
type PluginMacro struct {
	Name        string `toml:"name"`
	Command     string `toml:"command"`
	Dangerous   bool   `toml:"dangerous"`
	Explanation string `toml:"explanation"`
}

// PluginsDir returns path to installed plugins.
func PluginsDir() string {
	return filepath.Join(config.ConfigDir(), "plugins")
}

// Install installs a plugin from a URL or curated registry name.
func Install(source string) error {
	os.MkdirAll(PluginsDir(), 0755)

	var data []byte
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		// Download from URL
		ui.PrintStep(fmt.Sprintf("Downloading plugin from %s...", source))
		resp, err := http.Get(source)
		if err != nil {
			return fmt.Errorf("failed to download plugin: %w", err)
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read plugin: %w", err)
		}
	} else {
		// Try curated registry URL
		registryURL := fmt.Sprintf("https://raw.githubusercontent.com/GyaneshSamanta/cue-plugins/main/plugins/%s.toml", source)
		resp, err := http.Get(registryURL)
		if err != nil || resp.StatusCode != 200 {
			return fmt.Errorf("plugin '%s' not found in registry. Try a full URL.", source)
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read plugin: %w", err)
		}
	}

	// Parse the plugin TOML
	var plugin Plugin
	if err = toml.Unmarshal(data, &plugin); err != nil {
		return fmt.Errorf("invalid plugin TOML: %w", err)
	}

	// Save plugin file
	pluginPath := filepath.Join(PluginsDir(), plugin.Meta.Name+".toml")
	if err = os.WriteFile(pluginPath, data, 0644); err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Plugin '%s' v%s installed (%d macros)",
		plugin.Meta.Name, plugin.Meta.Version, len(plugin.Macros)))
	return nil
}

// Remove removes an installed plugin.
func Remove(name string) error {
	pluginPath := filepath.Join(PluginsDir(), name+".toml")
	if _, err := os.Stat(pluginPath); err != nil {
		return fmt.Errorf("plugin '%s' not installed", name)
	}
	os.Remove(pluginPath)
	ui.PrintSuccess(fmt.Sprintf("Plugin '%s' removed", name))
	return nil
}

// List lists all installed plugins.
func List() error {
	dir := PluginsDir()
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		ui.PrintInfo("No plugins installed. Use 'cue plugin install <name>' to add one.")
		return nil
	}

	ui.PrintHeader("Installed Plugins")
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}
		var plugin Plugin
		if _, err := toml.DecodeFile(filepath.Join(dir, entry.Name()), &plugin); err != nil {
			continue
		}
		fmt.Printf("  %-20s v%-8s %d macros  %s\n",
			plugin.Meta.Name, plugin.Meta.Version, len(plugin.Macros), plugin.Meta.Description)
	}
	return nil
}

// LoadAll loads all plugin macros into the macro registry.
func LoadAll() {
	dir := PluginsDir()
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}
		var plugin Plugin
		if _, err := toml.DecodeFile(filepath.Join(dir, entry.Name()), &plugin); err != nil {
			continue
		}
		for _, pm := range plugin.Macros {
			macro.Registry[pm.Name] = &macro.Macro{
				Name:        pm.Name,
				Command:     pm.Command,
				Description: pm.Explanation,
				Dangerous:   pm.Dangerous,
				Source:      "plugin:" + plugin.Meta.Name,
			}
		}
	}
}

// Create scaffolds a new plugin template.
func Create(name string) error {
	content := fmt.Sprintf(`[meta]
name = "%s"
version = "1.0.0"
description = "Description of your plugin"
author = "your-name"
min_cue_version = "2.0.0"

[[macro]]
name = "%s-example"
command = "echo 'Hello from %s plugin'"
dangerous = false
explanation = """
Example macro from the %s plugin.
Replace this with your actual macro.
"""
`, name, name, name, name)

	filename := name + ".toml"
	os.WriteFile(filename, []byte(content), 0644)
	ui.PrintSuccess(fmt.Sprintf("Plugin template created: %s", filename))
	ui.PrintInfo("Edit the file, then run 'cue plugin install ./" + filename + "' to test.")
	return nil
}
