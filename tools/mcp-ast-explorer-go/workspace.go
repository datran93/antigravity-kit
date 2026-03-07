package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var IGNORE_DIRS = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".venv":        true,
	"venv":         true,
	"dist":         true,
	"build":        true,
	".next":        true,
	".agent":       true,
}

func getProjectFiles(workspacePath string) []string {
	// Try git first
	cmd := exec.Command("git", "-C", workspacePath, "ls-files", "--cached", "--others", "--exclude-standard")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		var files []string
		for _, line := range lines {
			if line != "" {
				files = append(files, filepath.Join(workspacePath, line))
			}
		}
		if len(files) > 0 {
			return files
		}
	}

	// Fallback to walk
	var files []string
	filepath.Walk(workspacePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if IGNORE_DIRS[info.Name()] || strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files
}

func getMainLanguageFamily(files []string) string {
	familyCounts := make(map[string]int)

	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f))
		switch ext {
		case ".py":
			familyCounts["python"]++
		case ".go":
			familyCounts["go"]++
		case ".ts", ".cts", ".mts", ".tsx":
			familyCounts["typescript"]++
		case ".js", ".jsx", ".cjs", ".mjs":
			familyCounts["javascript"]++
		}
	}

	maxCount := 0
	mainFamily := ""
	for fam, count := range familyCounts {
		if count > maxCount {
			maxCount = count
			mainFamily = fam
		}
	}

	return mainFamily
}
