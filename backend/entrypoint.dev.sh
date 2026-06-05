#!/bin/sh
set -e

echo "==> Syncing Go modules..."
go mod tidy

echo "==> Starting air (hot reload)..."
exec air -c .air.toml
