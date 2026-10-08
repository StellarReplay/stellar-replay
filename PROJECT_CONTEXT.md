# Stellar Replay Project Context

Repository state is authoritative. Future agents must inspect the tree, `git status`, this file, and `ROADMAP.md` before acting.

## Purpose

Stellar Replay is developer infrastructure for turning controlled Stellar RPC interactions into deterministic local fixtures for tests. It is not an explorer, wallet, generic RPC client, verifier, or hosted service.

## GitHub organization and repository topology

The GitHub organization is `StellarReplay`. The current three-repository topology is:

```text
StellarReplay/
├── stellar-replay         # core Go product and CLI
├── stellar-replay-action  # GitHub Actions integration
└── stellar-replay-docs    # extended documentation and guides
```

`stellar-replay` is the source of truth for core behavior, the CLI, fixture schema and semantics, capture, normalization, sanitization, hashing, replay, and the local server. `stellar-replay-action` owns only Action metadata, wrappers, examples, and Action release behavior; it consumes a released core interface and must not copy the engine. `stellar-replay-docs` owns extended tutorials, integration guides, navigation, and conceptual material; it must link to the core repository for executable behavior.

The canonical fixture format belongs to the core repository. Other repositories consume a released schema version and cannot silently fork it. Cross-repository changes require coordinated documentation, compatibility notes, and independent repository checks.

Future candidates such as SDK/integrations, conformance fixtures, or ecosystem tooling are discovery-gated. A new repository is justified only by distinct ownership, release cadence, dependency boundary, contributor workflow, or security boundary. Repository count is not a goal.

## Frozen v0.1

Go CLI; human-readable JSON fixtures; schema validation; SHA-256 integrity; capture from a configured HTTPS RPC endpoint; inspect; replay; local HTTP JSON-RPC server; deterministic request matching; sanitization; offline tests; CI and cross-platform release packaging.

Supported methods: `getHealth`, `getLatestLedger`, `getNetwork`, and `getLedgerEntries`. `simulateTransaction`, `sendTransaction`, `getEvents`, historical transaction/ledger listing, and Horizon are out of scope for v0.1. The first four provide health, identity, current ledger, and state-read examples without transaction submission.

## Principles

- read-only and offline by default;
- exact, canonical request matching;
- fixtures are evidence of a captured interaction, never proof of chain truth;
- no arbitrary fallback or live-network fallback in replay;
- small Go standard-library-first components;
- documentation and code change together;
- no fake commands or placeholder implementations presented as complete.

## Architecture decision

Go is selected over TypeScript for a single cross-platform binary, simple HTTP concurrency, standard JSON tooling, and low operational footprint. TypeScript has excellent SDK ergonomics and ecosystem reach, but would add a runtime/package distribution decision for a CLI whose core needs are protocol and file handling. The implementation may later add an optional SDK adapter, but v0.1 core does not depend on one.

## Versioning and releases

The core product owns the `stellar-replay` CLI version and releases `v0.1.0` independently. The Action uses its own semantic version and should depend only on a released core binary/interface, never an unreleased commit. Documentation is versioned by compatibility with the core release and identifies the release or branch it describes; it is not a monolithic binary release.

## Trust and security

Fixture integrity is SHA-256 over canonical fixture content excluding the integrity field. Schema v1 is implemented in `internal/fixture`; it strictly validates the four frozen read-only methods, rejects unknown fields on load, bounds fixture size and ledger keys, and atomically stores fixtures with restrictive permissions. Controlled capture is implemented in `internal/capture` with explicit HTTPS, bounded request/response bodies, no redirects, no retries, no implicit proxy, restricted destinations, default response redaction, and an explicit additional response sanitization hook. `internal/sanitize` conservatively redacts sensitive response keys and obvious PEM/bearer material while documenting detection limits. Deterministic replay is implemented in `internal/replay`; it integrity-gates fixtures, matches exact canonical method/parameters, echoes request IDs, rejects duplicate keys, returns safe deterministic errors, and supports concurrent reads without network fallback. The loopback HTTP server is implemented in `internal/server`; it accepts one bounded `POST /`, rejects batches and malformed input, maps replay errors to JSON-RPC envelopes, handles concurrent clients, and shuts down cleanly without live fallback. The CLI is implemented in `internal/cli` and `cmd/stellar-replay`; `record`, `inspect`, `validate`, `replay`, and `serve` have stable exit codes and keep live network access limited to explicit `record`. Phase 7 includes four sanitized checked-in fixtures and a loopback-only integration example proving repeated, concurrent, network-guarded replay. Phase 8 security controls and limitations are recorded in `docs/SECURITY_CHECKLIST.md`. Provenance records endpoint URL, network label, capture time, and tool version; these are metadata, not proof. No signing or transaction submission is allowed in v0.1. See `docs/FIXTURE_SCHEMA.md`, `docs/CAPTURE.md`, `docs/REPLAY.md`, and `docs/SERVER.md` for the durable contracts and limitations.

## Workflow

For each phase: inspect state; read this file and `ROADMAP.md`; verify the prior phase; implement only the active phase; run the complete phase test/check set; inspect the tree and diff; update docs/changelog; create the phase report; commit the approved work; synchronize the commit and workflow files to the configured GitHub repository; verify the GitHub commit and hosted checks; then mark the phase `COMPLETE`. If GitHub synchronization or hosted verification cannot happen because no remote/access exists, the phase remains `READY_FOR_REVIEW` (or `BLOCKED` when progress cannot continue) rather than being called complete. Stop at the phase boundary and wait for explicit authorization before executing the next phase.

Statuses: `NOT_STARTED`, `IN_PROGRESS`, `BLOCKED`, `READY_FOR_REVIEW`, `COMPLETE`.

## Source-of-truth hierarchy

1. Actual repository state
2. This file
3. `ROADMAP.md`
4. Architecture/protocol/security docs
5. Tests
6. README
7. Chat history

## GitHub completion gate

Local code and documentation are not the whole phase. A completed phase must have its approved work committed, pushed to the intended GitHub remote, and checked through the applicable GitHub Actions workflow. The phase report must record the commit SHA, remote/repository, workflow run URL or an explicit hosted-verification limitation, tests run locally, acceptance criteria, changed files, and remaining risks. Never claim that GitHub has the work until the remote commit and checks have been directly verified. Never create or push a remote/release without explicit authorization and configured access.

## Current phase

Phase 0, repository and project constitution plus initial GitHub topology: `COMPLETE`. Phase 1 research and MVP freeze: `COMPLETE`; evidence is in `docs/PHASE_1_DECISION.md` and `PHASE_1_REPORT.md`. Phase 2 fixture schema and validation: `COMPLETE`; evidence is in `docs/FIXTURE_SCHEMA.md` and `PHASE_2_REPORT.md`. Phase 3 capture: `COMPLETE`; evidence is in `docs/CAPTURE.md` and `PHASE_3_REPORT.md`. Phase 4 replay engine: `COMPLETE`; evidence is in `docs/REPLAY.md` and `PHASE_4_REPORT.md`. Phase 5 local replay server: `COMPLETE`; evidence is in `docs/SERVER.md` and `PHASE_5_REPORT.md`. Phase 6 CLI integration: `COMPLETE`; evidence is in `PHASE_6_REPORT.md`. Phase 7 deterministic test workflow: `COMPLETE`; evidence is in `PHASE_7_REPORT.md`. Phase 8 security and sanitization hardening: `COMPLETE`; evidence is in `docs/SECURITY_CHECKLIST.md` and `PHASE_8_REPORT.md`. Phase 9 documentation and developer experience: `COMPLETE`; evidence is in `PHASE_9_REPORT.md`. Phase 10 CI and cross-platform release is next and remains `NOT_STARTED`.

