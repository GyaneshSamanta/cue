package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func init() {
	registerKubernetesMacros()
}

func registerKubernetesMacros() {
	macro.Register(&macro.Macro{
		Name:        "k8s-context",
		Command:     "kubectl config get-contexts && echo '' && kubectl config current-context",
		Description: "List and show current Kubernetes context",
		Explanation: `Shows all configured Kubernetes contexts and highlights the currently active one.
Useful when working with multiple clusters (dev, staging, production).
To switch contexts: kubectl config use-context <name>`,
		Dangerous: false,
		BuiltIn:   true,
	})

	macro.Register(&macro.Macro{
		Name:        "k8s-pod-shell",
		Commands:    []macro.Step{{OS: "all", Command: "kubectl exec -it $1 -- bash || kubectl exec -it $1 -- sh"}},
		Description: "Open an interactive shell in a Kubernetes pod",
		Explanation: `Connects to a running pod and opens a bash shell (falls back to sh).
Useful for debugging containers in a cluster.
WARNING: Changes made inside the pod are ephemeral — they disappear on restart.`,
		Dangerous: false,
		BuiltIn:   true,
		Flags:     []macro.Flag{{Name: "pod", Description: "Pod name"}},
	})

	macro.Register(&macro.Macro{
		Name:        "k8s-logs",
		Commands:    []macro.Step{{OS: "all", Command: "kubectl logs -f $1 --tail=100"}},
		Description: "Stream logs from a Kubernetes pod",
		Explanation: `Follows the log output of a running pod, showing the most recent 100 lines first.
Press Ctrl+C to stop streaming.
For multi-container pods, add --container=<name>.`,
		Dangerous: false,
		BuiltIn:   true,
		Flags:     []macro.Flag{{Name: "pod", Description: "Pod name"}},
	})

	macro.Register(&macro.Macro{
		Name:        "port-forward",
		Commands:    []macro.Step{{OS: "all", Command: "kubectl port-forward $1 $2"}},
		Description: "Forward a local port to a Kubernetes pod",
		Explanation: `Creates a tunnel from your local machine to a pod's port.
Example: port-forward my-pod 8080:80 forwards local port 8080 to pod port 80.
Press Ctrl+C to stop forwarding.`,
		Dangerous: false,
		BuiltIn:   true,
		Flags: []macro.Flag{
			{Name: "pod", Description: "Pod name"},
			{Name: "port", Description: "Port mapping (local:remote)"},
		},
	})
}
