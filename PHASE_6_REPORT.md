# Phase 6 Report — CLI integration

Status: COMPLETE

Phase 6 adds the `stellar-replay` command-line workflow over the existing capture,
fixture, replay, and server packages. It does not add new RPC methods, shell
completion, SDKs, or implicit live-network behavior.

## Delivered

- Added `cmd/stellar-replay` as the standalone executable entry point.
- Added `internal/cli` with stable exit codes: `0` success, `1` operational
  failure, and `2` usage error.
- Implemented `record` with explicit HTTPS endpoint, method, network, output,
  timeout, JSON params, params-file/stdin support, and request ID flags.
- Implemented `validate` for strict schema and integrity verification.
- Implemented `inspect` for safe fixture metadata and response-kind output.
- Implemented `replay` with fixture loading, flags or request-file/stdin input,
  JSON-RPC output, and offline deterministic errors.
- Implemented `serve` with repeatable fixture flags, loopback listen address,
  signal-aware lifecycle, and clean shutdown.
- Added command help/version output and useful usage/operational errors.
- Replaced the planned README walkthrough with executable commands and explicit
  live-network boundaries.
- Added end-to-end CLI tests for command parsing, exit codes, validate/inspect,
  offline replay, request files, missing fixtures, and capture endpoint policy.

## Verification

Executed locally from the core repository:

```text
gofmt -w internal/cli cmd/stellar-replay          PASS
go test -p 1 -count=1 ./...                       PASS
go vet ./...                                      PASS
go test -race -p 1 -count=1 ./...                 PASS
go build ./cmd/stellar-replay                     PASS
binary --help smoke test                          PASS
git diff --check                                  PASS
```

The Windows environment intermittently reports an access-denied cleanup warning
when Go removes a temporary test binary after a passing run. The complete suite
was rerun with a repository-local `GOCACHE`, serial package execution, and the
race detector; all packages passed.

## Acceptance result

- All five planned commands exist: PASS.
- README walkthrough uses real implemented flags: PASS.
- Offline validate/inspect/replay workflow works end to end: PASS.
- Replay CLI performs no network access: PASS.
- Record remains the only network-capable command and requires explicit HTTPS:
  PASS.
- Stable success, failure, and usage exit codes: PASS.
- Fixture output path and request/params file inputs: PASS.
- Server lifecycle is wired to interrupt/termination signals: PASS.

## Boundary and risks

The CLI is intentionally standard-library-only and remains an internal v0.1
surface while usage stabilizes. `record` may contact one configured HTTPS endpoint;
all other commands are offline. Cross-platform release packaging, deterministic
checked-in examples, and broader workflow proof remain later phases.

## Synchronization record

The implementation commit was pushed to the public `StellarReplay/stellar-replay`
repository on `main`:

- implementation commit: `e622ff521004ce23477ca83f1468981dfcc54336`
- implementation commit URL: <https://github.com/StellarReplay/stellar-replay/commit/e622ff521004ce23477ca83f1468981dfcc54336>
- hosted CI run: <https://github.com/StellarReplay/stellar-replay/actions/runs/37724646177>
- hosted CI result: `success`

This report checkpoint is committed separately after hosted verification so the
remote report contains the evidence above. Phase 7 deterministic test workflow
remains `NOT_STARTED`; no Phase 7 work is included here.
