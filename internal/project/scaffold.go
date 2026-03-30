package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/GyaneshSamanta/cue/internal/ui"
)

// ScaffoldType defines project scaffold options.
type ScaffoldType struct {
	Name        string
	Description string
	Stack       string
	InitFunc    func(dir string) error
}

var scaffolds = []ScaffoldType{
	{Name: "REST API (Node/Express)", Description: "Express.js REST API with basic structure", Stack: "mern",
		InitFunc: scaffoldNodeExpress},
	{Name: "REST API (Go/Gin)", Description: "Go Gin REST API with modules", Stack: "golang",
		InitFunc: scaffoldGoGin},
	{Name: "REST API (Python/FastAPI)", Description: "FastAPI with virtual environment", Stack: "data-science",
		InitFunc: scaffoldPythonFastAPI},
	{Name: "Static site (Next.js)", Description: "Next.js starter with TypeScript", Stack: "frontend",
		InitFunc: scaffoldNextJS},
	{Name: "Rust CLI", Description: "Rust CLI application skeleton", Stack: "rust",
		InitFunc: scaffoldRustCLI},
}

// ListScaffolds returns available scaffold types.
func ListScaffolds() []ScaffoldType {
	return scaffolds
}

// Scaffold creates a new project with the given scaffold type.
func Scaffold(name string, scaffoldIdx int) error {
	if scaffoldIdx < 0 || scaffoldIdx >= len(scaffolds) {
		return fmt.Errorf("invalid scaffold selection")
	}

	scaffold := scaffolds[scaffoldIdx]
	ui.PrintHeader(fmt.Sprintf("Creating project: %s", name))

	// Create directory
	if err := os.MkdirAll(name, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Init git
	ui.PrintStep("Initializing git repository...")
	cmd := exec.Command("git", "init", name)
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Write .gitignore
	ui.PrintStep("Writing .gitignore...")
	writeGitignore(name, scaffold.Stack)

	// Write .cue project file
	ui.PrintStep("Writing .cue config...")
	writeProjectConfig(name, scaffold.Stack)

	// Run scaffold
	ui.PrintStep(fmt.Sprintf("Scaffolding %s structure...", scaffold.Name))
	if err := scaffold.InitFunc(name); err != nil {
		ui.PrintWarning(fmt.Sprintf("Scaffold error: %v", err))
	}

	ui.PrintSuccess(fmt.Sprintf("Project '%s' ready! cd %s to start.", name, name))
	return nil
}

func writeGitignore(dir, stack string) {
	content := "node_modules/\n.env\n*.log\n.DS_Store\n"
	switch stack {
	case "golang":
		content = "*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n*.test\n*.out\nvendor/\n.env\n"
	case "rust":
		content = "/target\nCargo.lock\n*.pdb\n.env\n"
	case "data-science":
		content = "__pycache__/\n*.py[cod]\n*$py.class\n.venv/\nvenv/\n*.egg-info/\ndist/\n.env\n"
	case "mern":
		content = "node_modules/\n.env\n*.log\nbuild/\ndist/\n.DS_Store\n"
	case "frontend":
		content = "node_modules/\n.next/\nout/\n.env\n*.log\n.DS_Store\n"
	}
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(content), 0644)
}

func writeProjectConfig(dir, stack string) {
	name := filepath.Base(dir)
	content := fmt.Sprintf(`# .cue project config
tag = "%s"
stack = "%s"
`, name, stack)
	os.WriteFile(filepath.Join(dir, ".cue"), []byte(content), 0644)
}

func scaffoldNodeExpress(dir string) error {
	// Create package.json
	pkg := `{
  "name": "` + filepath.Base(dir) + `",
  "version": "1.0.0",
  "main": "src/index.js",
  "scripts": {
    "start": "node src/index.js",
    "dev": "nodemon src/index.js"
  },
  "dependencies": {
    "express": "^4.18.0",
    "cors": "^2.8.5",
    "dotenv": "^16.0.0"
  },
  "devDependencies": {
    "nodemon": "^3.0.0"
  }
}
`
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)

	// Create src directory
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	indexJS := `const express = require('express');
const cors = require('cors');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());

app.get('/', (req, res) => {
  res.json({ message: 'Hello from Express!' });
});

app.listen(PORT, () => {
  console.log('Server running on port ' + PORT);
});
`
	os.WriteFile(filepath.Join(dir, "src", "index.js"), []byte(indexJS), 0644)
	os.WriteFile(filepath.Join(dir, ".env.example"), []byte("PORT=3000\n"), 0644)
	return nil
}

func scaffoldGoGin(dir string) error {
	modName := filepath.Base(dir)
	exec.Command("go", "mod", "init", modName).Run()

	os.MkdirAll(filepath.Join(dir, "cmd"), 0755)
	os.MkdirAll(filepath.Join(dir, "internal"), 0755)

	main := `package main

import "fmt"

func main() {
	fmt.Println("Hello from Go!")
}
`
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(main), 0644)
	return nil
}

func scaffoldPythonFastAPI(dir string) error {
	os.MkdirAll(filepath.Join(dir, "app"), 0755)

	main := `from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def read_root():
    return {"message": "Hello from FastAPI!"}
`
	os.WriteFile(filepath.Join(dir, "app", "main.py"), []byte(main), 0644)

	reqs := "fastapi>=0.100.0\nuvicorn[standard]>=0.20.0\n"
	os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte(reqs), 0644)
	return nil
}

func scaffoldNextJS(dir string) error {
	ui.PrintStep("Running: npx create-next-app...")
	cmd := exec.Command("npx", "-y", "create-next-app@latest", dir, "--typescript", "--eslint", "--no-tailwind", "--src-dir", "--no-app", "--import-alias", "@/*")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func scaffoldRustCLI(dir string) error {
	cmd := exec.Command("cargo", "init", "--name", filepath.Base(dir), dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
