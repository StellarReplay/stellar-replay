# Soroban RPC integration example

This example shows a safe, offline path for a Soroban-oriented
`getLedgerEntries` workflow. It uses the checked-in synthetic fixture and never
requires a wallet, secret, credential, or live RPC endpoint.

## Validate the reviewed fixture

From the core repository root:

```text
go build -o stellar-replay ./cmd/stellar-replay
stellar-replay validate --fixture fixtures/get-ledger-entries.json
stellar-replay inspect --fixture fixtures/get-ledger-entries.json
```

The fixture contains one base64 ledger key and a synthetic response. Its
`testnet` provenance identifies the shape of the source interaction; it does
not prove current chain state.

## Replay through the local server

Start the loopback-only server:

```text
stellar-replay serve --fixture fixtures/get-ledger-entries.json \
  --listen 127.0.0.1:8787
```

In another terminal, send the same request shape with a new JSON-RPC ID:

```text
curl -s http://127.0.0.1:8787/ \
  -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":7,"method":"getLedgerEntries","params":{"keys":["AQIDBA=="],"xdrFormat":"base64"}}'
```

The replay response is deterministic and comes from the fixture. The recorded
endpoint is not contacted. Request matching preserves key-array order, scalar
types, and `xdrFormat`; object-key order is canonicalized.

## Boundaries

The v0.1 contract accepts at most 200 non-empty base64 keys and does not define
pagination or cursors. Split larger workloads into separate reviewed fixtures.
`sendTransaction`, simulation, contract execution, Horizon APIs, batch requests,
and generic RPC passthrough remain outside this product contract.

The executable offline proof also exercises this method:

```text
go test ./examples -run TestOfflineReplayMethodMatrix -count=1
```
