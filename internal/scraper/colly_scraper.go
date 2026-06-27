package scraper

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"marketcap-acquisition-engine/internal/domain"
	"marketcap-acquisition-engine/internal/logger"

	"github.com/gocolly/colly/v2"
)

// RunScraper executes the web scraping process and returns a slice of extracted Companies.
func RunScraper(targetPages int) ([]domain.Company, error) {
	var companies []domain.Company
	var mu sync.Mutex

	// Progress counters for pages
	var extractedPages int32
	var totalNumPages int32 = int32(targetPages)

	c := colly.NewCollector(
		colly.AllowedDomains("companiesmarketcap.com"),
		colly.Async(true),
		colly.CacheDir("./colly_cache"), // Cache enabled
		colly.CacheExpiration(24*time.Hour), // Cache expires after 24 hours
	)

	// Calculate the number of workers based on the system's logical cores
	// For I/O-bound tasks (like HTTP requests), NumCPU() * 2 is a standard formula
	numCores := runtime.NumCPU()
	workers := numCores * 2
	logger.Info("💻 Detected %d logical cores. Assigning %d concurrent workers.", numCores, workers)

	// Rate limiting to simulate human behavior and prevent server overload
	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*companiesmarketcap.com*",
		Parallelism: workers,
		// RandomDelay: 1 * time.Second,
	}); err != nil {
		return nil, fmt.Errorf("error setting limit rule: %w", err)
	}

	var pagesOnce sync.Once

	// HTML Handler: Total count of companies (Pagination logic)
	c.OnHTML("span.companies-count", func(e *colly.HTMLElement) {
		if targetPages != 0 {
			return // Specific pages requested, skip dynamic calculation
		}

		pagesOnce.Do(func() {
			text := strings.TrimSpace(e.Text)
			text = strings.ReplaceAll(text, ",", "")
			num, err := strconv.Atoi(text)
			if err == nil && num > 0 {
				totalPages := int(math.Ceil(float64(num) / 100.0))
				atomic.StoreInt32(&totalNumPages, int32(totalPages))
				logger.Info("📊 Detected %d companies. Will dynamically extract %d pages.", num, totalPages)

				for i := 2; i <= totalPages; i++ {
					e.Request.Visit(fmt.Sprintf("https://companiesmarketcap.com/page/%d/", i))
				}
			} else {
				logger.Warn("⚠️ Could not parse total companies count. Only page 1 will be extracted.")
			}
		})
	})

	// HTML Handler: Row Extraction
	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		tds := e.DOM.Find("td")
		if tds.Length() < 7 {
			return // Ignore ad rows or malformed rows
		}

		rankStr := strings.TrimSpace(tds.Eq(1).Text())
		name := strings.TrimSpace(e.DOM.Find("div.company-name").Text())
		code := strings.TrimSpace(e.DOM.Find("div.company-code").Text())
		if code != "" {
			name = fmt.Sprintf("%s (%s)", name, code)
		}

		country := strings.TrimSpace(e.DOM.Find("span.responsive-hidden").Text())
		marketCap := strings.TrimSpace(tds.Eq(3).Text())
		price := strings.TrimSpace(tds.Eq(4).Text())
		today := strings.TrimSpace(tds.Eq(5).Text())

		rank, _ := strconv.Atoi(rankStr)

		company := domain.Company{
			Rank:      rank,
			Name:      name,
			MarketCap: marketCap,
			Price:     price,
			Today:     today,
			Country:   country,
		}

		// Since Async=true, append must be thread-safe
		mu.Lock()
		companies = append(companies, company)
		mu.Unlock()
	})

	// Callback when a page is completely processed
	c.OnScraped(func(r *colly.Response) {
		current := atomic.AddInt32(&extractedPages, 1)
		total := atomic.LoadInt32(&totalNumPages)
		
		if total > 0 {
			logger.Trace("✅ [%d/%d] Successfully extracted page: %s", current, total, r.Request.URL.String())
		} else {
			logger.Trace("✅ [%d/?] Successfully extracted page: %s", current, r.Request.URL.String())
		}
	})

	// Request interceptor: logging and custom headers
	c.OnRequest(func(r *colly.Request) {
		logger.Trace("📡 Requesting: %s", r.URL.String())
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Error("❌ Error on %s: %v", r.Request.URL, err)
	})

	// Entry Point URL
	c.Visit("https://companiesmarketcap.com/page/1/")

	// If explicit pages are requested, enqueue them statically
	if targetPages > 1 {
		for i := 2; i <= targetPages; i++ {
			c.Visit(fmt.Sprintf("https://companiesmarketcap.com/page/%d/", i))
		}
	}

	logger.Info("⏳ Waiting for concurrent requests to finish...")
	c.Wait()

	return companies, nil
}
