package tools

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ToolResult map[string]interface{}

type ListDirectoryParams struct {
	Path string `json:"path"`
}

func findAllowedRoot(allowedDirs []string, rel string) (string, error) {
	cleanRel := filepath.Clean(rel)
	if filepath.IsAbs(cleanRel) {
		absRel := cleanRel
		for _, root := range allowedDirs {
			cleanRoot := filepath.Clean(root)
			if absRel == cleanRoot || strings.HasPrefix(absRel+string(os.PathSeparator), cleanRoot+string(os.PathSeparator)) || strings.HasPrefix(absRel, cleanRoot+string(os.PathSeparator)) {
				return absRel, nil
			}
		}
		return "", errors.New("access outside of allowed directories is not allowed")
	}
	for _, root := range allowedDirs {
		abs, err := secureJoin(root, rel)
		if err == nil {
			return abs, nil
		}
	}
	return "", errors.New("access outside of allowed directories is not allowed")
}

func ListDirectory(params ListDirectoryParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}
	var result []map[string]string
	for _, entry := range entries {
		typeStr := "file"
		if entry.IsDir() {
			typeStr = "directory"
		}
		result = append(result, map[string]string{
			"name": entry.Name(),
			"type": typeStr,
		})
	}
	return ToolResult{"entries": result}, nil
}

type ReadFileParams struct {
	Path string `json:"path"`
}

func ReadFile(params ReadFileParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return ToolResult{"content": string(data)}, nil
}

type WriteFileParams struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func WriteFile(params WriteFileParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(absPath, []byte(params.Content), 0644)
	if err != nil {
		return nil, err
	}
	return ToolResult{"ok": true}, nil
}

type CreateDirectoryParams struct {
	Path string `json:"path"`
}

func CreateDirectory(params CreateDirectoryParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return nil, err
	}
	return ToolResult{"ok": true}, nil
}

type GetFileInfoParams struct {
	Path string `json:"path"`
}

func GetFileInfo(params GetFileInfoParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	var creationTime, accessTime string
	// platform-specific: for compatibility only, on Unix you can only get modTime
	creationTime = info.ModTime().Format("2006-01-02T15:04:05Z07:00")
	accessTime = info.ModTime().Format("2006-01-02T15:04:05Z07:00")
	perms := info.Mode().Perm().String()
	return ToolResult{
		"size":         info.Size(),
		"mode":         info.Mode().String(),
		"modTime":      info.ModTime().Format("2006-01-02T15:04:05Z07:00"),
		"isDir":        info.IsDir(),
		"name":         info.Name(),
		"creationTime": creationTime,
		"accessTime":   accessTime,
		"permissions":  perms,
	}, nil
}

type MoveFileParams struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func MoveFile(params MoveFileParams, allowedDirs []string) (ToolResult, error) {
	src, err := findAllowedRoot(allowedDirs, params.Source)
	if err != nil {
		return nil, err
	}
	dst, err := findAllowedRoot(allowedDirs, params.Destination)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(dst); err == nil {
		return nil, errors.New("destination already exists")
	}
	err = os.Rename(src, dst)
	if err != nil {
		return nil, err
	}
	return ToolResult{"ok": true}, nil
}

type DeleteFileParams struct {
	Path string `json:"path"`
}

func DeleteFile(params DeleteFileParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	if _, statErr := os.Stat(absPath); os.IsNotExist(statErr) {
		return nil, errors.New("file does not exist")
	}
	err = os.RemoveAll(absPath)
	if err != nil {
		return nil, err
	}
	return ToolResult{"ok": true}, nil
}

type SearchFilesParams struct {
	Path            string   `json:"path"`
	Pattern         string   `json:"pattern"`
	ExcludePatterns []string `json:"excludePatterns"`
}

func SearchFiles(params SearchFilesParams, allowedDirs []string) (ToolResult, error) {
	startDir, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return ToolResult{"content": []map[string]interface{}{{"type": "text", "text": "Error: " + err.Error()}}, "isError": true}, nil
	}
	var matches []string
	excludes := params.ExcludePatterns
	isExcluded := func(rel string) bool {
		for _, ex := range excludes {
			if ok, _ := filepath.Match(ex, rel); ok {
				return true
			}
		}
		return false
	}
	err = filepath.WalkDir(startDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // пропускаем ошибки доступа
		}
		allowed := false
		for _, root := range allowedDirs {
			if strings.HasPrefix(path, filepath.Clean(root)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil
		}
		rel, _ := filepath.Rel(startDir, path)
		matched, _ := filepath.Match(params.Pattern, d.Name())
		if matched && !isExcluded(rel) {
			matches = append(matches, rel)
		}
		return nil
	})
	if err != nil {
		return ToolResult{"content": []map[string]interface{}{{"type": "text", "text": "Error: " + err.Error()}}, "isError": true}, nil
	}
	var text string
	if len(matches) > 0 {
		text = strings.Join(matches, "\n")
	} else {
		text = "No matches found"
	}
	return ToolResult{"content": []map[string]interface{}{{"type": "text", "text": text}}}, nil
}

type ReadMultipleFilesParams struct {
	Paths []string `json:"paths"`
}

func ReadMultipleFiles(params ReadMultipleFilesParams, allowedDirs []string) (ToolResult, error) {
	var results []string
	for _, p := range params.Paths {
		absPath, err := findAllowedRoot(allowedDirs, p)
		if err != nil {
			results = append(results, p+": Error - "+err.Error())
			continue
		}
		data, err := os.ReadFile(absPath)
		if err != nil {
			results = append(results, p+": Error - "+err.Error())
			continue
		}
		results = append(results, p+":\n"+string(data))
	}
	joined := strings.Join(results, "\n---\n")
	return ToolResult{"content": []map[string]interface{}{{"type": "text", "text": joined}}}, nil
}

func ListAllowedDirectories(allowedDirs []string) (ToolResult, error) {
	return ToolResult{"directories": allowedDirs}, nil
}

type EditFileParams struct {
	Path  string `json:"path"`
	Edits []struct {
		OldText string `json:"oldText"`
		NewText string `json:"newText"`
	} `json:"edits"`
	DryRun bool `json:"dryRun"`
}

func EditFile(params EditFileParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	origData, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(origData), "\n")
	newLines := make([]string, len(lines))
	copy(newLines, lines)
	var diff bytes.Buffer
	changed := false

	for _, edit := range params.Edits {
		old := edit.OldText
		newT := edit.NewText
		match := false
		for i, line := range newLines {
			if strings.Contains(line, old) {
				match = true
				before := newLines[i]
				newLines[i] = strings.ReplaceAll(line, old, newT)
				if before != newLines[i] {
					changed = true
					diff.WriteString(fmt.Sprintf("-%s\n+%s\n", before, newLines[i]))
				}
			}
		}
		if !match {
			diff.WriteString(fmt.Sprintf("~ not found: %q\n", old))
		}
	}

	result := ToolResult{
		"diff":    diff.String(),
		"changed": changed,
	}

	if params.DryRun {
		result["preview"] = strings.Join(newLines, "\n")
		return result, nil
	}

	if changed {
		joined := strings.Join(newLines, "\n")
		err = os.WriteFile(absPath, []byte(joined), 0644)
		if err != nil {
			return nil, err
		}
		result["ok"] = true
	} else {
		result["ok"] = false
	}
	return result, nil
}

type ListDirectoryWithSizesParams struct {
	Path   string `json:"path"`
	SortBy string `json:"sortBy"`
}

func ListDirectoryWithSizes(params ListDirectoryWithSizesParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	type entryInfo struct {
		Name  string
		IsDir bool
		Size  int64
	}
	var infos []entryInfo
	for _, entry := range entries {
		info := entryInfo{Name: entry.Name(), IsDir: entry.IsDir(), Size: 0}
		if !entry.IsDir() {
			stat, err := os.Stat(filepath.Join(absPath, entry.Name()))
			if err == nil {
				info.Size = stat.Size()
			}
		}
		infos = append(infos, info)
	}
	// Сортировка
	sortBy := params.SortBy
	if sortBy != "size" {
		sortBy = "name"
	}
	if sortBy == "size" {
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Size > infos[j].Size
		})
	} else {
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Name < infos[j].Name
		})
	}
	var result []map[string]interface{}
	var totalFiles, totalDirs int
	var totalSize int64
	for _, info := range infos {
		typeStr := "file"
		if info.IsDir {
			typeStr = "directory"
			totalDirs++
		} else {
			totalFiles++
			totalSize += info.Size
		}
		result = append(result, map[string]interface{}{
			"name": info.Name,
			"type": typeStr,
			"size": info.Size,
		})
	}
	return ToolResult{
		"entries":    result,
		"totalFiles": totalFiles,
		"totalDirs":  totalDirs,
		"totalSize":  totalSize,
	}, nil
}

type DirectoryTreeParams struct {
	Path string `json:"path"`
}

type TreeEntry struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"` // "file" или "directory"
	Children []TreeEntry `json:"children,omitempty"`
}

func DirectoryTree(params DirectoryTreeParams, allowedDirs []string) (ToolResult, error) {
	absPath, err := findAllowedRoot(allowedDirs, params.Path)
	if err != nil {
		return nil, err
	}
	entry, err := buildTree(absPath)
	if err != nil {
		return nil, err
	}
	return ToolResult{"tree": entry}, nil
}

func buildTree(path string) (TreeEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return TreeEntry{}, err
	}
	entry := TreeEntry{Name: info.Name()}
	if info.IsDir() {
		entry.Type = "directory"
		files, err := os.ReadDir(path)
		if err != nil {
			return entry, err
		}
		for _, f := range files {
			child, err := buildTree(filepath.Join(path, f.Name()))
			if err == nil {
				entry.Children = append(entry.Children, child)
			}
		}
	} else {
		entry.Type = "file"
	}
	return entry, nil
}

func secureJoin(root, rel string) (string, error) {
	cleanRoot := filepath.Clean(root)
	joined := filepath.Join(cleanRoot, rel)
	abs, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	if len(abs) < len(cleanRoot) || abs[:len(cleanRoot)] != cleanRoot {
		return "", errors.New("access outside of root directory is not allowed")
	}
	return abs, nil
}
