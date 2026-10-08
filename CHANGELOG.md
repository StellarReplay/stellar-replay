# Changelog

## Unreleased

- Established the Stellar Replay project constitution, research record, frozen v0.1 scope, architecture, protocol, security model, and executable roadmap.
- Added the mandatory phase-completion gate requiring an approved commit, GitHub synchronization, hosted-check verification, and a phase report before a phase can be marked complete.
- Completed Phase 1 research and MVP freeze with current Stellar RPC request, response, matching, SDK, and capture-endpoint decisions.
- Implemented Phase 2 fixture schema v1, canonical JSON hashing, strict validation, sensitive-field rejection, bounded loading, and atomic storage with tests.
- Implemented Phase 3 controlled capture with bounded HTTPS transport, endpoint restrictions, timeout/size enforcement, no redirects or retries, response sanitization hook, atomic recording, and local transport tests.
- Implemented Phase 4 deterministic replay with integrity-gated fixture loading, canonical exact matching, request-ID echoing, safe errors, duplicate-key rejection, and concurrent-read tests.

