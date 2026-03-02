import argparse
import shlex
import re
import random

import os
import time
import sqlite3
import subprocess
import fcntl

def acquire_global_lock(workspace):
    lock_file = os.path.join(workspace, ".agent_logs", "llm_request.lock")
    os.makedirs(os.path.dirname(lock_file), exist_ok=True)
    f = open(lock_file, "w")
    try:
        # Wait for the lock (blocking)
        fcntl.flock(f, fcntl.LOCK_EX)
        return f
    except Exception:
        return None

def release_global_lock(f):
    if f:
        fcntl.flock(f, fcntl.LOCK_UN)
        f.close()

def strip_ansi_codes(text):
    ansi_escape = re.compile(r'(?:\x1B[@-_]|[\x80-\x9F])[0-?]*[ -/]*[@-~]')
    return ansi_escape.sub('', text)

def get_db_connection(db_path):
    conn = sqlite3.connect(db_path)
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
    conn.row_factory = sqlite3.Row
    return conn

def read_unread_messages(db_path, role):
    try:
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        cursor.execute("SELECT id, sender_role, content FROM messages WHERE receiver_role = ? AND is_read = 0 ORDER BY created_at ASC", (role,))
        rows = cursor.fetchall()

        if rows:
            msg_ids = [r['id'] for r in rows]
            placeholders = ', '.join(['?'] * len(msg_ids))
            cursor.execute(f"UPDATE messages SET is_read = 1 WHERE id IN ({placeholders})", msg_ids)
            conn.commit()

        conn.close()
        return rows
    except Exception as e:
        print(f"DB Error: {e}")
        return []

def log_to_bus(db_path, sender, topic, content, receiver="all"):
    try:
        clean_content = strip_ansi_codes(content) if content else ""
        conn = get_db_connection(db_path)
        cursor = conn.cursor()
        cursor.execute("INSERT INTO messages (topic, sender_role, receiver_role, content) VALUES (?, ?, ?, ?)",
                      (topic, sender, receiver, clean_content))
        conn.commit()
        conn.close()
    except Exception as e:
        print(f"DB Error log_to_bus: {e}")

def run_engine_command(engine, prompt, workspace, db_path, role):
    try:
        if engine == "gemini" or engine == "copilot":
            cmd = f"{engine} -p {shlex.quote(prompt)}"
        else:
            # Default to 'run' command for kilocode and opencode
            cmd = f"{engine} run {shlex.quote(prompt)}"

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
                    # Optional: still log status to bus for dashboard but less frequently?
                    # For now, we collect everything to log as a single block at the end.

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

        # Check for rate limit keywords to alert user on the dashboard
        lower_content = final_content.lower()
        if "429" in lower_content or "rate limit" in lower_content:
            log_to_bus(db_path, "system", f"subagent_{role}", f"⚠️ RATE LIMIT DETECTED from {engine}! Taking a breather to avoid penalty...", receiver="all")
            time.sleep(10) # Auto-cool down a bit extra on 429

        log_to_bus(db_path, role, f"subagent_{role}", f"Task block execution finished. State saved.")
    except Exception as e:
        log_to_bus(db_path, "system", f"subagent_{role}", f"Failed to run {engine}: {e}")

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
    parser.add_argument('--engine', required=False, default="kilocode", choices=["kilocode", "opencode", "gemini", "copilot"])
    parser.add_argument('--resume', action='store_true', help="Resume agent with previous memory context")
    args = parser.parse_args()

    db_path = os.path.join(args.workspace, ".agent_logs", "multi_agent_bus.db")

    log_to_bus(db_path, "system", f"subagent_{args.role}", f"Daemon Worker for ROLE [{args.role}] started (Resume={args.resume}). Listening for messages...")

    # If planner and NOT resuming, execute initial task right away
    if args.role == 'planner' and args.task and not args.resume:
        log_to_bus(db_path, args.role, f"subagent_{args.role}", f"I am the Planner. Queueing LLM request for {args.engine}...")

        # Solution 1: Static instruction FIRST for Caching
        prompt = f"SYSTEM INSTRUCTION: {args.instruction}\n\n"
        prompt += f"TASK: {args.task}\n\n"

        # Solution 3: Strict Brevity
        prompt += "CRITICAL: Output ONLY a technical summary and the next delegation command. No conversational filler or greetings. Be ultra-concise."

        # Acquire global lock to serialize LLM requests
        lock = acquire_global_lock(args.workspace)
        try:
            run_engine_command(args.engine, prompt, args.workspace, db_path, args.role)
        finally:
            release_global_lock(lock)
            # Mandatory cooldown to prevent RPM spikes
            time.sleep(2)

    # Adaptive polling logic
    current_sleep = 2
    max_sleep = 15

    # Core event loop: wait for messages and process them
    while True:
        msgs = read_unread_messages(db_path, args.role)
        if msgs:
            # Active path: Reset sleep to fast mode
            current_sleep = 2
            log_to_bus(db_path, args.role, f"subagent_{args.role}", f"Received {len(msgs)} unread message(s). Batching to save tokens and avoid Rate Limits...")

            combined_msgs = ""
            for msg in msgs:
                combined_msgs += f"--- Message from '{msg['sender_role']}': ---\n{msg['content']}\n\n"

            # Solution 1: SYSTEM INSTRUCTION remains at the VERY TOP for prompt caching
            prompt = f"SYSTEM INSTRUCTION: {args.instruction}\n\n"

            if args.task:
                prompt += f"YOUR ORIGINAL MISSION/TASK WAS: {args.task}\n\n"

            # Inject recent historical context (only 2 since we use STATE.md now)
            recent_history = get_recent_history(db_path, args.role, limit=2)
            if recent_history:
                prompt += recent_history

            prompt += f"NEW UNREAD MESSAGES TO PROCESS:\n{combined_msgs}"
            prompt += "\n\nCRITICAL COMMUNICATION POLICY:\n"
            prompt += "1. EXECUTE YOUR TASK IN FULL: Finish your work before chatting.\n"
            prompt += "2. MINIMIZE CHATTER: No status updates. Only 'publish_message' for handoffs.\n"
            prompt += "3. BATCH RESPONSES: One message at the end of the run.\n"
            prompt += "4. NO QUESTIONS: Do not ask for confirmation.\n"
            prompt += "5. BE TOKEN-EFFICIENT: Use minimal words to achieve the goal.\n\n"

            # Solution 3: Strict Terse Output Instruction at the end
            prompt += "FINAL INSTRUCTION: Read history, take actions via tools, and output a technical summary ONLY. No filler. No intro. No outro."

            # Acquire global lock to serialize LLM requests
            lock = acquire_global_lock(args.workspace)
            try:
                run_engine_command(args.engine, prompt, args.workspace, db_path, args.role)
            finally:
                release_global_lock(lock)
                # Mandatory cooldown to prevent RPM spikes
                time.sleep(2)

        else:
            # Idle path: Backoff slowly to avoid eating CPU
            current_sleep = min(current_sleep + 2, max_sleep)

        # Sleep tight
        time.sleep(current_sleep)

if __name__ == "__main__":
    main()
