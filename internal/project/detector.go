package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Detection represents a detected project type and its recommendations.
type Detection struct {
	Files     []string
	StackType string
	StoreName string
	Missing   []string
	EnvVars   []EnvVarStatus
}

// EnvVarStatus shows whether a required env var is set.
type EnvVarStatus struct {
	Name string
	Set  bool
}

// DetectionSignal maps file patterns to stack types.
type DetectionSignal struct {
	Files     []string // files that must exist
	DepFiles  map[string][]string // package file → dependency keywords
	StackType string
	StoreName string
}

var signals = []DetectionSignal{
	{Files: []string{"package.json", "next.config.js"}, StackType: "Next.js", StoreName: "frontend"},
	{Files: []string{"package.json", "next.config.mjs"}, StackType: "Next.js", StoreName: "frontend"},
	{Files: []string{"package.json", "next.config.ts"}, StackType: "Next.js", StoreName: "frontend"},
	{Files: []string{"package.json", "vite.config.js"}, StackType: "Vite/Frontend", StoreName: "frontend"},
	{Files: []string{"package.json", "vite.config.ts"}, StackType: "Vite/Frontend", StoreName: "frontend"},
	{Files: []string{"Cargo.toml"}, StackType: "Rust", StoreName: "rust"},
	{Files: []string{"go.mod"}, StackType: "Go", StoreName: "golang"},
	{Files: []string{"pom.xml"}, StackType: "Java/Maven", StoreName: "java"},
	{Files: []string{"build.gradle"}, StackType: "Java/Gradle", StoreName: "java"},
	{Files: []string{"build.gradle.kts"}, StackType: "Kotlin/Gradle", StoreName: "java"},
	{Files: []string{"pubspec.yaml"}, StackType: "Flutter", StoreName: "mobile"},
	{Files: []string{"hardhat.config.js"}, StackType: "Web3/Hardhat", StoreName: "web3"},
	{Files: []string{"hardhat.config.ts"}, StackType: "Web3/Hardhat", StoreName: "web3"},
	{Files: []string{"foundry.toml"}, StackType: "Web3/Foundry", StoreName: "web3"},
	{Files: []string{"docker-compose.yml"}, StackType: "Containerized", StoreName: "backend"},
	{Files: []string{"docker-compose.yaml"}, StackType: "Containerized", StoreName: "backend"},
	{Files: []string{"Dockerfile"}, StackType: "Containerized", StoreName: "backend"},
	{Files: []string{"requirements.txt"}, StackType: "Python", StoreName: "data-science"},
	{Files: []string{"pyproject.toml"}, StackType: "Python", StoreName: "data-science"},
	{Files: []string{"Pipfile"}, StackType: "Python", StoreName: "data-science"},
	{Files: []string{"composer.json"}, StackType: "PHP", StoreName: "lamp"},
	{Files: []string{"Gemfile"}, StackType: "Ruby", StoreName: ""},
}

// Detect scans the given directory for project type indicators.
func Detect(dir string) []Detection {
	var detections []Detection

	for _, sig := range signals {
		allFound := true
		var found []string
		for _, f := range sig.Files {
			path := filepath.Join(dir, f)
			if _, err := os.Stat(path); err != nil {
				allFound = false
				break
			}
			found = append(found, f)
		}
		if allFound {
			det := Detection{
				Files:     found,
				StackType: sig.StackType,
				StoreName: sig.StoreName,
			}
			detections = append(detections, det)
		}
	}

	// Check for .env.example vs .env
	envExample := filepath.Join(dir, ".env.example")
	envFile := filepath.Join(dir, ".env")
	if data, err := os.ReadFile(envExample); err == nil {
		envData, _ := os.ReadFile(envFile)
		envContent := string(envData)

		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) > 0 {
				varName := strings.TrimSpace(parts[0])
				if varName != "" {
					set := strings.Contains(envContent, varName+"=")
					if len(detections) > 0 {
						detections[0].EnvVars = append(detections[0].EnvVars, EnvVarStatus{
							Name: varName,
							Set:  set,
						})
					}
				}
			}
		}
	}

	return detections
}

// PrintDetection renders the detection results.
func PrintDetection(detections []Detection) {
	if len(detections) == 0 {
		ui.PrintInfo("No recognizable project type detected in this directory.")
		return
	}

	ui.PrintHeader("Project Detection")

	for _, det := range detections {
		fmt.Printf("  Found: %-30s → %s project\n", strings.Join(det.Files, " + "), det.StackType)
	}

	fmt.Println()

	// Primary detection
	primary := detections[0]
	fmt.Printf("  Project type: %s\n\n", primary.StackType)

	if primary.StoreName != "" {
		fmt.Printf("  Recommended store   : cue store install %s\n", primary.StoreName)
	}

	if len(primary.EnvVars) > 0 {
		fmt.Println("\n  Environment variables (.env.example vs .env):")
		for _, ev := range primary.EnvVars {
			icon := "✔"
			status := "set"
			if !ev.Set {
				icon = "✗"
				status = "not set"
			}
			fmt.Printf("    %s  %-24s (%s)\n", icon, ev.Name, status)
		}
	}
}
