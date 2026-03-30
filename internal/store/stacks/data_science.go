package stacks

import "github.com/GyaneshSamanta/cue/internal/store"

func init() { store.RegisterStack(&DataScienceStack{}) }

type DataScienceStack struct{}

func (s *DataScienceStack) Name() string            { return "data-science" }
func (s *DataScienceStack) Description() string     { return "Complete Python/R/Jupyter data science environment" }
func (s *DataScienceStack) EstimatedSizeMB() int    { return 2500 }

func (s *DataScienceStack) Components() []store.Component {
	return []store.Component{
		{Name: "Python", Version: "3.11+", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"python3", "python3-pip", "python3-venv"}, Darwin: []string{"python@3.11"}, Windows: []string{"Python.Python.3.11"}}},
		{Name: "R", Version: "4.x", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Linux: []string{"r-base"}, Darwin: []string{"r"}, Windows: []string{"RProject.R"}}},
		{Name: "Miniconda", Version: "latest", OS: []string{"linux", "darwin", "windows"},
			InstallMethod: store.InstallMethod{Script: "conda --version || echo 'Install miniconda from https://docs.conda.io/en/latest/miniconda.html'"}},
		{Name: "JupyterLab", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install jupyterlab"}},
		{Name: "Core Libraries", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install numpy pandas matplotlib seaborn scipy scikit-learn"}},
		{Name: "DB Connectors", Version: "latest", DependsOn: []string{"Python"},
			InstallMethod: store.InstallMethod{Script: "pip install sqlalchemy psycopg2-binary pymysql"}},
		{Name: "R Packages", Version: "latest", Optional: true, OptionalPrompt: "(tidyverse, ggplot2, caret — ~500MB)", DependsOn: []string{"R"},
			InstallMethod: store.InstallMethod{Script: `Rscript -e "install.packages(c('tidyverse','ggplot2','caret'), repos='https://cloud.r-project.org')"`}},
		{Name: "VS Code + Python Extension", Version: "latest", Optional: true, OptionalPrompt: "(IDE integration)",
			InstallMethod: store.InstallMethod{Windows: []string{"Microsoft.VisualStudioCode"}, Darwin: []string{"visual-studio-code"}, Linux: []string{"code"}}},
	}
}

func (s *DataScienceStack) VerificationChecks() []store.Check {
	return []store.Check{
		{Name: "Python", Command: "python3 --version", Pattern: `Python 3\.`},
		{Name: "pip", Command: "pip3 --version"},
		{Name: "R", Command: "Rscript --version"},
		{Name: "Jupyter", Command: "jupyter --version"},
		{Name: "NumPy", Command: `python3 -c "import numpy; print(numpy.__version__)"`},
		{Name: "Pandas", Command: `python3 -c "import pandas; print(pandas.__version__)"`},
	}
}
