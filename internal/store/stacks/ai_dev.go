package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&AIDevStack{}) }

type AIDevStack struct{}

func (s *AIDevStack) Name() string            { return "ai-dev" }
func (s *AIDevStack) Description() string     { return "AI/ML dev: Python, CUDA, PyTorch, TensorFlow, HuggingFace, Ollama" }
func (s *AIDevStack) EstimatedSizeMB() int    { return 8000 }

func (s *AIDevStack) Components() []store.Component {
	return []store.Component{
		{Name: "Python", Version: "3.11+", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"python3", "python3-pip", "python3-venv"}, Darwin: []string{"python@3.11"}, Windows: []string{"Python.Python.3.11"}}},
		{Name: "CUDA Toolkit", Version: "12.x", Optional: true, OptionalPrompt: "(NVIDIA GPU required, ~3GB)",
			InstallMethod: store.InstallMethod{Linux: []string{"nvidia-cuda-toolkit"}, Windows: []string{"Nvidia.CUDA"}}},
		{Name: "PyTorch", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install torch torchvision torchaudio"}},
		{Name: "TensorFlow", Version: "2.x", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install tensorflow"}},
		{Name: "HuggingFace", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install transformers datasets accelerate huggingface_hub"}},
		{Name: "Ollama", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"ollama"}, Darwin: []string{"ollama"}, Windows: []string{"Ollama.Ollama"}, Script: "curl -fsSL https://ollama.ai/install.sh | sh"}},
		{Name: "FAISS", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install faiss-cpu"}},
		{Name: "MLflow", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install mlflow"}},
		{Name: "Jupyter", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install jupyterlab ipywidgets"}},
	}
}

func (s *AIDevStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Python", Command: "python3 --version", Pattern: `Python 3\.`},
		{Name: "PyTorch", Command: `python3 -c "import torch; print(torch.__version__)"`},
		{Name: "TensorFlow", Command: `python3 -c "import tensorflow; print(tensorflow.__version__)"`},
		{Name: "HuggingFace", Command: `python3 -c "import transformers; print(transformers.__version__)"`},
		{Name: "Ollama", Command: "ollama --version"},
		{Name: "MLflow", Command: "mlflow --version"},
		{Name: "GPU Available", Command: `python3 -c "import torch; print('CUDA:', torch.cuda.is_available())"`},
	}
}
