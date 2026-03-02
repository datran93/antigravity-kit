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
        -m|--model) MODEL="$2"; shift ;;
        *) TASK="$1" ;;
    esac
    shift
done

# Pass model argument if provided
MODEL_ARG=""
if [ -n "$MODEL" ]; then
    MODEL_ARG="--model $MODEL"
fi

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
1. STATE MANAGEMENT: Maintain a single source of truth at '.agent_logs/STATE.md'. It must track Architecture, Task Status, and Current Assignments.
2. DELEGATION: You have one coder: 'coder'. Send extremely brief messages to it (e.g., 'Task #2 assigned. Check STATE.md'). Do NOT put task details in the message.
3. COMMUNICATION: NO intermediate updates. Update STATE.md whenever tasks are completed or bugs are reported.
4. EXPLORE FIRST: Do not hallucinate files or context.
5. NEVER ask the user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role planner --instruction "$PLANNER_INST" --task "$TASK" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/planner.log" 2>&1 &
sleep 10

echo "🚀 Summoning [Coder Agent]..."
CODER_INST="You are a strict CODER. Your role name is 'coder'.
YOUR JOB:
1. Wait for task notifications from 'planner' via 'read_messages'.
2. READ STATE: First, read '.agent_logs/STATE.md' to understand your task and the architecture.
3. IMPLEMENT: Finish the task completely. Write CLEAN, TESTABLE code (SOLID, DRY). Ensure logic is isolated and use dependency injection to make testing easy.
4. NOTIFY REVIEWER: Once done, briefly note your changes in STATE.md, then call 'publish_message' to notify 'reviewer'. Message should just say 'Task ready for review from coder'.
5. FAIL-FAST: If blocked, update STATE.md with the blocker and escalate to 'planner'.
NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role coder --instruction "$CODER_INST" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/coder.log" 2>&1 &
sleep 10

echo "🚀 Summoning [Reviewer Agent]..."
REVIEWER_INST="You are a strict REVIEWER/AUDITOR.
YOUR JOB:
1. Wait for notifications from 'coder' via 'read_messages'.
2. READ STATE: Check '.agent_logs/STATE.md' to see what they built.
3. AUDIT: Perform a complete audit of code quality and security on their changes.
4. DECISION: If approved, update STATE.md as 'Approved' and notify 'tester'. If rejected, log feedback in STATE.md and notify 'coder'.
Prevent infinite loops. NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role reviewer --instruction "$REVIEWER_INST" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/reviewer.log" 2>&1 &
sleep 10

echo "🚀 Summoning [Tester Agent]..."
TESTER_INST="You are a strict TESTER/QA.
YOUR JOB:
1. Wait for audit approval from 'reviewer' via 'read_messages'.
2. READ STATE: Check '.agent_logs/STATE.md' for the features to test.
3. TEST & COVERAGE: Run comprehensive tests. You MUST write unit tests to test ALL individual functions implemented by the coder. Aim for 100% logic coverage.
4. COMPLETION: If pass, mark as COMPLETED in STATE.md and notify 'planner'. If fail, write bug details to STATE.md and notify 'planner' so it can assign a fix.
NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role tester --instruction "$TESTER_INST" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/tester.log" 2>&1 &

echo "🚀 Starting Web Dashboard..."
export MULTI_AGENT_DB_PATH="$WORKSPACE_DB"
cd "$SCRIPT_DIR" || exit
nohup "$PYTHON_ENV" dashboard.py > "$LOG_DIR/dashboard.log" 2>&1 &
cd "$WORKSPACE" || exit

echo "---------------------------------------------------------"
echo "✅ ALL 4 AGENTS HAVE BEEN RESET AND STARTED IN: $WORKSPACE!"
echo "📊 Monitoring Dashboard: http://localhost:6060"
echo "---------------------------------------------------------"
