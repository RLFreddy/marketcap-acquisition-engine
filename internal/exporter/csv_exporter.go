package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	"marketcap-acquisition-engine/internal/domain"
	"marketcap-acquisition-engine/internal/logger"
)

// ExportToCSV sorts the given companies by Rank and writes them to the specified CSV file.
func ExportToCSV(companies []domain.Company, fileName string) error {
	logger.Info("Sorting %d companies for CSV export...", len(companies))
	
	// Sort by Rank
	sort.Slice(companies, func(i, j int) bool {
		return companies[i].Rank < companies[j].Rank
	})

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating csv file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			logger.Error("CSV flush error: %v", err)
		}
	}()

	// Write Headers
	headers := []string{"Rank", "Name", "Market Cap", "Price", "Today", "Country"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}

	// Write Data
	for _, comp := range companies {
		row := []string{
			fmt.Sprintf("%d", comp.Rank),
			comp.Name,
			comp.MarketCap,
			comp.Price,
			comp.Today,
			comp.Country,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row %d: %w", comp.Rank, err)
		}
	}

	return nil
}
