---
description: Structured brainstorming for projects and features. Explores multiple options - NO implementation.
---

# /brainstorm - High-Leverage Idea Exploration

Guide the user through a structured exploration of ideas: Context → Divergence → Evaluation → Selection.

🚨 **CRITICAL RULE: This workflow generates OPTIONS ONLY. No implementation details, code snippets, or task plans.**

---

## Phase 1: Context & Skill Discovery 🔍

1. **Skill Mapping**: Use `@mcp:skill-router` (`search_skills`) to find domain-specific knowledge.
2. **Context Intake**: Use the **Socratic Gate** (Multiple Choice) to clarify:
   - **Problem**: What is the root cause?
   - **Constraints**: Budget, tech stack, or deadline limitations?
   - **Success**: What does a "perfect" solution look like?

---

## Phase 2: Divergent Options 🌐

Generate exactly **3 distinct options** with different trade-offs (e.g., Fast/Cheap, Scalable/Premium, Unconventional/Risky).

For each option, provide:
- **Concept**: High-level approach.
- **Mental Model**: How to think about this solution.
- **Why it fits**: Alignment with user constraints.

---

## Phase 3: Convergent Analysis 📊

Evaluate the options using a comparison matrix:

| Metric     | Option 1     | Option 2 | Option 3 |
| :--------- | :----------- | :------- | :------- |
| **Effort** | Low/Med/High | ...      | ...      |
| **Impact** | ...          | ...      | ...      |
| **Risk**   | ...          | ...      | ...      |

Highlight the **Non-Obvious Trade-offs** for each.

---

## Phase 4: Recommendation & Decision 💡

1. **Synthesize**: Recommend one option based on the weighed trade-offs.
2. **STOP**: Wait for the user to choose or iterate. **DO NOT provide code or implementation steps.**

---

## Phase 5: Decision Persistence 📝 (Post-Selection)

**🔔 Execute ONLY after a choice is made.**

Save the decision to `brainstorm-{slug}.md` for long-term memory.

### Template:
```markdown
# Decision: [Option Name]
**Date**: [YYYY-MM-DD]
**Rationale**: Why this won over other options.
**Trade-offs**: Explicitly accepted drawbacks.
```
