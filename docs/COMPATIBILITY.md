# Compatibility and deprecation policy

## Current compatibility line

Core `v0.1.1` defines fixture schema version `1` and the frozen four-method
read-only contract. Patch releases preserve the schema and CLI behavior except
for security fixes, documentation corrections, or fixes that reject invalid
data more clearly. The Action and npm adapter are versioned independently but
publish their supported core version explicitly.

| Component | Current version | Supported relationship |
|---|---:|---|
| Core CLI | `v0.1.1` | Owns schema, protocol, replay, and release binaries |
| GitHub Action | `v0.1.1` | Pins and verifies core `v0.1.1` archives |
| npm adapter | `0.1.1` | Invokes a locally installed core CLI `v0.1.1`; no second engine |
| Docs repository | current `main` | Documents the named core and Action versions |

## Compatibility rules

1. A fixture schema change requires a new `schemaVersion`, compatibility tests,
   migration guidance, and an explicit roadmap decision.
2. A core release must identify the supported Action and adapter versions.
3. The Action must consume a released core archive and verify its checksum; it
   must never depend on an unreleased commit.
4. The npm adapter may support more than one core version only after a matrix
   test demonstrates the behavior and the package documents that range.
5. Documentation changes must link to the core source of truth and must not
   silently redefine fixture or replay semantics.

## Deprecation policy

No v0.1 behavior is deprecated in `v0.1.1`. When a future release deprecates a
method, field, CLI flag, or integration contract, the change must include:

- the affected version and replacement behavior;
- a migration example and compatibility test;
- a minimum one-minor-release notice where practical;
- an update to the conformance matrix, changelog, and release notes; and
- a security review if the change affects capture, secrets, integrity, or
  network boundaries.

Breaking fixture or CLI changes belong in a minor or major compatibility
decision, not an unannounced patch release. `v1.0.0` remains a maturity gate,
not a calendar promise; it requires sustained compatibility and adoption
evidence as described in `ROADMAP.md`.
