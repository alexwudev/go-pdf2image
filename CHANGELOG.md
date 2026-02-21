# Changelog

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
