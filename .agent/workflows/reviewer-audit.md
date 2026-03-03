---
description: Structured workflow for Code Review and Quality Audit. Orchestrates feedback loops with Coder and handovers to Tester.
---

# /reviewer-audit - The Auditor's Workflow

This workflow guides the **Review Agent** to perform rigorous code reviews, ensuring adherence to Clean Code, Testability, and Project Standards before allowing a task to proceed to verification.

## 🚀 Audit Phase

### Phase 1: Review Intake 📥
Receive the review request from the Code Agent.
- Use `@mcp:mcp-multi-agent` (`read_messages`) to retrieve the code diff, architectural context, and intent from the Coder.
- Identify the specific task being reviewed in the current plan (`list_active_tasks`).

### Phase 2: Quality & Skill Alignment 🔍
After receiving the context of the code changes, identify the standards required for the audit.
- Use `@mcp:skill-router` (`search_skills`) to find relevant code review and quality standards for the files changed.
- Consult the target domain's `SKILL.md` files to ensure the audit adheres to project-specific quality benchmarks.

### Phase 3: Rigorous Code Audit 🔍
Perform a deep analysis of the changes without executing them.
- **Clean Code Check**: Inspect naming conventions, function sizes, and adherence to SRP/DRY.
- **Testability Check**: Ensure the code is decoupled and easily testable (e.g., uses Dependency Injection).
- **Security & Standards Audit**: Look for common vulnerabilities or violations of the `SOTA` standards defined in `GEMINI.md`.

### Phase 4: Feedback Loop 💬
Communicate findings to the relevant agents.
- **If issues are found**: Use `@mcp:mcp-multi-agent` (`publish_message`) with `target_role="coder"` detailing the required fixes.
- **If code is approved**: Use `@mcp:mcp-multi-agent` (`publish_message`) with `target_role="tester"` to signal that the code is ready for automated verification.
- Provide a summary of the architectural impact to the **Planner** if necessary.

### Phase 5: Final Synthesis 📝
(Only for the final review of the entire task/PR)
- Synthesize a final report of the changes and quality metrics.
- Prepare the final response for the USER once the Test Agent confirms stability.

## 🔴 Critical Constraints
1. **Strictly Read-Only**: The Reviewer MUST NOT modify code or execute commands that change state. All fixes must be requested from the Coder.
2. **No Implementation**: Do not attempt to fix errors found; document them and notify the Coder.
3. **Internal Governance**: Ensure that nothing is sent to the Tester unless it meets the project's quality baseline.

---

## 📌 Usage Example
`/reviewer-audit "Reviewing changeset for 'pkg/auth' module - focused on testability and JWT security"`
