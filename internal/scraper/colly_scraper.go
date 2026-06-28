package scraper

import (
	"fmt"
	"html"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"marketcap-acquisition-engine/internal/config"
	"marketcap-acquisition-engine/internal/domain"
	"marketcap-acquisition-engine/internal/logger"

	"github.com/gocolly/colly/v2"
)

const (
	companiesPerPage = 100
	defaultCacheTTL  = 24 * time.Hour
)

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"")
	s = html.UnescapeString(s)
	return s
}

func RunScraper(cfg *config.Config) ([]domain.Company, error) {
	var companies []domain.Company
	var mu sync.Mutex

	var extractedPages int32
	targetPages := cfg.Scraper.Pages
	var totalNumPages int32 = int32(targetPages)

	cacheTTL := cfg.Scraper.CacheTTL
	if cacheTTL == 0 {
		cacheTTL = defaultCacheTTL
	}

	c := colly.NewCollector(
		colly.AllowedDomains("companiesmarketcap.com"),
		colly.Async(true),
		colly.CacheDir("./colly_cache"),
		colly.CacheExpiration(cacheTTL),
	)

	numCores := runtime.NumCPU()
	workers := cfg.Scraper.Workers
	if workers <= 0 {
		workers = numCores * 2
	}
	logger.Info("Detected %d logical cores. Assigning %d concurrent workers.", numCores, workers)

	delay := cfg.Scraper.Delay
	if delay <= 0 {
		delay = 500 * time.Millisecond
	}

	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*companiesmarketcap.com",
		Parallelism: workers,
		RandomDelay: delay,
	}); err != nil {
		return nil, fmt.Errorf("error setting limit rule: %w", err)
	}

	userAgent := cfg.Scraper.UserAgent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	}

	baseURL := cfg.Scraper.BaseURL
	if baseURL == "" {
		baseURL = "https://companiesmarketcap.com"
	}

	var pagesOnce sync.Once

	c.OnHTML("span.companies-count", func(e *colly.HTMLElement) {
		if targetPages != 0 {
			return
		}

		pagesOnce.Do(func() {
			text := strings.TrimSpace(e.Text)
			text = strings.ReplaceAll(text, ",", "")
			num, err := strconv.Atoi(text)
			if err == nil && num > 0 {
				totalPages := int(math.Ceil(float64(num) / float64(companiesPerPage)))
				atomic.StoreInt32(&totalNumPages, int32(totalPages))
				logger.Info("Detected %d companies. Will dynamically extract %d pages.", num, totalPages)

				for i := 2; i <= totalPages; i++ {
					url := fmt.Sprintf("%s/page/%d/", baseURL, i)
					if err := e.Request.Visit(url); err != nil {
						logger.Warn("Failed to visit %s: %v", url, err)
					}
				}
			} else {
				logger.Warn("Could not parse total companies count. Only page 1 will be extracted.")
			}
		})
	})

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		tds := e.DOM.Find("td")
		if tds.Length() < 7 {
			return
		}

		rankStr := sanitize(tds.Eq(1).Text())
		name := sanitize(e.DOM.Find("div.company-name").Text())
		code := sanitize(e.DOM.Find("div.company-code").Text())
		if code != "" {
			name = fmt.Sprintf("%s (%s)", name, code)
		}

		country := sanitize(e.DOM.Find("span.responsive-hidden").Text())
		marketCap := sanitize(tds.Eq(3).Text())
		price := sanitize(tds.Eq(4).Text())
		today := sanitize(tds.Eq(5).Text())

		rank, _ := strconv.Atoi(rankStr)

		company := domain.Company{
			Rank:      rank,
			Name:      name,
			MarketCap: marketCap,
			Price:     price,
			Today:     today,
			Country:   country,
		}

		mu.Lock()
		companies = append(companies, company)
		mu.Unlock()
	})

	c.OnScraped(func(r *colly.Response) {
		current := atomic.AddInt32(&extractedPages, 1)
		total := atomic.LoadInt32(&totalNumPages)

		if total > 0 {
			logger.Trace("[%d/%d] Successfully extracted page: %s", current, total, r.Request.URL.String())
		} else {
			logger.Trace("[%d/?] Successfully extracted page: %s", current, r.Request.URL.String())
		}
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Trace("Requesting: %s", r.URL.String())
		r.Headers.Set("User-Agent", userAgent)
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Error("Error on %s: %v", r.Request.URL, err)
	})

	entryURL := fmt.Sprintf("%s/page/1/", baseURL)
	if err := c.Visit(entryURL); err != nil {
		logger.Warn("Failed to visit %s: %v", entryURL, err)
	}

	if targetPages > 1 {
		for i := 2; i <= targetPages; i++ {
			url := fmt.Sprintf("%s/page/%d/", baseURL, i)
			if err := c.Visit(url); err != nil {
				logger.Warn("Failed to visit %s: %v", url, err)
			}
		}
	}

	logger.Info("Waiting for concurrent requests to finish...")
	c.Wait()

	return companies, nil
}
