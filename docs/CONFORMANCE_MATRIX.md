# Stellar Replay conformance matrix

This is the release-facing compatibility matrix for core `v0.1.1`. It describes
what is verified by the repository's tests and what remains deliberately outside
the product contract. A checked box means the behavior is covered by automated
tests or hosted release verification; it does not claim that a provider is
truthful or that every provider implements the same edge case.

## Core protocol

| Area | v0.1.1 contract | Evidence | Status |
|---|---|---|---|
| JSON-RPC envelope | JSON-RPC 2.0 request and response envelopes | `internal/fixture`, `internal/server` tests | Verified |
| Request/response IDs | Stored IDs must match exactly, including scalar type | `internal/fixture/fixture_test.go` | Verified |
| Supported methods | `getHealth`, `getLatestLedger`, `getNetwork`, `getLedgerEntries` | `docs/METHODS.md`, fixture tests | Verified |
| Unsupported methods | Writes, simulation, Horizon, batches, and generic passthrough rejected | capture, replay, server tests | Verified |
| JSON-RPC errors | Error envelopes are stored and replayed deterministically | edge-case fixtures and replay tests | Verified |
| Pagination | Cursor/page fields rejected in schema v1 | invalid fixture corpus and schema tests | Verified boundary |

## Fixture and replay behavior

| Area | Verified behavior | Evidence | Status |
|---|---|---|---|
| Schema version | Only schema version `1` is accepted by core v0.1.1 | table-driven schema tests | Verified |
| Canonical hashing | Object-key order and whitespace do not change SHA-256 | fixture tests | Verified |
| Integrity gate | Missing or modified integrity prevents use | fixture and replay tests | Verified |
| Provenance | HTTPS endpoint, network, capture time, and tool version required | schema and capture tests | Verified |
| Secret handling | Credential-bearing request and endpoint fields rejected; response sanitization is explicit | security tests/checklist | Verified with documented limits |
| Replay fallback | No live-network fallback exists | replay/server tests and architecture docs | Verified |
| Storage safety | Bounded loads, strict fields, atomic writes, restrictive permissions | fixture tests | Verified |

## Release and integration surfaces

| Surface | Compatibility target | Verification |
|---|---|---|
| Core CLI | `stellar-replay` `v0.1.1` | GitHub release archives, checksums, and Linux/Windows/macOS smoke tests |
| GitHub Action | Action `v0.1.1`, core `v0.1.1` | Action CI metadata, negative-contract, and three-platform offline smoke tests |
| Node.js adapter | `@stellar-replay/cli-adapter` `0.1.1`, Node.js 18+ | adapter tests, npm package contents, downstream SDK example |
| Documentation | Describes core `v0.1.1` and Action `v0.1.1` | Docs CI and cross-repository link review |

## Known gaps and non-claims

- The matrix does not claim support for transaction submission, simulation,
  contract execution, Horizon APIs, or arbitrary RPC methods.
- Provider-specific pagination, retention codes, and response fields are not
  normalized into a universal chain model.
- Fixture hashes detect modification; they do not prove endpoint truth.
- Release archives are checksum-protected but not publisher-signed in v0.1.x.
- External adoption and third-party contribution evidence is tracked separately
  in `docs/ADOPTION_FEEDBACK.md`.
