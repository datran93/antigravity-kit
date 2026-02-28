import os
import sqlite3
import json
from datetime import datetime
from mcp.server.fastmcp import FastMCP
from typing import List, Optional

mcp = FastMCP("McpContextManager")

def get_db_connection(workspace_path: str):
    if not workspace_path:
        raise ValueError("workspace_path is required")

    # Force absolute paths
    workspace_path = os.path.abspath(workspace_path)

    db_path = os.path.join(workspace_path, "context.db")

    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS checkpoints (
            task_id TEXT PRIMARY KEY,
            description TEXT,
            status TEXT,
            completed_steps TEXT,
            next_steps TEXT,
            active_files TEXT,
            notes TEXT,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ''')
    conn.commit()
    return conn

def write_markdown_progress(workspace_path, task_id, description, status, completed_steps, next_steps, active_files, notes):
    try:
        md_path = os.path.join(workspace_path, "progress.md")

        total_steps = len(completed_steps) + len(next_steps)
        progress_pct = (len(completed_steps) / total_steps * 100) if total_steps > 0 else 0

        # Progress bar
        bar_len = 20
        filled = int(bar_len * progress_pct // 100)
        bar = "█" * filled + "░" * (bar_len - filled)

        content = f"# 🚀 Project Progress: {task_id}\n\n"
        content += f"**Status:** `{status.upper()}` | **Progress:** `[{bar}] {progress_pct:.1f}%` ({len(completed_steps)}/{total_steps})\n\n"
        content += f"> {description}\n\n"

        # Mermaid Graph
        if total_steps > 0:
            content += "## 📊 Task Dependency Graph\n"
            content += "```mermaid\ngraph TD\n"
            content += f"  Start((Start)) --> Task1\n"
            all_s = completed_steps + next_steps
            for i, s in enumerate(all_s):
                node_id = f"Task{i+1}"
                style = "fill:#4caf50,stroke:#2e7d32,color:#fff" if s in completed_steps else "fill:#fbc02d,stroke:#f9a825,color:#000"
                sanitized_name = s.replace('"', "'")
                content += f'  {node_id}["{sanitized_name}"]\n'
                content += f'  style {node_id} {style}\n'
                if i < len(all_s) - 1:
                    content += f"  {node_id} --> Task{i+2}\n"
                else:
                    content += f"  {node_id} --> End((End))\n"
            content += "```\n\n"

        if active_files:
            content += "### 📁 Active Files\n"
            for f in active_files: content += f"- `{f}`\n"
            content += "\n"

        content += "### ✅ Completed\n"
        if not completed_steps: content += "_None yet_\n"
        else:
            for s in completed_steps: content += f"- [x] {s}\n"
        content += "\n"

        content += "### ⏳ Next Steps\n"
        if not next_steps: content += "_All tasks done!_ 🎉\n"
        else:
            for s in next_steps: content += f"- [ ] {s}\n"
        content += "\n"

        if notes:
            content += "### 📝 Log & Notes\n"
            content += f"```text\n{notes}\n```\n\n"

        content += f"---\n*Last sync: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}*"

        with open(md_path, 'w', encoding='utf-8') as f:
            f.write(content)
    except Exception as e:
        print(f"Error writing progress.md: {e}")

@mcp.tool()
def save_checkpoint(
    workspace_path: str,
    task_id: str,
    description: str,
    status: str,
    completed_steps: List[str],
    next_steps: List[str],
    active_files: List[str],
    notes: str
) -> str:
    """
    Save or update a task checkpoint/context.
    """
    try:
        conn = get_db_connection(workspace_path)
        cursor = conn.cursor()
        now = datetime.now().isoformat()

        cursor.execute('''
            INSERT INTO checkpoints (task_id, description, status, completed_steps, next_steps, active_files, notes, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            ON CONFLICT(task_id) DO UPDATE SET
                description=excluded.description,
                status=excluded.status,
                completed_steps=excluded.completed_steps,
                next_steps=excluded.next_steps,
                active_files=excluded.active_files,
                notes=excluded.notes,
                updated_at=excluded.updated_at
        ''', (task_id, description, status, json.dumps(completed_steps), json.dumps(next_steps), json.dumps(active_files), notes, now))

        conn.commit()
        conn.close()

        write_markdown_progress(workspace_path, task_id, description, status, completed_steps, next_steps, active_files, notes)

        msg = f"✅ Checkpoint '{task_id}' saved."
        if not next_steps and completed_steps:
            msg += "\n\n🎉 ALL TASKS COMPLETED! Great job."
        return msg
    except Exception as e:
        return f"❌ Error saving checkpoint: {str(e)}"

@mcp.tool()
def initialize_task_plan(workspace_path: str, task_id: str, description: str, steps: List[str]) -> str:
    """Start a new task with a list of steps."""
    init_notes = f"[{datetime.now().strftime('%H:%M:%S')}] Task started."
    return save_checkpoint(
        workspace_path=workspace_path,
        task_id=task_id,
        description=description,
        status="in_progress",
        completed_steps=[],
        next_steps=steps,
        active_files=[],
        notes=init_notes
    )

@mcp.tool()
def complete_task_step(workspace_path: str, task_id: str, step_name: str, notes: Optional[str] = None) -> str:
    """Mark step as done, update graph and bar."""
    try:
        conn = get_db_connection(workspace_path)
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM checkpoints WHERE task_id = ?", (task_id,))
        row = cursor.fetchone()
        if not row:
            conn.close()
            return f"❌ Task '{task_id}' not found."

        comp = json.loads(row['completed_steps'])
        nxt = json.loads(row['next_steps'])

        if step_name in nxt:
            nxt.remove(step_name)
            comp.append(step_name)
        else:
            conn.close()
            return f"⚠️ Step '{step_name}' not in queue."

        stat = "completed" if not nxt else row['status']
        log = row['notes'] + f"\n[{datetime.now().strftime('%H:%M:%S')}] Done: {step_name}"
        if notes: log += f" - {notes}"

        conn.close()
        return save_checkpoint(
            workspace_path=workspace_path,
            task_id=task_id,
            description=row['description'],
            status=stat,
            completed_steps=comp,
            next_steps=nxt,
            active_files=json.loads(row['active_files']),
            notes=log
        )
    except Exception as e: return f"❌ Error: {str(e)}"

@mcp.tool()
def load_checkpoint(workspace_path: str, task_id: str) -> str:
    """
    Load a previously saved task checkpoint.
    """
    try:
        conn = get_db_connection(workspace_path)
        cursor = conn.cursor()
        cursor.execute("SELECT * FROM checkpoints WHERE task_id = ?", (task_id,))
        row = cursor.fetchone()
        conn.close()

        if not row:
            return f"❌ Checkpoint '{task_id}' not found."

        res = [f"🔄 {row['task_id']} [{row['status'].upper()}]", f"Last Update: {row['updated_at']}\n"]

        comp = json.loads(row['completed_steps'])
        nxt = json.loads(row['next_steps'])
        total = len(comp) + len(nxt)
        pct = (len(comp)/total*100) if total > 0 else 0
        res.append(f"Progress: {pct:.1f}% ({len(comp)}/{total} steps)\n")

        res.append("## ✅ Completed")
        for s in comp: res.append(f"- [x] {s}")

        res.append("\n## ⏳ Next")
        for s in nxt: res.append(f"- [ ] {s}")

        res.append(f"\n## 📝 Notes\n{row['notes']}")
        return "\n".join(res)
    except Exception as e:
        return f"❌ Error loading: {str(e)}"

@mcp.tool()
def list_active_tasks(workspace_path: str) -> str:
    """
    List all active tasks.
    """
    try:
        conn = get_db_connection(workspace_path)
        cursor = conn.cursor()
        cursor.execute("SELECT task_id, status, updated_at FROM checkpoints ORDER BY updated_at DESC")
        rows = cursor.fetchall()
        conn.close()

        if not rows: return "No tasks found."

        res = ["📋 WORKSPACE TASKS:"]
        for r in rows:
            res.append(f"- **{r['task_id']}** ({r['status']}) - {r['updated_at']}")
        return "\n".join(res)
    except Exception as e:
        return f"❌ Error: {str(e)}"

if __name__ == "__main__":
    mcp.run(transport='stdio')
