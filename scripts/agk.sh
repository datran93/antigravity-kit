#!/bin/bash

# =============================================================================
# Antigravity Kit (agk) - AI Agent Management Tool
# =============================================================================

VERSION="1.4.0"

# --- Resolve Script Directory ---
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
    DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
SCRIPT_DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# --- Configuration ---
CACHE_DIR="$HOME/.antigravity/cache"
REPO_URL="git@github.com:Dang-Hai-Tran/antigravity-kit.git"
REPO_NAME="antigravity-kit"
SOURCE_AGENT_DIR="$CACHE_DIR/$REPO_NAME/.agents"
TARGET_AGENT_DIR="./.agents"

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

# Add .agents to git exclude
add_git_exclude() {
    [ ! -d ".git" ] && return

    local exclude=".git/info/exclude"
    grep -q "^\.agents$" "$exclude" 2>/dev/null || echo ".agents" >> "$exclude"
}

# Remove .agents from git exclude
remove_git_exclude() {
    [ ! -d ".git" ] && return

    local exclude=".git/info/exclude"
    if [ -f "$exclude" ] && grep -q "^\.agents$" "$exclude"; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' '/^\.agents$/d' "$exclude"
        else
            sed -i '/^\.agents$/d' "$exclude"
        fi
    fi
}

# =============================================================================
# Commands
# =============================================================================

cmd_install() {
    if [ -d "$TARGET_AGENT_DIR" ]; then
        log_error ".agents already exists. Use 'agk update' or remove it first."
        exit 1
    fi

    sync_repo

    log_info "Installing .agents..."
    cp -R "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null

    add_git_exclude
    log_success ".agents installed!"
}

cmd_update() {
    sync_repo

    if [ ! -d "$TARGET_AGENT_DIR" ]; then
        log_info ".agents not found. Installing..."
        cp -R "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null
        add_git_exclude
        log_success ".agents installed!"
        return
    fi

    log_info "Updating .agents..."
    rm -rf "$TARGET_AGENT_DIR" 2>/dev/null
    cp -R "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null
    log_success ".agents updated!"
}

cmd_status() {
    sync_repo

    if [ ! -d "$TARGET_AGENT_DIR" ]; then
        log_error ".agents not found. Run 'agk install'."
        exit 1
    fi

    local diff_output
    diff_output=$(diff -r --brief "$SOURCE_AGENT_DIR" "$TARGET_AGENT_DIR" 2>/dev/null | grep -v "\.DS_Store")

    if [ -z "$diff_output" ]; then
        log_success ".agents is up to date."
    else
        log_info "Updates available. Run 'agk update'."
    fi
}

cmd_remove() {
    if [ ! -d "$TARGET_AGENT_DIR" ]; then
        log_error ".agents not found."
        exit 1
    fi

    echo -e "${RED}WARNING: This will delete .agents folder.${NC}"
    read -p "Continue? (y/N): " confirm
    [[ ! "$confirm" =~ ^[Yy]$ ]] && { log_info "Cancelled."; return; }

    rm -rf "$TARGET_AGENT_DIR" 2>/dev/null
    remove_git_exclude
    log_success ".agents removed!"
}


cmd_sync_skills() {
    local sync_script="$SCRIPT_DIR/sync-skills.sh"

    if [ ! -f "$sync_script" ]; then
        log_error "sync-skills.sh not found at $sync_script"
        exit 1
    fi

    log_info "Running sync-skills.sh..."
    bash "$sync_script" "$@"
}

cmd_sync_env() {
    local src_env="$(pwd)/.env"
    local dest_env="$SCRIPT_DIR/.env"

    if [ ! -f "$src_env" ]; then
        log_error ".env file not found in current directory: $src_env"
        exit 1
    fi

    log_info "Copying $src_env to $dest_env"
    mkdir -p "$(dirname "$dest_env")" 2>/dev/null
    cp "$src_env" "$dest_env"
    log_success "Environment configuration copied successfully!"
}

cmd_show_env() {
    local env_file="$SCRIPT_DIR/.env"
    if [ ! -f "$env_file" ]; then
        log_error "No .env file found at: $env_file. You need to run 'agk sync-env' first."
        exit 1
    fi

    log_info "Content of .env file:"
    echo "----------------------------------------"
    cat "$env_file"
    echo "----------------------------------------"
}

show_help() {
    cat << EOF
Antigravity Kit v$VERSION

Usage: agk <command>

Commands:
  install       Install .agents folder
  update        Update .agents to latest
  status        Check for updates
  remove        Remove .agents folder
  sync-skills   Sync local .agents with awesome-skills repo
  sync-env      Copy current directory .env to .agents/scripts/.env
  show-env      Print the current .env from .agents/scripts/.env
  help          Show this help

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
    sync-skills) shift; cmd_sync_skills "$@" ;;
    sync-env) cmd_sync_env ;;
    show-env) cmd_show_env ;;
    help|-h|--help|"") show_help ;;
    *) log_error "Unknown command: $1"; show_help; exit 1 ;;
esac
