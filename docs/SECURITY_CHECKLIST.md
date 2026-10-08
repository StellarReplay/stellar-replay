# Phase 8 security checklist

This checklist records the implemented v0.1 hardening boundary. It is not a
claim that arbitrary captured data is automatically safe to publish.

## Controls verified

- [x] Request and response bodies are bounded before or during transport/storage.
- [x] Fixture loads reject oversized input, unknown fields, malformed JSON, and
      integrity mismatches.
- [x] Capture rejects non-HTTPS endpoints, URL userinfo/fragments, sensitive
      query keys, literal private/loopback/link-local/unspecified addresses,
      redirects, and implicit proxy use.
- [x] Hostname destinations are checked at dial time before connection.
- [x] Request parameters reject sensitive field names recursively.
- [x] Captured response result/error data receives conservative default redaction
      for sensitive keys, PEM private-key material, and bearer values.
- [x] A custom sanitization hook can further redact or reject a response.
- [x] Fixture writes use temporary files, synchronization, rename installation,
      and restrictive `0600` permissions where the platform exposes them.
- [x] Replay and local serving have no live-network or fallback dependency.
- [x] CLI errors and replay miss diagnostics do not print raw request parameters.
- [x] Checked-in fixtures contain synthetic values and no credentials/private data.
- [x] Concurrent capture/replay/server paths have race-tested adversarial cases.

## Remaining limitations

Secret detection is conservative and key/pattern based; it cannot recognize every
secret encoding or semantic credential. Review captured fixtures before sharing.
Response redaction can change replay semantics, so the resulting fixture must be
validated and behaviorally reviewed. Authentication headers, private mirrors,
explicit proxies, transaction credentials, cryptographic chain proofs, and a full
security gateway remain out of scope.

User-selected filesystem paths are honored by the CLI. Callers must not direct
fixture output into sensitive directories or trust unreviewed symlinked/shared
locations. Atomic rename protects installation integrity but is not an access
control boundary for a directory already writable by other users.

## Test evidence

The Phase 8 tests cover secret-like keys and values, malformed response data,
oversized capture requests, oversized fixtures, endpoint policy, file permissions,
path creation, malformed fixtures, concurrent access, and no-network replay.
