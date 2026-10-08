# Fixtures

Checked-in fixtures must be sanitized, versioned, schema-valid, and safe to
publish. The canonical schema is documented in
[`docs/FIXTURE_SCHEMA.md`](../docs/FIXTURE_SCHEMA.md) and implemented by the
internal `fixture` package.

The repository includes four deterministic, sanitized examples covering
`getHealth`, `getLatestLedger`, `getNetwork`, and `getLedgerEntries`. They use
fixed testnet provenance and synthetic safe response values; they are not claims
about current chain state and contain no credentials or private data.

Phase 7 loads this fixture matrix in the offline integration example. New
fixtures must pass `stellar-replay validate` and remain safe to publish before
being checked in.

