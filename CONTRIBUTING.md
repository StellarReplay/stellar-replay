# Contributing

Stellar Replay is documentation-first and phase-driven. Read `PROJECT_CONTEXT.md`, `ROADMAP.md`, and the relevant design document before changing code.

Contributions should be narrow, tested, and honest about protocol limitations. Add deterministic local tests for behavior; keep live-network checks opt-in. Do not add transaction submission, secret handling, or new RPC methods without a roadmap decision and protocol/security documentation.

## Fresh checkout workflow

From the repository root:

```text
go test ./...
go vet ./...
go test -race ./...
go build ./cmd/stellar-replay
go test ./examples -run TestOfflineReplayWorkflow -count=1
```

No credentials or live RPC access are required. The checked-in fixtures and
offline example are the default proof path. Run `gofmt -w .` and
`git diff --check` before committing.

## Documentation and fixtures

Update the relevant documentation when behavior changes. Command flags belong in
[`docs/COMMANDS.md`](docs/COMMANDS.md); supported RPC shapes belong in
[`docs/METHODS.md`](docs/METHODS.md); operational issues belong in
[`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md). Review links and run the
documentation tests in `internal/doccheck`.

Checked-in fixtures must be synthetic or explicitly approved, sanitized, schema
valid, integrity sealed, and free of credentials or private data. Do not add a
live endpoint to default tests.

Before submitting a change, run the applicable Go formatting, test, vet, and diff checks. Update the changelog and phase checkpoint when a phase changes. Never commit secrets or captured data that is not safe to publish.

