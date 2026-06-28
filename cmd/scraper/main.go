package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketcap-acquisition-engine/internal/config"
	"marketcap-acquisition-engine/internal/exporter"
	"marketcap-acquisition-engine/internal/logger"
	"marketcap-acquisition-engine/internal/scraper"
)

func main() {
	pagesFlag := flag.Int("pages", 0, "Number of pages to extract (0 = dynamic, extract all)")
	configFlag := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configFlag)
	if err != nil {
		logger.Warn("Could not load config: %v, using defaults", err)
		cfg = config.Default()
	}
	cfg = config.MergeWithFlags(cfg, *pagesFlag)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Info("MarketCap Acquisition Engine")
	logger.Info("Target: %d pages (0 = dynamic)", cfg.Scraper.Pages)

	companies, err := scraper.RunScraper(ctx, cfg)
	if err != nil {
		logger.Fatal("Data extraction failed: %v", err)
	}

	logger.Success("Extraction finished. Total obtained: %d companies.", len(companies))

	outputFile := fmt.Sprintf("%s/%s%s.csv", cfg.Output.Dir, cfg.Output.FilenamePrefix, time.Now().Format("2006-01-02"))
	if err := exporter.ExportToCSV(companies, outputFile); err != nil {
		logger.Fatal("Failed exporting to CSV: %v", err)
	}

	logger.Success("File successfully exported at: %s", outputFile)
}
