---
description: An optimized Git workflow to safely stage, commit, and proactively rebase against develop to avoid conflicts early.
---
# /git-commit - Safe Git Commit and Sync Workflow

This workflow guides the agent to safely review, stage, commit your code, **and immediately pull the latest `develop` to rebase**, ensuring that your branch is always up-to-date right at the point of commit.

## When to Use
- `/git-commit`
- Keywords: "commit my changes", "save my work"

## Workflow Steps

1. **Check Current Status**
   Review modified files and untracked changes.
   ```bash
   git status
   ```

2. **Stage Changes**
   Add specific files or all changes as requested.
   ```bash
   git add .
   ```

3. **Commit Changes**
   Commit them with a meaningful, conventional commit message (e.g., `feat:`, `fix:`, `chore:`).
   *Note: If the message isn't obvious from the changes, quickly ask the user or propose one based on diff. **DO NOT add 'Co-authored-by' or any agent information to the commit message.***
   ```bash
   git commit -m "feat: your descriptive message"
   ```

4. **Fetch Latest Changes**
   // turbo
   Fetch the latest updates from the remote repository to ensure you have the most recent version of `develop`.
   ```bash
   git fetch origin
   ```

5. **Rebase Against Develop**
   // turbo
   Immediately rebase the current feature branch on top of `origin/develop`. This rewrites your newly created commit (and any other local commits) so they are applied after the latest changes in `develop`, minimizing future conflicts.
   ```bash
   git rebase origin/develop
   ```

6. **Auto-Resolve Conflicts (If Step 5 Fails)**
   If the `git rebase` fails due to conflicts, do NOT abort unless strictly necessary.
   - Run `git status` to identify the conflicted files.
   - Use the `view_file` tool to read the conflicted files (look for `<<<<<<<`, `=======`, `>>>>>>>` markers).
   - Analyze the conflicting changes logically. **CRITICAL:** Combine the remote changes (from `develop`) and local changes (from your branch) intelligently so that no logic or code is lost.
   - Use the `replace_file_content` or `multi_replace_file_content` tool to remove the conflict markers and write the resolved code.
   - Stage the resolved files: `git add <resolved_file>`
   - Continue the rebase:
     ```bash
     git rebase --continue
     ```
   *(Repeat Step 6 until the rebase completes successfully).*

## Next Steps
Once committed and successfully rebased, you can use the `/git-push` workflow to safely send the changes to the remote branch.
