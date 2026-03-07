package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

var (
	workspaceRoot = "/Users/datran/LearnDev/antigravity-kit"
	skillsDir     = filepath.Join(workspaceRoot, ".agent", "skills")
	dbDir         = filepath.Join(workspaceRoot, "tools", "mcp-skill-router-go", ".db")
	dbFile        = filepath.Join(dbDir, "skills_cache.json")
	mu            sync.Mutex
)

type SkillMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Path        string `json:"path"`
	Preview     string `json:"preview"`
	Hash        string `json:"hash"`
}

type SkillDoc struct {
	ID        string        `json:"id"`
	Text      string        `json:"text"`
	Metadata  SkillMetadata `json:"metadata"`
	Embedding []float32     `json:"embedding"`
}

func getOpenAIClient() (*openai.Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ERROR: OPENAI_API_KEY environment variable is not set")
	}
	return openai.NewClient(apiKey), nil
}

func parseSkillFile(path string) (*SkillDoc, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)
	hashBytes := md5.Sum(content)
	fileHash := hex.EncodeToString(hashBytes[:])

	re := regexp.MustCompile("(?s)^---\n(.*?)\n---")
	match := re.FindStringSubmatch(contentStr)

	var meta map[string]any
	if len(match) > 1 {
		_ = yaml.Unmarshal([]byte(match[1]), &meta)
	}

	desc := ""
	if d, ok := meta["description"].(string); ok {
		desc = d
	}

	tagsStr := ""
	if t, ok := meta["tags"]; ok {
		switch v := t.(type) {
		case []interface{}:
			var ts []string
			for _, item := range v {
				ts = append(ts, fmt.Sprintf("%v", item))
			}
			tagsStr = strings.Join(ts, ", ")
		case string:
			tagsStr = v
		}
	}

	textContent := re.ReplaceAllString(contentStr, "")
	textContent = strings.TrimSpace(textContent)

	words := strings.Fields(textContent)
	preview := strings.Join(words, " ")
	if len(preview) > 250 {
		preview = preview[:250] + "..."
	}

	skillName := filepath.Base(filepath.Dir(path))

	searchText := fmt.Sprintf("Skill: %s\nTags: %s\nDescription: %s\n\nPreview: %s", skillName, tagsStr, desc, preview)
	descClean := strings.ReplaceAll(desc, "\n", " ")
	if len(descClean) > 150 {
		descClean = descClean[:150]
	}

	return &SkillDoc{
		ID:   skillName,
		Text: searchText,
		Metadata: SkillMetadata{
			Name:        skillName,
			Description: descClean,
			Tags:        tagsStr,
			Path:        path,
			Preview:     preview,
			Hash:        fileHash,
		},
	}, nil
}

func loadCache() (map[string]SkillDoc, error) {
	cache := make(map[string]SkillDoc)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return cache, nil
	}
	data, err := os.ReadFile(dbFile)
	if err != nil {
		return nil, err
	}
	var docs []SkillDoc
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
	}
	for _, d := range docs {
		cache[d.ID] = d
	}
	return cache, nil
}

func saveCache(cache map[string]SkillDoc) error {
	os.MkdirAll(dbDir, 0755)
	var docs []SkillDoc
	for _, d := range cache {
		docs = append(docs, d)
	}
	data, err := json.MarshalIndent(docs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbFile, data, 0644)
}

func cosineSimilarity(a, b []float32) float32 {
	var dotProduct, normA, normB float32
	for i := 0; i < len(a) && i < len(b); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

func ensureIndex() (map[string]SkillDoc, error) {
	mu.Lock()
	defer mu.Unlock()

	client, err := getOpenAIClient()
	if err != nil {
		return nil, err
	}

	cache, err := loadCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
		cache = make(map[string]SkillDoc)
	}

	skillDirs, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading skills dir: %v", err)
	}

	currentIDs := make(map[string]bool)
	var toUpdate []SkillDoc
	var updateTexts []string

	for _, dir := range skillDirs {
		if !dir.IsDir() {
			continue
		}
		skillPath := filepath.Join(skillsDir, dir.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}
		doc, err := parseSkillFile(skillPath)
		if err != nil {
			continue
		}
		currentIDs[doc.ID] = true
		if cached, ok := cache[doc.ID]; ok && cached.Metadata.Hash == doc.Metadata.Hash && len(cached.Embedding) > 0 {
			doc.Embedding = cached.Embedding
			cache[doc.ID] = *doc
			continue
		}
		toUpdate = append(toUpdate, *doc)
		updateTexts = append(updateTexts, doc.Text)
	}

	changed := false
	for id := range cache {
		if !currentIDs[id] {
			delete(cache, id)
			changed = true
		}
	}

	if len(toUpdate) > 0 {
		fmt.Fprintf(os.Stderr, "Indexing %d new/modified skills...\n", len(toUpdate))

		// Batch requests to avoid rate limits
		batchSize := 20
		for i := 0; i < len(updateTexts); i += batchSize {
			end := i + batchSize
			if end > len(updateTexts) {
				end = len(updateTexts)
			}

			resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
				Input: updateTexts[i:end],
				Model: openai.SmallEmbedding3,
			})
			if err != nil {
				return nil, fmt.Errorf("error fetching embeddings: %v", err)
			}
			for j, emp := range resp.Data {
				idx := i + j
				toUpdate[idx].Embedding = emp.Embedding
				cache[toUpdate[idx].ID] = toUpdate[idx]
			}
			time.Sleep(100 * time.Millisecond) // small delay between batches
		}
		changed = true
	}

	if changed {
		saveCache(cache)
	}

	return cache, nil
}

type SearchResult struct {
	Doc   SkillDoc
	Score float32
}

func searchSkillsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	query, _ := args["query"].(string)
	tagsFilter, _ := args["tags_filter"].(string)

	topK := 3
	if tk, ok := args["top_k"].(float64); ok {
		topK = int(tk)
	}

	client, err := getOpenAIClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cache, err := ensureIndex()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error initializing index: %v", err)), nil
	}

	if len(cache) == 0 {
		return mcp.NewToolResultText("❌ No skills found in directory."), nil
	}

	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{query},
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("❌ Error embedding query: %v", err)), nil
	}
	queryEmb := resp.Data[0].Embedding

	var filterTags []string
	if tagsFilter != "" {
		parts := strings.Split(tagsFilter, ",")
		for _, p := range parts {
			f := strings.TrimSpace(strings.ToLower(p))
			if f != "" {
				filterTags = append(filterTags, f)
			}
		}
	}

	var results []SearchResult
	for _, doc := range cache {
		skillTags := strings.ToLower(doc.Metadata.Tags)
		valid := true
		if len(filterTags) > 0 {
			for _, ft := range filterTags {
				if !strings.Contains(skillTags, ft) {
					valid = false
					break
				}
			}
		}

		if valid && len(doc.Embedding) > 0 {
			score := cosineSimilarity(queryEmb, doc.Embedding)
			results = append(results, SearchResult{Doc: doc, Score: score})
		}
	}

	if len(results) == 0 {
		if tagsFilter != "" {
			return mcp.NewToolResultText(fmt.Sprintf("❌ No relevant skills matching tags: '%s'", tagsFilter)), nil
		}
		return mcp.NewToolResultText("❌ No relevant skills found."), nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score // Descending
	})

	if len(results) > topK {
		results = results[:topK]
	}

	out := []string{fmt.Sprintf("🎯 SEMANTIC SEARCH RESULTS FOR QUERY: '%s'", query)}
	for _, r := range results {
		res := fmt.Sprintf("\n🔹 **%s**", r.Doc.Metadata.Name)
		res += fmt.Sprintf("\n   - Path: `%s`", r.Doc.Metadata.Path)
		res += fmt.Sprintf("\n   - Tags: %s", r.Doc.Metadata.Tags)
		res += fmt.Sprintf("\n   - Description: %s", r.Doc.Metadata.Description)
		if r.Doc.Metadata.Preview != "" {
			res += fmt.Sprintf("\n   - Preview: %s", r.Doc.Metadata.Preview)
		}
		out = append(out, res)
	}

	out = append(out, "\n💡 ADVICE FOR AGENT: Use the `view_file` tool on the Path provided above to read the skill details before starting your task.")

	return mcp.NewToolResultText(strings.Join(out, "\n")), nil
}

func main() {
	s := server.NewMCPServer("McpSkillRouter", "1.0.0")

	searchSkills := mcp.NewTool("search_skills",
		mcp.WithDescription("Search for relevant AI agent skills based on a semantic query.\nUse this tool when you need to find which skills are best suited for a user's request."),
		mcp.WithString("query", mcp.Required(), mcp.Description("Semantic query to search for (e.g., 'beautiful ui design', 'fix database query')")),
		mcp.WithString("tags_filter", mcp.Description("Optional comma-separated tags to filter by (e.g., 'frontend, react'). Will only return skills containing these tags.")),
		mcp.WithNumber("top_k", mcp.Description("Number of skills to return (default is 3)")),
	)
	s.AddTool(searchSkills, searchSkillsTool)

	server.ServeStdio(s)
}
