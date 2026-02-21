# Changelog

## v1.4.0

### Improvements

- **Project restructure**: platform-specific files (libraries, binaries) organized into `platform/windows/` and `platform/linux/`
- **Interactive build script**: `./build.sh` now shows a menu to choose the target platform; also accepts `./build.sh windows` or `./build.sh linux`

## v1.3.0

### New Features

- **Linux support**: both GUI and CLI modes now work on Linux
  - GUI requires GTK 3 and WebKit2GTK 4.0
  - CLI requires only `libmupdf.so` in the library path
  - `build.sh` supports `./build.sh linux` for native Linux builds
- **Cross-platform window controls**: title bar buttons now use SVG icons instead of Windows-only Segoe MDL2 Assets font

## v1.2.0

### New Features

- **CLI mode**: run conversions from the command line without opening the GUI
  ```
  pdf2image.exe --cli --pdf FILE [--format jpg|png] [--dpi N] [--quality N] [--pages SPEC] [--output DIR] [--workers N] [--zip]
  ```

## v1.1.0

### New Features

- **Title bar progress**: thin progress bar at the top of the window and percentage display in header; window title updates with conversion progress (visible in taskbar)
- **Conversion timer**: live elapsed time display during conversion; final duration shown alongside "Done" upon completion
- **ZIP packaging**: optional checkbox to package all converted images into a single `.zip` file
- **Stop button**: cancel an in-progress conversion; kills all worker subprocesses and cleans up partial output files

## v1.0.0

### Initial Release

- PDF to JPG/PNG conversion with configurable DPI (72–600) and JPEG quality (10–100%)
- Flexible page selection: all pages, specific pages, or ranges
- Parallel multi-process conversion with 1–20 configurable worker subprocesses
- Live page preview with zoom (scroll) and pan (drag)
- Real-time page-by-page conversion progress
- Custom output directory
- Multi-language UI: English, 繁體中文
