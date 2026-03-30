package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&GolangStack{}) }

type GolangStack struct{}

func (s *GolangStack) Name() string         { return "golang" }
func (s *GolangStack) Description() string  { return "Go development: Go 1.22+, golangci-lint, air (hot reload), gvm" }
func (s *GolangStack) EstimatedSizeMB() int { return 350 }

func (s *GolangStack) Components() []store.Component {
	return []store.Component{
		{Name: "Go", Version: "1.22+", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"golang"},
				Darwin:  []string{"go"},
				Windows: []string{"GoLang.Go"},
			}},
		{Name: "gvm", Version: "latest", Optional: true, OptionalPrompt: "(Go version manager)", OS: []string{"linux", "darwin"},
			InstallMethod: store.InstallMethod{Script: "bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)"}},
		{Name: "golangci-lint", Version: "latest",
			InstallMethod: store.InstallMethod{
				Linux:   []string{},
				Darwin:  []string{"golangci-lint"},
				Windows: []string{},
				Script:  "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin",
			}},
		{Name: "air", Version: "latest", DependsOn: []string{"Go"},
			InstallMethod: store.InstallMethod{Script: "go install github.com/air-verse/air@latest"}},
		{Name: "gopls", Version: "latest", DependsOn: []string{"Go"},
			InstallMethod: store.InstallMethod{Script: "go install golang.org/x/tools/gopls@latest"}},
		{Name: "delve", Version: "latest", Optional: true, OptionalPrompt: "(Go debugger)", DependsOn: []string{"Go"},
			InstallMethod: store.InstallMethod{Script: "go install github.com/go-delve/delve/cmd/dlv@latest"}},
		{Name: "HTTPie", Version: "latest", Optional: true, OptionalPrompt: "(human-friendly HTTP client)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"httpie"},
				Darwin:  []string{"httpie"},
				Windows: []string{},
				Script:  "pip3 install httpie",
			}},
	}
}

func (s *GolangStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Go", Command: "go version", Pattern: `go1\.\d+`},
		{Name: "golangci-lint", Command: "golangci-lint --version"},
		{Name: "air", Command: "air -v"},
		{Name: "gopls", Command: "gopls version"},
	}
}
