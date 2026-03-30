package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerSystemMacros() {
	reg(&macro.Macro{
		Name: "env-check", Category: "system",
		Description: "Print all PATH entries, one per line",
		Commands: []macro.Step{
			{OS: "linux", Command: `echo $PATH | tr ':' '\n'`},
			{OS: "darwin", Command: `echo $PATH | tr ':' '\n'`},
			{OS: "windows", Command: `echo %PATH:;=& echo %`},
		},
		Explanation: `
✔ All PATH entries displayed.
─────────────────────────────────────────────────────
Each line is a directory where your OS looks for
executable programs. Tools must be in one of these
directories to be found by name.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "disk-check", Category: "system",
		Description: "Human-readable disk usage summary",
		Commands: []macro.Step{
			{OS: "linux", Command: "df -h --total"},
			{OS: "darwin", Command: "df -h"},
			{OS: "windows", Command: "wmic logicaldisk get size,freespace,caption"},
		},
		Explanation: `
✔ Disk usage displayed.
─────────────────────────────────────────────────────
Size = total capacity, Used = space consumed,
Avail = free space, Use% = how full the disk is.
Watch out for partitions above 90%.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "process-find", Category: "system",
		Description: "Find all running processes matching a name",
		Commands: []macro.Step{
			{OS: "linux", Command: `ps aux | grep -i "$1" | grep -v grep`},
			{OS: "darwin", Command: `ps aux | grep -i "$1" | grep -v grep`},
			{OS: "windows", Command: `tasklist /FI "IMAGENAME eq $1*"`},
		},
		Explanation: `
✔ Matching processes listed.
─────────────────────────────────────────────────────
The PID column shows the process ID. Use
'kill <PID>' (Unix) or 'taskkill /PID <PID>' (Windows)
to terminate a specific process.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "hosts-edit", Category: "system",
		Description: "Open /etc/hosts in default editor with elevation",
		Commands: []macro.Step{
			{OS: "linux", Command: "sudo ${EDITOR:-nano} /etc/hosts"},
			{OS: "darwin", Command: "sudo ${EDITOR:-nano} /etc/hosts"},
			{OS: "windows", Command: `notepad C:\Windows\System32\drivers\etc\hosts`},
		},
		Explanation: `
✔ Hosts file opened for editing.
─────────────────────────────────────────────────────
The hosts file maps hostnames to IP addresses locally.
Add entries like: 127.0.0.1  myapp.local
Save and close. Changes take effect immediately.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "path-add", Category: "system",
		Description: "Persist a new directory to PATH",
		Commands: []macro.Step{
			{OS: "linux", Command: `echo 'export PATH="$1:$PATH"' >> ~/.bashrc && source ~/.bashrc`},
			{OS: "darwin", Command: `echo 'export PATH="$1:$PATH"' >> ~/.zshrc && source ~/.zshrc`},
			{OS: "windows", Command: `setx PATH "%PATH%;$1"`},
		},
		Explanation: `
✔ Done. Directory added to PATH.
─────────────────────────────────────────────────────
The directory will be available in new terminal sessions.
On Windows, you may need to restart your terminal.
On Unix, it's added to your shell profile.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
