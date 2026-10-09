# StellarReplay Independent Organization Audit

Audit date: 2026-10-09  
Audit branch: `audit/2026-10-09-independent`  
Baseline: `bc2422477d01b0fbbfb643064b0767b6d7136052`

## Executive conclusion

The StellarReplay organization is usable and presents a credible early open-source
product, but it is not yet fully hardened for a high-trust external audit. The
core repository is a released v0.1.0 product with reproducible release assets,
passing local and hosted tests, and a clear deterministic RPC-testing contract.
The Action repository is intentionally pre-release and the docs repository is a
useful documentation foundation. No P0 defect or core-behavior P1 defect was
found in this audit.

The most important remaining work is governance and evidence: enable repository
security analysis, establish protected-main rules, improve Action/docs CI depth,
and prepare a concrete post-v0.1 productization milestone. These are owner or
organization controls, not changes that should be silently made during this
audit.

## Audit method and evidence standard

The audit inspected all repositories visible in the `StellarReplay` GitHub
organization, local checkouts, branches, remotes, commit history, tags, releases,
workflows, hosted workflow runs, source trees, tests, fixtures, documentation,
cross-repository links, and release archives. Statements are marked as verified
facts, strong inferences, or unknowns where relevant. No repository was deleted,
archived, force-pushed, merged, or released during this audit.

## Organization inventory

The authenticated GitHub inventory contains exactly three public, non-archived
repositories. All default to `main`, currently have no open issues or pull
requests, and have no published releases except the core v0.1.0 release.

| Repository | Classification | Current state | Lifecycle recommendation |
|---|---|---|---|
| `StellarReplay/stellar-replay` | Core product | Go implementation, tests, fixtures, CLI, local replay server, CI, and v0.1.0 release | ACTIVE / MAINTAIN |
| `StellarReplay/stellar-replay-action` | Integration | Professional foundation and metadata, but intentionally exits with a pre-release guard and is not advertised as working | PRE-RELEASE / BLOCKED on stable integration contract |
| `StellarReplay/stellar-replay-docs` | Documentation | Source-controlled extended documentation with independent CI and links to core | ACTIVE / MAINTAIN |

No fourth organization repository was found. No local fourth StellarReplay
checkout was found. The `.go-build-phase*` directories under the local parent
are ignored build caches, not repositories or product components.

### Repository evidence

* Core default-branch head before this audit: `bc2422477d01b0fbbfb643064b0767b6d7136052`.
* Action default-branch head: `a64d098216f236a27853ba120f0b0e4236d05aaf`.
* Docs default-branch head: `18a0b4a423b39aed318863ec71ecd966378e77ab`.
* All three remotes point to the intended `StellarReplay` organization.
* All three repositories are public. GitHub currently reports `NOASSERTION` for
  the API license field even though each repository contains an Apache-2.0
  license file; this is a metadata-quality issue, not evidence that the license
  file is absent.
* The organization profile metadata is sparse and no visible `.github` profile
  repository was found. Avatar/profile presentation was not independently
  verified through the API.

## Core product assessment

### Verified implementation contract

The core implements the documented v0.1 boundaries:

* Supported RPC methods are `getHealth`, `getLatestLedger`, `getNetwork`, and
  `getLedgerEntries`.
* Capture is explicit, read-only by default, HTTPS-only, bounded, no-redirect,
  no-retry, and no-proxy.
* Fixtures use a versioned schema, strict validation, canonical JSON, SHA-256
  integrity, sanitization, bounded loading, and atomic writes.
* Replay uses exact canonical request matching, checks fixture integrity, and has
  no live-network fallback.
* The local server is loopback-only, bounded, non-batch, and shuts down cleanly.
* The CLI provides `record`, `inspect`, `validate`, `replay`, and `serve`.
* Checked-in fixtures are synthetic and contain no credentials or private keys.

### Maturity decision

The core is beyond foundation and research. It is best described as a released
early developer product: v0.1.0 is published, but adoption, compatibility
feedback, and ecosystem integrations are still immature. The organization as a
whole has entered productization, not mature-product status.

## Health and verification results

### Core

Passed locally:

```text
go run ./cmd/stellar-replay help
go run ./cmd/stellar-replay version
go test -count=1 ./examples
go test -p 1 -count=1 ./...
go test -race -p 1 -count=1 ./...
go vet ./...
go build -trimpath ./cmd/stellar-replay
go test -count=1 ./internal/doccheck
git diff --check
```

The source build reports `0.1.0-dev`; tagged release builds inject `0.1.0`.
The hosted core CI run [37859361007](https://github.com/StellarReplay/stellar-replay/actions/runs/37859361007)
passed its Linux, macOS, and Windows matrix.

### Action

Hosted CI run [37859368105](https://github.com/StellarReplay/stellar-replay-action/actions/runs/37859368105)
passed its current foundation checks. The workflow currently verifies that
`action.yml` and `README.md` exist; it does not yet exercise a working Action.
This is consistent with the repository's explicit pre-release status.

### Docs

Hosted CI run [37859374858](https://github.com/StellarReplay/stellar-replay-docs/actions/runs/37859374858)
passed its current documentation-presence checks. Deeper link/build checking is
still a recommended improvement.

## Release verification

The core has a published, non-draft, non-prerelease `v0.1.0` release at
[`StellarReplay/stellar-replay/releases/tag/v0.1.0`](https://github.com/StellarReplay/stellar-replay/releases/tag/v0.1.0).
The tagged release workflow [37778472965](https://github.com/StellarReplay/stellar-replay/actions/runs/37778472965)
passed. Release archives were downloaded independently and their SHA-256
digests matched the published `SHA256SUMS` asset. Linux, Windows, and macOS
archives contained the expected executable plus `LICENSE` and `README.md`.

Release artifacts are not publisher-signed. Checksums detect alteration but do
not prove publisher authenticity. This remains an accepted v0.1 limitation and
should be revisited before higher-trust distribution.

## Security and governance findings

### P0 — none found

No exposed credential, private key, unsafe live-network fallback, destructive
default, or release-integrity failure was found in the inspected repositories or
release assets.

### P1 — owner-controlled hardening required

1. GitHub Dependabot security updates are disabled for all three repositories.
2. Secret scanning and push protection are disabled for all three repositories.
3. No repository-level branch protection was reported for `main`.

These are verified GitHub settings observations. They were not changed because
this audit must not silently alter organization security policy or bypass owner
approval. They should be addressed before presenting the organization as fully
audit-hardened.

### P2 — recommended improvements

* Add real Action metadata/YAML validation and tests when implementation begins.
* Add stronger docs link checking/build validation.
* Set GitHub license metadata where supported and add organization profile
  presentation (description, links, and visual identity) deliberately.
* Consider signed release provenance before a broader production release.
* Continue conservative fixture review and redaction before sharing captured RPC
  data.

No signing keys, blockchain credentials, GitHub tokens, or environment secrets
were added or committed during the audit.

## Cross-repository architecture

```text
stellar-replay
  owns the CLI, fixture schema, capture, validation, replay, and server
       |
       +--> stellar-replay-action
       |      consumes a stable released CLI/interface; must not copy Go core
       |
       +--> stellar-replay-docs
              explains and teaches the product; must not redefine behavior
```

The three-repository strategy remains justified. The core owns the canonical
fixture schema and replay semantics. The Action owns only GitHub Actions
integration. The docs repository owns extended guides and navigation. No
additional repository is justified by current evidence.

The Action is intentionally not coupled to an unreleased core commit. Before it
becomes functional, its compatibility contract must name the supported core
release range and pin or otherwise verify the binary/interface it invokes.

## Inactive or underworked repository diagnosis

### `stellar-replay-action`

* **Current state:** pre-release foundation.
* **Verified evidence:** metadata explicitly describes a future integration;
  `action.yml` contains a pre-release guard; no release or implementation test
  exists; CI passes only foundation checks.
* **Probable reason:** the stable core interface was established before building
  the integration. This is a strong inference, not a recorded maintainer claim.
* **Confidence:** high for current state; medium for the reason.
* **Impact:** useful to maintain as a contract/foundation, but not yet suitable
  for users to install as a working Action.
* **Recommendation:** retain as PRE-RELEASE; implement only after a released
  integration contract and end-to-end test plan exist.

### `stellar-replay-docs`

* **Current state:** active documentation foundation with limited CI depth.
* **Verified evidence:** independent repository, current links, documentation
  structure, and passing hosted foundation checks.
* **Probable reason for lighter code activity:** documentation is intentionally
  separated from implementation; no evidence supports calling it abandoned.
* **Confidence:** high.
* **Recommendation:** retain and maintain; deepen link/build checks as content
  grows.

## Product gap analysis

| Gap | Needed capability | Evidence before completion |
|---|---|---|
| v0.1 to reliable maintenance | Compatibility policy, issue intake, regression discipline, security settings | Protected main, security tooling, documented support policy, recurring CI |
| Reliable release to adoption | Fast demo, clearer onboarding, real-world fixture examples, feedback loop | External users can install and reproduce a deterministic replay |
| Adoption to ecosystem | Working GitHub Action and stable integration contract | End-to-end Action tests against a released core |
| Ecosystem to mature product | Compatibility guarantees, provenance/signing, richer observability only if demanded | Maintainer and user evidence, security review, release evidence |

The product should not become a wallet, explorer, signing system, transaction
submitter, Horizon replacement, generic RPC proxy, or generic mocking platform.

## Funding and external-readiness assessment

This audit does not guarantee acceptance by any funder. It identifies evidence
needed for a credible application.

### Drips

Drips' current Wave documentation describes recurring contribution cycles,
maintainer onboarding, issue scoping, review, and responsiveness expectations;
participation is described as free. The project should claim/configure its repo,
install the relevant GitHub integration, and prepare small, well-scoped issues
before applying or participating. See the [Wave overview](https://docs.drips.network/wave/),
[maintainer participation guide](https://docs.drips.network/wave/maintainers/participating-in-a-wave/),
and [maintainer FAQ](https://docs.drips.network/wave/maintainers/faq/).

Drips RPGF applications require a claimed Drips Project and round-specific
information; eligibility and selection remain organizer-controlled. See
[claiming a repository](https://docs.drips.network/get-support/claim-your-repository/)
and [applying to a round](https://docs.drips.network/rpgf/apply-to-a-round/).

**Assessment:** potentially suitable, but the application needs a concise user
problem, demonstrated v0.1 evidence, maintainer capacity, and concrete issues
that contributors can complete.

### GrantFox

The official [GrantFox site](https://grantfox.xyz/) and [documentation](https://docs.grantfox.xyz/)
describe an open-source contribution ecosystem centered on issues, bounties,
reputation, and rewards. The public material inspected did not establish a
current, authoritative acceptance rubric for StellarReplay.

**Assessment:** potentially suitable for issue/bounty-style contribution
discovery, but eligibility and funding likelihood are UNKNOWN. Treat it as an
outreach and contributor-readiness channel until current organizer guidance is
confirmed.

### Stellar Community Fund context

The official [SCF awards page](https://communityfund.stellar.org/awards)
describes Build rounds, tracks, submission quality, use of Stellar, integration
planning, readiness to build, and budget/tranches. StellarReplay is not an
obvious end-user Stellar application, so an application would need to frame its
developer-infrastructure value and show concrete Stellar/Soroban adoption or
integration evidence. **Assessment:** potentially suitable after specific fixes;
not an acceptance prediction.

## Long-term product direction

The roadmap should be understood as:

```text
Vision → Era → Milestone → Phase → Task → Verification evidence
```

Current position:

* **Era:** foundation completed; early developer-product productization.
* **Milestone:** v0.1 core release completed.
* **Next milestone:** harden maintenance/governance and validate real developer
  adoption before expanding the product surface.
* **Candidate later milestone:** a released, tested GitHub Action using the
  stable core interface.
* **Speculative later work:** ecosystem integrations, provenance/signing, or
  additional SDKs only after user evidence justifies them.

Future phases should retain explicit scope, dependencies, affected repositories,
security requirements, tests, documentation, CI, acceptance criteria, evidence,
rollback considerations, release impact, and maintainer-authorization gates.

## Required owner actions before claiming full audit readiness

1. Enable Dependabot security updates, secret scanning, and push protection where
   available and appropriate.
2. Protect `main` in each repository with required checks and review rules.
3. Decide whether to add organization profile metadata and correct GitHub license
   metadata.
4. Create a post-v0.1 milestone with a small number of externally verifiable
   tasks, including the Action contract and deeper docs checks.
5. For funding submissions, prepare the project narrative, maintainer capacity,
   milestone budget, contribution issues, and links to the published release and
   passing CI.

## Changes made and not made

Made in this audit:

* Added this evidence-backed audit report on the dedicated audit branch.
* Verified organization inventory, local synchronization, releases, assets,
  workflows, tests, cross-repository boundaries, and funding-readiness context.

Not made:

* No core product feature changes.
* No changes to GitHub security settings or branch protection.
* No release, tag, merge, force push, archive, repository deletion, or settings
  bypass.
* No edits to historical Phase 0–12 reports.
* No new repository.

## Current state and next step

**Where Stellar Replay is today:** a real, released, usable v0.1 developer tool
with a healthy core foundation and credible public presentation, transitioning
from implementation into productization. It is suitable to present for external
review as an early open-source developer-infrastructure project, but should not
be described as fully mature or fully security-hardened until the owner actions
above are complete.

**Next highest-value executable milestone:** repository governance and adoption
readiness: harden GitHub controls, create a small contributor-ready backlog,
improve Action/docs validation, and gather external usage evidence before
building additional product features.

