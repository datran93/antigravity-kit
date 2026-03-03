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


---

## 🚨 MULTI-AGENT ARCHITECTURE FLOW

The workflow follows a multi-agent architectural pattern to handle requests robustly:
**User Request -> Planner Agent <-> Code Agent <-> Test Agent <-> Review Agent**

### 1. Planner Agent (Initial & Orchestration)
- **Role**: Understand user requests, discover context, and plan the execution.
- **Action**: Follows the `[/planner-architect](file://.agent/workflows/planner-architect.md)` workflow. Uses `@mcp:skill-router` (`search_skills`) to find relevant skills. Uses `@mcp:ast-explorer` and other discovery tools to map the impact area.
- **Output**: Calls `@mcp:context-manager` (`initialize_task_plan`) to persist the plan. Hands over tasks to Code Agent via `@mcp:mcp-multi-agent` (`delegate_to_subagent`).

### 2. Code Agent (Execution)
- **Role**: Implement tasks following **Clean Code** and **Testability** standards.
- **Action**: Follows the `[/coder-implementation](file://.agent/workflows/coder-implementation.md)` workflow. Writes, refactors, and engineers the solution based on the plan.
- **Output**: Upon completing a task, sends a message to the Reviewer for feedback. Does **NOT** mark tasks as complete.

### 3. Test Agent (Verification)
- **Role**: Check functionality, stability, and quality.
- **Action**: Follows the `[/tester-verification](file://.agent/workflows/tester-verification.md)` workflow. Runs tests, verification commands (`run_command`), and ensures code behaves correctly under edge cases.
- **Output**: Exchanges feedback with Code Agent for fixes. Once passing, sends to Review Agent.

### 4. Review Agent (Audit & Handover)
- **Role**: Final quality, security, and standards validation.
- **Action**: Follows the `[/reviewer-audit](file://.agent/workflows/reviewer-audit.md)` workflow. Audits the PR/changes, ensures adherence to SOTA & coding standards.
- **Output**: Synthesizes the final report. Hands back to the user or back to Planner for missed requirements.

### ⛔ INTERNAL COMMUNICATION & GOVERNANCE (DEADLOCK PREVENTION)

To prevent deadlocks (where all agents sleep) and resource conflicts, follow these rules:

1.  **Star Topology (Planner as Dispatcher)**: Every subagent (Coder, Tester, Reviewer) MUST report back to the **Planner** or the next designated role after every work block. The Planner is the "Source of Truth" for the pipeline state.
2.  **Mandatory Activation Message**: An active Agent MUST NOT exit or idle without calling `publish_message` to signal the next participant. If a task is blocked, notify the Planner.
3.  **Exclusive Resource Ownership**: 
    - **Planner**: Task Plan (`.agent/tasks.md`) & Checkpoints.
    - **Code Agent**: Source code (`src/`, etc.).
    - **Test Agent**: Test files and execution environments.
    - **Reviewer**: Read-only access.
4.  **Sequential Handoffs**: Use `run_background=False` for critical state transitions to ensure the parent process (Planner/Coordinator) maintains execution flow.
5.  **Atomic State Updates**: All writes to shared state must use the `@mcp:context-manager` tools sequentially.

---

### 🛑 Socratic Gate Protocol

Minimize user friction by providing choices:

- **New Feature**: "How should we handle X? [A] Option A, [B] Option B".
- **Bug Fix**: "Confirming impact: This fix resets Y. Proceed? [Yes/No]".
- **Vague**: "Objective seems to be Z. Is this for [1] Perf, [2] Security, or [3] UX?".
- **Critical Action Verification**: Before executing any destructive or critical action (e.g., database writes, bulk file deletions, 
  **production deployments), the Agent **MUST** explicitly describe the action and ask for user confirmation.

---

## 📌 Metadata

- **Version**: 1.0.1
- **Last Updated**: 2026-03-03
- **Maintainer**: Antigravity Team
- **Related**: `.agent/skills/*.md`, `GEMINI.md`
