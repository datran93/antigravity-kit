import json
from mcp.server.fastmcp import FastMCP
from sqlalchemy import create_engine, inspect, text
from sqlalchemy.engine.reflection import Inspector

# Khởi tạo MCP Server
mcp = FastMCP("McpDatabaseInspector")

def get_inspector(connection_string: str) -> Inspector:
    """Helper method to create SQLAlchemy engine and return inspector."""
    # Xử lý an toàn: Có thể thêm một vài filter ở đây (phòng trường hợp SQL params độc hại)
    engine = create_engine(connection_string)
    return inspect(engine)

@mcp.tool()
def list_tables(connection_string: str) -> str:
    """
    List all tables and views in the database.
    For Redis, it returns a summary of keys by exploring the database.

    Args:
        connection_string: SQLAlchemy connection string.
                           Example Postgres: postgresql://user:password@localhost:5432/dbname
                           Example MySQL: mysql+pymysql://user:password@localhost:3306/dbname
                           Example Redis: redis://localhost:6379/0
    """
    try:
        if connection_string.startswith("redis://") or connection_string.startswith("rediss://"):
            import redis
            r = redis.from_url(connection_string)
            dbsize = r.dbsize()
            # Scan a few keys to give a sample
            _, keys = r.scan(0, count=100)

            result = [f"🗄 REDIS DATABASE (Size: {dbsize} keys)\n"]
            result.append("### Sample Keys (up to 100):")
            for k in keys:
                result.append(f"- {k.decode('utf-8')}")

            if not keys:
                return "❌ No keys found in this Redis database."

            return "\n".join(result)

        inspector = get_inspector(connection_string)
        tables = inspector.get_table_names()
        views = inspector.get_view_names()

        result = [f"🗄 DATABASE ENTITIES\n"]
        result.append("### Tables:")
        for t in tables:
            result.append(f"- {t}")

        result.append("\n### Views:")
        for v in views:
            result.append(f"- {v}")

        if not tables and not views:
            return "❌ No tables or views found in this database."

        return "\n".join(result)
    except Exception as e:
        return f"❌ Error connecting or reading database: {str(e)}"

@mcp.tool()
def inspect_schema(connection_string: str, table_name: str) -> str:
    """
    Get detailed schema information for a specific table. Includes columns, types, primary keys, and foreign keys.
    For Redis, table_name acts as the key name. Returns key type and properties.

    Args:
        connection_string: SQLAlchemy connection string.
        table_name: The exact name of the table to inspect.
    """
    try:
        if connection_string.startswith("redis://") or connection_string.startswith("rediss://"):
            import redis
            r = redis.from_url(connection_string)
            if not r.exists(table_name):
                return f"❌ Key '{table_name}' does not exist in Redis."

            key_type_bytes = r.type(table_name)
            key_type = key_type_bytes.decode('utf-8') if isinstance(key_type_bytes, bytes) else key_type_bytes
            ttl = r.ttl(table_name)

            result = [f"📋 SCHEMA FOR KEY: `{table_name}`\n"]
            result.append(f"### Metadata:")
            result.append(f"- **Type**: {key_type}")
            result.append(f"- **TTL**: {ttl} seconds")

            result.append("\n### Details:")
            if key_type == "string":
                # Check length
                val_len = r.strlen(table_name)
                result.append(f"- String Length: {val_len}")
            elif key_type == "hash":
                hlen = r.hlen(table_name)
                result.append(f"- Hash Entries: {hlen}")
            elif key_type == "list":
                llen = r.llen(table_name)
                result.append(f"- List Length: {llen}")
            elif key_type == "set":
                scard = r.scard(table_name)
                result.append(f"- Set Members: {scard}")
            elif key_type == "zset":
                zcard = r.zcard(table_name)
                result.append(f"- Sorted Set Members: {zcard}")

            return "\n".join(result)

        inspector = get_inspector(connection_string)

        if table_name not in inspector.get_table_names() and table_name not in inspector.get_view_names():
            return f"❌ Table '{table_name}' does not exist in the database."

        columns = inspector.get_columns(table_name)
        pk_constraint = inspector.get_pk_constraint(table_name)
        foreign_keys = inspector.get_foreign_keys(table_name)

        result = [f"📋 SCHEMA FOR TABLE: `{table_name}`\n"]

        # Setup Primary Keys display
        pks = pk_constraint.get('constrained_columns', [])

        result.append("### Columns:")
        for col in columns:
            col_name = col['name']
            col_type = str(col['type'])
            nullable = "NULL" if col.get('nullable') else "NOT NULL"

            # Format row
            prefix = "🔑 " if col_name in pks else "   "
            result.append(f"{prefix}`{col_name}` : **{col_type}** ({nullable})")

        if foreign_keys:
            result.append("\n### Foreign Keys:")
            for fk in foreign_keys:
                constrained_cols = ", ".join(fk['constrained_columns'])
                referred_table = fk['referred_table']
                referred_cols = ", ".join(fk['referred_columns'])
                result.append(f"- `{constrained_cols}` -> `{referred_table}`.`{referred_cols}`")

        return "\n".join(result)
    except Exception as e:
        return f"❌ Error inspecting table '{table_name}': {str(e)}"

@mcp.tool()
def run_read_query(connection_string: str, query: str, limit: int = 10) -> str:
    """
    Execute a read-only SQL query to preview raw data. DO NOT use this for INSERT/UPDATE/DELETE.
    Results are returned in JSON format.
    For Redis, query is the Redis command (e.g. 'GET mykey' or 'HGETALL myhash'). It only supports read-only commands.

    Args:
        connection_string: SQLAlchemy connection string.
        query: The raw SQL query to run (e.g., 'SELECT * FROM users'). For Redis: 'GET mykey'.
        limit: Max number of rows to return (defaults to 10 to prevent massive memory dumps).
    """
    try:
        if connection_string.startswith("redis://") or connection_string.startswith("rediss://"):
            import redis
            r = redis.from_url(connection_string)
            query_parts = query.strip().split()
            if not query_parts:
                return "❌ Empty query"
            command = query_parts[0].upper()

            safe_commands = [
                'GET', 'MGET', 'HGET', 'HGETALL', 'HMGET', 'HKEYS', 'HVALS', 'HLEN',
                'LRANGE', 'LLEN', 'LINDEX',
                'SMEMBERS', 'SCARD', 'SISMEMBER',
                'ZRANGE', 'ZCARD', 'ZSCORE', 'ZREVRANGE',
                'TYPE', 'TTL', 'PTTL', 'EXISTS',
                'SCAN', 'HSCAN', 'SSCAN', 'ZSCAN',
                'INFO', 'DBSIZE', 'PING', 'KEYS'
            ]

            if command not in safe_commands:
                return f"❌ SECURITY BLOCK: Command '{command}' is not allowed or is a write operation. Only Read-only Redis queries are allowed."

            # Execute command
            res = r.execute_command(*query_parts)

            # Decode recursive function for bytes
            def decode_redis(obj):
                if isinstance(obj, bytes):
                    return obj.decode('utf-8', errors='replace')
                elif isinstance(obj, list) or isinstance(obj, tuple) or isinstance(obj, set):
                    return [decode_redis(i) for i in obj]
                elif isinstance(obj, dict):
                    return {decode_redis(k): decode_redis(v) for k, v in obj.items()}
                return obj

            decoded_res = decode_redis(res)

            # Formatting output
            json_result = json.dumps({"result": decoded_res}, indent=2)

            output = [
                f"✅ REDIS QUERY SUCCESS",
                f"CMD = `{query}`",
                "```json",
                json_result,
                "```"
            ]
            return "\n".join(output)

        query_upper = query.strip().upper()
        if any(keyword in query_upper for keyword in ['INSERT', 'UPDATE', 'DELETE', 'DROP', 'ALTER', 'TRUNCATE', 'GRANT', 'REVOKE']):
            return "❌ SECURITY BLOCK: Only SELECT/Read-only queries are allowed."

        engine = create_engine(connection_string)

        # Enforce LIMIT inside query logically if possible, or visually check
        if "LIMIT " not in query_upper:
            # We strictly enforce limit on the cursor level instead of parsing complex SQL
            safe_query = f"{query} LIMIT {limit}"
        else:
            safe_query = query

        with engine.connect() as conn:
            result_proxy = conn.execute(text(safe_query))
            rows = result_proxy.fetchall()
            keys = result_proxy.keys()

            if not rows:
                return "✅ Query executed successfully. 0 rows returned."

            data = [dict(zip(keys, row)) for row in rows]

            # Convert to string with proper format handling for datetime/UUID
            def json_serial(obj):
                try:
                    return str(obj)
                except:
                    return "Unserializable Data"

            json_result = json.dumps(data, default=json_serial, indent=2)

            output = [
                f"✅ QUERY RESULTS (Limited to {limit} rows max)",
                f"SQL = `{safe_query}`",
                "```json",
                json_result,
                "```"
            ]
            return "\n".join(output)

    except Exception as e:
        return f"❌ Error executing query: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')
