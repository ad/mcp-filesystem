#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_step() {
  echo -e "${CYAN}==> $1${NC}"
}
print_ok() {
  echo -e "${GREEN}✔ $1${NC}"
}
print_fail() {
  echo -e "${RED}✖ $1${NC}"
}
print_warn() {
  echo -e "${YELLOW}! $1${NC}"
}

print_step "Тестирование MCP Server на Go (все инструменты)..."

# Проверяем наличие собранного сервера
if [ ! -f "./mcp-filesystem" ]; then
    print_warn "Сервер не найден. Запуск сборки..."
    make build-local
fi

print_step "Запуск unit тестов..."
if go test -v; then
    print_ok "Unit тесты пройдены."
else
    print_fail "Unit тесты не пройдены!"
    exit 1
fi

# Создаем временную директорию для тестов
TESTDIR="./test_tmp_dir_$$"
mkdir -p "$TESTDIR"
cd "$TESTDIR"

print_step "Создание тестовых файлов..."
mkdir subdir

echo "hello world" > file1.txt
echo "foo bar" > subdir/file2.txt

print_step "Генерация MCP-запросов для всех инструментов..."
cat > test_input.jsonrpc << EOF
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {"tools": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "list_directory", "arguments": {"path": "."}}}
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "read_file", "arguments": {"path": "file1.txt"}}}
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "write_file", "arguments": {"path": "file3.txt", "content": "test123"}}}
{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "read_file", "arguments": {"path": "file3.txt"}}}
{"jsonrpc": "2.0", "id": 6, "method": "tools/call", "params": {"name": "search_files", "arguments": {"path": ".", "pattern": "file*.txt"}}}
{"jsonrpc": "2.0", "id": 7, "method": "tools/call", "params": {"name": "get_file_info", "arguments": {"path": "file1.txt"}}}
{"jsonrpc": "2.0", "id": 8, "method": "tools/call", "params": {"name": "create_directory", "arguments": {"path": "newdir"}}}
{"jsonrpc": "2.0", "id": 9, "method": "tools/call", "params": {"name": "move_file", "arguments": {"source": "file3.txt", "destination": "newdir/file3moved.txt"}}}
{"jsonrpc": "2.0", "id": 10, "method": "tools/call", "params": {"name": "delete_file", "arguments": {"path": "newdir/file3moved.txt"}}}
# Ошибочный кейс: чтение несуществующего файла
{"jsonrpc": "2.0", "id": 11, "method": "tools/call", "params": {"name": "read_file", "arguments": {"path": "no_such_file.txt"}}}
EOF

print_step "Запуск сервера с тестовыми MCP данными..."

cat test_input.jsonrpc | ../mcp-filesystem -transport stdio "$PWD" | tee test_output.log

print_step "Проверка результатов..."

# Проверяем успешные ответы
check_ok() {
  local id="$1"
  local desc="$2"
  if grep -q '"id":'$id',"result"' test_output.log; then
    print_ok "$desc"
  else
    print_fail "$desc"
  fi
}
# Проверяем ошибку
check_error() {
  local id="$1"
  local desc="$2"
  if grep -q '"id":'$id',"error"' test_output.log; then
    print_ok "$desc (ошибка ожидаема)"
  else
    print_fail "$desc (ожидалась ошибка)"
  fi
}

check_ok 2 "list_directory"
check_ok 3 "read_file (file1.txt)"
check_ok 4 "write_file (file3.txt)"
check_ok 5 "read_file (file3.txt)"
check_ok 6 "search_files"
check_ok 7 "get_file_info"
check_ok 8 "create_directory"
check_ok 9 "move_file"
check_ok 10 "delete_file"
check_error 11 "read_file (no_such_file.txt)"

cd ..
rm -rf "$TESTDIR"

print_step "Тестирование завершено."
