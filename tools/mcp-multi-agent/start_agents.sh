#!/bin/bash

# Core directory for Multi-Agent model (fixed)
SCRIPT_DIR="/Users/datran/LearnDev/antigravity-kit/tools/mcp-multi-agent"
WORKER_SCRIPT="$SCRIPT_DIR/worker.py"
PYTHON_ENV="$SCRIPT_DIR/.venv/bin/python"
DB_PATH="$SCRIPT_DIR/multi_agent_bus.db"

# Current workspace
WORKSPACE=$(pwd)
LOG_DIR="$WORKSPACE/.agent_logs"

mkdir -p "$LOG_DIR"

ENGINE="kilocode"
TASK=""

# Parse options
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -e|--engine) ENGINE="$2"; shift ;;
        *) TASK="$1" ;;
    esac
    shift
done

if [ -z "$TASK" ]; then
  echo "⚠️ Error: No request provided for Agent."
  echo "💡 Usage: $0 [-e|--engine <engine>] \"<Your Request>\""
  exit 1
fi

echo "🧹 Cleaning up existing Agents and Dashboard in this workspace (if any)..."
pkill -f "worker.py"
pkill -f "dashboard.py"
sleep 1

WORKSPACE_DB="$LOG_DIR/multi_agent_bus.db"
if [ -f "$WORKSPACE_DB" ]; then
    echo "🗑 Resetting Agent memory to start a NEW project..."
    rm "$WORKSPACE_DB"
fi

echo "🚀 Summoning [Planner Agent] for project: $WORKSPACE..."
PLANNER_INST="You are a strict PLANNER and ARCHITECT.
CRITICAL RULES:
1. DO NOT WRITE CODE directly. Do not implement features yourself.
2. YOUR JOB: Analyze requirements using tools (view_file, list_dir), design the architecture, and break the mission into atomic tasks.
3. DELEGATION: Send tasks to the 'coder' role ONE AT A TIME using the 'publish_message' tool. Wait for 'tester' (the final gate) to mark a task as COMPLETED before sending the next one.
4. COMMUNICATION POLICY: NO intermediate updates. NO chat. Only publish messages to hand off tasks. Batch all instructions into one message.
5. EXPLORE FIRST: Do not hallucinate files or context.
6. NEVER ask the user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role planner --instruction "$PLANNER_INST" --task "$TASK" --engine "$ENGINE" > "$LOG_DIR/planner.log" 2>&1 &
sleep 10

echo "🚀 Summoning [Coder Agent]..."
CODER_INST="You are a strict CODER.
YOUR JOB:
1. Wait for tasks from 'planner' by calling 'read_messages(receiver_role=\"coder\")'.
2. EXPLORE: Use tools (list_dir, view_file) to understand existing code before writing. Do not hallucinate.
3. IMPLEMENT FULLY: Do not send messages during code execution. Finish the entire task or at least a significant milestone before talking.
4. NOTIFY REVIEWER: Once a large feature block is written and self-verified, call 'publish_message' to notify 'reviewer' for code audit. NO status updates like 'starting work' or 'file X updated'.
5. FAIL-FAST: If blocked, summarize the exact issue and escalate to 'planner'.
DO NOT plan the architecture. Just implement. NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role coder --instruction "$CODER_INST" --engine "$ENGINE" > "$LOG_DIR/coder.log" 2>&1 &
sleep 10

echo "🚀 Summoning [Reviewer Agent]..."
REVIEWER_INST="You are a strict REVIEWER/AUDITOR.
YOUR JOB:
1. Wait for implementation notification from 'coder' by calling 'read_messages(receiver_role=\"reviewer\")'.
2. FINAL AUDIT: Perform a complete audit of code quality, security, and standards.
3. DECISION: If approved, notify 'tester' to begin functional verification. If rejected, notify 'coder' with concrete feedback.
4. FEEDBACK: Only provide blocked feedback. Do not nitpick or start discussions.
Prevent infinite loops. NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role reviewer --instruction "$REVIEWER_INST" --engine "$ENGINE" > "$LOG_DIR/reviewer.log" 2>&1 &
sleep 10

echo "🚀 Summoning [Tester Agent]..."
TESTER_INST="You are a strict TESTER/QA.
YOUR JOB:
1. Wait for audit approval from 'reviewer' by calling 'read_messages(receiver_role=\"tester\")'.
2. TEST FULL SUITE: Run comprehensive functional tests and check edge cases.
3. SINGLE REPORT: Send ONE message with all test results.
4. COMPLETION: If tests pass, notify 'planner' that the task is COMPLETED. If they fail, notify 'coder' with a full bug list.
NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role tester --instruction "$TESTER_INST" --engine "$ENGINE" > "$LOG_DIR/tester.log" 2>&1 &

echo "🚀 Starting Web Dashboard..."
export MULTI_AGENT_DB_PATH="$WORKSPACE_DB"
cd "$SCRIPT_DIR" || exit
nohup "$PYTHON_ENV" dashboard.py > "$LOG_DIR/dashboard.log" 2>&1 &
cd "$WORKSPACE" || exit

echo "---------------------------------------------------------"
echo "✅ ALL 4 AGENTS HAVE BEEN RESET AND STARTED IN: $WORKSPACE!"
echo "📊 Monitoring Dashboard: http://localhost:6060"
echo "---------------------------------------------------------"
