package macro

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Macro defines a semantic shortcut command.
type Macro struct {
	Name        string
	Category    string
	Description string
	Command     string // Simple single-command (alternative to Commands)
	Commands    []Step
	Explanation string
	Flags       []Flag
	Dangerous   bool
	BuiltIn     bool
	Source      string // e.g. "builtin", "plugin:<name>", "user"
}

// Step is a single shell command, optionally OS-gated.
type Step struct {
	OS      string // "all" | "linux" | "darwin" | "windows"
	Command string
	Args    []string
}

// Flag is a named option for a macro.
type Flag struct {
	Name        string
	Description string
	Default     string
}

// Registry is the global macro map.
var Registry = map[string]*Macro{}

// Register adds a macro to the global registry.
func Register(m *Macro) {
	Registry[m.Name] = m
}

// StepsForOS returns only the steps applicable to the given OS.
func (m *Macro) StepsForOS(os string, flags map[string]string) []Step {
	var result []Step
	for _, s := range m.Commands {
		if s.OS == "all" || s.OS == os {
			result = append(result, s)
		}
	}
	return result
}

// tomlMacro is the TOML deserialization format for user macros.
type tomlMacro struct {
	Name        string   `toml:"name"`
	Command     string   `toml:"command"`
	Category    string   `toml:"category"`
	Description string   `toml:"description"`
	Explanation string   `toml:"explanation"`
	Tags        []string `toml:"tags"`
}

// LoadUserMacros reads custom macros from macros.toml and merges into Registry.
func LoadUserMacros(path string) error {
	var file struct {
		Macro []tomlMacro `toml:"macro"`
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(path, &file); err != nil {
		return err
	}
	for _, tm := range file.Macro {
		Registry[tm.Name] = &Macro{
			Name:        tm.Name,
			Category:    tm.Category,
			Description: tm.Description,
			Commands:    []Step{{OS: "all", Command: tm.Command}},
			Explanation: tm.Explanation,
			BuiltIn:     false,
		}
	}
	return nil
}
