# Fixture schema v1

The canonical fixture format is owned by `stellar-replay`. Other repositories may
consume this format, but must not silently redefine it. A schema change requires a
new `schemaVersion` and a compatibility note.

## Document shape

Each fixture is one JSON-RPC request/response interaction:

```json
{
  "schemaVersion": "1",
  "capturedAt": "2026-10-08T12:00:00Z",
  "provenance": {
    "endpoint": "https://soroban-testnet.stellar.org",
    "network": "testnet",
    "toolVersion": "dev",
    "chainContext": {}
  },
  "request": {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getLatestLedger",
    "params": {}
  },
  "response": {
    "jsonrpc": "2.0",
    "id": 1,
    "result": {}
  },
  "integrity": {
    "algorithm": "sha256",
    "sha256": "..."
  }
}
```

Required top-level fields are `schemaVersion`, `capturedAt`, `provenance`,
`request`, `response`, and `integrity`. `chainContext`, request `params`, and
response `error.data` are optional JSON values. The response contains exactly one
of `result` or `error`.

The request and response use JSON-RPC 2.0. IDs are scalar JSON values and the
fixture request and response IDs must match without coercing strings and numbers.
IDs are not used as the replay matching key. Replay will echo the incoming request ID. The
frozen v0.1 methods are `getHealth`, `getLatestLedger`, `getNetwork`, and
`getLedgerEntries`. Zero-parameter methods accept omitted params or `{}`. The
`getLedgerEntries` object requires a non-empty `keys` array of at most 200
non-empty base64 strings and optionally accepts `xdrFormat` of `base64` or
`json`.

## Integrity and canonical JSON

`Seal` computes SHA-256 over canonical JSON for the complete fixture with the
`integrity` field omitted. Canonical JSON sorts object keys recursively, preserves
array order and JSON scalar types, preserves JSON numbers, and rejects trailing or
multiple JSON values. The stored digest is lowercase hexadecimal. `Validate` and
`Load` recompute and verify the digest before a fixture can be used.

Pretty-printed storage is allowed; whitespace and object-key order do not affect
the digest. The hash detects modification; it does not prove that an endpoint was
truthful or that the captured state remains current.

## Validation and storage limits

- Provenance requires an HTTPS endpoint without URL userinfo or credential-bearing
  query parameters (including generic token keys and underscore/hyphen/case
  variants), plus a network label and tool version.
- Unknown fields are rejected when loading from disk.
- Fixtures larger than 8 MiB are rejected.
- Sensitive request parameter keys such as private keys, seed phrases, passwords,
  generic or authorization tokens, API keys, and bearer values are rejected
  recursively.
- Response JSON is preserved as opaque JSON. XDR decoding is intentionally outside
  the schema package; base64 and JSON response representations are not interpreted
  here.
- `Save` writes a restrictive `0600` temporary file, syncs it, and installs it by
  rename. On platforms where rename cannot replace an existing file, callers must
  choose a new path or remove/rotate the old file according to their own policy.

Fixtures are evidence of a captured interaction, not chain truth. Capture,
sanitization policy beyond the schema checks, replay matching, and the local server
are separate phases.

## Go API

The implementation is in `internal/fixture`:

- `Seal` and `Validate` establish and verify the contract.
- `Hash` and `CanonicalJSON` expose deterministic integrity primitives.
- `Save` and `Load` provide bounded, strict, atomic file handling.
- `Supported` reports the frozen method allowlist.

This package is internal until a stable public Go API is deliberately designed.
