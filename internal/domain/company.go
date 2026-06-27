package domain

// Company represents a single row extracted from the CompaniesMarketCap table.
type Company struct {
	Rank      int
	Name      string
	MarketCap string
	Price     string
	Today     string
	Country   string
}
