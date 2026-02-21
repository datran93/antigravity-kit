---
description: Redis interaction workflow using @mcp:database-inspector.
---

# /redis - Redis Interaction Workflow

Enable agents to explore, query, and manipulate Redis data safely via `@mcp:database-inspector`.

## When to Use

- `/redis [command or instruction]`
- Keywords: "redis", "cache", "key", "get", "hget", "expire", "ttl"

## 🔴 Critical Rules

1. **Safety First**: The `list_tables` tool will return a safe sample of keys. Never run raw commands like `KEYS *` through `run_read_query`.
2. **Limit Output**: For lists or sets via `run_read_query`, limit the range if possible.
3. **Environment**: Prefer using `REDIS_URL` from `.env` files (e.g. `redis://localhost:6379/0`).
4. **Data Types**: Use the `inspect_schema` tool to check the type and properties of a key before attempting to read it fully.
5. **No Write Operations**: Only read-only queries are supported by `@mcp:database-inspector`. Do not attempt to run SET, HSET, DEL, etc.

---

## Phase 1: Environment & Connectivity 🔐

1.  **Locate Connection**: Check `scripts/.env` or project root `.env` for `REDIS_URL`.
2.  **Verify Connection**: Use `list_tables(connection_string)` to ensure the server is reachable and valid.

---

## Phase 2: Key Discovery 🔍

If the user asks to "show keys" or "find keys related to X":

- **List sample keys**: `list_tables(connection_string=...)`
- **Inspect Specific Key (Type/TTL)**: `inspect_schema(connection_string=..., table_name="key_name")`

---

## Phase 3: Reading Data 📖

Execute the appropriate command via `run_read_query` based on the data type (which you found in Phase 2):

- **String**: `GET key`
- **Hash**: `HGETALL key` (or `HMGET key field1 field2`)
- **List**: `LRANGE key 0 10`
- **Set**: `SMEMBERS key` (careful with large sets, consider `SSCAN`)
- **Sorted Set**: `ZRANGE key 0 10 WITHSCORES`

Example: `run_read_query(connection_string=..., query="HGETALL my_hash")`

---

## Phase 4: Reporting 📊

Summarize the interaction:
- **Operation performed** (Read)
- **Data returned** (Formatted for readability)
- **Current status of the key** (Type, TTL)
