# PDF2Image

<p align="center">
  <img src="build/appicon.png" alt="PDF2Image" width="128">
</p>

<p align="center">
  English | <a href="docs/README.zh-TW.md">繁體中文</a> | <a href="docs/README.zh-CN.md">简体中文</a> | <a href="docs/README.ja.md">日本語</a>
</p>

A desktop application for converting PDF pages to high-quality images, built with [Wails](https://wails.io/) (Go backend + Vue 3 frontend). It uses [MuPDF](https://mupdf.com/) for fast, accurate PDF rendering. Supports **Windows** and **Linux** (GUI and CLI modes).

<h2 id="table-of-contents">Table of Contents</h2>

- [Features](#features)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [CLI Mode](#cli-mode)
- [Prerequisites](#prerequisites)
- [Building from Source](#building-from-source)
- [Project Structure](#project-structure)
- [License](#license)

<h2 id="features">Features <a href="#table-of-contents">⬆</a></h2>

- **Output formats**: JPG or PNG
- **Configurable DPI**: 72–600 (default 300)
- **JPEG quality control**: 10–100% (default 90%)
- **Flexible page selection**: convert all pages, specific pages, or ranges (e.g. `1-5, 8, 10-12`)
- **Parallel conversion**: configurable 1–20 worker processes for batch conversion; each worker is an independent subprocess with its own MuPDF instance for full isolation
- **Live page preview**: zoom (scroll wheel) and pan (drag) support; double-click to reset view
- **Real-time progress**: page-by-page progress display; custom title bar fills with color as conversion progresses; Windows taskbar button also shows progress
- **Conversion timer**: live elapsed time during conversion; final duration shown on completion
- **ZIP packaging**: optionally package all converted images into a single `.zip` file
- **Stop button**: cancel an in-progress conversion; worker subprocesses are terminated and partial files cleaned up
- **Custom output directory**: choose where to save, or default to the same directory as the PDF
- **Multi-language UI**: English, 繁體中文 — switchable from the dropdown in the top-right corner, preference saved automatically

<h2 id="quick-start">Quick Start <a href="#table-of-contents">⬆</a></h2>

<h3 id="option-a-download-pre-built-release-recommended">Option A: Download Pre-built Release (Recommended) <a href="#table-of-contents">⬆</a></h3>

1. Go to the [Releases](https://github.com/alexwudev/go-pdf2image/releases) page
2. Download the latest `go-pdf2image.zip`
3. Extract to any folder
4. Run `pdf2image.exe`

> **Note**: `libmupdf.dll` must be in the same directory as `pdf2image.exe`. It is included in the release package.

<h3 id="option-b-build-from-source">Option B: Build from Source <a href="#table-of-contents">⬆</a></h3>

See [Building from Source](#building-from-source) below.

<h2 id="usage">Usage <a href="#table-of-contents">⬆</a></h2>

1. Launch `pdf2image.exe`
2. Click **Browse Files** (or drag & drop) to open a PDF
3. Preview pages using the navigation buttons; zoom with scroll wheel, pan by dragging
4. Adjust output settings in the left panel:
   - **Output Format**: JPG or PNG
   - **DPI**: slide to set resolution (72–600)
   - **JPEG Quality**: slide to set compression level (JPG only)
   - **Concurrency**: slide to set number of parallel worker processes (1–20)
   - **Page Range**: convert all pages or specify a custom range
   - **Output Directory**: choose a destination folder
   - **Package as ZIP**: check to output a single `.zip` file instead of individual images
5. Click **Convert** (click **Stop** to cancel mid-conversion)
6. Converted images (or ZIP) are saved to the output directory

<h2 id="cli-mode">CLI Mode <a href="#table-of-contents">⬆</a></h2>

Run conversions from the command line without opening the GUI:

```bash
pdf2image.exe --cli --pdf INPUT.pdf [options]
```

| Option | Default | Description |
|---|---|---|
| `--pdf PATH` | *(required)* | Input PDF file |
| `--format jpg\|png` | `jpg` | Output image format |
| `--dpi N` | `300` | Resolution (72–600) |
| `--quality N` | `90` | JPEG quality (10–100, JPG only) |
| `--pages SPEC` | all | Page selection (e.g. `1-5,8,10-12`) |
| `--output DIR` | same as PDF | Output directory |
| `--workers N` | `4` | Number of parallel worker processes (1–20) |
| `--zip` | off | Package output into a single `.zip` file |

**Example:**

```bash
pdf2image.exe --cli --pdf report.pdf --format png --dpi 150 --pages 1-10 --workers 8 --output ./images
```

Progress is printed to stderr:

```
PDF: report.pdf (50 pages)
Converting 10 pages | format=png dpi=150 quality=90 workers=8
[10/10] 100% - Page 10 done
Done! 10 files in 5.2s → ./images
```

<h2 id="prerequisites">Prerequisites <a href="#table-of-contents">⬆</a></h2>

**Windows:**

- **Windows 10/11** (x64)
- **[Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)** (pre-installed on most Windows 10/11 systems)
- **`libmupdf.dll`** must be in the same directory as the executable (included in releases)

**Linux** (x64):

- **GTK 3** and **WebKit2GTK 4.0** (for GUI mode)
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
  ```
- **`libmupdf.so`** must be findable by the dynamic linker (included in releases)
  ```bash
  # Run with library in the same directory
  LD_LIBRARY_PATH=. ./pdf2image
  ```

<h2 id="building-from-source">Building from Source <a href="#table-of-contents">⬆</a></h2>

<h3 id="requirements">Requirements <a href="#table-of-contents">⬆</a></h3>

**Common (both platforms):**

- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/) (for building the frontend)

**Windows build (WSL cross-compile):**

```bash
# go-winres for embedding the app icon
go install github.com/tc-hib/go-winres@latest
```

**Linux build (native):**

```bash
# Ubuntu/Debian
sudo apt install gcc pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
```

<h3 id="wsl-cross-compile-to-windows">WSL (cross-compile to Windows) <a href="#table-of-contents">⬆</a></h3>

```bash
./scripts/build.sh            # or: ./scripts/build.sh windows
# Output: pdf2image.exe (project root)
```

<h3 id="linux-native">Linux (native) <a href="#table-of-contents">⬆</a></h3>

```bash
./scripts/build.sh linux
# Output: platform/linux/pdf2image
```

<h3 id="windows-native">Windows (native) <a href="#table-of-contents">⬆</a></h3>

```batch
scripts\build.bat
REM Output: pdf2image.exe (project root)
```

<h3 id="development-mode">Development Mode <a href="#table-of-contents">⬆</a></h3>

Requires [Wails CLI](https://wails.io/docs/gettingstarted/installation).

```bash
wails dev
```

<h3 id="libmupdf">libmupdf.dll / libmupdf.so <a href="#table-of-contents">⬆</a></h3>

The executable requires the MuPDF shared library (1.24.9, x64) in the same directory or library path.

**Windows** (`libmupdf.dll`) — cross-compile from WSL:

```bash
# Requires mingw-w64: sudo apt install gcc-mingw-w64-x86-64
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make OS=mingw64-cross shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# Copy to project: cp build/shared-release/libmupdf.dll /path/to/go-pdf2image/platform/windows/
```

**Linux** (`libmupdf.so`) — native build:

```bash
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# Copy to project: cp build/shared-release/libmupdf.so.24.9 /path/to/go-pdf2image/platform/linux/libmupdf.so
```

<h2 id="project-structure">Project Structure <a href="#table-of-contents">⬆</a></h2>

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

<h2 id="license">License <a href="#table-of-contents">⬆</a></h2>

[MIT](LICENSE)
