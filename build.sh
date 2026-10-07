#!/usr/bin/env bash
# LessonMode — Linux/macOS dan Windows .exe yasash (cross-compile).
#
# Foydalanish:
#   ./build.sh                      # oddiy exe (config.json kerak bo'ladi)
#   BOT_TOKEN=xxx ADMIN_ID=123 ./build.sh   # token/ID exe ichiga joylashtiriladi
#
# Tayyor fayl: dist/LessonMode.exe
set -euo pipefail

OUT="${OUT:-dist/LessonMode.exe}"
mkdir -p "$(dirname "$OUT")"

LDFLAGS="-s -w -H=windowsgui"
[ -n "${BOT_TOKEN:-}" ] && LDFLAGS="$LDFLAGS -X main.buildBotToken=${BOT_TOKEN}"
[ -n "${ADMIN_ID:-}" ] && LDFLAGS="$LDFLAGS -X main.buildAdminID=${ADMIN_ID}"

echo "LessonMode build qilinmoqda -> $OUT"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o "$OUT" .

ls -lh "$OUT"
echo "Tayyor."
