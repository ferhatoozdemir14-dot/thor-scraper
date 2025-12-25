package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func logStatus(file *os.File, status string, url string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s -> %s\n", timestamp, status, url)
	if _, err := file.WriteString(logLine); err != nil {
		fmt.Printf("[ERR] Rapor dosyasına yazılamadı: %v\n", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go targets.yaml")
		os.Exit(1)
	}

	targetFile := os.Args[1]

	file, err := os.Open(targetFile)
	if err != nil {
		log.Fatalf("Dosya okuma hatası: %v", err)
	}
	defer file.Close()

	// Çıktı klasörü
	outputDir := "scraped_data"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	reportFile, err := os.Create("scan_report.log")
	if err != nil {
		log.Fatalf("Rapor dosyası oluşturulamadı: %v", err)
	}
	defer reportFile.Close()

	scanner := bufio.NewScanner(file)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://127.0.0.1:9150"),
		chromedp.Flag("host-resolver-rules", "MAP * ~NOTFOUND , EXCLUDE 127.0.0.1"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("headless", true),
		chromedp.DisableGPU,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	fmt.Printf("--- TARAMA BAŞLATILIYOR: %s ---\n", targetFile)

	for scanner.Scan() {
		rawLine := scanner.Text()
		url := strings.TrimSpace(rawLine)

		if url == "" || strings.HasPrefix(url, "urls:") || !strings.HasPrefix(url, "http") {
			continue
		}

		fmt.Printf("[INFO] Taranıyor: %s ... -> ", url)

		ctx, cancelCtx := chromedp.NewContext(allocCtx)
		ctx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)

		var imageBuf []byte
		var htmlContent string

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(5*time.Second),
			chromedp.FullScreenshot(&imageBuf, 80),
			chromedp.OuterHTML("html", &htmlContent),
		)

		if err != nil {
			fmt.Printf("[ERR] %v\n", err)
			logStatus(reportFile, "PASSIVE/TIMEOUT", url)
		} else {

			safeName := strings.ReplaceAll(url, "http://", "")
			safeName = strings.ReplaceAll(safeName, ".onion", "")
			safeName = strings.ReplaceAll(safeName, "/", "_")
			if len(safeName) > 50 {
				safeName = safeName[:50]
			}
			timestamp := time.Now().Unix()

			imgFilename := fmt.Sprintf("%s/%s_%d.png", outputDir, safeName, timestamp)
			if err := os.WriteFile(imgFilename, imageBuf, 0644); err != nil {
				fmt.Printf("[ERR] Resim kaydedilemedi\n")
			}

			htmlFilename := fmt.Sprintf("%s/%s_%d.html", outputDir, safeName, timestamp)
			if err := os.WriteFile(htmlFilename, []byte(htmlContent), 0644); err != nil {
				fmt.Printf("[ERR] HTML kaydedilemedi\n")
			} else {
				fmt.Printf("SUCCESS\n")
				logStatus(reportFile, "ACTIVE", url)
			}
		}

		cancelTimeout()
		cancelCtx()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- Tarama Tamamlandı ---")
}
