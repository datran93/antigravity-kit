#!/bin/bash

# Core directory for Multi-Agent model (fixed)
SCRIPT_DIR="/Users/datran/LearnDev/antigravity-kit/tools/mcp-multi-agent"
WORKER_SCRIPT="$SCRIPT_DIR/worker.py"
PYTHON_ENV="$SCRIPT_DIR/.venv/bin/python"

WORKSPACE=$(pwd)
LOG_DIR="$WORKSPACE/.agent_logs"
mkdir -p "$LOG_DIR"

ENGINE="opencode"
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

echo "🔋 Resuming background functions (cleaning existing processes)..."
pkill -f "worker.py --workspace $WORKSPACE"
pkill -f "dashboard.py"
echo "📝 Resetting stale Agent statuses..."
DB_PATH="$LOG_DIR/multi_agent_bus.db"
if [ -f "$DB_PATH" ]; then
    sqlite3 "$DB_PATH" "UPDATE agent_status SET status = 'OFFLINE', current_task = '' ;"
fi
sleep 1

echo "♻️ Resuming [Lead Planner Agent] (Context Pre-loaded)..."
PLANNER_INST="You are the LEAD PLANNER and ORCHESTRATOR.
MANDATORY PROTOCOLS (SEQUENTIAL MODE):
1. WORKFLOW: Follow '.agent/workflows/planner-architect.md' strictly.
2. EPHEMERAL DELEGATION: You are the ONLY persistent agent. When you need a task done, you MUST call 'delegate_to_subagent' with run_background=False (Sequential Mode).
3. DISPATCH & DESTROY: Call subagents (coder, reviewer, or tester) for atomic tasks. Once they return their technical summary, analyze it, and then delegate the next task to a NEW subagent.
4. NO PERSISTENT TEAM: Do not expect coder/reviewer/tester to be running. You summon them only when needed via the tool.
5. REPORT: Maintain the project plan using @mcp:context-manager."

# Resume Planner as a persistent daemon with memory context and optional new task instruction
# We use --resume to signal worker.py to pull history from DB
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role planner --instruction "$PLANNER_INST" --task "$TASK" --resume --engine "$ENGINE" $MODEL_ARG > "$LOG_DIR/planner.log" 2>&1 &
sleep 2

echo "🚀 Starting Web Dashboard..."
export MULTI_AGENT_DB_PATH="$LOG_DIR/multi_agent_bus.db"
cd "$SCRIPT_DIR" || exit
nohup "$PYTHON_ENV" dashboard.py > "$LOG_DIR/dashboard.log" 2>&1 &
cd "$WORKSPACE" || exit

echo "---------------------------------------------------------"
echo "✅ SEQUENTIAL PIPELINE RESUMED AT: $WORKSPACE!"
echo "📍 Previous context is ACTIVE. New instruction (if any) sent to Planner."
echo "📊 Monitoring Dashboard: http://localhost:6060"
echo "---------------------------------------------------------"
