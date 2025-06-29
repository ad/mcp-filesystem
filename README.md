# mcp-filesystem

> **Note:** This is a Go port of the official [TypeScript version](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem) of the MCP filesystem server.

**Go implementation of MCP server filesystem** — fast, secure, and fully compatible server for the Model Context Protocol (MCP), implementing all MCP file tools with support for STDIO, HTTP, and SSE transports.

---

## 📑 Table of Contents
- [Features](#features)
- [Quick Start](#quick-start)
- [Modes and Architecture](#modes-and-architecture)
- [Integration](#integration)
- [API Reference (MCP Tools)](#api-reference-mcp-tools)
- [Testing and Automation](#testing-and-automation)
- [Project Structure](#project-structure)
- [License](#license)

---

## 🚀 Features
- **Full set of MCP tools**: list_directory, read_file, write_file, create_directory, get_file_info, move_file, delete_file, search_files, read_multiple_files, list_allowed_directories, edit_file (WIP), list_directory_with_sizes, directory_tree
- **Three transports**: STDIO (MCP), HTTP (REST), SSE (Server-Sent Events)
- **Concurrency**: parallel client handling (goroutines)
- **Access restriction**: works only in allowed directories
- **MIT license**
- **Linux and macOS support**
- **Logging to stderr only**

---

## ⚡ Quick Start

### Install via go install

To quickly install the latest version from the repository:

```fish
go install github.com/ad/mcp-filesystem@latest
```

The binary will appear in `$GOBIN` or `$HOME/go/bin` (make sure this path is in your `$PATH`).

### 1. Build from source
```fish
# Clone the repository
git clone https://github.com/ad/mcp-filesystem.git
cd mcp-filesystem

go mod tidy
# Local build
make build-local
# Or manually
go build -o mcp-filesystem main.go
# Docker build
make build
```

### 2. Run

#### STDIO (VS Code, Claude Desktop)
```fish
./mcp-filesystem -transport stdio
# or
make run-stdio
```

#### HTTP
```fish
./mcp-filesystem -transport http -port 8080
# or
make run-local
```

#### SSE
```fish
./mcp-filesystem -transport sse -port 8080
```

#### Docker
```fish
make run
# or manually
docker run --rm -p 8080:8080 danielapatin/mcp-filesystem:latest -transport http -port 8080
```

---

## 🌐 Modes and Architecture

- **STDIO** — MCP compatibility (JSON-RPC via stdin/stdout)
- **HTTP** — REST API (POST /mcp)
- **SSE** — Server-Sent Events (POST /sse)
- **Concurrency** — each client is handled in a separate goroutine
- **Access restriction** — operations only in allowed directories
- **Logging** — stderr only, same format as original

---

## 🔌 Integration

### VS Code

```
go install github.com/ad/mcp-filesystem@latest
````

Добавьте в `settings.json`:
```json
{
  "mcp": {
    "servers": {
      "mcp-filesystem": {
        "type": "stdio",
        "command": "/absolute/path/to/mcp-filesystem",
        "args": ["-transport", "stdio", "/Users/username/Desktop", "/path/to/other/allowed/dir"]
      }
    }
  }
}
```

### Docker (VS Code)
```json
{
  "mcp": {
    "servers": {
      "mcp-filesystem": {
        "type": "stdio",
        "command": "docker",
        "args": [
          "run", "--rm", "-i",
          "--mount", "type=bind,src=/Users/username/Desktop,dst=/projects/Desktop",
          "danielapatin/mcp-filesystem:latest",
          "-transport", "stdio"
        ]
      }
    }
  }
}
```

### Claude Desktop
```json
{
  "mcpServers": {
    "mcp-filesystem": {
      "command": "/absolute/path/to/mcp-filesystem",
      "args": ["-transport", "stdio", "/Users/username/Desktop", "/path/to/other/allowed/dir"]
    }
  }
}
```

### Web (SSE/HTTP)
```javascript
// SSE
const eventSource = new EventSource('http://localhost:8080/sse');
eventSource.onmessage = event => console.log(JSON.parse(event.data));

// HTTP
fetch('http://localhost:8080/mcp', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    jsonrpc: '2.0', id: 1, method: 'tools/call',
    params: { name: 'mcp-filesystem', arguments: { thought: 'Step', thoughtNumber: 1, totalThoughts: 3, nextThoughtNeeded: true } }
  })
});
```

---

## 🛠 API Reference (MCP Tools)

### list_directory
- **Input:** `{ "path": "subdir" }`
- **Output:** `{ "entries": [ { "name": "foo.txt", "type": "file" }, { "name": "bar", "type": "directory" } ] }`

### read_file
- **Input:** `{ "path": "file.txt" }`
- **Output:** `{ "content": "file contents..." }`

### write_file
- **Input:** `{ "path": "file.txt", "content": "new content" }`
- **Output:** `{ "ok": true }`

### create_directory
- **Input:** `{ "path": "newdir/subdir" }`
- **Output:** `{ "ok": true }`

### get_file_info
- **Input:** `{ "path": "file.txt" }`
- **Output:**
```json
{
  "size": 123,
  "mode": "-rw-r--r--",
  "modTime": "2025-06-29T12:00:00Z",
  "isDir": false,
  "name": "file.txt",
  "creationTime": "2025-06-29T12:00:00Z",
  "accessTime": "2025-06-29T12:00:00Z",
  "permissions": "rw-r--r--"
}
```

### move_file
- **Input:** `{ "source": "a.txt", "destination": "b.txt" }`
- **Output:** `{ "ok": true }`

### delete_file
- **Input:** `{ "path": "file.txt" }`
- **Output:** `{ "ok": true }`

### search_files
- **Input:** `{ "path": ".", "pattern": "*.go", "excludePatterns": ["*_test.go"] }`
- **Output:** `{ "matches": ["main.go", "tools/filesystem.go"] }`

### read_multiple_files
- **Input:** `{ "paths": ["a.txt", "b.txt"] }`
- **Output:** `{ "results": { "a.txt": "A", "b.txt": "B" } }`

### list_allowed_directories
- **Input:** `{}`
- **Output:** `{ "directories": ["/your/workdir"] }`

### edit_file
- **Input:** `{ "path": "file.txt", "edits": [ { "oldText": "foo", "newText": "bar" } ], "dryRun": true }`
- **Output:** `{ "error": "not implemented yet" }`

### list_directory_with_sizes
- **Input:** `{ "path": ".", "sortBy": "size" }`
- **Output:**
```json
{
  "entries": [
    { "name": "foo.txt", "type": "file", "size": 123 },
    { "name": "bar", "type": "directory", "size": 0 }
  ],
  "totalFiles": 1,
  "totalDirs": 1,
  "totalSize": 123
}
```

### directory_tree
- **Input:** `{ "path": "." }`
- **Output:**
```json
{
  "tree": {
    "name": "root",
    "type": "directory",
    "children": [
      { "name": "foo.txt", "type": "file" },
      { "name": "bar", "type": "directory", "children": [ ... ] }
    ]
  }
}
```

---

## 🧪 Testing and Automation

### Unit-тесты
```fish
go test -v
make test
```

### Скрипт тестирования
```fish
./test.sh
```

### Проверка разных режимов
```fish
./mcp-filesystem -transport stdio
./mcp-filesystem -transport sse -port 8080
./mcp-filesystem -transport http -port 8080
```

---

## 📁 Структура проекта
```
mcp-filesystem/
├── main.go              # Точка входа
├── main_test.go         # Unit-тесты
├── go.mod               # Go module
├── Makefile             # Сборка и тесты
├── test.sh              # Скрипт тестирования
├── Dockerfile           # Docker
├── tools/               # MCP tools (filesystem.go)
└── README.md            # Документация
```

---

## 📝 Лицензия

MIT License. См. файл [LICENSE](./LICENSE).

---
