---
description: Create new application command. Triggers interactive dialogue and skill mapping to build new projects systematically.
---

# /create - Universal Creation Workflow

Guide the AI to create new things systematically: services, systems, apps, components, pipelines.

## When to Use

- `/create [description]`
- Keywords: "create", "build", "make", "new", "implement", "set up"

---

## Phase 1: Classification & Skill Mapping 🔀

Analyze the request to identify the creation type. Check `.agent/CATALOG.md` to load required skills.

```markdown
🔀 **Creation Type:** [e.g., Web App, API, Model] 🛠️ **Primary Skills:** `@[skill-1]`, `@[skill-2]`
```

---

## Phase 2: Discovery 🔍 (Socratic Gate)

**🛑 Do NOT start building without understanding requirements.** Extract What, Why, Who, Scope, Constraints.

If unclear, **Ask 3+ strategic questions**:

- Must-have features?
- Target users?
- Tech preference?
- Expected scale/integrations?

Confirm understanding before planning.

---

## Phase 3: Planning 📋

Outline key architecture decisions. Provide a **Proposed File Structure**. Define a clear **Build Order** (e.g.,
Foundation → Data Layer → Business Logic → UI → Tests).

---

## Phase 4: Implementation 🔨

Follow the build order layer by layer. For each:

1. Read relevant `SKILL.md` files.
2. Build the component.
3. Verify it works before moving on.

Perform continuous Quality Checks (e.g., adherence to `clean-code`).

---

## Phase 5: Verification ✅

Run smoke tests (e.g., `npm run dev` or `go run`). Verify core functionality works without console errors or build
crashes.

---

## Phase 6: Documentation & Delivery 📝

Save a summary to `create-{slug}.md` detailing:

- What was created
- Tech Stack
- How to run
- Recommendations for next steps

> `✅ Creation complete! Summary saved: CREATE-{slug}.md`
