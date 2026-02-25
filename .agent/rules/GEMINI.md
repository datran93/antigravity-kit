---
trigger: always_on
---

# GEMINI.md - Antigravity Kit

## 🎯 UNIVERSAL CORE RULES

### 🌐 Communication

- **Bilingual Support**: Non-English prompt → Translate internally → Respond in user's language → Code and comments in
  **English**.

### 🛠️ MCP Tool Mastery (Priority Over Bash)

| MCP Server                       | Tool Name                  | Description & Usage                                                                       |
| :------------------------------- | :------------------------- | :---------------------------------------------------------------------------------------- |
| **`@mcp:ast-explorer`**          | `get_project_architecture` | Get AST-based structural overview. Extracts Classes/Functions/Methods for Py, Go, JS, TS. |
| **`@mcp:context-manager`**       | `save_checkpoint`          | Persist progress, completed/next steps, active files, and notes to `context.db`.          |
|                                  | `load_checkpoint`          | Recover context from a saved checkpoint ID.                                               |
|                                  | `list_active_tasks`        | List all tasks and their current status in the workspace.                                 |
| **`@mcp:database-inspector`**    | `list_tables`              | List all tables/views in SQL or sample keys in Redis.                                     |
|                                  | `inspect_schema`           | Get detailed schema (columns, types, PK/FK) for a table or Redis key.                     |
|                                  | `run_read_query`           | Execute read-only SQL/Redis queries.                                                      |
|                                  | `run_write_query`          | Execute write SQL/Redis queries. **MANDATORY**: Ask user first, then set `confirm=True`.  |
| **`@mcp:doc-researcher`**        | `search_latest_syntax`     | Real-time web search for SOTA syntax, best practices, and documentation.                  |
|                                  | `read_website_markdown`    | Scrape a specific URL and return its content as clean Markdown.                           |
| **`@mcp:figma-reader`**          | `read_figma_design`        | Read raw design structure, layers, and metadata from a Figma URL.                         |
|                                  | `export_figma_images`      | Render specific Figma nodes (layers) as temporary image URLs.                             |
|                                  | `get_design_details`       | Get deep JSON details for specific Figma node IDs.                                        |
| **`@mcp:gitlab-mr-discussions`** | `read_mr_discussions`      | Fetch all comment threads and resolution status from a GitLab MR.                         |
|                                  | `reply_to_mr_discussion`   | Post a reply to a specific GitLab discussion thread.                                      |
|                                  | `resolve_mr_discussion`    | Resolve or unresolve a discussion thread on GitLab.                                       |
| **`@mcp:skill-router`**          | `search_skills`            | Semantic search for the most relevant skills based on the task query.            |

### 🔧 MCP Tool Availability Check

- **Pre-flight**: Before using any MCP tool, verify availability with a simple call.
- **Graceful Degradation**: If MCP tool unavailable, fallback to standard tools:

  | MCP Tool                  | Fallback                                            |
  | ------------------------- | --------------------------------------------------- |
  | `@mcp:context-manager`    | Use `.agent/checkpoints/` directory with JSON files |
  | `@mcp:ast-explorer`       | Use `grep_search` + glob patterns                   |
  | `@mcp:doc-researcher`     | Use `codesearch` tool                               |
  | `@mcp:skill-router`       | Use `glob_search` for `.agent/skills/*.md`          |
  | `@mcp:database-inspector` | Use `usql` CLI tool via `run_command`               |

- **Log Warning**: Always inform user when falling back from MCP to standard tools.

---

## 🚨 MANDATORY 2-STEP PROCESSING FLOW

### Step 1: Skill Discovery & Loading

- **Semantic Search**: Use `@mcp:skill-router` (`search_skills`) as the **first action**.
- **Surgical Load**: Use `view_file` with `StartLine/EndLine` to read only necessary skill sections.
- **Fallback**: Use `clean-code` and general engineering if no specific skill matches.

### ⛔ ANTI-SKIP ENFORCEMENT

| Violation                 | Consequence                                                |
| :------------------------ | :--------------------------------------------------------- |
| **Skipped Step 1**        | No skills loaded → STOP, run `search_skills` first.        |
| **Skipped Context Phase** | Code is unguided → STOP, run architecture/usage discovery. |
| **No SOTA Research**      | Potential legacy code → STOP, run `search_latest_syntax`.  |
| **No Progress Plan**      | Workflow is untracked → STOP, provide atomic checklist.    |
| **No Checkpointing**      | No persistence → STOP, save checkpoint after major steps.  |

### Step 2: 4-Phase Execution Protocol

#### Phase 1: Context Discovery

- **Action**: Map the impact area.
- **Tools**: `@mcp:ast-explorer` (Architecture), `@mcp:database-inspector` (Data), `grep_search` (Usage).
- **SOTA**: Use `@mcp:doc-researcher` for modern syntax.

#### Phase 2: Progress Report (The Plan)

- **Action**: Provide a bulleted, atomic checklist.
- **Socratic Gate**: Ask clarifying questions via **Multiple Choice** options for vague specs.
- **Confirmation**: Wait for user approval for major changes.

#### Phase 3: Atomic Execution

- **Action**: Implement tasks one by one.
- **Persistence**: Call `@mcp:context-manager` (`save_checkpoint`) after completing each task.

#### Phase 4: Verification & Delivery

- **Action**: Verify changes via `run_command` (lint, build, test).
- **UI**: Instruct user to check UI in browser for design tasks.
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

### 🔐 Priority Hierarchy

1. **GEMINI.md** (Global Rules)
2. **SKILL.md** (Domain Specifics)
3. **Internal Documentation** (README, ADRs)

---

## ⚠️ Error Handling Protocol

### MCP Tool Failure

1. **Log** the error message clearly
2. **Fallback** to standard tools (see MCP Tool Availability Check)
3. **Continue** execution without blocking user workflow
4. **Report** to user: "MCP tool X unavailable, using Y instead"

### Partial Execution Recovery

1. **Checkpoint** files enable rollback to last known good state
2. **Resume** from last checkpoint ID: `load_checkpoint(id)`
3. **Document** failure point for debugging

### User Interruption

1. **Save** current progress immediately
2. **Create** checkpoint with descriptive name
3. **Summary** of completed vs pending tasks

---

## 📋 Summary Checklist

Before submitting any work, verify:

- [ ] Skills searched and loaded (Step 1)
- [ ] Context discovery completed (Phase 1)
- [ ] Progress plan shared and approved (Phase 2)
- [ ] Checkpoints saved after major tasks (Phase 3)
- [ ] Verification tests passed (Phase 4)

---

## 📌 Metadata

- **Version**: 1.0.0
- **Last Updated**: 2026-02-23
- **Maintainer**: Antigravity Team
- **Related**: `.agent/skills/*.md`, `AGENTS.md`
