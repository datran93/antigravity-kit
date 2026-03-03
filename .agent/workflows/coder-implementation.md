---
description: Structured workflow for Code Implementation. Handles task execution, snippet verification, and Reviewer feedback loops.
---

# 💻 Coder / Implementation Workflow

This workflow guides the **Code Agent** through implementing atomic tasks defined by the Planner. It emphasizes technical excellence, clean code, and strict adherence to the project's architectural standards.

## 🚀 Implementation Phase

### Phase 1: Task Intake 📥
Load the context provided by the Planner and the task plan.
- Use `@mcp:context-manager` (`list_active_tasks`) to see the current plan.
- Identify the next pending task.
- Read internal messages from the Planner using `@mcp:mcp-multi-agent` (`read_messages`) for specific implementation hints.

### Phase 2: Skill & Pattern Alignment 🔍
After receiving the task context, align the implementation with project-specific standards.
- Use `@mcp:skill-router` (`search_skills`) to find relevant clean code patterns and language-specific best practices for the task.
- Review the `SKILL.md` files referenced in the plan to ensure architectural compliance.

### Phase 3: Execution & Engineering 🛠️
Write high-quality, maintainable code.
- Implement the changes using `replace_file_content` or `write_to_file`.
- **MANDATORY**: Follow **Clean Code** principles (naming, small functions, SRP).
- **MANDATORY**: Ensure code is **Testable** (Dependency Injection, modularity).
- Use `@mcp:context-manager` (`save_checkpoint`) if the task is complex.

### Phase 4: Reviewer Feedback Loop 💬
Before marking the task as complete, seek architectural and quality confirmation.
- Use `@mcp:mcp-multi-agent` (`publish_message`) to send the code diff or a summary to the `target_role="reviewer"`.
- Focus on review for code quality, adherence to standards, and testability.
- Wait for feedback (`read_messages`).
- **If negative feedback**: Fix the issues and re-publish.
- **If positive feedback**: Proceed.

### Phase 5: Handover & Wait 🤝
Finalize the implementation snippet.
- Notify the **Review Agent** via internal message that the code is ready for audit.
- **DO NOT** call `complete_task_step`. The Planner will manage the task lifecycle.
- **WAIT**: Remain active and wait for feedback from the Reviewer or new instructions from the Planner.

## 🔴 Critical Constraints
1. **Exclusive Code Ownership**: During implementation, you are the ONLY agent allowed to modify the source code.
2. **Quality First**: Never ignore lint errors or violations of the design system.
3. **No Completion Ownership**: The Coder builds; the Planner manages the plan.

---

> [!TIP]
> Always verify that your changes don't break existing modularity. If you need to refactor common code, document it for the Reviewer.
