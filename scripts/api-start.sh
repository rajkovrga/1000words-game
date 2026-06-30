#!/bin/bash

set -e

APP_NAME="1000words-api"
DB_FILE="words.db"
MIGRATIONS_DIR="database/migrations"
BUILD_DIR="build"

echo "=================================="
echo "Starting $APP_NAME"
echo "=================================="

echo ""
echo "Checking Go installation..."
if ! command -v go >/dev/null 2>&1; then
    echo "ERROR: Go is not installed or not added to PATH."
    exit 1
fi

echo "Go found:"
go version

echo ""
echo "Checking goose CLI..."
if ! command -v goose >/dev/null 2>&1; then
    echo "goose is not installed. Installing goose CLI..."

    go install github.com/pressly/goose/v3/cmd/goose@latest

    export PATH="$PATH:$(go env GOPATH)/bin"

    if ! command -v goose >/dev/null 2>&1; then
        echo "ERROR: goose installed, but it is still not available in PATH."
        echo "Add this folder to PATH:"
        echo "$(go env GOPATH)/bin"
        exit 1
    fi
fi

echo "goose found:"
goose -version

echo ""
echo "Checking migrations folder..."
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "ERROR: Migrations folder does not exist: $MIGRATIONS_DIR"
    exit 1
fi

echo ""
echo "Downloading Go dependencies..."
go mod tidy

echo ""
echo "Running database migrations..."
goose -dir "$MIGRATIONS_DIR" sqlite3 "$DB_FILE" up

echo ""
echo "Migration status:"
goose -dir "$MIGRATIONS_DIR" sqlite3 "$DB_FILE" status

echo ""
echo "Building API application..."
mkdir -p "$BUILD_DIR"

if [[ "$OSTYPE" == "msys"* || "$OSTYPE" == "win32"* || "$OS" == "Windows_NT" ]]; then
    go build -o "$BUILD_DIR/$APP_NAME.exe" ./cmd/api
    APP_PATH="$BUILD_DIR/$APP_NAME.exe"
else
    go build -o "$BUILD_DIR/$APP_NAME" ./cmd/api
    APP_PATH="$BUILD_DIR/$APP_NAME"
fi

echo ""
echo "Build successful:"
echo "$APP_PATH"

echo ""
echo "Starting API application..."
"$APP_PATH"