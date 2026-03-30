package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerGitMacros() {
	reg(&macro.Macro{
		Name: "git-undo", Category: "git",
		Description: "Safely undo the last commit, keeping changes staged",
		Flags:       []macro.Flag{{Name: "hard", Description: "Discard changes too (DESTRUCTIVE)", Default: "false"}},
		Commands:    []macro.Step{{OS: "all", Command: "git reset --soft HEAD~1"}},
		Explanation: `
✔ Done. Here's what happened:
─────────────────────────────────────────────────────
Your last commit was "undone," but your file changes
are SAFE and still staged. The commit message is gone,
but your work is not. You can re-commit when ready.
This rewrites local history only.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-clean", Category: "git", Dangerous: true,
		Description: "Remove all untracked files and directories",
		Commands:    []macro.Step{{OS: "all", Command: "git clean -fd"}},
		Explanation: `
✔ Done. All untracked files and directories were removed.
─────────────────────────────────────────────────────
Only files NOT tracked by git were deleted. Your committed
and staged files are untouched. This cannot be undone.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-save", Category: "git",
		Description: "Stash current changes with a message",
		Commands:    []macro.Step{{OS: "all", Command: `git stash push -m "cue save"`}},
		Explanation: `
✔ Done. Your working directory changes were stashed.
─────────────────────────────────────────────────────
Your changes are safely stored. Retrieve them with:
  cue git-unsave
or: git stash pop
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-unsave", Category: "git",
		Description: "Pop the most recent stash",
		Commands:    []macro.Step{{OS: "all", Command: "git stash pop"}},
		Explanation: `
✔ Done. Your most recent stash was applied and removed.
─────────────────────────────────────────────────────
The stashed changes are now back in your working directory.
If there were conflicts, resolve them manually.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-whoops", Category: "git",
		Description: "Amend the last commit (add forgotten files)",
		Commands:    []macro.Step{{OS: "all", Command: "git add -A && git commit --amend --no-edit"}},
		Explanation: `
✔ Done. All current changes were added to your last commit.
─────────────────────────────────────────────────────
The commit message stays the same. If you already pushed
this commit, you'll need: cue git-oops-push
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-oops-push", Category: "git", Dangerous: true,
		Description: "Force push with lease (safe force push)",
		Commands:    []macro.Step{{OS: "all", Command: "git push --force-with-lease"}},
		Explanation: `
✔ Done. Force push completed with lease protection.
─────────────────────────────────────────────────────
--force-with-lease ensures you don't overwrite someone
else's work. If the remote has commits you don't have,
the push is rejected. This is safer than --force.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-log-pretty", Category: "git",
		Description: "Show a pretty one-line commit graph",
		Commands:    []macro.Step{{OS: "all", Command: "git log --oneline --graph --decorate -20"}},
		Explanation: `
✔ Showing the last 20 commits in a compact graph view.
─────────────────────────────────────────────────────
Each line is one commit. Branch merges are visualized
with ASCII art. Decorations show branch/tag labels.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-branch-clean", Category: "git", Dangerous: true,
		Description: "Delete all local branches merged into main/master",
		Commands: []macro.Step{
			{OS: "linux", Command: `git branch --merged | grep -v '\*\|main\|master\|develop' | xargs -r git branch -d`},
			{OS: "darwin", Command: `git branch --merged | grep -v '\*\|main\|master\|develop' | xargs git branch -d`},
			{OS: "windows", Command: `git branch --merged | findstr /V "* main master develop" > %TEMP%\gh_branches.txt & FOR /F "tokens=*" %b IN (%TEMP%\gh_branches.txt) DO git branch -d %b`},
		},
		Explanation: `
✔ Done. All local merged branches were deleted.
─────────────────────────────────────────────────────
Only branches already merged into main/master/develop
were removed. Remote branches are untouched.
To clean remote: git remote prune origin
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-diff-staged", Category: "git",
		Description: "Show diff of staged changes",
		Commands:    []macro.Step{{OS: "all", Command: "git diff --cached"}},
		Explanation: `
✔ Showing all changes currently staged for commit.
─────────────────────────────────────────────────────
Green (+) lines are additions, red (-) are deletions.
These are the changes that will be included in your
next git commit.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-pr", Category: "git",
		Description: "Create a PR in the browser using GitHub CLI",
		Commands:    []macro.Step{{OS: "all", Command: "gh pr create --web"}},
		Explanation: `
✔ GitHub CLI launched the Pull Request page in your browser.
─────────────────────────────────────────────────────
This uses the official 'gh' command line tool.
Ensure you have run 'gh auth login' previously.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-sync", Category: "git",
		Description: "Fetch main, rebase current branch, and push",
		Commands:    []macro.Step{{OS: "all", Command: "git pull --rebase origin main && git push origin HEAD"}},
		Explanation: `
✔ Branch successfully synced with remote main.
─────────────────────────────────────────────────────
This operation pulled the most recent changes from main,
rebased your local commits on top of them, and pushed
your updated branch back to the remote.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "git-contributors", Category: "git",
		Description: "Print a leaderboard of contributors",
		Commands:    []macro.Step{{OS: "all", Command: "git shortlog -sn --no-merges"}},
		Explanation: `
✔ Leaderboard created.
─────────────────────────────────────────────────────
This prints a summary of commits per author, sorted
by commit count. Merges are ignored.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
