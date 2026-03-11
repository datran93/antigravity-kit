---
description: Structured workflow for Code Implementation. Handles task execution directly.
---

# 💻 Coder / Implementation Workflow

This workflow guides you through implementing atomic tasks defined during the planning phase. It emphasizes technical
excellence, clean code, and strict adherence to the project's architectural standards.

## 🚀 Implementation Phase

### Phase 1: Task Intake & Intent Locking 📥

Load the context for the current implementation step.

- Analyze the task requirements, the defined **Verification Command**, and specific files to modify.
- **Intent Locking**: Call `@mcp:context-manager` (`check_intent_lock`) on the files you intend to modify. If the tool
  returns a Scope Creep ALARM, you MUST transition back to the `[Role: 🏗️ Planner]` to update the intent via
  `declare_intent` before proceeding.
- Review any feedback from previous review or testing cycles, and check for any **Ghost Context** to understand the
  current state.

### Phase 2: Skill & Pattern Alignment 🔍

Align the implementation with project-specific standards via parallel discovery.

- **Parallel Context Gathering**: Use multiple MCP tools in a single turn to gather knowledge instantly:
  - Use `@mcp:skill-router` (`search_skills`) to find relevant clean code patterns.
  - Use `@mcp:context7` (`query-docs`) to verify the latest API syntax, preventing legacy code hallucinations.
- Review specific `SKILL.md` files referenced in the instructions to ensure architectural compliance.

### Phase 3: Execution & Engineering 🛠️

Write high-quality, maintainable code.

- **NO BLIND WRITES**: You MUST explicitly read a file (`view_file`, `grep_search`, or `ast-explorer`) before attempting
  to modify it (`replace_file_content`, `write_to_file`). Attempting to edit a file without prior reading is strictly
  forbidden.
- Implement the changes using `replace_file_content` or `write_to_file`.
- **MANDATORY**: Follow **Clean Code** principles (naming, small functions, SRP).
- **MANDATORY**: Ensure code is **Testable** (Dependency Injection, modularity).
- Focus only on the atomic task assigned. Do not over-engineer or touch unrelated files.

### Phase 4: Verification & Drift Detection 📝

Execute the Acceptance Criteria and handle failures safely.

- Run the **Verification Command** defined by the Planner (e.g., `npm run test:models` or `go build`). Do not rely
  solely on visual inspection.
- **Panic Protocol (Drift Detection)**: If the Verification Command fails and you cycle through "fix -> verify -> fail"
  **3 times**, you MUST STOP. Call `@mcp:context-manager` (`record_failure`) and transition back to the
  `[Role: 🏗️ Planner]` to discuss the blocker with the USER. Do not perform endless blind fixes.
- Synthesize a concise mental summary of your work:
  - Files modified (Keep track of this list, it will be needed as `active_files` when calling `complete_task_step`
    later).
  - Logic implemented/refactored.
  - Any technical debt or edge cases identified.
- If you realize the current task is too large or requires new steps, transfer to the `planner` role and use
  `@mcp:context-manager` (`add_task_step`) to append them.

### Phase 5: Role Transition & Export Intelligence 🔄

Hand over your context cleanly to the next role.

- **Inject Ghost Context**: If you had to implement a complicated workaround or discovered a framework gotcha during
  coding, invoke `@mcp:context-manager` (`annotate_file`) to attach this lesson directly to the file.
- Before transitioning, extract key "intelligence" (e.g., "Library X had a bug, so I downgraded it to Y", or "I had to
  adjust the DB schema for Z").
- Pass this intelligence explicitly in your conversational response so the next role (e.g., Reviewer or Tester) doesn't
  start blind.
- Transition to the `reviewer` or `tester` role to validate your implementation.

## 🔴 Critical Constraints

1. **Quality First**: Never ignore lint errors or violations of the design system.
2. **Direct Feedback**: If you encounter issues that prevent completion, report them to the user and adjust the plan.
3. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 💻 Coder]`.

---

> [!TIP] Always verify that your changes don't break existing modularity. If you needed to refactor common code,
> document it clearly for the next testing phase.
