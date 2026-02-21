---
description: Debugging command. Activates DEBUG mode for systematic problem investigation.
---

# /debug - Systematic Problem Investigation

Guide agents to debug issues methodically: understand → hypothesize → investigate → fix → prevent.

## When to Use

- `/debug [issue description]`
- Keywords: "bug", "error", "not working", "broken", "fails", "crash", "issue"

## 🔴 Critical Rules

1. **Reproduce first** - Confirm the issue before fixing.
2. **One change at a time** - Don't change multiple things simultaneously.
3. **Verify fix works** - Test after each change.
4. **Understand root cause** - Don't just patch symptoms.

---

## Phase 1: Classification & Skill Mapping 🔀

Identify the core domain of the bug (Frontend, Backend, Database, Infrastructure, etc.) and map relevant skills from
`.agent/CATALOG.md`.

---

## Phase 2: Information Gathering 📋 (Socratic Gate)

Gather all available facts:

- **Expected vs Actual Behavior**
- **Error Messages & Stack Traces**
- **Logs & Recent Commits**
- **Reproduction Steps**

If anything is unclear, **ask clarifying questions** before forming hypotheses.

---

## Phase 3: Hypothesis Formation 🧠

List possible causes, ordered by likelihood (High, Medium, Low). Prioritize testing by: Most likely -> Easiest to verify
-> Most impactful.

---

## Phase 4: Systematic Investigation 🔍

Test each hypothesis one by one (using logs, breakpoints, minimal reproduction, etc.). Track your findings.

---

## Phase 5: Root Cause Analysis 🎯

State clearly:

- **What:** Technical description of the bug.
- **Why:** Why this caused the symptom.
- **Where:** File, line, function.

---

## Phase 6: Fix Implementation 🔧

Before coding, propose a **Fix Plan** detailing files to modify and risk levels. Apply the fix using loaded skill
patterns.

---

## Phase 7: Prevention & Documentation 🛡️

Ensure the fix is verified (tests pass, bug is gone). Add prevention measures (tests, validations, documentation).

Save the report to `debug-{slug}.md`.

### Output Template

```markdown
## 🔍 Debug Report: [Issue Title]

### 1. Problem

[Symptom vs Expected]

### 2. Root Cause

[Explanation]

### 3. Fix

[Diff or explanation of fix]

### 4. Prevention

[Measures added]
```
