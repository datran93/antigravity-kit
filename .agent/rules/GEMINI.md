---
trigger: always_on
---

# GEMINI.md - Antigravity Kit

## UNIVERSAL RULES

### 🌐 Language

Non-English prompt → Translate internally → Respond in user's language → Code and all comments in English

### 🛠️ Tool & MCP Server Mastery

> **CRITICAL RULE:** Always prioritize using specialized MCP servers and specific tools over generic terminal bash operations. 

**MCP Server Protocols:**
1. **`@mcp:skill-router`**: Use `search_skills` as your **very first action** when determining how to solve a problem to retrieve the right domain knowledge.
2. **`@mcp:database-inspector`**: Use `inspect_schema`, `list_tables`, and `run_read_query` to query databases (SQL, Redis, etc.) safely instead of running raw terminal commands like `psql` or `redis-cli`.
3. **`@mcp:doc-researcher`**: Use `search_latest_syntax` before writing code for new features to avoid generating legacy code or using deprecated APIs. Do not rely solely on your internal training data.
4. **`@mcp:context-manager`**: Use `list_active_tasks`, `load_checkpoint`, and `save_checkpoint` to persist memory during complex, multi-file agentic tasks.
5. **`@mcp:gitlab-mr-discussions`**: Use for interacting, reading, and resolving GitLab threads directly without needing the UI.

**System Tool Rules:**
- **NEVER** use `cat`, `grep`, `ls`, or `sed` inside a bash command if native tools (`view_file`, `grep_search`, `list_dir`, `replace_file_content`) are available.
- Always read existing files before modifying them (`view_file`).
- Limit `run_command` to actual application logic, testing, or building, not generic file parsing.
### 🧹 Clean Code

**ALL code follows `@[skills/clean-code]`.** Concise, self-documenting. Testing mandatory (Pyramid + AAA). Measure performance first.

## 🚨 MANDATORY: 3-STEP PROCESSING FLOW

You must complete the following steps in order and NEVER skip them:
1. **CLASSIFY REQUEST**: Determine Type and Tier.
2. **SKILL ROUTING**: Semantic search and apply SKILL.md.
3. **EXECUTE TASK**: Context Discovery -> SOTA Research -> Execute -> Verify -> Deliver.

### ⛔ ANTI-SKIP ENFORCEMENT

| Violation                          | Consequence                                                   |
| ---------------------------------- | ------------------------------------------------------------- |
| Skipped Step 1 (no classification) | Response is INVALID → Go back, classify first                 |
| Skipped Step 2 (no skills loaded)  | Response lacks depth → Re-run search_skills, enhance          |
| Skipped SOTA Research              | Code is potentially LEGACY → STOP, run `search_latest_syntax` |
| Started code before Context check  | Code is UNGUIDED → Stop, survey codebase first                |

### 🔐 Priority Hierarchy (BINDING)

```
P0: GEMINI.md (this file) → ALWAYS applies, cannot be overridden
P1: SKILL.md files        → Detailed patterns and techniques (Single Source of Truth)
```
---

## 📥 STEP 1: REQUEST CLASSIFIER

**Analyze keywords + context → Determine Type → Set Execution Tier:**

| Request Type      | Decision Heuristics (Rule of Thumb)             | Trigger Keywords                              | Tier / Mode     |
| :---------------- | :---------------------------------------------- | :-------------------------------------------- | :-------------- |
| **CLARIFICATION** | Informational/Conceptual. No code changes.      | "what is", "how", "why"                       | T0 (Direct)     |
| **SURVEY/INTEL**  | Analysis of state/code/repo. No implementation. | "analyze", "audit", "find"                    | T0 + Explorer   |
| **OPTIMIZATION**  | Improving existing code status/perf.            | "refactor", "cleanup", "optimize"             | T1+ (Execution) |
| **SIMPLE CODE**   | Specific fix/add restricted to 1 file.          | "fix", "add", "update"                        | T1+ (Execution) |
| **COMPLEX CODE**  | Feature creation affecting multiple files.      | "build", "create", "implement"                | Full (Agent)    |
| **DESIGN/UI**     | Visual/UX focus, dashboard/component styles.    | "design", "ui", "premium"                     | Full (Agent)    |
| **SYSTEM/SYNC**   | Infrastructure, ENV, or script automation.      | "sync", "setup", "env", "script"              | T1+ (Execution) |
| **SLASH CMD**     | Workflow trigger using /command syntax.         | /create, /orchestrate, /debug, /plan, /update | Workflow Mode   |

### 📊 Tier Assessment Matrix

| Tier     | Definition & Complexity                                              | Requirements                                       |
| :------- | :------------------------------------------------------------------- | :------------------------------------------------- |
| **T0**   | **Knowledge Only**. Pure information retrieval or basic explanation. | Direct response.                                   |
| **T1+**  | **Implementation Light**. Changes to existing logic or 1-3 files.    | Socratic Gate (Multiple Choice).                   |
| **Full** | **Systemic Build**. New features, complex logic, multi-file arch.    | Implementation Plan + Socratic Gate + Checkpoints. |

**Output format after classification:**

```markdown
📥 **Request Type:** [TYPE] → [TIER]
```
---

## 📚 STEP 2: SKILL DISCOVERY & LOADING PROTOCOL

> 🔴 **MANDATORY:** You MUST prioritize using the `search_skills` tool from the `@mcp:skill-router` server to find relevant skills based on semantic meaning.

**Find and load skills selectively:**

1.  **Semantic Search (Primary)**: Use the `search_skills` tool with the user's intent or query to find the top relevant skills. This is the Single Source of Truth.
2.  **Directory Search (Fallback)**: If MCP routing fails, use `find_by_name` on `.agent/skills/` to match user request keywords against folder names.
3.  **Audit (Discovery)**: Use `view_file_outline` on each identified `SKILL.md` path (returned by the router) to see its structure.
4.  **Precision Search**: If headers are ambiguous, use `grep_search` within the skill directory for user request keywords.
5.  **Targeted Load**: Use `view_file` with `StartLine` and `EndLine` to read ONLY the relevant sections.
6.  **Fallback Mechanism**: If no specific skill matches perfectly, fallback to `clean-code` and general standard engineering. DO NOT hallucinate a skill.

> [!IMPORTANT] **Avoid loading entire SKILL.md files.** Loading 500+ lines of generic patterns for a 5-line fix is a
> waste of tokens and context. Be surgical.

**Output format after skill selection:**

```markdown
🤖 **Applying skills: `[skill-names]`...**
```
---

## ⚡ STEP 3: TASK EXECUTION (4-Phase Execution)

**Now you may proceed with the actual work following these 4 phases:**

### Phase 1: Context First (Discovery)

- **Do not write code immediately.**
- **Working Memory:** Use `list_active_tasks` (from `@mcp:context-manager`) to check for pending background work, or `load_checkpoint` if you are explicitly resuming a previous task.
- **🚨 SOTA Research (MANDATORY)**: ALWAYS use `search_latest_syntax` or `read_website_markdown` (from `@mcp:doc-researcher`) before implementing new features or using external libraries. You MUST ensure you write modern code and avoid deprecated APIs. Skipping this for any Tier 1+ or Full request is a violation.
- Read existing `README.md`, `.cursorrules`, `.clinerules`, or scan relevant architecture files using `grep_search` / `view_file`.
- Understand the current state of the codebase to retain consistency.

### Phase 2: Execute with Skills

- For **T0 (Questions)**: Respond directly using loaded knowledge. No Socratic Gate required.
- For **T1+ and Full**: Apply Socratic Gate if vague. If clear, implement the solution using the patterns learned from
  the loaded skills.

### Phase 3: Trust, but Verify

- **You must verify your changes.**
- Run linter (e.g., `npm run lint`), compiler (e.g., `tsc`), or tests if available before claiming success.
- **UI/Web Apps:** If working on frontend code, instruct the user to verify the UI/layout in the browser before Checkpointing. Do not blindly assume visual correctness.
- Ensure no regressions are introduced.

### Phase 4: Deliver & Checkpoint

- For small tasks, deliver the final atomic change.
- For **Full (Agent)** multi-file tasks: Report progress after each file or component as a "Checkpoint" before
  proceeding. Don't process 5+ complex files in a single burst without saving state and asking/confirming progress.
- **Persist Memory:** MANDATORY for long-running workflows. Explicitly use the `save_checkpoint` tool from `@mcp:context-manager` to snapshot the current task ID, completed steps, and next pending actions into the workspace database.

---

### 🛑 Socratic Gate (Optimized)

When asking questions to clarify vague requirements or complex features:

- **Format as Multiple Choice:** Minimize the user's typing effort.

| Request Type       | Action                                       |
| ------------------ | -------------------------------------------- |
| **New Feature**    | ASK 3+ strategic questions (Multiple Choice) |
| **Bug Fix**        | Confirm understanding + impact questions     |
| **Vague**          | Ask Purpose, Users, Scope (Multiple Choice)  |
| **Orchestration**  | STOP until user confirms plan                |
| **Direct Proceed** | Ask 2 Edge Case questions first              |

**Protocol:** Never assume → Spec-heavy? Ask trade-offs via options → Wait for Gate clearance.
