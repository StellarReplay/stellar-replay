# Changelog

## v0.1.0 - 2026-10-08

- First verified Stellar Replay core release. The release includes the read-only
  capture, fixture, replay, local server, CLI, security, documentation, CI, and
  cross-platform packaging foundations described in this changelog.

- Established the Stellar Replay project constitution, research record, frozen v0.1 scope, architecture, protocol, security model, and executable roadmap.
- Added the mandatory phase-completion gate requiring an approved commit, GitHub synchronization, hosted-check verification, and a phase report before a phase can be marked complete.
- Completed Phase 1 research and MVP freeze with current Stellar RPC request, response, matching, SDK, and capture-endpoint decisions.
- Implemented Phase 2 fixture schema v1, canonical JSON hashing, strict validation, sensitive-field rejection, bounded loading, and atomic storage with tests.
- Implemented Phase 3 controlled capture with bounded HTTPS transport, endpoint restrictions, timeout/size enforcement, no redirects or retries, response sanitization hook, atomic recording, and local transport tests.
- Implemented Phase 4 deterministic replay with integrity-gated fixture loading, canonical exact matching, request-ID echoing, safe errors, duplicate-key rejection, and concurrent-read tests.
- Implemented Phase 5 loopback local replay server with bounded HTTP JSON-RPC transport, malformed/batch/error handling, concurrency tests, and clean shutdown.
- Implemented Phase 6 CLI commands for record, inspect, validate, replay, and serve with stable exit codes, file inputs, offline workflow tests, and executable documentation.
- Added Phase 7 sanitized fixture matrix and a loopback-only offline integration example covering repeated and concurrent deterministic replay.
- Hardened Phase 8 secret handling, request/response bounds, endpoint query policy, fixture permissions, and adversarial security coverage with a documented checklist.
- Completed Phase 9 documentation and developer experience with an executable offline demo, command/method/troubleshooting guides, contributor workflow, and automated documentation drift checks.
- Completed Phase 10 CI and cross-platform release foundation with Windows/Linux/macOS verification, reproducible tagged release packaging, archive smoke tests, embedded release versions, checksum generation, least-privilege workflow permissions, and release installation guidance.
- Completed Phase 11 final audit with resolved documentation/CI drift, explicit accepted limitations, release conditions, and a conditional GO decision for Phase 12.
- Published v0.1.0 with Linux amd64, Windows amd64, and macOS amd64 archives plus independently verified SHA-256 checksums and post-release smoke tests.

