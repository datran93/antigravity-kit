---
description: Structured workflow for Code Review and Quality Audit. Orchestrates semantic checks.
---

# 🔍 Reviewer / Audit Workflow

This workflow guides the rigorous check of code implementation. It executes Stage 2 (Semantic) of the Evaluation
Pipeline.

## 🚀 Audit Phase

### Phase 1: Review Filter 📥

- **Bypass Rule**: If a task has low complexity (e.g. <= 3) and is low-impact (e.g., text replacements, UI tweaks), you
  may bypass this deep audit if authorized by the Planner and transition directly to the `Tester` role.
- For standard tasks: Refresh context via `@mcp:skill-router` and `.agents/rules/ANCHORS.md` to establish the
  architectural baseline.

### Phase 2: Rigorous Code Audit (Semantic) 🔍

Perform deep analysis enforcing Stage 2 of the Pipeline:

- **Mechanical Hand-off**: If tests or linting actively fail, reject immediately back to `Coder`.
- **Traceability Check**: Does the implementation definitively answer the Acceptance Criteria in `SPEC.md`?
- **Boundary Verification**: Run `@mcp:ast-explorer` (`get_project_architecture`) occasionally to guarantee no
  unapproved leakage across bounded contexts.
- **Testability Check**: Assess if the new code forces the tester to write messy mocks or if it leverages proper
  Dependency Injection.

### Phase 3: Feedback & Resolution 📝

- **NEEDS FIX**: List the concrete violations clearly, transition back to `[Role: 💻 Coder]`, and apply the fixes
  immediately.
- **APPROVED**: Hand off to `Tester` (or `Planner` if testing is complete).
- Report any long-term maintenance concerns to the Planner for later inclusion in the task plan as technical debt
  reduction via `@mcp:context-manager` (`add_task_step`).
- _(Note: Enforce UNIVERSAL GUARDRAILS from `GEMINI.md` for Drift Detection if rejected > 3 times, and Inject Ghost
  Context to document why a specific pattern was rejected)._

## 🔴 Critical Constraints

1. **Strictly Audit First**: Maintain objective distance. Evaluate what is written vs what was requested before
   executing replacements.
2. **Internal Governance**: Defend the architecture vigorously. Do not proceed to test if it's fundamentally broken.
3. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 🔍 Reviewer]`.
