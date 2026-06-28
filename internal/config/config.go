package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Scraper ScraperConfig `yaml:"scraper"`
	Output  OutputConfig  `yaml:"output"`
}

type ScraperConfig struct {
	BaseURL   string        `yaml:"base_url"`
	Pages     int           `yaml:"pages"`
	CacheTTL  time.Duration `yaml:"cache_ttl"`
	Delay     time.Duration `yaml:"delay"`
	Workers   int           `yaml:"workers"`
	UserAgent string        `yaml:"user_agent"`
}

type OutputConfig struct {
	Dir            string `yaml:"dir"`
	FilenamePrefix string `yaml:"filename_prefix"`
}

func Default() *Config {
	return &Config{
		Scraper: ScraperConfig{
			BaseURL:   "https://companiesmarketcap.com",
			Pages:     0,
			CacheTTL:  24 * time.Hour,
			Delay:     500 * time.Millisecond,
			Workers:   0,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		Output: OutputConfig{
			Dir:            ".",
			FilenamePrefix: "companies_",
		},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func MergeWithFlags(cfg *Config, pages int) *Config {
	if cfg == nil {
		cfg = Default()
	}
	if pages != 0 {
		cfg.Scraper.Pages = pages
	}
	return cfg
}
