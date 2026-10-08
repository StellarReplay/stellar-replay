param(
    [Parameter(Mandatory = $true)][string]$Version,
    [Parameter(Mandatory = $true)][string]$GOOS,
    [Parameter(Mandatory = $true)][string]$GOARCH,
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"
$packageName = "stellar-replay-$Version-$GOOS-$GOARCH"
$stage = Join-Path $OutputDir ".stage-$packageName"
$archive = Join-Path $OutputDir "$packageName.zip"

if (Test-Path -LiteralPath $stage) { Remove-Item -LiteralPath $stage -Recurse -Force }
New-Item -ItemType Directory -Path $stage -Force | Out-Null
New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

$env:CGO_ENABLED = "0"
$env:GOOS = $GOOS
$env:GOARCH = $GOARCH
go build -trimpath -ldflags "-s -w -X github.com/stellar-replay/stellar-replay/internal/cli.Version=$Version" -o (Join-Path $stage "stellar-replay.exe") ./cmd/stellar-replay
Copy-Item LICENSE, README.md -Destination $stage
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $archive -Force
Remove-Item -LiteralPath $stage -Recurse -Force
Write-Output $archive
