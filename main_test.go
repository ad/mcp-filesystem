package main

import (
	"os"
	"testing"

	"github.com/ad/mcp-filesystem/tools"
)

func TestListDirectory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/foo.txt", []byte("abc"), 0644)
	res, err := tools.ListDirectory(tools.ListDirectoryParams{Path: "."}, []string{dir})
	if err != nil {
		t.Fatalf("ListDirectory error: %v", err)
	}
	if len(res["entries"].([]map[string]string)) == 0 {
		t.Error("ListDirectory: entries should not be empty")
	}
}

func TestReadWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := "file.txt"
	_, err := tools.WriteFile(tools.WriteFileParams{Path: path, Content: "hello"}, []string{dir})
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	res, err := tools.ReadFile(tools.ReadFileParams{Path: path}, []string{dir})
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if res["content"].(string) != "hello" {
		t.Errorf("ReadFile: expected 'hello', got '%s'", res["content"])
	}
}

func TestCreateDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := tools.CreateDirectory(tools.CreateDirectoryParams{Path: "foo/bar"}, []string{dir})
	if err != nil {
		t.Fatalf("CreateDirectory error: %v", err)
	}
	if _, err := os.Stat(dir + "/foo/bar"); err != nil {
		t.Error("Directory not created")
	}
}

func TestGetFileInfo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/f.txt", []byte("x"), 0644)
	res, err := tools.GetFileInfo(tools.GetFileInfoParams{Path: "f.txt"}, []string{dir})
	if err != nil {
		t.Fatalf("GetFileInfo error: %v", err)
	}
	if res["name"] != "f.txt" {
		t.Errorf("GetFileInfo: wrong name: %v", res["name"])
	}
}

func TestMoveFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.txt", []byte("1"), 0644)
	_, err := tools.MoveFile(tools.MoveFileParams{Source: "a.txt", Destination: "b.txt"}, []string{dir})
	if err != nil {
		t.Fatalf("MoveFile error: %v", err)
	}
	if _, err := os.Stat(dir + "/b.txt"); err != nil {
		t.Error("File not moved")
	}
}

func TestDeleteFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/del.txt", []byte("x"), 0644)
	_, err := tools.DeleteFile(tools.DeleteFileParams{Path: "del.txt"}, []string{dir})
	if err != nil {
		t.Fatalf("DeleteFile error: %v", err)
	}
	if _, err := os.Stat(dir + "/del.txt"); err == nil {
		t.Error("File not deleted")
	}
}

func TestSearchFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.go", []byte("x"), 0644)
	os.WriteFile(dir+"/b.txt", []byte("y"), 0644)
	res, err := tools.SearchFiles(tools.SearchFilesParams{Path: ".", Pattern: "*.go", ExcludePatterns: nil}, []string{dir})
	if err != nil {
		t.Fatalf("SearchFiles error: %v", err)
	}
	matches := res["matches"].([]string)
	if len(matches) == 0 {
		t.Error("SearchFiles: no matches")
	}
}

func TestReadMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.txt", []byte("A"), 0644)
	os.WriteFile(dir+"/b.txt", []byte("B"), 0644)
	res, err := tools.ReadMultipleFiles(tools.ReadMultipleFilesParams{Paths: []string{"a.txt", "b.txt", "no.txt"}}, []string{dir})
	if err != nil {
		t.Fatalf("ReadMultipleFiles error: %v", err)
	}
	results := res["results"].(map[string]interface{})
	if results["a.txt"].(string) != "A" || results["b.txt"].(string) != "B" {
		t.Error("ReadMultipleFiles: wrong content")
	}
	if _, ok := results["no.txt"].(map[string]interface{}); !ok {
		t.Error("ReadMultipleFiles: missing error for no.txt")
	}
}

func TestListAllowedDirectories(t *testing.T) {
	dir := t.TempDir()
	res, err := tools.ListAllowedDirectories([]string{dir})
	if err != nil {
		t.Fatalf("ListAllowedDirectories error: %v", err)
	}
	dirs := res["directories"].([]string)
	if len(dirs) == 0 || dirs[0] != dir {
		t.Error("ListAllowedDirectories: wrong result")
	}
}

func TestListDirectoryWithSizes(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.txt", []byte("abc"), 0644)
	os.WriteFile(dir+"/b.txt", []byte("defg"), 0644)
	os.Mkdir(dir+"/subdir", 0755)
	res, err := tools.ListDirectoryWithSizes(tools.ListDirectoryWithSizesParams{Path: ".", SortBy: "size"}, []string{dir})
	if err != nil {
		t.Fatalf("ListDirectoryWithSizes error: %v", err)
	}
	entries := res["entries"].([]map[string]interface{})
	if len(entries) < 2 {
		t.Error("ListDirectoryWithSizes: entries should not be empty")
	}
	if res["totalFiles"].(int) < 2 {
		t.Error("ListDirectoryWithSizes: totalFiles should be >= 2")
	}
}

func TestDirectoryTree(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.txt", []byte("abc"), 0644)
	os.Mkdir(dir+"/subdir", 0755)
	os.WriteFile(dir+"/subdir/b.txt", []byte("def"), 0644)
	res, err := tools.DirectoryTree(tools.DirectoryTreeParams{Path: "."}, []string{dir})
	if err != nil {
		t.Fatalf("DirectoryTree error: %v", err)
	}
	tree := res["tree"].(tools.TreeEntry)
	if tree.Type != "directory" {
		t.Error("DirectoryTree: root should be directory")
	}
	found := false
	for _, child := range tree.Children {
		if child.Name == "subdir" && child.Type == "directory" {
			found = true
		}
	}
	if !found {
		t.Error("DirectoryTree: subdir not found in tree")
	}
}
