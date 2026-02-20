---
description: Read GitLab Merge Request discussions, analyze requested changes, and apply fixes locally.
---
# Read and Fix GitLab MR Discussions Workflow

This workflow automates the first half of addressing a GitLab Merge Request (MR) feedback: reading the reviewer's feedback, analyzing the requested changes, and applying the fixes to the local codebase.

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

## Next Steps
Once verified, proceed to the `/gitlab-mr-reply-push` workflow to commit, push the changes, and resolve the discussions on GitLab.
