# Contributing

Stellar Replay is documentation-first and phase-driven. Read `PROJECT_CONTEXT.md`, `ROADMAP.md`, and the relevant design document before changing code.

Contributions should be narrow, tested, and honest about protocol limitations. Add deterministic local tests for behavior; keep live-network checks opt-in. Do not add transaction submission, secret handling, or new RPC methods without a roadmap decision and protocol/security documentation.

Before submitting a change, run the applicable Go formatting, test, vet, and diff checks. Update the changelog and phase checkpoint when a phase changes. Never commit secrets or captured data that is not safe to publish.

