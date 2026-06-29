package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Scraper.Pages != 0 {
		t.Errorf("expected default pages 0 (auto-detect), got %d", cfg.Scraper.Pages)
	}
	if cfg.Scraper.CacheTTL != 24*time.Hour {
		t.Errorf("expected default cache TTL 24h, got %v", cfg.Scraper.CacheTTL)
	}
	if cfg.Scraper.Workers != 0 {
		t.Errorf("expected default workers 0 (auto), got %d", cfg.Scraper.Workers)
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

func TestValidateRejectsNegativePages(t *testing.T) {
	content := []byte(`scraper:
  pages: -1
`)
	path := t.TempDir() + "/negative_pages.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for pages=-1")
	}
}

func TestValidateRejectsNegativeWorkers(t *testing.T) {
	content := []byte(`scraper:
  workers: -5
`)
	path := t.TempDir() + "/negative_workers.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for workers=-5")
	}
}

func TestValidateRejectsNegativeDelay(t *testing.T) {
	content := []byte(`scraper:
  delay: -1s
`)
	path := t.TempDir() + "/negative_delay.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for delay=-1s")
	}
}

func TestValidateRejectsNegativeRetryCount(t *testing.T) {
	content := []byte(`scraper:
  retry_count: -1
`)
	path := t.TempDir() + "/negative_retry.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for retry_count=-1")
	}
}

func TestValidateAcceptsZeroValues(t *testing.T) {
	content := []byte(`scraper:
  pages: 0
  delay: 0s
  workers: 0
  retry_count: 0
  retry_delay: 0s
  cache_ttl: 0s
`)
	path := t.TempDir() + "/zero_values.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err != nil {
		t.Fatalf("expected zero values to be valid, got: %v", err)
	}
}

func TestValidateAcceptsPositiveValues(t *testing.T) {
	content := []byte(`scraper:
  pages: 10
  delay: 2s
  workers: 8
  retry_count: 5
`)
	path := t.TempDir() + "/positive_values.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err != nil {
		t.Fatalf("expected positive values to be valid, got: %v", err)
	}
}

func TestLoadDetectsUnknownFields(t *testing.T) {
	content := []byte(`scraper:
  pages: 5
  unknown_field: true
`)
	path := t.TempDir() + "/unknown_field.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestLoadDetectsBadYAML(t *testing.T) {
	content := []byte(`scraper:
  pages: not_a_number
`)
	path := t.TempDir() + "/bad_yaml.yaml"
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for bad YAML, got nil")
	}
}
