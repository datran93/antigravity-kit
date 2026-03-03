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
pkill -f "worker.py --workspace $WORKSPACE"
pkill -f "dashboard.py"
sleep 1

WORKSPACE_DB="$LOG_DIR/multi_agent_bus.db"
if [ -f "$WORKSPACE_DB" ]; then
    echo "🗑 Resetting Agent memory to start a NEW project..."
    rm "$WORKSPACE_DB"
fi

echo "🚀 Summoning [Planner Agent]..."
PLANNER_INST="You are the LEAD PLANNER and ARCHITECT.
MANDATORY PROTOCOLS:
1. WORKFLOW: Follow '.agent/workflows/planner-architect.md' strictly.
2. COORDITANION: You are the Dispatcher. Assign tasks to 'coder' via 'publish_message'.
3. NO DEADLOCK: Never idle without a message.
4. PLAN: Maintain the project plan using @mcp:context-manager."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role planner --instruction "$PLANNER_INST" --task "$TASK" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/planner.log" 2>&1 &
sleep 5

echo "🚀 Summoning [Coder Agent]..."
CODER_INST="You are the CODER.
MANDATORY PROTOCOLS:
1. WORKFLOW: Follow '.agent/workflows/coder-implementation.md' strictly.
2. REPORT: Always report status back to the Planner.
3. QUALITY: Write Clean, Testable code."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role coder --instruction "$CODER_INST" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/coder.log" 2>&1 &
sleep 5

echo "🚀 Summoning [Reviewer Agent]..."
REVIEWER_INST="You are the REVIEWER.
MANDATORY PROTOCOLS:
1. WORKFLOW: Follow '.agent/workflows/reviewer-audit.md' strictly.
2. READ-ONLY: Never modify source code. Request fixes from the coder."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role reviewer --instruction "$REVIEWER_INST" --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/reviewer.log" 2>&1 &
sleep 5

echo "🚀 Summoning [Tester Agent]..."
TESTER_INST="You are the TESTER.
MANDATORY PROTOCOLS:
1. WORKFLOW: Follow '.agent/workflows/tester-verification.md' strictly.
2. TEST: Write and run tests for all code logic."
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
