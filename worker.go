package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gen2brain/go-fitz"
)

// runWorker is the headless subprocess entry point.
// Args: --pdf PATH --pages 0,2,5 --dpi 300 --quality 90 --format jpg --outdir DIR --basename NAME
// Output per page: OK\tpageIdx\toutPath\n
// On error:        ERR\tmessage\n
func runWorker(args []string) {
	var pdfPath, pagesStr, format, outDir, baseName string
	var dpi float64
	var quality int

	for i := 0; i < len(args)-1; i += 2 {
		switch args[i] {
		case "--pdf":
			pdfPath = args[i+1]
		case "--pages":
			pagesStr = args[i+1]
		case "--dpi":
			dpi, _ = strconv.ParseFloat(args[i+1], 64)
		case "--quality":
			quality, _ = strconv.Atoi(args[i+1])
		case "--format":
			format = args[i+1]
		case "--outdir":
			outDir = args[i+1]
		case "--basename":
			baseName = args[i+1]
		}
	}

	if pdfPath == "" || pagesStr == "" {
		fmt.Fprintf(os.Stdout, "ERR\tmissing required arguments\n")
		os.Exit(1)
	}

	doc, err := fitz.New(pdfPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "ERR\t%v\n", err)
		os.Exit(1)
	}
	defer doc.Close()

	// Parse page indices (0-based, comma-separated)
	for _, s := range strings.Split(pagesStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		pageIdx, err := strconv.Atoi(s)
		if err != nil {
			continue
		}

		outPath := filepath.Join(outDir, fmt.Sprintf("%s_%d.%s", baseName, pageIdx+1, format))

		img, err := doc.ImageDPI(pageIdx, dpi)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERR\t第 %d 頁渲染失敗：%v\n", pageIdx+1, err)
			os.Exit(1)
		}

		f, err := os.Create(outPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERR\t無法建立檔案 %s：%v\n", outPath, err)
			os.Exit(1)
		}

		switch format {
		case "jpg":
			err = jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
		case "png":
			err = png.Encode(f, img)
		}
		f.Close()

		if err != nil {
			os.Remove(outPath)
			fmt.Fprintf(os.Stdout, "ERR\t第 %d 頁編碼失敗：%v\n", pageIdx+1, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "OK\t%d\t%s\n", pageIdx, outPath)
	}
}
