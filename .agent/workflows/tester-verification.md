---
description: Structured workflow for the Test Agent. Handles test generation, execution, and verification of code from Coder.
---

# 🧪 Tester / Verification Workflow

This workflow guides the **Test Agent** through verifying the functionality and stability of code changes. It focuses on comprehensive coverage, edge-case validation, and ensuring the project remains bug-free.

## 🚀 Verification Phase

### Phase 1: Intake & Setup 📥
Receive signals from the Review Agent and prepare the environment.
- Use `@mcp:mcp-multi-agent` (`read_messages`) to confirm the code is ready for testing and understand the architectural intent.
- Ensure the codebase is in a testable state.

### Phase 2: Testing Strategy & Skill Discovery 🔍
After receiving the code changes and intent, identify the best testing approach.
- Use `@mcp:skill-router` (`search_skills`) to find relevant testing frameworks, patterns, and tools for the tech stack involved in the diff.
- Review `SKILL.md` files related to QA and E2E testing to plan a comprehensive verification suite.

### Phase 3: Test Engineering 🧪
Create comprehensive test cases to challenge the new implementation.
- **Unit Tests**: Write surgical tests for new functions/logic using the project's framework (Go test, Jest, Pytest, etc.).
- **Integration Tests**: Verify that new components interact correctly with existing services (e.g., Database, Redis).
- **Edge Cases**: Specifically target boundary conditions and error handling paths mentioned in the `PLAN`.
- Use `write_to_file` to create/update test files.

### Phase 4: Execution & Analysis ▶️
Run the test suite and evaluate the results.
- Execute tests using `run_command` (e.g., `npm test`, `go test ./...`).
- **If tests fail**:
    - Analyze the root cause (is it the code or the test?).
    - Use `@mcp:mcp-multi-agent` (`publish_message`) to notify the `target_role="coder"` with the failure logs and fix recommendations.
- **If tests pass**:
    - Verify code coverage if required.
    - Proceed to handover.

### Phase 5: Feedback & Signaling 🤝
Report the final status of the task.
- Use `@mcp:mcp-multi-agent` (`publish_message`) with `target_role="planner"` to signal that the task has been verified and is bug-free.
- Provide a brief summary of the test coverage and any performance observations.
- **WAIT**: Remain active and wait for the next test assignment from the Planner or Reviewer.

## 🔴 Critical Constraints
1. **Exclusive Test Ownership**: You are the ONLY agent allowed to write or modify files in the `tests/` directory.
2. **No Logic Changes**: Do NOT fix implementation code. Notify the Coder if tests fail.
3. **Automated Verification**: Always run the commands; never assume code works because it "looks right".

---

> [!IMPORTANT]
> If a test fails, provide the EXACT error logs to the Coder via internal message to expedite the fix loop.
