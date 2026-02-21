---
description: Create project plan based on user requirements. No code writing - only planning.
---

# /plan - Project Planning Workflow

Guide the AI to create comprehensive, actionable plans based on required skills. **NO CODE WRITING**.

## When to Use

- `/plan [description]`
- Keywords: "plan", "roadmap", "breakdown", "strategy"

## 🔴 Critical Rules

1. **NO CODE WRITING** - This command creates plan files only.
2. **Socratic Gate** - Ask clarifying questions before planning.
3. **Actionable Output** - Every task must be executable.

---

## Phase 1: Context Analysis 🔍 (Socratic Gate)

Extract What, Why, Scope, Constraints, and Context from the request.

If vague, **ASK** before planning:

1. What are the boundaries?
2. What is the most important outcome?
3. Any constraints (team, tech, budget)?

Confirm understanding and proceed only when clear.

---

## Phase 2: Classification & Skill Mapping 🔀

Identify logical phases and use the `search_skills` tool from the `@mcp:skill-router` to suggest which skills are required for each phase.

---

## Phase 3: Task Breakdown & Dependencies 📋

Break work into phases (e.g., Foundation, Core Logic, Integration, Testing).

Create a **Dependency Graph** or list indicating what must be done first, and identify the Critical Path.

---

## Phase 4: Risk Assessment ⚠️

Document assumed risks and mitigations: | Risk | Impact | Mitigation | | ---- | ------ | ---------- | | ... | ... | ...
|

---

## Phase 5: Success Criteria & Delivery 📝

Define what "Done" looks like.

Save the output to `PLAN-{slug}.md` containing:

- Overview & Goals
- Scope & Task Breakdown
- Dependencies & Risks
- Success Criteria

### Request User Review

> ✅ **Plan saved:** `plan-{slug}.md` **Please review the plan and:**
>
> 1. ✅ Approve (Run `/create` to start)
> 2. 📝 Edit
> 3. 💬 Discuss
