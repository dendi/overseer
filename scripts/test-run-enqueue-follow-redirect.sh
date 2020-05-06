#!/usr/bin/env bash
set -Eeuxmo pipefail
DIR="$(dirname "$(command -v greadlink >/dev/null 2>&1 && greadlink -f "$0" || readlink -f "$0")")"

TMP_FILE=$(mktemp)
cat >"$TMP_FILE" <<EOL
https://httpstat.us/301 must run http with status 200 with follow-redirect true
EOL

go run "$DIR/.." enqueue "$TMP_FILE"