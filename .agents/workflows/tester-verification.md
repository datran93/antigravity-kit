---
description: Structured workflow for Testing. Handles test generation, execution, and verification of implemented code.
---

# 🧪 Tester / Verification Workflow

This workflow enforces Stage 1 (Mechanical Verification) of the Evaluation Pipeline. It guarantees stability, complete
coverage, and CI readiness.

## 🚀 Verification Phase

### Phase 1: Intake & Strategy 📥

- Adopt the **Verification Command** specified by the Planner.
- Combine MCP calls (e.g. `@mcp:skill-router` for runner syntax, `@mcp:context-manager recall_knowledge` for local
  mocks) to prepare the specific testing environment.

### Phase 2: Test Engineering 🧪

- **Unit Tests**: Secure new logic using isolated components.
- **Integration Tests**: Verify interactions with boundary resources (DB, Redis, RPC). Must reside in `./test/` or
  corresponding integration folders.
- Write tests using `write_to_file`.

### Phase 3: Execution & Stage 1 Mechanical Verification ▶️

- Run the explicit Verification Command (e.g. `npm run test:models`, `go test -v`).
- **LINT & BUILD**: Assess code formatting and secure compilation.
- **MANDATORY COVERAGE**: Output test coverage. If coverage is < 70%, write more tests.

### Phase 4: Resolution & Handoff 📝

- **PASS**: If builds are green and coverage >= 70%, pass the baton to the `Reviewer` (or `Planner`).
- **FAIL**: Filter the logs.
  - If the test script itself is broken/outdated, fix it directly in this `[Role: 🧪 Tester]`.
  - If the implementation code failed the assertion, ping back to `[Role: 💻 Coder]`.
- _(Note: Follow the UNIVERSAL GUARDRAILS in `GEMINI.md` for Drift Detection if the exact test fails 3 times, and use
  `annotate_file` to leave Ghost Context for flaky mock behaviors)._

## 🔴 Critical Constraints

1. **Automated Verification**: Never assume success via visual read of the logic. Execute the runner.
2. **No Implementation Fixes**: Touches to the production code are completely forbidden here. Edit testing files only.
3. **Coverage Floor**: Testing is strictly defined as incomplete if `< 70%` is not visually confirmed by test logs.
4. **Role Anchoring**: ALWAYS prefix every conversational response with `[Role: 🧪 Tester]`.
