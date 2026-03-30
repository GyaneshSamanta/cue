package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&Web3Stack{}) }

type Web3Stack struct{}

func (s *Web3Stack) Name() string         { return "web3" }
func (s *Web3Stack) Description() string  { return "Web3/Blockchain: Hardhat, Foundry, ethers.js, IPFS" }
func (s *Web3Stack) EstimatedSizeMB() int { return 400 }

func (s *Web3Stack) Components() []store.Component {
	return []store.Component{
		{Name: "Node.js LTS", Version: "lts", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"nodejs", "npm"},
				Darwin:  []string{"node"},
				Windows: []string{"OpenJS.NodeJS.LTS"},
			}},
		{Name: "Hardhat", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g hardhat"}},
		{Name: "Foundry", Version: "latest", Optional: true, OptionalPrompt: "(alternative Solidity framework by Paradigm)",
			InstallMethod: store.InstallMethod{Script: "curl -L https://foundry.paradigm.xyz | bash && foundryup"}},
		{Name: "ethers.js", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g ethers"}},
		{Name: "web3.js", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g web3"}},
		{Name: "Solidity Compiler", Version: "latest", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g solc"}},
		{Name: "IPFS CLI", Version: "latest", Optional: true, OptionalPrompt: "(decentralized storage)",
			InstallMethod: store.InstallMethod{
				Darwin: []string{"ipfs"},
				Script: "curl -fsSL https://dist.ipfs.tech/kubo/v0.24.0/kubo_v0.24.0_linux-amd64.tar.gz | tar xz && cd kubo && sudo bash install.sh",
			}},
		{Name: "Ganache CLI", Version: "latest", Optional: true, OptionalPrompt: "(local Ethereum test chain)", DependsOn: []string{"Node.js LTS"},
			InstallMethod: store.InstallMethod{Script: "npm install -g ganache"}},
	}
}

func (s *Web3Stack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Node.js", Command: "node -v", Pattern: `v\d+`},
		{Name: "Hardhat", Command: "npx hardhat --version"},
		{Name: "solc", Command: "solcjs --version"},
	}
}
