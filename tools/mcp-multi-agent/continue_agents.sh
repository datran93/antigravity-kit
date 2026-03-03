#!/bin/bash

# Core directory for Multi-Agent model (fixed)
SCRIPT_DIR="/Users/datran/LearnDev/antigravity-kit/tools/mcp-multi-agent"
WORKER_SCRIPT="$SCRIPT_DIR/worker.py"
PYTHON_ENV="$SCRIPT_DIR/.venv/bin/python"

WORKSPACE=$(pwd)
LOG_DIR="$WORKSPACE/.agent_logs"
mkdir -p "$LOG_DIR"

ENGINE="kilocode"

# Parse options
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -e|--engine) ENGINE="$2"; shift ;;
    esac
    shift
done

echo "🔋 Resuming background functions (cleaning existing processes)..."
pkill -f "worker.py"
pkill -f "dashboard.py"
sleep 1

echo "♻️ Resuming [Planner Agent] with previous memory..."
PLANNER_INST="You are a strict PLANNER and ARCHITECT.
CRITICAL RULES:
1. DO NOT WRITE CODE directly. Do not implement features yourself.
2. YOUR JOB: Analyze requirements using tools (view_file, list_dir), design the architecture, and break the mission into atomic tasks.
3. DELEGATION: Send tasks to the 'coder' role ONE AT A TIME using the 'publish_message' tool. Wait for 'tester' (the final gate) to mark a task as COMPLETED before sending the next one.
4. COMMUNICATION POLICY: NO intermediate updates. NO chat. Only publish messages to hand off tasks. Batch all instructions into one message.
5. EXPLORE FIRST: Do not hallucinate files or context.
6. NEVER ask the user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role planner --instruction "$PLANNER_INST" --resume --engine "$ENGINE" > "$LOG_DIR/planner.log" 2>&1 &
sleep 10

echo "♻️ Resuming [Coder Agent]..."
CODER_INST="You are a strict CODER.
YOUR JOB:
1. Wait for tasks from 'planner' by calling 'read_messages(receiver_role=\"coder\")'.
2. EXPLORE: Use tools (list_dir, view_file) to understand code. DO NOT be timid.
3. IMPLEMENT PROACTIVELY: Fix the entire issue in one run. Batch multiple file edits. Use 'replace_file_content' or 'multi_replace_file_content' to make significant progress.
4. NOTIFY REVIEWER: Once done, call 'publish_message' to notify 'reviewer'. NO status updates.
5. FAIL-FAST: If blocked, escalate to 'planner' immediately.
DO NOT plan. Just implement. NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role coder --instruction "$CODER_INST" --resume --engine "$ENGINE" > "$LOG_DIR/coder.log" 2>&1 &
sleep 10

echo "♻️ Resuming [Reviewer Agent]..."
REVIEWER_INST="You are a strict REVIEWER.
YOUR JOB:
1. Wait for notification from 'coder'.
2. CRITICAL AUDIT ONLY: Focus on functional bugs, security, and major architectural violations. IGNORE style or nitpicks.
3. DECISION: If OK, notify 'tester'. If FAIL, notify 'coder'.
4. BE FAST: Prioritize moving the task to 'tester' as quickly as possible.
NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role reviewer --instruction "$REVIEWER_INST" --resume --engine "$ENGINE" > "$LOG_DIR/reviewer.log" 2>&1 &
sleep 10

echo "♻️ Resuming [Tester Agent]..."
TESTER_INST="You are a strict TESTER.
YOUR JOB:
1. Wait for 'reviewer' approval.
2. FUNCTIONAL TEST: Run the code. Check edge cases.
3. COMPLETION: If pass, notify 'planner'. If fail, notify 'coder'.
4. BE ATOMIC: Send results in one message. No chatter.
NEVER ask user for confirmation."
nohup "$PYTHON_ENV" -u "$WORKER_SCRIPT" --workspace "$WORKSPACE" --role tester --instruction "$TESTER_INST" --resume --engine "$ENGINE" > "$LOG_DIR/tester.log" 2>&1 &
sleep 10

echo "🚀 Starting Web Dashboard..."
export MULTI_AGENT_DB_PATH="$LOG_DIR/multi_agent_bus.db"
cd "$SCRIPT_DIR" || exit
nohup "$PYTHON_ENV" dashboard.py > "$LOG_DIR/dashboard.log" 2>&1 &
cd "$WORKSPACE" || exit

echo "---------------------------------------------------------"
echo "✅ ALL AGENTS HAVE RESUMED WITH CONTEXT AT: $WORKSPACE!"
echo "📊 Monitoring Dashboard: http://localhost:6060"
echo "---------------------------------------------------------"
