---
description: Push local code fixes to GitLab, then reply to and resolve the Merge Request discussions.
---
# Reply and Push GitLab MR Fixes Workflow

This workflow automates the second half of addressing a GitLab Merge Request (MR) feedback: committing and pushing the local changes, then replying to the reviewer and marking the discussions as resolved.

## Requirements
- You must have already completed the `/gitlab-mr-read-fix` workflow (or applied fixes manually).
- The `@mcp:gitlab-mr-discussions` server must be configured.
- You must know the Project ID, Merge Request IID, and the specific Discussion IDs you wish to reply to.

## Workflow Steps

1. **Commit and Push Changes**
   Stage the modified files, commit them with a descriptive message referencing the MR or discussion, and push to the remote branch. (You can also refer to the `/git-commit` and `/git-push` workflows).
   ```bash
   git add .
   git commit -m "fix: address review comments on MR"
   git push origin $(git rev-parse --abbrev-ref HEAD)
   ```

2. **Reply to Discussions and Resolve**
   Once the code is pushed, use the `@mcp:gitlab-mr-discussions` tools to inform the reviewer.
   - Use `reply_to_mr_discussion` with the `project_id`, `mr_iid`, `discussion_id`, and your `body` text (e.g., "Fixed in the latest commit.").
   - Use `resolve_mr_discussion` with `resolve=True` to explicitly close the thread if the feedback has been fully addressed.

## Final Verification
Confirm with the user that all targeted discussions have been replied to and marked as resolved, and that the changes were successfully pushed to GitLab.
