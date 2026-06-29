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

func TestLoadMultiplePaths(t *testing.T) {
	content := []byte(`scraper:
  pages: 10
`)
	first := t.TempDir() + "/first.yaml"
	if err := os.WriteFile(first, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(first, "/nonexistent/path.yaml")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Scraper.Pages != 10 {
		t.Errorf("expected pages 10, got %d", cfg.Scraper.Pages)
	}
}

func TestLoadFallbackToNextPath(t *testing.T) {
	content := []byte(`scraper:
  pages: 7
`)
	second := t.TempDir() + "/second.yaml"
	if err := os.WriteFile(second, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load("/nonexistent/path.yaml", second)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Scraper.Pages != 7 {
		t.Errorf("expected pages 7 from fallback, got %d", cfg.Scraper.Pages)
	}
}
