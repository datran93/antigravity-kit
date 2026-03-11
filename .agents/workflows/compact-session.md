---
description: Structured workflow for session compaction to prevent context drift and memory exhaustion.
---

# 🗜️ Context Compaction (KI Generation) Workflow

Use this workflow to compress the current working memory (context) when completing a major phase (`Tactic`) of the
project, or when shifting focus significantly. This mimics HiveMind's `compact_session` to prevent the agent from
overloading memory and ensures decisions are recorded as permanent Knowledge Items (KIs).

## 🚀 Execution Sequence

### Phase 1: Context Aggregation 📥

Gather all the outputs and knowledge from the recently completed sequence.

- Review the `active_files` that were modified during this Tactic.
- Synthesize the architectural decisions, patterns, or new library setups that emerged.
- Identify unresolved debts, bugs, or items explicitly delayed.

### Phase 2: KI (Knowledge Item) Generation & Memory Flush 🧠🧹

Persist the knowledge and prune the context automatically.

- Call `@mcp:context-manager` (`compact_memory`) with the `tactic_name`, `summary`, and `decisions` formed from Phase 1.
- This tool will **automatically**:
  - Generate the Markdown KI file inside the `knowledge/` directory.
  - Reset the `active_files` and update the checkpoint notes with the KI path.
  - Reset the drift failure counter and intent locks.
- **Mental Flush**: Explicitly state to the USER that context compaction is complete. Actively ignore previous tool
  outputs (CLI logs, debug traces, test outputs), retaining ONLY the global `.agents/rules/ANCHORS.md` and the objective
  of the next `Tactic`.

## 🔴 Critical Constraints

1. **Always Compact Between Tactics**: Do not start a completely new module without compacting the previous one.
2. **Actionable Knowledge**: KIs must be concise and actionable logic rules, not just a raw data dump of source code.
3. **No Code Edits**: This phase is strictly for writing documentation and managing state.

---

## 📌 Usage Example

`/compact-session "We finished the Auth module, let's compress context before starting the Payment module"`
