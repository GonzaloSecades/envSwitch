#!/bin/bash

# Detect OS
case "$(uname -s)" in
    Darwin*)    OS="darwin" ;;
    Linux*)     OS="linux" ;;
    MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
    *)          OS="unknown" ;;
esac

# Detect Architecture
case "$(uname -m)" in
    x86_64|amd64)   ARCH="amd64" ;;
    arm64|aarch64)  ARCH="arm64" ;;
    i386|i686)      ARCH="386" ;;
    *)              ARCH="unknown" ;;
esac

echo "======================================"
echo "  envswitch - Local Build"
echo "======================================"
echo ""
echo "Detected system:"
echo "  OS:           $OS"
echo "  Architecture: $ARCH"
echo ""

if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
    echo "ERROR: Could not detect OS or architecture"
    echo "Please use build.sh to manually select a target"
    exit 1
fi

# Set output filename
if [ "$OS" = "windows" ]; then
    OUTPUT="envswitch.exe"
else
    OUTPUT="envswitch"
fi

echo "Building for $OS/$ARCH..."
echo ""

GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT" .

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "Binary created: $OUTPUT"
    echo "Size: $(ls -lh "$OUTPUT" | awk '{print $5}')"
    echo ""
    echo "Usage:"
    if [ "$OS" = "windows" ]; then
        echo "  ./envswitch.exe --env test"
    else
        echo "  chmod +x $OUTPUT"
        echo "  ./$OUTPUT --env test"
    fi
else
    echo "✗ Build failed!"
    exit 1
fi

