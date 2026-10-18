# Submission Readiness Record — 2026-10-18

Status: READY FOR EXTERNAL SUBMISSION, subject to platform onboarding and
review. This record is a follow-up to the historical audit documents; it does
not rewrite their original findings or dates.

## Decision

Stellar Replay is suitable to submit as an early open-source developer
infrastructure project to Drips Wave and GrantFox. The core repository is a
released v0.1.0 product with deterministic offline replay, a documented
fixture contract, CI evidence, and an active contributor surface. The Action
repository is implemented against the released core but remains untagged. The
documentation repository is useful and synchronized but remains pre-release.

This is a readiness decision, not a claim of acceptance, funding, or platform
approval.

## Verified repository state

| Repository | Verified state | Evidence |
| --- | --- | --- |
| `StellarReplay/stellar-replay` | Public core product; v0.1.0 released; `main` protected; core security and JSON-RPC ID hardening merged | [PR #12](https://github.com/StellarReplay/stellar-replay/pull/12), [v0.1.0 release](https://github.com/StellarReplay/stellar-replay/releases/tag/v0.1.0) |
| `StellarReplay/stellar-replay-action` | Functional composite Action that downloads a pinned released core archive, verifies `SHA256SUMS`, and runs offline validate/replay; no Action release tag yet | [PR #3](https://github.com/StellarReplay/stellar-replay-action/pull/3), [Action CI](https://github.com/StellarReplay/stellar-replay-action/actions) |
| `StellarReplay/stellar-replay-docs` | Public documentation repository with offline quickstart and Action integration guide; pre-release documentation surface | [PR #3](https://github.com/StellarReplay/stellar-replay-docs/pull/3), [Docs CI](https://github.com/StellarReplay/stellar-replay-docs/actions) |

All three repositories use `main` as the default branch. Task branches were
merged through pull requests and removed after merge. No tag or release was
created for the Action or docs repository during this work.

## Verification performed

The core hardening baseline verified:

- credential-bearing endpoint query keys, including generic `token` variants,
  are rejected before storage and through fixture validation paths;
- request and response JSON-RPC scalar IDs must match without string/number
  coercion;
- supported methods remain `getHealth`, `getLatestLedger`, `getNetwork`, and
  `getLedgerEntries`;
- replay and local serving do not fall back to the live network;
- the core fixture format remains version 1 and integrity-gated.

The core checks passed before the documentation-only submission record:

```text
go test -p 1 -count=1 ./...
go test -race -p 1 -count=1 ./...
go vet ./...
go build -trimpath ./cmd/stellar-replay
go test -count=1 ./internal/doccheck
git diff --check
```

The hosted required checks passed for the merged core, Action, and docs pull
requests. The Action workflow also passed its offline fixture smoke test and
checksum contract checks. The docs workflow passed its Markdown/link
foundation checks.

## Contributor surface

The following scoped issues are open in the core repository and are suitable
for maintainer triage or external contribution:

- [#4 — Soroban RPC fixture integration example](https://github.com/StellarReplay/stellar-replay/issues/4)
- [#7 — getLedgerEntries pagination fixtures](https://github.com/StellarReplay/stellar-replay/issues/7)
- [#6 — prove replay network isolation](https://github.com/StellarReplay/stellar-replay/issues/6)
- [#11 — expand Action checksum and offline test matrix](https://github.com/StellarReplay/stellar-replay/issues/11)
- [#10 — cross-platform Action smoke tests](https://github.com/StellarReplay/stellar-replay/issues/10)
- [#9 — fixture schema compatibility tests](https://github.com/StellarReplay/stellar-replay/issues/9)
- [#5 — fuzz fixture sanitization boundaries](https://github.com/StellarReplay/stellar-replay/issues/5)
- [#8 — troubleshooting and migration guidance](https://github.com/StellarReplay/stellar-replay/issues/8)

These issues intentionally do not promise transaction signing, secret-key
handling, generic RPC proxying, Horizon behavior, or an expanded supported
method set.

## Drips Wave submission brief

Recommended submission unit: the core repository plus a small set of the
scoped issues above. The Action and docs repositories should be linked as
supporting project components, not presented as separate competing products.

Before submitting, a maintainer must onboard the `StellarReplay` organization
in Drips, install the Wave GitHub App, sync/apply the selected repository, and
wait for organizer approval. Select issues with clear scope and complexity;
review and assignment require active maintainer attention during a wave. The
platform's rules and organizer selection remain external to this repository.

## GrantFox submission brief

Recommended submission unit: an open-source contribution campaign centered on
the core repository and selected issues, with the Action and docs repositories
linked as integration and learning surfaces. Describe the user problem,
deterministic replay value, v0.1.0 evidence, contribution boundaries, and the
exact issue acceptance criteria.

GrantFox's public materials describe an ecosystem for publishing contribution
opportunities, reviewing work, and rewarding open-source contributions. Current
campaign eligibility, funding, and acceptance conditions must be confirmed in
the GrantFox submission flow; this record does not claim that a campaign will
be accepted or funded.

## Remaining limitations

1. The Action has no published tag/release and should not be advertised as a
   versioned marketplace action yet. Consumers should use the documented
   development reference only with that limitation understood.
2. Core archives provide checksums, not publisher-signature authenticity.
3. GitHub's API license metadata may remain `NOASSERTION` even though all three
   repositories contain the complete Apache-2.0 text; metadata refresh is a
   GitHub concern, not a missing license file.
4. External platform onboarding, approval, contributor matching, and funding
   are not controlled by Stellar Replay.

## Next executable milestone

The next product milestone is contributor adoption and release hardening, not
another repository split: publish a versioned Action only after its compatibility
and cross-platform matrix are proven, then use the scoped issues to gather
external contribution evidence. A maintainer-authorized Action release remains
required before calling the Action production-ready.

## Subsequent productization update

After this record's original submission checkpoint, the Action was released as
`v0.1.0` and the documentation examples were synchronized to the immutable tag.
Phases 13–15 are now recorded as complete in the roadmap. This addendum is a
later event and intentionally does not alter the historical state described
above.
