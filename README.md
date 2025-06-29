# mcp-filesystem

**Go-реализация MCP server filesystem** — быстрый, безопасный и полностью совместимый сервер для Model Context Protocol (MCP), реализующий все файловые инструменты MCP с поддержкой STDIO, HTTP и SSE транспортов.

---

## 📑 Оглавление
- [Возможности](#возможности)
- [Быстрый старт](#быстрый-старт)
- [Режимы работы и архитектура](#режимы-работы-и-архитектура)
- [Интеграция](#интеграция)
- [API Reference (MCP Tools)](#api-reference-mcp-tools)
- [Тестирование и автоматизация](#тестирование-и-автоматизация)
- [Структура проекта](#структура-проекта)
- [Лицензия](#лицензия)

---

## 🚀 Возможности
- **Полный набор MCP tools**: list_directory, read_file, write_file, create_directory, get_file_info, move_file, delete_file, search_files, read_multiple_files, list_allowed_directories, edit_file (WIP), list_directory_with_sizes, directory_tree
- **Три транспорта**: STDIO (MCP), HTTP (REST), SSE (Server-Sent Events)
- **Многопоточность**: параллельное обслуживание клиентов (goroutines)
- **Ограничение доступа**: работа только в разрешённых директориях
- **Только стандартная библиотека Go**
- **MIT лицензия**
- **Поддержка Linux и macOS**
- **Логирование только в stderr**

---

## ⚡ Быстрый старт

### Установка через go install

Для быстрой установки последней версии из репозитория используйте:

```fish
go install github.com/ad/mcp-filesystem@latest
```

Бинарный файл появится в `$GOBIN` или `$HOME/go/bin` (убедитесь, что этот путь есть в `$PATH`).

### 1. Сборка из исходников
```fish
# Клонируйте репозиторий
git clone https://github.com/ad/mcp-filesystem.git
cd mcp-filesystem

go mod tidy
# Локальная сборка
make build-local
# Или вручную
go build -o mcp-filesystem main.go
# Docker-сборка
make build
```

### 2. Запуск

#### STDIO (VS Code, Claude Desktop)
```fish
./mcp-filesystem -transport stdio
# или
make run-stdio
```

#### HTTP
```fish
./mcp-filesystem -transport http -port 8080
# или
make run-local
```

#### SSE
```fish
./mcp-filesystem -transport sse -port 8080
```

#### Docker
```fish
make run
# или вручную
docker run --rm -p 8080:8080 danielapatin/mcp-filesystem:latest -transport http -port 8080
```

---

## 🌐 Режимы работы и архитектура

- **STDIO** — MCP-совместимость (JSON-RPC через stdin/stdout)
- **HTTP** — REST API (POST /mcp)
- **SSE** — Server-Sent Events (POST /sse)
- **Многопоточность** — каждый клиент обслуживается в отдельной goroutine
- **Ограничение доступа** — операции только в разрешённых директориях
- **Логирование** — только stderr, формат как в оригинале

---

## 🔌 Интеграция

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

### edit_file (WIP)
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

## 🧪 Тестирование и автоматизация

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
