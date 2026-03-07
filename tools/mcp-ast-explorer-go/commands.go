package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type FileTask struct {
	Folder   string
	Filepath string
	Filename string
	Lang     string
}

type FileResult struct {
	Filename string
	Nodes    []NodeResult
}

func GetProjectArchitecture(workspacePath, subPath string, maxFiles int, includeDocs bool) (string, error) {
	baseDir := workspacePath
	if subPath != "" {
		baseDir = filepath.Join(workspacePath, subPath)
	}

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return fmt.Sprintf("❌ Path not found: %s", baseDir), nil
	}

	allFiles := getProjectFiles(workspacePath)
	mainFamily := getMainLanguageFamily(allFiles)

	folderToFiles := make(map[string][]FileTask)
	for _, f := range allFiles {
		if !strings.HasPrefix(f, baseDir) {
			continue
		}

		ext := filepath.Ext(f)
		family := getFamilyFromExt(ext)
		if family != mainFamily || family == "" {
			continue
		}

		relPath, err := filepath.Rel(workspacePath, f)
		if err != nil {
			relPath = f
		}

		folder := filepath.Dir(relPath)
		filename := filepath.Base(relPath)

		lang := getLanguageFromExt(ext)
		if lang == "" {
			continue
		}

		folderToFiles[folder] = append(folderToFiles[folder], FileTask{
			Folder:   folder,
			Filepath: f,
			Filename: filename,
			Lang:     lang,
		})
	}

	var sortedFolders []string
	for folder := range folderToFiles {
		sortedFolders = append(sortedFolders, folder)
	}
	sort.Strings(sortedFolders)

	var tasks []FileTask
	for _, folder := range sortedFolders {
		filesInFolder := folderToFiles[folder]
		sort.Slice(filesInFolder, func(i, j int) bool {
			return filesInFolder[i].Filepath < filesInFolder[j].Filepath
		})

		for _, task := range filesInFolder {
			tasks = append(tasks, task)
			if len(tasks) >= maxFiles {
				break
			}
		}
		if len(tasks) >= maxFiles {
			break
		}
	}

	subPathMsg := subPath
	if subPathMsg == "" {
		subPathMsg = "ROOT"
	}

	var output []string
	output = append(output, fmt.Sprintf("🏗 PROJECT ARCHITECTURE AST: %s (Main Lang: %s)\n", subPathMsg, mainFamily))

	// Process using a bounded thread pool / goroutines
	var wg sync.WaitGroup
	resultCh := make(chan struct {
		Folder string
		Res    FileResult
	}, len(tasks))

	sem := make(chan struct{}, 8) // Limit concurrency to 8

	for _, task := range tasks {
		wg.Add(1)
		go func(t FileTask) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			_, nodes := parseAndExtract(t.Filepath, workspacePath, t.Lang)
			if len(nodes) > 0 {
				resultCh <- struct {
					Folder string
					Res    FileResult
				}{t.Folder, FileResult{Filename: t.Filename, Nodes: nodes}}
			}
		}(task)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	resultsByFolder := make(map[string][]FileResult)
	for res := range resultCh {
		resultsByFolder[res.Folder] = append(resultsByFolder[res.Folder], res.Res)
	}

	for _, folder := range sortedFolders {
		folderResults, ok := resultsByFolder[folder]
		if !ok {
			continue
		}

		displayFolder := folder
		if displayFolder == "" || displayFolder == "." {
			displayFolder = "."
		}
		output = append(output, fmt.Sprintf("📁 %s", displayFolder))

		sort.Slice(folderResults, func(i, j int) bool {
			return folderResults[i].Filename < folderResults[j].Filename
		})

		for _, fileRes := range folderResults {
			output = append(output, fmt.Sprintf("  📄 %s", fileRes.Filename))
			for _, n := range fileRes.Nodes {
				indent := "    " + strings.Repeat("  ", n.Level)
				sig := ""
				if n.Signature != "" {
					sig = " " + n.Signature
				}
				output = append(output, fmt.Sprintf("%s▪ [%s] %s%s", indent, n.Type, n.Name, sig))
				if includeDocs && n.Doc != "" {
					firstLine := strings.Split(n.Doc, "\n")[0]
					output = append(output, fmt.Sprintf("%s  // %s", indent, firstLine))
				}
			}
		}
	}

	if len(tasks) >= maxFiles {
		output = append(output, fmt.Sprintf("\n⚠️ Reached limit of %d files.", maxFiles))
	}

	return strings.Join(output, "\n"), nil
}

func SearchSymbol(workspacePath, query string) (string, error) {
	allFiles := getProjectFiles(workspacePath)
	mainFamily := getMainLanguageFamily(allFiles)

	var tasks []FileTask
	for _, filepathStr := range allFiles {
		ext := filepath.Ext(filepathStr)
		family := getFamilyFromExt(ext)
		if family != mainFamily || family == "" {
			continue
		}

		lang := getLanguageFromExt(ext)
		if lang == "" {
			continue
		}

		tasks = append(tasks, FileTask{
			Filepath: filepathStr,
			Lang:     lang,
		})
	}

	var results []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 8)

	for _, task := range tasks {
		wg.Add(1)
		go func(t FileTask) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			relPath, nodes := parseAndExtract(t.Filepath, workspacePath, t.Lang)
			if len(nodes) == 0 {
				return
			}

			lowerQuery := strings.ToLower(query)
			for _, n := range nodes {
				if strings.Contains(strings.ToLower(n.Name), lowerQuery) {
					mu.Lock()
					if len(results) < 50 {
						results = append(results, fmt.Sprintf("📍 %s -> [%s] %s%s", relPath, n.Type, n.Name, n.Signature))
					}
					mu.Unlock()
				}
			}
		}(task)
	}

	wg.Wait()

	if len(results) == 0 {
		return fmt.Sprintf("🔍 No symbols matching '%s' found.", query), nil
	}

	return "🔎 SYMBOL SEARCH RESULTS:\n" + strings.Join(results, "\n"), nil
}
