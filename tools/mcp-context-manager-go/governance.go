package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func DeclareIntent(workspacePath, taskID, tactic string, lockedFiles []string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	lockedBytes, _ := json.Marshal(lockedFiles)

	query := `
        INSERT INTO intents (task_id, tactic, locked_files)
        VALUES (?, ?, ?)
        ON CONFLICT(task_id) DO UPDATE SET
            tactic=excluded.tactic,
            locked_files=excluded.locked_files,
            updated_at=CURRENT_TIMESTAMP
    `
	_, err = db.Exec(query, taskID, tactic, string(lockedBytes))
	if err != nil {
		return "", fmt.Errorf("failed to declare intent: %v", err)
	}

	return fmt.Sprintf("🔒 Intent declared. Lock applied to files: %s", strings.Join(lockedFiles, ", ")), nil
}

func CheckIntentLock(workspacePath, taskID, filePath string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var lockedFilesStr string
	err = db.QueryRow("SELECT locked_files FROM intents WHERE task_id = ?", taskID).Scan(&lockedFilesStr)
	if err != nil {
		return "⚠️ No intent declared for this task. Please declare_intent first.", nil
	}

	var lockedFiles []string
	if err := json.Unmarshal([]byte(lockedFilesStr), &lockedFiles); err != nil {
		return "", fmt.Errorf("failed to parse locked_files: %v", err)
	}

	isLocked := false
	for _, f := range lockedFiles {
		if strings.Contains(f, filePath) || strings.Contains(filePath, f) {
			isLocked = true
			break
		}
	}

	if isLocked {
		var gotchas string
		err := db.QueryRow("SELECT gotchas FROM file_annotations WHERE file_path = ?", filePath).Scan(&gotchas)
		ghostCtx := ""
		if err == nil && gotchas != "" {
			ghostCtx = fmt.Sprintf("\n👻 GHOST CONTEXT: %s", gotchas)
		}
		return fmt.Sprintf("✅ File '%s' is unlocked for current intent.%s", filePath, ghostCtx), nil
	}

	return fmt.Sprintf("❌ ALARM: Scope Creep! File '%s' is NOT in the active_files lock. Switch to Planner to update intent via declare_intent.", filePath), nil
}

func AnnotateFile(workspacePath, filePath, gotchas string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := `
        INSERT INTO file_annotations (file_path, gotchas)
        VALUES (?, ?)
        ON CONFLICT(file_path) DO UPDATE SET
            gotchas=excluded.gotchas,
            updated_at=CURRENT_TIMESTAMP
    `
	_, err = db.Exec(query, filePath, gotchas)
	if err != nil {
		return "", fmt.Errorf("failed to annotate file: %v", err)
	}

	return fmt.Sprintf("👻 Ghost Context added to '%s'.", filePath), nil
}

func RecordFailure(workspacePath, taskID string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := `
        INSERT INTO drift_tracker (task_id, failure_count, last_failed_at)
        VALUES (?, 1, CURRENT_TIMESTAMP)
        ON CONFLICT(task_id) DO UPDATE SET
            failure_count=failure_count + 1,
            last_failed_at=CURRENT_TIMESTAMP
    `
	if _, err := db.Exec(query, taskID); err != nil {
		return "", fmt.Errorf("failed to record failure: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT failure_count FROM drift_tracker WHERE task_id=?", taskID).Scan(&count); err != nil {
		return "", err
	}

	if count >= 3 {
		return fmt.Sprintf("🚨 DRIFT DETECTED (Failures: %d): You have failed 3 or more times. Trigger 'think_back' panic protocol immediately!", count), nil
	}
	return fmt.Sprintf("⚠️ Failure recorded. Count: %d", count), nil
}

func ClearDrift(workspacePath, taskID string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	if _, err := db.Exec("UPDATE drift_tracker SET failure_count=0 WHERE task_id=?", taskID); err != nil {
		return "", fmt.Errorf("failed to clear drift: %v", err)
	}

	return "🧹 Drift counter reset to 0.", nil
}

func ManageAnchors(workspacePath, action, key, value, rule string) (string, error) {
	db, err := GetDBConnection(workspacePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	switch action {
	case "set":
		if key == "" || value == "" {
			return "❌ Key and value required for 'set'.", nil
		}
		query := `
            INSERT INTO anchors (key, value, rule)
            VALUES (?, ?, ?)
            ON CONFLICT(key) DO UPDATE SET
                value=excluded.value,
                rule=excluded.rule,
                updated_at=CURRENT_TIMESTAMP
        `
		if _, err := db.Exec(query, key, value, rule); err != nil {
			return "", err
		}
		return fmt.Sprintf("⚓ Anchor '%s' secured successfully.", key), nil

	case "get":
		var val, r string
		err := db.QueryRow("SELECT value, rule FROM anchors WHERE key = ?", key).Scan(&val, &r)
		if err != nil {
			return fmt.Sprintf("⚠️ Anchor '%s' not found.", key), nil
		}
		return fmt.Sprintf("⚓ ANCHOR [%s]: %s (Rule: %s)", key, val, r), nil

	case "list":
		rows, err := db.Query("SELECT key, value, rule FROM anchors ORDER BY key")
		if err != nil {
			return "", err
		}
		defer rows.Close()

		var res []string
		res = append(res, "⚓ PROJECT ANCHORS:")
		count := 0
		for rows.Next() {
			var k, v, r string
			if err := rows.Scan(&k, &v, &r); err == nil {
				res = append(res, fmt.Sprintf("- **%s**: %s", k, v))
				if r != "" {
					res = append(res, fmt.Sprintf("  Rule: %s", r))
				}
				count++
			}
		}
		if count == 0 {
			return "No anchors defined.", nil
		}
		return strings.Join(res, "\n"), nil

	default:
		return "❌ Unknown action. Use 'set', 'get', or 'list'.", nil
	}
}
