# Demo Plan

This is an intended demo, not claimed output. It becomes executable after the corresponding roadmap phases are complete.

1. Use `record` against a documented read-only testnet RPC endpoint for `getLatestLedger`.
2. Write a versioned fixture with provenance and integrity hash.
3. `inspect` the human-readable interaction.
4. `validate` the schema and hash.
5. Start `serve` on loopback.
6. Send the same JSON-RPC request with curl or an SDK configured for the loopback URL.
7. Show the deterministic response and a deliberate mismatch error.
8. Run a local test against the server.
9. Disconnect the network and repeat the test to demonstrate no live dependency.

No output is written here until implementation produces it.

