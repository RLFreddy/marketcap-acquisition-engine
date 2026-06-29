# MarketCap Acquisition Engine

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev)
[![CI](https://github.com/RLFreddy/marketcap-acquisition-engine/actions/workflows/ci.yml/badge.svg)](https://github.com/RLFreddy/marketcap-acquisition-engine/actions)
[![golangci-lint](https://img.shields.io/badge/lint-golangci--lint-4BC51C)](#)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](#)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Author](https://img.shields.io/badge/by-RLFreddy-gray?logo=github)](https://github.com/RLFreddy)

> A concurrent data acquisition engine built in Go. Extracts, processes, and serializes market capitalization data from companiesmarketcap.com with caching, retry backoff, and race-condition-safe output.

![Demo](assets/demo.gif)

## Features

- **Concurrent extraction** — NumCPU × 2 workers with async Colly requests
- **Caching** — 24h TTL, avoids redundant HTTP requests
- **Auto-pagination** — reads total count, calculates pages dynamically
- **Retry with backoff** — exponential backoff on failed requests
- **Docker-ready** — multi-stage build, ~22MB image, non-root user

## Docker

```bash
docker build -t companies-scraper .

docker run --rm -t \
  -v ./output:/workspace \
  -v ./config.yaml:/etc/scraper/config.yaml \
  companies-scraper
```

The CSV file is written to `./output/companies_YYYY-MM-DD.csv`.

## Local Execution

```bash
# Ensure config.yaml is in the current directory
go build -o scraper ./cmd/scraper/

./scraper
```

## Configuration

The scraper looks for `config.yaml` in this order:

- `./config.yaml` (local development)
- `/etc/scraper/config.yaml` (Docker mount)

If neither is found, built-in defaults are used.

```yaml
scraper:
  base_url: "https://companiesmarketcap.com"
  pages: 0
  workers: 0
  delay: 500ms
  cache_dir: "./colly_cache"
  cache_ttl: 24h
  retry_count: 3
  retry_delay: 1s

output:
  dir: "."
  filename_prefix: "companies_"
```

## Commands

```bash
make lint
make test
make cover
make build
make clean
```

## Output

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
├── .github/workflows/ci.yml
├── .golangci.yml
├── Makefile
├── config.yaml
├── cmd/scraper/main.go
├── internal/
│   ├── config/config.go
│   ├── domain/company.go
│   ├── scraper/
│   │   ├── scraper.go
│   │   ├── colly_scraper.go
│   │   ├── scraper_test.go
│   │   └── testdata/
│   ├── exporter/
│   │   ├── csv_exporter.go
│   │   └── csv_exporter_test.go
│   └── logger/logger.go
├── Dockerfile
├── entrypoint.sh
├── go.mod
└── go.sum
```

## License

MIT

---

Built by [RLFreddy](https://github.com/RLFreddy)
