# Phase 1 Report — Research and MVP Freeze

Date: 2026-10-08
Status: `COMPLETE` after local and hosted verification

## Objective

Validate the v0.1 RPC assumptions and freeze the method, request, response, matching, endpoint, and security boundaries before implementation.

## Decisions completed

- v0.1 supports `getHealth`, `getLatestLedger`, `getNetwork`, and `getLedgerEntries` only.
- JSON-RPC 2.0 over HTTP POST is the only v0.1 transport.
- Zero-parameter requests accept omitted `params` and `{}` equivalently; `null`, arrays, and extra fields are invalid.
- `getLedgerEntries` uses named object parameters, base64 LedgerKey strings, a maximum of 200 keys, and optional `xdrFormat` of `base64` or `json`.
- Request ids are excluded from fixture matching and echoed from the incoming replay request.
- Results and errors are preserved as opaque JSON; XDR is not decoded by the core.
- Replay is offline-only, has no live fallback, rejects batch requests, and binds to loopback by default.
- Capture requires explicit HTTPS, rejects redirects and URL userinfo, rejects private/loopback/link-local destinations by default, disables implicit proxy use, and uses bounded timeouts/body sizes.
- Private mirrors, authenticated headers, explicit proxy support, batch replay, writes, simulation, Horizon, and additional methods are deferred.

## Research evidence

- JSON-RPC structure and error codes: https://developers.stellar.org/docs/data/apis/rpc/api-reference/structure/json-rpc
- `getHealth`: https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getHealth
- `getLatestLedger`: https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getLatestLedger
- `getNetwork`: https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getNetwork
- `getLedgerEntries`: https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getLedgerEntries
- SDK overview: https://developers.stellar.org/docs/tools/sdks/client-sdks
- Go RPC client URL constructor: https://pkg.go.dev/github.com/stellar/go-stellar-sdk/clients/rpcclient
- JavaScript RPC server URL constructor: https://stellar.github.io/js-stellar-sdk/reference/network-rpc/

## Acceptance results

- Current RPC method surface inspected: PASS
- Each supported method has request/response/matching/limitation policy: PASS
- JSON-RPC error and id behavior documented: PASS
- SDK endpoint compatibility documented: PASS
- Capture redirect/private/proxy policy documented: PASS
- MVP scope frozen before implementation: PASS
- Core protocol and decision record updated: PASS
- Action/docs cross-references audited: PASS

## Verification

- Core `git diff --check`: PASS
- Core `go test ./...`: PASS; no implementation packages exist by design
- Core `go vet ./...`: PASS; no implementation packages exist by design
- Cross-repository reference search: PASS; no stale repository name found
- Phase 0 GitHub gate: PASS before Phase 1

## Next gate

Phase 2 — Fixture Schema and Validation. It may begin only after this report and the pushed commits are verified. Phase 2 is the first implementation phase.

