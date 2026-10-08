# Supported RPC methods

Stellar Replay v0.1 supports four read-only JSON-RPC methods. Request parameters
are validated according to the frozen protocol; response objects remain opaque JSON
so provider fields can evolve without a second chain model.

| Method | Parameters | Purpose | Limitation |
|---|---|---|---|
| `getHealth` | omitted or `{}` | endpoint health | captures one observed health response |
| `getLatestLedger` | omitted or `{}` | latest known ledger | ledger/XDR fields are not decoded |
| `getNetwork` | omitted or `{}` | network identity/configuration | provenance is metadata, not proof |
| `getLedgerEntries` | `{ "keys": [base64...], "xdrFormat": "base64"\|"json" }` | read ledger state | at most 200 keys; XDR remains opaque |

For zero-parameter methods, omitted params and `{}` match equivalently. `null`,
arrays, extra fields, malformed JSON, and unsupported methods are rejected. For
`getLedgerEntries`, object key order is normalized, while key-array order, scalar
types, and `xdrFormat` remain significant.

The request ID is excluded from matching and echoed in replay responses. Stored
results and JSON-RPC errors are returned unchanged except for the capture-phase
sanitization that may have been applied before the fixture was written.

Out of scope for v0.1: `sendTransaction`, simulation, Horizon APIs, batch requests,
contract execution, and arbitrary RPC method passthrough. See
[`PROTOCOL.md`](PROTOCOL.md) for the full decision record.
