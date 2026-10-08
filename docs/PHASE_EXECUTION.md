# Phase Execution and GitHub Gate

Every phase is delivered as a reviewed checkpoint, not as an informal coding session.

## Required report

Before a phase can be marked complete, record:

- phase and final status;
- objective and exact scope completed;
- files/components changed;
- acceptance criterion result, one by one;
- local commands and results;
- security, determinism, and network-boundary checks;
- commit SHA and repository branch;
- GitHub remote and synchronized commit URL;
- GitHub Actions workflow/run URL and required check results;
- known limitations, deferred work, and rollback plan.

## Status rules

- `IN_PROGRESS`: implementation or verification is underway.
- `READY_FOR_REVIEW`: local work and checks are ready, but owner approval, commit, GitHub sync, or hosted evidence is pending.
- `COMPLETE`: acceptance criteria, local checks, approved commit, GitHub synchronization, and applicable hosted checks are verified.
- `BLOCKED`: a required external condition has prevented meaningful progress and is documented.

No phase may be marked complete because files exist or local tests pass alone. If the repository has no configured remote, report that fact and keep the phase `READY_FOR_REVIEW`; do not fabricate a GitHub result.

## Next-phase rule

The agent stops after the report. The next phase is executed only after the owner explicitly authorizes it. On authorization, the agent re-inspects the repository and GitHub state rather than relying on the previous conversation.
