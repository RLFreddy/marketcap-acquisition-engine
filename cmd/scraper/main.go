package main

import (
	"context"
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
	cfg, err := config.Load("config.yaml", "/etc/scraper/config.yaml")
	if err != nil {
		logger.Info("No config file found, using built-in defaults")
		cfg = config.Default()
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Info("MarketCap Acquisition Engine")
	logger.Info("Target: %d pages (0 = dynamic)", cfg.Scraper.Pages)

	companies, err := scraper.New().Run(ctx, cfg)
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
