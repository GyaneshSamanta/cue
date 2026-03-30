package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&BackendStack{}) }

type BackendStack struct{}

func (s *BackendStack) Name() string            { return "backend" }
func (s *BackendStack) Description() string     { return "Backend dev: Docker, DB clients, HTTPie, Make" }
func (s *BackendStack) EstimatedSizeMB() int    { return 1200 }

func (s *BackendStack) Components() []store.Component {
	return []store.Component{
		{Name: "Docker", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"docker.io", "docker-compose"}, Darwin: []string{"docker"}, Windows: []string{"Docker.DockerDesktop"}}},
		{Name: "PostgreSQL Client", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"postgresql-client"}, Darwin: []string{"libpq"}, Windows: []string{"PostgreSQL.PostgreSQL"}}},
		{Name: "MySQL Client", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"mysql-client"}, Darwin: []string{"mysql-client"}, Windows: []string{"Oracle.MySQL"}}},
		{Name: "Redis CLI", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"redis-tools"}, Darwin: []string{"redis"}, Windows: []string{"Redis.Redis"}}},
		{Name: "HTTPie", Version: "latest",
			InstallMethod: store.InstallMethod{Script: "pip install httpie"}},
		{Name: "Make", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"make"}, Darwin: []string{"make"}, Windows: []string{"GnuWin32.Make"}}},
		{Name: "DBeaver", Version: "latest", Optional: true, OptionalPrompt: "(universal database GUI)",
			InstallMethod: store.InstallMethod{Linux: []string{"dbeaver-ce"}, Darwin: []string{"dbeaver-community"}, Windows: []string{"dbeaver.dbeaver"}}},
	}
}

func (s *BackendStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Docker", Command: "docker --version"},
		{Name: "Docker Compose", Command: "docker compose version"},
		{Name: "psql", Command: "psql --version"},
		{Name: "mysql", Command: "mysql --version"},
		{Name: "redis-cli", Command: "redis-cli --version"},
		{Name: "HTTPie", Command: "http --version"},
		{Name: "Make", Command: "make --version"},
	}
}
