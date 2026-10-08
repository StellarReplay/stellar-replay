# Stellar Replay Roadmap

This is the executable master plan. Statuses are `NOT_STARTED`, `IN_PROGRESS`, `BLOCKED`, `READY_FOR_REVIEW`, and `COMPLETE`. A phase is complete only when its acceptance criteria, local verification, phase report, approved commit, GitHub synchronization, and applicable hosted checks pass. The deadline for v0.1 is 12:00 PM WAT, October 9, 2026; scope may be reduced, never silently expanded.

## Priority fence

- **P0 / v0.1:** research, frozen method scope, fixture schema/validation/hash, capture, replay, local server, CLI, deterministic tests, sanitization, docs, examples, CI, cross-platform release, final audit.
- **P1 / time permitting:** SDK adapters, richer diffing/metadata, more methods, UX polish.
- **P2 / cannot delay v0.1:** browser/UI, full proxy, record-all traffic, visualization, AI, unrelated integrations.

## Repository architecture

The initial open-source topology is intentionally three repositories under the `StellarReplay` organization:

```text
StellarReplay/
├── stellar-replay
├── stellar-replay-action
└── stellar-replay-docs
```

- `stellar-replay` owns the Go core, CLI, canonical fixture schema, replay semantics, local server, core tests, release binaries, and the minimum documentation required to use the product.
- `stellar-replay-action` owns GitHub Action metadata, Action-specific integration/wrapper code, examples, tests, and its own release lifecycle. It consumes a released core interface and never duplicates the Go engine.
- `stellar-replay-docs` owns extended guides, tutorials, integration documentation, conceptual explanations, and documentation navigation. It links to the core README and does not redefine protocol or fixture behavior.

The fixture schema and replay semantics remain centralized in the core repository. Repository splits are not a substitute for packages or services. A future repository is created only when distinct ownership, release cadence, dependency/security boundary, or contributor workflow creates measurable value. Candidates include SDK/integrations, conformance fixtures, and ecosystem tooling; none are committed now.

Core releases, Action releases, and docs versions are independent but compatibility-linked. A core schema or CLI change requires a versioned compatibility note; the Action must target a released core version; docs identify the core release or branch they describe. Cross-repository links and claims are checked before each phase is closed.

## Phase register

| Phase | Name | Status | Gate |
|---:|---|---|---|
| 0 | Constitution and repository foundation | COMPLETE | three repositories synchronized; core, Action, and docs foundation workflows pass |
| 1 | Research and MVP freeze | COMPLETE | decision record, protocol/network policy, local checks, and hosted verification |
| 2 | Fixture schema and validation | COMPLETE | versioned schema, deterministic integrity, strict validation, atomic storage, local and hosted verification |
| 3 | Capture | COMPLETE | bounded controlled capture, local transport tests, fixture validation, local and hosted verification |
| 4 | Replay engine | COMPLETE | exact canonical matching, integrity-gated loading, deterministic errors, concurrent-read verification, local and hosted verification |
| 5 | Local replay server | NOT_STARTED | HTTP/concurrency/shutdown tests |
| 6 | CLI integration | NOT_STARTED | end-to-end local workflow |
| 7 | Deterministic test workflow | NOT_STARTED | network-free example test |
| 8 | Security and sanitization hardening | NOT_STARTED | threat cases and review |
| 9 | Documentation and developer experience | NOT_STARTED | truthful walkthrough |
| 10 | CI and cross-platform release | NOT_STARTED | matrix and artifacts |
| 11 | Release candidate and final audit | NOT_STARTED | audit report and go/no-go |
| 12 | v0.1.0 release | NOT_STARTED | published, verified release |

## Phase execution protocol

For every phase: inspect repository state, `git status`, remotes, and the current GitHub synchronization state; read `PROJECT_CONTEXT.md` and this file; verify the previous phase's local and GitHub gates; set the phase `IN_PROGRESS`; implement only the authorized phase; run the complete phase test/check set; inspect the tree, diff, generated artifacts, and documentation claims; update docs and `CHANGELOG.md`; write a phase report; obtain required owner approval; commit the approved work; synchronize it to the configured GitHub repository; verify the GitHub commit and applicable Actions checks; mark `COMPLETE` only when all gates pass; otherwise mark `READY_FOR_REVIEW` or `BLOCKED` with the missing evidence; then stop. Do not push or create a remote without explicit authorization. If architecture changes, explain the conflict, update docs and dependencies, then implement.

## Phase contracts

Each phase below defines objective, dependencies, scope, non-scope, research, tasks, tests, docs, acceptance, commands, artifacts, risks, recovery, network policy, and review need. “Previous phase” is the required gate.

### Phase 0 — Constitution and repository foundation

**Objective/why:** make the project understandable and implementation-ready without fake product code. **Dependencies:** none. **Status:** COMPLETE after local and hosted verification.

**Scope:** identity, Go choice, context, README, architecture, protocol, security, competitive research, demo plan, OSS files, Go module, purposeful source/test directories, templates, and this roadmap. **Non-scope:** replay engine, commands, server, sample claims, live capture.

**Research:** official Stellar RPC method catalog, SDKs, RPC repository, JSON-RPC conventions, general VCR/record-replay, and Soroban fork/testing approaches. **Tasks:** freeze four read-only methods; define trust model and matching; record assumptions and citations; add contribution and release policy. **Tests/checks:** tree inspection, required-file check, `go mod tidy` only when code exists, `git diff --check`; no live network required. **Docs:** all foundation docs and changelog entry.

**Acceptance:** repository identity is Stellar Replay; Go rationale and MVP are explicit; unsupported methods and non-goals are explicit; architecture/protocol/security/competitive/demo docs exist; no document claims unimplemented functionality; GitHub foundation exists; repository ownership and topology are documented; phase workflow is executable; an approved commit is synchronized to GitHub and the foundation workflow is verified, or the phase remains `READY_FOR_REVIEW` until that is possible. **Commands:** `git status --short --branch`; `git remote -v`; `rg --files`; `git diff --check`; inspect required headings; after synchronization, inspect the GitHub commit and Actions run. **Artifacts:** foundation files, directory tree, topology decision, phase report, commit, and hosted evidence. **Risks/recovery:** accidental overlap with another repository; recover by preserving the unrelated repo and using a dedicated directory. **Network:** research only for local work; GitHub access is required for the completion gate. **Review:** independent review not required yet; architecture review may be requested after Phase 1.

### Phase 1 — Research and MVP freeze

**Objective/why:** validate assumptions before code. **Dependencies:** Phase 0. **Scope:** confirm current RPC semantics, JSON-RPC error/id behavior, SDK endpoint configurability, and method-specific parameter considerations; decide redirect/private endpoint policy. **Non-scope:** implementation or broad ecosystem survey. **Status:** COMPLETE; evidence is in `docs/PHASE_1_DECISION.md` and `PHASE_1_REPORT.md`.

**Tasks:** update landscape/protocol; write `docs/PHASE_1_DECISION.md`; freeze fixture fields and exact method list; define zero-parameter normalization, id handling, error mapping, endpoint restrictions, redirect behavior, proxy behavior, and local-test transport policy. **Tests:** document examples and negative cases; no live network required except opt-in research. **Acceptance:** each supported method has request shape, response policy, matching policy, and limitation; network policy is approved in docs; the official JSON-RPC and SDK sources are recorded with a research date; no implementation starts. **Commands:** `rg -n "supported|unsupported|redirect|private|fixture|proxy|params|request id" docs README.md`; `git diff --check`; `go test ./...`; `go vet ./...`. **Artifacts:** decision record and revised docs. **Risks:** protocol drift; record source URLs/date and preserve raw responses. **Recovery:** amend docs before schema work. **Review:** architecture review is justified if method/transport policy changes; no independent security review is required before implementation because the policy is deliberately restrictive.

### Phase 2 — Fixture schema and validation

**Objective/why:** establish the durable interchange contract. **Dependencies:** Phase 1. **Scope:** versioned JSON schema/model, canonical JSON, provenance, request/response/error, hash, validation, atomic storage. **Non-scope:** network capture and server.

**Tasks:** implement `internal/fixture`; add valid, malformed, missing-field, modified-hash, and secret-rejection cases. **Tests:** deterministic serialization/hash, schema validation, repeated load. **Acceptance:** identical logical fixtures hash identically; invalid or modified fixtures fail before replay; fields and limitations match docs. **Commands:** `gofmt -w .`; `go test ./...`; `go vet ./...`; `git diff --check`. **Artifacts:** package, tests, schema examples. **Risks/recovery:** schema churn; version it and record migration decision. **Network:** none. **Review:** not required unless canonicalization/security is contentious. **Status:** COMPLETE; implementation and limitations are documented in `docs/FIXTURE_SCHEMA.md`, with evidence in `PHASE_2_REPORT.md`.

### Phase 3 — Capture

**Objective/why:** obtain truthful fixtures from controlled read-only calls. **Dependencies:** Phase 2. **Scope:** bounded HTTPS client, supported methods, timeout, response capture, sanitization hook, atomic write. **Non-scope:** proxying all traffic, secrets, writes, retries that alter determinism.

**Tasks:** implement `internal/capture`; reject unsupported methods; preserve JSON-RPC errors; make endpoint explicit. **Tests:** local HTTP/TLS test server, timeout, malformed response, size bound, HTTPS policy, deterministic fixture output, redirect rejection, no-retry behavior, and sanitization hook. **Acceptance:** only the capture operation contacts the configured endpoint; captured data validates and hashes; request secrets are rejected and response sanitization is explicit. **Commands:** `go test ./internal/capture ./internal/fixture`; `go vet ./...`; `git diff --check`. **Artifacts:** capture package, tests, and `docs/CAPTURE.md`. **Risks/recovery:** endpoint behavior changes; keep raw response and limitation metadata. **Network:** opt-in live endpoint; tests local. **Review:** security review remains justified before release. **Status:** COMPLETE; evidence is in `PHASE_3_REPORT.md`.

### Phase 4 — Replay engine

**Objective/why:** make fixture behavior deterministic and explicit. **Dependencies:** Phase 2. **Scope:** canonical request matching, repeated calls, error responses, useful mismatch diagnostics, concurrent-safe reads. **Non-scope:** fuzzy matching, live fallback, response mutation.

**Tests:** exact match, key-order normalization, type/array mismatch, method mismatch, unknown method, missing fixture, repeated and concurrent calls. **Acceptance:** no unmatched request returns a successful response; same request always returns same stored response; fixture integrity is checked on load. **Commands:** `go test ./internal/replay ./internal/fixture`; `go test -race ./...`; `go vet ./...`. **Artifacts:** replay package/tests and `docs/REPLAY.md`. **Risks/recovery:** matching ambiguity; stop and update protocol before changing semantics. **Network:** none. **Review:** independent review optional; required if matching becomes fuzzy. **Status:** COMPLETE; evidence is in `PHASE_4_REPORT.md`.

### Phase 5 — Local replay server

**Objective/why:** let existing clients point at a local endpoint. **Dependencies:** Phases 2 and 4. **Scope:** HTTP JSON-RPC transport, malformed JSON, errors, concurrent requests, clean shutdown, loopback defaults. **Non-scope:** live fallback, batch unless specified, TLS termination.

**Tests:** startup, shutdown, invalid JSON, unsupported method, mismatch, concurrent requests, no-network assertion. **Acceptance:** server rejects invalid/unmatched input clearly and serves only validated fixtures; shutdown releases the listener. **Commands:** `go test ./internal/server ./internal/replay`; `go test -race ./...`; `go vet ./...`. **Artifacts:** server package and integration tests. **Risks/recovery:** HTTP behavior differs from SDK expectations; add a fixture-backed compatibility test and document gaps. **Network:** local only. **Review:** concurrency review justified.

### Phase 6 — CLI integration

**Objective/why:** expose a coherent developer workflow. **Dependencies:** Phases 3–5. **Scope:** `record`, `inspect`, `validate`, `replay`, `serve`; stable exit codes; useful errors; stdin/file options where documented. **Non-scope:** extra commands, shell completion, public SDK.

**Tests:** command parsing, end-to-end fixture workflow, clean exit codes, output paths, no-network replay. **Acceptance:** README walkthrough uses real commands/output; no placeholder command remains. **Commands:** `go build ./cmd/stellar-replay`; `go test ./...`; `go vet ./...`; `git diff --check`. **Artifacts:** binary entry point and CLI tests. **Risks/recovery:** CLI surface churn; preserve backwards-compatible flags after v0.1 freeze or document break. **Network:** record opt-in only. **Review:** UX review optional.

### Phase 7 — Deterministic test workflow

**Objective/why:** prove the product value in a normal test suite. **Dependencies:** Phase 6. **Scope:** checked-in safe fixtures, test helper or documented HTTP usage, offline integration example. **Non-scope:** live network in default tests, SDK-specific packages.

**Tests:** repeated offline run, network-disabled test, fixture matrix, concurrent clients. **Acceptance:** a new contributor can run tests without credentials/network and see the replay path. **Commands:** `go test ./...`; platform-equivalent test command; optionally block network in test harness. **Artifacts:** `fixtures/`, `examples/`, `tests/`. **Risks/recovery:** fixtures accidentally encode secrets; sanitize and rotate/remove them. **Network:** none for default checks. **Review:** not required.

### Phase 8 — Security and sanitization hardening

**Objective/why:** validate trust boundaries before release. **Dependencies:** Phases 3–7. **Scope:** secret detection/redaction, bounded input, URL policy, logs, file permissions, path handling, fixture sharing guidance. **Non-scope:** cryptographic chain proofs or transaction security.

**Tests:** secret-like fields, oversized body, malformed fixture, path traversal, private endpoint policy, logs, concurrent access. **Acceptance:** threat model matches code; no secret is written in tests; limitations are explicit. **Commands:** `go test ./...`; `go vet ./...`; selected static/security tools documented. **Artifacts:** updated docs, adversarial tests, security checklist. **Risks/recovery:** over-redaction breaks replay; fail closed and document explicit opt-out only if safe. **Network:** local adversarial tests; no live network required. **Review:** independent security review strongly justified.

### Phase 9 — Documentation and developer experience

**Objective/why:** make v0.1 understandable in 2–3 minutes. **Dependencies:** Phases 6–8. **Scope:** README walkthrough, demo, method docs, troubleshooting, fixture examples, contribution path. **Non-scope:** claims of adoption or ecosystem support not evidenced.

**Tests:** fresh checkout walkthrough, link/path check, command accuracy, docs drift review. **Acceptance:** every command is implemented or labeled planned; limitations and security model are visible. **Commands:** `go test ./...`; `rg -n "planned|implemented|experimental" README.md docs`. **Artifacts:** polished docs and changelog. **Risks/recovery:** stale examples; run them at each release checkpoint. **Network:** demo capture may be opt-in; validation/replay offline. **Review:** external reader review useful but not required.

### Phase 10 — CI and cross-platform release

**Objective/why:** produce credible binaries. **Dependencies:** Phases 6–9. **Scope:** GitHub Actions on Windows/Linux/macOS, Go formatting/test/vet, race on supported runner, tagged builds, archives, checksums, permissions. **Non-scope:** signing or registry publishing unless separately approved.

**Tests:** hosted matrix, archive smoke tests, checksum verification, clean-tree build. **Acceptance:** declared platforms pass; release artifacts contain only intended files; docs explain installation. **Commands:** `go test ./...`; `go vet ./...`; `go build ./cmd/stellar-replay`; workflow run IDs recorded. **Artifacts:** `.github/workflows`, release scripts, docs. **Risks/recovery:** runner drift; pin actions and record tool versions. **Network:** CI downloads dependencies; product tests offline. **Review:** release/security review justified.

### Phase 11 — Release candidate and final audit

**Objective/why:** prevent documentation, security, and behavior drift. **Dependencies:** Phase 10. **Scope:** create `FINAL_AUDIT.md`; verify methods, fixtures, matching, malformed input, concurrency, shutdown, secrets, sanitization, docs, CI, binaries, dead code/placeholders. **Non-scope:** new features unless audit blocks release.

**Tests:** full suite, race/static checks, release smoke tests, final network-boundary check. **Acceptance:** every finding is resolved, accepted with owner/evidence, or blocks release; independent Codex review if available; go/no-go is explicit. **Commands:** all project checks plus `git diff --check` and clean-tree inspection. **Artifacts:** `FINAL_AUDIT.md`, release checklist, updated changelog/context. **Risks/recovery:** audit discovers scope defect; defer release and update roadmap, never hide it. **Network:** hosted validation may be needed; no live network for default tests. **Review:** required independent review if environment supports it.

### Phase 12 — v0.1.0 release

**Objective/why:** publish only a verified, truthful release. **Dependencies:** Phase 11 go decision. **Scope:** tag, GitHub release, binaries/checksums, installation docs, post-release smoke test, stabilization note. **Non-scope:** unplanned features or adoption claims.

**Acceptance:** published artifacts match checksums; Windows/Linux/macOS smoke tests pass; changelog and known limitations are correct; release URL and evidence are recorded. **Commands:** clean-tree check, archive extraction, binary `--help`/documented smoke, checksum verification. **Artifacts:** release, notes, audit evidence. **Risks/recovery:** bad artifact; withdraw/replace with a documented patch release, preserve audit trail. **Network:** release hosting required; app tests remain offline. **Review:** maintainer go/no-go.

## Future direction

v0.x may add additional read-only methods, fixture diffing, SDK adapters, and richer sanitization only from real usage. v1.0 may consider a stable public Go API, broader RPC coverage, stronger provenance, and a documented network trust model. v2 may consider stateful forks, policy integrations, and broader ecosystem tooling. Browser UI, transaction writes, AI decisions, and generic blockchain mocking remain discovery-gated and cannot delay v0.1. Future contributors can work on method adapters, fixtures, docs, SDK compatibility, performance, sanitization, replay semantics, platform support, and test tooling.

