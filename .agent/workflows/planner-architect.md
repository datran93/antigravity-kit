---
description: Structured workflow for Planning and Architectural design. Orchestrates context discovery, task planning, and Coder delegation.
---

# /planner-architect - The Orchestrator's Workflow

This workflow guides the **Planner Agent** to analyze requirements, map the codebase, design the architecture, and initialize an actionable task plan before delegating to the **Code Agent**.

## 🚀 Execution Phase

### Phase 1: Contextual Discovery 🔍
Use MCP tools to build a comprehensive map of the impact area.
- Use `@mcp:ast-explorer` (`get_project_architecture`) to understand structural relationships (Py/Go/JS/TS).
- Use `find_by_name` and `grep_search` to locate relevant modules, APIs, and business logic.
- Use `@mcp:context7` (`query-docs`) to research latest syntax/libraries if needed.

### Phase 2: Skill & Pattern Discovery 🔍
After understanding the codebase context, identify the required domain expertise.
- Use `@mcp:skill-router` (`search_skills`) to find relevant skills based on the discovered modules and requirements.
- Read the identified `SKILL.md` files to align the architectural design with project-specific standards.

### Phase 3: Architectural Design 🏗️
Define the "North Star" for the implementation.
- Identify the core patterns (e.g., Clean Architecture, DDD, Hexagonal).
- Outline the schema changes or new API endpoints.
- Document architectural decisions in `DESIGN.md` or a local context if not already present.

### Phase 4: Task Plan Initialization 📋
Initialize the lifecycle of the task.
- Call `@mcp:context-manager` (`initialize_task_plan`) with a detailed list of atomic, executable steps.
- Set up checkpoints using `@mcp:context-manager` (`save_checkpoint`) at critical architectural milestones.
- Ensure each task in the plan is specific enough for a **Code Agent** to execute independently.

### Phase 5: Task Delegation 🤝 (Socratic Gate)
Before assigning the task to the Code Agent, perform a final validation with the USER.
- Describe the proposed architecture and plan summary.
- **MANDATORY**: Call `@mcp:mcp-multi-agent` (`enforce_socratic_gate`) if there are high-impact or ambiguous decisions.
- **COORDINATION**: Use `@mcp:mcp-multi-agent` (`publish_message`) with `target_role="coder"`, `topic="task"`, and a clear instruction of the atomic task to perform.
- The Planner remains active and monitors the bus for the next signal.

### Phase 6: Pipeline & Completion Management 🔄
Manage the lifecycle of the implementation as agents process the task.
- Periodically check for signals using `read_messages`.
- **WAIT** for the **Tester** to report success on the `topic="completion"` or similar.
- **ONLY AFTER** a success message from the Tester:
    - Call `@mcp:context-manager` (`complete_task_step`) to finalize the task in the plan.
    - If there are further tasks, repeat the **Phase 5** delegation for the next atomic step.
    - If all tasks are done, save the final context and report back to the USER.

## 🔴 Critical Constraints
1. **Ownership of Completion**: The Planner is the ONLY agent allowed to call `complete_task_step` and mark tasks as done.
2. **Persistent Team**: Use `publish_message` to communicate with the existing `coder`, `reviewer`, and `tester` daemons. Do not spawn new ones unless necessary.
3. **Atomic Steps**: Break down plans into independent, atomic steps that can be handed over one by one.

---

## 📌 Usage Example
`/planner-architect "Implementing a new user authentication flow with JWT and Redis caching"`
