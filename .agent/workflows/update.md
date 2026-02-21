---
description: Add or update features in existing application. Used for iterative development.
---

# /update - Update Existing Feature Workflow

Guide the AI to update, enhance, or modify existing features systematically with minimal risk.

## When to Use

- `/update [description]`
- Keywords: "update", "change", "modify", "enhance", "improve", "edit"

## 🔴 Critical Rules

1. **Understand before modifying** - Read existing code first!
2. **Impact analysis required** - Know what will be affected.
3. **Preserve working functionality** - Don't break what works.

---

## Phase 1: Classification & Context 🔍 (Socratic Gate)

Identify the feature being modified. Check `.agent/CATALOG.md` for necessary skills (e.g., framework patterns).

Read the existing codebase to understand Current State. **If request is vague, ask:** What exactly should change? Should
existing behavior be preserved?

---

## Phase 2: Impact Analysis & Risk ⚡

List affected areas (Direct, Dependencies, Downstream). Summarize files to modify, create, or delete. Assess Risk
(Low/Medium/High).

---

## Phase 3: Change Planning 📋

Break down changes into:

1. Preparation (Reading code)
2. Core Changes
3. Integration (Updating imports, deps)
4. Testing

For major changes, **present plan to user for approval** before executing.

---

## Phase 4: Implementation 🔨

Execute changes file by file based on loaded skills. Handle edge cases (Nulls, Errors, Empty states). Maintain backwards
compatibility where necessary (e.g., API signatures).

---

## Phase 5: Verification & Delivery ✅

Run smoke tests, linting, and unit tests to ensure nothing broke.

Save a summary to `update-{slug}.md`:

- What changed (Files)
- New Behavior
- Breaking Changes
- Testing Output

> `✅ Update complete! Report saved: UPDATE-{slug}.md`
