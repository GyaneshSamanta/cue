package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&GameDevStack{}) }

type GameDevStack struct{}

func (s *GameDevStack) Name() string         { return "game-dev" }
func (s *GameDevStack) Description() string  { return "Game development: Pygame, Godot, Phaser, ImageMagick, ffmpeg" }
func (s *GameDevStack) EstimatedSizeMB() int { return 500 }

func (s *GameDevStack) Components() []store.Component {
	return []store.Component{
		{Name: "Python 3", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"python3", "python3-pip"},
				Darwin:  []string{"python"},
				Windows: []string{"Python.Python.3.12"},
			}},
		{Name: "Pygame", Version: "latest", DependsOn: []string{"Python 3"},
			InstallMethod: store.InstallMethod{Script: "pip3 install pygame"}},
		{Name: "Node.js LTS", Version: "lts", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"nodejs", "npm"},
				Darwin:  []string{"node"},
				Windows: []string{"OpenJS.NodeJS.LTS"},
			}},
		{Name: "Phaser", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g phaser3-project-template"}},
		{Name: "Godot Engine", Version: "latest", Optional: true, OptionalPrompt: "(open-source game engine)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"godot3"},
				Darwin:  []string{"godot"},
				Windows: []string{"GodotEngine.GodotEngine"},
			}},
		{Name: "ImageMagick", Version: "latest",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"imagemagick"},
				Darwin:  []string{"imagemagick"},
				Windows: []string{"ImageMagick.ImageMagick"},
			}},
		{Name: "ffmpeg", Version: "latest",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"ffmpeg"},
				Darwin:  []string{"ffmpeg"},
				Windows: []string{"Gyan.FFmpeg"},
			}},
	}
}

func (s *GameDevStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Python", Command: "python3 --version", Pattern: `3\.\d+`},
		{Name: "Pygame", Command: "python3 -c \"import pygame; print(pygame.ver)\""},
		{Name: "ImageMagick", Command: "magick --version"},
		{Name: "ffmpeg", Command: "ffmpeg -version"},
	}
}
