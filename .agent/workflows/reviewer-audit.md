---
description: Structured workflow for Code Review and Quality Audit. Orchestrates feedback loops to the Coder role.
---

# 🔍 Reviewer / Audit Workflow

This workflow guides you to perform rigorous code reviews on your own implementations, ensuring adherence to Clean Code,
Testability, and Project Standards before advancing to testing.

## 🚀 Audit Phase

### Phase 1: Review Intake & Filtering 📥

Evaluate if a distinct review phase is necessary.

- **Bypass Rule**: If the task has low complexity (e.g. <= 3) and is low-impact (e.g., minor UI/Text changes), you may
  bypass this distinct audit and transition directly to the `tester` role.
- **For High-Impact / High-Complexity Tasks**: Load the context of the recent code changes. Identify the target quality
  standards using `@mcp:skill-router` (`search_skills`), re-read `.agent/rules/ANCHORS.md`, and understand the
  **Verification Command** to ensure the code's intended behavior meets Planner's acceptance criteria.

### Phase 2: Rigorous Code Audit (3-Stage Evaluation Pipeline) 🔍

Perform a deep analysis of the changes using parallel MCP calls, enforcing the 3-Stage Evaluation Pipeline.

- **Stage 1 (Mechanical / Tester Hand-off)**: Acknowledge if mechanical checks (LINT, BUILD, TEST) have successfully passed. If they haven't been run or fail, reject immediately and punt back to `coder`.
- **Stage 2 (Semantic Evaluation)**: Evaluate whether the output strictly satisfies the "Seed Specification" Acceptance Criteria.
  - Trace the functionality back to the original architecture requirements or plan.
  - Run `@mcp:ast-explorer` (`get_project_architecture`) to ensure the newly written code hasn't broken bounded contexts.
  - Review implementation against Clean Code principles (naming conventions, SRP/DRY).
  - **Testability Check**: Ensure the code is decoupled and easily testable (e.g., uses Dependency Injection).
- **Stage 3 (Consensus / Frontier Tier)**: For high-stakes decisions or highly ambiguous code, adopt a "Devil's Advocate" perspective. Challenge hidden assumptions before approving.
- **Security & Standards Audit**: Look for common vulnerabilities or violations of industry best practices.

### Phase 3: Feedback, Resolution & Drift Detection 📝

Act on the audit results and prevent endless revision loops.

- **NEEDS FIX**: If issues are found, list them clearly, transition back to the `coder` role, and fix them immediately.
- **Panic Protocol (Drift Detection)**: If you reject the code for the **same core issue 3 times** in a row, DO NOT send
  it back to the `coder`. Call `@mcp:context-manager` (`record_failure`) and transition back to the `[Role: 🏗️ Planner]`
  to re-evaluate the architecture or discuss with the USER.
- **APPROVED**: If code meets standards, proceed to the next testing or planning phase.
- Details any potential long-term maintenance concerns. If significant technical debt is found, transfer to `planner`
  role and add new refactoring tasks to the plan using `@mcp:context-manager` (`add_task_step`) but isn't a blocker
  right now.

### Phase 4: Role Transition & Export Intelligence 🔄

Hand over your context cleanly to the next role.

- **Inject Ghost Context**: If you detected a sneaky anti-pattern or a specific edge-case error during the review, call
  `@mcp:context-manager` (`annotate_file`) to log this lesson directly to the file so future Coders avoid repeating the
  same mistake.
- Before transitioning, extract key "intelligence" (e.g., "Refactored module X to comply with DRY", or "Noted technical
  debt Y for the future").
- Pass this intelligence explicitly to the next role via `@mcp:context-manager` (`save_checkpoint` notes) or your
  conversational response so the next role doesn't start blind.
- Successfully passing the review means transitioning to the `tester` role or continuing the plan.

## 🔴 Critical Constraints

1. **Strictly Audit First**: Evaluate the code objectively before trying to fix it.
2. **Internal Governance**: Do not proceed to testing if the code does not meet the basic quality baseline.
3. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 🔍 Reviewer]`.

---

> [!TIP] Focus your feedback on "Why" something should change, not just "What". Keep yourself accountable to the
> project's architectural standards.
