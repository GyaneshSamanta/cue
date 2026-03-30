package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&ClaudeStack{}) }

type ClaudeStack struct{}

func (s *ClaudeStack) Name() string            { return "claude" }
func (s *ClaudeStack) Description() string     { return "Claude/Anthropic dev: SDK, MCP, promptfoo, CLI tools" }
func (s *ClaudeStack) EstimatedSizeMB() int    { return 150 }

func (s *ClaudeStack) Components() []store.Component {
	return []store.Component{
		{Name: "Anthropic SDK", Version: "latest",
			InstallMethod: store.InstallMethod{Script: "pip install anthropic"}},
		{Name: "MCP SDK", Version: "latest",
			InstallMethod: store.InstallMethod{Script: "npm install -g @anthropic-ai/sdk"}},
		{Name: "promptfoo", Version: "latest",
			InstallMethod: store.InstallMethod{Script: "npm install -g promptfoo"}},
		{Name: "jq", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"jq"}, Darwin: []string{"jq"}, Windows: []string{"jqlang.jq"}}},
		{Name: "HTTPie", Version: "latest",
			InstallMethod: store.InstallMethod{Script: "pip install httpie"}},
	}
}

func (s *ClaudeStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Anthropic SDK", Command: `python3 -c "import anthropic; print(anthropic.__version__)"`},
		{Name: "promptfoo", Command: "promptfoo --version"},
		{Name: "jq", Command: "jq --version"},
		{Name: "HTTPie", Command: "http --version"},
	}
}
