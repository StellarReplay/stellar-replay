# JavaScript and TypeScript integration

The initial Node.js adapter lives at
[`integrations/typescript`](../integrations/typescript/). It translates
JavaScript/TypeScript calls into the released `stellar-replay` CLI contract.

The Go core remains authoritative for fixture schema v1, canonical hashing,
supported methods, validation, replay matching, and security policy. The
adapter does not copy the replay engine, add live-network fallback, or invent
pagination semantics.

The current compatibility target is core `v0.1.0` and Node.js `18+`.

Supported operations:

- `validateFixture(path)` validates schema and integrity offline.
- `inspectFixture(path)` prints fixture metadata offline.
- `replayFixture(path, request)` replays one fixture offline and returns the
  CLI's JSON-RPC response.

See [`integrations/typescript/README.md`](../integrations/typescript/README.md)
for usage and tests.
