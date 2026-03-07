package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CheckpointRow struct {
	TaskID         string
	Description    string
	Status         string
	CompletedSteps string // JSON array
	NextSteps      string // JSON array
	ActiveFiles    string // JSON array
	Notes          string
	UpdatedAt      string
}

func WriteMarkdownProgress(db *sql.DB, workspacePath, taskID, description, status string, completedSteps, nextSteps, activeFiles []string, notes string) error {
	mdPath := filepath.Join(workspacePath, "progress.md")

	// Get historical tasks
	rows, err := db.Query("SELECT task_id, description FROM checkpoints WHERE status = 'completed' AND task_id != ? ORDER BY updated_at DESC", taskID)
	var historicalTasks []map[string]string
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tID, tDesc string
			if err := rows.Scan(&tID, &tDesc); err == nil {
				historicalTasks = append(historicalTasks, map[string]string{"task_id": tID, "description": tDesc})
			}
		}
	}

	totalSteps := len(completedSteps) + len(nextSteps)
	progressPct := float64(0)
	if totalSteps > 0 {
		progressPct = float64(len(completedSteps)) / float64(totalSteps) * 100.0
	}

	barLen := 20
	filled := int(float64(barLen) * progressPct / 100.0)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)

	content := fmt.Sprintf("# 🚀 Project Progress: %s\n\n", taskID)
	content += fmt.Sprintf("**Status:** `%s` | **Progress:** `[%s] %.1f%%` (%d/%d)\n\n", strings.ToUpper(status), bar, progressPct, len(completedSteps), totalSteps)
	content += fmt.Sprintf("> %s\n\n", description)

	if len(activeFiles) > 0 {
		content += "### 📁 Active Files\n"
		for _, f := range activeFiles {
			content += fmt.Sprintf("- `%s`\n", f)
		}
		content += "\n"
	}

	content += "### ✅ Completed\n"
	if len(completedSteps) == 0 {
		content += "_None yet_\n"
	} else {
		for _, s := range completedSteps {
			content += fmt.Sprintf("- [x] %s\n", s)
		}
	}
	content += "\n"

	content += "### ⏳ Next Steps\n"
	if len(nextSteps) == 0 {
		content += "_All tasks done!_ 🎉\n"
	} else {
		for _, s := range nextSteps {
			content += fmt.Sprintf("- [ ] %s\n", s)
		}
	}
	content += "\n"

	if notes != "" {
		content += "### 📝 Log & Notes\n"
		content += fmt.Sprintf("```text\n%s\n```\n\n", notes)
	}

	if len(historicalTasks) > 0 {
		content += "---\n### 🏆 Historically Completed Tasks\n"
		for _, t := range historicalTasks {
			content += fmt.Sprintf("- **%s**: %s\n", t["task_id"], t["description"])
		}
		content += "\n"
	}

	content += fmt.Sprintf("---\n*Last sync: %s*", time.Now().Format("2006-01-02 15:04:05"))

	return os.WriteFile(mdPath, []byte(content), 0644)
}

func SaveCheckpoint(workspacePath, taskID, description, status string, completedSteps, nextSteps, activeFiles []string, notes string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	now := time.Now().Format(time.RFC3339)

	compBytes, _ := json.Marshal(completedSteps)
	nextBytes, _ := json.Marshal(nextSteps)
	activeBytes, _ := json.Marshal(activeFiles)

	query := `
        INSERT INTO checkpoints (task_id, description, status, completed_steps, next_steps, active_files, notes, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(task_id) DO UPDATE SET
            description=excluded.description,
            status=excluded.status,
            completed_steps=excluded.completed_steps,
            next_steps=excluded.next_steps,
            active_files=excluded.active_files,
            notes=excluded.notes,
            updated_at=excluded.updated_at
    `
	_, err = db.Exec(query, taskID, description, status, string(compBytes), string(nextBytes), string(activeBytes), notes, now)
	if err != nil {
		return "", fmt.Errorf("failed to save checkpoint table: %v", err)
	}

	if err := WriteMarkdownProgress(db, workspacePath, taskID, description, status, completedSteps, nextSteps, activeFiles, notes); err != nil {
		fmt.Printf("Error writing progress.md: %v\n", err)
	}

	msg := fmt.Sprintf("✅ Checkpoint '%s' saved.", taskID)
	if len(nextSteps) == 0 && len(completedSteps) > 0 {
		msg += "\n\n🎉 ALL TASKS COMPLETED! Great job."
	}
	return msg, nil
}
