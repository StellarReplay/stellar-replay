# Phase 2 Report — Fixture schema and validation

Status: COMPLETE

Phase 2 established the canonical, versioned fixture contract without starting
capture, replay, or server implementation.

## Delivered

- Added `internal/fixture` with schema v1 models for provenance, JSON-RPC request,
  response/error, and SHA-256 integrity.
- Added recursive canonical JSON serialization so object ordering and formatting
  do not change a fixture hash while array order and scalar types remain material.
- Added strict validation for the four Phase 1 methods and their parameter shapes.
- Added HTTPS provenance checks, sensitive request-field rejection, unknown-field
  rejection on load, an 8 MiB fixture limit, and a 200-key ledger-entry limit.
- Added `Save`/`Load` with restrictive permissions, bounded reads, validation, hash
  verification, and temporary-file-plus-rename installation.
- Added valid, malformed, parameter, secret, modified-hash, canonicalization, and
  load/unknown-field tests.
- Documented the durable contract and limitations in `docs/FIXTURE_SCHEMA.md`.

## Verification

Executed locally from the core repository:

```text
gofmt -w internal/fixture/fixture.go internal/fixture/fixture_test.go   PASS
go test -p 1 ./...                                                   PASS
go vet ./...                                                        PASS
go test -race -p 1 ./...                                            PASS
git diff --check                                                    PASS
```

The default Windows Go build cache returned access-denied errors from a stale
environment-owned cache entry. The same commands were rerun successfully with a
repository-local `GOCACHE`, leaving no project artifact in the source tree.

## Acceptance result

- Identical logical fixtures hash identically: PASS.
- Modified fixtures fail integrity verification before use: PASS.
- Invalid, unsupported, malformed, oversized, and sensitive fixtures fail
  validation: PASS.
- Schema fields and limitations match the implementation and documentation: PASS.
- Network access is used by none of the Phase 2 code or tests: PASS.

## Boundary and risks

Capture, replay matching, and the local server remain unimplemented for their
later phases. Response JSON remains opaque; this package does not decode XDR.
On platforms where rename cannot replace an existing file, callers must rotate or
choose a new destination before `Save`, as documented in the schema contract.

## Synchronization record

The implementation commit was pushed to the public `StellarReplay/stellar-replay`
repository on `main`:

- implementation commit: `9ebba9d24a09e1f0e3f6db9064a3777953cba18d`
- implementation commit URL: <https://github.com/StellarReplay/stellar-replay/commit/9ebba9d24a09e1f0e3f6db9064a3777953cba18d>
- hosted CI run: <https://github.com/StellarReplay/stellar-replay/actions/runs/37721884110>
- hosted CI result: `success`

This report checkpoint is committed separately after hosted verification so the
remote report contains the evidence above. Phase 3 capture remains `NOT_STARTED`;
no Phase 3 work is included here.
