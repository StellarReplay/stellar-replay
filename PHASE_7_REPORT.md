# Phase 7 Report — Deterministic test workflow

Status: COMPLETE

Phase 7 proves the normal developer value proposition with sanitized checked-in
fixtures and a fully offline integration example. It adds no live capture, SDK
dependency, or new RPC behavior.

## Delivered

- Added four checked-in, schema-valid, integrity-sealed fixtures for
  `getHealth`, `getLatestLedger`, `getNetwork`, and `getLedgerEntries`.
- Used fixed testnet provenance and synthetic safe response values; no credentials,
  private data, or claims of current chain state are included.
- Added `examples/offline_replay_test.go`, which loads the fixture matrix through
  the replay engine, starts the loopback server, verifies repeated byte-stable
  responses, and exercises concurrent clients.
- Added a transport guard that rejects every non-loopback destination during the
  integration example.
- Updated `examples/README.md`, `tests/README.md`, `fixtures/README.md`, README,
  roadmap, context, and changelog with the reproducible workflow.

## Verification

Executed locally from the core repository:

```text
go test -count=1 ./examples                               PASS
go test -p 1 -count=1 ./...                               PASS
go vet ./...                                               PASS
go test -race -p 1 -count=1 ./...                          PASS
git diff --check                                           PASS
```

The Windows environment intermittently reports an access-denied cleanup warning
when Go removes a temporary test binary after a passing run. The complete suite
was rerun with a repository-local `GOCACHE`, serial package execution, and the
race detector; all packages passed.

## Acceptance result

- New contributor can run the offline replay path without credentials: PASS.
- All four frozen methods have checked-in safe fixtures: PASS.
- Fixture matrix validates and verifies integrity: PASS.
- Repeated replay is byte-stable: PASS.
- Concurrent replay clients receive identical responses: PASS.
- Non-loopback network destinations are blocked by the example transport guard:
  PASS.
- Default tests use no live Stellar endpoint: PASS.

## Boundary and risks

Fixtures are synthetic examples of the schema and protocol shape, not chain-truth
claims. Live capture remains explicit through `record`. Future fixture additions
must be sanitized, validated, integrity-sealed, and reviewed for publish safety.

## Synchronization record

The implementation commit, hosted workflow, and final report checkpoint are added
after local verification and push. Phase 8 security and sanitization hardening
remains `NOT_STARTED`; no Phase 8 work is included here.
