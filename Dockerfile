# syntax=docker/dockerfile:1
ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath -ldflags="-s -w" -o /scraper ./cmd/scraper/

FROM alpine:3.20

RUN apk add --no-cache ca-certificates su-exec && \
    adduser -D -u 1001 appuser

COPY --from=builder /scraper /scraper
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
