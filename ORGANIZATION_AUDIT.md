# StellarReplay Organization Audit

Audit date: 2026-10-08  
Authority: authenticated GitHub organization state, repository contents, Git history, workflow history, release metadata, and local verification.

## Executive finding

The organization contains exactly three public, non-archived repositories. The
planned three-repository topology is valid, but only the core is currently a
released product. The Action and docs repositories are legitimate companion
boundaries in pre-release foundation state; they should not be filled with
placeholder implementation merely to make the organization look active.

## Inventory

| Repository | Category | Lifecycle | Evidence-backed state | Recommendation |
|---|---|---|---|---|
| [stellar-replay](https://github.com/StellarReplay/stellar-replay) | CORE PRODUCT | ACTIVE / RELEASED | Go implementation, tests, CI, release `v0.1.0`, tag and verified platform artifacts | KEEP / DEVELOP through post-release productization |
| [stellar-replay-action](https://github.com/StellarReplay/stellar-replay-action) | INTEGRATION | PRE-RELEASE / DORMANT | Two foundation commits, composite metadata with an intentional failing guard, no tag/release, successful metadata-only CI | KEEP / REPAIR documentation, then develop only after the Action contract is approved |
| [stellar-replay-docs](https://github.com/StellarReplay/stellar-replay-docs) | DOCUMENTATION | PRE-RELEASE / DORMANT | Two foundation commits, navigation placeholder, no site release, successful foundation-only CI | KEEP / REPAIR content plan, then develop adoption documentation |

All three repositories are public, use `main`, are not archived or forks, and
have no open issues or pull requests at audit time. The Action and docs
repositories have no tags or releases. No fourth repository was found.

## Repository health

### stellar-replay

**Verified:** The repository contains the Go core, CLI, fixture implementation,
capture, replay, local server, tests, examples, security controls, CI, release
workflow, and documentation. Local tests, race tests, vet, build, documentation
checks, and release smoke checks passed in the phase evidence. GitHub release
`v0.1.0` exists at the linked release page. Release run
[37778472965](https://github.com/StellarReplay/stellar-replay/actions/runs/37778472965)
and post-release CI run
[37778852876](https://github.com/StellarReplay/stellar-replay/actions/runs/37778852876)
passed. Main is at commit `f9b5624` and is synchronized with origin.

**Assessment:** Released v0.1 product, not merely foundation or research. The
next risk is adoption and maintenance rather than missing v0.1 core machinery.

### stellar-replay-action

**Verified:** The repository contains `action.yml`, governance files, examples
and test guidance, but its composite action intentionally exits with a
pre-release error. Its CI checks metadata presence and uses read-only workflow
permissions; it does not claim a working Action release. The latest commit is
`64d64c7`, and Action CI run
[37720763406](https://github.com/StellarReplay/stellar-replay-action/actions/runs/37720763406)
passed.

**Reason for limited work:** The repository README and history explicitly state
that it was established as a thin boundary while the core interface was being
built. It was created before the core release and has only two same-day
foundation commits.

**Classification:** VERIFIED pre-release foundation; LIKELY intentionally
dormant pending a contract decision; UNKNOWN whether maintainers want the first
integration to be a downloaded binary, a container, or another stable interface.

### stellar-replay-docs

**Verified:** The repository contains ownership, contribution, security,
versioning, and a documentation navigation placeholder, but no guide set or
documentation site. Its CI checks the foundation files only. The latest commit
is `e987eb5`, and Docs CI run
[37720763216](https://github.com/StellarReplay/stellar-replay-docs/actions/runs/37720763216)
passed.

**Reason for limited work:** Its README and history explicitly state that it was
created as an extended-documentation boundary before the core behavior became
real. This is evidence of sequencing, not evidence of abandonment.

**Classification:** VERIFIED pre-release documentation foundation; LIKELY
intentionally dormant pending a content plan; UNKNOWN whether a site framework
will eventually be justified.

## Cross-repository findings

The architecture is sound when ownership remains centralized:

```text
stellar-replay
  ├─ canonical fixture schema, CLI, core behavior, releases
  ├─ consumed by the Action through a released interface
  └─ explained in depth by the docs repository
```

The Action must never copy the Go engine, and the docs repository must not
redefine fixture or protocol semantics. The core owns schema versioning and
compatibility notes. Core, Action, and docs releases are separate artifacts
linked by explicit compatibility statements.

The only material documentation drift found was release-state wording: the
core README and companion repositories still described the core as unreleased.
That drift is repaired in this audit. `docs/history/audits/FINAL_AUDIT.md` is
not edited because it is historical evidence from before the release.

## Security and governance findings

- All inspected workflows use `permissions: contents: read`; no credentials,
  signing keys, or blockchain secrets were introduced.
- Checked-in core fixtures are synthetic according to the completed phase
  evidence; fixture review remains a maintainer responsibility before sharing.
- Release archives are checksum-verified but unsigned in v0.1.0, as documented.
- GitHub reports `NOASSERTION`/Other license metadata even though each repository
  contains an Apache-2.0 `LICENSE` file. This is a GitHub metadata-detection
  inconsistency to monitor, not evidence that the license files are absent.
- Branch protection is not enabled on the small companion repositories. This is
  a governance improvement candidate, not a reason to invent product code.

## Evidence limits

Issue and pull-request counts were zero at audit time. That proves no visible
open work was found; it does not prove that no private planning or external
discussion exists. “Abandoned,” “blocked,” and “intentionally dormant” are not
used as interchangeable claims.
