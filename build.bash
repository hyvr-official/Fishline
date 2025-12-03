#!/usr/bin/env bash

echo "Building Go binaries"

if [ -d "build" ]; then
    rm -rf build
fi

export GORUN=true

goreleaser release --snapshot --clean

ESC="\033"

echo
echo -e "${ESC}[32mBuild complete! Files are in the /build folder${ESC}[0m"

read -p "Press Enter to continue..."