package domain

import "testing"

func TestCompanyCreation(t *testing.T) {
	c := Company{Rank: 1, Name: "Apple", MarketCap: "$3T", Price: "$200", Today: "+1%", Country: "US"}
	if c.Rank != 1 {
		t.Errorf("expected Rank 1, got %d", c.Rank)
	}
	if c.Name != "Apple" {
		t.Errorf("expected Name 'Apple', got %s", c.Name)
	}
	if c.MarketCap != "$3T" {
		t.Errorf("expected MarketCap '$3T', got %s", c.MarketCap)
	}
	if c.Price != "$200" {
		t.Errorf("expected Price '$200', got %s", c.Price)
	}
	if c.Today != "+1%" {
		t.Errorf("expected Today '+1%%', got %s", c.Today)
	}
	if c.Country != "US" {
		t.Errorf("expected Country 'US', got %s", c.Country)
	}
}

func TestCompanyZeroValue(t *testing.T) {
	var c Company
	if c.Rank != 0 || c.Name != "" {
		t.Errorf("expected zero value Company, got %+v", c)
	}
}
