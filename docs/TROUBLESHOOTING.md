# Troubleshooting

## `fixture ... is invalid`

Run:

```text
stellar-replay validate --fixture path/to/fixture.json
```

Check schema version, required provenance, supported method parameters, response
shape, and the SHA-256 integrity field. Do not hand-edit a fixture after sealing;
re-record or regenerate it through the fixture API.

## `no replay fixture matched`

Replay matches `(method, canonical(params))`, not request IDs. Confirm the method,
parameter object, array order, scalar types, and `xdrFormat`. For zero-parameter
methods, omitted params and `{}` are equivalent. Use `inspect` to confirm the
fixture method and provenance.

## `endpoint must be an HTTPS URL` or endpoint policy errors

Only `record` contacts a live endpoint. Use an explicit HTTPS URL without userinfo,
fragments, sensitive query fields, redirects, implicit proxy use, or private/
loopback/link-local destinations. Local testing belongs behind the replay server,
not the capture endpoint policy.

## `request body exceeds` or `RPC response exceeds`

The capture and server boundaries are intentionally bounded to protect memory and
avoid accidental large-data capture. Reduce the request/response or split the
workflow; there is no unrestricted mode in v0.1.

## `address must be loopback-only`

`serve` is local infrastructure and binds loopback by design. Use
`127.0.0.1:PORT` or `[::1]:PORT`; public binding is not supported in v0.1.

## I need a live capture

Live capture is opt-in and requires a supported read-only method, an explicit
network label, and a safe HTTPS endpoint:

```text
stellar-replay record --rpc-url https://soroban-testnet.stellar.org \
  --network testnet --method getLatestLedger --output fixtures/latest-ledger.json
```

Review and sanitize the resulting fixture before sharing it. Never provide secret
keys, authorization values, or signing material.

## CI or test failure

Run the same local gates used by CI:

```text
gofmt -w .
go test ./...
go vet ./...
go test -race ./...
git diff --check
```

The default suite is network-free. A Windows temporary test-binary cleanup warning
can come from the local Go environment; rerun with a repository-local `GOCACHE`
and serial execution if needed. Do not treat an unverified hosted result as a
passing phase gate.
