# Stellar Replay Roadmap

This is the executable master plan. Statuses are `NOT_STARTED`, `IN_PROGRESS`, `BLOCKED`, `READY_FOR_REVIEW`, and `COMPLETE`. A phase is complete only when its acceptance criteria, local verification, phase report, approved commit, GitHub synchronization, and applicable hosted checks pass. v0.1.0 is released; future work is productization-gated and must not silently expand the frozen core contract.

## Priority fence

- **P0 / v0.1:** research, frozen method scope, fixture schema/validation/hash, capture, replay, local server, CLI, deterministic tests, sanitization, docs, examples, CI, cross-platform release, final audit.
- **P1 / time permitting:** SDK adapters, richer diffing/metadata, more methods, UX polish.
- **P2 / cannot delay v0.1:** browser/UI, full proxy, record-all traffic, visualization, AI, unrelated integrations.

## Repository architecture

The current open-source topology is intentionally three repositories under the `StellarReplay` organization:

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

The full ownership map is maintained in [docs/REPOSITORY_ARCHITECTURE.md](docs/REPOSITORY_ARCHITECTURE.md). The organization audit and lifecycle decisions are recorded in [ORGANIZATION_AUDIT.md](ORGANIZATION_AUDIT.md) and [REPOSITORY_ACTION_PLAN.md](REPOSITORY_ACTION_PLAN.md). Historical phase evidence is indexed by [docs/PHASE_HISTORY.md](docs/PHASE_HISTORY.md); reports are evidence, not a second source of product truth.

## Strategy hierarchy

The long-term plan is interpreted as **product era → milestone → phase → task → verification evidence**. Completed phases 0–12 form Era 1 Foundation/MVP/Release. The organization is now in Era 2, Developer Productization: the core is released, while the Action and docs repositories are pre-release companions. Future phase cards must include objective, dependencies, affected repositories, scope and non-scope, tests, security, documentation, CI, acceptance, commands, evidence, risks, rollback, release impact, agent/review guidance, network needs, and maintainer authorization.

## Phase register

| Phase | Name | Status | Gate |
|---:|---|---|---|
| 0 | Constitution and repository foundation | COMPLETE | three repositories synchronized; core, Action, and docs foundation workflows pass |
| 1 | Research and MVP freeze | COMPLETE | decision record, protocol/network policy, local checks, and hosted verification |
| 2 | Fixture schema and validation | COMPLETE | versioned schema, deterministic integrity, strict validation, atomic storage, local and hosted verification |
| 3 | Capture | COMPLETE | bounded controlled capture, local transport tests, fixture validation, local and hosted verification |
| 4 | Replay engine | COMPLETE | exact canonical matching, integrity-gated loading, deterministic errors, concurrent-read verification, local and hosted verification |
| 5 | Local replay server | COMPLETE | loopback HTTP JSON-RPC transport, malformed-input/error handling, concurrency, clean shutdown, local and hosted verification |
| 6 | CLI integration | COMPLETE | five-command workflow, stable exit codes, end-to-end offline tests, local and hosted verification |
| 7 | Deterministic test workflow | COMPLETE | checked-in sanitized fixture matrix, loopback-only offline example, repeated/concurrent proof, local and hosted verification |
| 8 | Security and sanitization hardening | COMPLETE | adversarial controls, default redaction, security checklist, local and hosted verification |
| 9 | Documentation and developer experience | COMPLETE | executable walkthrough, command/method/troubleshooting docs, contributor path, automated link and drift checks, local and hosted verification |
| 10 | CI and cross-platform release | COMPLETE | matrix and artifacts, local package smoke tests, hosted verification |
| 11 | Release candidate and final audit | COMPLETE | final audit, findings, go/no-go, local and hosted verification |
| 12 | v0.1.0 release | COMPLETE | published artifacts, checksums, platform smoke tests, release evidence |

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

**Objective/why:** validate assumptions before code. **Dependencies:** Phase 0. **Scope:** confirm current RPC semantics, JSON-RPC error/id behavior, SDK endpoint configurability, and method-specific parameter considerations; decide redirect/private endpoint policy. **Non-scope:** implementation or broad ecosystem survey. **Status:** COMPLETE; evidence is in `docs/PHASE_1_DECISION.md` and `docs/history/phases/PHASE-001.md`.

**Tasks:** update landscape/protocol; write `docs/PHASE_1_DECISION.md`; freeze fixture fields and exact method list; define zero-parameter normalization, id handling, error mapping, endpoint restrictions, redirect behavior, proxy behavior, and local-test transport policy. **Tests:** document examples and negative cases; no live network required except opt-in research. **Acceptance:** each supported method has request shape, response policy, matching policy, and limitation; network policy is approved in docs; the official JSON-RPC and SDK sources are recorded with a research date; no implementation starts. **Commands:** `rg -n "supported|unsupported|redirect|private|fixture|proxy|params|request id" docs README.md`; `git diff --check`; `go test ./...`; `go vet ./...`. **Artifacts:** decision record and revised docs. **Risks:** protocol drift; record source URLs/date and preserve raw responses. **Recovery:** amend docs before schema work. **Review:** architecture review is justified if method/transport policy changes; no independent security review is required before implementation because the policy is deliberately restrictive.

### Phase 2 — Fixture schema and validation

**Objective/why:** establish the durable interchange contract. **Dependencies:** Phase 1. **Scope:** versioned JSON schema/model, canonical JSON, provenance, request/response/error, hash, validation, atomic storage. **Non-scope:** network capture and server.

**Tasks:** implement `internal/fixture`; add valid, malformed, missing-field, modified-hash, and secret-rejection cases. **Tests:** deterministic serialization/hash, schema validation, repeated load. **Acceptance:** identical logical fixtures hash identically; invalid or modified fixtures fail before replay; fields and limitations match docs. **Commands:** `gofmt -w .`; `go test ./...`; `go vet ./...`; `git diff --check`. **Artifacts:** package, tests, schema examples. **Risks/recovery:** schema churn; version it and record migration decision. **Network:** none. **Review:** not required unless canonicalization/security is contentious. **Status:** COMPLETE; implementation and limitations are documented in `docs/FIXTURE_SCHEMA.md`, with evidence in `docs/history/phases/PHASE-002.md`.

### Phase 3 — Capture

**Objective/why:** obtain truthful fixtures from controlled read-only calls. **Dependencies:** Phase 2. **Scope:** bounded HTTPS client, supported methods, timeout, response capture, sanitization hook, atomic write. **Non-scope:** proxying all traffic, secrets, writes, retries that alter determinism.

**Tasks:** implement `internal/capture`; reject unsupported methods; preserve JSON-RPC errors; make endpoint explicit. **Tests:** local HTTP/TLS test server, timeout, malformed response, size bound, HTTPS policy, deterministic fixture output, redirect rejection, no-retry behavior, and sanitization hook. **Acceptance:** only the capture operation contacts the configured endpoint; captured data validates and hashes; request secrets are rejected and response sanitization is explicit. **Commands:** `go test ./internal/capture ./internal/fixture`; `go vet ./...`; `git diff --check`. **Artifacts:** capture package, tests, and `docs/CAPTURE.md`. **Risks/recovery:** endpoint behavior changes; keep raw response and limitation metadata. **Network:** opt-in live endpoint; tests local. **Review:** security review remains justified before release. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-003.md`.

### Phase 4 — Replay engine

**Objective/why:** make fixture behavior deterministic and explicit. **Dependencies:** Phase 2. **Scope:** canonical request matching, repeated calls, error responses, useful mismatch diagnostics, concurrent-safe reads. **Non-scope:** fuzzy matching, live fallback, response mutation.

**Tests:** exact match, key-order normalization, type/array mismatch, method mismatch, unknown method, missing fixture, repeated and concurrent calls. **Acceptance:** no unmatched request returns a successful response; same request always returns same stored response; fixture integrity is checked on load. **Commands:** `go test ./internal/replay ./internal/fixture`; `go test -race ./...`; `go vet ./...`. **Artifacts:** replay package/tests and `docs/REPLAY.md`. **Risks/recovery:** matching ambiguity; stop and update protocol before changing semantics. **Network:** none. **Review:** independent review optional; required if matching becomes fuzzy. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-004.md`.

### Phase 5 — Local replay server

**Objective/why:** let existing clients point at a local endpoint. **Dependencies:** Phases 2 and 4. **Scope:** HTTP JSON-RPC transport, malformed JSON, errors, concurrent requests, clean shutdown, loopback defaults. **Non-scope:** live fallback, batch unless specified, TLS termination.

**Tests:** startup, shutdown, invalid JSON, unsupported method, mismatch, malformed/batch input, transport misuse, bounded body, concurrent requests, loopback configuration, and no-network assertion. **Acceptance:** server rejects invalid/unmatched input clearly and serves only validated fixtures; shutdown releases the listener. **Commands:** `go test ./internal/server ./internal/replay`; `go test -race ./...`; `go vet ./...`. **Artifacts:** server package, tests, and `docs/SERVER.md`. **Risks/recovery:** HTTP behavior differs from SDK expectations; fixture-backed compatibility tests and documented gaps are required. **Network:** local only. **Review:** concurrency review justified. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-005.md`.

### Phase 6 — CLI integration

**Objective/why:** expose a coherent developer workflow. **Dependencies:** Phases 3–5. **Scope:** `record`, `inspect`, `validate`, `replay`, `serve`; stable exit codes; useful errors; stdin/file options where documented. **Non-scope:** extra commands, shell completion, public SDK.

**Tests:** command parsing, end-to-end fixture workflow, clean exit codes, output paths, request-file replay, binary help, and no-network replay. **Acceptance:** README walkthrough uses real commands/output; no placeholder command remains. **Commands:** `go build ./cmd/stellar-replay`; `go test ./...`; `go vet ./...`; `git diff --check`. **Artifacts:** binary entry point, CLI package/tests, and executable workflow documentation. **Risks/recovery:** CLI surface churn; preserve backwards-compatible flags after v0.1 freeze or document break. **Network:** record opt-in only. **Review:** UX review optional. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-006.md`.

### Phase 7 — Deterministic test workflow

**Objective/why:** prove the product value in a normal test suite. **Dependencies:** Phase 6. **Scope:** checked-in safe fixtures, test helper or documented HTTP usage, offline integration example. **Non-scope:** live network in default tests, SDK-specific packages.

**Tests:** repeated offline run, network-disabled transport guard, fixture matrix, concurrent clients. **Acceptance:** a new contributor can run tests without credentials/network and see the replay path. **Commands:** `go test ./...`; platform-equivalent test command; optionally block network in test harness. **Artifacts:** sanitized `fixtures/`, `examples/` integration test, and `tests/` guidance. **Risks/recovery:** fixtures accidentally encode secrets; sanitize and rotate/remove them. **Network:** none for default checks. **Review:** not required. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-007.md`.

### Phase 8 — Security and sanitization hardening

**Objective/why:** validate trust boundaries before release. **Dependencies:** Phases 3–7. **Scope:** secret detection/redaction, bounded input, URL policy, logs, file permissions, path handling, fixture sharing guidance. **Non-scope:** cryptographic chain proofs or transaction security.

**Tests:** secret-like fields and values, oversized request/response/fixture bodies, malformed fixture/response, path/output handling, private endpoint policy, diagnostic/log safety, concurrent access, and no-network replay. **Acceptance:** threat model matches code; no secret is written in tests; limitations are explicit. **Commands:** `go test ./...`; `go vet ./...`; selected static/security tools documented. **Artifacts:** updated docs, adversarial tests, default redaction, and security checklist. **Risks/recovery:** over-redaction breaks replay; fail closed and document explicit opt-out only if safe. **Network:** local adversarial tests; no live network required. **Review:** independent security review strongly justified. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-008.md`.

### Phase 9 — Documentation and developer experience

**Objective/why:** make v0.1 understandable in 2–3 minutes. **Dependencies:** Phases 6–8. **Scope:** README walkthrough, demo, method docs, troubleshooting, fixture examples, contribution path. **Non-scope:** claims of adoption or ecosystem support not evidenced.

**Tests:** fresh checkout walkthrough, Markdown link/path check, command accuracy, docs drift review, offline demo, and full Go checks. **Acceptance:** every command is implemented or labeled planned; limitations and security model are visible. **Commands:** `go test ./...`; `go test ./internal/doccheck`; `rg -n "planned|implemented|experimental" README.md docs`. **Artifacts:** polished README/demo, command/method/troubleshooting guides, contributor path, documentation checks, and changelog. **Risks/recovery:** stale examples; run them at each release checkpoint. **Network:** demo capture may be opt-in; validation/replay offline. **Review:** external reader review useful but not required. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-009.md`.

### Phase 10 — CI and cross-platform release

**Objective/why:** produce credible binaries. **Dependencies:** Phases 6–9. **Scope:** GitHub Actions on Windows/Linux/macOS, Go formatting/test/vet, race on supported runner, tagged builds, archives, checksums, permissions. **Non-scope:** signing or registry publishing unless separately approved. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-010.md`.

**Tests:** hosted matrix, archive smoke tests, checksum verification, clean-tree build. **Acceptance:** declared platforms pass; release artifacts contain only intended files; docs explain installation. **Commands:** `go test ./...`; `go vet ./...`; `go build ./cmd/stellar-replay`; workflow run IDs recorded. **Artifacts:** `.github/workflows`, release scripts, docs. **Risks/recovery:** runner drift; pin actions and record tool versions. **Network:** CI downloads dependencies; product tests offline. **Review:** release/security review justified.

### Phase 11 — Release candidate and final audit

**Objective/why:** prevent documentation, security, and behavior drift. **Dependencies:** Phase 10. **Scope:** create `docs/history/audits/FINAL_AUDIT.md`; verify methods, fixtures, matching, malformed input, concurrency, shutdown, secrets, sanitization, docs, CI, binaries, dead code/placeholders. **Non-scope:** new features unless audit blocks release. **Status:** COMPLETE; evidence and the conditional Phase 12 go decision are in `docs/history/audits/FINAL_AUDIT.md`; later release evidence is in `docs/history/phases/PHASE-012.md`.

**Tests:** full suite, race/static checks, release smoke tests, final network-boundary check. **Acceptance:** every finding is resolved, accepted with owner/evidence, or blocks release; independent Codex review if available; go/no-go is explicit. **Commands:** all project checks plus `git diff --check` and clean-tree inspection. **Artifacts:** `docs/history/audits/FINAL_AUDIT.md`, release checklist, updated changelog/context. **Risks/recovery:** audit discovers scope defect; defer release and update roadmap, never hide it. **Network:** hosted validation may be needed; no live network for default tests. **Review:** required independent review if environment supports it.

### Phase 12 — v0.1.0 release

**Objective/why:** publish only a verified, truthful release. **Dependencies:** Phase 11 go decision. **Scope:** tag, GitHub release, binaries/checksums, installation docs, post-release smoke test, stabilization note. **Non-scope:** unplanned features or adoption claims. **Status:** COMPLETE; evidence is in `docs/history/phases/PHASE-012.md`.

**Acceptance:** published artifacts match checksums; Windows/Linux/macOS smoke tests pass; changelog and known limitations are correct; release URL and evidence are recorded. **Commands:** clean-tree check, archive extraction, binary `--help`/documented smoke, checksum verification. **Artifacts:** release, notes, audit evidence. **Risks/recovery:** bad artifact; withdraw/replace with a documented patch release, preserve audit trail. **Network:** release hosting required; app tests remain offline. **Review:** maintainer go/no-go.

## Long-term product roadmap

The roadmap is hierarchical rather than an endless list of numbered phases:

```text
Product vision
  Era 1 — Foundation, MVP, and release (complete)
    Milestones 1–9 — constitution through v0.1.0 (Phases 0–12)
  Era 2 — Developer productization (current)
    Milestone 10 — post-release stabilization and adoption readiness
    Milestone 11 — GitHub Action integration
    Milestone 12 — documentation and onboarding depth
  Era 3 — Ecosystem (candidate)
    Milestone 13 — selected integrations and conformance assets
  Era 4 — Trust and scale (candidate)
    Milestone 14 — provenance, compatibility, and operational maturity
  Era 5 — Mature product (speculative)
    Milestone 15 — only evidence-backed expansion
```

The current gap is not missing foundation code. It is the path from a verified
MVP to reliable adoption: post-release feedback, installation/onboarding
quality, a secure Action contract, useful extended docs, support practice, and
evidence that additional methods or repositories are needed.

### Future phase cards

Every card below is a planning contract. A future session must first **Check
Phase N** against actual repository state and then execute only after explicit
authorization. “Codex” is the primary agent for ordinary work; an independent
security or architecture review is used only when the card says the risk or
scope justifies it.

#### Phase 13 — Post-release stabilization and adoption readiness

- **Parent:** Era 2 / Milestone 10. **Status:** PLANNED. **Objective:** turn the
  verified v0.1.0 release into a maintainable first-user experience.
- **Why/dependencies:** follows Phase 12 and depends on the published artifacts;
  no core contract expansion.
- **Repositories:** core first; cross-repository link checks for Action/docs.
- **In scope:** installation walkthrough verification, release/support issue
  templates, compatibility notes, stale-claim scan, reproducible smoke script,
  and a small feedback log. **Out of scope:** new RPC methods, Action runtime,
  UI, signing, or hosted service.
- **Tests/security/docs/CI:** full Go suite, race, vet, doccheck, archive smoke;
  review fixture/secrets guidance; update README, RELEASES, CHANGELOG, and a
  phase report; retain least-privilege workflows.
- **Acceptance/evidence:** clean fresh-checkout walkthrough, release URL and
  checksums still valid, no stale “unreleased” claims, hosted core CI passes.
- **Commands:** `go test -p 1 -count=1 ./...`; `go test -race -p 1 -count=1 ./...`;
  `go vet ./...`; `go test ./internal/doccheck`; `git diff --check`.
- **Risks/rollback/release:** documentation-only or patch-level rollback; no
  release unless a separately authorized fix is required. Live network: no.
  Maintainer authorization: yes. Recommended review: Codex; independent UX
  review optional.

#### Phase 14 — Action contract and secure integration

- **Parent:** Era 2 / Milestone 11. **Status:** CANDIDATE. **Objective:** define
  and, if justified, implement the smallest useful Action around released core.
- **Why/dependencies:** depends on Phase 13 feedback and an approved binary
  acquisition/version-pinning design. **Repositories:** Action plus core docs;
  core only if a documented stable CLI interface must be clarified.
- **In scope:** input contract, pinned core version, checksum/provenance policy,
  exit/log behavior, offline fixture tests, least-privilege workflow, and Action
  release notes. **Out of scope:** copied Go code, live credentials, implicit
  network capture, or unpinned `main` downloads.
- **Tests/security/docs/CI:** action metadata/YAML validation, fixture-based
  integration, malicious-input/path tests, permission audit, hosted Action CI;
  document compatibility and rollback to the pre-release guard.
- **Acceptance/evidence:** reproducible Action run against a released core,
  explicit compatibility matrix, no secret exposure, tagged Action release.
- **Commands:** repository-specific metadata/YAML checks, core fixture validation,
  and hosted workflow inspection. Live network: no. Maintainer authorization:
  yes. Review: independent security review strongly recommended.

#### Phase 15 — Documentation adoption baseline

- **Parent:** Era 2 / Milestone 12. **Status:** CANDIDATE. **Objective:** make
  the docs repository useful without duplicating core truth.
- **Why/dependencies:** depends on v0.1.0 and Phase 13 terminology. **Repositories:**
  docs primarily; core links and compatibility notes.
- **In scope:** quickstart, fixture concepts, CI integration guide, troubleshooting
  navigation, release compatibility, and link checking. **Out of scope:** heavy
  framework, generated API duplication, or claims of user adoption.
- **Tests/security/docs/CI:** local Markdown/link check, version-target review,
  no credentials in examples, docs CI with meaningful checks. Evidence is a
  versioned docs commit and guide review.
- **Acceptance/evidence:** a new user can follow core README then docs guides;
  every behavior claim links to core authority. Live network: no. Authorization:
  yes. Agent: Codex; content review may use an independent reader.

#### Phase 16 — Evidence-backed compatibility expansion

- **Parent:** Era 3 / Milestone 13. **Status:** CANDIDATE. **Objective:** add
  only a method, SDK adapter, or conformance asset demonstrated by usage.
- **Dependencies/preconditions:** adoption evidence, a written protocol decision,
  fixture schema compatibility plan, and security review. **Repositories:** core;
  a new repository is not presumed.
- **In scope:** one approved read-only capability at a time with fixture,
  capture, replay, docs, and tests. **Out of scope:** writes, generic proxying,
  batch redesign, or speculative ecosystem breadth.
- **Acceptance/evidence:** protocol source, versioned schema/compatibility note,
  adversarial tests, docs, CI, and a phase report. Live network: opt-in research
  only. Maintainer authorization: yes. Agent: Codex plus independent review if
  semantics or security are contentious.

#### Phase 17 — Provenance and artifact trust evaluation

- **Parent:** Era 4 / Milestone 14. **Status:** SPECULATIVE. **Objective:**
  evaluate whether stronger fixture/release provenance is needed by real users.
- **Dependencies:** adoption evidence and threat model; no implementation is
  implied. **Repositories:** core, possibly release infrastructure.
- **In scope:** threat model, signed-release or signed-fixture options, key
  custody, verification UX, and migration impact. **Out of scope:** claiming
  signatures prove chain truth or introducing keys without governance.
- **Acceptance/evidence:** documented decision, security review, cost/benefit,
  rollback and key-compromise plan. Live network: no. Maintainer authorization:
  required before any signing implementation. Recommended agent: Codex for
  analysis, independent security review before design approval.

#### Phase 18 — Maturity and support review

- **Parent:** Era 5 / Milestone 15. **Status:** SPECULATIVE. **Objective:**
  reassess whether v1.0 criteria are met rather than mechanically pursuing them.
- **Dependencies:** sustained releases, compatibility evidence, support patterns,
  and clear non-goals. **Repositories:** all existing repositories; no new repo
  unless the audit proves a boundary.
- **In scope:** API/CLI stability, deprecation policy, platform support, CI
  reliability, governance, security response, and user evidence. **Out of scope:**
  enterprise infrastructure or unrelated product expansion.
- **Acceptance/evidence:** independent maturity audit, v1.0 decision record,
  migration plan, release checklist, and repository action plan update. Live
  network: no by default. Maintainer authorization: yes; independent audit is
  recommended.

## Out-of-scope guardrail

Browser UI, transaction signing/submission, generic RPC proxying, AI decisions,
hosted replay, and broad blockchain support remain outside the product unless a
future strategy review explicitly changes the constitution. They cannot be
smuggled in as “developer experience” work.

