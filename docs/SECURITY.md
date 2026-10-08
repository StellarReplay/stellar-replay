# Product Security Model

## Scope

Stellar Replay captures read-only RPC interactions. It never signs, submits, or executes transactions in v0.1. A fixture is not a cryptographic proof of historical chain state or endpoint honesty.

## Threats and controls

| Threat | Control |
|---|---|
| secret capture | reject secret-like fields and document redaction; no signing inputs |
| fixture tampering | canonical content hash and validation before replay |
| accidental live fallback | replay has no live client dependency |
| request confusion | exact method/params matching and useful misses |
| malicious fixture | bounded parsing, schema checks, no execution |
| sensitive shared data | explicit sanitization and provenance review |
| network abuse | HTTPS, timeouts, bounded bodies, endpoint policy |

Account IDs, contract IDs, request data, response data, URLs, timestamps, and logs may still be sensitive. Sanitization must be opt-in/explicitly visible and must not claim that redacted fixtures reproduce the original response unless matching semantics remain valid.

## Out of scope

No private-key support, transaction submission, package execution, sandbox escape claims, cryptographic chain proofs, or full RPC security gateway is planned for v0.1.

