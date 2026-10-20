# Downstream integration: `stellar/js-stellar-sdk`

This example proves that a real Stellar/Soroban consumer can use Stellar
Replay's loopback server without contacting a live RPC provider. The downstream
project is Stellar's official [`stellar/js-stellar-sdk`](https://github.com/stellar/js-stellar-sdk),
which exposes the Soroban `rpc.Server#getLedgerEntries` API and tests that API
in its own repository.

The example uses the same public SDK method and a real SDK-compatible Soroban
ledger-key/XDR response. The only changed boundary is the server URL: it points
to `127.0.0.1`, where Stellar Replay serves the checked-in fixture.

## Run it

From the Stellar Replay repository root:

```text
go run ./cmd/stellar-replay serve \
  --fixture fixtures/edge-cases/get-ledger-entries-stellar-js-sdk.json \
  --listen 127.0.0.1:8787
```

In a second terminal:

```text
cd examples/downstream/stellar-js-sdk
npm install
npm run verify
```

Expected output:

```text
stellar/js-stellar-sdk consumed Stellar Replay fixture successfully
```

The example is intentionally read-only and offline after the fixture is
checked in. It does not deploy or invoke a contract, submit a transaction, or
expand Stellar Replay's frozen v0.1 RPC method set. It demonstrates the
consumer boundary that matters for deterministic Soroban test setup:
`rpc.Server#getLedgerEntries` receives a normal parsed SDK response from the
local replay server.

## Upstream evidence

The API shape is grounded in the upstream SDK source and tests:

- [`src/rpc/index.ts`](https://github.com/stellar/js-stellar-sdk/blob/main/src/rpc/index.ts)
  exports `RpcServer` as `rpc.Server`.
- [`get_ledger_entries.test.ts`](https://github.com/stellar/js-stellar-sdk/blob/main/test/unit/server/soroban/get_ledger_entries.test.ts)
  exercises `Server#getLedgerEntries`, ledger-entry XDR, and TTL metadata.
