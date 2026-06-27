# MarketCap Acquisition Engine

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![Colly](https://img.shields.io/badge/Framework-Colly-blue)](#)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](#)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Author](https://img.shields.io/badge/by-RLFreddy-gray?logo=github)](https://github.com/RLFreddy)

> A concurrent data acquisition engine built in Go. Extracts, processes, and serializes market capitalization data from companiesmarketcap.com with caching and race-condition-safe output.

## Overview

This project uses native HTTP requests via the Colly framework with Go's concurrency model. It discovers pagination dynamically, handles state in memory to avoid I/O bottlenecks, and outputs a sorted CSV dataset.

## Features

- **Concurrent extraction:** NumCPU × 2 workers, async requests via Colly
- **Caching layer:** Built-in cache with 24-hour TTL, prevents redundant requests
- **Pagination discovery:** Reads total company count and calculates pages automatically
- **Thread-safe serialization:** Results aggregated via mutex-guarded slice, sorted by rank, written sequentially
- **Docker-ready:** Multi-stage build, ~8MB image, runs as non-root user (UID 1001)

## Quick Start (Docker)

**1. Clone**

```bash
git clone https://github.com/RLFreddy/marketcap-acquisition-engine
cd marketcap-acquisition-engine
```

**2. Build**

```bash
docker build -t companies-scraper .
```

**3. Run (first page)**

```bash
docker run --rm -t -v ./output:/workspace companies-scraper -pages 1
```

**3.1. Or extract all pages**

```bash
docker run --rm -t -v ./output:/workspace companies-scraper
```

## Local Execution (Go)

```bash
go mod download
go build -o scraper ./cmd/scraper/

# Extract first page (100 companies)
./scraper -pages 1

# Extract all pages
./scraper
```

## Output

Generates `companies_YYYY-MM-DD.csv`:

| Column     | Type    | Example       |
| ---------- | ------- | ------------- |
| Rank       | Integer | 1             |
| Name       | String  | NVIDIA (NVDA) |
| Market Cap | String  | $4.663 T      |
| Price      | String  | $192.53       |
| Today      | String  | 1.64%         |
| Country    | String  | USA           |

## Project Structure

```
├── cmd/scraper/main.go          # Entry point, CLI flags
├── internal/
│   ├── domain/company.go        # Data structures
│   ├── scraper/colly_scraper.go # Extraction, caching, concurrency
│   ├── exporter/csv_exporter.go # CSV output, sorting
│   └── logger/logger.go         # Colored stdout logging
├── Dockerfile
├── entrypoint.sh
└── go.mod
```

## License

MIT

---

Built by [RLFreddy](https://github.com/RLFreddy)
