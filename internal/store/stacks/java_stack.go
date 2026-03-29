package stacks

import "github.com/GyaneshSamanta/gyanesh-help/internal/store"

func init() { store.RegisterStack(&JavaStack{}) }

type JavaStack struct{}

func (s *JavaStack) Name() string         { return "java" }
func (s *JavaStack) Description() string  { return "Java/JVM: OpenJDK 21, SDKMAN!, Maven, Gradle, Spring Boot, Kotlin" }
func (s *JavaStack) EstimatedSizeMB() int { return 500 }

func (s *JavaStack) Components() []store.Component {
	return []store.Component{
		{Name: "SDKMAN!", Version: "latest", OS: []string{"linux", "darwin"},
			InstallMethod: store.InstallMethod{Script: `curl -s "https://get.sdkman.io" | bash`}},
		{Name: "OpenJDK 21 LTS", Version: "21", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{},
				Darwin:  []string{},
				Windows: []string{"EclipseAdoptium.Temurin.21.JDK"},
				Script:  `source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install java 21-tem`,
			}},
		{Name: "Maven", Version: "3.x", DependsOn: []string{"SDKMAN!"},
			InstallMethod: store.InstallMethod{
				Windows: []string{"Apache.Maven"},
				Script:  `source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install maven`,
			}},
		{Name: "Gradle", Version: "8.x", DependsOn: []string{"SDKMAN!"},
			InstallMethod: store.InstallMethod{
				Windows: []string{"Gradle.Gradle"},
				Script:  `source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install gradle`,
			}},
		{Name: "Spring Boot CLI", Version: "latest", Optional: true, OptionalPrompt: "(Spring framework CLI)", DependsOn: []string{"SDKMAN!"},
			InstallMethod: store.InstallMethod{
				Script: `source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install springboot`,
			}},
		{Name: "Kotlin", Version: "latest", Optional: true, OptionalPrompt: "(JVM language by JetBrains)", DependsOn: []string{"SDKMAN!"},
			InstallMethod: store.InstallMethod{
				Script: `source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install kotlin`,
			}},
		{Name: "IntelliJ IDEA Community", Version: "latest", Optional: true, OptionalPrompt: "(JetBrains Java IDE, ~800MB)",
			InstallMethod: store.InstallMethod{
				Darwin:  []string{"intellij-idea-ce"},
				Windows: []string{"JetBrains.IntelliJIDEA.Community"},
			}},
	}
}

func (s *JavaStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Java", Command: "java -version"},
		{Name: "javac", Command: "javac -version"},
		{Name: "Maven", Command: "mvn --version"},
		{Name: "Gradle", Command: "gradle --version"},
	}
}
