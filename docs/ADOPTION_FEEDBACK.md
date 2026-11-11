# Adoption feedback

This is a lightweight maintainer record for evidence gathered after the v0.1.1
release. It is not a substitute for GitHub issues, release notes, or private
user data. Do not record credentials, private fixture contents, or identifying
information here.

## Current evidence

As of 2026-11-10, the repository has a verified core v0.1.1 release, a
versioned Action v0.1.1, the independently published
`@stellar-replay/cli-adapter` 0.1.1 package, an extended documentation
repository, and a closed initial issue queue. No independently verified
external usage report or external contribution has yet been recorded in this
repository. The release and integration surfaces are ready for external
evaluation; adoption evidence remains an open productization gate.

## Maintainer evidence

| Evidence | Status | Meaning |
|---|---|---|
| Core release | Verified | `v0.1.1` archives, checksums, and platform smoke tests published |
| Action release | Verified | Action `v0.1.1` pins and verifies core `v0.1.1` |
| npm adapter | Verified | `@stellar-replay/cli-adapter@0.1.1` is independently installable |
| Downstream example | Verified | Official Stellar JS SDK consumes a checked-in replay fixture |
| External user report | Not yet recorded | No public, reproducible third-party workflow has been documented |
| External contribution | Not yet recorded | No independently authored contribution has been accepted yet |

## Feedback entry template

```text
Date:
Source: public issue / pull request / anonymized maintainer note
User problem:
Workflow attempted:
Observed result:
Requested change:
Security or compatibility impact:
Follow-up issue or decision:
```

Feedback becomes roadmap evidence only after it is reproducible, safely
sanitized, and linked to a public issue or decision record where appropriate.
