package exporter

import (
	"os"
	"testing"

	"marketcap-acquisition-engine/internal/domain"
)

func TestExportToCSV(t *testing.T) {
	companies := []domain.Company{
		{Rank: 2, Name: "Beta Inc", MarketCap: "$200B", Price: "$100", Today: "+1%", Country: "US"},
		{Rank: 1, Name: "Alpha Corp", MarketCap: "$500B", Price: "$250", Today: "-0.5%", Country: "CN"},
		{Rank: 3, Name: "Gamma LLC", MarketCap: "$50B", Price: "$50", Today: "+2%", Country: "JP"},
	}

	path := t.TempDir() + "/test_output.csv"
	if err := ExportToCSV(companies, path); err != nil {
		t.Fatalf("ExportToCSV failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	expectedHeaders := "Rank,Name,Market Cap,Price,Today,Country\n"
	if content[:len(expectedHeaders)] != expectedHeaders {
		t.Errorf("expected headers %q, got %q", expectedHeaders, content[:len(expectedHeaders)])
	}

	rows := content[len(expectedHeaders):]
	expectedBody := "1,Alpha Corp,$500B,$250,-0.5%,CN\n2,Beta Inc,$200B,$100,+1%,US\n3,Gamma LLC,$50B,$50,+2%,JP\n"
	if rows != expectedBody {
		t.Errorf("unexpected CSV body:\ngot:  %q\nwant: %q", rows, expectedBody)
	}
}

func TestExportToCSVEmpty(t *testing.T) {
	path := t.TempDir() + "/empty.csv"
	if err := ExportToCSV(nil, path); err != nil {
		t.Fatalf("ExportToCSV with nil slice failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	expectedHeaders := "Rank,Name,Market Cap,Price,Today,Country\n"
	if content != expectedHeaders {
		t.Errorf("expected only headers, got %q", content)
	}
}

func TestExportToCSVBadPath(t *testing.T) {
	err := ExportToCSV(nil, "/nonexistent/dir/file.csv")
	if err == nil {
		t.Fatal("expected error for bad path, got nil")
	}
}
