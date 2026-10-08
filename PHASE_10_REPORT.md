# Phase 10 Report — CI and cross-platform release

Status: COMPLETE

Phase 10 establishes credible cross-platform verification and a controlled
release-artifact path without publishing an unapproved v0.1.0 release.

## Delivered

- Replaced the bootstrap CI workflow with a Linux, Windows, and macOS matrix.
- Added formatting, tests, vet, CLI build, CLI smoke, and clean-tree checks.
- Added Linux race-detector coverage; Windows and macOS run the full non-race
  verification matrix.
- Added reproducible static release packaging for Linux amd64, Windows amd64,
  and macOS amd64.
- Added platform-specific packaging and archive-content smoke-test scripts.
- Embedded the tagged release version in the CLI binary through release build
  flags and verified it in the packaging workflow.
- Added checksum generation and verification before the tagged release publish
  step, with least-privilege read permissions for packaging and write permission
  only on the publish job.
- Added `docs/RELEASES.md` with archive contents, installation, checksum,
  compatibility, and signing limitations.

## Local verification

Executed from the core repository:

```text
go test -p 1 -count=1 ./...                         PASS
go test -race -p 1 -count=1 ./...                   PASS
go vet ./...                                         PASS
go test -count=1 ./internal/doccheck                PASS
git diff --check                                     PASS
```

Packaging smoke checks also passed locally for the Windows ZIP and Linux
tarball. The Windows archive was extracted and its embedded version executed;
the Linux archive contents and checksum were verified locally. The Linux binary
was not executed on Windows because it is intentionally a Linux target.

## Hosted verification

The final implementation checkpoint is commit
[`5ef43f20f5c39e966bf0eaf13999fa51bc16c881`](https://github.com/StellarReplay/stellar-replay/commit/5ef43f20f5c39e966bf0eaf13999fa51bc16c881).

The final hosted matrix passed:

- [GitHub Actions run 37773103501](https://github.com/StellarReplay/stellar-replay/actions/runs/37773103501)
  — success.
- Linux: formatting, tests, vet, build, CLI smoke, race tests, and clean tree.
- Windows: build, tests, vet, CLI smoke, and clean tree; formatting coverage is
  enforced on Unix runners because Windows checkout line endings make the direct
  `gofmt -l` directory check report a false difference.
- macOS: formatting, tests, vet, build, CLI smoke, and clean tree.

The hosted run emitted only routine runner/action deprecation, cache, and image
migration annotations; no verification step failed.

## Release boundary

The tagged release workflow is implemented and locally exercised through its
packaging and smoke-test scripts. It is intentionally not triggered during this
phase because doing so would publish a GitHub release; public v0.1.0 publication
belongs to Phase 12 after the Phase 11 audit and go/no-go decision. Signing and
registry publishing remain out of scope.

## Acceptance result

- Declared Windows/Linux/macOS CI matrix: PASS.
- Go formatting, tests, vet, build, and race coverage: PASS.
- Archive contents contain only the intended binary, license, and README: PASS.
- Embedded version and archive smoke tests: PASS.
- Checksum generation and verification path: implemented and locally reviewed.
- Least-privilege workflow permissions: PASS.
- Installation and release limitations documented: PASS.

Phase 11 — release candidate and final audit — remains `NOT_STARTED`.
