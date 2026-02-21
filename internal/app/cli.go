package app

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gen2brain/go-fitz"
)

func RunCLI(args []string) {
	var pdfPath, format, pagesSpec, outputDir string
	var dpi float64
	var quality, workers int
	var zipOutput bool

	// Defaults
	format = "jpg"
	dpi = 300
	quality = 90
	workers = 4

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--pdf":
			if i+1 < len(args) {
				pdfPath = args[i+1]
				i++
			}
		case "--format":
			if i+1 < len(args) {
				format = strings.ToLower(args[i+1])
				i++
			}
		case "--dpi":
			if i+1 < len(args) {
				dpi, _ = strconv.ParseFloat(args[i+1], 64)
				i++
			}
		case "--quality":
			if i+1 < len(args) {
				quality, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--pages":
			if i+1 < len(args) {
				pagesSpec = args[i+1]
				i++
			}
		case "--output":
			if i+1 < len(args) {
				outputDir = args[i+1]
				i++
			}
		case "--workers":
			if i+1 < len(args) {
				workers, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--zip":
			zipOutput = true
		}
	}

	if pdfPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --pdf is required")
		fmt.Fprintln(os.Stderr, "Usage: pdf2image.exe --cli --pdf FILE [--format jpg|png] [--dpi N] [--quality N] [--pages SPEC] [--output DIR] [--workers N] [--zip]")
		os.Exit(1)
	}

	// Validate PDF
	doc, err := fitz.New(pdfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot open PDF: %v\n", err)
		os.Exit(1)
	}
	totalInPDF := doc.NumPage()
	doc.Close()

	fmt.Fprintf(os.Stderr, "PDF: %s (%d pages)\n", filepath.Base(pdfPath), totalInPDF)

	// Parse pages
	pages := ParsePages(pagesSpec, totalInPDF)
	if len(pages) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no valid pages to convert")
		os.Exit(1)
	}

	// Output directory
	if outputDir == "" {
		outputDir = filepath.Dir(pdfPath)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot create output directory: %v\n", err)
		os.Exit(1)
	}

	// Clamp parameters
	if format != "jpg" && format != "png" {
		format = "jpg"
	}
	if dpi <= 0 {
		dpi = 300
	}
	if quality <= 0 || quality > 100 {
		quality = 90
	}
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

	fmt.Fprintf(os.Stderr, "Converting %d pages | format=%s dpi=%.0f quality=%d workers=%d\n", totalPages, format, dpi, quality, workers)

	// Split pages into chunks
	chunks := SplitIntoChunks(pages, workers)

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot get executable path: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()

	var done int64
	var mu sync.Mutex
	var firstErr string
	allFiles := make(map[int]string)

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}

		pageStrs := make([]string, len(chunk))
		for i, p := range chunk {
			pageStrs[i] = strconv.Itoa(p)
		}

		wg.Add(1)
		go func(pageList string) {
			defer wg.Done()

			cmd := exec.Command(exePath, "--worker",
				"--pdf", pdfPath,
				"--pages", pageList,
				"--dpi", fmt.Sprintf("%.0f", dpi),
				"--quality", fmt.Sprintf("%d", quality),
				"--format", format,
				"--outdir", outputDir,
				"--basename", baseName,
			)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				mu.Lock()
				if firstErr == "" {
					firstErr = fmt.Sprintf("pipe error: %v", err)
				}
				mu.Unlock()
				return
			}

			if err := cmd.Start(); err != nil {
				mu.Lock()
				if firstErr == "" {
					firstErr = fmt.Sprintf("worker start error: %v", err)
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
						fmt.Fprintf(os.Stderr, "\r[%d/%d] %.0f%% - Page %d done", cur, totalPages, pct, pageIdx+1)
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
	fmt.Fprintln(os.Stderr) // newline after progress

	if firstErr != "" {
		fmt.Fprintf(os.Stderr, "Error: %s\n", firstErr)
		os.Exit(1)
	}

	// Build output file list in page order
	outputFiles := make([]string, 0, totalPages)
	for _, pageIdx := range pages {
		if path, ok := allFiles[pageIdx]; ok {
			outputFiles = append(outputFiles, path)
		}
	}

	// ZIP packaging
	if zipOutput && len(outputFiles) > 0 {
		zipPath := filepath.Join(outputDir, baseName+".zip")
		fmt.Fprintf(os.Stderr, "Creating ZIP: %s\n", zipPath)
		if err := CreateZip(zipPath, outputFiles); err != nil {
			fmt.Fprintf(os.Stderr, "Error: ZIP creation failed: %v\n", err)
			os.Exit(1)
		}
		for _, f := range outputFiles {
			os.Remove(f)
		}
		outputFiles = []string{zipPath}
	}

	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "Done! %d files in %s → %s\n", len(outputFiles), elapsed.Round(time.Millisecond), outputDir)
}
