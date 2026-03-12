---
description:
  Structured workflow for Code Implementation. Reads DESIGN.md and the Planner's task list, executes each task in order,
  and reports a clear summary of all changes to the USER upon completion.
---

# 💻 Coder Workflow

This workflow guides the Coder through implementing each task defined by the Planner. It starts from reading the design
and task list, executes work atomically, and ends with a structured report to the USER.

---

## 🚀 Implementation Phases

### Phase 0: Read Design & Task List 📖

Before touching any code:

1. **Read `DESIGN.md`** — Understand the architecture, data models, file structure, and constraints the Planner defined.
2. **Load the task list** — Call `@mcp:context-manager` (`load_checkpoint`) to retrieve the active task plan and ordered
   Action steps.
3. **Confirm scope** — Identify which files will be created or modified, and in what order.

> ❌ Do NOT start writing code before completing this phase.

---

### Phase 1: Task Intake & Intent Lock 📥

For each Action in the task list (execute one at a time, in order):

- Read the Action's description, target files, and **Verification Command**.
- **Intent Locking**: Call `@mcp:context-manager` (`declare_intent`) to lock the files for this Action.
- Call `@mcp:context-manager` (`check_intent_lock`) before modifying any file. If a Scope Creep ALARM is returned,
  **stop and ask the USER** — do not proceed or expand scope independently.

---

### Phase 2: Skill & Pattern Alignment 🔍

Before implementing each Action:

- `@mcp:skill-router` (`search_skills`) — find relevant patterns or tech-specific best practices.
- `@mcp:context7` (`query-docs`) — verify the latest API specs for any library involved (avoid syntax hallucinations).
- Cross-reference `DESIGN.md` for contracts that must be respected.

---

### Phase 3: Execution 🛠️

- **NO BLIND WRITES**: Always read a file (`view_file`, `grep_search`, `ast-explorer`) before modifying it.
- Follow **Clean Code** principles: clear naming, small focused functions, SOLID.
- Ensure code is **testable**: use Dependency Injection, avoid hardcoded globals.
- Do not touch files or logic outside the locked Bounded Context for this Action.

---

### Phase 4: Verification ✅

After completing each Action:

- Run the **Verification Command** defined by the Planner (e.g., `go test ./...`, `npm run lint`).
- If the command **fails**:
  - Call `@mcp:context-manager` (`record_failure`).
  - Fix the issue and re-run. After **3 consecutive failures on the same issue**, stop and ask the USER.
- If the command **passes**:
  - Note the result — do NOT commit, do NOT mark the task as complete.
  - Committing and closing tasks is the **Planner's responsibility** after review and tests pass.

Repeat **Phase 1 → Phase 4** for each remaining Action in the task list.

---

### Phase 5: Final Report to USER 📋

When **all Actions** in the task list are complete, present a structured report:

```
## ✅ Implementation Complete

### Changes Made
| File | Action | Purpose |
|------|--------|---------|
| path/to/file.go | Created | Implements X service to handle Y |
| path/to/other.go | Modified | Added Z function to support W |
| ... | ... | ... |

### What was built
<1-2 sentence summary of the feature/fix implemented>

### Verification
- All Verification Commands passed: ✅
- Files committed: <list of commits>

### Notes / Known Limitations
<Any technical debt, deferred items, or edge cases to watch>
```

> 🛑 **STOP HERE.** After the report, the USER decides whether to proceed with `/tester-verification`,
> `/reviewer-audit`, or other workflows.

---

## 🔴 Critical Constraints

1. **Read before write**: Never modify a file without reading it first.
2. **Task order matters**: Complete Actions in the order defined by the Planner. Do not skip or reorder.
3. **No hidden failures**: If an unexpected blocker appears (dependency conflict, missing contract), stop and ask the
   USER immediately — do not make assumptions or attempt to resolve it by reinterpreting the design.
4. **Stay in scope**: Do not refactor, rename, or improve code outside the current Action's Bounded Context.
5. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 💻 Coder]`.
