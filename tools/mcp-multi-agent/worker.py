import argparse
import shlex
import re
import random

import os
import time
import sqlite3
import subprocess
import fcntl
import sys

def acquire_token_bucket(workspace, max_rpm=10, max_concurrent=2):
    db_path = os.path.join(workspace, ".agent_logs", "rate_limit.db")
    os.makedirs(os.path.dirname(db_path), exist_ok=True)

    lock_file = os.path.join(workspace, ".agent_logs", "rate_limit.lock")
    f = open(lock_file, "w")

    while True:
        fcntl.flock(f, fcntl.LOCK_EX)
        try:
            conn = sqlite3.connect(db_path)
            conn.execute('''CREATE TABLE IF NOT EXISTS requests (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            timestamp REAL,
                            status TEXT
                        )''')

            curr_time = time.time()

            # Clean up old timestamps (older than 60s)
            conn.execute("DELETE FROM requests WHERE timestamp < ?", (curr_time - 60,))

            # Reset 'running' stuck requests (older than 10 mins)
            conn.execute("DELETE FROM requests WHERE status = 'running' AND timestamp < ?", (curr_time - 600,))

            cursor = conn.cursor()
            cursor.execute("SELECT COUNT(*) FROM requests")
            rpm_count = cursor.fetchone()[0]

            cursor.execute("SELECT COUNT(*) FROM requests WHERE status = 'running'")
            concurrent_count = cursor.fetchone()[0]

            if rpm_count < max_rpm and concurrent_count < max_concurrent:
                cursor.execute("INSERT INTO requests (timestamp, status) VALUES (?, 'running')", (curr_time,))
                req_id = cursor.lastrowid
                conn.commit()
                conn.close()
                fcntl.flock(f, fcntl.LOCK_UN)
                return req_id, f, db_path
            else:
                conn.commit()
                conn.close()
        except Exception:
            pass

        fcntl.flock(f, fcntl.LOCK_UN)
        time.sleep(1)

def release_token_bucket(req_id, f, db_path):
    if f and req_id:
        fcntl.flock(f, fcntl.LOCK_EX)
        try:
            conn = sqlite3.connect(db_path)
            conn.execute("UPDATE requests SET status = 'done' WHERE id = ?", (req_id,))
            conn.commit()
            conn.close()
        except:
            pass
        fcntl.flock(f, fcntl.LOCK_UN)
        f.close()

def strip_ansi_codes(text):
    ansi_escape = re.compile(r'(?:\x1B[@-_]|[\x80-\x9F])[0-?]*[ -/]*[@-~]')
    return ansi_escape.sub('', text)

def get_db_connection(db_path):
    conn = sqlite3.connect(db_path, timeout=30.0)
    conn.execute('PRAGMA journal_mode=WAL')
    conn.execute('PRAGMA synchronous=NORMAL')
    conn.execute('''
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            topic TEXT NOT NULL,
            sender_role TEXT NOT NULL,
            receiver_role TEXT NOT NULL,
            content TEXT NOT NULL,
            is_read BOOLEAN DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ''')
    conn.execute('''
        CREATE TABLE IF NOT EXISTS agent_status (
            role TEXT PRIMARY KEY,
            status TEXT NOT NULL,
            last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            current_task TEXT,
            assigned_by TEXT DEFAULT ''
        )
    ''')
    conn.commit()
    conn.row_factory = sqlite3.Row
    return conn

def update_agent_status(db_path, role, status, current_task=None, assigned_by=None):
    try:
        conn = get_db_connection(db_path)
        conn.execute("BEGIN IMMEDIATE")
        cursor = conn.cursor()
        if current_task is not None:
            cursor.execute('''
                INSERT INTO agent_status (role, status, current_task, assigned_by, last_seen)
                VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
                ON CONFLICT(role) DO UPDATE SET
                    status=excluded.status,
                    current_task=excluded.current_task,
                    assigned_by=COALESCE(excluded.assigned_by, agent_status.assigned_by),
                    last_seen=CURRENT_TIMESTAMP
            ''', (role, status, current_task, assigned_by))
        else:
            cursor.execute('''
                INSERT INTO agent_status (role, status, current_task, last_seen)
                VALUES (?, ?, '', CURRENT_TIMESTAMP)
                ON CONFLICT(role) DO UPDATE SET
                    status=excluded.status,
                    current_task=excluded.current_task,
                    assigned_by=CASE WHEN excluded.status='IDLE' THEN '' ELSE agent_status.assigned_by END,
                    last_seen=CURRENT_TIMESTAMP
            ''', (role, status))
        conn.commit()
        conn.close()
    except Exception as e:
        print(f"Error updating agent status: {e}")

def read_unread_messages(db_path, role):
    try:
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        cursor.execute("SELECT id, sender_role, content FROM messages WHERE receiver_role = ? AND is_read = 0 ORDER BY created_at ASC", (role,))
        rows = cursor.fetchall()
        conn.close()
        return rows
    except Exception as e:
        print(f"DB Error read_unread_messages: {e}")
        return []

def mark_messages_as_read(db_path, msg_ids):
    if not msg_ids: return
    conn = None
    try:
        conn = get_db_connection(db_path)
        conn.execute("BEGIN IMMEDIATE")
        cursor = conn.cursor()
        placeholders = ', '.join(['?'] * len(msg_ids))
        cursor.execute(f"UPDATE messages SET is_read = 1 WHERE id IN ({placeholders})", msg_ids)
        conn.commit()
    except Exception as e:
        print(f"DB Error mark_messages_as_read: {e}")
        if conn: conn.rollback()
    finally:
        if conn: conn.close()

def log_to_bus(db_path, sender, topic, content, receiver="all"):
    max_retries = 5
    conn = None
    clean_content = strip_ansi_codes(content) if content else ""
    for attempt in range(max_retries):
        try:
            conn = get_db_connection(db_path)
            conn.execute("BEGIN IMMEDIATE")
            cursor = conn.cursor()
            cursor.execute("INSERT INTO messages (topic, sender_role, receiver_role, content) VALUES (?, ?, ?, ?)",
                          (topic, sender, receiver, clean_content))
            conn.commit()
            return
        except sqlite3.OperationalError as e:
            if conn: conn.rollback()
            if "locked" in str(e).lower():
                time.sleep(0.5 * (attempt + 1))
                continue
            print(f"DB Operational Error (log_to_bus): {e}")
            return
        except Exception as e:
            print(f"DB Error log_to_bus: {e}")
            return
        finally:
            if conn: conn.close()

def run_engine_command(engine, prompt, workspace, db_path, role, model=None):
    try:
        if engine in ["copilot", "gemini"]:
            flag = " --allow-all-tools" if engine == "copilot" else " --yolo"
            cmd = f"{engine} -p {shlex.quote(prompt)}{flag}"
        else:
            # Default to 'run' command for kilocode and opencode
            cmd = f"{engine} run {shlex.quote(prompt)}"
            if model:
                cmd += f" --model {shlex.quote(model)}"

            # Auto-approve if engine is kilocode (to prevent hangs)
            if engine == "kilocode":
                cmd += " --auto"

        my_env = os.environ.copy()

        log_to_bus(db_path, role, f"subagent_{role}", f"Executing task block with {engine}...")
        process = subprocess.Popen(
            cmd,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
            cwd=workspace,
            env=my_env
        )

        output_lines = []
        for line in iter(process.stdout.readline, ''):
            if line:
                clean_line = strip_ansi_codes(line.strip())
                if clean_line:
                    output_lines.append(clean_line)
                    # Print for server.py to capture and bridge to bus in real-time
                    print(clean_line, flush=True)

        process.stdout.close()
        process.wait()

        # Solution 4: Smart Truncation
        if len(output_lines) > 60:
            summary = f"--- [TECHNICAL LOG TRUNCATED: {len(output_lines)} lines] ---\n"
            summary += "\n".join(output_lines[:20]) # Top 20 lines
            summary += "\n\n... [... TRUNCATED ...] ...\n\n"
            summary += "\n".join(output_lines[-20:]) # Bottom 20 lines
            final_content = summary
        else:
            final_content = "\n".join(output_lines)

        # Solution 1: Use 'log' receiver to prevent other agents from reading noisy history
        if final_content.strip():
            log_to_bus(db_path, role, f"subagent_{role}", final_content, receiver="log")

        # Enhanced Rate Limit & Error Detection - check FULL original output
        full_text_joined = "\n".join(output_lines)
        full_text_lower = full_text_joined.lower()

        # Comprehensive keywords for various LLM providers/proxies
        rate_limit_keywords = [
            "429", "rate limit", "quota exceeded", "too many requests",
            "capacity reached", "credit balance", "exhausted", "retry after"
        ]

        if any(kw in full_text_lower for kw in rate_limit_keywords):
            log_to_bus(db_path, "system", f"alert_{role}", f"🛑 RATE LIMIT: {engine} reported quota/limit issues! Pausing 15s...", receiver="all")
            time.sleep(15)
        elif process.returncode != 0:
            log_to_bus(db_path, "system", f"alert_{role}", f"⚠️ ENGINE ERROR: {engine} exited with code {process.returncode}.", receiver="all")

        log_to_bus(db_path, role, f"subagent_{role}", f"Task block execution finished. State saved.")

        # For sequential mode to capture: print result summary to stdout
        print("\n=== AGENT_FINAL_SUMMARY ===", flush=True)
        print(final_content, flush=True)
        print("=== END_AGENT_FINAL_SUMMARY ===", flush=True)
    except Exception as e:
        log_to_bus(db_path, "system", f"alert_{role}", f"❌ FATAL: Failed to spawn {engine}: {e}", receiver="all")

def get_recent_history(db_path, role, limit=5):
    try:
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        # Solution 1: Filter out 'log' messages from LLM history
        cursor.execute("""
            SELECT sender_role, receiver_role, content
            FROM messages
            WHERE (sender_role = ? OR receiver_role = ? OR receiver_role = 'all')
              AND receiver_role != 'log'
            ORDER BY created_at DESC LIMIT ?
        """, (role, role, limit))
        rows = cursor.fetchall()
        conn.close()
        if not rows: return ""
        history = "=== RECENT HISTORY (MEMORY) ===\n"
        for r in reversed(rows):
            # Limit content to 500 chars to save tokens and prevent rate limits
            history += f"From {r['sender_role']} to {r['receiver_role']}: {r['content'][:500]}\n---\n"
        history += "================================\n\n"
        return history
    except Exception as e:
        return ""

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--workspace', required=True)
    parser.add_argument('--role', required=True)
    parser.add_argument('--instruction', required=True)
    parser.add_argument('--task', required=False, default="")
    parser.add_argument('--engine', required=False, default="opencode", choices=["kilocode", "opencode", "gemini", "copilot"])
    parser.add_argument('--model', required=False, default=None, help="Specific model to use (provider/model)")
    parser.add_argument('--resume', action='store_true', help="Resume agent with previous memory context")
    args = parser.parse_args()

    db_path = os.path.join(args.workspace, ".agent_logs", "multi_agent_bus.db")
    print(f"--- DAEMON BOOTSTRAP [{args.role}] ---", flush=True)
    print(f"Workspace: {args.workspace}", flush=True)
    print(f"DB Path: {db_path}", flush=True)

    log_to_bus(db_path, "system", f"subagent_{args.role}", f"Daemon Worker for ROLE [{args.role}] started (Resume={args.resume}). Listening for messages...")
    print("Logged start to bus.", flush=True)
    update_agent_status(db_path, args.role, "IDLE", current_task="Waiting for mission" if not args.task else args.task, assigned_by="system" if args.task else "")
    print("Updated agent status.", flush=True)

    # Process initial mission/task if provided
    if args.task and not args.resume:
        log_to_bus(db_path, args.role, f"subagent_{args.role}", f"Task started for ROLE [{args.role}]: {args.task[:100]}...")

        # Build execution prompt
        prompt = f"SYSTEM INSTRUCTION: {args.instruction}\n\n"
        prompt += f"TASK: {args.task}\n\n"
        prompt += "CRITICAL: Output ONLY a technical summary and the next command. No conversational filler. Be ultra-concise."

        # Acquire token bucket to allow concurrency but prevent rate limiting
        req_id, lock_f, rate_db = acquire_token_bucket(args.workspace, max_rpm=10, max_concurrent=2)
        try:
            update_agent_status(db_path, args.role, "WORKING", current_task=args.task, assigned_by="user")
            run_engine_command(args.engine, prompt, args.workspace, db_path, args.role, model=args.model)
            update_agent_status(db_path, args.role, "IDLE", current_task=None)
        finally:
            release_token_bucket(req_id, lock_f, rate_db)
            # Mandatory cooldown to prevent RPM spikes
            time.sleep(2)

        # EPHEMERAL SUBAGENT POLICY: If this was a one-off task for a subagent (non-planner), exit now.
        # This prevents process multiplication and 'zombie' listeners.
        if args.role != 'planner':
            log_to_bus(db_path, "system", f"subagent_{args.role}", f"Subagent for ROLE [{args.role}] finished ephemeral task and is exiting safely.")
            sys.exit(0)

    # Adaptive polling logic for persistent daemons
    current_sleep = 1
    max_sleep = 5
    idle_cycles = 0

    # Core event loop: wait for messages and process them
    while True:
        try:
            msgs = read_unread_messages(db_path, args.role)
            if msgs:
                # Active path: Reset sleep to fast mode
                current_sleep = 2
                idle_cycles = 0
                log_to_bus(db_path, args.role, f"subagent_{args.role}", f"Received {len(msgs)} unread message(s). Batching...")

                combined_msgs = ""
                senders = set()
                msg_ids = []
                for msg in msgs:
                    combined_msgs += f"--- Message from '{msg['sender_role']}': ---\n{msg['content']}\n\n"
                    senders.add(msg['sender_role'])
                    msg_ids.append(msg['id'])

                assigned_by = ", ".join(senders)

                # SYSTEM INSTRUCTION remains at the VERY TOP for prompt caching
                prompt = f"SYSTEM INSTRUCTION: {args.instruction}\n\n"

                if args.task:
                    prompt += f"YOUR ORIGINAL MISSION/TASK WAS: {args.task}\n\n"

                # Inject recent historical context
                recent_history = get_recent_history(db_path, args.role, limit=2)
                if recent_history:
                    prompt += recent_history

                prompt += f"NEW UNREAD MESSAGES TO PROCESS:\n{combined_msgs}"

                # ROLE-SPECIFIC ENFORCEMENT
                if args.role == 'reviewer':
                    prompt += "\n\nCRITICAL HANDOVER PROTOCOL (REVIEWER):\n"
                    prompt += "1. If you FIND ISSUES or BUGS: You MUST call 'publish_message' to notify the 'coder' role with detailed feedback so they can fix it.\n"
                    prompt += "2. If you APPROVE the changes (Review Passed): You MUST call 'publish_message' to notify the 'tester' role to start verification. DO NOT report back to 'planner' yet.\n"
                elif args.role == 'tester':
                    prompt += "\n\nCRITICAL REPORTING PROTOCOL (TESTER):\n"
                    prompt += "1. If tests FAIL: You MUST call 'publish_message' to notify the 'coder' role detailing which tests failed so they can fix the code.\n"
                    prompt += "2. If tests PASS: You MUST call 'publish_message' to notify the 'planner' role with the results summary. This is the only way to advance the project.\n"
                elif args.role == 'planner':
                    prompt += "\n\nCRITICAL ORCHESTRATION PROTOCOL (PLANNER):\n"
                    prompt += "1. You are the DISPATCHER. You MUST proactively ask 'coder' for updates if the project stalls.\n"
                    prompt += "2. If the project pipeline is stalled (everyone is idling), you MUST immediately assign a new task to the 'coder'. Avoid re-triggering 'reviewer' or 'tester' unless new code changes are ready.\n"
                    prompt += "3. DO NOT wait for them to reach out to you; you lead the coordination through the coder first.\n"

                prompt += "\n\nCRITICAL COMMUNICATION POLICY:\n"
                prompt += "1. EXECUTE YOUR TASK IN FULL.\n"
                prompt += "2. NO CHATTER: Provide only technical output.\n"
                prompt += "3. TERMINATE: Do not ask questions; finish and exit.\n\n"
                prompt += "FINAL INSTRUCTION: Use tools, then output a Technical Summary ONLY."

                # Acquire token bucket
                req_id, lock_f, rate_db = acquire_token_bucket(args.workspace, max_rpm=10, max_concurrent=2)
                try:
                    update_agent_status(db_path, args.role, "WORKING", current_task=combined_msgs[:500], assigned_by=assigned_by)
                    run_engine_command(args.engine, prompt, args.workspace, db_path, args.role, model=args.model)
                    update_agent_status(db_path, args.role, "IDLE", current_task=None)

                    # ONLY MARK AS READ AFTER SUCCESSFUL FINISH
                    mark_messages_as_read(db_path, msg_ids)
                finally:
                    release_token_bucket(req_id, lock_f, rate_db)
                    time.sleep(2)

            else:
                # Idle path: Backoff slowly
                idle_cycles += 1
                current_sleep = min(current_sleep + 2, max_sleep)

                # In sequential mode, the planner just waits for user or task signals.
                # No proactive team pulse needed as workers are ephemeral.
                pass

        except Exception as loop_err:
            log_to_bus(db_path, "system", f"alert_{args.role}", f"Error in agent loop: {str(loop_err)}")
            time.sleep(5)

        # Sleep tight
        time.sleep(current_sleep)
        if random.random() < 0.05: # Occasional heartbeat to show daemon is alive
             log_to_bus(db_path, "system", f"status_{args.role}", "Agent is idling, waiting for messages.", receiver="log")

if __name__ == "__main__":
    main()
