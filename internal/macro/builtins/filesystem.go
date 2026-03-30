package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerFilesystemMacros() {
	reg(&macro.Macro{
		Name: "find-big-files", Category: "filesystem",
		Description: "Find files larger than 100MB in current directory",
		Commands: []macro.Step{
			{OS: "linux", Command: `find . -size +100M -exec ls -lh {} \; 2>/dev/null`},
			{OS: "darwin", Command: `find . -size +100M -exec ls -lh {} \; 2>/dev/null`},
			{OS: "windows", Command: `forfiles /S /C "cmd /c if @fsize GTR 104857600 echo @path @fsize"`},
		},
		Explanation: `
✔ Listed all files larger than 100MB.
─────────────────────────────────────────────────────
These large files might be candidates for .gitignore,
Git LFS, or deletion if they're build artifacts.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "find-old-logs", Category: "filesystem",
		Description: "Find .log files older than 7 days",
		Commands: []macro.Step{
			{OS: "linux", Command: `find . -name "*.log" -mtime +7 -ls`},
			{OS: "darwin", Command: `find . -name "*.log" -mtime +7 -ls`},
			{OS: "windows", Command: `forfiles /S /M *.log /D -7 /C "cmd /c echo @path @fdate"`},
		},
		Explanation: `
✔ Listed all .log files older than 7 days.
─────────────────────────────────────────────────────
These logs are likely safe to delete. Add -delete
(Unix) to the find command to remove them, or use
'find . -name "*.log" -mtime +7 -delete'
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "nuke-node", Category: "nodejs", Dangerous: true,
		Description: "Remove node_modules and package-lock.json",
		Commands: []macro.Step{
			{OS: "linux", Command: "rm -rf node_modules package-lock.json"},
			{OS: "darwin", Command: "rm -rf node_modules package-lock.json"},
			{OS: "windows", Command: "rmdir /s /q node_modules 2>nul & del /q package-lock.json 2>nul"},
		},
		Explanation: `
✔ Done. node_modules and package-lock.json removed.
─────────────────────────────────────────────────────
Run 'npm install' to regenerate them from package.json.
This is useful when dependencies get corrupted or you
want a completely clean install.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
