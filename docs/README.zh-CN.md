# PDF2Image

<p align="center">
  <img src="../build/appicon.png" alt="PDF2Image" width="128">
</p>

<p align="center">
  <a href="../README.md">English</a> | <a href="README.zh-TW.md">繁體中文</a> | 简体中文 | <a href="README.ja.md">日本語</a>
</p>

一款桌面应用程序，用于将 PDF 页面转换为高质量图片。使用 [Wails](https://wails.io/)（Go 后端 + Vue 3 前端）开发，通过 [MuPDF](https://mupdf.com/) 进行快速、精确的 PDF 渲染。支持 **Windows** 和 **Linux**（GUI 与 CLI 模式）。

<h2 id="目录">目录</h2>

- [功能](#功能)
- [快速开始](#快速开始)
- [使用方式](#使用方式)
- [命令行模式](#命令行模式)
- [前置需求](#前置需求)
- [从源码构建](#从源码构建)
- [项目结构](#项目结构)
- [许可证](#许可证)

<h2 id="功能">功能 <a href="#目录">⬆</a></h2>

- **输出格式**：JPG 或 PNG
- **可调节 DPI**：72–600（默认 300）
- **JPEG 质量控制**：10–100%（默认 90%）
- **灵活页面选择**：转换全部页面、指定页面或范围（例如 `1-5, 8, 10-12`）
- **并行转换**：可设置 1–20 个 worker 进程同时转换；每个 worker 为独立子进程，各自拥有独立的 MuPDF 实例，完全隔离
- **实时页面预览**：支持缩放（滚轮）和平移（拖拽）；双击重置视图
- **实时转换进度**：逐页显示转换进度；自定义标题栏随进度填色；Windows 任务栏按钮同步显示进度
- **转换计时**：转换期间实时显示已用时间；完成时显示总耗时
- **ZIP 打包**：可选将所有转换图片打包为单一 `.zip` 压缩文件
- **停止按钮**：可随时取消进行中的转换；终止所有 worker 子进程并清理未完成的文件
- **自定义输出目录**：可选择保存位置，或默认与 PDF 同目录
- **多语言界面**：English、繁體中文 — 从右上角下拉菜单切换，偏好设置自动保存

<h2 id="快速开始">快速开始 <a href="#目录">⬆</a></h2>

<h3 id="方式-a下载预编译版本推荐">方式 A：下载预编译版本（推荐） <a href="#目录">⬆</a></h3>

1. 前往 [Releases](https://github.com/alexwudev/go-pdf2image/releases) 页面
2. 下载最新的 `go-pdf2image.zip`
3. 解压到任意文件夹
4. 运行 `pdf2image.exe`

> **注意**：`libmupdf.dll` 必须与 `pdf2image.exe` 放在同一目录下，发行版本中已包含此文件。

<h3 id="方式-b从源码编译">方式 B：从源码编译 <a href="#目录">⬆</a></h3>

请参阅下方[从源码构建](#从源码构建)。

<h2 id="使用方式">使用方式 <a href="#目录">⬆</a></h2>

1. 启动 `pdf2image.exe`
2. 点击**浏览文件**（或拖放文件）打开 PDF
3. 使用导航按钮翻页；滚轮缩放，拖拽平移
4. 在左侧面板调整输出设置：
   - **输出格式**：JPG 或 PNG
   - **DPI**：滑块调整分辨率（72–600）
   - **JPEG 质量**：滑块调整压缩程度（仅 JPG）
   - **同时处理数**：滑块调整并行 worker 进程数量（1–20）
   - **页面范围**：转换全部页面或自定义范围
   - **输出目录**：选择保存目标文件夹
   - **打包为 ZIP**：勾选后输出为单一 `.zip` 压缩文件
5. 点击**开始转换**（点击**停止**可随时取消）
6. 转换完成的图片（或 ZIP）将保存至输出目录

<h2 id="命令行模式">命令行模式 <a href="#目录">⬆</a></h2>

不打开 GUI，直接从命令行执行转换：

```bash
pdf2image.exe --cli --pdf INPUT.pdf [选项]
```

| 选项 | 默认值 | 说明 |
|---|---|---|
| `--pdf PATH` | *（必填）* | 输入 PDF 文件 |
| `--format jpg\|png` | `jpg` | 输出图片格式 |
| `--dpi N` | `300` | 分辨率（72–600） |
| `--quality N` | `90` | JPEG 质量（10–100，仅 JPG） |
| `--pages SPEC` | 全部 | 页面选取（例如 `1-5,8,10-12`） |
| `--output DIR` | 与 PDF 同目录 | 输出目录 |
| `--workers N` | `4` | 并行 worker 进程数（1–20） |
| `--zip` | 关闭 | 将输出打包为单一 `.zip` 文件 |

**示例：**

```bash
pdf2image.exe --cli --pdf report.pdf --format png --dpi 150 --pages 1-10 --workers 8 --output ./images
```

进度输出至 stderr：

```
PDF: report.pdf (50 pages)
Converting 10 pages | format=png dpi=150 quality=90 workers=8
[10/10] 100% - Page 10 done
Done! 10 files in 5.2s → ./images
```

<h2 id="前置需求">前置需求 <a href="#目录">⬆</a></h2>

**Windows：**

- **Windows 10/11**（x64）
- **[Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)**（大多数 Windows 10/11 系统已预装）
- **`libmupdf.dll`** 必须与可执行文件放在同一目录（发行版本中已包含）

**Linux**（x64）：

- **GTK 3** 和 **WebKit2GTK 4.0**（GUI 模式所需）
  ```bash
  # Ubuntu/Debian
  sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37
  ```
- **`libmupdf.so`** 必须在动态链接器可搜索的路径中（发行版本已包含）
  ```bash
  # 以同目录的库运行
  LD_LIBRARY_PATH=. ./pdf2image
  ```

<h2 id="从源码构建">从源码构建 <a href="#目录">⬆</a></h2>

<h3 id="需求">需求 <a href="#目录">⬆</a></h3>

- [Go](https://go.dev/) 1.24+
- [Node.js](https://nodejs.org/)
- [go-winres](https://github.com/tc-hib/go-winres)（仅 Windows 构建，嵌入应用图标用）：`go install github.com/tc-hib/go-winres@latest`
- **Linux 构建**另需：`sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev`

<h3 id="wsl交叉编译为-windows">WSL（交叉编译为 Windows） <a href="#目录">⬆</a></h3>

```bash
./build.sh            # 或：./build.sh windows
```

<h3 id="linux原生编译">Linux（原生编译） <a href="#目录">⬆</a></h3>

```bash
./build.sh linux
```

<h3 id="windows原生编译">Windows（原生编译） <a href="#目录">⬆</a></h3>

```batch
build.bat
```

<h3 id="开发模式">开发模式 <a href="#目录">⬆</a></h3>

需要 [Wails CLI](https://wails.io/docs/gettingstarted/installation)。

```bash
wails dev
```

<h3 id="libmupdf">libmupdf.dll / libmupdf.so <a href="#目录">⬆</a></h3>

可执行文件需要 MuPDF 共享库（1.24.9, x64）在同一目录或库路径中。

**Windows**（`libmupdf.dll`）— 从 WSL 交叉编译：

```bash
# 需要 mingw-w64：sudo apt install gcc-mingw-w64-x86-64
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make OS=mingw64-cross shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 复制到项目：cp build/shared-release/libmupdf.dll /path/to/go-pdf2image/platform/windows/
```

**Linux**（`libmupdf.so`）— 原生编译：

```bash
git clone --recursive --branch 1.24.9 --depth 1 https://github.com/ArtifexSoftware/mupdf.git
cd mupdf
make shared=yes build=release \
  HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no \
  -j$(nproc)
# 复制到项目：cp build/shared-release/libmupdf.so.24.9 /path/to/go-pdf2image/platform/linux/libmupdf.so
```

<h2 id="项目结构">项目结构 <a href="#目录">⬆</a></h2>

```
go-pdf2image/
├── main.go              # 入口：无边框 GUI、--cli 或 --worker 子进程模式
├── app.go               # Go 后端：PDF 信息、预览、多进程转换
├── cli.go               # 命令行模式：不打开 GUI 直接转换
├── worker.go            # 无界面 worker 子进程：渲染与编码页面
├── taskbar_windows.go   # Windows 任务栏进度（ITaskbarList3）与图标
├── taskbar_stub.go      # 非 Windows 平台的空实现
├── go.mod / go.sum      # Go 依赖
├── wails.json           # Wails 项目配置
├── build.sh             # Quickstart 构建脚本（交互菜单或参数）
├── build.bat            # Windows 原生编译脚本
├── platform/
│   ├── windows/
│   │   ├── libmupdf.dll     # MuPDF 动态链接库（运行时依赖）
│   │   ├── winres.json      # go-winres 配置（图标与 manifest 嵌入）
│   │   └── pdf2image.exe    # 构建产出
│   └── linux/
│       ├── libmupdf.so      # MuPDF 共享库（运行时依赖）
│       └── pdf2image         # 构建产出
├── build/
│   ├── appicon.png      # 应用图标
│   └── windows/         # Windows manifest 与图标资源
├── docs/                # 翻译版 README
├── CHANGELOG.md
├── LICENSE
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    └── src/
        ├── main.ts          # Vue 应用初始化
        ├── App.vue          # 根组件布局、自定义标題栏 + 语言切换
        ├── style.css        # 全局样式
        ├── i18n/            # 国际化（en、zh-TW）
        ├── stores/
        │   └── appStore.ts  # Pinia 状态管理
        └── components/
            ├── PdfImport.vue       # PDF 文件选择器
            ├── SettingsPanel.vue   # 输出设置
            ├── ActionBar.vue       # 转换按钮与状态消息
            ├── PreviewPanel.vue    # 页面预览（缩放/平移）
            └── ConvertProgress.vue # 转换进度条
```

<h2 id="许可证">许可证 <a href="#目录">⬆</a></h2>

[MIT](../LICENSE)
