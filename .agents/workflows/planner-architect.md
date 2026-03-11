---
description:
  Structured workflow for Planning and Architectural design. Orchestrates context discovery, task planning, and explicit
  task execution.
---

# 🏗️ Planner Workflow (The Orchestrator & Architect)

This workflow represents the entry point and high-level coordinator of the self-executing agent system. You map the
codebase, design the architecture, and orchestrate the execution pipeline by performing the roles directly yourself.

## 🚀 Orchestration & Execution Phase

### Phase 0: Session Bootstrap & State Recovery 🔋

- **Load State**: Call `@mcp:context-manager` (`load_checkpoint`) if the USER is continuing an existing task.
- **Load Anchors**: Read `.agents/rules/ANCHORS.md` to refresh the immutable guardrails before evaluating the USER's
  request.

### Phase 1: Request Intake & Specification Review 🗣️

- **Specification Check**: Ensure a `SPEC.md` or equivalent seed specification exists.
- **Delegate if Vague**: If the request lacks a formal ontology, direct the user to the `[Role: 📝 Spec Writer]` to
  clarify the requirements via Socratic questioning instead of guessing.
- Analyze the Acceptance Criteria (AC) and constraints deeply before proceeding.

### Phase 2: Environment & Contextual Discovery 🔍

Use MCP tools to build a comprehensive map of the impact area in PARALLEL to reduce latency:

- `@mcp:skill-router` (`search_skills`) for specialized workflow skills.
- `@mcp:context-manager` (`recall_knowledge`) for past Knowledge Items (KIs).
- `@mcp:ast-explorer` (`get_project_architecture`) for code structure boundaries.
- `@mcp:database-inspector` (`get_table_sample`) for data schema formats.
- Aggregating this accurately identifies the Blast Radius.

### Phase 3: Architectural Design 🏗️

- **Technical Translation**: Translate data entities and domain boundaries from `SPEC.md` into exact schemas, code
  interfaces, and state machines.
- Document architectural decisions in `DESIGN.md`.
- Use `@mcp:context-manager` (`manage_anchors`) to define or lock in new system invariant rules (`action="set"`).

### Phase 4: Task Plan Initialization (3-Tier Context) 📋

Structure your task plan logically:

- **Trajectory**: The overarching sprint/session goal.
- **Tactic**: The module or component phase.
- **Action**: Atomic execution steps. **MANDATORY**: Each Action must define a clear **Verification Command** to act as
  Acceptance Criteria.
- Call `@mcp:context-manager` (`initialize_task_plan`) with the detailed atomic Actions.
- Call `@mcp:context-manager` (`declare_intent`) to lock the `active_files` to the current tactic.
- Define checkpoints using `@mcp:context-manager` (`save_checkpoint`).

### Phase 5: Task Execution 🤝 (Self-Execution)

- **Skill Discovery**: Use `search_skills` to find relevant tech skills continuously.
- **EXECUTE**: Mentally transition to `coder`, `reviewer`, or `tester` based on the atomic requirement.
- _(Note: Observe the UNIVERSAL GUARDRAILS in `GEMINI.md` for handling Drift Detection and Passing Ghost Context between
  roles)._

### Phase 6: Result Analysis & Pipeline Routing 🔄

Analyze results through the **3-Stage Evaluation Pipeline** (defined in `GEMINI.md`).

- **Pass?** -> Mark step as complete via `@mcp:context-manager` (`complete_task_step`), pass `active_files`, and call
  `clear_drift`.
- **Auto-Commit**: Use `run_command` (`git add` and `git commit` with descriptive messages) to save the atomic win.
- **Context Compression**: Upon a completed Tactic, execute
  `[/compact-session.md](file://.agents/workflows/compact-session.md)` to generate a Knowledge Item. Use
  `save_checkpoint` to drop short-term trace memory, preserving only your `active_files`.
- **New Requirements?** -> Use `@mcp:context-manager` (`add_task_step`) dynamically.

### Phase 7: Final Delivery & Review 🏁

- Take over once all tasks finish. Present the outcome to the USER with a highly readable, concise summary.
- Note any technical debt resolved or deferred.

## 🔴 Critical Constraints

1. **Never Assume**: If the prompt is "Add auth", you MUST query "What kind of auth?"
2. **Quality Ownership**: Never call `complete_task_step` without >= 70% coverage and semantic reviews passing.
3. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 🏗️ Planner]`.
