$ErrorActionPreference = "Stop"

go test ./examples -run TestOfflineReplayWorkflow -count=1
go test ./internal/doccheck -count=1
go build -trimpath ./cmd/stellar-replay
if (Test-Path -LiteralPath "stellar-replay.exe") {
    Remove-Item -LiteralPath "stellar-replay.exe" -Force
}
Write-Output "offline smoke test passed"
