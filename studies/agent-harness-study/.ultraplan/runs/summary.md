# Study Run Summary

- Study: `agent-harness-study`
- Updated: `2026-08-29T18:23:49Z`
- Study progress state: `.ultraplan/run-state.json`
- Ledger: `.ultraplan/runs/tasks.jsonl`

## Status

| Metric | Value |
| --- | ---: |
| Runs recorded | 1730 |
| Completed | 1482 |
| Failed | 113 |
| Cancelled | 135 |
| Skipped | 0 |
| Remaining tasks | 0 |
| Dimensions seen | 118 |
| Sources seen | 10 |

## Remaining Work

No remaining tasks in the current run state.

## Dimensions

| Dimension | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 01.01-execution-model-taxonomy | 11 | 11 | 0 | 6811m36s | - | - |
| 01.02-control-flow-ownership | 11 | 11 | 0 | 6811m32s | - | - |
| 01.03-step-turn-task-atomicity | 18 | 11 | 7 | 4169m19s | - | - |
| 01.04-termination-and-loop-bounds | 11 | 11 | 0 | 6820m48s | - | - |
| 01.05-pause-resume-interrupt-semantics | 8 | 8 | 0 | 9354m17s | - | - |
| 01.06-scheduling-and-trigger-semantics | 13 | 6 | 3 | 114m43s | - | - |
| 01.07-concurrency-and-parallel-advancement | 11 | 8 | 0 | 67m38s | - | - |
| 01.08-streaming-execution-semantics | 8 | 8 | 0 | 10m22s | - | - |
| 01.09-delivery-guarantees-idempotency | 7 | 5 | 2 | 7m02s | - | - |
| 01.10-replay-and-determinism | 3 | 3 | 0 | 10m58s | - | - |
| 02.01-state-taxonomy-and-ownership | 29 | 7 | 12 | 4m01s | - | - |
| 02.02-snapshot-and-checkpoint-architecture | 6 | 6 | 0 | 7m21s | - | - |
| 02.03-event-sourcing-and-replay-state | 3 | 3 | 0 | 12m05s | - | - |
| 02.04-mutation-discipline-and-state-transitions | 8 | 6 | 0 | 10m58s | - | - |
| 02.05-persistence-durability-tiers | 13 | 7 | 0 | 14m28s | - | - |
| 02.06-schema-versioning-and-migration | 10 | 6 | 0 | 6m56s | - | - |
| 02.07-session-thread-user-boundaries | 10 | 7 | 2 | 10m14s | - | - |
| 02.08-crash-recovery-and-reconstruction | 8 | 6 | 2 | 12m01s | - | - |
| 02.09-state-pruning-compaction-retention | 9 | 7 | 0 | 177m02s | 9182246 | 0.0000 USD |
| 03.01-llm-turn-loop-structure | 9 | 8 | 0 | 16m12s | - | - |
| 03.02-reason-act-observe-cadence | 9 | 8 | 1 | 124m02s | - | - |
| 03.03-tool-calling-roundtrip-control | 10 | 8 | 0 | 187m54s | - | - |
| 03.04-planner-executor-separation | 9 | 7 | 0 | 14m50s | - | - |
| 03.05-reflection-reask-self-correction | 8 | 8 | 0 | 104m27s | 13267339 | 0.0000 USD |
| 03.06-stuck-doom-loop-detection | 10 | 8 | 0 | 83m55s | 11559601 | 0.0000 USD |
| 03.07-context-refresh-inside-loop | 8 | 8 | 0 | 11m36s | 10665613 | 0.0000 USD |
| 03.08-subagent-forked-loop-design | 10 | 7 | 0 | 7m21s | 8669596 | 0.0000 USD |
| 03.09-completion-and-finalization-semantics | 9 | 8 | 1 | 10m30s | 23056032 | 0.0000 USD |
| 04.01-tool-definition-and-registration | 10 | 8 | 1 | 11m03s | - | - |
| 04.02-tool-schema-generation-validation | 12 | 9 | 1 | 8m37s | - | - |
| 04.03-tool-catalog-discovery-routing | 10 | 8 | 1 | 8m39s | 24001024 | 0.0000 USD |
| 04.04-tool-context-dependency-injection | 10 | 8 | 0 | 92m14s | 10124106 | 0.0000 USD |
| 04.05-tool-permissions-approval-metadata | 12 | 8 | 0 | 147m35s | - | - |
| 04.06-tool-result-contract-error-envelope | 10 | 9 | 1 | 8m09s | 345185 | 0.0000 USD |
| 04.07-external-tool-protocols-mcp | 10 | 9 | 0 | 300m40s | - | - |
| 04.08-agent-as-tool-composition | 14 | 14 | 0 | 1280m14s | - | - |
| 05.01-short-term-conversation-memory | 17 | 16 | 1 | 2249m02s | - | - |
| 05.02-working-memory-scratchpad | 17 | 16 | 1 | 2300m45s | - | - |
| 05.03-long-term-user-project-domain-memory | 11 | 10 | 1 | 2028m18s | - | - |
| 05.04-retrieval-augmented-memory | 15 | 14 | 1 | 2173m30s | - | - |
| 05.05-memory-write-policy | 11 | 10 | 1 | 2287m56s | - | - |
| 05.06-memory-compression-summarization | 12 | 12 | 0 | 2113m13s | - | - |
| 05.07-memory-privacy-scope-deletion | 11 | 10 | 1 | 2033m19s | - | - |
| 05.08-memory-evaluation-freshness | 4 | 4 | 0 | 1271m01s | - | - |
| 06.01-planning-location-responsibility | 21 | 16 | 1 | 2628m37s | - | - |
| 06.02-task-decomposition-representation | 33 | 14 | 1 | 1011m03s | - | - |
| 06.03-plan-lifecycle-and-revision | 20 | 18 | 0 | 4m00s | - | - |
| 06.04-planner-executor-contract | 14 | 14 | 0 | 1m27s | - | - |
| 06.05-objective-progress-tracking | 19 | 18 | 1 | 4m52s | - | - |
| 06.06-search-backtracking-alternatives | 10 | 10 | 0 | 1m26s | - | - |
| 06.07-plan-observability-evaluation | 14 | 14 | 0 | 1m30s | - | - |
| 07.01-tool-scheduling-and-dispatch | 19 | 18 | 1 | 4m18s | - | - |
| 07.02-sequential-vs-parallel-tool-execution | 16 | 16 | 0 | 6m05s | - | - |
| 07.03-idempotency-and-retry-semantics | 15 | 14 | 1 | 3m59s | - | - |
| 07.04-timeouts-and-cancellation | 16 | 16 | 0 | 994m27s | - | - |
| 07.05-resource-locking-and-isolation | 8 | 8 | 0 | 732m30s | - | - |
| 07.06-side-effect-ledger-transaction-boundaries | 8 | 8 | 0 | 6m37s | - | - |
| 07.07-tool-output-streaming | 17 | 16 | 1 | 6m50s | - | - |
| 07.08-tool-failure-escalation | 17 | 16 | 1 | 4m07s | - | - |
| 08.01-capability-model-trust-boundaries | 16 | 16 | 0 | 7m29s | - | - |
| 08.02-permission-policy-approval-gates | 14 | 14 | 0 | 5m45s | - | - |
| 08.03-secrets-identity-environment | 17 | 16 | 1 | 5m26s | - | - |
| 08.04-security-auditability | 16 | 16 | 0 | 5m58s | - | - |
| 09.01-policy-injection-points | 12 | 12 | 0 | 4m25s | - | - |
| 09.02-risk-taxonomy-control-mapping | 25 | 12 | 3 | 4m10s | - | - |
| 09.03-governance-ux-operator-workflow | 24 | 12 | 2 | 3m02s | - | - |
| 09.04-governance-evidence-generation | 27 | 12 | 9 | 2m28s | - | - |
| 10.01-span-hierarchy-run-tree | 29 | 20 | 8 | 45m08s | - | - |
| 10.02-event-schema-lifecycle-events | 23 | 22 | 1 | 791m53s | - | - |
| 10.03-causal-links-lineage | 21 | 20 | 1 | 6m11s | - | - |
| 10.04-export-interoperability-observability | 6 | 6 | 0 | 2m43s | - | - |
| 11.01-context-selection-policy | 16 | 16 | 0 | 4m39s | - | - |
| 11.02-token-budgeting-and-compression | 16 | 16 | 0 | 4m39s | - | - |
| 11.03-repository-workspace-context-maps | 8 | 8 | 0 | 4m42s | - | - |
| 11.04-context-provenance-integrity | 21 | 20 | 1 | 5m50s | - | - |
| 12.01-prompt-storage-versioning | 19 | 18 | 1 | 3m49s | - | - |
| 12.02-prompt-templating-variable-contracts | 21 | 18 | 1 | 5m53s | - | - |
| 12.03-prompt-evaluation-experiments | 8 | 8 | 0 | 2m24s | 960072 | - |
| 12.04-prompt-rollback-change-control | 3 | 3 | 0 | 1m46s | 392198 | - |
| 13.01-error-taxonomy | 22 | 22 | 0 | 13m03s | - | - |
| 13.02-retry-fallback-degraded-mode | 23 | 22 | 1 | 6m18s | - | - |
| 13.03-failure-visibility | 25 | 23 | 2 | 6m44s | - | - |
| 13.04-recovery-vs-escalation | 21 | 20 | 1 | 6m35s | - | - |
| 14.01-human-in-the-loop-trigger-policy | 12 | 12 | 0 | 6m23s | - | - |
| 14.02-approval-session-design | 12 | 12 | 0 | 4m13s | - | - |
| 14.03-human-intervention-takeover | 22 | 16 | 0 | 5m42s | - | - |
| 15.01-coordination-topology | 15 | 14 | 1 | 4m07s | - | - |
| 15.02-message-routing-termination | 15 | 14 | 1 | 6m52s | - | - |
| 15.03-shared-state-conflict-resolution | 14 | 14 | 0 | 7m11s | - | - |
| 16.01-artifact-lifecycle | 7 | 5 | 0 | 5m18s | - | - |
| 16.02-diff-review-rollback | 4 | 4 | 0 | 3m03s | - | - |
| 16.03-artifact-provenance-reproducibility | 5 | 5 | 0 | 2m46s | - | - |
| 17.01-sandbox-boundary | 8 | 8 | 0 | 4m26s | - | - |
| 17.02-filesystem-network-process-controls | 11 | 10 | 1 | 4m00s | - | - |
| 17.03-isolation-observability-escape-handling | 6 | 5 | 0 | 2m05s | 495267 | - |
| 18.01-dataset-golden-task-management | 19 | 18 | 1 | 4m24s | - | - |
| 18.02-trajectory-evaluation | 9 | 9 | 0 | 1m57s | 1060258 | - |
| 18.03-regression-gating-ci-integration | 15 | 14 | 1 | 1043m04s | - | - |
| 18.04-cost-latency-quality-evaluation | 15 | 14 | 1 | 3m19s | - | - |
| 19.01-protocol-compatibility | 33 | 22 | 8 | 2m32s | - | - |
| 19.02-portable-trace-eval-prompt-schemas | 7 | 7 | 0 | 12m58s | - | - |
| 19.03-adapter-interop-boundary-design | 33 | 22 | 1 | 4m36s | - | - |
| 20.01-token-cost-accounting | 17 | 16 | 1 | 5m06s | - | - |
| 20.02-caching-batching-reuse | 17 | 16 | 1 | 5m07s | - | - |
| 20.03-quality-cost-routing | 27 | 16 | 9 | 3m25s | - | - |
| 21.01-plugin-extension-points | 11 | 11 | 0 | 2m45s | - | - |
| 21.02-provider-backend-adapters | 23 | 22 | 1 | 5m39s | - | - |
| 21.03-extension-compatibility-testing | 26 | 22 | 0 | 4m44s | - | - |
| 22.01-package-module-boundaries | 22 | 22 | 0 | 31m56s | - | - |
| 22.02-configuration-deployment-shape | 22 | 22 | 0 | 5m23s | - | - |
| 22.03-docs-examples-contributor-workflow | 23 | 22 | 1 | 4m17s | - | - |
| 23.01-autonomy-boundary | 17 | 16 | 1 | 3m37s | - | - |
| 23.02-persistence-vs-escalation-philosophy | 17 | 16 | 1 | 6m41s | - | - |
| 23.03-responsibility-accountability-model | 8 | 8 | 0 | 2m18s | 1128066 | - |
| 24.01-public-api-surface | 22 | 22 | 0 | 9m56s | - | - |
| 24.02-interface-contract-design | 23 | 22 | 1 | 71m14s | - | - |
| 24.03-api-versioning-compatibility | 22 | 22 | 0 | 1m56s | - | - |
| 24.04-embedding-and-host-integration-ergonomics | 26 | 22 | 3 | 5m47s | - | - |

## Sources

| Source | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| agent-framework | 223 | 186 | 11 | 300m10s | 16196051 | 0.0000 USD |
| crewai | 161 | 139 | 5 | 401m13s | 320501 | - |
| langfuse | 97 | 79 | 8 | 32m58s | 9666408 | 0.0000 USD |
| langgraph | 169 | 154 | 5 | 293m14s | 36374024 | 0.0000 USD |
| letta | 135 | 114 | 8 | 447m48s | 12479429 | 0.0000 USD |
| opa | 71 | 62 | 3 | 78m29s | 143511 | - |
| openai-agents-sdk | 160 | 145 | 7 | 229m51s | 19841523 | 0.0000 USD |
| openhands | 211 | 183 | 7 | 307m51s | 776818 | 0.0000 USD |
| pydantic-ai | 133 | 123 | 2 | 247m21s | 18236003 | 0.0000 USD |
| temporal | 107 | 95 | 1 | 91m45s | 67313 | - |

## Runtime And Model

| Runtime / Provider / Model | Runs | Completed | Failed | Avg Duration | Tokens | Cost |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| minimax-coding-plan / minimax-coding-plan / MiniMax-M3 | 390 | 314 | 24 | 995m05s | 110870742 | 0.0000 USD |
| opencode / - / minimax-coding-plan/MiniMax-M3 | 55 | 1 | 49 | 5s | - | - |
| opencode / - / openrouter/stealth/ox-alpha | 452 | 449 | 3 | 11m33s | - | - |
| opencode / opencode / muse-spark-1.2-contributor-free | 215 | 178 | 0 | 6m27s | 4035861 | - |
| openrouter / openrouter / minimax/minimax-m3:free | 2 | 1 | 0 | 8m10s | - | - |
| openrouter / openrouter / stealth/ox-alpha | 616 | 539 | 37 | 612m41s | - | - |

## Recent Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-29T10:52:01Z | 18.02-trajectory-evaluation | (synthesis) | synthesis | completed | 2m50s | muse-spark-1.2-contributor-free | 110588 | - |
| 2026-08-29T10:49:11Z | 18.02-trajectory-evaluation | pydantic-ai | analysis | completed | 1m37s | muse-spark-1.2-contributor-free | 124720 | - |
| 2026-08-29T10:47:34Z | 18.02-trajectory-evaluation | openhands | analysis | completed | 1m21s | muse-spark-1.2-contributor-free | 93467 | - |
| 2026-08-29T10:46:13Z | 18.02-trajectory-evaluation | openai-agents-sdk | analysis | completed | 1m57s | muse-spark-1.2-contributor-free | 144389 | - |
| 2026-08-29T10:44:15Z | 18.02-trajectory-evaluation | letta | analysis | completed | 1m49s | muse-spark-1.2-contributor-free | 102002 | - |
| 2026-08-29T10:42:24Z | 18.02-trajectory-evaluation | langgraph | analysis | completed | 3m22s | muse-spark-1.2-contributor-free | 120288 | - |
| 2026-08-29T10:39:02Z | 18.02-trajectory-evaluation | langfuse | analysis | completed | 1m39s | muse-spark-1.2-contributor-free | 147617 | - |
| 2026-08-29T10:37:22Z | 18.02-trajectory-evaluation | crewai | analysis | completed | 1m17s | muse-spark-1.2-contributor-free | 91692 | - |
| 2026-08-29T10:36:04Z | 18.02-trajectory-evaluation | agent-framework | analysis | completed | 1m38s | muse-spark-1.2-contributor-free | 125495 | - |
| 2026-08-29T10:34:26Z | 12.04-prompt-rollback-change-control | (synthesis) | synthesis | completed | 1m38s | muse-spark-1.2-contributor-free | 76278 | - |
| 2026-08-29T10:32:46Z | 12.04-prompt-rollback-change-control | langfuse | analysis | completed | 2m16s | muse-spark-1.2-contributor-free | 153233 | - |
| 2026-08-29T10:30:29Z | 12.04-prompt-rollback-change-control | agent-framework | analysis | completed | 1m23s | muse-spark-1.2-contributor-free | 162687 | - |
| 2026-08-29T10:29:05Z | 12.03-prompt-evaluation-experiments | (synthesis) | synthesis | completed | 2m06s | muse-spark-1.2-contributor-free | 71821 | - |
| 2026-08-29T10:26:58Z | 12.03-prompt-evaluation-experiments | pydantic-ai | analysis | completed | 1m20s | muse-spark-1.2-contributor-free | 169017 | - |
| 2026-08-29T10:25:37Z | 12.03-prompt-evaluation-experiments | openhands | analysis | completed | 1m59s | muse-spark-1.2-contributor-free | 95614 | - |
| 2026-08-29T10:23:37Z | 12.03-prompt-evaluation-experiments | openai-agents-sdk | analysis | completed | 1m25s | muse-spark-1.2-contributor-free | 97642 | - |
| 2026-08-29T10:22:10Z | 12.03-prompt-evaluation-experiments | langgraph | analysis | completed | 1m34s | muse-spark-1.2-contributor-free | 112571 | - |
| 2026-08-29T10:20:35Z | 12.03-prompt-evaluation-experiments | langfuse | analysis | completed | 2m26s | muse-spark-1.2-contributor-free | 183312 | - |
| 2026-08-29T10:18:08Z | 12.03-prompt-evaluation-experiments | crewai | analysis | completed | 6m04s | muse-spark-1.2-contributor-free | 99166 | - |
| 2026-08-29T10:12:04Z | 12.03-prompt-evaluation-experiments | agent-framework | analysis | completed | 2m13s | muse-spark-1.2-contributor-free | 130929 | - |

## Slowest Runs

| Completed | Dimension | Source | Kind | Status | Duration | Model | Tokens | Cost |
| --- | --- | --- | --- | --- | ---: | --- | ---: | ---: |
| 2026-08-23T17:39:50Z | 01.01-execution-model-taxonomy | (synthesis) | synthesis | completed | 74865m00s | MiniMax-M3 | - | - |
| 2026-08-23T18:59:17Z | 01.02-control-flow-ownership | (synthesis) | synthesis | completed | 74812m19s | MiniMax-M3 | - | - |
| 2026-08-23T18:22:27Z | 01.05-pause-resume-interrupt-semantics | (synthesis) | synthesis | completed | 74775m29s | MiniMax-M3 | - | - |
| 2026-08-23T18:08:32Z | 01.04-termination-and-loop-bounds | (synthesis) | synthesis | completed | 74761m34s | MiniMax-M3 | - | - |
| 2026-08-23T17:59:49Z | 01.03-step-turn-task-atomicity | (synthesis) | synthesis | completed | 74677m54s | MiniMax-M3 | - | - |
| 2026-08-27T10:46:15Z | 06.02-task-decomposition-representation | openai-agents-sdk | analysis | completed | 8409m13s | stealth/ox-alpha | - | - |
| 2026-08-27T10:40:51Z | 06.02-task-decomposition-representation | crewai | analysis | completed | 8403m49s | stealth/ox-alpha | - | - |
| 2026-08-27T10:38:06Z | 06.02-task-decomposition-representation | agent-framework | analysis | completed | 8401m04s | stealth/ox-alpha | - | - |
| 2026-08-27T10:25:40Z | 06.01-planning-location-responsibility | openhands | analysis | completed | 8388m38s | stealth/ox-alpha | - | - |
| 2026-08-27T00:19:30Z | 06.01-planning-location-responsibility | agent-framework | analysis | completed | 7789m04s | stealth/ox-alpha | - | - |

## Failed Or Cancelled Runs

| Completed | Dimension | Source | Status | Error |
| --- | --- | --- | --- | --- |
| 2026-08-28T17:00:43Z | 17.03-isolation-observability-escape-handling | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T15:08:38Z | 16.01-artifact-lifecycle | langgraph | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/minimax/minimax-m3:free] |
| 2026-08-27T14:27:51Z | 16.01-artifact-lifecycle | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 06.01-planning-location-responsibility | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 06.02-task-decomposition-representation | crewai | cancelled | cancellation: persist runtime event: context canceled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 06.02-task-decomposition-representation | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 06.03-plan-lifecycle-and-revision | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 06.02-task-decomposition-representation | agent-framework | cancelled | cancellation: persist runtime event: context canceled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 06.02-task-decomposition-representation | pydantic-ai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T09:59:12Z | 21.03-extension-compatibility-testing | (synthesis) | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:45Z | 06.02-task-decomposition-representation | pydantic-ai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:45Z | 06.02-task-decomposition-representation | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:45Z | 06.01-planning-location-responsibility | openhands | cancelled | cancellation: persist runtime event: context canceled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:45Z | 06.02-task-decomposition-representation | openai-agents-sdk | cancelled | cancellation: persist runtime event: context canceled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:44Z | 21.03-extension-compatibility-testing | (synthesis) | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:44Z | 06.02-task-decomposition-representation | crewai | cancelled | cancellation: persist runtime event: context canceled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:46:43Z | 06.02-task-decomposition-representation | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 06.02-task-decomposition-representation | pydantic-ai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 06.02-task-decomposition-representation | agent-framework | cancelled | cancellation: persist runtime event: context canceled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 06.02-task-decomposition-representation | openai-agents-sdk | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 21.03-extension-compatibility-testing | (synthesis) | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 06.02-task-decomposition-representation | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 06.01-planning-location-responsibility | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:34:11Z | 06.02-task-decomposition-representation | crewai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 06.02-task-decomposition-representation | pydantic-ai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 21.03-extension-compatibility-testing | (synthesis) | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 06.02-task-decomposition-representation | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 06.03-plan-lifecycle-and-revision | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 06.02-task-decomposition-representation | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 06.01-planning-location-responsibility | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-27T08:14:09Z | 06.02-task-decomposition-representation | crewai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 20.03-quality-cost-routing | (synthesis) | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 19.03-adapter-interop-boundary-design | langgraph | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 19.01-protocol-compatibility | pydantic-ai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 19.03-adapter-interop-boundary-design | letta | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 19.03-adapter-interop-boundary-design | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 19.03-adapter-interop-boundary-design | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T18:54:46Z | 19.03-adapter-interop-boundary-design | crewai | cancelled | cancellation: validation: cancellation: validation run was cancelled [opencode/muse-spark-1.2-contributor-free] |
| 2026-08-26T14:13:08Z | 19.03-adapter-interop-boundary-design | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:13:07Z | 19.03-adapter-interop-boundary-design | letta | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:13:03Z | 19.03-adapter-interop-boundary-design | crewai | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:13:03Z | 19.03-adapter-interop-boundary-design | langgraph | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:13:03Z | 19.03-adapter-interop-boundary-design | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:13:03Z | 19.01-protocol-compatibility | pydantic-ai | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:13:02Z | 19.03-adapter-interop-boundary-design | opa | cancelled | context canceled |
| 2026-08-26T14:13:02Z | 19.01-protocol-compatibility | temporal | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:12:58Z | 19.01-protocol-compatibility | openhands | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:12:12Z | 19.01-protocol-compatibility | openai-agents-sdk | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:12:12Z | 19.01-protocol-compatibility | opa | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:12:11Z | 19.01-protocol-compatibility | letta | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:12:11Z | 19.01-protocol-compatibility | langgraph | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:12:03Z | 19.01-protocol-compatibility | langfuse | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:11:46Z | 19.01-protocol-compatibility | crewai | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:11:32Z | 19.01-protocol-compatibility | agent-framework | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:31Z | 20.03-quality-cost-routing | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-26T14:10:31Z | 20.03-quality-cost-routing | openai-agents-sdk | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:31Z | 20.03-quality-cost-routing | langgraph | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:30Z | 20.03-quality-cost-routing | openhands | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:30Z | 20.03-quality-cost-routing | pydantic-ai | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:29Z | 20.03-quality-cost-routing | langfuse | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:18Z | 20.03-quality-cost-routing | crewai | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:10:12Z | 20.03-quality-cost-routing | agent-framework | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:33Z | 09.02-risk-taxonomy-control-mapping | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-26T14:08:33Z | 09.02-risk-taxonomy-control-mapping | langfuse | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:33Z | 09.04-governance-evidence-generation | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-26T14:08:31Z | 09.03-governance-ux-operator-workflow | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-26T14:08:31Z | 09.04-governance-evidence-generation | agent-framework | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:29Z | 09.04-governance-evidence-generation | openai-agents-sdk | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:29Z | 09.03-governance-ux-operator-workflow | openhands | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:28Z | 09.04-governance-evidence-generation | langfuse | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:25Z | 09.04-governance-evidence-generation | opa | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:08:22Z | 09.04-governance-evidence-generation | openhands | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.) → openrouter/stealth/ox-alpha:runtime_exit (Thank you for participating in the Stealth Ox Alpha testing period. This model will be revealed today, August 26th.)] |
| 2026-08-26T14:06:58Z | 09.02-risk-taxonomy-control-mapping | langfuse | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:06:56Z | 20.03-quality-cost-routing | agent-framework | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:06:55Z | 09.04-governance-evidence-generation | langfuse | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:06:54Z | 09.04-governance-evidence-generation | opa | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:06:51Z | 09.04-governance-evidence-generation | openhands | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T14:06:48Z | 09.04-governance-evidence-generation | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:11Z | 09.04-governance-evidence-generation | agent-framework | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:11Z | 09.02-risk-taxonomy-control-mapping | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:11Z | 09.03-governance-ux-operator-workflow | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:11Z | 09.04-governance-evidence-generation | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:11Z | 09.04-governance-evidence-generation | openai-agents-sdk | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:11Z | 09.04-governance-evidence-generation | opa | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:22:07Z | 20.03-quality-cost-routing | agent-framework | cancelled | context canceled |
| 2026-08-26T13:22:06Z | 09.04-governance-evidence-generation | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.03-governance-ux-operator-workflow | opa | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.03-governance-ux-operator-workflow | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.02-risk-taxonomy-control-mapping | openhands | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.03-governance-ux-operator-workflow | openai-agents-sdk | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.03-governance-ux-operator-workflow | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.02-risk-taxonomy-control-mapping | langfuse | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.03-governance-ux-operator-workflow | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:10:36Z | 09.02-risk-taxonomy-control-mapping | openai-agents-sdk | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:08:59Z | 09.02-risk-taxonomy-control-mapping | openai-agents-sdk | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:08:59Z | 09.03-governance-ux-operator-workflow | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:08:59Z | 09.02-risk-taxonomy-control-mapping | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:08:59Z | 09.03-governance-ux-operator-workflow | agent-framework | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:08:54Z | 09.02-risk-taxonomy-control-mapping | openhands | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:08:54Z | 09.03-governance-ux-operator-workflow | opa | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:04:00Z | 09.02-risk-taxonomy-control-mapping | langfuse | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T13:03:56Z | 09.02-risk-taxonomy-control-mapping | openhands | cancelled | context canceled |
| 2026-08-26T13:03:56Z | 09.03-governance-ux-operator-workflow | agent-framework | cancelled | context canceled |
| 2026-08-26T13:03:55Z | 09.02-risk-taxonomy-control-mapping | openai-agents-sdk | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T11:43:42Z | 14.03-human-intervention-takeover | crewai | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T11:43:42Z | 14.03-human-intervention-takeover | agent-framework | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T11:40:32Z | 14.03-human-intervention-takeover | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T11:40:32Z | 14.03-human-intervention-takeover | crewai | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-26T11:36:24Z | 14.03-human-intervention-takeover | agent-framework | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-26T11:36:24Z | 14.03-human-intervention-takeover | crewai | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-25T22:16:03Z | 12.02-prompt-templating-variable-contracts | langgraph | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-25T22:16:03Z | 12.02-prompt-templating-variable-contracts | langfuse | cancelled | cancellation: persist runtime event: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-25T13:08:50Z | 10.01-span-hierarchy-run-tree | langfuse | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:timeout (OpenCode run timed out) → openrouter/stealth/ox-alpha:runtime_exit ([Stealth] stealth/ox-alpha is temporarily rate-limited upstream. Please retry shortly.)] |
| 2026-08-25T12:54:46Z | 10.01-span-hierarchy-run-tree | agent-framework | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit ([Stealth] stealth/ox-alpha is temporarily rate-limited upstream. Please retry shortly.) → openrouter/stealth/ox-alpha:runtime_exit ({"code":429,"message":"Provider returned error","metadata":{"error_type":"rate_limit_exceeded"}})] |
| 2026-08-25T12:52:51Z | 24.04-embedding-and-host-integration-ergonomics | (synthesis) | failed | opencode event: runtime_exit: OpenCode reported a fatal session error [openrouter/stealth/ox-alpha:runtime_exit ([Stealth] stealth/ox-alpha is temporarily rate-limited upstream. Please retry shortly.) → openrouter/stealth/ox-alpha:runtime_exit ([Stealth] stealth/ox-alpha is temporarily rate-limited upstream. Please retry shortly.)] |
| 2026-08-25T12:51:15Z | 10.01-span-hierarchy-run-tree | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T10:09:07Z | 24.04-embedding-and-host-integration-ergonomics | opa | cancelled | cancellation: validation: cancellation: validation run was cancelled [openrouter/stealth/ox-alpha] |
| 2026-08-25T10:09:06Z | 10.01-span-hierarchy-run-tree | crewai | cancelled | context canceled |
| 2026-08-25T10:09:06Z | 10.01-span-hierarchy-run-tree | agent-framework | failed | durable run acceptance failed: context canceled [openrouter/stealth/ox-alpha] |
| 2026-08-25T10:08:59Z | 10.01-span-hierarchy-run-tree | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T10:08:59Z | 13.03-failure-visibility | openai-agents-sdk | failed | opencode run: runtime_exit: OpenCode exited before a successful final result [openrouter/stealth/ox-alpha:runtime_exit (exit_code=1 stderr=[91m[1mError: [0mUnexpected error  unable to open database file ) → openrouter/stealth/ox-alpha:runtime_exit (exit_code=1 stderr=[91m[1mError: [0mUnexpected error  unable to open database file )] |
| 2026-08-25T10:08:51Z | 24.04-embedding-and-host-integration-ergonomics | agent-framework | failed | durable run acceptance failed: quota_read: inspect run-control storage usage failed [openrouter/stealth/ox-alpha] |
| 2026-08-25T10:08:45Z | 24.04-embedding-and-host-integration-ergonomics | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T10:08:45Z | 13.03-failure-visibility | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T10:05:22Z | 15.02-message-routing-termination | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T10:00:04Z | 06.05-objective-progress-tracking | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T08:50:44Z | 15.01-coordination-topology | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T05:44:40Z | 05.07-memory-privacy-scope-deletion | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T05:30:54Z | 05.05-memory-write-policy | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T05:16:59Z | 05.03-long-term-user-project-domain-memory | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T04:55:53Z | 05.04-retrieval-augmented-memory | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T04:25:23Z | 05.02-working-memory-scratchpad | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T03:37:29Z | 20.02-caching-batching-reuse | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T01:46:46Z | 05.01-short-term-conversation-memory | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-25T00:32:02Z | 11.04-context-provenance-integrity | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T23:52:03Z | 23.02-persistence-vs-escalation-philosophy | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T22:29:04Z | 23.01-autonomy-boundary | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T20:28:26Z | 08.03-secrets-identity-environment | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T19:56:54Z | 17.02-filesystem-network-process-controls | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T18:12:44Z | 07.01-tool-scheduling-and-dispatch | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T17:29:20Z | 20.01-token-cost-accounting | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T17:10:02Z | 07.03-idempotency-and-retry-semantics | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T17:08:10Z | 18.01-dataset-golden-task-management | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T17:05:45Z | 12.01-prompt-storage-versioning | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T17:03:17Z | 12.02-prompt-templating-variable-contracts | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T17:00:52Z | 10.03-causal-links-lineage | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:58:17Z | 21.02-provider-backend-adapters | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:55:11Z | 18.04-cost-latency-quality-evaluation | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:53:19Z | 07.07-tool-output-streaming | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:51:03Z | 13.02-retry-fallback-degraded-mode | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:47:59Z | 13.04-recovery-vs-escalation | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:45:17Z | 07.08-tool-failure-escalation | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T16:34:26Z | 22.03-docs-examples-contributor-workflow | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-24T15:19:56Z | 10.01-span-hierarchy-run-tree | openai-agents-sdk | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-24T14:58:24Z | 10.01-span-hierarchy-run-tree | langfuse | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-24T14:47:34Z | 10.01-span-hierarchy-run-tree | crewai | failed | repair_exhausted: validation repair: repair_exhausted: validation repair attempts exhausted [openrouter/stealth/ox-alpha] |
| 2026-08-23T13:37:17Z | 06.02-task-decomposition-representation | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-23T13:31:44Z | 06.01-planning-location-responsibility | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-21T16:53:39Z | 18.03-regression-gating-ci-integration | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-21T16:52:47Z | 10.02-event-schema-lifecycle-events | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-21T16:51:02Z | 24.02-interface-contract-design | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-08-15T15:53:10Z | 04.07-external-tool-protocols-mcp | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-08-15T15:29:18Z | 04.06-tool-result-contract-error-envelope | openai-agents-sdk | failed | validation_failed |
| 2026-07-28T08:56:28Z | 04.05-tool-permissions-approval-metadata | crewai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-28T08:56:25Z | 04.05-tool-permissions-approval-metadata | langgraph | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-28T08:56:25Z | 04.04-tool-context-dependency-injection | langgraph | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-28T08:48:45Z | 04.04-tool-context-dependency-injection | langgraph | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-28T08:48:45Z | 04.05-tool-permissions-approval-metadata | langgraph | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-28T08:48:45Z | 04.05-tool-permissions-approval-metadata | crewai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-22T19:08:30Z | 04.03-tool-catalog-discovery-routing | crewai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-22T19:07:19Z | 04.03-tool-catalog-discovery-routing | langgraph | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-19T22:22:23Z | 04.02-tool-schema-generation-validation | opa | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-19T22:22:23Z | 04.01-tool-definition-and-registration | pydantic-ai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-19T22:22:23Z | 04.02-tool-schema-generation-validation | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-19T20:47:08Z | 04.02-tool-schema-generation-validation | agent-framework | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-19T18:38:09Z | 04.01-tool-definition-and-registration | agent-framework | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-19T18:38:09Z | 03.09-completion-and-finalization-semantics | pydantic-ai | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-15T11:11:32Z | 03.08-subagent-forked-loop-design | openai-agents-sdk | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-15T11:11:32Z | 03.08-subagent-forked-loop-design | langgraph | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-15T11:11:32Z | 03.08-subagent-forked-loop-design | crewai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-15T09:38:16Z | 03.06-stuck-doom-loop-detection | pydantic-ai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-15T08:58:47Z | 03.06-stuck-doom-loop-detection | crewai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-14T16:07:00Z | 03.03-tool-calling-roundtrip-control | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-14T16:07:00Z | 03.04-planner-executor-separation | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-14T15:59:38Z | 03.03-tool-calling-roundtrip-control | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-14T15:59:38Z | 03.04-planner-executor-separation | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-14T14:31:13Z | 03.02-reason-act-observe-cadence | (synthesis) | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-13T18:45:17Z | 02.07-session-thread-user-boundaries | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-13T18:45:17Z | 02.09-state-pruning-compaction-retention | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-12T18:22:05Z | 02.09-state-pruning-compaction-retention | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-12T18:22:01Z | 03.01-llm-turn-loop-structure | crewai | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-12T16:30:06Z | 02.07-session-thread-user-boundaries | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-11T10:11:04Z | 02.08-crash-recovery-and-reconstruction | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-10T20:40:00Z | 02.08-crash-recovery-and-reconstruction | openhands | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-10T19:54:25Z | 02.07-session-thread-user-boundaries | letta | failed | opencode run: runtime_exit: OpenCode exited before a successful final result |
| 2026-07-10T16:21:57Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-10T16:21:57Z | 02.01-state-taxonomy-and-ownership | letta | failed | validation_failed |
| 2026-07-10T14:47:30Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:47:30Z | 02.06-schema-versioning-and-migration | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:47:30Z | 02.05-persistence-durability-tiers | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:43:56Z | 02.06-schema-versioning-and-migration | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:43:56Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:43:56Z | 02.05-persistence-durability-tiers | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:43:27Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:43:27Z | 02.06-schema-versioning-and-migration | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:43:27Z | 02.05-persistence-durability-tiers | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:41:19Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:41:19Z | 02.06-schema-versioning-and-migration | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:41:19Z | 02.05-persistence-durability-tiers | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:39:29Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:36:28Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:35:43Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:29:50Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-10T14:29:36Z | 02.05-persistence-durability-tiers | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:19:18Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-10T14:19:18Z | 02.01-state-taxonomy-and-ownership | letta | failed | validation_failed |
| 2026-07-10T14:10:50Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T14:08:32Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-10T14:07:58Z | 02.05-persistence-durability-tiers | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T10:37:58Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-10T10:37:58Z | 02.01-state-taxonomy-and-ownership | letta | failed | validation_failed |
| 2026-07-10T09:04:15Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T08:58:32Z | 02.01-state-taxonomy-and-ownership | letta | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T08:57:30Z | 02.04-mutation-discipline-and-state-transitions | langgraph | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T08:57:30Z | 02.04-mutation-discipline-and-state-transitions | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-10T08:19:53Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-06T21:14:53Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-06T17:27:11Z | 02.01-state-taxonomy-and-ownership | (synthesis) | failed | synthesis dependencies failed or were cancelled |
| 2026-07-06T17:24:45Z | 02.01-state-taxonomy-and-ownership | letta | failed | validation_failed |
| 2026-07-04T13:04:29Z | 01.09-delivery-guarantees-idempotency | (synthesis) | failed | validation: validation: required output validation failed |
| 2026-07-03T16:44:39Z | 01.09-delivery-guarantees-idempotency | (synthesis) | failed | validation: validation: required output validation failed |
| 2026-07-03T08:52:04Z | 01.07-concurrency-and-parallel-advancement | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-03T08:52:04Z | 01.06-scheduling-and-trigger-semantics | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-03T08:52:04Z | 01.06-scheduling-and-trigger-semantics | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-02T17:54:38Z | 01.07-concurrency-and-parallel-advancement | crewai | cancelled | context canceled |
| 2026-07-02T17:54:38Z | 01.06-scheduling-and-trigger-semantics | temporal | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-02T17:54:38Z | 01.07-concurrency-and-parallel-advancement | agent-framework | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-02T17:54:38Z | 01.06-scheduling-and-trigger-semantics | openhands | cancelled | validation: cancellation: validation run was cancelled |
| 2026-07-02T17:48:33Z | 01.06-scheduling-and-trigger-semantics | langgraph | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T17:47:43Z | 01.06-scheduling-and-trigger-semantics | crewai | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T17:47:15Z | 01.06-scheduling-and-trigger-semantics | agent-framework | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T16:49:10Z | 01.03-step-turn-task-atomicity | temporal | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T15:00:12Z | 01.03-step-turn-task-atomicity | openai-agents-sdk | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T15:00:12Z | 01.03-step-turn-task-atomicity | openhands | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T14:48:14Z | 01.03-step-turn-task-atomicity | letta | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T14:36:03Z | 01.03-step-turn-task-atomicity | langgraph | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T14:36:03Z | 01.03-step-turn-task-atomicity | crewai | failed | opencode run: timeout: OpenCode run timed out |
| 2026-07-02T14:24:05Z | 01.03-step-turn-task-atomicity | agent-framework | failed | opencode run: timeout: OpenCode run timed out |

