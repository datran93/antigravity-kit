---
trigger: always_on
---

# GEMINI.md - Antigravity Kit

## 🎯 UNIVERSAL CORE RULES

### 🚀 Mandatory Initialization

- **Session Start**: At the beginning of **EVERY** session, the Agent MUST read `@/.agents/rules/GEMINI.md` to ensure
  all core protocols are fresh, and MUST read `@/.agents/rules/ANCHORS.md` to load the immutable facts and project
  guardrails.
- **Proactive MCP Tooling**: Always use MCP tools for discovery, research, and analysis tasks **automatically**. Do NOT
  ask for permission to use read-only discovery tools.

### 🛠️ MCP Tool Mastery (Priority Over Bash)

| MCP Server                       | Tool Name                                                                        | Description & Usage                                                                                        |
| :------------------------------- | :------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------- |
| **`@mcp:ast-explorer`**          | `get_project_architecture`, `search_symbol`                                      | Get AST structural overview or search for class/function symbols. Omits docs by default for large repos.   |
| **`@mcp:context-manager`**       | `save_checkpoint`, `initialize_task_plan`, `complete_task_step`, `add_task_step` | Persist progress, manage task checklists with **time-stamped logs**. Add new steps sequentially if needed. |
|                                  | `load_checkpoint`, `list_active_tasks`                                           | Recover context from a saved checkpoint ID or list tasks.                                                  |
|                                  | `declare_intent`, `check_intent_lock`                                            | Apply Intent Locking to prevent blind writes and scope creep.                                              |
|                                  | `compact_memory`, `recall_knowledge`                                             | Automatically extract and index KIs into Local RAG, or recall them via SQLite FTS5.                        |
|                                  | `record_failure`, `clear_drift`                                                  | Track successive test or logic failures to detect drift and trigger `think_back`.                          |
|                                  | `manage_anchors`, `annotate_file`                                                | Retrieve/set immutable Anchor rules, or inject Ghost Context (gotchas) into specific files.                |
| **`@mcp:context7`**              | `query-docs`, `resolve-library-id`                                               | Retrieves and queries up-to-date documentation and code examples for any library.                          |
| **`@mcp:database-inspector`**    | `list_tables`, `get_table_sample`                                                | List tables or get **DDL + 5 sample rows** for rapid context matching.                                     |
|                                  | `inspect_schema`, `explain_query`                                                | Get schema detail or run **EXPLAIN ANALYZE** (Postgres/MySQL) for performance analysis.                    |
|                                  | `run_read_query`                                                                 | Execute read-only SQL/Redis queries with **mandatory pagination** (limit/offset).                          |
|                                  | `run_write_query`                                                                | Execute write SQL/Redis queries. **MANDATORY**: Ask user first, then set `confirm=True`.                   |
| **`@mcp:doc-researcher`**        | `search_latest_syntax`, `read_doc_file`                                          | Real-time web search for SOTA syntax, and read local `.doc`, `.docx`, `.pdf` files with pagination.        |
|                                  | `read_website_markdown`                                                          | Scrape a specific URL and return its content as clean Markdown.                                            |
| **`@mcp:figma-reader`**          | `read_figma_design`                                                              | Read raw design structure, layers, and metadata from a Figma URL.                                          |
|                                  | `export_figma_images`                                                            | Render specific Figma nodes (layers) as temporary image URLs.                                              |
|                                  | `get_design_details`                                                             | Get deep JSON details for specific Figma node IDs.                                                         |
| **`@mcp:gitlab-mr-discussions`** | `read_mr_discussions`                                                            | Fetch all comment threads and resolution status from a GitLab MR.                                          |
|                                  | `reply_to_mr_discussion`                                                         | Post a reply to a specific GitLab discussion thread.                                                       |
|                                  | `resolve_mr_discussion`                                                          | Resolve or unresolve a discussion thread on GitLab.                                                        |
| **`@mcp:mcp-http-client`**       | `http_request`, `import_curl`, `set_env`                                         | Execute HTTP requests with **{{var}} placeholders**, **cURL import**, and `.rest` logging.                 |
|                                  | `clear_history`, `list_history`, `set_config`                                    | View or clear request history, configure base URL and auth token.                                          |
| **`@mcp:skill-router`**          | `search_skills`                                                                  | Semantic search for skills. Supports `tags_filter` for exact matching and returns Mini-RAG previews.       |
| **`@mcp:stitch`**                | `generate_screen_from_text`, `edit_screens`, `generate_variants`                 | Generate, edit, and create variants of UI screens using Google's Stitch AI design tool.                    |
|                                  | `create_project`, `get_project`, `list_projects`                                 | Manage Stitch AI projects (create, get details, list).                                                     |
|                                  | `get_screen`, `list_screens`                                                     | Read specific UI screens from existing Stitch AI projects.                                                 |

## 🚨 SELF-EXECUTING AGENT ARCHITECTURE (ROLE TRANSITIONS)

The platform uses a **Role Transition Architecture** where you (the Agent) directly perform all tasks by switching your
mindset, rather than delegating to external subagents.

**User Request -> [Spec Writer] -> [Planner] <-> [Coder] <-> [Reviewer] <-> [Tester]**

### 1. The Roles

- **Pre-Planning (Spec Writer)**: `[/specifications-writer.md](file://.agents/workflows/specifications-writer.md)` -
  Requirements engineering & Socratic questioning.
- **Planner (Orchestrator)**: `[/planner-architect.md](file://.agents/workflows/planner-architect.md)` - Architecture,
  state management, and 3-Tier Task orchestration.
- **Coder**: `[/coder-implementation.md](file://.agents/workflows/coder-implementation.md)` - File edits, Intent
  Locking, testability constraints.
- **Reviewer**: `[/reviewer-audit.md](file://.agents/workflows/reviewer-audit.md)` - Semantic tracing against SPEC and
  Code Quality validation.
- **Tester**: `[/tester-verification.md](file://.agents/workflows/tester-verification.md)` - Generating tests, measuring
  coverage (>= 70%), checking CI mechanics.

## ⛔ CORE COMMUNICATION & PROTOCOLS

1. **Strict Sequential Flow**: Complete one role's responsibilities perfectly before transitioning.
2. **Explicit Resource Ownership**:
   - **Planner**: Owns the Task Plan & Context Architecture.
   - **Coder**: Owns the Source Code & Implementation.
   - **Tester**: Owns the Test Suite.
3. **Role Anchoring**: ALWAYS prefix every conversational response with your current role tag (e.g.
   `[Role: 🏗️ Planner]`, `[Role: 💻 Coder]`).
4. **Skill Transparency**: Explicitly print out utilized specialized skills BEFORE executing tasks.

## 🛡️ UNIVERSAL GUARDRAILS (Applies to ALL Roles)

To prevent infinite loops, duplicated instructions, and technical debt, enforce these guardrails across any role:

### A. Drift Detection (Panic Protocol)

If a role is stuck in a failure cycle for the **same core issue 3 times** (e.g. Coder failing the same test, Reviewer
rejecting the same design):

- Stop endless fixing natively.
- Call `@mcp:context-manager` (`record_failure`).
- Transition to `[Role: 🏗️ Planner]` and execute a **Lateral Pivot** (Simplifier, Contrarian, or Hacker personas) to
  discuss the systemic blocker with the USER.

### B. Export Intelligence & Ghost Context

You must never hand over a blank slate to the next role.

- Before a role transition, use `@mcp:context-manager` (`annotate_file`) to inject non-obvious gotchas, architectural
  quirks, or library limitations _directly_ to the involved files.
- Summarize findings and explicitly pass them in your transition message so the next role processes it instantly.

### C. 3-Stage Evaluation Pipeline

Every Tactic completion must pass this pipeline seamlessly:

- **Stage 1 (Mechanical)**: Managed by `Tester`. Lint, Compile, Tests == PASS. Coverage >= 70%.
- **Stage 2 (Semantic)**: Managed by `Reviewer`. Does the output directly fulfill the Bounded Contexts of the `SPEC.md`?
- **Stage 3 (Consensus/Frontier)**: Managed by `Planner`. High-impact security and architectural coherence checks.

### D. Governance Modes

If `--mode=coach` or `--mode=strict` is applied, block direct modifications without rigid testing and Socratic
interrogation.

---

## 📌 Metadata

- **Version**: 1.3.0
- **Last Updated**: 2026-03-11
- **Related**: `.agents/workflows/*.md`, `GEMINI.md`
