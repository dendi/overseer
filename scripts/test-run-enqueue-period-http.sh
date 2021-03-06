#!/usr/bin/env bash
set -Eeuxmo pipefail
DIR="$(dirname "$(command -v greadlink >/dev/null 2>&1 && greadlink -f "$0" || readlink -f "$0")")"

TMP_FILE=$(mktemp)
cat >"$TMP_FILE" <<EOL
https://google.com must run http with status any with pt-duration 5s with pt-sleep 100ms with pt-threshold 0%
EOL

go run "$DIR/.." enqueue "$TMP_FILE"