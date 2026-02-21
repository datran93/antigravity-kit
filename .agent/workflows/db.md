---
description: Database interaction workflow using usql CLI.
---

# /db - Database Interaction Workflow

Enable agents to explore, query, and manipulate databases using the `usql` CLI.

## When to Use

- `/db [query or instruction]`
- Keywords: "database", "table", "schema", "query", "sql", "rows", "column", "select", "insert", "update", "delete"

## 🔴 Critical Rules

1. **Read-Only First**: Always start with `SELECT` or schema exploration before making changes.
2. **Limit Results**: Always append `LIMIT 10` or similar to prevent overwhelming the context window.
3. **Connection Management**: Use `DATABASE_URL` from `.env` files if available, or ask for the connection string.
4. **Data Privacy**: Avoid printing sensitive PII (Personally Identifiable Information) unless explicitly required.
5. **Verify Impact**: Before `UPDATE` or `DELETE`, run a `SELECT` with the same `WHERE` clause to count affected rows.

---

## Phase 1: Environment Setup 🔐

1.  **Locate Credentials**: Check for `DATABASE_URL` or relevant environment variables in `.env` files.
2.  **Verify CLI**: Ensure `usql` is accessible.
3.  **Test Connection**: Run a simple `usql $DATABASE_URL -c "SELECT 1"` to verify access.

---

## Phase 2: Schema Exploration 🗂️

If the query is generic ("show me the db", "describe tables"), use:

- **List Tables**: `usql $DATABASE_URL -c "\dt"`
- **Describe Table**: `usql $DATABASE_URL -c "\d table_name"`
- **List Databases**: `usql $DATABASE_URL -c "\l"`

---

## Phase 3: Query Execution ⚡

1.  **Format SQL**: Ensure SQL is valid for the target database type (Postgres, MySQL, SQLite, etc.).
2.  **Execute**: Run `usql $DATABASE_URL -c "YOUR_SQL_QUERY"`.
3.  **Output Format**: Use `-F csv` or default table format depending on the required level of detail.
4.  **Error Handling**: If a query fails, analyze the error (syntax, missing table, etc.) and suggest a fix.

---

## Phase 4: Data Manipulation (Safe Mode) 🛡️

For `INSERT`, `UPDATE`, or `DELETE`:

1.  **Dry Run / Count**: Run `SELECT count(*) FROM ... WHERE ...` first.
2.  **Explicit Verification**: Report the count to the user and confirm before proceeding if the impact is large.
3.  **Execute & Verify**: Run the command and then verify the change with a subsequent `SELECT`.

---

## Phase 5: Reporting 📊

Summarize the results:
- **Rows affected/returned**
- **Schema changes (if any)**
- **Next suggested steps** (e.g., "Would you like me to analyze the trends in this data?")
