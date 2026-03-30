package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&FrontendStack{}) }

type FrontendStack struct{}

func (s *FrontendStack) Name() string            { return "frontend" }
func (s *FrontendStack) Description() string     { return "Modern frontend dev: Node.js, npm, yarn, pnpm, Vite, ESLint, TypeScript" }
func (s *FrontendStack) EstimatedSizeMB() int    { return 400 }

func (s *FrontendStack) Components() []store.Component {
	return []store.Component{
		{Name: "Node.js LTS", Version: "lts", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"nodejs", "npm"}, Darwin: []string{"node"}, Windows: []string{"OpenJS.NodeJS.LTS"}}},
		{Name: "Yarn", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "corepack enable && corepack prepare yarn@stable --activate"}},
		{Name: "pnpm", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g pnpm"}},
		{Name: "Vite", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g vite"}},
		{Name: "ESLint", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g eslint"}},
		{Name: "Prettier", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g prettier"}},
		{Name: "TypeScript", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g typescript"}},
		{Name: "serve", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g serve"}},
		{Name: "Playwright", Version: "latest", Optional: true, OptionalPrompt: "(browser testing, ~100MB)", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g playwright && npx playwright install"}},
	}
}

func (s *FrontendStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Node.js", Command: "node -v", Pattern: `v\d+`},
		{Name: "npm", Command: "npm -v"},
		{Name: "yarn", Command: "yarn -v"},
		{Name: "pnpm", Command: "pnpm -v"},
		{Name: "Vite", Command: "vite --version"},
		{Name: "TypeScript", Command: "tsc -v"},
	}
}
