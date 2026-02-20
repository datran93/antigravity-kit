---
description: An optimized Git workflow to safely commit, merge, and push changes to GitLab while minimizing conflicts and auto-resolving them if they happen.
---
# Git Commit and Push Workflow

This workflow is designed to safely commit your local changes, pull the latest changes from the remote GitLab repository to avoid conflicts, and push your updates. If conflicts happen, the agent will analyze and resolve them automatically.

## Requirements
- You must be on the branch you want to push.
- The remote should be correctly set up (usually `origin`).

## Workflow Steps

1. **Check Current Status**
   Check the current branch and see what files have been modified.
   ```bash
   git status
   ```

2. **Stage and Commit Changes**
   Add all changes and commit them with a meaningful, conventional commit message.
   *Note: If taking input from the user, ask for the commit message if it's not obvious.*
   ```bash
   git add .
   git commit -m "chore: update files"
   ```

3. **Fetch and Rebase (Prevent Conflicts)**
   // turbo
   Pull the latest changes from the remote branch using rebase to maintain a clean linear history and avoid unnecessary merge commits. Replace `<branch_name>` with the current branch name.
   ```bash
   git pull --rebase origin $(git rev-parse --abbrev-ref HEAD)
   ```

4. **Auto-Resolve Conflicts (If Step 3 Fails)**
   If the `git pull --rebase` fails due to conflicts:
   - Run `git status` to identify the files with conflicts.
   - Use the `view_file` tool to read the conflicted files (they will have `<<<<<<<`, `=======`, `>>>>>>>` markers).
   - Analyze the conflicting changes logically. Combine the remote changes (theirs) and local changes (ours) appropriately so no logic is lost.
   - Use the `replace_file_content` or `multi_replace_file_content` tool to remove the conflict markers and write the resolved code.
   - Stage the resolved files: `git add <resolved_file>`
   - Continue the rebase:
     ```bash
     git rebase --continue
     ```
   *(Repeat Step 4 until the rebase is completely finished).*

5. **Push to GitLab**
   Once the rebase is successful and history is clean, push the changes to GitLab.
   ```bash
   git push origin $(git rev-parse --abbrev-ref HEAD)
   ```

## Final Verification
Verify that the branch is up to date and clean.
```bash
git status
```
