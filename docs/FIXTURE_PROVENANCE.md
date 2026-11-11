# Fixture provenance guidance

Fixture provenance explains how a fixture was obtained. It is evidence about
the capture process, not proof that the provider response is correct or that
the chain state remains current.

## Required provenance

Every fixture records:

- `endpoint`: an HTTPS RPC endpoint without URL userinfo or credential-bearing
  query parameters;
- `network`: a human-readable network label such as `testnet`;
- `toolVersion`: the Stellar Replay version or development identifier used to
  capture the interaction; and
- `capturedAt`: an RFC3339 capture timestamp.

`chainContext` is optional and may contain reviewed, non-secret context such as
a ledger or network identifier. It must never contain credentials, headers,
seed phrases, private keys, or unreviewed private URLs.

## Review before sharing

Before committing or publishing a fixture:

1. Confirm the endpoint is public or explicitly approved for sharing.
2. Remove credentials, authorization material, private URLs, and identifying
   request data.
3. Run `stellar-replay validate --fixture PATH` and inspect the fixture.
4. Confirm the response is synthetic or safe to redistribute.
5. Record the intended network and why the fixture exists.

The validator rejects known credential-bearing fields, but sanitization is not a
universal data-loss prevention system. Human review remains required. The
canonical schema and integrity rules are defined in `FIXTURE_SCHEMA.md`.

## Reproducibility record

For a maintained fixture, the surrounding documentation should identify the
core version, method, network label, fixture purpose, and expected replay
outcome. Do not record live credentials or assume a fixture can be recaptured
identically after ledger state or provider retention changes.
