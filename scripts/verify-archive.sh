#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 ARCHIVE" >&2
  exit 2
fi

archive="$1"
case "$archive" in
  *.tar.gz)
    actual=$(tar -tzf "$archive" | sed 's#^\./##; s#/$##' | sort | grep -v '^$' | paste -sd ' ' -)
    expected="LICENSE README.md stellar-replay"
    ;;
  *.zip)
    actual=$(unzip -Z1 "$archive" | sort | paste -sd ' ' -)
    expected="LICENSE README.md stellar-replay.exe"
    ;;
  *) echo "unsupported archive: $archive" >&2; exit 2 ;;
esac

if [[ "$actual" != "$expected" ]]; then
  printf 'unexpected archive contents for %s:\n' "$archive" >&2
  printf 'actual: %s\nexpected: %s\n' "$actual" "$expected" >&2
  exit 1
fi
echo "archive smoke test passed: $archive"
