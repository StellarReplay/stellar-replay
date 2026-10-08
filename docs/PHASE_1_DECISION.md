# Phase 1 Decision Record — Research and MVP Freeze

Date: 2026-10-08
Status: `COMPLETE`

## Decision

Freeze v0.1 around four read-only Stellar RPC methods: `getHealth`, `getLatestLedger`, `getNetwork`, and `getLedgerEntries`. The core captures and replays JSON-RPC envelopes, not a simulated ledger or decoded XDR model.

## Evidence reviewed

- Official Stellar RPC method catalog and current method pages.
- Official JSON-RPC structure and error-code guidance.
- Official client SDK documentation and URL-configurable JavaScript/Go RPC clients.
- Current Stellar RPC naming and bounded-history context.

## Decisions

1. POST JSON-RPC 2.0 is the only v0.1 transport.
2. Zero-parameter methods accept omitted `params` and `{}` as equivalent; `null`, arrays, and extra fields do not match.
3. `getLedgerEntries` requires named object parameters, supports up to 200 base64 LedgerKeys, and preserves `xdrFormat` exactly.
4. Request ids are ignored for matching and echoed from the incoming request during replay.
5. Response results and errors are preserved as opaque JSON values; XDR is not decoded by the core.
6. Replay never contacts a live endpoint and never falls back when a fixture is missing.
7. Capture requires an explicit HTTPS endpoint, rejects redirects and URL userinfo, rejects private/loopback/link-local destinations by default, disables implicit proxy use, and uses bounded timeouts/body sizes.
8. Local replay serves loopback by default. Local test servers are injected in tests rather than reached through the public capture policy.
9. Batch requests, writes, simulation, Horizon, and unsupported RPC methods are excluded from v0.1.

## Deferred questions

Private mirrors, authenticated RPC headers, explicit proxy support, batch replay, response normalization, and additional methods require a later security/protocol decision. They must not enter implementation by implication.

## Sources

- https://developers.stellar.org/docs/data/apis/rpc/api-reference/structure/json-rpc
- https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getHealth
- https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getLatestLedger
- https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getNetwork
- https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getLedgerEntries
- https://developers.stellar.org/docs/tools/sdks/client-sdks
- https://pkg.go.dev/github.com/stellar/go-stellar-sdk/clients/rpcclient
- https://stellar.github.io/js-stellar-sdk/reference/network-rpc/

