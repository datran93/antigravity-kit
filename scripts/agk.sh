#!/bin/bash

# =============================================================================
# Antigravity Kit (agk) - AI Agent Management Tool
# =============================================================================

VERSION="1.4.0"

# Configuration
CACHE_DIR="$HOME/.antigravity/cache"
REPO_URL="git@github.com:Dang-Hai-Tran/antigravity-kit.git"
REPO_NAME="antigravity-kit"
SOURCE_AGENT_DIR="$CACHE_DIR/$REPO_NAME/.agent"
TARGET_AGENT_DIR="./.agent"
TARGET_DOCS_DIR="./agent-docs"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

# =============================================================================
# Logging
# =============================================================================

log_info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $1" >&2; }
log_warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }

# =============================================================================
# Core Functions
# =============================================================================

# Sync repository: compare HEAD, if different -> delete and re-clone
sync_repo() {
    mkdir -p "$CACHE_DIR" 2>/dev/null

    if [ -d "$CACHE_DIR/$REPO_NAME/.git" ]; then
        log_info "Checking for updates..."

        local local_head remote_head
        local_head=$(cd "$CACHE_DIR/$REPO_NAME" && git rev-parse HEAD 2>/dev/null)
        remote_head=$(git ls-remote "$REPO_URL" HEAD 2>/dev/null | cut -f1)

        if [ -z "$remote_head" ]; then
            log_warn "Cannot reach remote. Using cached version."
            return 0
        fi

        if [ "$local_head" = "$remote_head" ]; then
            log_success "Cache is up to date."
            return 0
        fi

        log_info "Updates available. Re-cloning..."
        rm -rf "$CACHE_DIR/$REPO_NAME" 2>/dev/null
    fi

    log_info "Cloning repository..."
    if ! git clone --depth 1 "$REPO_URL" "$CACHE_DIR/$REPO_NAME" 2>/dev/null; then
        log_error "Failed to clone repository."
        exit 1
    fi

    log_success "Repository synced."
}

# Add .agent and agent-docs to git exclude
add_git_exclude() {
    [ ! -d ".git" ] && return

    local exclude=".git/info/exclude"
    grep -q "^\.agent$" "$exclude" 2>/dev/null || echo ".agent" >> "$exclude"
    grep -q "^agent-docs$" "$exclude" 2>/dev/null || echo "agent-docs" >> "$exclude"
}

# Remove .agent from git exclude
remove_git_exclude() {
    [ ! -d ".git" ] && return

    local exclude=".git/info/exclude"
    if [ -f "$exclude" ] && grep -q "^\.agent$" "$exclude"; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' '/^\.agent$/d' "$exclude"
        else
            sed -i '/^\.agent$/d' "$exclude"
        fi
    fi
}

# =============================================================================
# Commands
# =============================================================================

cmd_install() {
    if [ -d "$TARGET_AGENT_DIR" ]; then
        log_error ".agent already exists. Use 'agk update' or remove it first."
        exit 1
    fi

    sync_repo

    log_info "Installing .agent..."
    cp -R "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null

    [ ! -d "$TARGET_DOCS_DIR" ] && mkdir -p "$TARGET_DOCS_DIR" 2>/dev/null

    add_git_exclude
    log_success ".agent installed!"
}

cmd_update() {
    sync_repo

    if [ ! -d "$TARGET_AGENT_DIR" ]; then
        log_info ".agent not found. Installing..."
        cp -R "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null
        [ ! -d "$TARGET_DOCS_DIR" ] && mkdir -p "$TARGET_DOCS_DIR" 2>/dev/null
        add_git_exclude
        log_success ".agent installed!"
        return
    fi

    log_info "Updating .agent..."
    [ ! -d "$TARGET_DOCS_DIR" ] && mkdir -p "$TARGET_DOCS_DIR" 2>/dev/null
    rm -rf "$TARGET_AGENT_DIR" 2>/dev/null
    cp -R "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null
    log_success ".agent updated!"
}

cmd_status() {
    sync_repo

    if [ ! -d "$TARGET_AGENT_DIR" ]; then
        log_error ".agent not found. Run 'agk install'."
        exit 1
    fi

    local diff_output
    diff_output=$(diff -r --brief "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null | grep -v "\.DS_Store")

    if [ -z "$diff_output" ]; then
        log_success ".agent is up to date."
    else
        log_info "Updates available. Run 'agk update'."
    fi

    if [ -d "$TARGET_DOCS_DIR" ]; then
        local count
        count=$(find "$TARGET_DOCS_DIR" -name "*.md" -type f 2>/dev/null | wc -l | tr -d ' ')
        log_info "agent-docs: $count document(s)"
    fi
}

cmd_remove() {
    if [ ! -d "$TARGET_AGENT_DIR" ]; then
        log_error ".agent not found."
        exit 1
    fi

    echo -e "${RED}WARNING: This will delete .agent folder.${NC}"
    echo -e "${YELLOW}Note: agent-docs will be preserved.${NC}"
    read -p "Continue? (y/N): " confirm
    [[ ! "$confirm" =~ ^[Yy]$ ]] && { log_info "Cancelled."; return; }

    rm -rf "$TARGET_AGENT_DIR" 2>/dev/null
    remove_git_exclude
    log_success ".agent removed!"
}

cmd_toolbox() {
    local current_dir=$(pwd)
    local source="${BASH_SOURCE[0]}"

    # Resolve symlinks
    while [ -h "$source" ]; do
        local dir="$( cd -P "$( dirname "$source" )" && pwd )"
        source="$(readlink "$source")"
        [[ $source != /* ]] && source="$dir/$source"
    done
    local script_dir="$( cd -P "$( dirname "$source" )" && pwd )"
    local tool_script="$script_dir/toolbox-mcp.sh"

    if [ ! -f "$tool_script" ]; then
        log_error "toolbox-mcp.sh not found at $tool_script"
        exit 1
    fi

    log_info "Updating toolbox project directory to: $current_dir"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|^PROJECT_DIR=.*|PROJECT_DIR=\"$current_dir\"|" "$tool_script"
    else
        sed -i "s|^PROJECT_DIR=.*|PROJECT_DIR=\"$current_dir\"|" "$tool_script"
    fi

    log_success "Updated $tool_script"
}

show_help() {
    cat << EOF
Antigravity Kit v$VERSION

Usage: agk <command>

Commands:
  install   Install .agent folder
  update    Update .agent to latest
  status    Check for updates
  remove    Remove .agent folder
  toolbox   Set current directory as PROJECT_DIR for toolbox-mcp
  help      Show this help

EOF
}

# =============================================================================
# Main
# =============================================================================

case "$1" in
    install) cmd_install ;;
    update)  cmd_update ;;
    status)  cmd_status ;;
    remove)  cmd_remove ;;
    toolbox) cmd_toolbox ;;
    help|-h|--help|"") show_help ;;
    *) log_error "Unknown command: $1"; show_help; exit 1 ;;
esac
