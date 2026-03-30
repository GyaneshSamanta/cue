package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func init() {
	registerTerraformMacros()
}

func registerTerraformMacros() {
	macro.Register(&macro.Macro{
		Name:        "tf-plan-clean",
		Command:     "terraform init && terraform fmt -recursive && terraform validate && terraform plan",
		Description: "Run full Terraform workflow: init → fmt → validate → plan",
		Explanation: `Executes the complete Terraform pre-apply workflow in sequence:
1. terraform init — initializes providers and modules
2. terraform fmt -recursive — formats all .tf files
3. terraform validate — checks configuration syntax
4. terraform plan — shows what changes would be applied
This is the standard "check before apply" workflow.`,
		Dangerous: false,
	})
}
