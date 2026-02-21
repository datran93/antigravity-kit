---
description: Structured brainstorming for projects and features. Explores multiple options - NO implementation.
---

# /brainstorm - Structured Idea Exploration

Guide agents to brainstorm user ideas systematically: understand → explore → evaluate → **user decides → save**.

🚨 **CRITICAL RULE: This workflow generates OPTIONS ONLY. No implementation suggestions!**

## When to Use

- `/brainstorm [topic]`
- Keywords: "explore", "options", "ideas", "approach", "how should I", "what if"

---

## Phase 1: Classification & Skill Mapping 🔀

### Identify Relevant Skills

Analyze the topic and list related skills from `.agent/CATALOG.md`.

```markdown
🔀 **Topic:** [identified topic] 🛠️ **Primary Skills:** `@[skill-1]`, `@[skill-2]`
```

---

## Phase 2: Understanding Lock 🔒 (Socratic Gate)

**🛑 MANDATORY: Fully understand the problem before generating options.**

Extract Context:

- **Problem:** What needs solving?
- **Goal:** Desired outcome?
- **Constraints:** Any tech, time, or resource limitations?
- **Context:** Greenfield / Existing system?

If insufficient information, **ASK** clarifying questions before brainstorming. Once clear, summarize the problem and
wait for user confirmation.

---

## Phase 3: Divergent Thinking 🌐 (Generate Options)

Generate a minimum of **3 options** (include at least 1 unconventional approach). For each option provide:

- **Clear name**
- **Short description**
- **How it works** (high-level conceptually)

---

## Phase 4: Convergent Analysis 📊 (Evaluate)

Evaluate each option objectively based on Effort, Impact, Risk, Scalability, and Reversibility.

Create a Pros/Cons list and a Comparison Matrix.

---

## Phase 5: Recommendation & Decision 💡

Synthesize a recommendation explaining why it best fits the constraints.

**🛑 STOP HERE. Wait for the user to choose an option.** Do NOT suggest implementation details.

---

## Phase 6: Save Decision 📝

**🔔 Only execute this AFTER user selects an option.**

Save the accepted approach to `brainstorm-{slug}.md` and notify the user.

### Output Template:

```markdown
## 📝 Decision Record

**Date:** [YYYY-MM-DD] **Decision:** Option [X] - [Name] **Rationale:** [Why chosen]

### Context

[Problem being solved]

### Chosen Approach

[Detailed description]

### Trade-offs Accepted

- [Trade-off 1]
```
