# Protocol

Stellar Replay speaks the HTTP JSON-RPC shape used by Stellar RPC. It does not reimplement Stellar RPC; it stores and serves selected request/response interactions.

## v0.1 method scope

| Method | Purpose | v0.1 |
|---|---|---|
| `getHealth` | endpoint health | supported |
| `getLatestLedger` | latest known ledger | supported |
| `getNetwork` | network configuration | supported |
| `getLedgerEntries` | read ledger state by keys | supported |
| `getEvents`, `getLedgers`, `getTransactions`, `getTransaction`, `getFeeStats` | indexed/history or fee data | future |
| `simulateTransaction` | transaction simulation | future, review required |
| `sendTransaction` | transaction submission | explicitly excluded |

The selected set is deliberately small. Stellar's official method list and response shapes evolve; responses are preserved as JSON values, not converted into a home-grown chain model.

## Matching

The key is `(method, canonical(params))`. Missing params and JSON `null` are distinct according to the fixture schema. Object key order is normalized; array order and scalar types are preserved. A method mismatch, parameter mismatch, malformed request, unknown method, or missing fixture returns a clear JSON-RPC error and never a random fixture. Repeated identical requests return the same response. Concurrent requests are independent and deterministic.

The request id is echoed by the server but is not part of matching, so a test may use its own id. Batch JSON-RPC requests are out of scope for v0.1 unless explicitly added during implementation; the server must reject them clearly rather than partially process them.

## Errors

Malformed JSON uses a parse error. Unsupported methods use a method-not-found error. Valid but unmatched requests use a replay-miss error with method and a safe parameter diff. Fixture validation and integrity failures prevent server startup.

## Network boundary

`serve`, `replay`, and tests must not contact live Stellar RPC. `record` alone may contact one configured endpoint over HTTPS, with timeouts and no secret headers. Redirects, private endpoints, and proxy behavior require the explicit policy in Phase 1 before implementation.

