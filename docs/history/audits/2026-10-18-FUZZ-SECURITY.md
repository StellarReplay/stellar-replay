# Fuzz and security follow-up — 2026-10-18

## Scope

This follow-up exercises the two highest-risk JSON boundaries in the v0.1 core:
fixture request validation and response sanitization. It complements, but does
not replace, the historical hardening audit and its accepted limitations.

## Changes

- Added `FuzzValidateRequestNoPanic` for arbitrary `getLedgerEntries` parameter
  bytes, including sensitive-key and malformed-JSON seeds.
- Added `FuzzResponseNoPanic` for arbitrary response-result bytes, including
  credential-shaped keys, bearer material, arrays, and `null`.
- Preserved the existing explicit regression tests for endpoint credential keys,
  fixture sensitive fields, response redaction, JSON-RPC ID matching, bounds,
  and integrity validation.

## Evidence

The bounded fuzz runs completed without a panic or finding:

```text
go test ./internal/fixture -run=^$ -fuzz=FuzzValidateRequestNoPanic -fuzztime=10s
go test ./internal/sanitize -run=^$ -fuzz=FuzzResponseNoPanic -fuzztime=10s
```

Unit tests for both packages passed before fuzzing. The hosted core matrix is
the cross-platform verification gate for the merged change.

## Limits

Fuzzing demonstrates parser robustness for the exercised boundaries; it does
not prove universal secret detection or that captured data is safe without
maintainer review. Signed release/fixture provenance remains a future design
evaluation, not a v0.1 claim.
