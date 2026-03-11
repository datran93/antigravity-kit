---
description: Run a command in the terminal
---
# 💻 Terminal Workflow

This workflow represents the standard procedure for safely running commands in the terminal using the `run_command` and `command_status` tools.

## 🚀 Execution Steps

1. **Interpret Command**:
   - Analyze the command you intend to execute. Ensure it is syntactically correct for `mac` OS and `zsh` shell.
   - Determine if the command is safe to execute automatically via `SafeToAutoRun`. Highly destructive operations MUST NOT be auto-run. Avoid commands that require interactive prompts when possible.

// turbo
2. **Execute Command**:
   - Call the `run_command` tool passing the exact `CommandLine`.
   - Set `Cwd` to the appropriate working directory.
   - Use `WaitMsBeforeAsync` appropriately depending on how long you expect the command to take before it backgrounds.

// turbo
3. **Status Monitoring & Input Delivery**:
   - If the command string returns a `CommandId` due to executing in the background, monitor its execution using the `command_status` tool.
   - If the interactive shell or command requires input, use the `send_command_input` tool, providing necessary input strings and `\n` characters to proceed.
   - Monitor the command until its status is "done".

4. **Verify Execution Outcome**:
   - Inspect the command output lines.
   - If the command failed or produced errors, self-correct the command and retry, or gracefully report the failure context to the User.
