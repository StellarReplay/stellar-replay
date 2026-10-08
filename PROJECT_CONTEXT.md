# Stellar Replay Project Context

Repository state is authoritative. Future agents must inspect the tree, `git status`, this file, and `ROADMAP.md` before acting.

## Purpose

Stellar Replay is developer infrastructure for turning controlled Stellar RPC interactions into deterministic local fixtures for tests. It is not an explorer, wallet, generic RPC client, verifier, or hosted service.

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

## Trust and security

Fixture integrity is SHA-256 over canonical fixture content excluding the integrity field. Provenance records endpoint URL, network label, capture time, and tool version; these are metadata, not proof. Secrets are rejected from explicit sensitive fields and sanitized according to documented rules. No signing or transaction submission is allowed in v0.1.

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

Phase 0, repository and project constitution: `READY_FOR_REVIEW` locally. The foundation files exist and local checks pass, but this new repository has no commit or GitHub remote yet. Phase 1 remains `NOT_STARTED`. No product core implementation is authorized by this bootstrap task.

