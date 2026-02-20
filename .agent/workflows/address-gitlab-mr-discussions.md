---
description: Workflow to read GitLab Merge Request discussions, address the code review feedback, push the changes, and reply to or resolve the discussions.
---
# Address GitLab MR Discussions Workflow

This workflow automates the process of reading feedback from a GitLab Merge Request (MR), making the necessary code changes locally, pushing those changes, and responding to the reviewer.

## Requirements
- The user must provide the Merge Request IID (e.g., `!123` or just `123`) or the Agent should be able to deduce the current MR from the branch context.
- The GitLab MCP server must be configured and authenticated.

## Workflow Steps

1. **Get Discussions & Notes from the Merge Request**
   Use the GitLab MCP tools to fetch the discussions for the specified MR.
   - Tool: `mcp_gitlab_mr_discussions` or `mcp_gitlab_get_merge_request_notes`
   - Parameters: `merge_request_iid`
   
   *Agent Action*: Read through unresolved threads and notes. Identify the specific files and lines of code mentioned by the reviewers and the requested changes.

2. **Analyze and Plan the Fixes**
   For each piece of feedback:
   - Identify the target file.
   - Use `view_file` to read the current state of the code.
   - Plan how to implement the requested change without breaking existing functionality.

3. **Implement Code Changes**
   Apply the changes to the codebase based on the plan.
   - Tools: `replace_file_content` or `multi_replace_file_content`
   - Ensure you follow the project's coding standards.

4. **Verify Changes (Lint/Test)**
   Run local linters or tests to ensure the new code is correct and hasn't introduced regressions.
   ```bash
   # Example verification commands (adapt based on project)
   # npm run lint
   # npm test
   # go test ./...
   ```

5. **Commit and Push Changes**
   Stage the modified files, commit them with a descriptive message referencing the MR or discussion, and push to the remote branch. (You can also refer to the `/git-commit-push` workflow).
   ```bash
   git add .
   git commit -m "fix: address review comments on MR"
   git push origin $(git rev-parse --abbrev-ref HEAD)
   ```

6. **Reply to and Resolve Discussions**
   Once the code is pushed, go back to the GitLab discussions and inform the reviewer.
   - To reply to a thread, use `mcp_gitlab_create_merge_request_discussion_note` with the `discussion_id` and a message like: *"Fixed in the latest commit."*
   - To resolve the thread, use `mcp_gitlab_resolve_merge_request_thread` with `resolved: true`.

## Final Verification
Confirm with the user that all targeted discussions have been addressed, resolved, and changes are successfully pushed to GitLab.
