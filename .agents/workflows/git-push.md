---
description: An optimized Git workflow to safely pull, rebase against develop, and push changes to the remote branch, with automatic conflict resolution.
---
# /git-push - Safe Git Push and Rebase Workflow

This workflow guides the agent to safely fetch the latest changes from the remote repository, **rebase the local changes against the `develop` branch** to prevent conflicts during merging, resolve any conflicts intelligently, and push to the remote branch.

## When to Use
- `/git-push`
- Keywords: "push my code", "sync with develop", "rebase and push"

## Requirements
- You must have committed your changes (see `/git-commit`).
- The remote should be correctly set up (e.g., `origin`).

## Workflow Steps

1. **Fetch Latest Changes**
   // turbo
   Fetch the latest updates from the remote repository to ensure you have the most recent version of `develop`.
   ```bash
   git fetch origin
   ```

2. **Rebase Against Develop**
   // turbo
   Rebase the current feature branch on top of `origin/develop`. This rewrites your local commits so they are applied after the latest changes in `develop`, minimizing conflicts when the Merge Request is eventually merged.
   ```bash
   git rebase origin/develop
   ```

3. **Auto-Resolve Conflicts (If Step 2 Fails)**
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
   *(Repeat Step 3 until the rebase completes successfully).*

4. **Push to Remote**
   Once the rebase is successful and history is clean, push the changes. Because the history has been rewritten during the rebase, you must use `--force-with-lease` for safety.
   ```bash
   git push --force-with-lease origin $(git rev-parse --abbrev-ref HEAD)
   ```
   *Note: If pushing to a completely new branch, a simple `git push -u origin <branch>` might be required if `--force-with-lease` complains about a missing upstream.*

## Final Verification
Verify that the branch is successfully pushed and clean.
```bash
git status
```
