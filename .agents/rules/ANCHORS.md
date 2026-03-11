# ⚓ ANCHORS (Immutable Facts & Guardrails)

This file contains the absolute, non-negotiable project facts, architectural constraints, and guardrails for the agent
system. These **Anchors** must survive context compaction and session restarts.

As a Self-Executing Agent, you MUST respect these facts before executing any code changes or planning any project
tactics.

## �️ Technology Stack

- **Backend Stack**: Golang >= 1.25.
- **Frontend Stack**: React >= 19 / Next.js >= 16 App Router.
- **Package Manager**: Use `pnpm` instead of `npm` for better performance.
- **Database**: PostgreSQL 15+. DO NOT use MongoDB or MySQL.

## 🛡️ Agentic Guardrails & Execution Constraints

- **State Boundaries**: `active_files` MUST be explicitly tracked across role transitions and locked via Intent
  Declaration to prevent blind writes and scope creep.
- **No Destruction**: Do not delete existing API contracts, database columns, or core functionality without explicit
  confirmation from the USER.
- **Panic Protocol**: If you fail at the same step 3 times (e.g., test failures looping between Coder and Tester), you
  MUST STOP, record failure (`record_failure`), and drop back to the Planner role to discuss with the USER instead of
  blindly hallucinating fixes.
- **No Assumptions**: If Requirements are vague, stop and clarify via Socratic questioning before writing any code.

## 🧪 Quality & Verification

- **TDD Requirement**: Every action in a Task Plan MUST include a clear Verification Command as acceptance criteria.
  Code is not complete until verification passes.
- **Language**: All source code, variables, functions, comments, and commit messages MUST be in English.
- **Coverage**: You must ensure high quality and stability (aim for >= 70% test coverage) before calling
  `complete_task_step`.

## 🧠 Memory Governance

- **Session Compaction**: At the end of a major Tactic/Phase, run the `compact_session` workflow to distill
  architectural decisions and lessons learned into long-term Knowledge Items (KIs).
- **Ghost Context**: When encountering complex file-specific logic or tricky quirks, leverage `annotate_file` to attach
  localized lessons directly to the file to prevent recurring mistakes.

## 🛑 Restricted Terminal Commands

To prevent destructive actions, the agent MUST set `SafeToAutoRun: false` and explicitly ask the USER for confirmation before executing any of the following commands:

- **File & Directory Destruction**: `rm` (especially `rm -rf` or `rm -f`), `mv` (if overwriting critical files), `shred`
- **Destructive Git Operations**: `git push -f`, `git reset --hard`, `git clean -fd`, `git rebase` (on shared branches), `git branch -D`
- **Database & Infrastructure Alterations**: SQL `drop`/`truncate`/`delete`, `terraform destroy`/`apply`, `aws * delete-*`/`terminate-*`, `kubectl delete`, `docker system prune`/`docker rm -f`
- **System Permissions & Security**: `chmod` (e.g., recursive changes), `chown`, `chgrp`, `sudo`
- **Publishing & Deployments**: `npm publish`, `docker push`, package registry uploads, direct production deployments (`vercel deploy --prod`)
- **Process Management**: `kill`, `killall`, `pkill`, `systemctl stop`/`restart`
