package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerRustMacros() {
	macro.Register(&macro.Macro{
		Name:        "cargo-release",
		Category:    "rust",
		Command:     "cargo fmt -- --check && cargo clippy -- -D warnings && cargo test && cargo build --release",
		Description: "Full Rust release workflow: fmt → clippy → test → build --release",
		Explanation: `Runs the complete Rust release preparation workflow:
1. cargo fmt -- --check — verifies formatting (fails if not formatted)
2. cargo clippy -- -D warnings — runs linter with strict warnings
3. cargo test — runs all tests
4. cargo build --release — builds optimized release binary
If any step fails, the pipeline stops. Fix issues before re-running.`,
		Dangerous: false,
		BuiltIn:   true,
	})
}
