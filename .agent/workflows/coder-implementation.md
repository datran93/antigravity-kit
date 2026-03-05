---
description: Structured workflow for Code Implementation. Handles task execution directly.
---

# 💻 Coder / Implementation Workflow

This workflow guides you through implementing atomic tasks defined during the planning phase. It emphasizes technical excellence, clean code, and strict adherence to the project's architectural standards.

## 🚀 Implementation Phase

### Phase 1: Task Intake 📥
Load the context for the current implementation step.
- Analyze the task requirements and specific files to modify.
- Review any feedback from previous review or testing cycles to understand the current state.

### Phase 2: Skill & Pattern Alignment 🔍
Align the implementation with project-specific standards.
- Use `@mcp:skill-router` (`search_skills`) to find relevant clean code patterns and language-specific best practices for the task.
- Review specific `SKILL.md` files referenced in the instructions to ensure architectural compliance.

### Phase 3: Execution & Engineering 🛠️
Write high-quality, maintainable code.
- Implement the changes using `replace_file_content` or `write_to_file`.
- **MANDATORY**: Follow **Clean Code** principles (naming, small functions, SRP).
- **MANDATORY**: Ensure code is **Testable** (Dependency Injection, modularity).
- Focus only on the atomic task assigned. Do not over-engineer or touch unrelated files.

### Phase 4: Verification 📝
Prepare for review and testing.
- Run lightweight verification commands (e.g., `go build`, `tsc --noEmit`, `python -m py_compile`) to ensure no syntax errors.
- Synthesize a concise mental summary of your work:
    - Files modified.
    - Logic implemented/refactored.
    - Any technical debt or edge cases identified.

### Phase 5: Role Transition 🔄
- Transition to the `reviewer` or `tester` role to validate your implementation.

## 🔴 Critical Constraints
1. **Quality First**: Never ignore lint errors or violations of the design system.
2. **Direct Feedback**: If you encounter issues that prevent completion, report them to the user and adjust the plan.

---

> [!TIP]
> Always verify that your changes don't break existing modularity. If you needed to refactor common code, document it clearly for the next testing phase.
