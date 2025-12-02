@echo off
echo Building Go binaries

if not exist build mkdir build

set GOOS=windows
set GOARCH=amd64
go build -o build\fishline.exe

set GOOS=linux
set GOARCH=amd64
go build -o build\fishline

echo Build complete! Files are in the /build folder
pause