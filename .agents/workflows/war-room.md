---
description: Structured workflow for emergency incident response. Enforces Evidence Checklists before arbitrary fixes.
---

# 🚨 War-Room (Incident Response) Workflow

Use this workflow to prevent rushed code changes when a critical bug or production incident is reported by the USER.
**DO NOT propose immediate code fixes.** You must enforce rigorous data collection before identifying the root cause.

## 🚀 Execution Sequence

### Phase 1: Context Isolation & Evidence Checklist 🗂️

Gather undeniable proof of what failed.

- Acknowledge the incident. Do not guess what the problem is.
- Create an **Evidence Checklist** (Logs, DB Queries, Network Traces, Stack Traces, Terminal outputs).
- Read the necessary log files using `view_file`.
- Check monitoring/infrastructure tools or databases using `@mcp:database-inspector` or `@mcp:mcp-http-client`.

### Phase 2: Root Cause Formularization 🧪

Propose the exact mechanism of failure instead of patching surface symptoms.

- Summarize the gathered evidence in markdown.
- Formulate a clear hypothesis on why the incident occurred based on system constraints or `.agents/rules/ANCHORS.md`.
- E.g., "The logs show a duplicate key error in PostgreSQL, suggesting an atomic constraint was violated during a
  concurrent POST request."

### Phase 3: Socratic Resolution Gate 🛑

Validate the plan with the USER.

- Tell the USER your root cause theory.
- Present a concrete, multi-stage remediation plan. Include immediate safe mitigations (e.g., rolling back the commit,
  reverting a config flag) vs long-term structural fixes.
- **WAIT** explicitly for USER approval.

### Phase 4: Resolution & Retrospective Execution 🛠️

Once authorized, execute the patch.

- Take on the `[Role: 💻 Coder]` to apply the code fix cleanly. Keep changes scoped exclusively to the Incident
  Resolution.
- After fixing, immediately switch to the `[Role: 🏗️ Planner]` and trigger **Context Compression (KI Generation)**:
  write an incident retrospective in `/knowledge/` containing what broke, why, and how to prevent it in the future.

## 🔴 Critical Constraints

1. **Blind Fixes are Forbidden**: Do not edit files (`replace_file_content`) until you have completed the Evidence
   Checklist and Formularization.
2. **Prioritize Stabilization**: Before optimizing code, ensure the service is returned to a working state securely.
3. **Save Postmortems**: A session must end with an incident Ki written as a Markdown file in `/knowledge/`.

---

## 📌 Usage Example

`/war-room "Production API is currently returning 500 internal server errors on user login"`
