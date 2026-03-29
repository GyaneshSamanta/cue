package stacks

import "github.com/GyaneshSamanta/gyanesh-help/internal/store"

func init() { store.RegisterStack(&MobileStack{}) }

type MobileStack struct{}

func (s *MobileStack) Name() string         { return "mobile" }
func (s *MobileStack) Description() string  { return "Mobile development: React Native, Flutter, Expo, Android toolchain" }
func (s *MobileStack) EstimatedSizeMB() int { return 1800 }

func (s *MobileStack) Components() []store.Component {
	return []store.Component{
		{Name: "Node.js LTS", Version: "lts", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"nodejs", "npm"}, Darwin: []string{"node"}, Windows: []string{"OpenJS.NodeJS.LTS"}}},
		{Name: "React Native CLI", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g react-native-cli"}},
		{Name: "Expo CLI", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g expo-cli"}},
		{Name: "React DevTools", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g react-devtools"}},
		{Name: "Flutter SDK", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"flutter"},
				Darwin:  []string{"flutter"},
				Windows: []string{""},
				Script:  "git clone https://github.com/flutter/flutter.git -b stable ~/flutter && echo 'export PATH=\"$HOME/flutter/bin:$PATH\"' >> ~/.bashrc",
			}},
		{Name: "Android Studio", Version: "latest", Optional: true, OptionalPrompt: "(full Android IDE + SDK + AVD, ~3 GB)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"android-studio"},
				Darwin:  []string{"android-studio"},
				Windows: []string{"Google.AndroidStudio"},
			}},
		{Name: "Xcode CLI Tools", Version: "latest", Optional: true, OptionalPrompt: "(macOS only — iOS development)", OS: []string{"darwin"},
			InstallMethod: store.InstallMethod{Darwin: []string{""}, Script: "xcode-select --install"}},
		{Name: "Watchman", Version: "latest", Optional: true, OptionalPrompt: "(file watcher for React Native)",
			InstallMethod: store.InstallMethod{
				Linux:  []string{"watchman"},
				Darwin: []string{"watchman"},
			}},
	}
}

func (s *MobileStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Node.js", Command: "node -v", Pattern: `v\d+`},
		{Name: "npm", Command: "npm -v"},
		{Name: "Expo CLI", Command: "expo --version"},
		{Name: "Flutter", Command: "flutter --version"},
	}
}
