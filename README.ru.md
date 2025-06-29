# mcp-filesystem (Русская версия)

> **Примечание:** Это Go-порт официальной [TypeScript-версии](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem) сервера MCP filesystem.

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

(Остальной текст см. оригинальный README.md)
