# Study Run Summary

- Study: `ultraplan-daemon-events-study`
- Updated: `2026-09-03T19:11:58Z`
- Study progress state: `.ultraplan/run-state.json`
- Ledger: `.ultraplan/runs/tasks.jsonl`

## Status

| Metric | Value |
| --- | ---: |
| Runs recorded | 66 |
| Completed | 50 |
| Failed | 5 |
| Cancelled | 11 |
| Skipped | 0 |
| Remaining tasks | 0 |
| Dimensions seen | 10 |
| Sources seen | 4 |

## Remaining Work

No remaining tasks in the current run state.

## Dimensions

| Dimension | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 01.01-daemon-rpc-and-execution-ownership | 18 | 5 | 4 | 459m08s | 477767 | - |
| 01.02-operation-step-attempt-process-identity | 8 | 5 | 1 | 686m39s | 559475 | - |
| 01.03-durable-submission-and-idempotent-acceptance | 5 | 5 | 0 | 3m18s | 685210 | - |
| 01.04-leases-heartbeats-fencing-stale-worker-rejection | 5 | 5 | 0 | 2m34s | 564913 | - |
| 01.05-atomic-state-transition-and-event-journal | 5 | 5 | 0 | 5m10s | 600964 | - |
| 01.06-event-envelope-ordering-causation-and-replay | 5 | 5 | 0 | 15m44s | 505945 | - |
| 01.07-event-delivery-backpressure-and-retention-tiers | 5 | 5 | 0 | 1m38s | 340900 | - |
| 01.08-crash-recovery-reconciliation-and-checkpoints | 5 | 5 | 0 | 1m25s | 348071 | - |
| 01.09-cancellation-shutdown-and-process-cleanup | 5 | 5 | 0 | 1m47s | 348365 | - |
| 01.10-scheduler-controller-retries-and-at-least-once | 5 | 5 | 0 | 1m58s | 316570 | - |

## Sources

| Source | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| buildkit | 19 | 10 | 4 | 290m59s | 1080380 | - |
| containerd | 15 | 10 | 1 | 365m44s | 974775 | - |
| dagster | 12 | 10 | 0 | 229m56s | 1092364 | - |
| temporal | 10 | 10 | 0 | 10m02s | 951756 | - |

## Runtime And Model

| Runtime / Provider / Model | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| opencode / - / openrouter/stealth/ox-alpha | 5 | 0 | 0 | 2728m52s | - | - |
| opencode / opencode / muse-spark-1.2-contributor-free | 29 | 29 | 0 | 5m43s | 3394274 | - |
| opencode / opencode / muse-spark-1.3-contributor-free | 20 | 20 | 0 | 1m42s | 1353906 | - |
| openrouter / openrouter / stealth/ox-alpha | 12 | 1 | 5 | 6m45s | - | - |

## Recent Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-09-03T19:11:58Z | 01.10-scheduler-controller-retries-and-at-least-once | (synthesis) | synthesis | completed | 1m13s | muse-spark-1.3-contributor-free | 53067 | - |
| 2026-09-03T19:10:45Z | 01.10-scheduler-controller-retries-and-at-least-once | temporal | analysis | completed | 2m27s | muse-spark-1.3-contributor-free | 50110 | - |
| 2026-09-03T19:09:49Z | 01.10-scheduler-controller-retries-and-at-least-once | dagster | analysis | completed | 2m56s | muse-spark-1.3-contributor-free | 46965 | - |
| 2026-09-03T19:09:02Z | 01.09-cancellation-shutdown-and-process-cleanup | (synthesis) | synthesis | completed | 1m11s | muse-spark-1.3-contributor-free | 46722 | - |
| 2026-09-03T19:08:18Z | 01.10-scheduler-controller-retries-and-at-least-once | containerd | analysis | completed | 1m33s | muse-spark-1.3-contributor-free | 69457 | - |
| 2026-09-03T19:07:50Z | 01.09-cancellation-shutdown-and-process-cleanup | temporal | analysis | completed | 3m00s | muse-spark-1.3-contributor-free | 44751 | - |
| 2026-09-03T19:06:53Z | 01.10-scheduler-controller-retries-and-at-least-once | buildkit | analysis | completed | 1m40s | muse-spark-1.3-contributor-free | 96971 | - |
| 2026-09-03T19:06:44Z | 01.09-cancellation-shutdown-and-process-cleanup | dagster | analysis | completed | 2m00s | muse-spark-1.3-contributor-free | 95336 | - |
| 2026-09-03T19:05:13Z | 01.09-cancellation-shutdown-and-process-cleanup | containerd | analysis | completed | 1m30s | muse-spark-1.3-contributor-free | 76602 | - |
| 2026-09-03T19:04:51Z | 01.08-crash-recovery-reconciliation-and-checkpoints | (synthesis) | synthesis | completed | 1m17s | muse-spark-1.3-contributor-free | 44600 | - |
| 2026-09-03T19:04:44Z | 01.09-cancellation-shutdown-and-process-cleanup | buildkit | analysis | completed | 1m13s | muse-spark-1.3-contributor-free | 84954 | - |
| 2026-09-03T19:03:44Z | 01.07-event-delivery-backpressure-and-retention-tiers | (synthesis) | synthesis | completed | 1m00s | muse-spark-1.3-contributor-free | 46691 | - |
| 2026-09-03T19:03:34Z | 01.08-crash-recovery-reconciliation-and-checkpoints | temporal | analysis | completed | 1m09s | muse-spark-1.3-contributor-free | 71781 | - |
| 2026-09-03T19:03:31Z | 01.08-crash-recovery-reconciliation-and-checkpoints | dagster | analysis | completed | 1m28s | muse-spark-1.3-contributor-free | 66085 | - |
| 2026-09-03T19:02:44Z | 01.07-event-delivery-backpressure-and-retention-tiers | temporal | analysis | completed | 2m20s | muse-spark-1.3-contributor-free | 39493 | - |
| 2026-09-03T19:02:24Z | 01.08-crash-recovery-reconciliation-and-checkpoints | buildkit | analysis | completed | 1m48s | muse-spark-1.3-contributor-free | 91249 | - |
| 2026-09-03T19:02:03Z | 01.08-crash-recovery-reconciliation-and-checkpoints | containerd | analysis | completed | 1m20s | muse-spark-1.3-contributor-free | 74356 | - |
| 2026-09-03T19:00:43Z | 01.07-event-delivery-backpressure-and-retention-tiers | dagster | analysis | completed | 1m45s | muse-spark-1.3-contributor-free | 94527 | - |
| 2026-09-03T19:00:37Z | 01.07-event-delivery-backpressure-and-retention-tiers | containerd | analysis | completed | 1m39s | muse-spark-1.3-contributor-free | 72079 | - |
| 2026-09-03T19:00:24Z | 01.07-event-delivery-backpressure-and-retention-tiers | buildkit | analysis | completed | 1m27s | muse-spark-1.3-contributor-free | 88110 | - |

## Slowest Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | containerd | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | containerd | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | buildkit | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.01-daemon-rpc-and-execution-ownership | dagster | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-08-24T10:14:10Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | cancelled | 2728m52s | openrouter/stealth/ox-alpha | - | - |
| 2026-09-02T19:34:44Z | 01.06-event-envelope-ordering-causation-and-replay | temporal | analysis | completed | 45m41s | muse-spark-1.2-contributor-free | 119499 | - |
| 2026-08-22T12:26:16Z | 01.01-daemon-rpc-and-execution-ownership | temporal | analysis | completed | 32m48s | stealth/ox-alpha | - | - |
| 2026-09-02T20:00:59Z | 01.06-event-envelope-ordering-causation-and-replay | (synthesis) | synthesis | completed | 26m15s | muse-spark-1.2-contributor-free | 57535 | - |
| 2026-08-22T12:33:32Z | 01.02-operation-step-attempt-process-identity | buildkit | analysis | failed | 15m30s | stealth/ox-alpha | - | - |
| 2026-09-02T18:31:10Z | 01.05-atomic-state-transition-and-event-journal | buildkit | analysis | completed | 14m39s | muse-spark-1.2-contributor-free | 121020 | - |

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

