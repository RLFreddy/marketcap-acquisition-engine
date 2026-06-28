package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Scraper.BaseURL != "https://companiesmarketcap.com" {
		t.Errorf("expected default base URL, got %s", cfg.Scraper.BaseURL)
	}
	if cfg.Scraper.Pages != 0 {
		t.Errorf("expected default pages 0, got %d", cfg.Scraper.Pages)
	}
	if cfg.Scraper.CacheTTL != 24*time.Hour {
		t.Errorf("expected default cache TTL 24h, got %v", cfg.Scraper.CacheTTL)
	}
	if cfg.Scraper.Workers != 0 {
		t.Errorf("expected default workers 0, got %d", cfg.Scraper.Workers)
	}
	if cfg.Output.Dir != "." {
		t.Errorf("expected default output dir '.', got %s", cfg.Output.Dir)
	}
	if cfg.Output.FilenamePrefix != "companies_" {
		t.Errorf("expected default filename prefix 'companies_', got %s", cfg.Output.FilenamePrefix)
	}
}

func TestLoad(t *testing.T) {
	content := []byte(`scraper:
  pages: 5
  delay: 1s
  workers: 4
output:
  dir: /tmp
  filename_prefix: test_
`)
	path := t.TempDir() + "/test_config.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Scraper.Pages != 5 {
		t.Errorf("expected pages 5, got %d", cfg.Scraper.Pages)
	}
	if cfg.Scraper.Delay != time.Second {
		t.Errorf("expected delay 1s, got %v", cfg.Scraper.Delay)
	}
	if cfg.Scraper.Workers != 4 {
		t.Errorf("expected workers 4, got %d", cfg.Scraper.Workers)
	}
	if cfg.Output.Dir != "/tmp" {
		t.Errorf("expected output dir /tmp, got %s", cfg.Output.Dir)
	}
	if cfg.Output.FilenamePrefix != "test_" {
		t.Errorf("expected prefix 'test_', got %s", cfg.Output.FilenamePrefix)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestMergeWithFlags(t *testing.T) {
	cfg := Default()
	cfg.Scraper.Pages = 100

	result := MergeWithFlags(cfg, 5)
	if result.Scraper.Pages != 5 {
		t.Errorf("expected pages 5 from flag override, got %d", result.Scraper.Pages)
	}
	if result != cfg {
		t.Error("expected MergeWithFlags to return the same pointer")
	}
}

func TestMergeWithFlagsZero(t *testing.T) {
	cfg := Default()
	cfg.Scraper.Pages = 100

	result := MergeWithFlags(cfg, 0)
	if result.Scraper.Pages != 100 {
		t.Errorf("expected pages 100 unchanged when flag is 0, got %d", result.Scraper.Pages)
	}
}

func TestMergeWithFlagsNil(t *testing.T) {
	result := MergeWithFlags(nil, 5)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Scraper.Pages != 5 {
		t.Errorf("expected pages 5, got %d", result.Scraper.Pages)
	}
}
