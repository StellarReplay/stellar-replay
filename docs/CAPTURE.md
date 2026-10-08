# Capture contract

Phase 3 implements one controlled capture operation in `internal/capture`. It is
the only current component that can contact a configured Stellar RPC endpoint.
There is not yet a CLI command; command integration is Phase 6 work.

## Operation

`capture.Config` requires an explicit HTTPS `Endpoint`, network label, tool
version, supported method, and scalar JSON-RPC id. `Capture` performs exactly one
HTTP POST and returns a sealed [fixture](FIXTURE_SCHEMA.md). `Record` uses the same
operation and writes the result through the fixture package's atomic storage.

The capture client:

- sends JSON-RPC 2.0 with `Content-Type: application/json`;
- supports only `getHealth`, `getLatestLedger`, `getNetwork`, and
  `getLedgerEntries`;
- uses a 30-second default timeout and an 8 MiB response bound;
- makes no retries and never falls back to another endpoint;
- rejects redirects, URL userinfo, URL fragments, sensitive query fields, and
  literal private, loopback, link-local, or unspecified addresses;
- disables implicit proxy use for its default transport;
- resolves hostnames at dial time and skips private, loopback, link-local, and
  unspecified destination addresses;
- preserves JSON-RPC error envelopes, including error data;
- invokes an optional response sanitization hook before validation and hashing.

The production transport verifies TLS with TLS 1.2 as its minimum. Tests inject a
transport or rewrite requests to a local TLS test server; this keeps default tests
offline without weakening the production endpoint policy.

## Response handling

HTTP status alone does not replace JSON-RPC semantics. A bounded response body must
be valid JSON and must contain exactly one JSON-RPC `result` or `error` when the
fixture is sealed. The response result and error data remain opaque JSON. XDR is
not decoded by capture.

The sanitization hook is deliberately explicit. Phase 3 does not claim that every
possible secret representation is automatically redacted. A hook may mutate or
reject the response, and fixture validation rejects prohibited sensitive request
fields before a fixture can be written.

## Trust boundary

Capture is opt-in and network-bound. Replay and validation must remain offline.
Credentials, authorization headers, proxies, writes, transaction submission,
simulation, batch requests, and unsupported methods are outside this phase and
require a later protocol/security decision.
