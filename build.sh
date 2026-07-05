#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

echo "==> 构建前端"
cd frontend
pnpm install
pnpm run build

echo "==> 拷贝 dist 到后端"
rm -rf ../backend/web/dist
cp -r dist ../backend/web/dist

echo "==> 构建后端二进制"
cd ../backend
go build -o ../jiaocheng-web ./cmd/server

echo "==> 构建抓取脚本"
go build -o ../scrape ./cmd/scrape

echo "==> 完成"
ls -lh "$ROOT/jiaocheng-web" "$ROOT/scrape"
