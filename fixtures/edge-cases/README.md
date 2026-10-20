# Edge-case fixtures

These fixtures are synthetic, sanitized protocol cases. They exercise how
Stellar Replay preserves or rejects difficult provider behavior; they are not
claims about a current network response or a universal provider error code.

## Included cases

- `get-ledger-entries-retention.json` is a valid fixture containing a
  provider-style ledger-retention error. Replay preserves it as a deterministic
  JSON-RPC error response.
- `get-ledger-entries-stellar-js-sdk.json` contains SDK-compatible ledger-key
  and entry XDR used by the real downstream example in
  `examples/downstream/stellar-js-sdk/`.

Malformed envelopes and unsupported pagination shapes live under
`tests/fixtures/invalid/` because they are intentionally rejected and must not
be mistaken for usable fixtures.
