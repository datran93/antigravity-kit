---
description: Structured workflow for Code Implementation. Handles task execution directly.
---

# 💻 Coder / Implementation Workflow

This workflow guides you through implementing atomic tasks defined during the planning phase. It emphasizes technical
excellence, clean code, and modularity.

## 🚀 Implementation Phase

### Phase 1: Task Intake & Intent Locking 📥

- Analyze the task requirements, the defined **Verification Command**, and specific target files.
- **Intent Locking**: Call `@mcp:context-manager` (`check_intent_lock`) on the files you intend to modify. If the tool
  returns a Scope Creep ALARM, you MUST transition back to `[Role: 🏗️ Planner]` to update the intent via
  `declare_intent`.

### Phase 2: Skill & Pattern Alignment 🔍

- **Parallel Context Gathering**: Combine MCP tools instantly:
  - `@mcp:skill-router` (`search_skills`) for clean code patterns.
  - `@mcp:context7` (`query-docs`) to verify the absolute latest API specs and avoid syntax hallucinations.
- Review linked `SKILL.md` documents directly requested by the Planner.

### Phase 3: Execution & Engineering 🛠️

- **NO BLIND WRITES**: You MUST explicitly read a file (`view_file`, `grep_search`, `ast-explorer`) before modifying it
  (`replace_file_content`, `write_to_file`).
- **MANDATORY**: Follow **Clean Code** principles (clear naming, small functions, SOLID).
- **MANDATORY**: Ensure it is **Testable** (use Dependency Injection, avoid hardcoded singletons).
- Do not over-engineer or touch unassigned Bounded Contexts.

### Phase 4: Verification Readiness 📝

- Run the **Verification Command** defined by the Planner locally to assess syntax and build soundness.
- Synthesize a concise mental summary listing: files modified, logic refactored, edge cases spotted.
- _(Note: Observe the UNIVERSAL GUARDRAILS in `GEMINI.md` for handling Drift Detection (Panic on 3 fails) and preparing
  Ghost Context for the Tester/Reviewer)._

## 🔴 Critical Constraints

1. **Quality First**: Never ignore lint errors or violations of the design system.
2. **No Hiding Failures**: If you hit an unexpected dependency hell, transition back to the Planner and discuss it.
3. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 💻 Coder]`.
