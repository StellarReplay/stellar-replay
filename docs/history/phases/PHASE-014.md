# Phase 14 — Action contract and secure integration

Status: COMPLETE

## Objective

Provide a thin, versioned GitHub Action around the released core binary without
copying the Go engine or introducing live RPC behavior.

## Completed work

- Published `StellarReplay/stellar-replay-action@v0.1.0`.
- Pinned the Action to released core v0.1.0 and fixture schema v1.
- Verified downloaded archives against `SHA256SUMS`.
- Supported Linux amd64, Windows amd64, and macOS amd64.
- Preserved validate/replay-only behavior and fail-closed input handling.
- Added immutable-tag usage examples and a line-ending policy for shell tests.

## Verification

Local metadata and contract checks passed. The Windows path downloaded the real
core release archive, verified its checksum, and validated the fixture offline.
Hosted Action CI passed on the protected repository workflow.

## Evidence

- [Action PR #4](https://github.com/StellarReplay/stellar-replay-action/pull/4)
- [Action v0.1.0 release](https://github.com/StellarReplay/stellar-replay-action/releases/tag/v0.1.0)

## Non-goals

No wallet credentials, signing, live capture, copied Go implementation, or
unreleased-core dependency was added.
