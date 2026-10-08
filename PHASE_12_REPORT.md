# Phase 12 Report — v0.1.0 release

Status: COMPLETE

## Release

- Repository: StellarReplay/stellar-replay
- Tag: v0.1.0
- Release: [Stellar Replay v0.1.0](https://github.com/StellarReplay/stellar-replay/releases/tag/v0.1.0)
- Release commit: [cbbd6f0](https://github.com/StellarReplay/stellar-replay/commit/cbbd6f0433af9441bded45bb5f69e3544901d785)
- Release workflow: [GitHub Actions run 37778472965](https://github.com/StellarReplay/stellar-replay/actions/runs/37778472965) — success
- Release is public, non-draft, and non-prerelease.

## Published artifacts

| Artifact | SHA-256 |
|---|---|
| `stellar-replay-0.1.0-linux-amd64.tar.gz` | `8a66874705ff3eac4e7222dec78ac68edb0e521129e02f5acf934e3a2906d29e` |
| `stellar-replay-0.1.0-windows-amd64.zip` | `a0b18b6283cb797df17cbe73277e34a7f95858596c657ccda11bc8fb2bb7d6cb` |
| `stellar-replay-0.1.0-darwin-amd64.tar.gz` | `60ceb8a56b409ac88e3d4713b80a1482274c398a78b87f2072a72af2032ebcae` |

`SHA256SUMS` was downloaded from the published release and independently
verified successfully against all three artifacts.

## Independent post-release verification

- Linux archive contents: PASS.
- macOS archive contents: PASS.
- Windows archive contents: PASS.
- Extracted Windows binary reported `0.1.0`: PASS.
- Extracted Windows binary `help` output listed the CLI commands: PASS.
- Release archive contents are limited to the intended binary, `LICENSE`, and
  `README.md`.

## Pre-release verification

The tagged source passed before publication:

    go test -p 1 -count=1 ./...
    go test -race -p 1 -count=1 ./...
    go vet ./...
    go build -trimpath ./cmd/stellar-replay
    go test -count=1 ./internal/doccheck
    git diff --check

The release preparation commit’s hosted matrix passed on Linux, Windows, and
macOS in [run 37777979553](https://github.com/StellarReplay/stellar-replay/actions/runs/37777979553).

## Recovery note

The first tag workflow attempt was blocked by a macOS Bash portability defect in
the archive smoke script. No release was published from that attempt. The
portable fix was committed as `cbbd6f0`, passed hosted CI in [run
37778329773](https://github.com/StellarReplay/stellar-replay/actions/runs/37778329773),
and the tag was moved to that corrected commit before the successful release
workflow.

The workflow emitted a non-blocking Node.js 20 deprecation annotation for
`actions/upload-artifact@v4`; all release jobs and artifact checks passed. This
can be addressed as routine workflow maintenance without changing the v0.1.0
artifacts.

## Scope and next state

This release does not claim signing, package-registry distribution, additional
platform architectures, ecosystem adoption, or unsupported RPC methods. Phase
12 is complete. No further phase is authorized by this request.
