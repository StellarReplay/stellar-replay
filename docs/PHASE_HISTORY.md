# Phase History

This index separates durable product truth from historical execution evidence.
Reports record what was verified at a checkpoint and must not be rewritten to
make later events appear earlier. The current roadmap is authoritative for
future work; this page is the navigable evidence index.

## Completed foundation, MVP, and release

| Phase | Milestone / era | Result | Evidence |
|---:|---|---|---|
| 0 | Foundation / constitution | Repository identity, topology, governance, and architecture established | [Phase 0](history/phases/PHASE-000.md) |
| 1 | Foundation / MVP freeze | Supported methods, protocol, fixture and network boundaries frozen | [Phase 1](history/phases/PHASE-001.md), [decision record](PHASE_1_DECISION.md) |
| 2 | MVP / fixture contract | Versioned validation, canonical hashing, and atomic storage implemented | [Phase 2](history/phases/PHASE-002.md), [schema](FIXTURE_SCHEMA.md) |
| 3 | MVP / capture | Bounded, explicit, read-only capture implemented | [Phase 3](history/phases/PHASE-003.md), [capture](CAPTURE.md) |
| 4 | MVP / replay | Exact deterministic matching and integrity-gated replay implemented | [Phase 4](history/phases/PHASE-004.md), [replay](REPLAY.md) |
| 5 | MVP / local server | Loopback JSON-RPC replay server implemented | [Phase 5](history/phases/PHASE-005.md), [server](SERVER.md) |
| 6 | MVP / CLI | `record`, `inspect`, `validate`, `replay`, and `serve` integrated | [Phase 6](history/phases/PHASE-006.md), [commands](COMMANDS.md) |
| 7 | MVP / deterministic workflow | Sanitized fixtures and offline integration proof added | [Phase 7](history/phases/PHASE-007.md) |
| 8 | MVP / security | Sanitization, bounds, secret controls, and security checklist completed | [Phase 8](history/phases/PHASE-008.md), [security checklist](SECURITY_CHECKLIST.md) |
| 9 | MVP / developer experience | Executable docs, demo, troubleshooting, and doc checks completed | [Phase 9](history/phases/PHASE-009.md) |
| 10 | Release / CI | Cross-platform CI, archives, checksums, and packaging completed | [Phase 10](history/phases/PHASE-010.md) |
| 11 | Release / audit | Conditional release GO and accepted limitations recorded | [Final audit](history/audits/FINAL_AUDIT.md) |
| 12 | Release / v0.1.0 | Tag, published release, checksums, artifact inspection, and smoke tests completed | [Phase 12](history/phases/PHASE-012.md), [release](RELEASES.md) |

## Historical interpretation

The final audit is intentionally pre-release evidence: it says no release had
yet been published at its audit checkpoint. Phase 12 is the later evidence that
the authorized release was actually published and verified. Root-level report
clutter was moved into `docs/history/phases/`; the numbered files retain their
historical content and Git history.

## Future evidence rule

Each future phase gets one `PHASE-NNN.md` report under `docs/history/phases/`.
Audits belong under `docs/history/audits/`. Release evidence belongs in the
phase report and `docs/RELEASES.md`; GitHub releases, commits, and workflow URLs
are linked rather than copied into many documents.
