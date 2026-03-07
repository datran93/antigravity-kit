# ⚓ ANCHORS (Immutable Facts & Guardrails)

This file contains the absolute, non-negotiable project facts, architectural constraints, and guardrails for the agent
system. These **Anchors** must survive context compaction and session restarts.

As a Self-Executing Agent, you MUST respect these facts before executing any code changes or planning any project
tactics.

## 🛡️ Core Guardrails

_(Add your project-specific guardrails here. Example below:)_

- **Backend Stack**: Golang >= 1.21.
- **Frontend Stack**: React 18 / Next.js 14 App Router.
- **Database**: PostgreSQL 15+. DO NOT use MongoDB or MySQL.
- **State Boundaries**: `active_files` must be explicitly tracked across transitions to prevent blind writes.
- **No Destruction**: Do not delete existing API contracts without explicit confirmation from the USER.
- **MCP Go Porting**: All new MCP tools must be written in Native Go using `github.com/mark3labs/mcp-go`. Strict type
  assertion for `map[string]any` is required when parsing tool arguments. Avoid heavy external databases/servers if an
  in-memory or purely native math equivalent exists (e.g. replacing ChromaDB with Go vector math). Target paths in
  `mcp_config.json` must point directly to the compiled Go binary.
