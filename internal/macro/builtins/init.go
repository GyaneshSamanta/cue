package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func init() { RegisterAll() }

// RegisterAll registers every built-in macro.
func RegisterAll() {
	registerGitMacros()
	registerDockerMacros()
	registerFilesystemMacros()
	registerNetworkMacros()
	registerSystemMacros()
	registerPythonMacros()
	registerNodejsMacros()
	registerSSHMacros()
	registerWorkspaceMacros()
}

func reg(m *macro.Macro) { macro.Register(m) }
