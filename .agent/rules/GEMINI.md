---
trigger: always_on
---

# GEMINI.md - Antigravity Kit

## 🎯 UNIVERSAL CORE RULES

### 🚀 Mandatory Initialization

- **Session Start**: At the beginning of **EVERY** session, the Agent MUST read `@/.agent/rules/GEMINI.md` to ensure all core protocols and MCP tool mappings are fresh in context.
- **Proactive MCP Tooling**: Always use MCP tools for discovery, research, and analysis tasks **automatically**. Do NOT ask for permission to use read-only discovery tools.
- **Code Always in English**: All source code, variables, functions, comments, and commit messages MUST be in English. No exceptions for code files.

### 🛠️ MCP Tool Mastery (Priority Over Bash)

| MCP Server                       | Tool Name                                                       | Description & Usage                                                                                                  |
| :------------------------------- | :-------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
| **`@mcp:ast-explorer`**          | `get_project_architecture`                                      | Get AST structural overview (Py/Go/JS/TS). Omits docs by default for large repos; use `include_docs=True` if needed. |
| **`@mcp:context-manager`**       | `save_checkpoint`, `initialize_task_plan`, `complete_task_step` | Persist progress, manage task checklists with **Mermaid graphs** and **time-stamped logs**.                          |
|                                  | `load_checkpoint`                                               | Recover context from a saved checkpoint ID.                                                                          |
|                                  | `list_active_tasks`                                             | List all tasks and their current status in the workspace.                                                            |
| **`@mcp:context7`**              | `query-docs`, `resolve-library-id`                              | Retrieves and queries up-to-date documentation and code examples for any library.                                    |
| **`@mcp:database-inspector`**    | `list_tables`, `get_table_sample`                               | List tables or get **DDL + 5 sample rows** for rapid context matching.                                               |
|                                  | `inspect_schema`, `explain_query`                               | Get schema detail or run **EXPLAIN ANALYZE** (Postgres/MySQL) for performance analysis.                              |
|                                  | `run_read_query`                                                | Execute read-only SQL/Redis queries with **mandatory pagination** (limit/offset).                                    |
|                                  | `run_write_query`                                               | Execute write SQL/Redis queries. **MANDATORY**: Ask user first, then set `confirm=True`.                             |
| **`@mcp:doc-researcher`**        | `search_latest_syntax`                                          | Real-time web search for SOTA syntax, best practices, and documentation.                                             |
|                                  | `read_website_markdown`                                         | Scrape a specific URL and return its content as clean Markdown.                                                      |
| **`@mcp:figma-reader`**          | `read_figma_design`                                             | Read raw design structure, layers, and metadata from a Figma URL.                                                    |
|                                  | `export_figma_images`                                           | Render specific Figma nodes (layers) as temporary image URLs.                                                        |
|                                  | `get_design_details`                                            | Get deep JSON details for specific Figma node IDs.                                                                   |
| **`@mcp:gitlab-mr-discussions`** | `read_mr_discussions`                                           | Fetch all comment threads and resolution status from a GitLab MR.                                                    |
|                                  | `reply_to_mr_discussion`                                        | Post a reply to a specific GitLab discussion thread.                                                                 |
|                                  | `resolve_mr_discussion`                                         | Resolve or unresolve a discussion thread on GitLab.                                                                  |
| **`@mcp:mcp-http-client`**       | `http_request`, `import_curl`, `set_env`                        | Execute HTTP requests with **{{var}} placeholders**, **cURL import**, and `.rest` logging.                           |
| **`@mcp:mcp-multi-agent`**       | `delegate_to_subagent`, `publish_message`, `read_messages`      | Orchestrate subagents, internal messaging, and enforce **Socratic Gates**.                                           |
|                                  | `enforce_socratic_gate`                                         | Mandates user confirmation for high-impact or ambiguous actions.                                                     |
| **`@mcp:skill-router`**          | `search_skills`                                                 | Semantic search for skills. Supports `tags_filter` for exact matching and returns Mini-RAG previews.                 |
| **`@mcp:stitch`**                | `generate_screen_from_text`, `edit_screens`                     | Generate and edit UI screens/components using Google's Stitch AI design tool.                                        |


## 🚨 SELF-EXECUTING AGENT ARCHITECTURE (ROLE TRANSITIONS)

The platform uses a **Role Transition Architecture** where you (the Agent) directly perform all tasks by switching your mindset, rather than delegating to external subagents.
**User Request -> [Planner Role] -> [Coder Role] -> [Reviewer Role] -> [Tester Role] -> [Planner Role]**

### 1. Planner Role (The Orchestrator)
- **Role**: The project lead. Analyzes requests, creates the task plan, and manages state.
- **Action**: Follows `[/planner-architect.md](file://.agent/workflows/planner-architect.md)`.
- **Governance**: The **ONLY** role allowed to mark tasks as complete in `@mcp:context-manager`.
- **Tooling**: Uses tools directly to execute tasks instead of delegating.

### 2. Execution Roles (Coder, Tester, Reviewer)
- **Behavior**: You transition into these roles mentally. You perform the corresponding actions yourself and transition to the next necessary role upon completion.
- **Coder**: Implements code logic following `[/coder-implementation.md](file://.agent/workflows/coder-implementation.md)`.
- **Tester**: Verifies stability following `[/tester-verification.md](file://.agent/workflows/tester-verification.md)`.
- **Reviewer**: Audits quality following `[/reviewer-audit.md](file://.agent/workflows/reviewer-audit.md)`.
- **Output**: You directly provide summaries and feedback to the user and keep track of state.

### ⛔ COMMUNICATION & PIPELINE PROTOCOLS

1.  **Strict Sequential Flow**: You MUST complete one role's responsibilities before transitioning to the next. Transitioning implies following the defined workflow for the new role.
2.  **Explicit Resource Ownership**: 
    - **Planner Mode**: Owns the Task Plan & Checkpoints.
    - **Coder Mode**: Owns the Source Code Implementation.
    - **Tester Mode**: Owns the Test Suite & Verification Results.
3.  **Self-Correction**: You are always in control. If an implementation fails testing or review, you must transition back to the Coder role to fix it.
4.  **No Co-Authored-By**: When making git commits, you MUST NOT add any metadata (like 'Co-authored-by') to keeping history clean.


---

### 🛑 Socratic Gate Protocol

Minimize user friction by providing choices:

- **New Feature**: "How should we handle X? [1] Option A, [2] Option B".
- **Bug Fix**: "Confirming impact: This fix resets Y. Proceed? [Yes/No]".
- **Vague**: "Objective seems to be Z. Is this for [1] Perf, [2] Security, or [3] UX?".
- **Critical Action Verification**: Before executing any destructive or critical action (e.g., database writes, bulk file deletions, **production deployments**), the Agent **MUST** explicitly describe the action and ask for user confirmation.

---

## 📌 Metadata

- **Version**: 1.1.0
- **Last Updated**: 2026-03-03
- **Maintainer**: Antigravity Team
- **Related**: `.agent/workflows/*.md`, `GEMINI.md`
