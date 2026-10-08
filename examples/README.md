# Examples

The executable offline workflow is covered by
[`offline_replay_test.go`](offline_replay_test.go). It loads the four checked-in
fixtures, starts the loopback replay server, verifies repeated deterministic
responses, exercises concurrent clients, and blocks non-loopback destinations.

Run it with:

```text
go test ./examples -run TestOfflineReplayWorkflow -count=1
```

Normal examples use only checked-in sanitized fixtures. Live capture is opt-in
through the CLI `record` command and is never required for tests.

