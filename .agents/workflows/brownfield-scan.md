---
description: Structured workflow for analyzing and orchestrating legacy or existing codebases before refactoring.
---

# 🔎 Brownfield Reconnaissance Workflow

Use this workflow when the USER asks to "refactor," "scan," or "modernize" an existing codebase. **DO NOT start coding
or initialize a task plan immediately.** You must perform a structured analysis to prevent destructive changes and
contextual drift.

## 🚀 Execution Sequence

### Phase 1: Analyze 🧬

Scan the environment to detect frameworks, structure, and potential risks.

- Use `@mcp:ast-explorer` (`get_project_architecture`, `search_symbol`) to build a mental map of the project
  architecture and locate key definitions.
- Use `find_by_name` and `list_dir` to find configurations (`package.json`, `go.mod`, `pom.xml`, etc.).
- Identify "**Stale Context**" or "Code smells" (e.g., duplicated logic, deprecated library usage, mixed formatting).

### Phase 2: Recommend 💡

Formulate a remediation strategy based on your analysis.

- Summarize the identified architecture.
- List out 2-3 step-by-step strategic approaches to refactor or modernize the project without breaking existing
  operations (e.g., Strangler Fig pattern).
- Present this recommendation to the USER as a Socratic Gate, asking which strategy they prefer before generating the
  actual code.

### Phase 3: Orchestrate 🏗️

Isolate baseline constraints before beginning execution.

- Extract any invariant facts discovered (e.g., "Must communicate over gRPC", "Framework X version Y is strictly used").
- Use `@mcp:context-manager` (`manage_anchors` with `action="set"`) to persist these findings. This ensures they survive
  context compaction.
- Once the USER approves, transition to the `[Role: 🏗️ Planner]` to invoke `initialize_task_plan` and proceed with the
  execution.

## 🔴 Critical Constraints

1. **No Destructive Action**: Never delete or rewrite large chunks of legacy code during the Analyze phase.
2. **Anchor the Facts**: Ensure you leave `ANCHORS.md` updated so future `Coder` loops do not violate the baseline
   logic.
3. **Wait for Approval**: Do NOT proceed to Orchestration / Planning until the USER has explicitly confirmed the
   Recommendations from Phase 2.

---

## 📌 Usage Example

`/brownfield-scan "Please scan the repository and modernize the API layer"`
