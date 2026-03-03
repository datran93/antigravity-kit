import os
import sys
import sqlite3
import json
import shlex
from datetime import datetime
from typing import List, Optional
import subprocess
import threading
from mcp.server.fastmcp import FastMCP

mcp = FastMCP("McpMultiAgent")

def get_db_connection(workspace_path: str):
    if not workspace_path:
        raise ValueError("workspace_path is required")

    # Create DB in .agent_logs folder of the respective workspace
    log_dir = os.path.join(workspace_path, ".agent_logs")
    os.makedirs(log_dir, exist_ok=True)
    db_path = os.path.join(log_dir, "multi_agent_bus.db")

    conn = sqlite3.connect(db_path, timeout=30.0)
    conn.execute('PRAGMA journal_mode=WAL')
    conn.execute('PRAGMA synchronous=NORMAL')
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    cursor.execute('''
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
    conn.commit()
    return conn
def internal_log(workspace_path: str, sender: str, topic: str, content: str, receiver: str = "all"):
    conn = None
    try:
        conn = get_db_connection(workspace_path)
        conn.execute("BEGIN IMMEDIATE")
        cursor = conn.cursor()
        cursor.execute('INSERT INTO messages (topic, sender_role, receiver_role, content) VALUES (?, ?, ?, ?)',
                      (topic, sender, receiver, content))
        conn.commit()
    except:
        if conn: conn.rollback()
    finally:
        if conn: conn.close()

@mcp.tool()
def publish_message(workspace_path: str, topic: str, sender_role: str, receiver_role: str, content: str) -> str:
    """
    Publish a message to the internal message bus.
    Allows agents to communicate asynchronously without polluting the main user conversation.
    """
    # ZERO-TOKEN PRE-CHECK (Idea 2): Validate syntax before bothering the Reviewer
    if sender_role.startswith("coder") and receiver_role == "reviewer":
        try:
            # Check for modified or new Python files
            git_cmd_py = "git ls-files -m --others --exclude-standard | grep '\\.py$'"
            py_files = subprocess.run(git_cmd_py, shell=True, cwd=workspace_path, capture_output=True, text=True).stdout.strip().split('\n')
            py_files = [f for f in py_files if f]
            if py_files:
                python_exec = sys.executable if sys.executable else "python3"
                py_res = subprocess.run([python_exec, "-m", "py_compile"] + py_files, cwd=workspace_path, capture_output=True, text=True)
                if py_res.returncode != 0:
                    err_msg = f"❌ AUTO-PRE-CHECK FAILED: Python syntax errors found:\n{py_res.stderr}\n{py_res.stdout}"
                    internal_log(workspace_path, "system", f"lint_{sender_role}", err_msg, receiver="all")
                    return err_msg

            # Check for modified TypeScript/JS files if tsconfig exists
            if os.path.exists(os.path.join(workspace_path, "tsconfig.json")):
                ts_res = subprocess.run(["npx", "tsc", "--noEmit"], cwd=workspace_path, capture_output=True, text=True)
                if ts_res.returncode != 0:
                    err_msg = f"❌ AUTO-PRE-CHECK FAILED: TypeScript errors found:\n{ts_res.stdout}"
                    internal_log(workspace_path, "system", f"lint_{sender_role}", err_msg, receiver="all")
                    return err_msg

            # Check for Go files
            git_cmd_go = "git ls-files -m --others --exclude-standard | grep '\\.go$'"
            go_files = subprocess.run(git_cmd_go, shell=True, cwd=workspace_path, capture_output=True, text=True).stdout.strip().split('\n')
            go_files = [f for f in go_files if f]
            if go_files:
                go_res = subprocess.run(["go", "vet", "./..."], cwd=workspace_path, capture_output=True, text=True)
                if go_res.returncode != 0:
                    err_msg = f"❌ AUTO-PRE-CHECK FAILED: Go syntax/vet errors found:\n{go_res.stderr}"
                    internal_log(workspace_path, "system", f"lint_{sender_role}", err_msg, receiver="all")
                    return err_msg

        except Exception as precheck_err:
            pass # If pre-check fails to run (e.g. no git, no node, no go), just ignore it and proceed handoff

    # Solution: Retry logic for DB locking
    max_retries = 5
    for attempt in range(max_retries):
        try:
            conn = get_db_connection(workspace_path)
            conn.execute("BEGIN IMMEDIATE")
            cursor = conn.cursor()
            cursor.execute('''
                INSERT INTO messages (topic, sender_role, receiver_role, content)
                VALUES (?, ?, ?, ?)
            ''', (topic, sender_role, receiver_role, content))
            conn.commit()
            msg_id = cursor.lastrowid
            conn.close()
            return f"✅ Message {msg_id} published from {sender_role} to {receiver_role} on topic '{topic}'."
        except sqlite3.OperationalError as e:
            if "locked" in str(e).lower():
                import time
                time.sleep(0.5 * (attempt + 1))
                continue
            return f"❌ SQLite Operational Error: {str(e)}"
        except Exception as e:
            return f"❌ Error publishing message: {str(e)}"

    return "❌ Error: Database remained locked after multiple retries."

@mcp.tool()
def read_messages(workspace_path: str, receiver_role: str, topic: str = "", unread_only: bool = True) -> str:
    """
    Read messages from the internal message bus for a specific role.
    """
    try:
        conn = get_db_connection(workspace_path)
        cursor = conn.cursor()
        query = "SELECT * FROM messages WHERE receiver_role = ?"
        params = [receiver_role]

        if topic:
            query += " AND topic = ?"
            params.append(topic)

        if unread_only:
            query += " AND is_read = 0"

        query += " ORDER BY created_at ASC"

        cursor.execute(query, params)
        rows = cursor.fetchall()

        if not rows:
            conn.close()
            return "No messages found."

        res = [f"📬 Messages for {receiver_role}:"]
        msg_ids = []
        for r in rows:
            msg_ids.append(r['id'])
            res.append(f"--- [ID: {r['id']} | Topic: {r['topic']} | From: {r['sender_role']} | Time: {r['created_at']}] ---")
            res.append(r['content'])
            res.append("")

        if unread_only and msg_ids:
            # Mark as read
            placeholders = ', '.join(['?'] * len(msg_ids))
            cursor.execute(f"UPDATE messages SET is_read = 1 WHERE id IN ({placeholders})", msg_ids)
            conn.commit()

        conn.close()
        return "\n".join(res)
    except Exception as e:
        return f"❌ Error reading messages: {str(e)}"

@mcp.tool()
def clear_topic(workspace_path: str, topic: str) -> str:
    """Clear all messages for a specific topic."""
    try:
        conn = get_db_connection(workspace_path)
        cursor = conn.cursor()
        cursor.execute("DELETE FROM messages WHERE topic = ?", (topic,))
        deleted = cursor.rowcount
        conn.commit()
        conn.close()
        return f"✅ Cleared {deleted} messages for topic '{topic}'."
    except Exception as e:
        return f"❌ Error clearing topic: {str(e)}"

@mcp.tool()
def enforce_socratic_gate(action_name: str, impact_description: str, options: List[str]) -> str:
    """
    Enforce a socratic gate to verify critical actions with the user.
    The agent MUST output the returned string to the user and wait for their response.
    """
    formatted_options = "\n".join([f"- [{i+1}] {opt}" for i, opt in enumerate(options)])

    prompt = f"SYSTEM_OVERRIDE_STOP: CRITICAL ACTION VERIFICATION REQUIRED\\n\\n"
    prompt += f"The following action requires user confirmation before proceeding:\\n"
    prompt += f"**Action**: {action_name}\\n"
    prompt += f"**Impact**: {impact_description}\\n\\n"
    prompt += f"Please select an option to proceed:\\n"
    prompt += f"{formatted_options}\\n\\n"
    prompt += f"(Agent Instruction: You MUST echo this exact message to the user and stop execution. Wait for the user to select an option before continuing.)"

    return prompt

@mcp.tool()
def delegate_to_subagent(workspace_path: str, target_role: str, task_description: str, context_files: List[str], timeout_mins: int = 10, run_background: bool = True, engine: str = "kilocode", model: str = "") -> str:
    """
    Delegate a task to a subagent running.
    If run_background is True, spawns the agent and returns immediately (parallel mode), sending logs to the bus.
    If run_background is False, waits for the agent to finish before returning (sequential mode).
    """
    try:
        # Publish an initial message indicating the subagent is starting
        publish_message(workspace_path, f"subagent_{target_role}", "system", "all", f"Spawning Subagent '{target_role}' for task: {task_description}")

        role_requirements = ""

        workflow_file = ""
        role_lower = target_role.lower()
        if "planner" in role_lower:
            workflow_file = ".agent/workflows/planner-architect.md"
        elif "coder" in role_lower:
            workflow_file = ".agent/workflows/coder-implementation.md"
        elif "reviewer" in role_lower:
            workflow_file = ".agent/workflows/reviewer-audit.md"
        elif "tester" in role_lower:
            workflow_file = ".agent/workflows/tester-verification.md"

        system_prompt = (
            f"You are the {target_role.upper()} in a Multi-Agent architecture. Your current task is: {task_description}.\n\n"
            "MANDATORY PROTOCOLS:\n"
            f"1. WORKFLOW ADHERENCE: You MUST follow the instructions in '[{workflow_file}](file:///{workspace_path}/{workflow_file})'. Read it before taking action.\n"
            "2. GLOBAL RULES: Follow all rules in '[GEMINI.md](file:///.agent/rules/GEMINI.md)'.\n"
            "3. STAR TOPOLOGY: Report status back to the PLANNER or the next role in the workflow via `publish_message` after every work block.\n"
            "4. NO DEADLOCKS: You MUST NOT exit or idle without sending an activation message to the next participant. If stuck, notify the Planner.\n"
            "5. CONCISE COMMUNICATION: Focus on technical execution. Only send ONE final message summarizing your work.\n"
            f"{role_requirements}"
            "Once task is finished, output a technical summary and wait for instructions."
        )

        # Prepare context by referring to files if specified
        files_str = ""
        if context_files:
            files_str = f"Please refer to these files: {', '.join(context_files)}. "

        worker_script = os.path.join(os.path.dirname(os.path.abspath(__file__)), "worker.py")
        python_exec = sys.executable if sys.executable else "python"
        # FIX: use task_description instead of undefined prompt, and pass --engine and --model
        cmd = f"{python_exec} {shlex.quote(worker_script)} --workspace {shlex.quote(workspace_path)} --role {shlex.quote(target_role)} --instruction {shlex.quote(system_prompt)} --task {shlex.quote(task_description)} --engine {shlex.quote(engine)}"
        if model:
            cmd += f" --model {shlex.quote(model)}"

        publish_message(workspace_path, f"subagent_{target_role}", "system", "all", f"Spawning Daemon: {cmd}")

        my_env = os.environ.copy()

        if run_background:
            # Spawn and detach completely
            process = subprocess.Popen(
                cmd,
                shell=True,
                cwd=workspace_path,
                env=my_env,
                stdout=subprocess.DEVNULL, # Worker logs directly to DB
                stderr=subprocess.DEVNULL,
                start_new_session=True # Detach from parent
            )
            return f"🚀 Subagent Daemon {target_role} spawned in BACKGROUND. (PID: {process.pid}). Listening continuously for messages. View logs in Dashboard."
        else:
            # Sequential mode (for compatibility if needed)
            process = subprocess.Popen(
                cmd,
                shell=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                bufsize=1,
                cwd=workspace_path,
                env=my_env
            )
            for line in iter(process.stdout.readline, ''):
                if line:
                    publish_message(workspace_path, f"subagent_{target_role}", target_role, "all", line.strip())
            process.wait()
            return f"✅ Subagent {target_role} execution completed."

    except Exception as e:
        return f"❌ Subagent execution failed: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')
