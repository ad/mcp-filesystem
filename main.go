package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ad/mcp-filesystem/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func makeHandleListDirectory(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] list_directory: %v", request.Params.Arguments)
		var params tools.ListDirectoryParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] list_directory: %v", err)
			return nil, err
		}
		res, err := tools.ListDirectory(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] list_directory: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleReadFile(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] read_file: %v", request.Params.Arguments)
		var params tools.ReadFileParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] read_file: %v", err)
			return nil, err
		}
		res, err := tools.ReadFile(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] read_file: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleWriteFile(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] write_file: %v", request.Params.Arguments)
		var params tools.WriteFileParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] write_file: %v", err)
			return nil, err
		}
		res, err := tools.WriteFile(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] write_file: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleCreateDirectory(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] create_directory: %v", request.Params.Arguments)
		var params tools.CreateDirectoryParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] create_directory: %v", err)
			return nil, err
		}
		res, err := tools.CreateDirectory(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] create_directory: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleGetFileInfo(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] get_file_info: %v", request.Params.Arguments)
		var params tools.GetFileInfoParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] get_file_info: %v", err)
			return nil, err
		}
		res, err := tools.GetFileInfo(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] get_file_info: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleMoveFile(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] move_file: %v", request.Params.Arguments)
		var params tools.MoveFileParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] move_file: %v", err)
			return nil, err
		}
		res, err := tools.MoveFile(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] move_file: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleDeleteFile(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] delete_file: %v", request.Params.Arguments)
		var params tools.DeleteFileParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] delete_file: %v", err)
			return nil, err
		}
		res, err := tools.DeleteFile(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] delete_file: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleSearchFiles(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] search_files: %v", request.Params.Arguments)
		var params tools.SearchFilesParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] search_files: %v", err)
			return nil, err
		}
		res, err := tools.SearchFiles(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] search_files: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleReadMultipleFiles(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] read_multiple_files: %v", request.Params.Arguments)
		var params tools.ReadMultipleFilesParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] read_multiple_files: %v", err)
			return nil, err
		}
		res, err := tools.ReadMultipleFiles(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] read_multiple_files: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleListAllowedDirectories(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] list_allowed_directories")
		res, err := tools.ListAllowedDirectories(allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] list_allowed_directories: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleEditFile(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] edit_file: %v", request.Params.Arguments)
		var params tools.EditFileParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] edit_file: %v", err)
			return nil, err
		}
		res, err := tools.EditFile(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] edit_file: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleListDirectoryWithSizes(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] list_directory_with_sizes: %v", request.Params.Arguments)
		var params tools.ListDirectoryWithSizesParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] list_directory_with_sizes: %v", err)
			return nil, err
		}
		res, err := tools.ListDirectoryWithSizes(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] list_directory_with_sizes: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

func makeHandleDirectoryTree(allowedDirs []string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("[MCP] directory_tree: %v", request.Params.Arguments)
		var params tools.DirectoryTreeParams
		if err := decodeParams(request.Params.Arguments, &params); err != nil {
			log.Printf("[MCP][ERROR] directory_tree: %v", err)
			return nil, err
		}
		res, err := tools.DirectoryTree(params, allowedDirs)
		if err != nil {
			log.Printf("[MCP][ERROR] directory_tree: %v", err)
			return nil, err
		}
		return wrapResult(res), nil
	}
}

// decodeParams — универсальный декодер параметров MCP tool
func decodeParams(args interface{}, out interface{}) error {
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

// wrapResult — оборачивает ToolResult в MCP CallToolResult
func wrapResult(res tools.ToolResult) *mcp.CallToolResult {
	// Если есть ключ isError, выставляем его явно
	isError := false
	if v, ok := res["isError"]; ok {
		if b, ok := v.(bool); ok && b {
			isError = true
		}
	}
	b, _ := json.Marshal(res)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(b),
			},
		},
		IsError: isError,
	}
}

func main() {
	var transport = flag.String("transport", "stdio", "Transport type: stdio, sse, or http")
	var port = flag.String("port", "8080", "Port for SSE/HTTP servers")
	flag.Parse()

	allowedDirs := flag.Args()
	if len(allowedDirs) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-transport stdio|sse|http] [-port PORT] <allowed-directory> [additional-directories...]\n", os.Args[0])
		os.Exit(1)
	}

	mcpServer := server.NewMCPServer(
		"filesystem",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Регистрация MCP tools (файловая система) с поддержкой allowedDirs
	mcpServer.AddTool(
		mcp.NewTool("list_directory",
			mcp.WithDescription("Get a detailed listing of all files and directories in a specified path. "+
				"Results clearly distinguish between files and directories with [FILE] and [DIR] "+
				"prefixes. This tool is essential for understanding directory structure and "+
				"finding specific files within a directory. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("Directory path"), mcp.Required()),
		),
		makeHandleListDirectory(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("read_file",
			mcp.WithDescription("Read the complete contents of a file from the file system. "+
				"Handles various text encodings and provides detailed error messages "+
				"if the file cannot be read. Use this tool when you need to examine "+
				"the contents of a single file. Use the 'head' parameter to read only "+
				"the first N lines of a file, or the 'tail' parameter to read only "+
				"the last N lines of a file. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("File path"), mcp.Required()),
		),
		makeHandleReadFile(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("write_file",
			mcp.WithDescription("Create a new file or completely overwrite an existing file with new content. "+
				"Use with caution as it will overwrite existing files without warning. "+
				"Handles text content with proper encoding. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("File path"), mcp.Required()),
			mcp.WithString("content", mcp.Description("File content"), mcp.Required()),
		),
		makeHandleWriteFile(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("create_directory",
			mcp.WithDescription("Create a new directory or ensure a directory exists. Can create multiple "+
				"nested directories in one operation. If the directory already exists, "+
				"this operation will succeed silently. Perfect for setting up directory "+
				"structures for projects or ensuring required paths exist. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("Directory path"), mcp.Required()),
		),
		makeHandleCreateDirectory(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("get_file_info",
			mcp.WithDescription("Retrieve detailed metadata about a file or directory. Returns comprehensive "+
				"information including size, creation time, last modified time, permissions, "+
				"and type. This tool is perfect for understanding file characteristics "+
				"without reading the actual content. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("Path"), mcp.Required()),
		),
		makeHandleGetFileInfo(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("move_file",
			mcp.WithDescription("Move or rename files and directories. Can move files between directories "+
				"and rename them in a single operation. If the destination exists, the "+
				"operation will fail. Works across different directories and can be used "+
				"for simple renaming within the same directory. Both source and destination must be within allowed directories."),
			mcp.WithString("source", mcp.Description("Source path"), mcp.Required()),
			mcp.WithString("destination", mcp.Description("Destination path"), mcp.Required()),
		),
		makeHandleMoveFile(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("delete_file",
			mcp.WithDescription("Delete file or directory"),
			mcp.WithString("path", mcp.Description("Path to delete"), mcp.Required()),
		),
		makeHandleDeleteFile(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("search_files",
			mcp.WithDescription("Recursively search for files and directories matching a pattern. "+
				"Searches through all subdirectories from the starting path. The search "+
				"is case-insensitive and matches partial names. Returns full paths to all "+
				"matching items. Great for finding files when you don't know their exact location. "+
				"Only searches within allowed directories."),
			mcp.WithString("path", mcp.Description("Start directory"), mcp.Required()),
			mcp.WithString("pattern", mcp.Description("Glob pattern"), mcp.Required()),
			mcp.WithArray("excludePatterns", mcp.Items(map[string]any{"type": "string"})),
		),
		makeHandleSearchFiles(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("read_multiple_files",
			mcp.WithDescription("Read the contents of multiple files simultaneously. This is more "+
				"efficient than reading files one by one when you need to analyze "+
				"or compare multiple files. Each file's content is returned with its "+
				"path as a reference. Failed reads for individual files won't stop "+
				"the entire operation. Only works within allowed directories."),
			mcp.WithArray("paths", mcp.Items(map[string]any{"type": "string"})),
		),
		makeHandleReadMultipleFiles(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("edit_file",
			mcp.WithDescription("Make line-based edits to a text file. Each edit replaces exact line sequences "+
				"with new content. Returns a git-style diff showing the changes made. "+
				"Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("File to edit"), mcp.Required()),
			mcp.WithArray("edits", mcp.Items(map[string]any{"type": "string"})),
			mcp.WithBoolean("dryRun", mcp.Description("Preview changes without applying")),
		),
		makeHandleEditFile(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("list_allowed_directories",
			mcp.WithDescription("Returns the list of directories that this server is allowed to access. "+
				"Use this to understand which directories are available before trying to access files."),
		),
		makeHandleListAllowedDirectories(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("list_directory_with_sizes",
			mcp.WithDescription("Get a detailed listing of all files and directories in a specified path, including sizes. Results clearly distinguish between files and directories with [FILE] and [DIR] prefixes. This tool is useful for understanding directory structure and finding specific files within a directory. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("Directory path"), mcp.Required()),
			mcp.WithString("sortBy", mcp.Description("Sort entries by name or size (name|size)")),
		),
		makeHandleListDirectoryWithSizes(allowedDirs),
	)
	mcpServer.AddTool(
		mcp.NewTool("directory_tree",
			mcp.WithDescription("Get a recursive tree view of files and directories as a JSON structure. Each entry includes 'name', 'type' (file/directory), and 'children' for directories. Files have no children array, while directories always have a children array (which may be empty). The output is formatted with 2-space indentation for readability. Only works within allowed directories."),
			mcp.WithString("path", mcp.Description("Directory path"), mcp.Required()),
		),
		makeHandleDirectoryTree(allowedDirs),
	)

	switch *transport {
	case "stdio":
		log.Println("Starting MCP server with STDIO transport...")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatal("STDIO server error:", err)
		}

	case "sse":
		log.Printf("Starting MCP server with SSE transport on port %s...", *port)
		sseServer := server.NewSSEServer(mcpServer)

		http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
			sseServer.ServeHTTP(w, r)
		})

		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatal("SSE server error:", err)
		}

	case "http":
		log.Printf("Starting MCP server with streamable HTTP transport on port %s...", *port)
		httpServer := server.NewStreamableHTTPServer(mcpServer)

		log.Printf("HTTP server listening on :%s/mcp", *port)
		if err := httpServer.Start(":" + *port); err != nil {
			log.Fatal("HTTP server error:", err)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown transport: %s\n", *transport)
		fmt.Fprintf(os.Stderr, "Usage: %s [-transport stdio|sse|http] [-port PORT] <allowed-directory> [additional-directories...]\n", os.Args[0])
		os.Exit(1)
	}
}
