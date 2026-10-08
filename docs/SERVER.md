# Local replay server contract

Phase 5 implements `internal/server`, a loopback-only HTTP JSON-RPC boundary over
the Phase 4 replay engine. It is intended to let existing Stellar RPC clients
point at a local fixture-backed endpoint. The CLI command that will manage this
server is Phase 6 work.

## Boundary

- The default bind address is `127.0.0.1:0`; configured addresses must be
  loopback-only.
- The server accepts `POST /` with one JSON-RPC request.
- Non-POST requests receive HTTP `405` and other paths receive HTTP `404`.
- JSON-RPC success and protocol-error envelopes use HTTP `200`; the envelope's
  JSON-RPC error code is authoritative.
- Request bodies are bounded to 8 MiB by default.
- Batch requests, trailing JSON, malformed JSON, unknown methods, invalid
  parameters, and replay misses are rejected clearly.
- The server never calls capture, a live RPC endpoint, a proxy, or a fallback
  source. Its only data source is the validated immutable replay engine.

The server uses `http.Server.Shutdown` and exposes `Wait` so shutdown releases the
listener and waits for active handlers. The HTTP handler is safe for concurrent
requests because the replay engine is safe for concurrent reads.

## Error behavior

JSON-RPC errors use the Phase 4 mapping: `-32600` for invalid requests, `-32601`
for unsupported methods, `-32602` for invalid parameters or misses, and `-32603`
for internal failures. A valid scalar request ID is echoed in the error envelope;
malformed or batch input uses JSON `null`.

Transport misuse is handled before JSON-RPC processing with normal HTTP `404` or
`405` responses. The server does not invent a successful response for a mismatch.

## Compatibility and limitations

The endpoint speaks the JSON-RPC request/response shape expected by the supported
Stellar RPC clients, but only the four frozen read-only methods are available.
TLS termination, public binding, batch processing, live fallback, authentication,
and transaction submission are intentionally excluded. Local clients should use
the server's returned loopback endpoint and the same method/parameter shapes as
the captured fixture.
