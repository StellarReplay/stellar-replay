# Phase 5 Report — Local replay server

Status: COMPLETE

Phase 5 adds a loopback-only HTTP JSON-RPC server over the validated Phase 4 replay
engine. It does not add the CLI, public binding, TLS termination, batch support,
live fallback, or transaction behavior.

## Delivered

- Added `internal/server` with loopback-only configuration and a default
  `127.0.0.1:0` listener.
- Added one-request `POST /` JSON-RPC handling with an 8 MiB default body bound.
- Added HTTP `404` path and `405` method handling with `Allow: POST`.
- Added JSON-RPC error envelopes for malformed JSON, batches, unsupported methods,
  invalid parameters, and replay misses.
- Preserved valid request IDs in error envelopes and used JSON `null` when an ID
  cannot be safely extracted.
- Kept protocol-level errors at HTTP 200; transport misuse uses normal HTTP status.
- Added concurrent request handling through the concurrency-safe replay engine.
- Added `Start`, `Addr`, `Endpoint`, `Wait`, and context-aware `Shutdown` APIs.
- Ensured shutdown closes the listener and waits for active handlers.
- Added startup, success, malformed, batch, unsupported, mismatch, body-limit,
  path/method, concurrency, loopback-policy, and shutdown tests.
- Documented the server contract in `docs/SERVER.md` and synchronized project
  documentation and phase status.

## Verification

Executed locally from the core repository:

```text
gofmt -w internal/server                                      PASS
go test -p 1 -count=1 ./...                                  PASS
go vet ./...                                                  PASS
go test -race -p 1 -count=1 ./...                            PASS
git diff --check                                              PASS
```

The Windows environment intermittently reports an access-denied cleanup warning
when Go removes a temporary test binary after a passing run. The complete suite
was rerun with a repository-local `GOCACHE`, serial package execution, and the
race detector; all packages passed.

## Acceptance result

- Server starts on loopback and returns a fixture-backed JSON-RPC response: PASS.
- Invalid JSON, batches, unsupported methods, mismatches, and oversized bodies are
  rejected clearly: PASS.
- Request IDs are preserved or safely represented as `null`: PASS.
- Concurrent requests are deterministic and race-free: PASS.
- Shutdown releases the listener and prevents further requests: PASS.
- Server startup rejects non-loopback addresses: PASS.
- Server code has no live-network or fallback path: PASS.

## Boundary and risks

The server accepts only `POST /` and the four frozen read-only methods. HTTP error
status conventions for protocol errors are intentionally simple and documented;
the JSON-RPC envelope is authoritative. CLI lifecycle management, richer client
compatibility, TLS, authentication, and public binding remain later or deferred
work.

## Synchronization record

The implementation commit was pushed to the public `StellarReplay/stellar-replay`
repository on `main`:

- implementation commit: `43763576db7a3e8c69647c4992048713f242bd14`
- implementation commit URL: <https://github.com/StellarReplay/stellar-replay/commit/43763576db7a3e8c69647c4992048713f242bd14>
- hosted CI run: <https://github.com/StellarReplay/stellar-replay/actions/runs/37724060283>
- hosted CI result: `success`

This report checkpoint is committed separately after hosted verification so the
remote report contains the evidence above. Phase 6 CLI integration remains
`NOT_STARTED`; no Phase 6 work is included here.
