#/bin/bash

set -e

mkdir -p ./dist/uploads

# Enter the directory where Tania's Golang project is
cd ./backend

# Running any go tests
echo "Running go test.."
go test ./...

# Build the binary
echo "Building golang binaries..."
go build -o ../dist/tania main.go

# Copy all config files and the database file to the dist folder
cd ../
cp ./configuration/conf.json ./dist/conf.json
cp -rf ./database ./dist/database/