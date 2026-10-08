# Replay engine contract

Phase 4 implements deterministic, offline replay in `internal/replay`. It indexes
validated Phase 2 fixtures and returns stored responses without contacting a live
endpoint. The CLI and local HTTP server are later phases.

## Matching

The replay key is:

```text
(method, canonical(params))
```

Request IDs are deliberately excluded from the key. The incoming ID is echoed in
the returned response while the stored result or error remains unchanged.

Canonical matching recursively sorts JSON object keys, preserves array order and
scalar types, and rejects multiple/trailing JSON values. Omitted params and `{}`
are equivalent only for `getHealth`, `getLatestLedger`, and `getNetwork`.
`getLedgerEntries` retains its required object and key-array semantics. Fuzzy
matching, response mutation, request ordering, and live fallback are not used.

## Fixture loading

`New` validates fixture schema and SHA-256 integrity before indexing. `Load` reads
fixture files through the strict Phase 2 loader. Tampered, malformed, unsupported,
or duplicate fixtures fail before an engine is returned. Duplicate canonical keys
are rejected instead of selecting an arbitrary response.

The engine is immutable after construction. Replay clones response JSON before
returning it, so callers cannot mutate the stored fixture through a returned
slice. Read access is protected for concurrent callers.

## Errors

The engine exposes deterministic JSON-RPC-compatible error codes for the future
server boundary:

- `-32600`: malformed or invalid JSON-RPC request;
- `-32601`: unsupported method;
- `-32602`: invalid parameters or no matching fixture;
- `-32603`: internal replay failure.

Miss diagnostics include only the method and a safe statement that canonical
parameters did not match. Raw parameters are never placed in an error message.
`ErrorResponse` converts a typed replay error to an error envelope for the later
HTTP server.

## Go API

- `New(fixtures)` validates and indexes in-memory fixtures.
- `Load(paths...)` validates and indexes fixture files.
- `Replay(request)` handles a parsed non-batch request.
- `ReplayJSON(raw)` parses exactly one request and rejects batches/trailing data.
- `ErrorResponse(id, err)` creates a JSON-RPC error envelope.

The package is internal until the public CLI and server behavior stabilize.
