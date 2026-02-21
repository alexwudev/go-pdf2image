#!/bin/bash
# Build script — supports Windows (cross-compile from WSL) and Linux (native)
# Usage: ./build.sh [windows|linux]   (default: windows)
cd "$(dirname "$0")"

TARGET="${1:-windows}"

# Load nvm for Node.js
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

echo "=== Building frontend ==="
cd frontend && npm install && npm run build && cd ..

if [ "$TARGET" = "linux" ]; then
    echo "=== Building Linux binary ==="
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
        go build -tags desktop,production -ldflags "-s -w" -o pdf2image .

    if [ $? -eq 0 ]; then
        echo "Build OK! -> pdf2image"
        ls -lh pdf2image
    else
        echo "Build FAILED!"
        exit 1
    fi
else
    echo "=== Embedding icon resource ==="
    if command -v go-winres &> /dev/null || [ -f ~/go/bin/go-winres ]; then
        WINRES="${HOME}/go/bin/go-winres"
        command -v go-winres &> /dev/null && WINRES="go-winres"
        $WINRES make --in winres.json --product-version 1.0.0.0 --file-version 1.0.0.0
    else
        echo "go-winres not found, skipping icon embed (install: go install github.com/tc-hib/go-winres@latest)"
    fi

    echo "=== Building Windows exe ==="
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
        go build -tags desktop,production -ldflags "-s -w -H windowsgui" -o pdf2image.exe .

    if [ $? -eq 0 ]; then
        echo "Build OK! -> pdf2image.exe"
        ls -lh pdf2image.exe
    else
        echo "Build FAILED!"
        exit 1
    fi
fi
