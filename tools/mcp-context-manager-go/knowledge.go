package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func cleanQuery(query string) string {
	var builder strings.Builder
	for _, ch := range query {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || unicode.IsSpace(ch) {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func RecallKnowledge(workspacePath, query string, topK int) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	if topK <= 0 {
		topK = 3
	}

	cleanedQuery := cleanQuery(query)
	parts := strings.Fields(cleanedQuery)
	var tokens []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			tokens = append(tokens, fmt.Sprintf("%s*", p))
		}
	}
	ftsQuery := strings.Join(tokens, " OR ")

	if ftsQuery == "" {
		return "🔍 Please provide a valid search query.", nil
	}

	sqlQuery := `
        SELECT tactic_name, ki_path, summary, decisions
        FROM knowledge_fts
        WHERE knowledge_fts MATCH ?
        ORDER BY rank
        LIMIT ?
    `
	rows, err := db.Query(sqlQuery, ftsQuery, topK)
	if err != nil {
		return "", fmt.Errorf("failed to query knowledge_fts: %v", err)
	}
	defer rows.Close()

	var res []string
	res = append(res, fmt.Sprintf("🧠 Recalled Knowledge for '%s':\n", query))
	count := 0

	for rows.Next() {
		var tacticName, kiPath, summary, decisions string
		if err := rows.Scan(&tacticName, &kiPath, &summary, &decisions); err == nil {
			count++
			res = append(res, fmt.Sprintf("### KI: %s", tacticName))
			res = append(res, fmt.Sprintf("**Path**: `%s`", kiPath))
			res = append(res, fmt.Sprintf("**Summary**: %s", summary))
			res = append(res, fmt.Sprintf("**Decisions**: %s\n---", decisions))
		}
	}

	if count == 0 {
		return fmt.Sprintf("🔍 No relevant Knowledge Items found for query: '%s'", query), nil
	}

	return strings.Join(res, "\n"), nil
}

func CompactMemory(workspacePath, taskID, tacticName, summary, decisions string) (string, error) {
	knowledgeDir := filepath.Join(workspacePath, "knowledge")
	if err := os.MkdirAll(knowledgeDir, 0755); err != nil {
		return "", fmt.Errorf("failed to make knowledge dir: %v", err)
	}

	safeName := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(tacticName)), " ", "_")
	if safeName == "" {
		safeName = "unknown_tactic"
	}
	kiPath := filepath.Join(knowledgeDir, fmt.Sprintf("%s.md", safeName))

	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var lockedFilesStr string
	err = db.QueryRow("SELECT locked_files FROM intents WHERE task_id = ?", taskID).Scan(&lockedFilesStr)
	var files []string
	if err == nil && lockedFilesStr != "" {
		json.Unmarshal([]byte(lockedFilesStr), &files)
	}

	content := fmt.Sprintf("# KI: %s\n\n## Summary\n%s\n\n## Affected Files\n", tacticName, summary)
	for _, f := range files {
		content += fmt.Sprintf("- `%s`\n", f)
	}
	content += fmt.Sprintf("\n## Architecture & Decisions\n%s\n", decisions)

	if err := os.WriteFile(kiPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write KI file: %v", err)
	}

	insertFTS := `
        INSERT INTO knowledge_fts (tactic_name, ki_path, summary, decisions)
        VALUES (?, ?, ?, ?)
    `
	if _, err := db.Exec(insertFTS, tacticName, kiPath, summary, decisions); err != nil {
		return "", fmt.Errorf("failed to index FTS: %v", err)
	}

	var notes string
	if err := db.QueryRow("SELECT notes FROM checkpoints WHERE task_id = ?", taskID).Scan(&notes); err == nil {
		newNotes := notes + fmt.Sprintf("\n[COMPACTION] Tactic '%s' completed. KI saved to %s", tacticName, kiPath)
		db.Exec("UPDATE checkpoints SET active_files='[]', notes=? WHERE task_id=?", newNotes, taskID)
		db.Exec("UPDATE intents SET locked_files='[]' WHERE task_id=?", taskID)
		db.Exec("UPDATE drift_tracker SET failure_count=0 WHERE task_id=?", taskID)
	}

	return fmt.Sprintf("🗜️ Context Compaction successful. Knowledge Item indexed into local RAG and saved to %s. Memory flushed.", kiPath), nil
}
