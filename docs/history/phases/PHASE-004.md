# Phase 4 Report — Replay engine

Status: COMPLETE

Phase 4 adds deterministic offline replay over validated fixtures. It does not add
the HTTP server, CLI integration, fuzzy matching, response mutation, or live
fallback behavior.

## Delivered

- Added `internal/replay` with immutable, concurrent-safe fixture indexing.
- Added integrity-gated `New` and file-based `Load` construction.
- Implemented exact `(method, canonical(params))` matching.
- Normalized omitted and empty params only for the three zero-parameter methods.
- Preserved method-specific `getLedgerEntries` validation and array/scalar
  semantics through the Phase 2 request validator.
- Excluded request IDs from matching and echoed the incoming ID in cloned response
  envelopes.
- Rejected duplicate canonical request keys rather than selecting nondeterministically.
- Added deterministic JSON-RPC-compatible errors for invalid requests, unsupported
  methods, invalid parameters/misses, and internal failures.
- Added safe miss diagnostics that do not expose raw request parameters.
- Added exact-one-request parsing with batch and trailing-data rejection.
- Added concurrent-read tests and response cloning to prevent returned data from
  mutating indexed fixtures.
- Added the shared fixture validation/canonicalization seam and documented the
  replay contract in `docs/REPLAY.md`.

## Verification

Executed locally from the core repository:

```text
gofmt -w internal/fixture/fixture.go internal/replay             PASS
go test -p 1 -count=1 ./...                                     PASS
go vet ./...                                                     PASS
go test -race -p 1 -count=1 ./...                               PASS
git diff --check                                                 PASS
```

The Windows environment intermittently reports an access-denied cleanup warning
when Go removes a temporary test binary, even after the test process passes. The
complete suite was rerun with a repository-local `GOCACHE`, serial package
execution, and the race detector; all packages passed.

## Acceptance result

- Exact canonical matching and object-key normalization: PASS.
- Array order and scalar/type distinctions remain material: PASS.
- Method mismatch, unsupported method, malformed request, and replay miss errors:
  PASS.
- Repeated requests return the same stored content with the incoming ID echoed:
  PASS.
- Tampered or invalid fixtures fail before engine creation: PASS.
- Duplicate fixture keys fail closed: PASS.
- Concurrent replay reads are deterministic and race-free: PASS.
- Replay code performs no network access: PASS.

## Boundary and risks

The engine returns typed errors and an error-envelope helper for the future server;
it does not open sockets or decide HTTP status codes. Fuzzy matching, response
normalization, request batching, fixture mutation, and live fallback remain out of
scope. The public CLI and local replay server are later phases.

## Synchronization record

The implementation commit was pushed to the public `StellarReplay/stellar-replay`
repository on `main`:

- implementation commit: `75b5088e0f16c39a6b3541291593fdc7839e468b`
- implementation commit URL: <https://github.com/StellarReplay/stellar-replay/commit/75b5088e0f16c39a6b3541291593fdc7839e468b>
- hosted CI run: <https://github.com/StellarReplay/stellar-replay/actions/runs/37723243370>
- hosted CI result: `success`

This report checkpoint is committed separately after hosted verification so the
remote report contains the evidence above. Phase 5 local replay server remains
`NOT_STARTED`; no Phase 5 work is included here.
