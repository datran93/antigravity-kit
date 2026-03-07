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
		"McpDatabaseInspector",
		"1.0.0",
		server.WithLogging(),
	)

	// Tool: list_tables
	listTablesTool := mcp.NewTool("list_tables",
		mcp.WithDescription("List all tables and views in the database.\\nFor Redis, it returns a summary of keys by exploring the database.\\n\\nArgs:\\n    connection_string: SQLAlchemy connection string.\\n                       Example Postgres: postgresql://user:password@localhost:5432/dbname\\n                       Example MySQL: mysql+pymysql://user:password@localhost:3306/dbname\\n                       Example Redis: redis://localhost:6379/0"),
		mcp.WithString("connection_string", mcp.Required(), mcp.Description("SQLAlchemy connection string.")),
	)

	s.AddTool(listTablesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		connStr, ok := args["connection_string"].(string)
		if !ok {
			return mcp.NewToolResultError("connection_string is required"), nil
		}

		if isRedis(connStr) {
			res, err := handleRedisListTables(connStr)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
			}
			return mcp.NewToolResultText(res), nil
		}

		res, err := handleSQLListTables(connStr)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Tool: inspect_schema
	inspectSchemaTool := mcp.NewTool("inspect_schema",
		mcp.WithDescription("Get detailed schema information for a specific table. Includes columns, types, primary keys, and foreign keys.\\nFor Redis, table_name acts as the key name. Returns key type and properties."),
		mcp.WithString("connection_string", mcp.Required(), mcp.Description("SQLAlchemy connection string.")),
		mcp.WithString("table_name", mcp.Required(), mcp.Description("The exact name of the table to inspect.")),
	)

	s.AddTool(inspectSchemaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		connStr, _ := args["connection_string"].(string)
		tableName, _ := args["table_name"].(string)

		if isRedis(connStr) {
			res, err := handleRedisInspectSchema(connStr, tableName)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
			}
			return mcp.NewToolResultText(res), nil
		}

		res, err := handleSQLInspectSchema(connStr, tableName)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Tool: explain_query
	explainQueryTool := mcp.NewTool("explain_query",
		mcp.WithDescription("Analyze query execution plan (EXPLAIN ANALYZE) for Postgres/MySQL/SQLite.\\nHelps identify performance bottlenecks and index usage."),
		mcp.WithString("connection_string", mcp.Required(), mcp.Description("SQLAlchemy connection string.")),
		mcp.WithString("query", mcp.Required(), mcp.Description("Query to run.")),
	)

	s.AddTool(explainQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		connStr, _ := args["connection_string"].(string)
		query, _ := args["query"].(string)

		if isRedis(connStr) {
			return mcp.NewToolResultText("⚠️ EXPLAIN is not supported for Redis in this server."), nil
		}

		res, err := handleSQLExplainQuery(connStr, query)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Tool: get_table_sample
	getTableSampleTool := mcp.NewTool("get_table_sample",
		mcp.WithDescription("Retrieve schema (DDL) and 5 sample rows for matching a table structure.\\nReturns result as a clean Markdown report."),
		mcp.WithString("connection_string", mcp.Required(), mcp.Description("SQLAlchemy connection string.")),
		mcp.WithString("table_name", mcp.Required(), mcp.Description("The exact name of the table.")),
	)

	s.AddTool(getTableSampleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		connStr, _ := args["connection_string"].(string)
		tableName, _ := args["table_name"].(string)

		if isRedis(connStr) {
			return mcp.NewToolResultText("⚠️ get_table_sample is not supported for Redis. Use list_tables instead."), nil
		}

		res, err := handleSQLGetTableSample(connStr, tableName)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Tool: run_read_query
	runReadQueryTool := mcp.NewTool("run_read_query",
		mcp.WithDescription("Execute a read-only SQL query to preview raw data. Results are returned in JSON format."),
		mcp.WithString("connection_string", mcp.Required(), mcp.Description("SQLAlchemy connection string.")),
		mcp.WithString("query", mcp.Required(), mcp.Description("SQL string (usually SELECT).")),
		mcp.WithNumber("limit", mcp.Description("Max rows to return (default 500, max 2000).")),
		mcp.WithNumber("offset", mcp.Description("Number of rows to skip.")),
	)

	s.AddTool(runReadQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		connStr, _ := args["connection_string"].(string)
		query, _ := args["query"].(string)

		limit := 500
		if l, ok := args["limit"].(float64); ok {
			limit = int(l)
		}

		offset := 0
		if o, ok := args["offset"].(float64); ok {
			offset = int(o)
		}

		if isRedis(connStr) {
			// Redis just ignores limits/offset for Do wrapper generally
			res, err := handleRedisRunQuery(connStr, query, false, false)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
			}
			return mcp.NewToolResultText(res), nil
		}

		res, err := handleSQLRunReadQuery(connStr, query, limit, offset)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	// Tool: run_write_query
	runWriteQueryTool := mcp.NewTool("run_write_query",
		mcp.WithDescription("Execute a write SQL query (INSERT/UPDATE/DELETE/ALTER/DROP/CREATE).\\nThis tool allows modifying the database.\\nAGENT: You MUST ask the user for explicit confirmation in the chat before calling this tool with confirm=True."),
		mcp.WithString("connection_string", mcp.Required(), mcp.Description("SQLAlchemy connection string.")),
		mcp.WithString("query", mcp.Required(), mcp.Description("The raw SQL query or Redis command to run.")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be True to execute. If False, the tool will return a request for confirmation.")),
	)

	s.AddTool(runWriteQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		connStr, _ := args["connection_string"].(string)
		query, _ := args["query"].(string)
		confirm, _ := args["confirm"].(bool)

		if isRedis(connStr) {
			res, err := handleRedisRunQuery(connStr, query, true, confirm)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
			}
			return mcp.NewToolResultText(res), nil
		}

		res, err := handleSQLRunWriteQuery(connStr, query, confirm)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("❌ Error: %v", err)), nil
		}
		return mcp.NewToolResultText(res), nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
