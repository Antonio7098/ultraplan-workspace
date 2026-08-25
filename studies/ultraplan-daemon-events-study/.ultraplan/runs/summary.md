# Study Run Summary

- Study: `ultraplan-daemon-events-study`
- Updated: `2026-08-24T12:39:37Z`
- Study progress state: `.ultraplan/run-state.json`
- Ledger: `.ultraplan/runs/tasks.jsonl`

## Status

| Metric | Value |
| --- | ---: |
| Runs recorded | 17 |
| Completed | 1 |
| Failed | 5 |
| Cancelled | 11 |
| Skipped | 0 |
| Remaining tasks | 47 |
| Dimensions seen | 2 |
| Sources seen | 4 |

## Remaining Work

| Dimension | Source | Kind | Status |
| --- | --- | --- | --- |
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
| 01.01-daemon-rpc-and-execution-ownership | 14 | 1 | 4 | 589m26s | - | - |
| 01.02-operation-step-attempt-process-identity | 3 | 0 | 1 | 1824m25s | - | - |

## Sources

| Source | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| buildkit | 9 | 0 | 4 | 610m26s | - | - |
| containerd | 5 | 0 | 1 | 1092m43s | - | - |
| dagster | 2 | 0 | 0 | 1367m32s | - | - |
| temporal | 1 | 1 | 0 | 32m48s | - | - |

## Runtime And Model

| Runtime / Provider / Model | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| opencode / - / openrouter/stealth/ox-alpha | 5 | 0 | 0 | 2728m52s | - | - |
| openrouter / openrouter / stealth/ox-alpha | 12 | 1 | 5 | 6m45s | - | - |

## Recent Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-24T12:39:37Z | 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | cancelled | 1m56s | stealth/ox-alpha | - | - |
| 2026-08-24T12:37:41Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | failed | 1m37s | stealth/ox-alpha | - | - |
| 2026-08-24T12:27:07Z | 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | cancelled | 3m01s | stealth/ox-alpha | - | - |
| 2026-08-24T12:24:06Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | failed | 40s | stealth/ox-alpha | - | - |
| 2026-08-24T12:23:07Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 5m20s | stealth/ox-alpha | - | - |
| 2026-08-24T12:17:08Z | 01.01-daemon-rpc-and-execution-ownership | dagster | analysis | cancelled | 6m12s | stealth/ox-alpha | - | - |
| 2026-08-24T12:10:56Z | 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | failed | 51s | stealth/ox-alpha | - | - |
| 2026-08-24T12:10:04Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | failed | 1m11s | stealth/ox-alpha | - | - |
| 2026-08-24T10:46:16Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 3m02s | stealth/ox-alpha | - | - |
| 2026-08-24T10:38:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 8m54s | stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | containerd | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | dagster | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | failed | 15m30s | stealth/ox-alpha | - | - |
| 2026-08-22T12:26:16Z | 01.01-daemon-rpc-and-execution-ownership | temporal | analysis | completed | 32m48s | stealth/ox-alpha | - | - |

## Slowest Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | containerd | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | dagster | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-22T12:26:16Z | 01.01-daemon-rpc-and-execution-ownership | temporal | analysis | completed | 32m48s | stealth/ox-alpha | - | - |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | failed | 15m30s | stealth/ox-alpha | - | - |
| 2026-08-24T10:38:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 8m54s | stealth/ox-alpha | - | - |
| 2026-08-24T12:17:08Z | 01.01-daemon-rpc-and-execution-ownership | dagster | analysis | cancelled | 6m12s | stealth/ox-alpha | - | - |
| 2026-08-24T12:23:07Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 5m20s | stealth/ox-alpha | - | - |

## Failed Or Cancelled Runs

| Completed | Dimension | Source | Status | Error |
| --- | --- | --- | --- | --- |
| 2026-08-24T12:39:37Z | 01.01-daemon-rpc-and-execution-ownership | containerd | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:37:41Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:27:07Z | 01.01-daemon-rpc-and-execution-ownership | containerd | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:24:06Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:23:07Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:17:08Z | 01.01-daemon-rpc-and-execution-ownership | dagster | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:10:56Z | 01.01-daemon-rpc-and-execution-ownership | containerd | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-24T12:10:04Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-24T10:46:16Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-24T10:38:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | containerd | cancelled | task belonged to a stopped process; inspect durable study state before resuming |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | buildkit | cancelled | task belonged to a stopped process; inspect durable study state before resuming |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | dagster | cancelled | task belonged to a stopped process; inspect durable study state before resuming |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | containerd | cancelled | task belonged to a stopped process; inspect durable study state before resuming |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | cancelled | task belonged to a stopped process; inspect durable study state before resuming |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:timeout (OpenCode run timed out) → openrouter/stealth/ox-alpha:runtime_exit (Failed to execute statement)] |

