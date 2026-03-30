package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&MERNStack{}) }

type MERNStack struct{}

func (s *MERNStack) Name() string            { return "mern" }
func (s *MERNStack) Description() string     { return "Full MERN stack: MongoDB, Express, React, Node.js, PM2" }
func (s *MERNStack) EstimatedSizeMB() int    { return 800 }

func (s *MERNStack) Components() []store.Component {
	return []store.Component{
		{Name: "Node.js LTS", Version: "lts", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"nodejs", "npm"}, Darwin: []string{"node"}, Windows: []string{"OpenJS.NodeJS.LTS"}}},
		{Name: "MongoDB", Version: "7.x", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"mongodb-org"}, Darwin: []string{"mongodb/brew/mongodb-community"}, Windows: []string{"MongoDB.Server"}}},
		{Name: "MongoDB Compass", Version: "latest", Optional: true, OptionalPrompt: "(graphical MongoDB GUI, ~200MB)",
			InstallMethod: store.InstallMethod{Linux: []string{"mongodb-compass"}, Darwin: []string{"mongodb-compass"}, Windows: []string{"MongoDB.Compass"}}},
		{Name: "Express Generator", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g express-generator"}},
		{Name: "PM2", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g pm2"}},
		{Name: "Newman", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g newman"}},
		{Name: "Redux Toolkit", Version: "latest", Optional: true, OptionalPrompt: "(state management for React)", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g @reduxjs/toolkit"}},
	}
}

func (s *MERNStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Node.js", Command: "node -v", Pattern: `v\d+`},
		{Name: "npm", Command: "npm -v"},
		{Name: "MongoDB", Command: "mongod --version"},
		{Name: "PM2", Command: "pm2 --version"},
		{Name: "Newman", Command: "newman --version"},
	}
}
