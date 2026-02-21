---
description: Database interaction workflow using @mcp:database-inspector.
---

# /db - Database Interaction Workflow

Enable agents to explore, query, and manipulate databases using the `@mcp:database-inspector` server tools.

## When to Use

- `/db [query or instruction]`
- Keywords: "database", "table", "schema", "query", "sql", "rows", "column", "select"

## 🔴 Critical Rules

1. **Read-Only First**: Always start with schema exploration before querying.
2. **Limit Results**: The `run_read_query` tool defaults to LIMIT 10, but explicitly state limits if you need more data safely.
3. **Connection Management**: Use `DATABASE_URL` from `.env` files if available (using `grep_search` or `view_file` to find it), or ask for the SQLAlchemy connection string.
4. **Data Privacy**: Avoid printing sensitive PII (Personally Identifiable Information).
5. **No Write Operations**: Only read-only queries are supported by `@mcp:database-inspector`. Do not attempt to run INSERT/UPDATE/DELETE.

---

## Phase 1: Environment Setup 🔐

1.  **Locate Credentials**: Check for `DATABASE_URL` or relevant environment variables in `.env` files.
2.  **Verify Connection**: Use `list_tables(connection_string)` to ensure the server is reachable and valid.

---

## Phase 2: Schema Exploration 🗂️

If the query is generic ("show me the db", "describe tables"), use:

- **List Tables**: `list_tables(connection_string=...)`
- **Describe Table**: `inspect_schema(connection_string=..., table_name=...)`

---

## Phase 3: Query Execution ⚡

1.  **Format SQL**: Ensure SQL is valid for the target database type (Postgres, MySQL, etc.).
2.  **Execute**: Run `run_read_query(connection_string=..., query="YOUR_SQL_QUERY", limit=10)`.
3.  **Error Handling**: If a query fails, analyze the error (e.g., from the tool output string) and fix the SQL.

---

## Phase 4: Reporting 📊

Summarize the results:
- **Rows returned**
- **Schema outline (if queried)**
- **Next suggested steps** (e.g., "Would you like me to analyze the trends in this data?")
