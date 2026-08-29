# Source Selection and Snapshot Manifest

## Selection Gates

Each source had to meet all of these conditions:

- The implementation is predominantly Go and the relevant path is implemented in Go.
- It is used as a real product, runtime, SDK, or storage component, or is a compact specialist library with unusually clear semantics.
- It contains tests for failure, cancellation, concurrency, recovery, overload, or another Aren-relevant edge.
- It answers a named Aren roadmap question better than a second generic framework would.
- Its applicability can be narrowed so the study does not waste runs on irrelevant dimensions.

The set intentionally mixes two direct agent runtimes with specialist systems. Agent repositories answer domain questions. The specialist repositories expose mature failure handling that agent projects often hide behind provider SDKs, databases, containers, or server frameworks.

## Selected Sources

The clones were created with `git clone --depth 1` on 26 August 2026. Commit IDs make the evidence reproducible even when the upstream default branch advances.

| Source | Snapshot | Primary role in this study | Why it is not sufficient alone |
|---|---|---|---|
| [sourcegraph/conc](https://github.com/sourcegraph/conc) | `5f936abd7ae8` | Small, readable structured-concurrency and panic-propagation reference | It does not own persistence, model streams, tools, or processes. |
| [charmbracelet/crush](https://github.com/charmbracelet/crush) | `abc246403b95` | Coding-agent loop, providers, tool permission, context, and sessions | Product and TUI concerns are intertwined, and its licence requires separate review before reuse. |
| [docker/docker-agent](https://github.com/docker/docker-agent) | `34ce4fa713ac` | Multi-provider agent runtime, tools, permissions, sandboxing, composition, and server sessions | It is broad and comparatively young, so patterns need corroboration from narrower systems. |
| [ollama/ollama](https://github.com/ollama/ollama) | `e2c6c7e89493` | Live model serving, cancellation, model residency, scheduling, and local API ownership | Local model scheduling has different cost and transport assumptions from hosted providers. |
| [opencontainers/runc](https://github.com/opencontainers/runc) | `9674194cc2b9` | Process lifecycle, signalling, containment, descendant cleanup, and executor seams | Container lifecycle is stronger and lower-level than Aren's default local-tool requirement. |
| [nats-io/nats-server](https://github.com/nats-io/nats-server) | `e91c5c6adb11` | Slow-consumer handling, bounded queues, persistence, recovery, and server shutdown | It is a network messaging server, so its distributed machinery is not a default Aren architecture. |
| [cockroachdb/pebble](https://github.com/cockroachdb/pebble) | `0457a364b877` | WAL, atomic commit, checkpoints, format evolution, and adversarial storage tests | It is a storage engine, not an application persistence design for runs and events. |
| [temporalio/sdk-go](https://github.com/temporalio/sdk-go) | `d9e04e81963b` | Durable cancellation contracts, deterministic replay, retry, children, and workflow tests | Its workflow model belongs late in Aren's roadmap and must not shape Phase 1. |
| [hashicorp/nomad](https://github.com/hashicorp/nomad) | `aa026cc99cfb` | Admission, evaluation queues, feasibility, fairness, daemon ownership, and remote drivers | Cluster scheduling and its source-available licence both argue against direct transplantation. |

## Web Research Findings That Drove Inclusion

- `conc` explicitly targets safer goroutine ownership, panic handling, bounded pools, and context cancellation. That makes it a better Phase 1 and composition source than a large server that happens to use goroutines.
- Crush exposes a concrete Go coding-agent orchestration path, multiple providers, tool permission hooks, context compaction, and SQLite sessions. It supplies direct evidence for Aren Phases 3 through 10.
- Docker Agent adds structured output, MCP and built-in tools, permission modes, sandbox execution, multi-agent composition, and server sessions. It provides a second implementation for the agent-specific dimensions where infrastructure analogies would be weak.
- Ollama's server scheduler owns model loading and request channels under memory pressure. It directly tests the shape of admission, cancellation, streaming, and resource accounting around model work.
- runc is the OCI reference runtime and handles signals, namespaces, cgroups, descriptor hygiene, and process cleanup. It is included for the process boundary only, not as an argument that Aren needs containers.
- NATS documents and tests explicit slow-consumer behaviour, and JetStream adds persistence and replay. This gives the observation study a system where backpressure is a primary concern rather than an incidental channel buffer.
- Pebble is production-ready inside CockroachDB and exposes WAL, checkpoints, format versions, commit code, recovery code, metamorphic tests, and fault-oriented test infrastructure in one Go repository.
- Temporal's Go SDK makes cancellation modes, retry policy, child ownership, determinism, and replay part of its public contract. It is a bounded late-roadmap source compared with cloning the much larger Temporal server again.
- Nomad documents its persisted evaluation queue, long-lived scheduler workers, feasibility checks, plan validation, and Ack or Nack ownership. Those are direct references for admission and remote-executor questions.

## Deliberate Exclusions

- The existing `ultraplan-daemon-events-study` already covers Temporal Server, containerd, BuildKit, and Dagster in depth. Duplicating those clones would add cost without adding a new comparison.
- Tailscale is an excellent Go daemon, but its strongest contribution overlaps the existing daemon study and the selected local-server sources.
- Very recent independent Go agent frameworks were excluded from the primary set because they lack the operating history requested by “elite repos”. They may be added later as compact contrast sources if Crush and Docker Agent leave an agent-loop question unanswered.
- Kubernetes, CockroachDB, and other very large Go monorepos were not cloned wholesale when a smaller repository exposed the relevant mechanism more directly.

## Interpretation Rule

A high score means the source handles its own problem well. It does not mean Aren should adopt the same subsystem. Every synthesis must distinguish:

1. a requirement Aren already has;
2. a small Go pattern that helps meet it;
3. complexity justified only by the source's scale or product;
4. a later trigger that would make that complexity relevant to Aren.
