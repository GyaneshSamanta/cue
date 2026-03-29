package stacks

import "github.com/GyaneshSamanta/gyanesh-help/internal/store"

func init() { store.RegisterStack(&DevOpsStack{}) }

type DevOpsStack struct{}

func (s *DevOpsStack) Name() string         { return "devops" }
func (s *DevOpsStack) Description() string  { return "DevOps/Platform Engineering: Terraform, Ansible, Docker, K8s, cloud CLIs" }
func (s *DevOpsStack) EstimatedSizeMB() int { return 600 }

func (s *DevOpsStack) Components() []store.Component {
	return []store.Component{
		{Name: "Docker", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"docker.io", "docker-compose-v2"},
				Darwin:  []string{"docker"},
				Windows: []string{"Docker.DockerDesktop"},
			}},
		{Name: "Terraform", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"terraform"},
				Darwin:  []string{"terraform"},
				Windows: []string{"Hashicorp.Terraform"},
			}},
		{Name: "tfenv", Version: "latest", Optional: true, OptionalPrompt: "(Terraform version manager)",
			InstallMethod: store.InstallMethod{
				Darwin: []string{"tfenv"},
				Script: "git clone https://github.com/tfutils/tfenv.git ~/.tfenv && ln -s ~/.tfenv/bin/* /usr/local/bin/",
			}},
		{Name: "Ansible", Version: "latest", OS: []string{"linux", "darwin"},
			InstallMethod: store.InstallMethod{Script: "pip3 install ansible"}},
		{Name: "kubectl", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"kubectl"},
				Darwin:  []string{"kubernetes-cli"},
				Windows: []string{"Kubernetes.kubectl"},
			}},
		{Name: "Helm", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{
				Linux:   []string{"helm"},
				Darwin:  []string{"helm"},
				Windows: []string{"Helm.Helm"},
				Script:  "curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash",
			}},
		{Name: "k9s", Version: "latest", Optional: true, OptionalPrompt: "(terminal UI for Kubernetes)",
			InstallMethod: store.InstallMethod{
				Linux:  []string{"k9s"},
				Darwin: []string{"k9s"},
			}},
		{Name: "minikube", Version: "latest", Optional: true, OptionalPrompt: "(local Kubernetes cluster)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"minikube"},
				Darwin:  []string{"minikube"},
				Windows: []string{"Kubernetes.minikube"},
			}},
		{Name: "Vault CLI", Version: "latest", Optional: true, OptionalPrompt: "(HashiCorp secrets management)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"vault"},
				Darwin:  []string{"vault"},
				Windows: []string{"Hashicorp.Vault"},
			}},
		{Name: "AWS CLI", Version: "latest", Optional: true, OptionalPrompt: "(Amazon Web Services CLI)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"awscli"},
				Darwin:  []string{"awscli"},
				Windows: []string{"Amazon.AWSCLI"},
			}},
		{Name: "GCP CLI", Version: "latest", Optional: true, OptionalPrompt: "(Google Cloud Platform CLI)",
			InstallMethod: store.InstallMethod{
				Darwin:  []string{"google-cloud-sdk"},
				Windows: []string{"Google.CloudSDK"},
				Script:  "curl https://sdk.cloud.google.com | bash",
			}},
		{Name: "Azure CLI", Version: "latest", Optional: true, OptionalPrompt: "(Microsoft Azure CLI)",
			InstallMethod: store.InstallMethod{
				Linux:   []string{"azure-cli"},
				Darwin:  []string{"azure-cli"},
				Windows: []string{"Microsoft.AzureCLI"},
			}},
		{Name: "tflint", Version: "latest", Optional: true, OptionalPrompt: "(Terraform linter)",
			InstallMethod: store.InstallMethod{
				Darwin: []string{"tflint"},
				Script: "curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash",
			}},
		{Name: "checkov", Version: "latest", Optional: true, OptionalPrompt: "(IaC security scanner)",
			InstallMethod: store.InstallMethod{Script: "pip3 install checkov"}},
	}
}

func (s *DevOpsStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Docker", Command: "docker --version"},
		{Name: "Terraform", Command: "terraform --version"},
		{Name: "kubectl", Command: "kubectl version --client --short"},
		{Name: "Helm", Command: "helm version --short"},
	}
}
