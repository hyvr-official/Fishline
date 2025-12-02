#!/bin/bash

echo "Building Go binaries"

mkdir -p build

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o build/fishline.exe

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/fishline

echo "Build complete! Files are in the ./build folder"