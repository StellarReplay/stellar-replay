# Command reference

The binary is built with:

```text
go build -o stellar-replay ./cmd/stellar-replay
```

Exit codes are stable: `0` means success, `1` means an operational failure, and
`2` means invalid command usage.

## `record`

Captures exactly one supported read-only RPC interaction. This is the only command
that may contact a live endpoint.

```text
stellar-replay record --rpc-url URL --network NAME --method METHOD --output PATH
```

Options:

- `--rpc-url`: required explicit HTTPS RPC URL;
- `--network`: required provenance label;
- `--method`: one of `getHealth`, `getLatestLedger`, `getNetwork`, or
  `getLedgerEntries`;
- `--output`: required fixture path;
- `--params JSON` or `--params-file PATH`/`-`: optional request parameters;
- `--id JSON`: scalar request ID, default `1`;
- `--tool-version VALUE`: provenance tool version;
- `--timeout DURATION`: bounded capture timeout, default `30s`.

## `validate`

Loads a fixture, rejects unknown/malformed/tampered content, and prints its
integrity digest:

```text
stellar-replay validate --fixture fixtures/get-health.json
```

## `inspect`

Loads and validates a fixture, then prints safe metadata without dumping response
payloads:

```text
stellar-replay inspect --fixture fixtures/get-health.json
```

## `replay`

Replays one request entirely offline. Fixtures can be repeated. Use flags or a
complete JSON-RPC request file/stdin:

```text
stellar-replay replay --fixture fixtures/get-health.json --method getHealth
stellar-replay replay --fixture fixtures/get-health.json --request-file request.json
```

`--params JSON`/`--params-file PATH` and `--id JSON` apply to flag-based requests.
Replay never falls back to a live endpoint.

## `serve`

Serves one or more validated fixtures at a loopback HTTP endpoint until interrupted:

```text
stellar-replay serve --fixture fixtures/get-health.json --listen 127.0.0.1:8787
```

Only `POST /` is accepted. The server uses JSON-RPC envelopes for protocol errors
and does not expose a public interface by default.
