#!/bin/bash

# --- CONFIGURATION ---
# Set this to the folder containing your project .env file the script's directory.
PROJECT_DIR="/Users/datran/LearnDev/antigravity-kit"
# ---------------------

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ENV_FILE=""

# Smart .env discovery:
# 1. Custom Project Directory
if [ -n "$PROJECT_DIR" ] && [ -f "$PROJECT_DIR/.env" ]; then
    ENV_FILE="$PROJECT_DIR/.env"
# 3. Fallback to script's own directory
elif [ -f "$SCRIPT_DIR/.env" ]; then
    ENV_FILE="$SCRIPT_DIR/.env"
fi

if [ -n "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

# Execute toolbox with original arguments
exec /opt/homebrew/bin/toolbox "$@"
