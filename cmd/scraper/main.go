package main

import (
	"flag"
	"fmt"
	"time"

	"marketcap-acquisition-engine/internal/exporter"
	"marketcap-acquisition-engine/internal/logger"
	"marketcap-acquisition-engine/internal/scraper"
)

func main() {
	// Flags configuration
	pagesFlag := flag.Int("pages", 0, "Number of pages to extract (0 = dynamic, extract all)")
	flag.Parse()

	logger.Info("🚀 Starting CompaniesMarketCap Professional Scraper")

	// 1. Scrape Data
	companies, err := scraper.RunScraper(*pagesFlag)
	if err != nil {
		logger.Fatal("Data extraction failed: %v", err)
	}

	logger.Success("✅ Extraction finished. Total obtained: %d companies.", len(companies))

	// 2. Export Data
	outputFile := fmt.Sprintf("companies_%s.csv", time.Now().Format("2006-01-02"))
	if err := exporter.ExportToCSV(companies, outputFile); err != nil {
		logger.Fatal("Failed exporting to CSV: %v", err)
	}

	logger.Success("🎉 File successfully exported at: %s", outputFile)
}
