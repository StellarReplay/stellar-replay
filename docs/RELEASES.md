# Releases and installation

Phase 10 establishes the release boundary for the core `stellar-replay` CLI.
Tagged core releases use semantic version tags such as `v0.1.0` and publish
platform archives plus `SHA256SUMS` on the GitHub release page.

## Supported release targets

The initial release workflow declares:

| Target | Archive |
|---|---|
| Linux amd64 | `stellar-replay-VERSION-linux-amd64.tar.gz` |
| Windows amd64 | `stellar-replay-VERSION-windows-amd64.zip` |
| macOS amd64 | `stellar-replay-VERSION-darwin-amd64.tar.gz` |

Each archive contains only the platform binary, `LICENSE`, and `README.md`.
The release workflow smoke-tests archive contents, checks the embedded CLI
version, generates SHA-256 checksums, and verifies the checksum file before
publishing. Release artifacts are not signed in v0.1.

## Install from an archive

Download the archive matching the operating system and CPU architecture, extract
it, and place the binary on `PATH`. Verify the checksum before installation:

```text
sha256sum -c SHA256SUMS
stellar-replay version
stellar-replay help
```

On Windows, use the platform checksum utility or PowerShell's
`Get-FileHash` and verify the extracted `stellar-replay.exe` with
`stellar-replay.exe version`.

The release workflow is triggered only by a `v*.*.*` tag. It does not publish to
a package registry and does not sign artifacts. A release is not claimed until
the tagged workflow and its artifact checks pass.

## Compatibility

The core CLI version, fixture schema version, and release documentation are
related but distinct. The canonical fixture schema remains owned by this
repository. The GitHub Action has its own release version and must consume a
released core interface rather than an unreleased commit.
