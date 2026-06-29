#!/bin/sh
set -e
mkdir -p /workspace 2>/dev/null || true
chown -R appuser:appuser /workspace 2>/dev/null || true
cd /workspace
exec su-exec appuser /scraper
