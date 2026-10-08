param([Parameter(Mandatory = $true)][string]$Archive)
$ErrorActionPreference = "Stop"
$temp = Join-Path ([System.IO.Path]::GetTempPath()) ("stellar-replay-archive-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $temp | Out-Null
try {
    Expand-Archive -LiteralPath $Archive -DestinationPath $temp
    $entries = @(Get-ChildItem -LiteralPath $temp -File | Select-Object -ExpandProperty Name | Sort-Object)
} finally {
    Remove-Item -LiteralPath $temp -Recurse -Force -ErrorAction SilentlyContinue
}
$expected = @("LICENSE", "README.md", "stellar-replay.exe")
if (($entries -join "`n") -cne ($expected -join "`n")) {
    throw "Unexpected archive contents. Actual: $($entries -join ', '); expected: $($expected -join ', ')"
}
Write-Output "archive smoke test passed: $Archive"
