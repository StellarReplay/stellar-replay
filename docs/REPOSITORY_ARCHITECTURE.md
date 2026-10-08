# Cross-Repository Architecture

## Ownership map

| Repository | Owns | Consumes | Must not duplicate |
|---|---|---|---|
| `stellar-replay` | Go core, CLI, fixture schema/semantics, capture, normalization, sanitization, hashing, replay, local server, core tests, release binaries | Stellar RPC protocol decisions and its own compatibility inputs | Nothing outside its core contract; it remains the canonical behavior source |
| `stellar-replay-action` | Action metadata, wrapper/integration behavior, Action examples/tests, Action releases | A released core CLI/interface and the canonical fixture format | The Go replay engine, fixture semantics, or an unreleased core commit |
| `stellar-replay-docs` | Tutorials, conceptual guides, integration material, navigation, contributor education | Released core behavior and explicit compatibility targets | Protocol/schema truth, executable core instructions, or contradictory security claims |

## Version relationship

The core uses semantic version tags and released `v0.1.0`. The Action will use
its own semantic version and major-version tags once it is functional; each
Action release must state the supported core version and fixture schema. Docs
commits and any future site versions identify the core release they describe.
There is no monolithic three-repository release.

## Change flow

Core behavior or schema changes are authored and tested in `stellar-replay`
first. The Action and docs repositories then consume the released interface or
documented schema version, add compatibility notes, run their own checks, and
release independently. Cross-repository links are audited at phase closure.

## Future repository criteria

Do not create a repository for convenience, empty planning, or cosmetic scale.
A split requires a distinct owner, release cadence, dependency/security
boundary, contributor workflow, or artifact lifecycle that cannot remain clear
inside an existing repository. Candidate areas are SDK/integrations,
conformance fixtures, and ecosystem tooling; none is currently approved.
