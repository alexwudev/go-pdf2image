# PDF2Image

<p align="center">
  <img src="../build/appicon.png" alt="PDF2Image" width="128">
</p>

<p align="center">
  <a href="../README.md">English</a> | 繁體中文 | <a href="README.zh-CN.md">简体中文</a> | <a href="README.ja.md">日本語</a>
</p>

一款桌面應用程式，用於將 PDF 頁面轉換為高品質圖片。使用 [Wails](https://wails.io/)（Go 後端 + Vue 3 前端）開發，透過 [MuPDF](https://mupdf.com/) 進行快速、精確的 PDF 渲染。支援 **Windows** 及 **Linux**（GUI 與 CLI 模式）。

<h2 id="目錄">目錄</h2>

- [功能](#功能)
- [快速開始](#快速開始)
- [使用方式](#使用方式)
- [命令列模式](#命令列模式)
- [前置需求](#前置需求)
- [從原始碼建置](#從原始碼建置)
- [專案結構](#專案結構)
- [授權條款](#授權條款)

<h2 id="功能">功能 <a href="#目錄">⬆</a></h2>

- **輸出格式**：JPG 或 PNG
- **可調整 DPI**：72–600（預設 300）
- **JPEG 品質控制**：10–100%（預設 90%）
- **彈性頁面選取**：轉換全部頁面、指定頁面或範圍（例如 `1-5, 8, 10-12`）
- **並行轉換**：可設定 1–20 個 worker 進程同時轉換；每個 worker 為獨立子進程，各自擁有獨立的 MuPDF 實例，完全隔離
- **即時頁面預覽**：支援縮放（滾輪）和平移（拖曳）；雙擊重置視圖
- **即時轉換進度**：逐頁顯示轉換進度；自訂標題列隨進度填色；Windows 工作列按鈕同步顯示進度
- **轉換計時**：轉換期間即時顯示已用時間；完成時顯示總耗時
- **ZIP 打包**：可選將所有轉換圖片打包為單一 `.zip` 壓縮檔
- **停止按鈕**：可隨時取消進行中的轉換；終止所有 worker 子進程並清理未完成的檔案
- **自訂輸出目錄**：可選擇儲存位置，或預設與 PDF 同目錄
- **多語言介面**：English、繁體中文 — 從右上角下拉選單切換，偏好設定自動儲存

<h2 id="快速開始">快速開始 <a href="#目錄">⬆</a></h2>

<h3 id="方式-a下載預編譯版本推薦">方式 A：下載預編譯版本（推薦） <a href="#目錄">⬆</a></h3>

1. 前往 [Releases](https://github.com/alexwudev/go-pdf2image/releases) 頁面
2. 下載最新的 `go-pdf2image.zip`
3. 解壓縮到任意資料夾
4. 執行 `pdf2image.exe`

> **注意**：`libmupdf.dll` 必須與 `pdf2image.exe` 放在同一目錄下，發行版本中已包含此檔案。

<h3 id="方式-b從原始碼編譯">方式 B：從原始碼編譯 <a href="#目錄">⬆</a></h3>

請參閱下方[從原始碼建置](#從原始碼建置)。

<h2 id="使用方式">使用方式 <a href="#目錄">⬆</a></h2>

1. 啟動 `pdf2image.exe`
2. 點擊**瀏覽檔案**（或拖放檔案）開啟 PDF
3. 使用導覽按鈕翻頁；滾輪縮放，拖曳平移
4. 在左側面板調整輸出設定：
   - **輸出格式**：JPG 或 PNG
   - **DPI**：滑桿調整解析度（72–600）
   - **JPEG 品質**：滑桿調整壓縮程度（僅 JPG）
   - **同時處理數**：滑桿調整並行 worker 進程數量（1–20）
   - **頁面範圍**：轉換全部頁面或自訂範圍
   - **輸出目錄**：選擇儲存目的資料夾
   - **打包為 ZIP**：勾選後輸出為單一 `.zip` 壓縮檔
5. 點擊**開始轉換**（點擊**停止**可隨時取消）
6. 轉換完成的圖片（或 ZIP）將儲存至輸出目錄

<h2 id="命令列模式">命令列模式 <a href="#目錄">⬆</a></h2>

不開啟 GUI，直接從命令列執行轉換：

```bash
pdf2image.exe --cli --pdf INPUT.pdf [選項]
```

| 選項 | 預設值 | 說明 |
|---|---|---|
| `--pdf PATH` | *（必填）* | 輸入 PDF 檔案 |
| `--format jpg\|png` | `jpg` | 輸出圖片格式 |
| `--dpi N` | `300` | 解析度（72–600） |
| `--quality N` | `90` | JPEG 品質（10–100，僅 JPG） |
| `--pages SPEC` | 全部 | 頁面選取（例如 `1-5,8,10-12`） |
| `--output DIR` | 與 PDF 同目錄 | 輸出目錄 |
| `--workers N` | `4` | 並行 worker 進程數（1–20） |
| `--zip` | 關閉 | 將輸出打包為單一 `.zip` 檔案 |

**範例：**

```bash
pdf2image.exe --cli --pdf report.pdf --format png --dpi 150 --pages 1-10 --workers 8 --output ./images
```

進度輸出至 stderr：

```
PDF: report.pdf (50 pages)
Converting 10 pages | format=png dpi=150 quality=90 workers=8
[10/10] 100% - Page 10 done
Done! 10 files in 5.2s → ./images
```

<h2 id="前置需求">前置需求 <a href="#目錄">⬆</a></h2>

**Windows：**

- **Windows 10/11**（x64）
- **[Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)**（大多數 Windows 10/11 系統已預裝）
- **`libmupdf.dll`** 必須與執行檔放在同一目錄（發行版本中已包含）

**Linux**（x64）：

- **GTK 3** 和 **WebKit2GTK 4.0**（GUI 模式所需）
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
  ```
- **`libmupdf.so`** 必須在動態連結器可搜尋的路徑中（發行版本已包含）
  ```bash
  # 以同目錄的程式庫執行
  LD_LIBRARY_PATH=. ./pdf2image
  ```

<h2 id="從原始碼建置">從原始碼建置 <a href="#目錄">⬆</a></h2>

<h3 id="需求">需求 <a href="#目錄">⬆</a></h3>

**共通（兩個平台都需要）：**

- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/)（建置前端用）

**Windows 建置（WSL 交叉編譯）：**

```bash
# go-winres 用於嵌入應用程式圖示
go install github.com/tc-hib/go-winres@latest
```

**Linux 建置（原生編譯）：**

```bash
# Ubuntu/Debian
sudo apt install gcc pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
```

<h3 id="wsl交叉編譯為-windows">WSL（交叉編譯為 Windows） <a href="#目錄">⬆</a></h3>

```bash
./scripts/build.sh            # 或：./scripts/build.sh windows
# 產出：pdf2image.exe (project root)
```

<h3 id="linux原生編譯">Linux（原生編譯） <a href="#目錄">⬆</a></h3>

```bash
./scripts/build.sh linux
# 產出：platform/linux/pdf2image
```

<h3 id="windows原生編譯">Windows（原生編譯） <a href="#目錄">⬆</a></h3>

```batch
scripts\build.bat
REM 產出：pdf2image.exe (project root)
```

<h3 id="開發模式">開發模式 <a href="#目錄">⬆</a></h3>

需要 [Wails CLI](https://wails.io/docs/gettingstarted/installation)。

```bash
wails dev
```

<h3 id="libmupdf">libmupdf.dll / libmupdf.so <a href="#目錄">⬆</a></h3>

執行檔需要 MuPDF 共用程式庫（1.24.9, x64）放在同一目錄或程式庫路徑中。

**Windows**（`libmupdf.dll`）— 從 WSL 交叉編譯：

```bash
# 需要 mingw-w64：sudo apt install gcc-mingw-w64-x86-64
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make OS=mingw64-cross shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 複製到專案：cp build/shared-release/libmupdf.dll /path/to/go-pdf2image/platform/windows/
```

**Linux**（`libmupdf.so`）— 原生編譯：

```bash
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 複製到專案：cp build/shared-release/libmupdf.so.24.9 /path/to/go-pdf2image/platform/linux/libmupdf.so
```

<h2 id="專案結構">專案結構 <a href="#目錄">⬆</a></h2>

```
go-pdf2image/
├── main.go              # Entry point: frameless GUI, --cli, or --worker subprocess mode
├── go.mod / go.sum      # Go dependencies
├── wails.json           # Wails project config
├── LICENSE
├── README.md
├── internal/
│   ├── app/
│   │   ├── app.go       # App struct, PDF info, preview, multi-process conversion
│   │   ├── cli.go       # CLI mode: command-line conversion without GUI
│   │   └── worker.go    # Headless worker subprocess: render & encode pages
│   └── taskbar/
│       ├── taskbar_windows.go  # Windows taskbar progress (ITaskbarList3) & icon
│       └── taskbar_stub.go     # No-op stub for non-Windows builds
├── scripts/
│   ├── build.sh         # Quickstart build script (interactive menu or argument)
│   └── build.bat        # Windows native build script
├── platform/
│   ├── windows/
│   │   └── winres.json      # go-winres config (icon & manifest)
│   └── linux/
├── build/
│   ├── appicon.png      # App icon
│   └── windows/         # Windows manifest & icon resources
├── docs/                # Translated READMEs, CHANGELOG
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    └── src/
        ├── main.ts          # Vue app init
        ├── App.vue          # Root layout, custom title bar + language switcher
        ├── style.css        # Global styles
        ├── i18n/            # Internationalization (en, zh-TW)
        ├── stores/
        │   └── appStore.ts  # Pinia state management
        └── components/
            ├── PdfImport.vue       # PDF file picker
            ├── SettingsPanel.vue   # Output settings
            ├── ActionBar.vue       # Convert button & status messages
            ├── PreviewPanel.vue    # Page preview with zoom/pan
            └── ConvertProgress.vue # Conversion progress bar
```

<h2 id="授權條款">授權條款 <a href="#目錄">⬆</a></h2>

[MIT](../LICENSE)
