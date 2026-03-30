package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func init() {
	registerCloudMacros()
}

func registerCloudMacros() {
	macro.Register(&macro.Macro{
		Name:        "gcloud-project-switch",
		Command:     "gcloud projects list && echo '' && read -p 'Project ID: ' proj && gcloud config set project $proj",
		Description: "Interactive GCP project selector",
		Explanation: `Lists all GCP projects you have access to, then prompts you to enter a Project ID.
Sets the selected project as the active project for all subsequent gcloud commands.
Equivalent to: gcloud config set project <PROJECT_ID>`,
		Dangerous: false,
	})
}
