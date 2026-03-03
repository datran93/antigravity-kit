#!/bin/bash
WORKSPACE=$(pwd)
echo "🛑 Stopping all background Agents and Dashboard running on project: $WORKSPACE..."
pkill -f "worker.py --workspace $WORKSPACE"
pkill -f "dashboard.py"
echo "🧹 Cleaning up resource locks..."
rm -f "$WORKSPACE/.agent_logs"/*.lock
rm -f "$WORKSPACE/.agent_logs"/*.db-shm
rm -f "$WORKSPACE/.agent_logs"/*.db-wal
echo "✅ Agents and Dashboard stopped."
