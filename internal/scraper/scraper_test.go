package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketcap-acquisition-engine/internal/config"

	"github.com/gocolly/colly/v2"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{`"quoted"`, "quoted"},
		{"&amp;", "&"},
		{"&lt;tag&gt;", "<tag>"},
		{"  &quot;spaced&quot;  ", "spaced"},
	}
	for _, tt := range tests {
		result := sanitize(tt.input)
		if result != tt.expected {
			t.Errorf("sanitize(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func testCollector() *colly.Collector {
	return colly.NewCollector(colly.AllowURLRevisit())
}

func testConfig(tsURL string) *config.Config {
	cfg := config.Default()
	cfg.Scraper.BaseURL = tsURL
	cfg.Scraper.Pages = 1
	cfg.Scraper.Delay = 0
	cfg.Scraper.Workers = 1
	return cfg
}

func TestScraper_ParsesCompanies(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/page-2rows.html")
	}))
	defer ts.Close()

	s := &collyScraper{collector: testCollector()}
	companies, err := s.Run(context.Background(), testConfig(ts.URL))
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if len(companies) != 2 {
		t.Fatalf("expected 2 companies, got %d", len(companies))
	}

	c1 := companies[0]
	if c1.Rank != 1 {
		t.Errorf("expected rank 1, got %d", c1.Rank)
	}
	if c1.Name != "NVIDIA (NVDA)" {
		t.Errorf("expected name 'NVIDIA (NVDA)', got %q", c1.Name)
	}
	if c1.MarketCap != "$4.663 T" {
		t.Errorf("expected market cap '$4.663 T', got %q", c1.MarketCap)
	}
	if c1.Price != "$192.53" {
		t.Errorf("expected price '$192.53', got %q", c1.Price)
	}
	if c1.Today != "1.64%" {
		t.Errorf("expected today '1.64%%', got %q", c1.Today)
	}
	if c1.Country != "USA" {
		t.Errorf("expected country 'USA', got %q", c1.Country)
	}

	c2 := companies[1]
	if c2.Rank != 2 {
		t.Errorf("expected rank 2, got %d", c2.Rank)
	}
	if c2.Name != "Apple (AAPL)" {
		t.Errorf("expected name 'Apple (AAPL)', got %q", c2.Name)
	}
	if c2.MarketCap != "$4.167 T" {
		t.Errorf("expected market cap '$4.167 T', got %q", c2.MarketCap)
	}
	if c2.Price != "$283.78" {
		t.Errorf("expected price '$283.78', got %q", c2.Price)
	}
	if c2.Today != "3.14%" {
		t.Errorf("expected today '3.14%%', got %q", c2.Today)
	}
	if c2.Country != "USA" {
		t.Errorf("expected country 'USA', got %q", c2.Country)
	}
}

func TestScraper_DynamicPages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/page/1/" {
			http.ServeFile(w, r, "testdata/page-with-count.html")
		} else {
			http.ServeFile(w, r, "testdata/empty.html")
		}
	}))
	defer ts.Close()

	s := &collyScraper{collector: testCollector()}
	cfg := testConfig(ts.URL)
	cfg.Scraper.Pages = 0
	cfg.Scraper.Workers = 1

	companies, err := s.Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(companies) < 1 {
		t.Fatal("expected at least 1 company from dynamic pages")
	}
}

func TestScraper_ContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		http.ServeFile(w, r, "testdata/page-2rows.html")
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())

	s := &collyScraper{collector: testCollector()}
	cfg := testConfig(ts.URL)
	cfg.Scraper.Delay = 0

	resultCh := make(chan error, 1)
	go func() {
		_, err := s.Run(ctx, cfg)
		resultCh <- err
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-resultCh:
		if err != nil {
			t.Errorf("expected nil error on cancellation, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestScraper_EmptyTable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/empty.html")
	}))
	defer ts.Close()

	s := &collyScraper{collector: testCollector()}
	companies, err := s.Run(context.Background(), testConfig(ts.URL))
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if len(companies) != 0 {
		t.Errorf("expected 0 companies, got %d", len(companies))
	}
}

func TestScraper_MalformedHTML(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/malformed.html")
	}))
	defer ts.Close()

	s := &collyScraper{collector: testCollector()}
	companies, err := s.Run(context.Background(), testConfig(ts.URL))
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if len(companies) != 0 {
		t.Errorf("expected 0 companies, got %d", len(companies))
	}
}

func TestScraper_SkipPagesLessThanTwoRows(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/page-skip.html")
	}))
	defer ts.Close()

	s := &collyScraper{collector: testCollector()}
	companies, err := s.Run(context.Background(), testConfig(ts.URL))
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if len(companies) != 0 {
		t.Errorf("expected 0 companies (row skipped, <7 tds), got %d", len(companies))
	}
}
