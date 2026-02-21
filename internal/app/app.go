package app

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gen2brain/go-fitz"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pdf2image/internal/taskbar"
)

type App struct {
	ctx       context.Context
	cancelMu  sync.Mutex
	cancelFn  context.CancelFunc
	cancelled bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	taskbar.Init()
}

// --- File dialogs ---

func (a *App) OpenPDFDialog() (string, error) {
	return wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "選擇 PDF 檔案",
		Filters: []wailsRuntime.FileFilter{{
			DisplayName: "PDF Files (*.pdf)",
			Pattern:     "*.pdf",
		}},
	})
}

func (a *App) SelectOutputDir() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "選擇輸出目錄",
	})
}

// --- PDF Info ---

type PDFInfo struct {
	PageCount int    `json:"page_count"`
	Error     string `json:"error"`
}

func (a *App) GetPDFInfo(path string) PDFInfo {
	doc, err := fitz.New(path)
	if err != nil {
		return PDFInfo{Error: fmt.Sprintf("無法開啟 PDF：%v", err)}
	}
	defer doc.Close()
	return PDFInfo{PageCount: doc.NumPage()}
}

// --- Page Preview ---

func (a *App) GetPagePreview(path string, page int) (string, error) {
	doc, err := fitz.New(path)
	if err != nil {
		return "", fmt.Errorf("open PDF: %w", err)
	}
	defer doc.Close()

	if page < 0 || page >= doc.NumPage() {
		return "", fmt.Errorf("page %d out of range (0-%d)", page, doc.NumPage()-1)
	}

	img, err := doc.ImageDPI(page, 72)
	if err != nil {
		return "", fmt.Errorf("render page %d: %w", page, err)
	}

	var buf strings.Builder
	buf.WriteString("data:image/jpeg;base64,")
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	if err := jpeg.Encode(enc, img, &jpeg.Options{Quality: 75}); err != nil {
		return "", fmt.Errorf("encode preview: %w", err)
	}
	enc.Close()

	return buf.String(), nil
}

// --- Convert ---

type ConvertConfig struct {
	DPI       float64 `json:"dpi"`
	Quality   int     `json:"quality"`
	Format    string  `json:"format"`
	Pages     string  `json:"pages"`
	OutputDir string  `json:"output_dir"`
	Workers   int     `json:"workers"`
	ZipOutput bool    `json:"zip_output"`
}

type ConvertResult struct {
	OutputFiles []string `json:"output_files"`
	Error       string   `json:"error"`
}

func (a *App) ConvertPDF(pdfPath string, cfg ConvertConfig) ConvertResult {
	// Create cancellable context for this conversion
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelMu.Lock()
	a.cancelFn = cancel
	a.cancelled = false
	a.cancelMu.Unlock()
	defer func() {
		a.cancelMu.Lock()
		a.cancelFn = nil
		a.cancelMu.Unlock()
		cancel()
		taskbar.SetProgress(0) // clear taskbar progress
	}()

	// Validate PDF
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return ConvertResult{Error: fmt.Sprintf("無法開啟 PDF：%v", err)}
	}
	total := doc.NumPage()
	doc.Close()

	pages := ParsePages(cfg.Pages, total)
	if len(pages) == 0 {
		return ConvertResult{Error: "沒有有效的頁面可轉換"}
	}

	outDir := cfg.OutputDir
	if outDir == "" {
		outDir = filepath.Dir(pdfPath)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return ConvertResult{Error: fmt.Sprintf("無法建立輸出目錄：%v", err)}
	}

	ext := strings.ToLower(cfg.Format)
	if ext != "jpg" && ext != "png" {
		ext = "jpg"
	}
	dpi := cfg.DPI
	if dpi <= 0 {
		dpi = 300
	}
	quality := cfg.Quality
	if quality <= 0 || quality > 100 {
		quality = 90
	}
	workers := cfg.Workers
	if workers <= 0 {
		workers = 1
	}
	if workers > 20 {
		workers = 20
	}
	if workers > len(pages) {
		workers = len(pages)
	}

	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	totalPages := len(pages)

	// Split pages into chunks for each worker
	chunks := SplitIntoChunks(pages, workers)

	// Get exe path for spawning worker subprocesses
	exePath, err := os.Executable()
	if err != nil {
		return ConvertResult{Error: fmt.Sprintf("無法取得執行檔路徑：%v", err)}
	}

	var done int64
	var mu sync.Mutex
	var firstErr string
	allFiles := make(map[int]string) // pageIdx -> outPath

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}

		// Build comma-separated 0-based page indices
		pageStrs := make([]string, len(chunk))
		for i, p := range chunk {
			pageStrs[i] = strconv.Itoa(p)
		}

		wg.Add(1)
		go func(pageList string) {
			defer wg.Done()

			cmd := exec.CommandContext(ctx, exePath, "--worker",
				"--pdf", pdfPath,
				"--pages", pageList,
				"--dpi", fmt.Sprintf("%.0f", dpi),
				"--quality", fmt.Sprintf("%d", quality),
				"--format", ext,
				"--outdir", outDir,
				"--basename", baseName,
			)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				mu.Lock()
				if firstErr == "" {
					firstErr = fmt.Sprintf("無法建立管道：%v", err)
				}
				mu.Unlock()
				return
			}

			if err := cmd.Start(); err != nil {
				mu.Lock()
				if firstErr == "" {
					firstErr = fmt.Sprintf("無法啟動 worker：%v", err)
				}
				mu.Unlock()
				return
			}

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "OK\t") {
					parts := strings.SplitN(line, "\t", 3)
					if len(parts) == 3 {
						pageIdx, _ := strconv.Atoi(parts[1])
						outPath := parts[2]

						mu.Lock()
						allFiles[pageIdx] = outPath
						mu.Unlock()

						cur := atomic.AddInt64(&done, 1)
						pct := float64(cur) / float64(totalPages) * 100
						taskbar.SetProgress(pct)
						wailsRuntime.EventsEmit(a.ctx, "convert:progress", map[string]interface{}{
							"current": cur,
							"total":   totalPages,
							"page":    pageIdx + 1,
							"percent": pct,
						})
					}
				} else if strings.HasPrefix(line, "ERR\t") {
					msg := strings.TrimPrefix(line, "ERR\t")
					mu.Lock()
					if firstErr == "" {
						firstErr = msg
					}
					mu.Unlock()
				}
			}

			cmd.Wait()
		}(strings.Join(pageStrs, ","))
	}

	wg.Wait()

	// Check if cancelled
	a.cancelMu.Lock()
	wasCancelled := a.cancelled
	a.cancelMu.Unlock()

	if wasCancelled {
		// Clean up partial output files
		for _, path := range allFiles {
			os.Remove(path)
		}
		return ConvertResult{Error: "cancelled"}
	}

	if firstErr != "" {
		return ConvertResult{Error: firstErr}
	}

	// Build output file list in original page order
	outputFiles := make([]string, 0, totalPages)
	for _, pageIdx := range pages {
		if path, ok := allFiles[pageIdx]; ok {
			outputFiles = append(outputFiles, path)
		}
	}

	// ZIP packaging
	if cfg.ZipOutput && len(outputFiles) > 0 {
		wailsRuntime.EventsEmit(a.ctx, "convert:progress", map[string]interface{}{
			"current": totalPages,
			"total":   totalPages,
			"page":    0,
			"percent": 99.0,
			"status":  "zipping",
		})
		zipPath := filepath.Join(outDir, baseName+".zip")
		if err := CreateZip(zipPath, outputFiles); err != nil {
			return ConvertResult{Error: fmt.Sprintf("建立 ZIP 失敗：%v", err)}
		}
		for _, f := range outputFiles {
			os.Remove(f)
		}
		outputFiles = []string{zipPath}
	}

	// Done
	wailsRuntime.EventsEmit(a.ctx, "convert:progress", map[string]interface{}{
		"current": totalPages,
		"total":   totalPages,
		"page":    pages[len(pages)-1] + 1,
		"percent": 100.0,
	})

	return ConvertResult{OutputFiles: outputFiles}
}

// CancelConvert stops the current conversion by killing all worker subprocesses.
func (a *App) CancelConvert() {
	a.cancelMu.Lock()
	defer a.cancelMu.Unlock()
	if a.cancelFn != nil {
		a.cancelled = true
		a.cancelFn()
	}
}

// CreateZip packages the given files into a ZIP archive.
func CreateZip(zipPath string, files []string) error {
	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()

	w := zip.NewWriter(zf)
	defer w.Close()

	for _, fpath := range files {
		fw, err := w.Create(filepath.Base(fpath))
		if err != nil {
			return err
		}
		f, err := os.Open(fpath)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, f)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// SplitIntoChunks divides pages into n roughly equal chunks.
func SplitIntoChunks(pages []int, n int) [][]int {
	chunks := make([][]int, n)
	for i, p := range pages {
		chunks[i%n] = append(chunks[i%n], p)
	}
	return chunks
}

// ParsePages parses a page specification string.
// "all" or "" → all pages; "1,3,5" → specific pages; "2-5" → range.
// Returns 0-based page indices.
func ParsePages(spec string, total int) []int {
	spec = strings.TrimSpace(spec)
	if spec == "" || strings.ToLower(spec) == "all" {
		pages := make([]int, total)
		for i := range pages {
			pages[i] = i
		}
		return pages
	}

	seen := make(map[int]bool)
	var pages []int

	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err1 := strconv.Atoi(strings.TrimSpace(bounds[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err1 != nil || err2 != nil {
				continue
			}
			for p := start; p <= end; p++ {
				idx := p - 1
				if idx >= 0 && idx < total && !seen[idx] {
					seen[idx] = true
					pages = append(pages, idx)
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			idx := p - 1
			if idx >= 0 && idx < total && !seen[idx] {
				seen[idx] = true
				pages = append(pages, idx)
			}
		}
	}
	return pages
}
