#!/bin/sh
set -e
chown -R appuser:appuser /workspace 2>/dev/null || true
cd /workspace
exec su-exec appuser /scraper "$@"
