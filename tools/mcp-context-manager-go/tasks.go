package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func InitializeTaskPlan(workspacePath, taskID, description string, steps []string) (string, error) {
	initNotes := fmt.Sprintf("[%s] Task started.", time.Now().Format("15:04:05"))
	return SaveCheckpoint(
		workspacePath,
		taskID,
		description,
		"in_progress",
		[]string{},
		steps,
		[]string{},
		initNotes,
	)
}

func CompleteTaskStep(workspacePath, taskID, stepName string, activeFiles []string, notes string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow("SELECT description, status, completed_steps, next_steps, active_files, notes FROM checkpoints WHERE task_id = ?", taskID)
	var description, status, completedStepsStr, nextStepsStr, activeFilesStr, currentNotes string
	if err := row.Scan(&description, &status, &completedStepsStr, &nextStepsStr, &activeFilesStr, &currentNotes); err != nil {
		return fmt.Sprintf("❌ Task '%s' not found.", taskID), nil
	}

	var comp []string
	var nxt []string
	var currActiveFiles []string
	json.Unmarshal([]byte(completedStepsStr), &comp)
	json.Unmarshal([]byte(nextStepsStr), &nxt)
	if activeFilesStr != "" {
		json.Unmarshal([]byte(activeFilesStr), &currActiveFiles)
	}

	stepFound := false
	var newNxt []string
	for _, s := range nxt {
		if s == stepName {
			stepFound = true
			comp = append(comp, s)
		} else {
			newNxt = append(newNxt, s)
		}
	}

	if !stepFound {
		return fmt.Sprintf("⚠️ Step '%s' not in queue.", stepName), nil
	}

	stat := "completed"
	if len(newNxt) > 0 {
		stat = status
	}

	endTimeStr := time.Now().Format("15:04:05")
	log := currentNotes + fmt.Sprintf("\n[%s] ✅ Done: %s", endTimeStr, stepName)

	if len(activeFiles) > 0 {
		var activeFilesLog []string
		for _, f := range activeFiles {
			activeFilesLog = append(activeFilesLog, f)
			found := false
			for _, cf := range currActiveFiles {
				if cf == f {
					found = true
					break
				}
			}
			if !found {
				currActiveFiles = append(currActiveFiles, f)
			}
		}
		fromJSON, _ := json.Marshal(activeFilesLog)
		log += fmt.Sprintf("\n  - Files: %s", string(fromJSON))
	}

	if notes != "" {
		log += fmt.Sprintf("\n  - Notes: %s", notes)
	}

	return SaveCheckpoint(workspacePath, taskID, description, stat, comp, newNxt, currActiveFiles, log)
}

func AddTaskStep(workspacePath, taskID, newStep string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow("SELECT description, status, completed_steps, next_steps, active_files, notes FROM checkpoints WHERE task_id = ?", taskID)
	var description, status, completedStepsStr, nextStepsStr, activeFilesStr, currentNotes string
	if err := row.Scan(&description, &status, &completedStepsStr, &nextStepsStr, &activeFilesStr, &currentNotes); err != nil {
		return fmt.Sprintf("❌ Task '%s' not found.", taskID), nil
	}

	var comp []string
	var nxt []string
	var currActiveFiles []string
	json.Unmarshal([]byte(completedStepsStr), &comp)
	json.Unmarshal([]byte(nextStepsStr), &nxt)
	if activeFilesStr != "" {
		json.Unmarshal([]byte(activeFilesStr), &currActiveFiles)
	}

	for _, s := range nxt {
		if s == newStep {
			return fmt.Sprintf("⚠️ Step '%s' already exists in task '%s'.", newStep, taskID), nil
		}
	}
	for _, s := range comp {
		if s == newStep {
			return fmt.Sprintf("⚠️ Step '%s' already exists in task '%s'.", newStep, taskID), nil
		}
	}

	nxt = append(nxt, newStep)
	log := currentNotes + fmt.Sprintf("\n[%s] Added new step: %s", time.Now().Format("15:04:05"), newStep)

	return SaveCheckpoint(workspacePath, taskID, description, status, comp, nxt, currActiveFiles, log)
}

func LoadCheckpoint(workspacePath, taskID string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow("SELECT status, updated_at, completed_steps, next_steps, notes FROM checkpoints WHERE task_id = ?", taskID)
	var status, updatedAt, completedStepsStr, nextStepsStr, notes string
	if err := row.Scan(&status, &updatedAt, &completedStepsStr, &nextStepsStr, &notes); err != nil {
		return fmt.Sprintf("❌ Checkpoint '%s' not found.", taskID), nil
	}

	var comp []string
	var nxt []string
	json.Unmarshal([]byte(completedStepsStr), &comp)
	json.Unmarshal([]byte(nextStepsStr), &nxt)

	total := len(comp) + len(nxt)
	pct := float64(0)
	if total > 0 {
		pct = float64(len(comp)) / float64(total) * 100.0
	}

	res := fmt.Sprintf("🔄 %s [%s]\nLast Update: %s\n\nProgress: %.1f%% (%d/%d steps)\n\n## ✅ Completed\n", taskID, strings.ToUpper(status), updatedAt, pct, len(comp), total)
	for _, s := range comp {
		res += fmt.Sprintf("- [x] %s\n", s)
	}

	res += "\n## ⏳ Next\n"
	for _, s := range nxt {
		res += fmt.Sprintf("- [ ] %s\n", s)
	}

	res += fmt.Sprintf("\n## 📝 Notes\n%s", notes)
	return res, nil
}

func stringsToUpper(s string) string {
	return strings.ToUpper(s)
}

func ListActiveTasks(workspacePath string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	rows, err := db.Query("SELECT task_id, status, updated_at FROM checkpoints ORDER BY updated_at DESC")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var res []string
	res = append(res, "📋 WORKSPACE TASKS:")
	count := 0
	for rows.Next() {
		var tID, stat, updated string
		if err := rows.Scan(&tID, &stat, &updated); err == nil {
			res = append(res, fmt.Sprintf("- **%s** (%s) - %s", tID, stat, updated))
			count++
		}
	}

	if count == 0 {
		return "No tasks found.", nil
	}

	return strings.Join(res, "\n"), nil
}
