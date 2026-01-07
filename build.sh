#!/bin/bash
echo "Building envswitch for all platforms..."

mkdir -p dist

echo ""
echo "[1/4] Windows (amd64)"
GOOS=windows GOARCH=amd64 go build -o dist/envswitch.exe .

echo "[2/4] macOS Intel (amd64)"
GOOS=darwin GOARCH=amd64 go build -o dist/envswitch-mac-intel .

echo "[3/4] macOS Apple Silicon (arm64)"
GOOS=darwin GOARCH=arm64 go build -o dist/envswitch-mac-arm .

echo "[4/4] Linux (amd64)"
GOOS=linux GOARCH=amd64 go build -o dist/envswitch-linux .

echo ""
echo "Done! Binaries are in the dist/ folder:"
ls -la dist/

echo ""
echo "Distribution:"
echo "  Windows users:      dist/envswitch.exe"
echo "  macOS Intel:        dist/envswitch-mac-intel"
echo "  macOS M1/M2/M3:     dist/envswitch-mac-arm"
echo "  Linux:              dist/envswitch-linux"

