package scraper

import (
	"context"

	"marketcap-acquisition-engine/internal/config"
	"marketcap-acquisition-engine/internal/domain"
)

type Scraper interface {
	Run(ctx context.Context, cfg *config.Config) ([]domain.Company, error)
}
