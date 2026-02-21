---
description: Redis interaction workflow using redis-cli.
---

# /redis - Redis Interaction Workflow

Enable agents to explore, query, and manipulate Redis data using the `redis-cli`.

## When to Use

- `/redis [command or instruction]`
- Keywords: "redis", "cache", "key", "get", "set", "hget", "hset", "expire", "ttl", "flush"

## 🔴 Critical Rules

1. **Safety First**: Never use `KEYS *` on production-like environments. Use `SCAN 0 COUNT 10` to discover keys.
2. **Limit Output**: For lists or sets, always limit the range (e.g., `LRANGE key 0 10`).
3. **Environment**: Prefer using `REDIS_URL` from `.env` files.
4. **Data Types**: Always check the type of a key (`TYPE key`) before attempting to read it with a specific command (e.g., don't use `GET` on a Hash).
5. **Mutation Safety**: Before `DEL`, verify the key exists and check its content.

---

## Phase 1: Environment & Connectivity 🔐

1.  **Locate Connection**: Check `scripts/.env` or project root `.env` for `REDIS_URL`.
2.  **Test Ping**: Run `redis-cli -u $REDIS_URL PING` to ensure the server is reachable.
3.  **Authentication**: If `REDIS_URL` is missing, look for `REDIS_HOST`, `REDIS_PORT`, and `REDIS_PASSWORD`.

---

## Phase 2: Key Discovery 🔍

If the user asks to "show keys" or "find keys related to X":

- **Scan Keys**: `redis-cli -u $REDIS_URL --scan --pattern "*X*"`
- **Check Type**: `redis-cli -u $REDIS_URL TYPE "key_name"`
- **Check TTL**: `redis-cli -u $REDIS_URL TTL "key_name"`

---

## Phase 3: Reading Data 📖

Execute the appropriate command based on the data type:

- **String**: `GET key`
- **Hash**: `HGETALL key` (or `HMGET key field1 field2`)
- **List**: `LRANGE key 0 10`
- **Set**: `SMEMBERS key` (careful with large sets, consider `SSCAN`)
- **Sorted Set**: `ZRANGE key 0 10 WITHSCORES`

Example: `redis-cli -u $REDIS_URL HGETALL my_hash`

---

## Phase 4: Manipulating Data ✍️

1.  **Set/Update**: `SET key value` or `HSET key field value`.
2.  **Delete**: `DEL key`.
3.  **Verification**: After any write operation, immediately run a read command to verify the change.

---

## Phase 5: Reporting 📊

Summarize the interaction:
- **Operation performed** (Read/Write/Delete)
- **Data returned** (Formatted for readability)
- **Current status of the key** (Type, TTL)
