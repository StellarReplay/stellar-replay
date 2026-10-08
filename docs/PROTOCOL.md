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

## Frozen request contract

The v0.1 request contract is based on the official Stellar RPC reference checked on 2026-10-08:

| Method | Parameters | Request considerations | Response policy |
|---|---|---|---|
| `getHealth` | none | omit `params` or use `{}`; `null` and positional arrays are invalid | preserve the complete `result`, including retention fields and future fields |
| `getLatestLedger` | none | omit `params` or use `{}` | preserve ledger id, protocol version, sequence, close time, and XDR fields opaquely |
| `getNetwork` | none | omit `params` or use `{}` | preserve passphrase, protocol version, optional Friendbot URL, and future fields |
| `getLedgerEntries` | object with required `keys` and optional `xdrFormat` | `keys` is base64 LedgerKey strings, maximum 200; `xdrFormat` is `base64` or `json` | preserve entries and XDR/JSON values without decoding or imposing a response schema |

Stellar RPC requires named parameter objects rather than positional arrays, and its documentation states that XDR-encoded values are base64 strings. The capture layer records the wire request and response; it does not validate that a provider's response is semantically truthful.

For zero-parameter methods, replay canonicalization treats omitted `params` and an empty object as the same matching form to accommodate the official examples and SDK clients. `null`, arrays, extra parameters, and malformed JSON remain distinct invalid requests. For `getLedgerEntries`, object key order is normalized, but array order, key bytes, `xdrFormat`, and scalar types are preserved.

The response envelope is preserved as either a JSON-RPC `result` or `error`. During replay, the incoming request id is echoed in the response envelope; the id is not part of fixture matching. Stored result/error content is otherwise returned unchanged.

## Matching

The key is `(method, canonical(params))`, after the zero-parameter normalization above. Object key order is normalized; array order and scalar types are preserved. A method mismatch, parameter mismatch, malformed request, unknown method, or missing fixture returns a clear JSON-RPC error and never a random fixture. Repeated identical requests return the same response. Concurrent requests are independent and deterministic.

The request id is echoed by the server but is not part of matching, so a test may use its own id. Batch JSON-RPC requests are out of scope for v0.1 unless explicitly added during implementation; the server must reject them clearly rather than partially process them.

## Errors

Malformed JSON uses a parse error. Unsupported methods use a method-not-found error. Valid but unmatched requests use a replay-miss error with method and a safe parameter diff. Fixture validation and integrity failures prevent server startup.

## Network boundary

`serve`, `replay`, and tests must not contact live Stellar RPC. `record` alone may contact one explicitly configured endpoint over HTTPS, with bounded timeouts, bounded response size, and no secret headers. v0.1 does not follow redirects, accepts no URL userinfo, and rejects loopback, link-local, and private-network destinations by default; local test servers use injected transports rather than the public capture path. Implicit environment proxy use is disabled for the capture transport in v0.1. Private mirrors, authenticated endpoints, HTTP endpoints, and proxy configuration are future policy work, not silent exceptions.

The server binds to loopback by default. It never falls back to the live endpoint, even when a fixture is missing. Batch JSON-RPC requests are out of scope for v0.1 and are rejected as unsupported rather than partially processed.

## Error mapping

The replay server uses JSON-RPC errors where possible: `-32600` for an invalid request, `-32601` for an unsupported method, `-32602` for invalid parameters or a replay miss, and a distinct implementation error for invalid fixture state. HTTP transport status is not used as the fixture matching key; the JSON-RPC envelope is authoritative for replay semantics.

## Research record

The decisions above are grounded in the official [JSON-RPC structure](https://developers.stellar.org/docs/data/apis/rpc/api-reference/structure/json-rpc), [getHealth](https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getHealth), [getLatestLedger](https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getLatestLedger), [getNetwork](https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getNetwork), and [getLedgerEntries](https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods/getLedgerEntries) references, plus the official [client SDK documentation](https://developers.stellar.org/docs/tools/sdks/client-sdks). The JavaScript SDK exposes an RPC server constructor accepting a URL, and the Go SDK exposes `rpcclient.NewClient(url, httpClient)`, so a loopback HTTP endpoint is a compatible integration boundary without making an SDK dependency part of the core.

