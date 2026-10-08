# Offline demo

This demo is executable from a fresh checkout without credentials or a live
network. It uses the sanitized checked-in fixtures.

## 1. Build and validate

```text
go build -o stellar-replay ./cmd/stellar-replay
stellar-replay validate --fixture fixtures/get-latest-ledger.json
stellar-replay inspect --fixture fixtures/get-latest-ledger.json
```

The validation command verifies the schema and integrity hash. Inspection prints
metadata without dumping the response payload.

## 2. Replay offline

```text
stellar-replay replay --fixture fixtures/get-latest-ledger.json \
  --method getLatestLedger --id 7
```

The response is served from the fixture and remains unchanged across repeated
runs. Replay does not construct a capture client or contact the recorded endpoint.

## 3. Serve locally

In one terminal:

```text
stellar-replay serve --fixture fixtures/get-latest-ledger.json \
  --listen 127.0.0.1:8787
```

In another terminal, send the same JSON-RPC shape:

```text
curl -s http://127.0.0.1:8787/ \
  -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":7,"method":"getLatestLedger"}'
```

An unmatched request returns a JSON-RPC error rather than a random fixture:

```text
curl -s http://127.0.0.1:8787/ \
  -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":8,"method":"getNetwork"}'
```

Stop the server with Ctrl-C. The server is loopback-only and has no live fallback.

## 4. Run the automated proof

```text
go test ./examples -run TestOfflineReplayWorkflow -count=1
```

The test loads all four checked-in fixtures, uses a transport guard that rejects
non-loopback destinations, verifies repeated byte-stable responses, and exercises
concurrent clients.

## Optional live capture

Only if explicitly authorized and network access is available, `record` can create
a new fixture from a supported HTTPS read-only endpoint. Review and validate the
result before sharing it:

```text
stellar-replay record --rpc-url https://soroban-testnet.stellar.org \
  --network testnet --method getLatestLedger \
  --output fixtures/latest-ledger.json
stellar-replay validate --fixture fixtures/latest-ledger.json
```

This optional step is not part of the default demo or test suite.

