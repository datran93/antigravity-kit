---
description: Push local code fixes to GitLab, then reply to and resolve the Merge Request discussions.
---
# Reply and Push GitLab MR Fixes Workflow

This workflow automates the second half of addressing a GitLab Merge Request (MR) feedback: committing and pushing the local changes, then replying to the reviewer and marking the discussions as resolved.

## Requirements
- You must have already completed the `/gitlab-mr-read-fix` workflow (or applied fixes manually).
- The `glab` CLI must be installed and authenticated.
- You must know the Merge Request IID (or be on the correct branch).

## Workflow Steps

1. **Commit and Push Changes**
   Stage the modified files, commit them with a descriptive message referencing the MR or discussion, and push to the remote branch. (You can also refer to the `/git-commit-push` workflow).
   ```bash
   git add .
   git commit -m "fix: address review comments on MR"
   git push origin $(git rev-parse --abbrev-ref HEAD)
   ```

2. **Reply to Discussions**
   Once the code is pushed, use `glab` to inform the reviewer.
   - To add a comment to the MR, use: `glab mr note [iid] -m "Fixed in the latest commit."`
   - If you are on the MR branch, you can omit the IID: `glab mr note -m "Fixed in the latest commit."`
   - To view the MR and resolve threads in the browser: `glab mr view -w`

## Final Verification
Confirm with the user that all targeted discussions have been replied to and marked as resolved, and that the changes were successfully pushed to GitLab.
