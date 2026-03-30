package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&RustStack{}) }

type RustStack struct{}

func (s *RustStack) Name() string         { return "rust" }
func (s *RustStack) Description() string  { return "Rust development: rustup, cargo, clippy, rust-analyzer, cross-compilation" }
func (s *RustStack) EstimatedSizeMB() int { return 300 }

func (s *RustStack) Components() []store.Component {
	return []store.Component{
		{Name: "rustup", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{},
				Darwin:  []string{},
				Windows: []string{"Rustlang.Rustup"},
				Script:  "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y",
			}},
		{Name: "Rust stable", Version: "stable", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "rustup default stable"}},
		{Name: "rustfmt", Version: "bundled", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "rustup component add rustfmt"}},
		{Name: "clippy", Version: "bundled", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "rustup component add clippy"}},
		{Name: "rust-analyzer", Version: "bundled", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "rustup component add rust-analyzer"}},
		{Name: "cross", Version: "latest", Optional: true, OptionalPrompt: "(cross-compilation tool)", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "cargo install cross --git https://github.com/cross-rs/cross"}},
		{Name: "cargo-audit", Version: "latest", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "cargo install cargo-audit"}},
		{Name: "cargo-watch", Version: "latest", DependsOn: []string{"rustup"},
			InstallMethod: store.InstallMethod{Script: "cargo install cargo-watch"}},
	}
}

func (s *RustStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "rustup", Command: "rustup --version"},
		{Name: "rustc", Command: "rustc --version", Pattern: `\d+\.\d+`},
		{Name: "cargo", Command: "cargo --version"},
		{Name: "clippy", Command: "cargo clippy --version"},
		{Name: "rustfmt", Command: "rustfmt --version"},
	}
}
