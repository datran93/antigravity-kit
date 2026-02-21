#!/bin/bash

# sync-skills.sh
# Synchronizes local .agent with the antigravity-awesome-skills repository.

# --- Configuration ---
AGENT_DIR="/Users/datran/LearnDev/antigravity-kit/.agent"
REPO_DIR="/Users/datran/LearnDev/antigravity-awesome-skills"

# --- 1. Pull the latest content ---
echo "🔄 Updating antigravity-awesome-skills repository..."
cd "$REPO_DIR" || { echo "❌ Error: Could not change directory to $REPO_DIR"; exit 1; }
git pull || { echo "⚠️ Warning: git pull failed, attempting to continue anyway..."; }

# Return to original directory
cd - > /dev/null || exit 1

# --- Validation ---
if [ ! -d "$AGENT_DIR" ]; then
  echo "❌ Error: Local .agent directory not found at $AGENT_DIR"
  exit 1
fi

# --- 2. Copy and replace skills directory ---
echo "🔄 Copying skills directory to $AGENT_DIR..."

# Replace skills directory
if [ -d "$REPO_DIR/skills" ]; then
  # Remove existing skills directory if it exists
  if [ -d "$AGENT_DIR/skills" ]; then
    rm -rf "$AGENT_DIR/skills"
  fi

  # Copy the new skills directory
  cp -R "$REPO_DIR/skills" "$AGENT_DIR/skills"
  echo "✅ Copied skills directory"
else
  echo "❌ Error: $REPO_DIR/skills not found."
  exit 1
fi

echo "🎉 Sync complete!"
