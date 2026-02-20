---
description: Push local code fixes to GitLab, then reply to and resolve the Merge Request discussions.
---
# Reply and Push GitLab MR Fixes Workflow

This workflow automates the second half of addressing a GitLab Merge Request (MR) feedback: committing and pushing the local changes, then replying to the reviewer and marking the discussions as resolved.

## Requirements
- You must have already completed the `/gitlab-mr-read-fix` workflow (or applied fixes manually).
- The GitLab MCP server must be configured and authenticated.
- You must know the Merge Request IID.

## Workflow Steps

1. **Commit and Push Changes**
   Stage the modified files, commit them with a descriptive message referencing the MR or discussion, and push to the remote branch. (You can also refer to the `/git-commit-push` workflow).
   ```bash
   git add .
   git commit -m "fix: address review comments on MR"
   git push origin $(git rev-parse --abbrev-ref HEAD)
   ```

2. **Reply to and Resolve Discussions**
   Once the code is pushed, go back to the GitLab discussions and inform the reviewer.
   - To reply to a thread, use `mcp_gitlab_create_merge_request_discussion_note` with the `discussion_id` and a message like: *"Fixed in the latest commit."*
   - To resolve the thread, use `mcp_gitlab_resolve_merge_request_thread` with `resolved: true`.

## Final Verification
Confirm with the user that all targeted discussions have been replied to and marked as resolved, and that the changes were successfully pushed to GitLab.
