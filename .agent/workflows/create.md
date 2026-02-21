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

Analyze the request to identify the creation type. Use the `search_skills` tool from the `@mcp:skill-router` server to load required skills based on semantic understanding.

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

1. Read relevant `SKILL.md` files downloaded from the router.
2. If using new libraries/latest versions, use `search_latest_syntax` from `@mcp:doc-researcher` to ensure you are not writing legacy code.
3. Build the component.
4. Verify it works before moving on.
5. For complex tasks, use `save_checkpoint` from `@mcp:context-manager` to persist your progress and active files.

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
