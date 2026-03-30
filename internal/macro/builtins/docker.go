package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerDockerMacros() {
	reg(&macro.Macro{
		Name: "docker-nuke", Category: "docker", Dangerous: true,
		Description: "Stop all containers and prune images/volumes",
		Commands: []macro.Step{
			{OS: "all", Command: "docker stop $(docker ps -aq) 2>/dev/null; docker system prune -af --volumes"},
		},
		Explanation: `
✔ Done. Docker has been completely cleaned.
─────────────────────────────────────────────────────
All containers stopped. All unused images, networks,
and volumes removed. This freed disk space but means
you'll need to re-pull images on next docker run.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "docker-shell", Category: "docker",
		Description: "Open a bash shell in a running container",
		Commands: []macro.Step{
			{OS: "all", Command: "docker exec -it $1 bash || docker exec -it $1 sh"},
		},
		Explanation: `
✔ Opened interactive shell in the container.
─────────────────────────────────────────────────────
You're now inside the container. Type 'exit' to leave.
Falls back to 'sh' if 'bash' is not available.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "docker-compose-restart", Category: "docker",
		Description: "Restart all docker-compose services gracefully",
		Commands: []macro.Step{
			{OS: "all", Command: "docker-compose down && docker-compose up -d"},
		},
		Explanation: `
✔ Services restarted successfully.
─────────────────────────────────────────────────────
This brings down all running containers defined in your
docker-compose.yml and immediately starts them again
in detached mode.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "nuke-docker-volume", Category: "docker", Dangerous: true,
		Description: "Remove all dangling docker volumes to free space",
		Commands: []macro.Step{
			{OS: "all", Command: "docker volume rm $(docker volume ls -qf dangling=true)"},
		},
		Explanation: `
✔ Orphaned volumes destroyed.
─────────────────────────────────────────────────────
This forcibly removes any volumes that are no longer
attached to a container, freeing up disk space.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
