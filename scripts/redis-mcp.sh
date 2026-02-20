#!/bin/bash

# --- Resolve Script Directory ---
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
    DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
SCRIPT_DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# --- Colors & Logging ---
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

log_info()    { echo -e "${BLUE}[INFO]${NC} $1" >&2; }
log_error()   { echo -e "${RED}[ERROR]${NC} $1" >&2; }

# --- Load Config ---
env_file="$SCRIPT_DIR/.env"
if [ ! -f "$env_file" ]; then
    log_error ".env file not found at $env_file"
    exit 1
fi

log_info "Loading config from $env_file..."
set -a
source "$env_file"
set +a

if ! command -v uvx &> /dev/null; then
    log_error "uvx command not found. Please install uv: curl -LsSf https://astral.sh/uv/install.sh | sh"
    exit 1
fi

host="${REDIS_HOST:-127.0.0.1}"
port="${REDIS_PORT:-6379}"
user="${REDIS_USER}"
pass="${REDIS_PASSWORD}"
db="${REDIS_DB:-0}"

auth=""
if [ -n "$user" ] || [ -n "$pass" ]; then
    auth="${user}:${pass}@"
fi

url="redis://${auth}${host}:${port}/${db}"

log_info "Starting redis-mcp-server connected to ${url}..."
redis-mcp-server --url "$url"
