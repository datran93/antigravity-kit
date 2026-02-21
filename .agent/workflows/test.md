---
description: Test generation and test running command. Creates and executes tests for code.
---

# /test - Test Generation and Execution

Guide the AI to generate, run, and analyze tests systematically.

## When to Use

- `/test` - Run all tests
- `/test [file/feature]` - Generate tests for specific target
- `/test coverage` - Show test coverage report

---

## Phase 1: Classification & Skill Mapping 🔀

Identify the existing testing framework (Jest, Pytest, Go Test, etc.) and load appropriate testing skills from
`.agent/CATALOG.md` (e.g., `testing-patterns`, `e2e-testing-patterns`).

---

## Phase 2: Test Analysis & Target 🔍 (Socratic Gate)

Analyze the target code (inputs, outputs, dependencies, edge cases). Identify test cases to create: Happy path, Edge
cases, Error cases. **Ask for clarification** if the behavior of the targeted code is ambiguous.

---

## Phase 3: Test Generation 🧪

Follow the AAA pattern (Arrange-Act-Assert):

- **Arrange:** Set up test data and mocks.
- **Act:** Execute the code under test.
- **Assert:** Verify the outcome.

Maintain one logical assertion per test where practical, and name tests descriptively.

---

## Phase 4: Test Execution & Debugging ▶️

Run the tests using the detected framework command.

Analyze results. If tests fail:

- Provide a **Failure Analysis** with Expected vs Actual.
- Suggest and apply fixes to code or tests.

---

## Phase 5: Coverage & Delivery 📝

Generate coverage report if applicable. Identify gaps.

Save the test report to `test-{slug}.md` summarizing Pass/Fail rates, Coverage, and Recommendations.

> `✅ Test report saved: test-{slug}.md`
