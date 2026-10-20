# Stellar Replay CLI adapter

This package provides a small Node.js/TypeScript integration for the released
`stellar-replay` CLI. It delegates fixture validation, inspection, and replay
to the Go binary, so the core repository remains the source of truth for the
fixture schema, hashing, matching, and supported RPC methods.

Compatibility target: adapter `0.1.1`, core `v0.1.1`, and Node.js `18+`.
The adapter is offline-only and does not implement a second replay engine.

The package is independently distributable from npm as
`@stellar-replay/cli-adapter`. It requires a compatible `stellar-replay`
binary on `PATH` (or an explicit executable path); installing the package does
not download a binary or contact a live RPC endpoint.

```js
import { createAdapter } from "@stellar-replay/cli-adapter";

const replay = createAdapter({ executable: "stellar-replay" });
replay.validateFixture("fixtures/get-health.json");
const response = replay.replayFixture("fixtures/get-health.json", {
  method: "getHealth", params: {}, id: "node-test",
});
console.log(response);
```

Run package tests with `npm test`.
