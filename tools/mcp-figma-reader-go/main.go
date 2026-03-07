package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func getToken() string {
	token := os.Getenv("FIGMA_ACCESS_TOKEN")
	if token != "" {
		return token
	}
	// Try fallback to .env if needed, though mostly passed via mcp_config.json env
	body, err := os.ReadFile(".env")
	if err == nil {
		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "FIGMA_ACCESS_TOKEN=") {
				parts := strings.SplitN(line, "=", 2)
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"'")
				return val
			}
		}
	}
	return ""
}

func parseFigmaURL(u string) map[string]string {
	res := make(map[string]string)

	fileKeyRe := regexp.MustCompile(`figma\.com/(?:file|design)/([a-zA-Z0-9]+)`)
	nodeIdRe := regexp.MustCompile(`node-id=([a-zA-Z0-9%-]+)`)

	fMatch := fileKeyRe.FindStringSubmatch(u)
	if len(fMatch) > 1 {
		res["file_key"] = fMatch[1]
	}

	nMatch := nodeIdRe.FindStringSubmatch(u)
	if len(nMatch) > 1 {
		res["node_id"] = strings.ReplaceAll(nMatch[1], "-", ":")
	}

	return res
}

func fetchFigmaAPI(endpoint string, queryParams map[string]string) (any, error) {
	token := getToken()
	if token == "" {
		return nil, fmt.Errorf("Missing FIGMA_ACCESS_TOKEN. Please set it in your environment variables or mcp_config.json")
	}

	apiUrl := fmt.Sprintf("https://api.figma.com/v1/%s", endpoint)
	if len(queryParams) > 0 {
		q := url.Values{}
		for k, v := range queryParams {
			q.Add(k, v)
		}
		apiUrl = fmt.Sprintf("%s?%s", apiUrl, q.Encode())
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Figma-Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		var errData map[string]any
		errMsg := string(body)
		if err := json.Unmarshal(body, &errData); err == nil {
			if eMsg, ok := errData["err"].(string); ok {
				errMsg = eMsg
			}
		}
		return nil, fmt.Errorf("Figma API Error (%d): %s", resp.StatusCode, errMsg)
	}

	var result any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func readFigmaDesignTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	u, _ := args["url"].(string)
	depth := 1
	if dVal, ok := args["depth"].(float64); ok {
		depth = int(dVal)
	}

	parsed := parseFigmaURL(u)
	fileKey, ok := parsed["file_key"]
	if !ok {
		return mcp.NewToolResultError("❌ Invalid Figma URL. Could not extract file key."), nil
	}

	var result any
	var err error

	if nodeID, ok := parsed["node_id"]; ok {
		result, err = fetchFigmaAPI(fmt.Sprintf("nodes/%s", fileKey), map[string]string{
			"ids": nodeID,
		})
	} else {
		result, err = fetchFigmaAPI(fmt.Sprintf("files/%s", fileKey), map[string]string{
			"depth": fmt.Sprintf("%d", depth),
		})
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error: %v", err)), nil
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func exportFigmaImagesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	u, _ := args["url"].(string)
	format := "png"
	if f, ok := args["format"].(string); ok && f != "" {
		format = strings.ToLower(f)
	}
	scale := 1.0
	if s, ok := args["scale"].(float64); ok && s > 0 {
		scale = s
	}

	parsed := parseFigmaURL(u)
	fileKey, ok := parsed["file_key"]
	if !ok {
		return mcp.NewToolResultError("❌ Invalid Figma URL. Could not extract file key."), nil
	}

	nodeID, ok := parsed["node_id"]
	if !ok {
		return mcp.NewToolResultError("❌ No node-id found in URL. Please select a specific element in Figma and copy the link to that element."), nil
	}

	result, err := fetchFigmaAPI(fmt.Sprintf("images/%s", fileKey), map[string]string{
		"ids":    nodeID,
		"format": format,
		"scale":  fmt.Sprintf("%f", scale),
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error: %v", err)), nil
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func getDesignDetailsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	u, _ := args["url"].(string)

	parsed := parseFigmaURL(u)
	fileKey, ok := parsed["file_key"]
	if !ok {
		return mcp.NewToolResultError("❌ Invalid Figma URL."), nil
	}

	targetIDs := ""
	if ids, ok := args["ids"].(string); ok && ids != "" {
		targetIDs = ids
	} else if nodeID, ok := parsed["node_id"]; ok {
		targetIDs = nodeID
	}

	if targetIDs == "" {
		return mcp.NewToolResultError("❌ No Node IDs provided. Use `read_figma_design` first to find IDs."), nil
	}

	result, err := fetchFigmaAPI(fmt.Sprintf("nodes/%s", fileKey), map[string]string{
		"ids": targetIDs,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error: %v", err)), nil
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(b)), nil
}

func main() {
	s := server.NewMCPServer("FigmaReader", "1.0.0")

	readFigmaDesign := mcp.NewTool("read_figma_design",
		mcp.WithDescription("Reads raw design data from a Figma URL.\nReturns the file structure and metadata."),
		mcp.WithString("url", mcp.Required(), mcp.Description("The full Figma design URL.")),
		mcp.WithNumber("depth", mcp.Description("How deep to traverse the file tree (default 1 for high-level structure).")),
	)
	s.AddTool(readFigmaDesign, readFigmaDesignTool)

	exportFigmaImages := mcp.NewTool("export_figma_images",
		mcp.WithDescription("Renders specific nodes from a Figma URL as images.\nReturns a map of node IDs to temporary image URLs."),
		mcp.WithString("url", mcp.Required(), mcp.Description("The full Figma design URL (must contain a node-id to render specifically).")),
		mcp.WithString("format", mcp.Description("Image format (png, jpg, svg, pdf).")),
		mcp.WithNumber("scale", mcp.Description("Scale factor (0.01 to 4).")),
	)
	s.AddTool(exportFigmaImages, exportFigmaImagesTool)

	getDesignDetails := mcp.NewTool("get_design_details",
		mcp.WithDescription("Gets detailed JSON data for specific nodes in a Figma file."),
		mcp.WithString("url", mcp.Required(), mcp.Description("The Figma file URL.")),
		mcp.WithString("ids", mcp.Description("Comma-separated list of node IDs (optional if URL already has node-id).")),
	)
	s.AddTool(getDesignDetails, getDesignDetailsTool)

	server.ServeStdio(s)
}
