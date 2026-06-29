package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Scraper ScraperConfig `yaml:"scraper" validate:"required"`
	Output  OutputConfig  `yaml:"output"`
}

type ScraperConfig struct {
	Pages      int           `yaml:"pages"      validate:"min=0"`
	CacheTTL   time.Duration `yaml:"cache_ttl"  validate:"min=0"`
	CacheDir   string        `yaml:"cache_dir"`
	Delay      time.Duration `yaml:"delay"      validate:"min=0"`
	Workers    int           `yaml:"workers"    validate:"min=0"`
	UserAgent  string        `yaml:"user_agent"`
	RetryCount int           `yaml:"retry_count" validate:"min=0"`
	RetryDelay time.Duration `yaml:"retry_delay" validate:"min=0"`
}

type OutputConfig struct {
	Dir            string `yaml:"dir"`
	FilenamePrefix string `yaml:"filename_prefix"`
}

var validate = validator.New()

func (c *Config) Validate() error {
	return validate.Struct(c)
}

func Default() *Config {
	return &Config{
		Scraper: ScraperConfig{
			Pages:      0, // 0 = auto-detect all pages
			CacheTTL:   24 * time.Hour,
			CacheDir:   "./colly_cache",
			Delay:      500 * time.Millisecond,
			Workers:    0, // 0 = auto (NumCPU * 2)
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			RetryCount: 3,
			RetryDelay: 1 * time.Second,
		},
		Output: OutputConfig{
			Dir:            ".",
			FilenamePrefix: "companies_",
		},
	}
}

func Load(paths ...string) (*Config, error) {
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		cfg := Default()
		dec := yaml.NewDecoder(bytes.NewReader(data))
		dec.KnownFields(true)
		if err := dec.Decode(cfg); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("%s: validation failed: %w", p, err)
		}
		return cfg, nil
	}
	return nil, fmt.Errorf("no config file found (tried: %v)", paths)
}
