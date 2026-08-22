# Study Run Summary

- Study: `ultraplan-daemon-events-study`
- Updated: `2026-08-22T12:45:18Z`
- Study progress state: `.ultraplan/run-state.json`
- Ledger: `.ultraplan/runs/tasks.jsonl`

## Status

| Metric | Value |
| --- | ---: |
| Runs recorded | 2 |
| Completed | 1 |
| Failed | 1 |
| Cancelled | 0 |
| Skipped | 0 |
| Remaining tasks | 49 |
| Dimensions seen | 2 |
| Sources seen | 2 |

## Remaining Work

| Dimension | Source | Kind | Status |
| --- | --- | --- | --- |
| 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | pending |
| 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | pending |
| 01.01-daemon-rpc-and-execution-ownership | dagster | analysis | pending |
| 01.02-operation-step-attempt-process-identity | buildkit | analysis | pending |
| 01.02-operation-step-attempt-process-identity | containerd | analysis | pending |
| 01.02-operation-step-attempt-process-identity | dagster | analysis | pending |
| 01.02-operation-step-attempt-process-identity | temporal | analysis | pending |
| 01.03-durable-submission-and-idempotent-acceptance | buildkit | analysis | pending |
| 01.03-durable-submission-and-idempotent-acceptance | containerd | analysis | pending |
| 01.03-durable-submission-and-idempotent-acceptance | dagster | analysis | pending |
| 01.03-durable-submission-and-idempotent-acceptance | temporal | analysis | pending |
| 01.04-leases-heartbeats-fencing-stale-worker-rejection | buildkit | analysis | pending |
| 01.04-leases-heartbeats-fencing-stale-worker-rejection | containerd | analysis | pending |
| 01.04-leases-heartbeats-fencing-stale-worker-rejection | dagster | analysis | pending |
| 01.04-leases-heartbeats-fencing-stale-worker-rejection | temporal | analysis | pending |
| 01.05-atomic-state-transition-and-event-journal | buildkit | analysis | pending |
| 01.05-atomic-state-transition-and-event-journal | containerd | analysis | pending |
| 01.05-atomic-state-transition-and-event-journal | dagster | analysis | pending |
| 01.05-atomic-state-transition-and-event-journal | temporal | analysis | pending |
| 01.06-event-envelope-ordering-causation-and-replay | buildkit | analysis | pending |
| 01.06-event-envelope-ordering-causation-and-replay | containerd | analysis | pending |
| 01.06-event-envelope-ordering-causation-and-replay | dagster | analysis | pending |
| 01.06-event-envelope-ordering-causation-and-replay | temporal | analysis | pending |
| 01.07-event-delivery-backpressure-and-retention-tiers | buildkit | analysis | pending |
| 01.07-event-delivery-backpressure-and-retention-tiers | containerd | analysis | pending |
| 01.07-event-delivery-backpressure-and-retention-tiers | dagster | analysis | pending |
| 01.07-event-delivery-backpressure-and-retention-tiers | temporal | analysis | pending |
| 01.08-crash-recovery-reconciliation-and-checkpoints | buildkit | analysis | pending |
| 01.08-crash-recovery-reconciliation-and-checkpoints | containerd | analysis | pending |
| 01.08-crash-recovery-reconciliation-and-checkpoints | dagster | analysis | pending |
| 01.08-crash-recovery-reconciliation-and-checkpoints | temporal | analysis | pending |
| 01.09-cancellation-shutdown-and-process-cleanup | buildkit | analysis | pending |
| 01.09-cancellation-shutdown-and-process-cleanup | containerd | analysis | pending |
| 01.09-cancellation-shutdown-and-process-cleanup | dagster | analysis | pending |
| 01.09-cancellation-shutdown-and-process-cleanup | temporal | analysis | pending |
| 01.10-scheduler-controller-retries-and-at-least-once | buildkit | analysis | pending |
| 01.10-scheduler-controller-retries-and-at-least-once | containerd | analysis | pending |
| 01.10-scheduler-controller-retries-and-at-least-once | dagster | analysis | pending |
| 01.10-scheduler-controller-retries-and-at-least-once | temporal | analysis | pending |
| 01.01-daemon-rpc-and-execution-ownership | (synthesis) | synthesis | pending |
| 01.02-operation-step-attempt-process-identity | (synthesis) | synthesis | pending |
| 01.03-durable-submission-and-idempotent-acceptance | (synthesis) | synthesis | pending |
| 01.04-leases-heartbeats-fencing-stale-worker-rejection | (synthesis) | synthesis | pending |
| 01.05-atomic-state-transition-and-event-journal | (synthesis) | synthesis | pending |
| 01.06-event-envelope-ordering-causation-and-replay | (synthesis) | synthesis | pending |
| 01.07-event-delivery-backpressure-and-retention-tiers | (synthesis) | synthesis | pending |
| 01.08-crash-recovery-reconciliation-and-checkpoints | (synthesis) | synthesis | pending |
| 01.09-cancellation-shutdown-and-process-cleanup | (synthesis) | synthesis | pending |
| 01.10-scheduler-controller-retries-and-at-least-once | (synthesis) | synthesis | pending |

## Dimensions

| Dimension | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 01.01-daemon-rpc-and-execution-ownership | 1 | 1 | 0 | 32m48s | - | - |
| 01.02-operation-step-attempt-process-identity | 1 | 0 | 1 | 15m30s | - | - |

## Sources

| Source | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| buildkit | 1 | 0 | 1 | 15m30s | - | - |
| temporal | 1 | 1 | 0 | 32m48s | - | - |

## Runtime And Model

| Runtime / Provider / Model | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| openrouter / openrouter / stealth/ox-alpha | 2 | 1 | 1 | 24m09s | - | - |

## Recent Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | failed | 15m30s | stealth/ox-alpha | - | - |
| 2026-08-22T12:26:16Z | 01.01-daemon-rpc-and-execution-ownership | temporal | analysis | completed | 32m48s | stealth/ox-alpha | - | - |

## Slowest Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-22T12:26:16Z | 01.01-daemon-rpc-and-execution-ownership | temporal | analysis | completed | 32m48s | stealth/ox-alpha | - | - |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | failed | 15m30s | stealth/ox-alpha | - | - |

## Failed Or Cancelled Runs

| Completed | Dimension | Source | Status | Error |
| --- | --- | --- | --- | --- |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:timeout (OpenCode run timed out) → openrouter/stealth/ox-alpha:runtime_exit (Failed to execute statement)] |

