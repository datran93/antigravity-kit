---
trigger: always_on
---

# GEMINI.md - Antigravity Kit

## 🎯 UNIVERSAL CORE RULES

### 🌐 Communication

- **Bilingual Support**: Non-English prompt → Translate internally → Respond in user's language → Code and comments in
  **English**.
- **Style**: Professional, technically accurate, concise. Use headers, bolding, and backticks.
- **Proactiveness**: Take follow-up actions (lint/test) automatically **UNLESS in Plan Mode**.

### 🛠️ MCP Tool Mastery (Priority Over Bash)

1. **Memory**: `@mcp:context-manager` (`save_checkpoint`, `load_checkpoint`) for any multi-file task.
2. **Discovery**: `@mcp:ast-explorer` (`get_project_architecture`) for mapping; `@mcp:database-inspector` for DB schema/data.
3. **Research**: `@mcp:doc-researcher` (`search_latest_syntax`) before writing new logic to avoid legacy code.
4. **Standard Tools**: Use `view_file`, `grep_search`, `list_dir`, `replace_file_content` instead of raw bash (`cat`, `ls`, `sed`).

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

### 🚦 Mode Awareness

- **Plan Mode (Read-Only)**: Only observe, analyze, plan. NO file edits or system changes.
- **Execute Mode**: Full execution capabilities including edits, commits, and commands.
- **Detection**: Check system reminders for active mode constraints.

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
- **Persistence**: Call `@mcp:context-manager` (`save_checkpoint`) after completing each major task or file.

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
production deployments), the Agent **MUST** explicitly describe the action and ask for user confirmation.

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
