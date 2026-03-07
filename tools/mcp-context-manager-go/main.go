package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"McpContextManager",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// save_checkpoint
	mcpServer.AddTool(mcp.NewTool("save_checkpoint",
		mcp.WithDescription("Save or update a task checkpoint/context."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Description of the task")),
		mcp.WithString("status", mcp.Required(), mcp.Description("Status")),
		mcp.WithArray("completed_steps", mcp.Items(map[string]interface{}{"type": "string"}), mcp.Description("Completed steps array")),
		mcp.WithArray("next_steps", mcp.Items(map[string]interface{}{"type": "string"}), mcp.Description("Next steps array")),
		mcp.WithArray("active_files", mcp.Items(map[string]interface{}{"type": "string"}), mcp.Description("Active files array")),
		mcp.WithString("notes", mcp.Required(), mcp.Description("Notes")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		description, _ := args["description"].(string)
		status, _ := args["status"].(string)

		compStepsIf, _ := args["completed_steps"].([]interface{})
		var compSteps []string
		for _, s := range compStepsIf {
			if str, ok := s.(string); ok {
				compSteps = append(compSteps, str)
			}
		}

		nextStepsIf, _ := args["next_steps"].([]interface{})
		var nextSteps []string
		for _, s := range nextStepsIf {
			if str, ok := s.(string); ok {
				nextSteps = append(nextSteps, str)
			}
		}

		activeFilesIf, _ := args["active_files"].([]interface{})
		var activeFiles []string
		for _, s := range activeFilesIf {
			if str, ok := s.(string); ok {
				activeFiles = append(activeFiles, str)
			}
		}

		notes, _ := args["notes"].(string)

		res, err := SaveCheckpoint(workspacePath, taskID, description, status, compSteps, nextSteps, activeFiles, notes)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// initialize_task_plan
	mcpServer.AddTool(mcp.NewTool("initialize_task_plan",
		mcp.WithDescription("Start a new task with a list of steps."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Description of the task")),
		mcp.WithArray("steps", mcp.Items(map[string]interface{}{"type": "string"}), mcp.Description("Steps collection")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		description, _ := args["description"].(string)

		stepsIf, _ := args["steps"].([]interface{})
		var steps []string
		for _, s := range stepsIf {
			if str, ok := s.(string); ok {
				steps = append(steps, str)
			}
		}

		res, err := InitializeTaskPlan(workspacePath, taskID, description, steps)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// complete_task_step
	mcpServer.AddTool(mcp.NewTool("complete_task_step",
		mcp.WithDescription("Mark step as done, track active files, update graph and bar."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("step_name", mcp.Required(), mcp.Description("Step name")),
		mcp.WithArray("active_files", mcp.Items(map[string]interface{}{"type": "string"}), mcp.Description("Active files array")),
		mcp.WithString("notes", mcp.Description("Notes")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		stepName, _ := args["step_name"].(string)

		var activeFiles []string
		if val, ok := args["active_files"]; ok && val != nil {
			if activeFilesIf, ok2 := val.([]interface{}); ok2 {
				for _, s := range activeFilesIf {
					if str, ok := s.(string); ok {
						activeFiles = append(activeFiles, str)
					}
				}
			}
		}

		var notes string
		if val, ok := args["notes"]; ok && val != nil {
			notes, _ = val.(string)
		}

		res, err := CompleteTaskStep(workspacePath, taskID, stepName, activeFiles, notes)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// add_task_step
	mcpServer.AddTool(mcp.NewTool("add_task_step",
		mcp.WithDescription("Add a new task step to the next_steps list of an existing task."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("new_step", mcp.Required(), mcp.Description("New step")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		newStep, _ := args["new_step"].(string)

		res, err := AddTaskStep(workspacePath, taskID, newStep)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// load_checkpoint
	mcpServer.AddTool(mcp.NewTool("load_checkpoint",
		mcp.WithDescription("Load a previously saved task checkpoint."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)

		res, err := LoadCheckpoint(workspacePath, taskID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// list_active_tasks
	mcpServer.AddTool(mcp.NewTool("list_active_tasks",
		mcp.WithDescription("List all active tasks."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)

		res, err := ListActiveTasks(workspacePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// declare_intent
	mcpServer.AddTool(mcp.NewTool("declare_intent",
		mcp.WithDescription("Declare working intention and restrict changes to specific locked_files."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("tactic", mcp.Required(), mcp.Description("Tactic name")),
		mcp.WithArray("locked_files", mcp.Items(map[string]interface{}{"type": "string"}), mcp.Description("Locked files")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		tactic, _ := args["tactic"].(string)

		lockedIf, _ := args["locked_files"].([]interface{})
		var locked []string
		for _, v := range lockedIf {
			if str, ok := v.(string); ok {
				locked = append(locked, str)
			}
		}

		res, err := DeclareIntent(workspacePath, taskID, tactic, locked)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// check_intent_lock
	mcpServer.AddTool(mcp.NewTool("check_intent_lock",
		mcp.WithDescription("Check if a file is unlocked for the current intent."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("File Path")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		filePath, _ := args["file_path"].(string)

		res, err := CheckIntentLock(workspacePath, taskID, filePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// annotate_file
	mcpServer.AddTool(mcp.NewTool("annotate_file",
		mcp.WithDescription("Add 'Ghost Context' (lessons/gotchas) to a specific file"),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("File Path")),
		mcp.WithString("gotchas", mcp.Required(), mcp.Description("Gotchas context")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		filePath, _ := args["file_path"].(string)
		gotchas, _ := args["gotchas"].(string)

		res, err := AnnotateFile(workspacePath, filePath, gotchas)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// record_failure
	mcpServer.AddTool(mcp.NewTool("record_failure",
		mcp.WithDescription("Record a failure (e.g. test failing, compile error) to detect context drift."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)

		res, err := RecordFailure(workspacePath, taskID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// clear_drift
	mcpServer.AddTool(mcp.NewTool("clear_drift",
		mcp.WithDescription("Reset the failure counter after a success or planner intervention."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)

		res, err := ClearDrift(workspacePath, taskID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// manage_anchors
	mcpServer.AddTool(mcp.NewTool("manage_anchors",
		mcp.WithDescription("Manage architectural anchors. Action can be 'set', 'get', or 'list'."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action (set, get, list)")),
		mcp.WithString("key", mcp.Description("Anchor key")),
		mcp.WithString("value", mcp.Description("Anchor value")),
		mcp.WithString("rule", mcp.Description("Anchor rule")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		action, _ := args["action"].(string)

		var key, value, rule string
		if val, ok := args["key"]; ok && val != nil {
			key, _ = val.(string)
		}
		if val, ok := args["value"]; ok && val != nil {
			value, _ = val.(string)
		}
		if val, ok := args["rule"]; ok && val != nil {
			rule, _ = val.(string)
		}

		res, err := ManageAnchors(workspacePath, action, key, value, rule)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// recall_knowledge
	mcpServer.AddTool(mcp.NewTool("recall_knowledge",
		mcp.WithDescription("Recall relevant Knowledge Items (Local RAG) based on a search query using SQLite FTS5."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithNumber("top_k", mcp.Description("Top K results")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		query, _ := args["query"].(string)

		topK := 3
		if val, ok := args["top_k"]; ok && val != nil {
			if num, ok2 := val.(float64); ok2 {
				topK = int(num)
			}
		}

		res, err := RecallKnowledge(workspacePath, query, topK)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// compact_memory
	mcpServer.AddTool(mcp.NewTool("compact_memory",
		mcp.WithDescription("Generate a Knowledge Item (KI) and prune context at the end of a tactic. Also indexes KI into Semantic SQLite."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("Workspace path")),
		mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID")),
		mcp.WithString("tactic_name", mcp.Required(), mcp.Description("Tactic name")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("KI summary")),
		mcp.WithString("decisions", mcp.Required(), mcp.Description("Architecture decisions")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		workspacePath, _ := args["workspace_path"].(string)
		taskID, _ := args["task_id"].(string)
		tacticName, _ := args["tactic_name"].(string)
		summary, _ := args["summary"].(string)
		decisions, _ := args["decisions"].(string)

		res, err := CompactMemory(workspacePath, taskID, tacticName, summary, decisions)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Run standard I/O server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
