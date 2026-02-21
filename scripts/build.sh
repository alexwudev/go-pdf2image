#!/bin/bash
# PDF2Image Quickstart Build Script
# Usage: ./scripts/build.sh [windows|linux]
# 若不帶參數，會顯示互動式選單
cd "$(dirname "$0")/.."

# 載入 nvm
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# 平台選擇
TARGET="$1"
if [ -z "$TARGET" ]; then
    echo ""
    echo "  PDF2Image Build Script"
    echo "  ======================"
    echo "  1) Windows (x64)"
    echo "  2) Linux   (x64)"
    echo ""
    read -p "  選擇目標平台 [1/2]: " choice
    case $choice in
        1) TARGET=windows ;;
        2) TARGET=linux ;;
        *) echo "  無效選擇"; exit 1 ;;
    esac
fi

# 建置前端
echo ""
echo "=== 建置前端 ==="
cd frontend && npm install --silent && npm run build && cd ..
if [ $? -ne 0 ]; then
    echo "前端建置失敗！"
    exit 1
fi

if [ "$TARGET" = "windows" ]; then
    echo ""
    echo "=== 嵌入圖示資源 ==="
    WINRES=""
    if command -v go-winres &> /dev/null; then
        WINRES="go-winres"
    elif [ -f ~/go/bin/go-winres ]; then
        WINRES="${HOME}/go/bin/go-winres"
    fi

    if [ -n "$WINRES" ]; then
        $WINRES make --in platform/windows/winres.json \
            --product-version 1.0.0.0 --file-version 1.0.0.0
    else
        echo "未找到 go-winres，跳過圖示嵌入"
        echo "安裝方式：go install github.com/tc-hib/go-winres@latest"
    fi

    echo ""
    echo "=== 建置 Windows 執行檔 ==="
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
        go build -tags desktop,production -ldflags "-s -w -H windowsgui" -o pdf2image.exe .

    if [ $? -eq 0 ]; then
        rm -f rsrc_windows_*.syso
        echo ""
        echo "建置完成！"
        ls -lh pdf2image.exe
    else
        echo "建置失敗！"
        exit 1
    fi

elif [ "$TARGET" = "linux" ]; then
    echo ""
    echo "=== 建置 Linux 執行檔 ==="
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
        go build -tags desktop,production -ldflags "-s -w" -o pdf2image .

    if [ $? -eq 0 ]; then
        mv pdf2image platform/linux/
        echo ""
        echo "建置完成！"
        ls -lh platform/linux/pdf2image
        echo ""
        echo "輸出目錄：platform/linux/"
        echo "  pdf2image    — 主程式"
        echo "  libmupdf.so  — MuPDF 函式庫（執行時需在 LD_LIBRARY_PATH 中）"
        echo ""
        echo "執行方式：cd platform/linux && LD_LIBRARY_PATH=. ./pdf2image"
    else
        echo "建置失敗！"
        exit 1
    fi

else
    echo "未知的目標平台：$TARGET"
    echo "用法：./build.sh [windows|linux]"
    exit 1
fi
