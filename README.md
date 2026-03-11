# Antigravity Kit

> AI Agent templates with Skills, Agents, and Workflows

## Quick Install

Clone the repository and link the `agk` script to your path for global access:

```bash
# 1. Clone the repository (if not already done)
git clone git@github.com:Dang-Hai-Tran/antigravity-kit.git ~/.antigravity/cache/antigravity-kit

# 2. Link the script to your bin
sudo ln -s ~/.antigravity/cache/antigravity-kit/scripts/agk /usr/local/bin/agk

# 3. Verify installation
agk
```

### ⚠️ Important Note on `.gitignore`

If you are using AI-powered editors like **Cursor** or **Windsurf**, adding the `.agents/` folder to your `.gitignore`
may prevent the IDE from indexing the workflows. This results in slash commands (like `/plan`, `/debug`) not appearing
in the chat suggestion dropdown.

**Recommended Solution:** To keep the `.agents/` folder local (not tracked by Git) while maintaining AI functionality:

1. Ensure `.agents/` is **NOT** in your project's `.gitignore`.
2. Instead, add it to your local exclude file: `.git/info/exclude`

## What's Included

| Component        | Description                                                               |
| ---------------- | ------------------------------------------------------------------------- |
| **Integrations** | Powerful MCP Servers (`mcp:database-inspector`, `mcp:skill-router`, etc.) |
| **Skills**       | 880+ Domain-specific knowledge modules                                    |
| **Workflows**    | 15 Slash command procedures                                               |

## Usage

### 1. Managing .agents Folder

Run these commands in any project root:

| Command                | Description                                    |
| ---------------------- | ---------------------------------------------- |
| `agk install`          | Install .agents folder (auto-excludes from git) |
| `agk install --force`  | Force install, overwrite existing              |
| `agk update`           | Update .agents folder to latest version         |
| `agk update --dry-run` | Preview changes without applying               |
| `agk update --offline` | Update using cached repo (no network)          |
| `agk status`           | Check if updates are available                 |
| `agk remove`           | Remove .agents folder (with confirmation)       |
| `agk remove --force`   | Remove without confirmation                    |
| `agk version`          | Show version and cache commit                  |
| `agk --help`           | Show full help with examples                   |

```bash
# Install .agents folder
agk install

# Force install (overwrite existing)
agk install --force

# Update to latest version
agk update

# Preview what would change
agk update --dry-run

# Update without internet (use cache)
agk update --offline

# Check for updates
agk status

# Remove .agents folder
agk remove
```

> **Note:** The `docs/` folder inside `.agents/` is preserved during updates.

### 2. Using Agents

**No need to mention agents explicitly!** The system automatically detects and applies the right specialist(s):

```
You: "Add JWT authentication"
AI: 🤖 Applying @security-auditor + @backend-specialist...

You: "Fix the dark mode button"
AI: 🤖 Using @frontend-specialist...

You: "Login returns 500 error"
AI: 🤖 Using @debugger for systematic analysis...
```

**How it works:**

- Analyzes your request silently

- Detects domain(s) automatically (frontend, backend, security, etc.)
- Selects the best specialist(s)
- Informs you which expertise is being applied
- You get specialist-level responses without needing to know the system architecture

**Benefits:**

- ✅ Zero learning curve - just describe what you need
- ✅ Always get expert responses
- ✅ Transparent - shows which agent is being used
- ✅ Can still override by mentioning agent explicitly

### 3. Using Workflows

Invoke workflows with slash commands:

| Command                 | Description                                           |
| ----------------------- | ----------------------------------------------------- |
| `/brainstorm`           | Explore options before implementation                 |
| `/create`               | Create new features or apps                           |
| `/db`                   | Act on databases with `@mcp:database-inspector`       |
| `/debug`                | Systematic debugging                                  |
| `/deploy`               | Deploy application                                    |
| `/design`               | UI/UX discovery, design system generation             |
| `/git-commit`           | Git stage, commit, and pre-rebase workflow            |
| `/git-push`             | Git pull rebase, conflict resolution, push            |
| `/gitlab-mr-read-fix`   | Read MR discussions via `@mcp:gitlab-mr-discussions`  |
| `/gitlab-mr-reply-push` | Push fixes and reply via `@mcp:gitlab-mr-discussions` |
| `/orchestrate`          | Multi-agent coordination for complex tasks            |
| `/plan`                 | Create comprehensive task breakdown                   |
| `/redis`                | Act on Redis datastores safely                        |
| `/test`                 | Generate and run tests                                |
| `/update`               | Update or enhance existing codebase                   |

Example:

```
/db explain the users table
/create landing page with hero section
/redis GET active_accounts
```

### 4. MCP Server Integration

The kit relies heavily on the **Model Context Protocol (MCP)** for secure, structured access to resources:

- `@mcp:skill-router`: Automatically searches and loads the perfect skills for your query.
- `@mcp:database-inspector`: Connects to Postgres, MySQL, and Redis to execute safe read-only queries.
- `@mcp:context-manager`: Persists progress checkpoints during long-running multi-file tasks.
- `@mcp:doc-researcher`: Fetches the latest bleeding-edge SOTA syntax before composing new features.

### 5. Using Skills

Skills are loaded automatically based on task context (using the `search_skills` MCP). The AI reads skill descriptions
and applies relevant knowledge precisely.

## License

MIT © datran
