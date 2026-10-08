#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 3 || $# -gt 4 ]]; then
  echo "usage: $0 VERSION GOOS GOARCH [OUTPUT_DIR]" >&2
  exit 2
fi

version="$1"
goos="$2"
goarch="$3"
output_dir="${4:-dist}"
package_name="stellar-replay-${version}-${goos}-${goarch}"
stage="${output_dir}/.stage-${package_name}"
go_command="${GO_COMMAND:-go}"
if ! command -v "$go_command" >/dev/null 2>&1 && command -v go.exe >/dev/null 2>&1; then
  go_command="go.exe"
fi

rm -rf "$stage"
mkdir -p "$stage" "$output_dir"
binary="${stage}/stellar-replay"
CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" "$go_command" build \
  -trimpath \
  -ldflags="-s -w -X github.com/stellar-replay/stellar-replay/internal/cli.Version=${version}" \
  -o "$binary" ./cmd/stellar-replay
cp LICENSE README.md "$stage/"

if [[ "$goos" == "windows" ]]; then
  mv "$binary" "${stage}/stellar-replay.exe"
  powershell.exe -NoProfile -NonInteractive -Command \
    "Compress-Archive -Path '${stage}/*' -DestinationPath '${output_dir}/${package_name}.zip' -Force"
else
  tar -C "$stage" -czf "${output_dir}/${package_name}.tar.gz" .
fi

rm -rf "$stage"
echo "${output_dir}/${package_name}"
