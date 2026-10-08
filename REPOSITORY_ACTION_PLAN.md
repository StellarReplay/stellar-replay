# Repository Action Plan

This plan records non-destructive lifecycle decisions from the 2026-10-08
organization audit. “Dormant” means intentionally awaiting a justified next
milestone; it does not mean deleted or abandoned.

| Repository | Decision | Why | Resume/develop gate |
|---|---|---|---|
| `StellarReplay/stellar-replay` | KEEP / DEVELOP | Released v0.1.0 with working core, tests, CI, and artifacts | Execute the post-release stabilization/adoption phase with maintainer authorization |
| `StellarReplay/stellar-replay-action` | KEEP / REPAIR / DORMANT | Legitimate integration boundary, but the Action intentionally fails closed and has no release | Approve invocation/download contract, compatibility policy, and security model first |
| `StellarReplay/stellar-replay-docs` | KEEP / REPAIR / DORMANT | Legitimate depth-documentation boundary, but currently only foundation/navigation | Approve initial guide set and core-version target; choose framework only if content requires it |

## Decisions not taken

No repository is archived, merged, renamed, deleted, or replaced. No fourth
repository is justified by current evidence. A future SDK/integration or
conformance repository may be proposed only when its ownership, release, and
contributor boundary is materially distinct from the core.

## Priority order

1. Maintain the released core and keep its v0.1 contract truthful.
2. Define and implement the smallest secure Action integration only after the
   interface is approved.
3. Build a useful docs adoption baseline without duplicating core truth.
4. Reassess repository count only after real usage produces a boundary problem.
