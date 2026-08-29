> **Inputs Used:** `projects/ultraplan-go/sprints/36-read-only-qa/requirements.md`, `projects/ultraplan-go/sprints/36-read-only-qa/sprint-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture: read-only QA decomposition and synthesis

This area decides where Sprint 36 behavior and state belong, how runtime-backed work crosses existing boundaries, and how detailed QA state survives cancellation and restart without becoming a second run-control system. It does not define frontend layout or public transport schemas.

## Area Decisions

### Proceed with a focused sprint-module extension

The current module architecture still fits, but Sprint 36 needs a small refactor before feature work. `PlanningStage` currently carries some verification compatibility values and `Service.stageRuntime` is keyed by that type. QA must not inherit either shape.

The refactor introduces `VerificationPhase` in `internal/sprint/verification_phase.go` with these values:

```text
conformance-review
qa
repair
```

`VerificationPhase` owns verification identity and compatibility mapping. `PlanningStages()` remains unchanged. Existing `StageReview`, `StageSmoke`, `review`, `review.md`, `smoke`, `smoke.md`, verdicts, and JSON fields remain compatibility contracts. Human-facing review labels become "Conformance Review." The `conformance-review` alias maps to the existing review capability and never creates a second reviewer or state record.

Verification runtime configuration is keyed by `VerificationPhase`, not `PlanningStage`. The sprint service receives a small `QASettings` value and a verification runtime selection at construction. Runtime-free mapping and status use a service that has no need to initialize an investigator runtime. This keeps `qa --dry-run` honest and avoids adding QA to `map[PlanningStage]StageRuntime`.

### Keep all QA semantics in `internal/sprint`

Sprint 36 uses focused files in the existing package:

| File | Ownership |
| --- | --- |
| `verification_phase.go` | Phase values and review compatibility mapping. |
| `qa_types.go` | Schema, map, shard, theory, outcome, budget, synthesis, fingerprint, and ID types. |
| `qa_state.go` | Contained paths, strict load/validation, atomic writes, current pointer, resume, cancellation, and reconciliation. |
| `qa_map.go` | Deterministic input normalization and behavioral shard construction. |
| `qa_prompt.go` | Bounded prompt packets, structured runtime contracts, and enforced permission policy. |
| `qa.go` | Attempt orchestration, bounded scheduling, focused execution, cancellation, resume, invalidation, progress, and next action. |
| `qa_synthesis.go` | Theory deduplication, contradiction retention, interactions, and bounded parent-linked follow-up. |

No `internal/qa`, generic workflow package, scheduler package, issue package, plugin registry, or sandbox framework is introduced. Similarity to review and smoke is not enough reason to extract their product rules. Small mechanical helpers may be reused when they already have a stable owner, such as workspace containment, generic runtime execution, process execution, and atomic file mechanics.

`internal/app` owns typed QA use cases, operation preparation, durable acceptance, run correlation, and adapter-neutral projections. CLI, TUI, and web translate requests and render results. `internal/web` continues importing only `internal/app` and the standard library. Neither adapter reads verification files, constructs prompts, calls runtimes, or decides outcomes.

Dependency direction remains:

```text
cmd/ultraplan -> internal/app
internal/tui -> internal/app
internal/web -> internal/app
internal/app -> internal/sprint + internal/runcontrol + platform composition
internal/sprint -> internal/project + internal/workspace + internal/platform/runtime + internal/platform/process
internal/runcontrol -> no sprint QA package
internal/platform/* -> no product package
```

### Separate semantic, summary, operational, and source authority

Four authorities remain distinct:

| Concern | Authority |
| --- | --- |
| Map, shard, theory, attempt, synthesis, freshness, and QA next action | Sprint-owned files under `verification/`. |
| Bounded phase summary and contained pointers or digests | `flow-state.json`. |
| Acceptance, owner lease, fencing, event order, cancellation request, liveness, and terminal execution result | Sprint 35 run control. |
| Production source, tests, governed inputs, verification code, harness content, and Git identity | The contained workspace and target checkout. |

The detailed layout is fixed:

```text
verification/state.json
verification/attempts/<attempt-id>/map.json
verification/attempts/<attempt-id>/shards/<shard-id>.json
verification/attempts/<attempt-id>/synthesis.json
```

`verification/state.json` is the semantic entry point. It contains schema version, project and sprint identity, the current attempt pointer, current fingerprints, verdict-neutral status, freshness, bounded counts, blocker, next action, cancellation summary, run correlation, and digests of referenced files. It does not embed maps, theories, command output, or full history.

`flow-state.json` receives only a bounded QA projection and a contained pointer to `verification/state.json`. It cannot become the theory database. A planning refresh must preserve the QA projection just as current stage writes preserve review and smoke records.

Run control never decides whether a theory is confirmed, whether synthesis is complete, or whether evidence is current. Conversely, QA state never claims an owner is live or arbitrates a run terminal outcome. The app projection joins those facts without changing either authority.

### Use schema major version 1 and verification-scoped IDs

Sprint 36 writes QA schema major version `1`. There is no legacy QA schema to migrate. The loader accepts version 1 only, strictly rejects unknown fields where the file contract is closed, rejects missing or unknown versions, rejects trailing JSON, and returns recovery guidance without rewriting the file.

Future migration support must be an explicit pure `vN -> vN+1` function with fixture coverage and atomic promotion. Unknown future major versions remain unreadable by older binaries. Additive public projections can evolve independently, but durable file changes require a schema decision.

Map, shard, theory, and attempt IDs include the `qa-v1` scope and derive from canonical normalized inputs. They are stable inside schema version 1 and the selected project/sprint verification scope. They are not global content IDs and make no promise across an unsupported major migration.

Observation timestamps, run IDs, sessions, provider metadata, worker order, and current time are excluded from map, shard, and theory identity. Every unordered collection is sorted before ID generation, fingerprinting, and persistence.

### Publish detailed state atomically in dependency order

Every path resolves through sprint-root containment before any read or write. Verification directories use private permissions, `0700` for directories and `0600` for detailed JSON files. Each candidate is fully validated before publication.

Atomic publication uses this sequence:

1. Create a temporary file in the destination directory.
2. Write normalized indented JSON with a trailing newline.
3. Flush and close the temporary file.
4. Rename it over the destination.
5. Sync the containing directory.
6. Update any file that points to the newly durable file only after that file is durable.

An attempt map is immutable after publication. Shard and synthesis records replace only their own prior valid versions. `verification/state.json` advances its pointer or digest last. The bounded `flow-state.json` projection is written after detailed state. If flow-state projection fails, detailed QA state remains authoritative and explicit recovery repairs the summary. A failed rename or validation keeps the prior valid file and pointer unchanged.

Multi-file recovery follows references and digests from `verification/state.json`. Unreferenced complete files may be retained as stale diagnostic records within attempt and storage bounds, but status never promotes them by directory discovery alone. A referenced missing or mismatched file makes the affected state invalid or stale and gives a recovery action; read-only status does not silently repair it.

### Bind product writes to accepted-run fencing without importing run control

Runtime-backed QA is durably accepted and claimed before `internal/sprint` starts an attempt. `internal/app` converts accepted run ownership into an opaque sprint execution token containing the run ID, operational attempt ID, and fencing generation. Sprint types own this token shape and do not import `internal/runcontrol`.

`verification/state.json` records the current writer token as correlation and stale-writer protection, not as an operational lifecycle database. Every detailed-state promotion checks that its expected token still matches the current token. A new accepted resume or recovery run installs a newer token before scheduling work. A worker with an old token records no further state after ownership loss.

The existing in-process sprint mutation lock still serializes local mutations. The writer token prevents an older accepted worker from overwriting newer QA state across processes. Run control remains the authority for whether the token's owner is live, fenced, cancelled, or terminal.

Runtime-free `qa --dry-run` and status use no writer token and write no detailed state. Explicit recovery obtains sprint mutation ownership and, when it starts child work, a durable accepted run.

### Give mapping, investigation, and synthesis separate authority

The mapper is deterministic and runtime-free. It consumes the governed inputs required by the sprint contract and produces one primary owner for every changed path. Boundary shards may overlap primary paths only for named cross-package, interface, producer/consumer, state-transition, public-API, or cross-cutting behavior. Unknown paths become visible blocked primary work rather than disappearing.

The map fixes all budgets, approved context, known checks, fingerprints, and shard relationships before investigation. An investigator cannot add a path, command, iteration, theory, or follow-up budget to itself.

Investigators receive one shard packet and approved common context. They return structured theories, static evidence, requested context expansion, requested existing checks, and stop reasons. The orchestrator validates every response before state promotion. Negative and blocked outcomes are complete evidence and remain durable.

Context expansion is a request to the orchestrator. The orchestrator checks the map budget, containment, relevance to the shard, and target fingerprint before granting it. A denial becomes inspectable outcome data, not an excuse to broaden the prompt.

Synthesis reads only validated shard records for the current map fingerprint. It deduplicates equivalent theories, retains contradictions and negative outcomes, identifies interactions, and may request only the map's remaining parent-linked follow-up budget. It cannot create issues, modify review verdicts, write `qa.md`, generate checks, or authorize repair.

### Execute approved checks outside the agent tool policy

Agentwrap remains the investigator runtime supervisor. The runtime request uses `read_only`, restricted permissions, default deny, and required permission capability. The agent receives read, list, and search access only to assigned paths. Unsupported permission enforcement blocks the shard before useful output can be promoted.

The agent never receives a general shell tool. When it requests an existing check, it selects a stable check ID from the map and supplies no executable, path, shell text, or environment. `internal/sprint` resolves that ID to a product-owned typed descriptor, validates current fingerprints and remaining budgets, then calls `internal/platform/process` with explicit executable, argv, working directory, environment allowlist, timeout, and output limits.

`internal/platform/process` knows nothing about QA, shards, or theory outcomes. It executes an already approved request and returns bounded process facts. Sprint code interprets those facts and feeds safe evidence into the next investigator iteration. Shell wrappers, command substitution, interpreters used as indirection, Git commands, output redirection, and caller-supplied environment are rejected.

This split gives investigators useful existing checks without pretending that prompt wording makes a shell read-only.

### Check target identity before and after every investigator attempt

The identity manifest covers these categories separately:

```text
production source
production tests
governed sprint and project inputs
verification implementation code
smoke-harness content
Git HEAD, index, and worktree state when the target is a Git checkout
```

Each category hashes sorted normalized relative paths and file digests. Symlinks are recorded by link identity and resolved containment; links that escape an approved root block the attempt. Declared QA state files and normal run-control records are excluded from target identity.

The pre-attempt manifest becomes part of the shard fingerprint. After the runtime and any approved checks stop, the service recomputes the same manifest before promoting outcomes. Any drift blocks the attempt, records the changed identity category and bounded paths, preserves the prior canonical outcome, and stops new scheduling. It never repairs or resets the checkout.

For non-Git targets, content identity still runs and Git identity records `not_applicable`. Gated Sprint 36 dogfood against the implementation repository requires Git identity and blocks if it is unavailable.

### Bound scheduling and classify stop policy

One QA run owns an instance-scoped scheduler. It has finite worker and queue capacity; it never creates one goroutine per shard. The scheduler stops admission before cancellation, propagates the accepted-run context to workers and subprocesses, and waits within the whole-run deadline.

These conditions stop the whole current attempt immediately:

- Missing, stale, or invalid governed inputs.
- Unsupported QA schema or corrupt current pointer.
- Writer-token or accepted-run ownership loss.
- Runtime permission enforcement unavailable.
- Target identity drift.
- Cancellation or whole-run timeout.
- Inability to persist current authoritative state.

These conditions remain shard-local while independent shards continue:

- Malformed investigator output.
- An approved check unavailable or nonzero.
- Context expansion denied or exhausted.
- Shard iteration, command, output, or wall-clock budget exhausted.
- A theory ending blocked, invalid, inconclusive, refuted, or not applicable.

Synthesis may complete over valid finished shards while reporting incomplete or blocked siblings. It cannot label partial work as globally complete without blockers.

### Separate cancellation from bounded cleanup

Cancellation follows one canonical path through run control. Once requested, the scheduler stops admission, cancels active investigator and process contexts, and preserves completed valid shards. Browser or SSE disconnect only removes an observer.

Work cancellation and cleanup use different contexts. After the work context is cancelled, the owner gets a separate 30-second cleanup context to stop process groups, record truthful active-shard outcomes, flush detailed state, update the bounded summary, emit final safe events, and release local resources. Cleanup timeout produces `cleanup_uncertain`; it never produces success.

A restart with unchanged fingerprints creates a new durable operational run and writer token. It reuses the same semantic QA attempt and completed valid shards. It does not adopt dead goroutines, sessions, or process handles. Fingerprint changes mark only dependent work stale and require a new map when map identity changes.

### Keep configuration immutable for an attempt

Effective QA configuration follows this precedence:

```text
explicit local CLI option > environment-backed workspace configuration > workspace configuration > product default
```

Only model selection and lower bounded limits may be configurable in Sprint 36. Browser callers cannot override models, budgets, commands, permissions, environment, or paths. The complete effective settings are validated after merge and frozen into the map and policy fingerprints. Zero or negative limits fail closed. A lower configured value is valid; values above product maxima fail validation rather than being silently clamped.

The attempt stores effective values and their sources, not only config paths. A restart reuses completed work only when those values and the investigator policy fingerprint remain unchanged.

### Use narrow fault seams, not a new persistence abstraction

The core mapper, validator, ID generation, synthesis, and state transitions remain concrete functions and sprint-owned types. Interfaces or function seams exist only around volatile behavior:

- Agent runtime execution.
- Existing-check process execution.
- Clock and owner identity.
- Target identity reads.
- Atomic write fault hooks used by tests.

The filesystem remains the selected product authority. Sprint 36 does not add a generic repository, virtual filesystem, alternate SQLite product store, dual writes, or import/export layer. Tests use temporary real directories for path, symlink, permission, fsync, and rename behavior. Small injected hooks cover failures that are hard to induce portably.

### Preserve one progress and event path

Sprint code emits bounded semantic progress containing attempt ID, shard ID, phase, status, completed, total, blocker category, and safe message. `internal/app` maps it to the existing operation event vocabulary. Run control assigns durable correlation, order, retention, cancellation, and terminal facts before transient delivery.

Redaction occurs before a value enters QA state, run-control events, logs, CLI progress, TUI history, or browser delivery. Raw prompts, provider payloads, absolute paths, command environment, unrestricted argv, and full command output are never operation events. Existing 16 KiB durable event and 250 ms progress coalescing limits remain unchanged.

### Freeze hard resource limits

Sprint 36 uses these defaults and maxima. Configuration can lower a value but cannot raise it above the maximum.

| Resource | Default | Maximum |
| --- | ---: | ---: |
| Changed paths in a map | 512 | 512 |
| Primary shards | 32 | 32 |
| Boundary shards | 8 | 8 |
| Changed paths per primary shard | 32 | 64 |
| Contextual paths per shard | 64 | 128 |
| Context expansions per shard | 2 | 4 |
| Paths per expansion | 16 | 32 |
| Behavioral concerns per shard | 12 | 24 |
| Theories per shard | 12 | 24 |
| Investigator iterations per shard attempt | 4 | 8 |
| Approved commands per shard attempt | 8 | 16 |
| Runtime retries after the initial call | 1 | 2 |
| Concurrent investigators | 3 | 8 |
| Follow-up shards per synthesis | 4 | 4 |
| Command wall-clock duration | 5 minutes | 10 minutes |
| Shard wall-clock duration | 20 minutes | 30 minutes |
| Whole-run wall-clock duration | 60 minutes | 90 minutes |
| Cleanup duration | 30 seconds | 30 seconds |
| Captured output per command | 256 KiB | 512 KiB |
| Stored investigator output per shard attempt | 1 MiB | 2 MiB |
| Investigator or synthesizer prompt | 512 KiB | 1 MiB |
| Recent progress records in a product summary | 100 | 200 |

Hitting a limit records the exhausted resource and stops the affected scope. It never widens scope, drops authoritative terminal data, or truncates content while claiming completeness. Run-control event, subscriber, retention, and workspace quota limits remain Sprint 35 values.

### Testing follows the ownership boundaries

Normal tests stay offline and deterministic:

- Pure tests freeze normalized maps, IDs, ordering, fingerprints, budgets, theory validation, synthesis, and affected-work invalidation.
- Temporary-workspace tests cover containment, symlink escape, strict JSON, unknown schema, atomic rename failure, directory sync behavior, pointer ordering, stale writer tokens, and prior-state preservation.
- Fake runtime and process seams cover permission rejection, malformed output, queue saturation, cancellation while queued and running, command timeout, process cleanup, output exhaustion, and unsupported capabilities.
- App fixtures project one semantic result through CLI, JSON, TUI, HTML, HTTP, and durable run inspection.
- Race tests cover scheduler bounds, cancellation, state promotion, and stale writers.
- Gated dogfood is the only normal path to a real runtime. Missing prerequisites produce blocked evidence and do not satisfy the exit gate.

The implementation proceeds after the small phase/runtime refactor. It does not require a larger package reorganization.

## Trade-Offs

### Separate detailed files instead of one state document

Multiple files add pointer ordering, reconciliation, and migration work. They keep status bounded, make immutable maps inspectable, isolate shard promotion failures, and prevent `flow-state.json` from growing with every theory. A single large JSON document was rejected because every shard update would rewrite unrelated evidence and increase corruption risk.

### Filesystem authority instead of product SQLite

Atomic multi-file publication is harder than one database transaction. The current product still treats workspace artifacts as authoritative, and Sprint 35 SQLite has a different job. Reusing it for theories would blur product and operational authority. Alternate product persistence remains gated until real browser workflows justify it.

### Opaque writer token instead of importing run control into sprint

Passing a small accepted-run token adds correlation fields to QA requests and state. It preserves one-way dependencies and gives stale workers a product-level write check without teaching sprint code about leases, event repositories, or SQLite. A direct `internal/sprint -> internal/runcontrol` import was rejected because it would make product semantics depend on operational storage types.

### Product-executed checks instead of agent shell access

The request/validate/execute/feed-back loop costs extra iterations and supports fewer existing checks than a shell. It makes command policy code-reviewable, keeps explicit argv, and gives output and timeout enforcement one owner. A read-only prompt plus shell access was rejected because shell behavior can mutate state indirectly.

### Independent shard collection instead of universal fail-fast

Independent continuation preserves negative and inconclusive evidence when one shard fails locally. It consumes more runtime after a local failure. Whole-attempt conditions still fail fast when state authority, permissions, identity, or ownership are unsafe.

### Concrete sprint state store instead of broad interfaces

A concrete store couples QA to current filesystem authority. That is deliberate. There is no second product persistence implementation, and path, mode, fsync, rename, and recovery semantics need real tests. Narrow write hooks and temporary workspaces provide fault coverage without a fake filesystem that may lie about atomicity.

### One schema version with rejection instead of speculative migrations

Version 1 cannot demonstrate migration from a shipped QA format because none exists. Inventing version 0 would create false compatibility. Sprint 36 freezes v1 fixtures, rejects every other major, and leaves a clear pure-function migration seam for the first real format change.

### Fixed conservative limits instead of adaptive scaling

Fixed bounds may block unusually large changes. They make behavior reproducible and keep memory, runtime, output, and state size reviewable. Adaptive widening would change map identity based on runtime conditions and could silently expand investigator authority.

### Rejected alternatives

- Add QA to `PlanningStage` or planning flow order.
- Create a second Conformance Review implementation for the new label.
- Put detailed QA records in `flow-state.json` or run-control events.
- Let browser memory, SSE, TUI state, or runtime sessions own recovery.
- Add a general-purpose QA package, scheduler, workflow engine, repository layer, plugin system, remote worker, or daemon.
- Let investigators create checks, write fixtures, promote issues, or repair code.
- Treat synthesis as adjudication or use QA to improve a failed review verdict.
- Use package globals, context values as service lookup, or `init()` registration for QA dependencies.
- Use in-memory mutexes or sessions as cross-process stale-writer proof.
- Make Git mutation or cleanup part of recovery.
- Add OpenTelemetry, a knowledge graph, retrieval, or content identity as part of this architecture.

## Evidence

### Governed project evidence

`projects/ultraplan-go/sprints/36-read-only-qa/requirements.md` fixes the product boundary. It assigns QA semantics to `internal/sprint`, adapter-neutral use cases to `internal/app`, and presentation to CLI, TUI, and web. It requires detailed state outside flow state, durable acceptance before child work, read-only target access, deterministic mapping, bounded synthesis, and no issue promotion or repair.

`projects/ultraplan-go/docs/ARCHITECTURE.md` establishes module ownership, one-way product/platform dependencies, current filesystem authority, and the distinction between product state and Sprint 35 operational run records. Its Phase 5 gate puts read-only QA before writable isolation, adjudication, and repair.

`projects/ultraplan-go/docs/PRD.md` and `projects/ultraplan-go/docs/TRD.md` establish the local product topology, Conformance Review compatibility, same-host durable run model, strict adapter boundaries, schema and atomic-write expectations, cancellation and recovery behavior, and the later gates for smoke integration and repair.

### Report evidence and project inference

The selected reports compare other Go CLIs. They support patterns and expose failure modes, but the Sprint 36 architecture above is an UltraPlan decision.

- `studies/go-cli-study/reports/final/01-project-structure.md` finds thin entrypoints and one-way imports in chezmoi at `main.go:16`, yq at `cmd/root.go:9`, and Helm's `pkg/cmd` to `pkg/action` split. This supports keeping adapters thin and product behavior in `internal/sprint`. The exact Sprint 36 file split comes from the sprint contract.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` finds manual composition roots and narrow replaceable dependencies in gh-cli at `internal/ghcmd/cmd.go:52-132` and restic at `internal/restic/repository.go:18-66`. It warns about global configuration and singleton caches such as rclone `fs/config.go:793` and `fs/cache/cache.go:16-21`. This supports explicit service construction and no package-global QA runtime.
- `studies/go-cli-study/reports/final/04-configuration-management.md` documents explicit-value preservation in chezmoi `internal/cmd/config.go:2253-2287`, changed-flag tracking in restic `internal/global/global.go:139,147`, and post-merge validation in opencode `internal/config/config.go:609-641`. This supports immutable effective attempt settings. The exact precedence and limits are Sprint decisions.
- `studies/go-cli-study/reports/final/05-error-handling.md` supports typed, cause-preserving failure categories through rclone `fs/fserrors/error.go:22-192`, restic `internal/errors/fatal.go:10-53`, and Helm's aggregate cleanup errors at `pkg/action/uninstall.go:232-254`. This supports separate blocked, stale, cancelled, ownership, persistence, and cleanup-uncertain handling rather than string matching.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` finds useful seams at terminal, filesystem, backend, and process boundaries, including restic `internal/fs/interface.go:10-31` and lazygit `pkg/commands/oscommands/cmd_obj_runner.go:18-23`. It also shows the cost of broad interfaces. This supports narrow runtime/process/identity seams and a concrete QA filesystem store.
- `studies/go-cli-study/reports/final/07-state-context.md` distinguishes root context cancellation, process-local sessions, persistent state, in-flight maps, and locks. Helm `pkg/cmd/install.go:333-347` supports root cancellation, while restic `internal/restic/lock.go:105,290-305` supports explicit ownership and separately bounded cleanup. The writer-token and restart reconstruction design are UltraPlan inferences.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports bounded admission and timed shutdown through k9s `internal/pool.go:21,30,37` and opencode `cmd/root.go:261-279`. It identifies gh-cli `pkg/cmd/extension/manager.go:196-206` as unbounded fan-out. This supports an instance-owned scheduler with separate worker and queue limits.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports structured fields, stream separation, and instrumentation wrappers through Helm `internal/logging/logging.go:31-71`, k9s `internal/slogs/keys.go:6-231`, and restic `internal/backend/logger/log.go:22-77`. It does not prove redaction or durable parity. The decision to redact before state and event persistence follows the sprint security contract.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports command scenarios, process fakes, persistent-state fakes, and fault seams through chezmoi `internal/cmd/main_test.go:64-174`, lazygit `pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, and restic `internal/backend/mock/backend.go:14-26`. Real temporary files remain necessary for containment, symlink, permission, and atomic-write tests.
- `studies/go-cli-study/reports/final/12-extensibility.md` shows that narrow injected adapters and static in-tree registration can support substitution without a plugin system, including dive `cmd/dive/cli/internal/command/adapter/analyzer.go:13-15` and restic `cmd/restic/main.go:77-106`. It also identifies hidden `init()` registries as a cost. Sprint 36 therefore keeps a closed set of in-tree QA checks.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit argv, permission gates, path containment, private storage, and secret-safe values through restic `internal/backend/shell_split.go:45-76`, opencode `internal/permission/permission.go:44-108`, and chezmoi `internal/cmd/gpgencryption.go:151-165`. This supports denying agent shell access and checking identity in addition to permission policy.
- `studies/go-cli-study/reports/final/14-performance.md` supports hard limits and disk-backed long-lived state through age `internal/stream/stream.go:20,195-219`, restic `internal/archiver/file_saver.go:56-58`, rclone `lib/pool/pool.go:17-24,52-53`, and opencode `internal/message/message.go:37-42`. Disk backing alone does not bound storage, so Sprint 36 also freezes output, history, worker, and duration limits.

The key conclusion is that QA needs one sprint-owned semantic state machine and one existing operational run system. They correlate through app-owned acceptance and an opaque writer token, but neither becomes the other.

## Risks

- The writer-token design must integrate with existing durable acceptance without leaking run-control storage types into `internal/sprint`. If composition cannot supply the accepted generation before work starts, stale-writer safety is not proven and runtime QA must remain blocked.
- Filesystem atomic rename protects one file, not the whole attempt. Pointer-last publication and digest reconciliation must be covered by failure tests at every step.
- `flow-state.json` currently embeds detailed review and smoke state. Adding a pointer-only QA projection creates an intentionally different pattern that validation and planning refresh code must preserve.
- Product-owned check descriptors can drift from actual test commands. Mapping records descriptor identity and fingerprint; missing or changed commands invalidate affected shards rather than falling back to shell.
- Some existing checks may write caches, coverage files, temporary files, or generated output despite looking read-only. The allowlist needs an adversarial corpus and before/after identity proof. Unknown checks remain blocked.
- Identity hashing can become expensive on large repositories. The implementation should stream sorted manifests and reuse nothing across attempts unless the cache itself has validated identity and bounds. Correctness wins over a fast but stale cache in Sprint 36.
- A target can change between the pre-attempt identity read and process launch. Revalidate the implementation fingerprint immediately before dispatch, then compare the full manifest after completion. Any remaining race produces drift and blocks promotion.
- Cancellation may race with a final valid shard write. Writer-token, context, and state-transition checks must establish one order and preserve a completed record only when validation and atomic promotion finished before cancellation ownership changed.
- A run can reach a terminal operational result while summary projection persistence fails. Status must show the run fact and invalid or stale QA summary together, then direct explicit recovery. It must not infer QA completion.
- Local shard failure continuation can waste runtime after a broad systemic problem is misclassified as local. Permission, identity, state, writer, and governed-input failures are whole-attempt failures; tests must freeze that boundary.
- Version 1 rejection is safe but gives no automatic recovery for future formats. Documentation must tell users to use a compatible binary or an explicit future migration command, never delete state casually.
- Private `0600` detailed files may differ from current workspace artifact permissions. Packaging, preview, backup, and browser paths must read them through contained app use cases rather than broad static serving.
- Fixed limits may block legitimate large changes. Exhaustion remains truthful and measurable. A later sprint may revise the defaults only with fixture, documentation, and fingerprint changes.
- A current failed or blocked Conformance Review may still feed diagnostic QA, while missing or stale review evidence blocks it. Readiness code must test freshness separately from verdict acceptability.
- The final sprint reasoning must reconcile this architecture with API and frontend area decisions and reference this file. `reasoning.md` does not yet exist.
