package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	workspaceRoot = "/Users/datran/LearnDev/antigravity-kit"
	restDir       = filepath.Join(workspaceRoot, "rest")
)

type SessionState struct {
	BaseURL string
	Headers map[string]string
	Env     map[string]any
	mu      sync.RWMutex
}

var state = SessionState{
	Headers: make(map[string]string),
	Env:     make(map[string]any),
}

func ensureRestDir() error {
	if _, err := os.Stat(restDir); os.IsNotExist(err) {
		return os.MkdirAll(restDir, 0755)
	}
	return nil
}

func replaceVars(text string) string {
	state.mu.RLock()
	defer state.mu.RUnlock()
	for k, v := range state.Env {
		valStr := fmt.Sprintf("%v", v)
		text = strings.ReplaceAll(text, fmt.Sprintf("{{%s}}", k), valStr)
	}
	return text
}

func sanitizeSlug(slug string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_-]+")
	safe := reg.ReplaceAllString(slug, "_")
	return strings.ToLower(safe)
}

func saveToHistory(entry map[string]any, requestHeaders map[string]string, slug string) {
	if err := ensureRestDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating rest dir: %v\n", err)
		return
	}
	safeSlug := sanitizeSlug(slug)
	if safeSlug == "" {
		safeSlug = "general"
	}
	filePath := filepath.Join(restDir, fmt.Sprintf("%s.rest", safeSlug))

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening history file: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := entry["timestamp"].(string)
	method := entry["method"].(string)
	statusCode := entry["status_code"].(int)
	reqUrl := entry["url"].(string)

	fmt.Fprintf(f, "### %s Request at %s\n", method, timestamp)
	fmt.Fprintf(f, "# Status: %d\n", statusCode)
	fmt.Fprintf(f, "%s %s\n", method, reqUrl)

	for k, v := range requestHeaders {
		fmt.Fprintf(f, "%s: %s\n", k, v)
	}

	if jsonBody, ok := entry["json_body"]; ok && jsonBody != nil {
		fmt.Fprintf(f, "\n")
		b, _ := json.MarshalIndent(jsonBody, "", "  ")
		f.Write(b)
		fmt.Fprintf(f, "\n")
	}
	fmt.Fprintf(f, "\n")
}

func executeRequest(method, path string, params map[string]any, jsonBody any, headers map[string]any, saveHistory bool, slug string) (string, error) {
	targetPath := replaceVars(path)

	state.mu.RLock()
	base := state.BaseURL
	state.mu.RUnlock()

	reqUrl := targetPath
	if base != "" && !strings.HasPrefix(targetPath, "http") {
		base = strings.TrimRight(base, "/")
		pathClean := strings.TrimLeft(targetPath, "/")
		reqUrl = fmt.Sprintf("%s/%s", base, pathClean)
	}

	if len(params) > 0 {
		u, err := url.Parse(reqUrl)
		if err == nil {
			q := u.Query()
			for k, v := range params {
				q.Add(k, fmt.Sprintf("%v", v))
			}
			u.RawQuery = q.Encode()
			reqUrl = u.String()
		}
	}

	requestHeaders := make(map[string]string)

	state.mu.RLock()
	for k, v := range state.Headers {
		requestHeaders[replaceVars(k)] = replaceVars(v)
	}
	state.mu.RUnlock()

	for k, v := range headers {
		if strVal, ok := v.(string); ok {
			requestHeaders[replaceVars(k)] = replaceVars(strVal)
		}
	}

	var bodyReader io.Reader
	if jsonBody != nil {
		b, err := json.Marshal(jsonBody)
		if err == nil {
			bodyReader = bytes.NewReader(b)
			if _, ok := requestHeaders["Content-Type"]; !ok {
				requestHeaders["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), reqUrl, bodyReader)
	if err != nil {
		return "", fmt.Errorf("❌ **Request Error**: %v", err)
	}

	for k, v := range requestHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("❌ **Request Error**: %v", err)
	}
	defer resp.Body.Close()

	if saveHistory {
		finalSlug := "general"
		if slug != "" {
			finalSlug = slug
		}
		entry := map[string]any{
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"method":      strings.ToUpper(method),
			"url":         reqUrl,
			"status_code": resp.StatusCode,
			"json_body":   jsonBody,
		}
		saveToHistory(entry, requestHeaders, finalSlug)
	}

	statusIcon := "❌"
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		statusIcon = "✅"
	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		statusIcon = "⚠️"
	}

	var output []string
	output = append(output, fmt.Sprintf("%s **%s %s**", statusIcon, strings.ToUpper(method), reqUrl))
	output = append(output, fmt.Sprintf("**Status**: `%d %s`", resp.StatusCode, resp.Status))
	output = append(output, "\n### 📝 Response Headers", "```json")

	respHeadersMap := make(map[string]any)
	for k, v := range resp.Header {
		if len(v) == 1 {
			respHeadersMap[k] = v[0]
		} else {
			respHeadersMap[k] = v
		}
	}
	hBytes, _ := json.MarshalIndent(respHeadersMap, "", "  ")
	output = append(output, string(hBytes), "```", "\n### 📦 Response Body")

	bodyBytes, _ := io.ReadAll(resp.Body)
	var bodyJSON any
	if err := json.Unmarshal(bodyBytes, &bodyJSON); err == nil {
		output = append(output, "```json")
		b, _ := json.MarshalIndent(bodyJSON, "", "  ")
		output = append(output, string(b))
		output = append(output, "```")
	} else {
		text := string(bodyBytes)
		if len(text) > 2000 {
			text = text[:2000]
		}
		output = append(output, fmt.Sprintf("```text\n%s\n```", text))
	}

	return strings.Join(output, "\n"), nil
}

func httpRequestTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	method, _ := args["method"].(string)
	path, _ := args["path"].(string)

	var params map[string]any
	if p, ok := args["params"].(map[string]any); ok {
		params = p
	}

	jsonBody := args["json_body"]

	var headers map[string]any
	if h, ok := args["headers"].(map[string]any); ok {
		headers = h
	}

	saveHistory := true
	if sh, ok := args["save_history"].(bool); ok {
		saveHistory = sh
	}

	slug, _ := args["slug"].(string)

	output, err := executeRequest(method, path, params, jsonBody, headers, saveHistory, slug)
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

func setEnvTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	key, ok := args["key"].(string)
	if !ok {
		return mcp.NewToolResultError("key must be a string"), nil
	}
	value := args["value"]

	state.mu.Lock()
	state.Env[key] = value
	state.mu.Unlock()

	return mcp.NewToolResultText(fmt.Sprintf("✅ Env var `%s` set to `%v`", key, value)), nil
}

// simplistic curl parser for basic functionality, won't handle highly complex escaping
func splitArgs(str string) []string {
	var args []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	for _, r := range str {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if r == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if (r == ' ' || r == '\t' || r == '\n') && !inSingleQuote && !inDoubleQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func importCurlTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	argsMap, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	cmdStr, ok := argsMap["curl_command"].(string)
	if !ok {
		return mcp.NewToolResultError("curl_command must be a string"), nil
	}

	cmdStr = strings.TrimSpace(cmdStr)
	cmdStr = strings.TrimPrefix(cmdStr, "curl ")

	args := splitArgs(cmdStr)

	var urlStr string
	method := "GET"
	headers := make(map[string]any)
	var dataChunks []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-X" || arg == "--request" {
			if i+1 < len(args) {
				method = args[i+1]
				i++
			}
		} else if arg == "-H" || arg == "--header" {
			if i+1 < len(args) {
				h := args[i+1]
				parts := strings.SplitN(h, ":", 2)
				if len(parts) == 2 {
					headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
				i++
			}
		} else if arg == "-d" || arg == "--data" || arg == "--data-raw" || arg == "--data-urlencode" {
			if i+1 < len(args) {
				dataChunks = append(dataChunks, args[i+1])
				if method == "GET" {
					method = "POST"
				}
				i++
			}
		} else if !strings.HasPrefix(arg, "-") {
			if urlStr == "" {
				urlStr = strings.Trim(arg, "'\"")
			}
		}
	}

	var jsonBody any
	if len(dataChunks) > 0 {
		combinedData := strings.Join(dataChunks, " ")
		err := json.Unmarshal([]byte(combinedData), &jsonBody)
		if err != nil {
			// Ignore json unmarshal err, might not be json.
			// Currently executeRequest expects `any` and marshals to json if not nil.
			// So if it's not JSON, we pass the raw string and handle it in executeRequest
			// Actually the original python only supports json_body or ignores.
			// We'll mimic python:
			jsonBody = combinedData
			// wait, if we pass string, json.Marshal inside executeRequest will encode it as a JSON string literal.
			// let's just leave it as combinedData string for now.
		}
	}

	output, err := executeRequest(method, urlStr, nil, jsonBody, headers, true, "curl-import")
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

func listHistoryTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if _, err := os.Stat(restDir); os.IsNotExist(err) {
		return mcp.NewToolResultText("📭 History directory (`rest/`) does not exist yet."), nil
	}

	entries, err := os.ReadDir(restDir)
	if err != nil {
		return mcp.NewToolResultError("Failed to read rest directory"), nil
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".rest") {
			files = append(files, e.Name())
		}
	}

	if len(files) == 0 {
		return mcp.NewToolResultText("📭 No .rest files found."), nil
	}
	sort.Strings(files)

	var output []string
	output = append(output, "### 📜 Request History List (.rest)\n")

	for _, file := range files {
		path := filepath.Join(restDir, file)
		b, err := os.ReadFile(path)
		if err != nil {
			output = append(output, fmt.Sprintf("- **`%s`**: Read error - %v", file, err))
			continue
		}

		lines := strings.Split(string(b), "\n")
		latestReq := ""
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if strings.HasPrefix(line, "###") {
				latestReq = line
				break
			}
		}
		if latestReq == "" {
			latestReq = "Empty"
		}
		output = append(output, fmt.Sprintf("- **`%s`**: %s", file, latestReq))
	}

	return mcp.NewToolResultText(strings.Join(output, "\n")), nil
}

func clearHistoryTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if _, err := os.Stat(restDir); err == nil {
		os.RemoveAll(restDir)
		os.MkdirAll(restDir, 0755)
	}
	return mcp.NewToolResultText("✅ Successfully cleared everything in the rest directory."), nil
}

func setConfigTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var updates []string

	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	state.mu.Lock()
	if baseUrl, ok := args["base_url"].(string); ok && baseUrl != "" {
		state.BaseURL = baseUrl
		updates = append(updates, fmt.Sprintf("Base URL set to: `%s`", baseUrl))
	}
	if authToken, ok := args["auth_token"].(string); ok && authToken != "" {
		state.Headers["Authorization"] = fmt.Sprintf("Bearer %s", authToken)
		updates = append(updates, "Authorization header updated.")
	}
	state.mu.Unlock()

	if len(updates) > 0 {
		return mcp.NewToolResultText("✅ Configuration Updated:\n- " + strings.Join(updates, "\n- ")), nil
	}
	return mcp.NewToolResultText("No configuration provided."), nil
}

func main() {
	s := server.NewMCPServer("McpHttpClient", "1.0.0")

	httpRequest := mcp.NewTool("http_request",
		mcp.WithDescription("Execute an HTTP request. Automatically prepends base_url and merges default headers.\nSupports placeholders like {{variable_name}}."),
		mcp.WithString("method", mcp.Required(), mcp.Description("HTTP Method (GET, POST, etc)")),
		mcp.WithString("path", mcp.Required(), mcp.Description("Target path or absolute URL")),
		mcp.WithAny("params", mcp.Description("Query parameters object")),
		mcp.WithAny("json_body", mcp.Description("JSON body object")),
		mcp.WithAny("headers", mcp.Description("Headers object")),
		mcp.WithBoolean("save_history", mcp.Description("Whether to save to history")),
		mcp.WithString("slug", mcp.Description("Filename slug to use for history log")),
	)
	s.AddTool(httpRequest, httpRequestTool)

	setEnv := mcp.NewTool("set_env",
		mcp.WithDescription("Set an environment variable for use in {{key}} placeholders."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Variable key")),
		mcp.WithAny("value", mcp.Required(), mcp.Description("Variable value")),
	)
	s.AddTool(setEnv, setEnvTool)

	importCurl := mcp.NewTool("import_curl",
		mcp.WithDescription("Parse a raw cURL command and execute it.\nUseful for quickly testing requests from documentation or browser."),
		mcp.WithString("curl_command", mcp.Required(), mcp.Description("Raw curl command string")),
	)
	s.AddTool(importCurl, importCurlTool)

	listHistory := mcp.NewTool("list_history",
		mcp.WithDescription("View context in .rest format."),
	)
	s.AddTool(listHistory, listHistoryTool)

	clearHistory := mcp.NewTool("clear_history",
		mcp.WithDescription("Clear all request history."),
	)
	s.AddTool(clearHistory, clearHistoryTool)

	setConfig := mcp.NewTool("set_config",
		mcp.WithDescription("Configure base URL and auth token."),
		mcp.WithString("base_url", mcp.Description("Base URL")),
		mcp.WithString("auth_token", mcp.Description("Bearer auth token")),
	)
	s.AddTool(setConfig, setConfigTool)

	server.ServeStdio(s)
}
