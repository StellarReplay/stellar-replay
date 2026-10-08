# Competitive Landscape and Research Notes

Research checked 2026-10-08 against official Stellar documentation and representative public projects.

## What exists

- Official Stellar RPC is a lightweight JSON-RPC access layer with a bounded recent-history window, not a Horizon replacement. Its current method catalog includes health, network, ledger, transaction, event, fee, simulation, and submission methods.
- Official SDKs, including JavaScript and Go, provide RPC clients. They are the ecosystem compatibility layer, not fixture systems.
- The official `stellar/stellar-rpc` repository is the server implementation and test/integration reference; Stellar Replay is not another RPC node.
- VCR-style tools record HTTP interactions into inspectable cassettes and replay them with configurable matching. They are general-purpose and do not understand Stellar method scope or fixture trust semantics.
- Soroban-focused forks and local networks can provide stateful chain-like environments. They are heavier than deterministic response replay and solve a different problem: executing or forking chain state.

## Wedge

Stellar Replay focuses on captured Stellar JSON-RPC request/response fixtures, canonical matching, a local RPC-compatible server, provenance limitations, and a short developer workflow. It does not claim general record/replay novelty, full chain simulation, or contract verification.

## Assumptions to validate

1. Developers need response-level fixtures before they need a local chain.
2. Exact parameters are sufficiently stable for useful replay.
3. Read-only methods cover a credible first workflow.
4. SDK clients can point at a local HTTP endpoint without adapter-specific behavior.
5. Sanitization can preserve matching for common fixtures.

## Research sources

- [Stellar RPC introduction](https://developers.stellar.org/docs/data/apis/rpc)
- [RPC method reference](https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods)
- [Stellar client SDKs](https://developers.stellar.org/docs/tools/sdks/client-sdks)
- [Stellar RPC repository](https://github.com/stellar/stellar-rpc)
- [VCR record/replay model](https://github.com/vcr/vcr)
- [Soroban fork example](https://github.com/lobotomoe/soroban-fork)

