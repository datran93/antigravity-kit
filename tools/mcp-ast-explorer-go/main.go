package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"McpAstExplorer",
		"1.0.0",
		server.WithLogging(),
	)

	// Tool: get_project_architecture
	getProjectArchTool := mcp.NewTool("get_project_architecture",
		mcp.WithDescription("Get a structural overview (AST-based) of the project architecture.\nExtracts Classes, Functions, and Methods with signatures."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("The base absolute path of the project workspace.")),
		mcp.WithString("sub_path", mcp.Description("Optional sub-path within the workspace to limit scope.")),
		mcp.WithNumber("max_files", mcp.Description("Max number of files to process, default 1000.")),
		mcp.WithBoolean("include_docs", mcp.Description("If true, includes the first line of docstrings/comments.")),
	)

	s.AddTool(getProjectArchTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		workspacePath, ok := args["workspace_path"].(string)
		if !ok {
			return mcp.NewToolResultError("workspace_path is required"), nil
		}

		subPath, _ := args["sub_path"].(string)

		maxFiles := 1000
		if maxFilesFloat, ok := args["max_files"].(float64); ok {
			maxFiles = int(maxFilesFloat)
		}

		includeDocs := false
		if incDocsBool, ok := args["include_docs"].(bool); ok {
			includeDocs = incDocsBool
		}

		res, err := GetProjectArchitecture(workspacePath, subPath, maxFiles, includeDocs)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Tool: search_symbol
	searchSymbolTool := mcp.NewTool("search_symbol",
		mcp.WithDescription("Search for a class or function symbol across the project using AST.\nUseful for finding definitions quickly."),
		mcp.WithString("workspace_path", mcp.Required(), mcp.Description("The base absolute path of the project workspace.")),
		mcp.WithString("query", mcp.Required(), mcp.Description("The symbol name to search for (case-insensitive substring match).")),
	)

	s.AddTool(searchSymbolTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		workspacePath, ok := args["workspace_path"].(string)
		if !ok {
			return mcp.NewToolResultError("workspace_path is required"), nil
		}

		query, ok := args["query"].(string)
		if !ok {
			return mcp.NewToolResultError("query is required"), nil
		}

		res, err := SearchSymbol(workspacePath, query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
