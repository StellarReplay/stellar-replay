#!/usr/bin/env bash
set -euo pipefail

go_command="${GO_COMMAND:-go}"
if ! command -v "$go_command" >/dev/null 2>&1 && command -v go.exe >/dev/null 2>&1; then
  go_command="go.exe"
fi

"$go_command" test ./examples -run TestOfflineReplayWorkflow -count=1
"$go_command" test ./internal/doccheck -count=1
"$go_command" build -trimpath ./cmd/stellar-replay
echo "offline smoke test passed"
