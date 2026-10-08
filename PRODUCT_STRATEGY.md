# Stellar Replay Product Strategy

## Product problem

Stellar/Soroban tests that depend directly on live RPC inherit network failure,
provider drift, rate limits, retention windows, and changing chain state.
Stellar Replay turns a controlled read-only interaction into a reviewable,
versioned fixture that can be replayed deterministically offline.

## Primary user

The primary user is a Stellar/Soroban developer or CI maintainer who needs
repeatable RPC behavior while developing or testing an application. The user
values inspectable artifacts, deterministic failures, low setup cost, and a
clear boundary around captured data.

## Smallest useful product

The smallest useful product is the released core CLI: explicitly capture one of
the supported read-only RPC methods, validate and inspect the fixture, replay it
without live-network fallback, and expose it through a loopback local server.
That product is v0.1.0.

## Product principles

- deterministic before convenient;
- read-only and offline by default;
- fixture format and replay semantics have one source of truth;
- evidence is not chain truth;
- secure boundaries are explicit and testable;
- a small useful tool is preferable to a fragmented platform;
- future capability must be justified by real usage.

## Current maturity

The core is a **released pre-1.0 product**: the v0.1.0 MVP is implemented,
tested, packaged, and published. The organization as a whole has entered
**productization**, not maturity: the Action and docs companions are still
pre-release, and user adoption evidence is not yet established.

## Version meanings

- **v0.1.0 — released:** frozen read-only MVP, deterministic fixtures/replay,
  local server, CLI, security boundaries, CI, and platform archives.
- **v0.2 — planned only after evidence:** targeted usability or compatibility
  improvements demonstrated by users; no automatic expansion of method scope.
- **v0.3 — candidate:** selected integrations, richer fixture inspection or
  diffing, and additional read-only coverage if demand justifies them.
- **v1.0 — candidate:** stable public contracts and compatibility policy with
  evidence of sustained use, support expectations, and migration discipline.

These are direction markers, not release promises.

## Product eras

1. **Foundation and MVP (complete):** constitution, research, fixture contract,
   capture, replay, server, CLI, tests, security, and documentation.
2. **Release and productization (current):** stabilize v0.1, improve onboarding,
   establish the Action contract, and create useful adoption documentation.
3. **Ecosystem (candidate):** integrations, SDK adapters, conformance fixtures,
   or tooling only when distinct ownership and user demand exist.
4. **Trust and scale (candidate):** provenance, optional signing, performance,
   compatibility, and support practices proportional to real use.

## Non-goals

Stellar Replay is not a wallet, explorer, payment application, generic RPC
client, generic mocking framework, transaction signer/submission tool, Horizon
replacement, hosted replay service, or AI-driven test oracle. Stateful forks,
browser interfaces, and broad blockchain support are discovery-gated ideas, not
commitments.

## Decision test for future capability

Before accepting a future feature, ask: does it solve a demonstrated user
problem; strengthen deterministic replay; remain optional where possible; add a
real security or ownership boundary; justify its maintenance cost; and have a
clear verification artifact? If not, it remains speculative or is rejected.
