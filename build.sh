#!/bin/sh

set -e

mkdir -p ./dist/uploads/areas
mkdir -p ./dist/uploads/crops
mkdir -p ./dist/database/mysql
mkdir -p ./dist/database/sqlite

# Enter the directory where Tania's Golang project is
cd ./backend

# Running any go tests
echo "Running go test.."
go test ./...

# Build the binary
echo "Building golang binaries..."
go build -o ../dist/taniad cmd/taniad/main.go

# Copy all config files and the database file to the dist folder
cp ./conf.json ../dist/conf.json
cp ./database/mysql/ddl.sql ../dist/database/mysql/ddl.sql
cp ./database/sqlite/ddl.sql ../dist/database/sqlite/ddl.sql
