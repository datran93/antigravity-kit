package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	s := server.NewMCPServer(
		"McpDocResearcher",
		"1.0.0",
		server.WithLogging(),
	)

	// Tool: search_latest_syntax
	searchSyntaxTool := mcp.NewTool("search_latest_syntax",
		mcp.WithDescription("Search the real-time internet for the absolute latest SOTA (State Of The Art) syntax, best practices, and documentation for a specific programming topic or library.\\nUse this tool before writing any logic to ensure you are not generating legacy code or using deprecated APIs."),
		mcp.WithString("topic", mcp.Required(), mcp.Description("The specific concept to research (e.g., 'React server components data fetching', 'Next.js 14 App Router layout constraints', 'Zustand slices 2026').")),
		mcp.WithArray("libraries", mcp.Description("Optional list of specific libraries being used to narrow down the context."), mcp.Items(map[string]interface{}{"type": "string"})),
	)

	s.AddTool(searchSyntaxTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		topic, ok := args["topic"].(string)
		if !ok {
			return mcp.NewToolResultError("topic is required"), nil
		}

		searchQuery := topic
		if libsRaw, ok := args["libraries"].([]interface{}); ok && len(libsRaw) > 0 {
			var libs []string
			for _, l := range libsRaw {
				if s, ok := l.(string); ok {
					libs = append(libs, s)
				}
			}
			if len(libs) > 0 {
				searchQuery += " " + strings.Join(libs, " ")
			}
		}

		cacheKey := generateCacheKey("search", searchQuery)
		if cachedResult := getCache(cacheKey); cachedResult != nil {
			return mcp.NewToolResultText(fmt.Sprintf("⚡ (Cached - Loaded instantly from memory)\n%s", *cachedResult)), nil
		}

		results, err := SearchDDGLite(searchQuery+" tutorial OR documentation", 3)
		if err != nil || len(results) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("❌ No recent results found for: %s", topic)), nil
		}

		var finalReport strings.Builder
		finalReport.WriteString(fmt.Sprintf("🔍 REAL-TIME RESEARCH RESULTS FOR: '%s'\n\n", topic))
		finalReport.WriteString("### 1. QUICK SNIPPETS (SEARCH ENGINE RESULTS)\n")

		for i, res := range results {
			finalReport.WriteString(fmt.Sprintf("%d. [%s](%s)\n", i+1, res.Title, res.Href))
			finalReport.WriteString(fmt.Sprintf("   Snippet: %s\n\n", res.Body))
		}

		topURL := results[0].Href
		if topURL != "" {
			finalReport.WriteString("### 2. DEEP DIVE (PARTIAL EXTRACTION OF TOP RESULT)\n")
			finalReport.WriteString(fmt.Sprintf("Source: %s\n", topURL))

			cacheKeyFull := generateCacheKey("url_full", topURL)
			var fullContent string
			if cachedFull := getCache(cacheKeyFull); cachedFull != nil {
				fullContent = *cachedFull
			} else {
				fullContent = FetchJinaMarkdown(topURL, "")
				if !strings.HasPrefix(fullContent, "Error fetching markdown") {
					setCache(cacheKeyFull, fullContent)
				}
			}

			totalLen := len(fullContent)
			truncatedLen := 4000
			if totalLen < 4000 {
				truncatedLen = totalLen
			}
			truncated := fullContent[:truncatedLen]

			if totalLen > 4000 {
				finalReport.WriteString(fmt.Sprintf("Reading content... (Previewing first %d/%d characters)\n", truncatedLen, totalLen))
				finalReport.WriteString(fmt.Sprintf("\n```markdown\n%s\n...\n```\n", truncated))
				finalReport.WriteString(fmt.Sprintf("\n💡 NOTICE: The full document has %d characters. To read beyond this preview, use `read_website_markdown(url=\"%s\", page=1)`.\n", totalLen, topURL))
			} else {
				finalReport.WriteString("Reading content...\n")
				finalReport.WriteString(fmt.Sprintf("\n```markdown\n%s\n```\n", truncated))
			}
		}

		finalReport.WriteString("\n💡 ADVICE FOR AGENT: Synthesize these latest patterns and strictly apply them to your code generation. DO NOT use legacy patterns from your original training data if they conflict with these new docs.")

		finalReportStr := finalReport.String()
		setCache(cacheKey, finalReportStr)

		return mcp.NewToolResultText(finalReportStr), nil
	})

	// Tool: read_website_markdown
	readWebTool := mcp.NewTool("read_website_markdown",
		mcp.WithDescription("Scrape any specific documentation URL or website and return its content perfectly formatted as clean Markdown.\\nSupports pagination for large documents. Each page returns up to 8000 characters."),
		mcp.WithString("url", mcp.Required(), mcp.Description("The absolute URL including https:// (e.g. 'https://react.dev/reference/react/useActionState')")),
		mcp.WithNumber("page", mcp.Description("The page number to read (default 1).")),
		mcp.WithString("cookie", mcp.Description("Optional cookie string to authenticate or bypass protections when scraping the webpage (e.g. 'session_id=1234')")),
	)

	s.AddTool(readWebTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		url, ok := args["url"].(string)
		if !ok {
			return mcp.NewToolResultError("url is required"), nil
		}

		pageFloat, ok := args["page"].(float64)
		page := 1
		if ok {
			page = int(pageFloat)
		}

		cookie, _ := args["cookie"].(string)

		cacheKey := generateCacheKey("url_full", url+cookie)
		var content string
		if cachedFull := getCache(cacheKey); cachedFull != nil {
			content = *cachedFull
		} else {
			content = FetchJinaMarkdown(url, cookie)
			if !strings.HasPrefix(content, "Error fetching markdown") {
				setCache(cacheKey, content)
			}
		}

		if strings.HasPrefix(content, "Error fetching markdown") {
			return mcp.NewToolResultText(content), nil
		}

		chunkSize := 8000
		totalLength := len(content)
		totalPages := (totalLength + chunkSize - 1) / chunkSize
		if totalPages == 0 {
			totalPages = 1
		}

		if page < 1 {
			page = 1
		}
		if page > totalPages {
			page = totalPages
		}

		startIdx := (page - 1) * chunkSize
		endIdx := startIdx + chunkSize
		if endIdx > totalLength {
			endIdx = totalLength
		}

		pageContent := content[startIdx:endIdx]

		header := fmt.Sprintf("📄 Source: %s | Page %d/%d\n--------------------------------------------------\n", url, page, totalPages)
		footer := "\n--------------------------------------------------\n"

		if page < totalPages {
			footer += fmt.Sprintf("💡 (Page %d/%d. There is more content. Extract the next page by calling this tool again with page=%d)\n", page, totalPages, page+1)
		} else {
			footer += "✅ (End of document)\n"
		}

		finalStr := header + pageContent + footer
		return mcp.NewToolResultText(finalStr), nil
	})

	// Tool: read_doc_file
	readDocTool := mcp.NewTool("read_doc_file",
		mcp.WithDescription("Read the contents of a local .doc, .docx, or .pdf document.\\nSupports pagination for large documents. Each page returns up to 8000 characters."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("The absolute path to the local .doc, .docx, or .pdf file (e.g. '/Users/doc.pdf')")),
		mcp.WithNumber("page", mcp.Description("The page number to read (default 1).")),
	)

	s.AddTool(readDocTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		filePath, ok := args["file_path"].(string)
		if !ok {
			return mcp.NewToolResultError("file_path is required"), nil
		}

		pageFloat, ok := args["page"].(float64)
		page := 1
		if ok {
			page = int(pageFloat)
		}

		cacheKey := generateCacheKey("doc_full", filePath)
		var content string
		if cachedFull := getCache(cacheKey); cachedFull != nil {
			content = *cachedFull
		} else {
			var err error
			content, err = ReadLocalDoc(filePath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to read document: %v", err)), nil
			}
			setCache(cacheKey, content)
		}

		chunkSize := 8000
		totalLength := len(content)
		totalPages := (totalLength + chunkSize - 1) / chunkSize
		if totalPages == 0 {
			totalPages = 1
		}

		if page < 1 {
			page = 1
		}
		if page > totalPages {
			page = totalPages
		}

		startIdx := (page - 1) * chunkSize
		endIdx := startIdx + chunkSize
		if endIdx > totalLength {
			endIdx = totalLength
		}

		pageContent := content[startIdx:endIdx]

		header := fmt.Sprintf("📄 Source: %s | Page %d/%d\n--------------------------------------------------\n", filePath, page, totalPages)
		footer := "\n--------------------------------------------------\n"

		if page < totalPages {
			footer += fmt.Sprintf("💡 (Page %d/%d. There is more content. Extract the next page by calling this tool again with page=%d)\n", page, totalPages, page+1)
		} else {
			footer += "✅ (End of document)\n"
		}

		finalStr := header + pageContent + footer
		return mcp.NewToolResultText(finalStr), nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
