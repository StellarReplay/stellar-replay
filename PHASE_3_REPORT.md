# Phase 3 Report — Capture

Status: COMPLETE

Phase 3 adds one controlled, read-only capture operation. It does not add a CLI,
replay engine, local server, retries, proxy support, authentication, or write
operations.

## Delivered

- Added `internal/capture` with explicit HTTPS endpoint configuration and one-shot
  JSON-RPC POST capture.
- Restricted capture to the four Phase 1 methods and scalar request IDs.
- Added default 30-second timeout and 8 MiB response bound, with no retry or live
  fallback behavior.
- Disabled implicit proxy use in the default transport and rejected redirects,
  URL userinfo/fragments, sensitive query fields, and private/loopback/link-local/
  unspecified destinations.
- Added hostname resolution checks at dial time for private destinations.
- Preserved JSON-RPC result/error envelopes and error data as opaque JSON.
- Added an explicit response sanitization hook before fixture sealing.
- Added `Record`, which uses Phase 2 validation, hashing, and atomic storage.
- Added local TLS-server and injected-transport tests for success, deterministic
  output, RPC errors, sanitization, unsupported methods, malformed/oversized
  responses, timeout, redirect rejection, no retries, and endpoint policy.
- Added `docs/CAPTURE.md` and synchronized README, architecture, roadmap,
  project context, and changelog claims.

## Verification

Executed locally from the core repository:

```text
gofmt -w internal/capture                         PASS
go test -p 1 -count=1 ./...                      PASS
go vet ./...                                     PASS
go test -race -p 1 -count=1 ./...                PASS
git diff --check                                 PASS
```

The default Windows Go cache produced intermittent environment-owned test-binary
cleanup/access warnings. The complete suite was rerun with a repository-local
`GOCACHE`, serial package execution, and the race detector; all packages passed.

## Acceptance result

- Only the configured capture operation contacts the endpoint: PASS.
- Unsupported methods are rejected before transport invocation: PASS.
- Captured successful and JSON-RPC error responses validate and hash: PASS.
- Timeout and response-size bounds are enforced: PASS.
- Redirects, unsafe endpoint forms, and private literal destinations are rejected: PASS.
- No retries or fallback endpoint behavior: PASS.
- Response sanitization is explicit and tested: PASS.
- Default tests use local/injected transports and require no live network: PASS.

## Boundary and risks

The sanitization hook is explicit rather than a claim of universal automatic
secret detection. Authentication headers, explicit proxies, private mirrors,
batch calls, simulation, writes, and additional RPC methods remain deferred.
Hostname policy is enforced by the default dialer; injected transports are a
testing/control seam and are not the default production path.

## Synchronization record

The implementation commit, hosted workflow, and final report checkpoint are added
after local verification and push. Phase 4 replay remains `NOT_STARTED`; no Phase
4 work is included here.
