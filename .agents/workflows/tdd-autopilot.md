---
description: Structured workflow for Test-Driven Development (TDD) execution loops.
---

# 🤖 TDD Autopilot Workflow

Use this workflow to enforce strict Test-Driven Development loops. It shifts the pipeline from `spec -> build -> validate` to `spec -> validate (Red) -> build (Green) -> optimize (Refactor)`.

## 🚀 Execution Sequence

### Phase 1: Test Specification (Red) 🔴
Write the tests for the required feature first.
- As `[Role: 💻 Coder]`, understand the requirements and create/write the test cases (e.g., `feature_test.go`, `feature.spec.ts`) utilizing the exact contract the system must fulfill.
- Write tests that intentionally **fail** because the implementation does not yet exist.
- Transition to `[Role: 👨‍🔬 Tester]` and run the test suite to verify the failing state (Output must literally show failures or missing dependencies).
- Add failing logs to the context notes via `save_checkpoint`.

### Phase 2: Implementation Execution (Green) 🟢
Write the minimum exact code necessary to pass the test cases.
- As `[Role: 💻 Coder]`, build out the feature implementation (`feature.go`, `feature.ts`).
- **NO BLIND WRITES**: Write clean, concise code strictly focused on the requested atomic test cases.
- Transition back to `[Role: 👨‍🔬 Tester]` and verify that the tests are now passing 100%. If any test fails, cycle back to `Coder`.
- **GATE**: Verify test coverage is **>= 70%** for the new module. If it is lower, write more edge case tests.

### Phase 3: Code Optimization (Refactor) 🛠️
Ensure the written code is clean and adheres to architectural standards.
- After all tests pass, assume the `[Role: 🧐 Reviewer]`.
- Enforce strict `Clean Architecture` or `SOLID` checking. Check for Dependency Injection, magic strings, and duplicated logic.
- Switch back to `[Role: 💻 Coder]` to apply the cleanup.
- Final test suite evaluation. Must be 100% Green.

### Phase 4: Autopilot Transition 🔄
Export the completed scope and update context.
- Once the cycle completes, switch to the `[Role: 🏗️ Planner]`.
- Summarize the implemented API/logic into `save_checkpoint` as Intelligence context.
- Mark the current task complete (`complete_task_step`) and prepare or initiate the next atomic TDD component cycle.

## 🔴 Critical Constraints
1. **Never Reverse the Cycle**: The initial loop MUST always fail before any implementation files are created or modified.
2. **Strict Coverage**: 70% bounds act as a hard constraint.
3. **Patience**: This workflow requires strict cyclical discipline. Do not conflate step 1 with step 2 simultaneously.

---

## 📌 Usage Example
`/tdd-autopilot "Implement the JWT token extraction parser middleware"`
