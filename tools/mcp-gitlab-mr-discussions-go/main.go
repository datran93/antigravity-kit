package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func getGitlabClient() (*gitlab.Client, error) {
	token := os.Getenv("GITLAB_PRIVATE_TOKEN")
	url := os.Getenv("GITLAB_URL")
	if url == "" {
		url = "https://gitlab.com"
	}
	if token == "" {
		return nil, fmt.Errorf("GITLAB_PRIVATE_TOKEN environment variable is not set. Please set it in your environment or in mcp_config.json env mapping")
	}

	return gitlab.NewClient(token, gitlab.WithBaseURL(url))
}

func readMRDiscussionsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	projectID, _ := args["project_id"].(string)
	mrIIDFloat, ok := args["mr_iid"].(float64)
	if !ok {
		return mcp.NewToolResultError("mr_iid must be a number"), nil
	}
	mrIID := int64(mrIIDFloat)

	gl, err := getGitlabClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error: %v", err)), nil
	}

	opt := &gitlab.ListMergeRequestDiscussionsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	var allDiscussions []*gitlab.Discussion
	for {
		discussions, resp, err := gl.Discussions.ListMergeRequestDiscussions(projectID, mrIID, opt)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("❌ Error reading MR discussions: %v", err)), nil
		}
		allDiscussions = append(allDiscussions, discussions...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var output []string
	output = append(output, fmt.Sprintf("💬 DISCUSSIONS FOR MR #%d (Project: %s)\n", mrIID, projectID))

	count := 0
	for _, disc := range allDiscussions {
		if len(disc.Notes) == 0 {
			continue
		}

		isSystemOnly := true
		for _, note := range disc.Notes {
			if !note.System {
				isSystemOnly = false
				break
			}
		}
		if isSystemOnly {
			continue
		}

		count++

		statusText := ""
		if disc.Notes[0].Resolvable {
			status := "❌ UNRESOLVED"
			if disc.Notes[0].Resolved {
				status = "✅ RESOLVED"
			}
			statusText = fmt.Sprintf(" | Status: %s", status)
		}

		output = append(output, fmt.Sprintf("--- Discussion ID: %s%s ---", disc.ID, statusText))

		for _, note := range disc.Notes {
			author := note.Author.Username
			if author == "" {
				author = "unknown"
			}
			body := note.Body
			created := ""
			if note.CreatedAt != nil {
				created = note.CreatedAt.Format("2006-01-02 15:04:05")
			}
			isSystem := ""
			if note.System {
				isSystem = "[SYSTEM]"
			}

			output = append(output, fmt.Sprintf("[%s] @%s %s (Note ID: %d):\n%s\n", created, author, isSystem, note.ID, body))
		}
	}

	if count == 0 {
		return mcp.NewToolResultText("✅ No user discussions found on this Merge Request."), nil
	}

	return mcp.NewToolResultText(strings.Join(output, "\n")), nil
}

func replyToMRDiscussionTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	projectID, _ := args["project_id"].(string)
	mrIIDFloat, ok := args["mr_iid"].(float64)
	if !ok {
		return mcp.NewToolResultError("mr_iid must be a number"), nil
	}
	mrIID := int64(mrIIDFloat)
	discussionID, _ := args["discussion_id"].(string)
	body, _ := args["body"].(string)

	gl, err := getGitlabClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error: %v", err)), nil
	}

	opt := &gitlab.AddMergeRequestDiscussionNoteOptions{
		Body: gitlab.Ptr(body),
	}

	_, _, err = gl.Discussions.AddMergeRequestDiscussionNote(projectID, mrIID, discussionID, opt)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error replying to discussion: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully replied to discussion '%s' on MR #%d.", discussionID, mrIID)), nil
}

func resolveMRDiscussionTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	projectID, _ := args["project_id"].(string)
	mrIIDFloat, ok := args["mr_iid"].(float64)
	if !ok {
		return mcp.NewToolResultError("mr_iid must be a number"), nil
	}
	mrIID := int64(mrIIDFloat)
	discussionID, _ := args["discussion_id"].(string)

	resolve := true
	if r, ok := args["resolve"].(bool); ok {
		resolve = r
	}

	gl, err := getGitlabClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error: %v", err)), nil
	}

	// First check if it exists and is resolvable
	disc, _, err := gl.Discussions.GetMergeRequestDiscussion(projectID, mrIID, discussionID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Discussion '%s' not found or error: %v", discussionID, err)), nil
	}

	if len(disc.Notes) == 0 || !disc.Notes[0].Resolvable {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Discussion '%s' is not resolvable.", discussionID)), nil
	}

	opt := &gitlab.ResolveMergeRequestDiscussionOptions{
		Resolved: gitlab.Ptr(resolve),
	}

	_, _, err = gl.Discussions.ResolveMergeRequestDiscussion(projectID, mrIID, discussionID, opt)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error resolving discussion: %v", err)), nil
	}

	action := "resolved"
	if !resolve {
		action = "unresolved"
	}
	return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully %s discussion '%s' on MR #%d.", action, discussionID, mrIID)), nil
}

func main() {
	s := server.NewMCPServer("McpGitLabMRDiscussions", "1.0.0")

	readMRDiscussions := mcp.NewTool("read_mr_discussions",
		mcp.WithDescription("Read all discussions (threads) from a specific GitLab Merge Request."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("The ID or URL-encoded path of the project")),
		mcp.WithNumber("mr_iid", mcp.Required(), mcp.Description("The internal ID of the merge request")),
	)
	s.AddTool(readMRDiscussions, readMRDiscussionsTool)

	replyToMRDiscussion := mcp.NewTool("reply_to_mr_discussion",
		mcp.WithDescription("Reply to an existing discussion thread on a GitLab Merge Request."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("The ID or URL-encoded path of the project")),
		mcp.WithNumber("mr_iid", mcp.Required(), mcp.Description("The internal ID of the merge request")),
		mcp.WithString("discussion_id", mcp.Required(), mcp.Description("The ID of the discussion thread to reply to")),
		mcp.WithString("body", mcp.Required(), mcp.Description("The text content of your reply")),
	)
	s.AddTool(replyToMRDiscussion, replyToMRDiscussionTool)

	resolveMRDiscussion := mcp.NewTool("resolve_mr_discussion",
		mcp.WithDescription("Resolve or unresolve a discussion thread on a GitLab Merge Request."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("The ID or URL-encoded path of the project")),
		mcp.WithNumber("mr_iid", mcp.Required(), mcp.Description("The internal ID of the merge request")),
		mcp.WithString("discussion_id", mcp.Required(), mcp.Description("The ID of the discussion thread")),
		mcp.WithBoolean("resolve", mcp.Required(), mcp.Description("True to resolve, False to unresolve")),
	)
	s.AddTool(resolveMRDiscussion, resolveMRDiscussionTool)

	server.ServeStdio(s)
}
