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
| **`@mcp:notebooklm`**            | `notebook_query`, `research_start`                              | Query NotebookLM notebooks for source-grounded insights and start deep research.                                     |
| **`@mcp:skill-router`**          | `search_skills`                                                 | Semantic search for the most relevant skills based on the task query.                                                |
| **`@mcp:stitch`**                | `generate_screen_from_text`                                     | Generate and edit UI screens/components using Google's Stitch AI design tool.                                        |

- **Graceful Degradation**: If MCP tool unavailable, fallback to standard tools:

---

## 🚨 MANDATORY 2-STEP PROCESSING FLOW

### Step 1: Skill Discovery & Loading

- **Semantic Search**: Use `@mcp:skill-router` (`search_skills`) as the **first action**.
- **Surgical Load**: Use `view_file` with `StartLine/EndLine` to read only necessary skill sections.
- **Fallback**: Use `clean-code` and general engineering if no specific skill matches.

### ⛔ ANTI-SKIP ENFORCEMENT

| Violation                 | Consequence                                                               |
| :------------------------ | :------------------------------------------------------------------------ |
| **Skipped Step 1**        | No skills loaded → STOP, run `search_skills` first.                       |
| **Skipped Context Phase** | Code is unguided → STOP, run architecture/usage discovery.                |
| **No SOTA Research**      | Potential legacy code → STOP, run `query-docs` or `search_latest_syntax`. |
| **No Progress Plan**      | Workflow is untracked → STOP, manage task checklists with progress bars.  |
| **No Checkpointing**      | No persistence → STOP, save checkpoint after major steps.                 |

### Step 2: 4-Phase Execution Protocol

#### Phase 1: Context Discovery

- **Action**: Map the impact area.
- **Tools**: `@mcp:ast-explorer` (Architecture), `@mcp:database-inspector` (Data), `grep_search` (Usage).
- **SOTA Research Hierarchy**:
    1. **`@mcp:context7`**: Primary source for modern framework/library syntax.
    2. **`@mcp:doc-researcher`**: Target specific documentation URLs.
    3. **`search_web`**: Troubleshooting latest bugs/issues/community fixes.

#### Phase 2: Progress Report (The Plan)

- **Socratic Gate**: Ask clarifying questions via **Multiple Choice** options for vague specs.
- **Action**: Brainstorm a plan, then call `@mcp:context-manager` (`initialize_task_plan`) to persist the plan and show the progress bar.
- **Confirmation**: Wait for user approval for major changes.

#### Phase 3: Atomic Execution

- **Action**: Implement tasks one by one.
- **Persistence**: Call `@mcp:context-manager` (`complete_task_step`) after completing each task to update the checklist and progress.
- **Checkpointing**: Use `save_checkpoint` for major architectural changes or state captures.

#### Phase 4: Verification & Delivery

- **Action**: Verify changes via `run_command` (lint, build, test).
- **Handover**: Summarize work done and link to relevant files.

---

### 🛑 Socratic Gate Protocol

Minimize user friction by providing choices:

- **New Feature**: "How should we handle X? [A] Option A, [B] Option B".
- **Bug Fix**: "Confirming impact: This fix resets Y. Proceed? [Yes/No]".
- **Vague**: "Objective seems to be Z. Is this for [1] Perf, [2] Security, or [3] UX?".
- **Critical Action Verification**: Before executing any destructive or critical action (e.g., database writes, bulk file deletions, 
  **production deployments), the Agent **MUST** explicitly describe the action and ask for user confirmation.

---

### 🌐 Linguistic & Operational Standards

- **Progress-First Strategy**: For any task with **3+ steps**, MUST call `initialize_task_plan` immediately before starting Phase 3.
- **API History Persistence**: All HTTP requests MUST be logged to `rest/{slug}.rest` for auditability and manual replay.
- **Critical Action Verification**: Before executing `run_write_query` or data-modifying `http_request` (DELETE/PUT/PATCH), MUST provide a summary table of changes and wait for explicit confirmation.

---

### Partial Execution Recovery

1. **Checkpoint** files enable rollback to last known good state
2. **Resume** from last checkpoint ID: `load_checkpoint(id)`
3. **Document** failure point for debugging

### User Interruption

1. **Save** current progress immediately
2. **Create** checkpoint with descriptive name
3. **Summary** of completed vs pending tasks

---

## 📌 Metadata

- **Version**: 1.0.0
- **Last Updated**: 2026-02-28
- **Maintainer**: Antigravity Team
- **Related**: `.agent/skills/*.md`, `GEMINI.md`
