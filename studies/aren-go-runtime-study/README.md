# aren-go-runtime-study

Evidence study of elite Go repositories for Aren's roadmap, ordered from the Phase 1 lifecycle contract through providers, tools, bounded agents, persistence, hosting, workflows, and later execution seams.

## Study Method

This is a code-evidence study, not a catalogue of popular Go projects. A source was included only when it provides an inspectable production path for an Aren requirement, meaningful failure tests, and a useful contrast with at least one other source. Repository size or star count was not enough.

The 15 dimensions are execution priorities. `01.01` through `01.03` protect the Phase 1 contract and its verification. The remaining dimensions follow the dependency order in Aren's roadmap. Later infrastructure patterns must not be pulled into an earlier phase merely because they score well in their native system.

For every report:

1. Cite repository-relative code and test paths with line numbers.
2. Separate observed behaviour from interpretation.
3. State the source's scale and product assumptions before mapping a pattern to Aren.
4. Record the smallest useful pattern, the complexity it carries, and the trigger that would justify it.
5. Link the finding to the Aren requirement named in the dimension.
6. Include negative evidence, unresolved questions, and tests Aren would still need.

The applicability map yields 43 focused source reports rather than 135 repository-by-dimension combinations. A source must not be analysed for dimensions absent from its metadata.

## Priority Map

| Priority | Aren roadmap concern | Dimensions |
|---:|---|---|
| 1 | Phase 1 lifecycle correctness | `01.01`-`01.03` |
| 2 | Phase 1 observation and Phases 3-6 model execution | `01.04`-`01.05` |
| 3 | Phases 7-8 tools and local effects | `01.06`-`01.07` |
| 4 | Phases 9-10 bounded agents and context | `01.08` |
| 5 | Cross-cutting resource evidence and Phase 13 pressure | `01.09` |
| 6 | Phases 11-12 persistence and composition | `01.10`-`01.11` |
| 7 | Phases 13-15 governance and hosting | `01.12`-`01.13` |
| 8 | Phases 16-18 workflows and later execution seams | `01.14`-`01.15` |

The numbers define priority, but they do not authorise roadmap progression. Aren's phase exit criteria and real-use evidence still decide whether a later dimension becomes implementation work.

## Sources

- `conc`: Compact structured-concurrency library whose small surface makes goroutine ownership, cancellation, panic propagation, bounded pools, and composition easy to inspect in full. (sources/conc)
- `crush`: Production Go coding agent with provider adapters, streaming sessions, tool permissions, MCP and LSP integration, context management, and SQLite-backed session state. (sources/crush)
- `docker-agent`: Go agent builder and runtime from Docker with multiple providers, MCP and built-in tools, sandboxing, permissions, structured output, multi-agent composition, and server sessions. (sources/docker-agent)
- `ollama`: Widely deployed Go model server with live response streaming, request cancellation, model scheduling, bounded loading, resource accounting, and a long-lived local API. (sources/ollama)
- `runc`: Reference OCI runtime implementation in Go, selected for process lifecycle, signal forwarding, namespace and cgroup containment, descendant cleanup, and executor boundaries. (sources/runc)
- `nats-server`: Mature Go server selected for ordered event delivery, slow-consumer handling, bounded queues, JetStream persistence, restart behaviour, and graceful server ownership. (sources/nats-server)
- `pebble`: Production Go storage engine selected for WAL and commit behaviour, checkpoints, recovery, format evolution, fault injection, metamorphic testing, and bounded storage work. (sources/pebble)
- `temporal-sdk-go`: Mature Go SDK for durable workflow execution, selected for cancellation states, deterministic replay, retry policy, child ownership, testing, and recovery-facing contracts. (sources/temporal-sdk-go)
- `nomad`: Production Go scheduler selected for evaluation queues, admission, resource feasibility, fairness, Ack and Nack ownership, persisted work, daemon operation, and remote executor seams. (sources/nomad)

## Dimensions

- `01.01`: Lifecycle Transition Ownership and Terminal Arbitration (`dimensions/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md`)
- `01.02`: Cancellation, Goroutine Ownership, and Cleanup (`dimensions/01.02-cancellation-goroutine-ownership-and-cleanup.md`)
- `01.03`: Adversarial Concurrency and Failure Verification (`dimensions/01.03-adversarial-concurrency-and-failure-verification.md`)
- `01.04`: Ordered Observation, Live Streaming, and Backpressure (`dimensions/01.04-ordered-observation-live-streaming-and-backpressure.md`)
- `01.05`: Provider Boundaries, Structured Results, and Retry (`dimensions/01.05-provider-boundaries-structured-results-and-retry.md`)
- `01.06`: Tool Contracts, Permissions, and Bounded Results (`dimensions/01.06-tool-contracts-permissions-and-bounded-results.md`)
- `01.07`: Subprocess Supervision and Process-Tree Containment (`dimensions/01.07-subprocess-supervision-and-process-tree-containment.md`)
- `01.08`: Bounded Agent Loop, Context, and Evidence (`dimensions/01.08-bounded-agent-loop-context-and-evidence.md`)
- `01.09`: Resource Accounting, Overload, and Bounded Work (`dimensions/01.09-resource-accounting-overload-and-bounded-work.md`)
- `01.10`: Persistence, Atomic Recovery, and Schema Evolution (`dimensions/01.10-persistence-atomic-recovery-and-schema-evolution.md`)
- `01.11`: Execution Composition and Child Ownership (`dimensions/01.11-execution-composition-and-child-ownership.md`)
- `01.12`: Policy, Admission, Scheduling, and Fairness (`dimensions/01.12-policy-admission-scheduling-and-fairness.md`)
- `01.13`: Daemon, Local API, and Multi-Client Lifetime (`dimensions/01.13-daemon-local-api-and-multi-client-lifetime.md`)
- `01.14`: Durable Workflows, Retry, Idempotency, and Replay (`dimensions/01.14-durable-workflows-retry-idempotency-and-replay.md`)
- `01.15`: Executor Boundaries and Remote Readiness (`dimensions/01.15-executor-boundaries-and-remote-readiness.md`)

## Generated Paths

- `study-init.yml`
- `study.json`
- `dimensions/`
- `sources/`
- `reports/source/`
- `reports/final/`
- `SOURCE_SELECTION.md`

Edit `study.json` to run selected dimensions before the remaining dimensions.

## Next Commands

- `ultraplan study list`
- `ultraplan study aren-go-runtime-study list`

## Scope and Licensing

The repositories are evidence sources, not dependencies or copy-ready templates. In particular, current Crush and Nomad snapshots use source-available licences rather than permissive open-source licences. Review the licence at the pinned source commit before copying any implementation. Architectural observations and independently written Aren code remain the intended use of this study.
