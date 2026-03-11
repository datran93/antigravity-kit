---
description: Generate hierarchical codebase visualizations and execution flow maps, similar to Windsurf Codemaps.
---

# 🗺️ Codemap Workflow (Codebase Visualization)

This workflow guides you to create a visual representation of the project's architecture, data flows, and module
relationships. It helps on-board new developers (and refresh yourself) on how the system is structured.

## 🚀 Execution Phase

### Phase 1: Structural Discovery 🔍

- Use `@mcp:ast-explorer` (`get_project_architecture`, `search_symbol`) to map the high-level project structure and find
  exact symbol locations.
- Identify the main entry points (e.g., `main.go`, `index.ts`, `app/server.py`).
- List the primary modules, services, and directories.

### Phase 2: Execution Flow Mapping 🌊

- Pick a core user journey (e.g., "User Login" or "Data Ingestion").
- Trace the call stack across files and services.
- Note the data transformations and external dependencies (DB, Redis, Third-party APIs).

### Phase 3: Visual Generation (Mermaid) 🎨

- Generate a comprehensive `CODEMAP.md` file in the root.
- **MANDATORY**: Include at least one **Architecture Diagram** (Mermaid `graph TD`).
- **MANDATORY**: Include at least one **Sequence Diagram** (Mermaid `sequenceDiagram`) for a primary execution flow.
- Use distinct nodes for different layers (UI, API, Service, DAO).

### Phase 4: Narrative Summary 📝

- Beneath each diagram, provide a narrative explanation of "The Life of a Request".
- Link to the actual files mentioned using `file_path:line_number`.

### Phase 5: Result Delivery 🏁

- Present the `CODEMAP.md` to the user.
- Offer to deep-dive into specific components.

## 🔴 Critical Constraints

1. **Fact-Based**: Do not guess names. Verify every function call and import.
2. **Visual Focus**: The goal is clarity. Use subgraphs in Mermaid to group related items.
3. **No Placeholders**: Every node in the diagram must correspond to a real code entity.

---

## 📌 Usage Example

`/codemap "Visualize the entire module dependency graph"`
`/codemap "Map the flow of an incoming HTTP request from route to database"`
