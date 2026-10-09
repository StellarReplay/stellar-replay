# Final Hardening Follow-up

Date: 2026-10-09  
Scope: `stellar-replay`, `stellar-replay-action`, and `stellar-replay-docs`

This is a follow-up to the historical [independent organization audit](../../../AUDIT_REPORT.md).
The historical report is preserved unchanged. This record documents the
confirmed findings, remediation, and final hosted verification.

## Findings and dispositions

| Finding | Result | Evidence |
|---|---|---|
| Generic credential-bearing query keys such as `token` could pass endpoint validation | Fixed | Core now rejects normalized token keys, including underscore, hyphen, and case variants; capture and fixture regression tests cover them. |
| Capture could accept a response with a different JSON-RPC ID | Fixed | Capture seals through fixture validation, which now requires request/response scalar IDs to match without string/number coercion. |
| Action and docs repositories contained abbreviated license text | Fixed | Both repositories now contain the complete Apache-2.0 text matching the core license content. |
| Action CI checked only file existence | Fixed | Deterministic metadata and fail-closed contract validation is implemented in `tests/validate-action.sh`. |
| Docs CI checked only two non-empty files | Fixed | Deterministic Markdown fence and local-relative-link validation is implemented in `tests/check_markdown.py`; the workflow uses Ubuntu's `python3`. |

## Pull requests and merged commits

| Repository | Pull request | Merge commit | Main CI |
|---|---|---|---|
| Core | [#2](https://github.com/StellarReplay/stellar-replay/pull/2) | `dceef4493ff11007538d62aaaf5d2077caf9be99` | [run 37968752044](https://github.com/StellarReplay/stellar-replay/actions/runs/37968752044) |
| Action | [#1](https://github.com/StellarReplay/stellar-replay-action/pull/1) | `b2e20d6654f1317510c1eae803aa9c410285e1e4` | [run 37968752638](https://github.com/StellarReplay/stellar-replay-action/actions/runs/37968752638) |
| Docs | [#1](https://github.com/StellarReplay/stellar-replay-docs/pull/1) | `2868b99329ac5258db69388bd669e214919ad656` | [run 37968763648](https://github.com/StellarReplay/stellar-replay-docs/actions/runs/37968763648) |

All three PRs were merged only after their required checks passed. Each merged
main workflow completed successfully. Temporary remote and local hardening
branches were deleted. All three local checkouts are clean and synchronized to
their respective `main` branches.

## Verification performed

Core local verification passed:

```text
go test -p 1 -count=1 ./...
go test -race -p 1 -count=1 ./...
go vet ./...
go build -trimpath ./cmd/stellar-replay
git diff --check
```

Focused regression tests covered generic token query keys, sensitive fixture
fields, matching IDs, mismatched IDs, and scalar type boundaries. The Action
metadata and pre-release guard check passed. A local Markdown relative-link
audit passed for the docs repository.

## Main-branch governance

Each repository has an active repository ruleset named `protect-main` targeting
the default branch. The verified rules require pull requests, prevent deletion
and non-fast-forward updates, require conversation resolution, and require
up-to-date status checks. No bypass actors are configured and the approving
review count remains zero for the solo-maintainer workflow.

Required checks are:

* Core: `Verify (ubuntu-latest)`, `Verify (windows-latest)`, and
  `Verify (macos-latest)`.
* Action: `metadata`.
* Docs: `foundation`.

## Security settings

The following settings are verified enabled for all three repositories:

* Dependabot alerts/security updates.
* Secret scanning.
* Secret scanning push protection.

GitHub's repository metadata still reports `NOASSERTION` for the SPDX license
field. This is a GitHub metadata refresh/recognition limitation; it does not
mean the complete Apache-2.0 `LICENSE` files are absent. The files and README
license references are consistent.

## Remaining limitations

The Action remains intentionally pre-release and fail-closed. This hardening
pass validates its metadata and contract without advertising a working
end-to-end integration. The release archives remain checksum-protected but not
publisher-signed, as previously documented for v0.1.0. Neither limitation was
silently expanded or reclassified by this follow-up.

