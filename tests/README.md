# Tests

The default test suite is deterministic and network-free. The integration example
in `../examples/offline_replay_test.go` is the contributor-facing proof of the
fixture-to-replay workflow and uses a transport guard that rejects non-loopback
destinations.

Run all checks with:

```text
go test ./...
```

Live RPC checks, if ever added, must be explicit, separate, and never required for
CI success.

