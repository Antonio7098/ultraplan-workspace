> **Inputs Used:** `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `projects/ultraplan-go/sprints/38-bounded-repair/code-context.md`, `projects/ultraplan-go/sprints/38-bounded-repair/sprint-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/technical-handbook.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/api-design.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `internal/sprint/qa_types.go`, `internal/sprint/qa_state.go`, `internal/sprint/state.go`, `internal/sprint/verification_lock.go`, `internal/platform/process/isolation.go`, `internal/app/durable_operations.go`, `internal/app/operation_runner.go`

# Architecture reasoning

Sprint 38 adds a production mutation to the existing verification system. The smallest honest design keeps repair inside `internal/sprint`, reuses run control and process isolation as mechanisms, and adds one product-owned protocol from an adjudicated issue to a terminal repair result. It does not add a repair package tree, plugin system, second automatic engine, new workflow framework, or alternate persistence authority.

## Area Decisions

### Repair remains one sprint-owned verification protocol

`repair` remains a `VerificationPhase`, separate from `PlanningStage`, Conformance Review, and QA. `internal/sprint` owns repair admission, packet construction, confirmation validation, scope policy, production application, progressive reverification, progress calculation, cleanup interpretation, manual-proof qualification, and semantic outcomes. These rules stay in focused files in the existing package, led by `qa_repair.go` and the repair records in `qa_types.go` and `qa_state.go`.

`internal/platform/process` continues to own bounded tree copy, identity, comparison, explicit executable and argument execution, process-tree cancellation, and workspace removal. It returns facts. It must not know whether an issue is eligible, whether a changed path is authorized, whether progress occurred, or which repair outcome applies. `internal/runcontrol` remains authoritative only for durable acceptance, owner liveness, cancellation routing, event order, fencing, and operational terminal state. `internal/app` maps shared use cases and operations to adapters. CLI, TUI, and web code do not read repair files or derive outcomes.

The existing module architecture therefore fits with a small refactor, not a new package hierarchy. The refactor separates durable acceptance from dispatch and gives repair focused domain types and persistence methods. Pure validators and outcome functions remain concrete package functions. Interfaces are limited to volatile runtime, process, clock, and persistence failure seams already needed by tests.

### The authority chain is fixed and cannot be skipped

Every mutable run follows this order:

1. Load the current digest-validated QA state, adjudication, issue, evidence, map, review, smoke, target identity, check catalog, governed inputs, and effective policy.
2. Build one normalized packet from one current repair-eligible issue. Validate all referenced evidence and commands, freeze finite production scope and limits, and publish `issue-packet.json` immutably.
3. Show the packet through a bounded app projection. Packet preparation is runtime-free and cannot mutate the target.
4. Revalidate packet freshness, then durably accept and claim the operation without launching a goroutine or runtime.
5. Build and publish `confirmation.json` over the packet digest, target identity, canonical request, mode, effective limits, governed fingerprints, durable run ID, operational attempt ID, and fencing generation.
6. Mark the accepted operation running, start its ownership controller, acquire the sprint mutation lease, and revalidate confirmation, target, policy, manual-proof requirements, and writer ownership.
7. Create a bounded isolated copy, let the repair runtime edit only that copy, derive and retain a canonical proposal, validate it, and cleanly stop the runtime.
8. Recheck target identity and writer ownership. Apply the validated proposal through product code, then compare actual production changes with packet scope.
9. Run the fixed reverification ladder in order. Record skipped wider gates when a narrow gate does not pass.
10. Cancel remaining work, prove process-tree and isolated-workspace cleanup, transition through terminalization, release the mutation lease, and publish exactly one terminal semantic result.
11. For a qualifying manual result only, publish or refresh the bounded manual-proof pointer. Finish the durable operation separately.

The current `durableOperationManager.AcceptOperation` starts its ownership goroutine immediately after claim. Repair needs acceptance split into `accept/claim` and `dispatch`: acceptance supplies the identity used by confirmation, while dispatch marks the run active and starts ownership control only after confirmation is durable. If confirmation publication fails, the accepted operation closes with a persistence failure and no repair worker starts.

### Packet authority is immutable and assembled from current QA records

The promoted `QAIssue` does not currently contain the full repair contract. Packet preparation joins only current, digest-validated facts from the QA map, shard and theory records, accepted evidence, adjudication, assessment, approved checks, review, containing smoke, governed inputs, and target identity. Caller text cannot supply a reproducer, expectation, command, path, or acceptance rule that is absent from those records.

Packet normalization sorts and deduplicates all finite collections, uses repository-relative slash-separated paths, and rejects empty, absolute, escaping, linked, generated-evidence, test, governance, state, configuration, and Git-control paths. Product policy computes `forbidden_paths`; the model cannot reduce that set. `allowed_paths` is the finite intersection of adjudicated production scope and product-safe path classes. Any ambiguity blocks preparation instead of asking the runtime to decide.

Repeated preparation with identical current inputs returns the current prepared packet and digest. A new attempt after confirmation or terminalization receives a new repair-run identity even if content happens to match. Packet bytes never change in place. A changed issue, QA attempt, target, review, smoke, policy, check catalog, governed input, or evidence digest makes the prepared packet stale and requires a new packet and confirmation.

### Confirmation is single-use durable authority

Confirmation is not adapter state. The immutable confirmation record is the authority, and its digest is also the durable operation deduplication key. A manual confirmation authorizes manual mode only. Automatic mode requires its own explicit mode and opt-in facts and cannot be derived from configuration or a prior manual confirmation.

The repair service consumes confirmation once when it enters mutable work. Replay can return the correlated canonical run but cannot create a new owner or mutation. Resume requires a new confirmation for the same packet, mode, target lineage, policy, original maxima, consumed counters, deadline, and latest proven cycle boundary. Resume can lower no persisted fact, reset no counter, and repeat no committed apply.

### Repair state extends the strict verification store

Detailed records live only under the required `verification/attempts/<qa-attempt-id>/repairs/<repair-run-id>/` tree. Packet, confirmation, cycle proposal, scope, reverification, cleanup, and result records are private, schema-versioned, bounded, unknown-field rejecting, path-contained, symlink rejecting, digest-addressed, and immutable once committed. `verification/repair-state.json` is the only mutable current repair pointer. It contains phase, freshness, correlation, counters, blocker, next action, and digest-bound references, not evidence bodies.

Repair publication follows the existing QA pattern: write immutable detail first, verify its digest, publish current repair state second, and publish the bounded flow summary last. Writer ownership is checked before every record write and every rename. Canonical current files use snapshot and rollback behavior for in-process publication failures. A crash may leave unreferenced immutable evidence, but it cannot make that evidence current. Recovery may retain or prune such evidence after validation; it never treats file presence as success.

`flow-state.json` gains an optional bounded repair summary with phase, freshness, mode, cycle counts, outcome, blocker, next action, and a digest-bound pointer to `verification/repair-state.json`. It does not gain packets, paths, commands, findings, or event history. Because flow-state decoding is strict and existing workspaces persist it, this change requires one explicit schema-version increment and a migration that initializes the repair summary as absent. `SaveFlowState` must preserve existing review, smoke, QA, and repair summaries when another writer updates its own concern. There is no dual-write compatibility format.

### The proposal is isolated and production application is product-owned

The repair runtime receives a private isolated copy and packet-approved context. It never receives the production root as a writable path. Product code compares the original isolated baseline with the runtime result and creates a canonical, bounded unified patch. It rejects NUL data, unsupported binary changes, links, special files, hard links, ambiguous paths, malformed hunks, mode changes outside policy, and unsupported rename semantics. A rename may be represented only as an independently authorized delete plus add; otherwise it is rejected.

Product code parses the canonical patch again before apply. It validates target identity, changed-file and byte totals, every old and new path, allowed and forbidden sets, writer fence, mutation lease, and expected pre-image digests. `proposal.patch` is published before production mutation. Runtime output is evidence, not an apply instruction.

The production apply boundary uses direct filesystem operations, never Git, a shell, a formatter, or runtime tools. It stages replacement bytes privately, records pre-images needed for compensation, rechecks each destination without following links, and uses same-directory atomic renames for individual files. On an in-process failure it restores committed pre-images and reports rollback failure separately. Multi-file replacement is not globally atomic on a normal filesystem, so state records an apply journal and recovery treats a crash or uncertain compensation as `escalated`; it never retries or completes an apply from file presence alone.

After apply, product code recomputes the full target identity and the actual changed path and byte set. The actual set must equal the applied operation set, remain a subset of `allowed_paths`, and remain disjoint from all protected classes. A mismatch ends mutable work before reverification and escalates. Pre-existing unrelated changes are part of the frozen target identity, are reported, and are never absorbed or reverted.

### Two ownership controls remain necessary

The sprint mutation lease serializes production and governed-state mutation across processes. The run-control fence identifies the exact accepted owner allowed to publish. Neither replaces the other. Every proposal, apply, command, cleanup, state publication, result publication, and resume boundary checks both where applicable.

The existing process-ID verification lock remains a mechanical exclusion mechanism, but repair state adds durable correlation to run ID, operational attempt ID, and fencing generation. Lock acquisition must re-read canonical repair state after exclusive acquisition, so a `terminalizing` state blocks a new mutation even after the prior owner releases the file lock.

Terminalization uses this order to make lock cleanup truthful:

1. While holding the mutation lease, stop work, clean descendants and isolation, recheck target identity, and publish a writer-fenced `terminalizing` repair state with no final result.
2. Release the owned mutation lease and verify release. If release fails, publish cleanup uncertainty and an `escalated` result while ownership is still identifiable.
3. After successful release, publish immutable cleanup and result records under the run-control fence, then publish terminal repair state and flow summary. New mutation acquisition must acquire the lease and then reject the `terminalizing` state, so it cannot race this publication.
4. If the process dies during terminalization, recovery may publish interruption, cleanup uncertainty, or escalation from proven facts. It cannot publish `verified` or `verified_with_findings` on behalf of the dead owner.

A stale or cancelled fence cannot apply, publish progress, install a result, release another owner's lock, or refresh manual proof. Result creation uses immutable create-or-compare semantics. An existing result wins over late runtime completion, duplicate cancellation, or recovery.

### Reverification is sequential and uses frozen commands

The ladder is one ordered product workflow: exact reproducer, affected primary shards, linked theory confirmation and refutation checks, boundary and follow-up shards, containing approved QA and smoke suites, then a focused Conformance Review delta. Gates do not run concurrently. This preserves the first authoritative failure, keeps command and target identities easy to audit, and avoids doing wider expensive work after a narrow failure.

Every command comes from the packet's immutable approved descriptors and runs through `platform/process` with explicit executable and arguments, contained working directory, allowlisted environment, timeout, output limit, cancellation, target checks, and process-tree cleanup. The runtime cannot replace a command or reinterpret an exit. Truncation, cleanup uncertainty, stale identity, or missing executables blocks or escalates according to product rules; it is not a failed assertion unless deterministic evidence proves the assertion false.

The Conformance Review delta reuses the independent reviewer capability but returns a bounded structured delta into `reverification.json`. It does not replace `review.md`, edit the existing verdict, or claim global conformance. Repair derives `verified` only from the complete ladder and current target, never from a successful reproducer alone.

### Automatic mode is a loop around the manual cycle, not another engine

One cycle function owns proposal, validation, evidence publication, apply, actual-scope comparison, reverification, and cleanup. Manual mode invokes it once and can never emit `stalled`. Automatic mode invokes the same function repeatedly only while current manual proof, confirmation, scope, target, policy, ownership, deadlines, and every lower-only budget remain valid.

Consumed counters and the absolute deadline are committed at each immutable cycle boundary before another cycle can be scheduled. Checks occur before proposal, apply, every command, publication, retry, and resume. Automatic progress is computed by product code from normalized issue sets, severity, exact-failure state, and bounded new facts. Runtime prose, reordered evidence, and repeated explanations are not progress. Repeated patch digests, repeated target identities, unchanged issue sets, reopening bounds, cycle limits, or stagnation limits stop before another mutation.

`verification/manual-repair-proof.json` is a bounded digest-bound pointer, not hand-authored proof. It may point only to a manual run with a real production apply, complete ladder, proven cleanup, current target, and `verified` or `verified_with_findings` result. Its protocol fingerprint covers repair schema, executable code identity, immutable maxima, effective policy, isolation capability facts, approved-check catalog, governed inputs, runtime policy, and qualifying result digest. Automatic admission recomputes each component and reports exact mismatches. A failed proof-pointer publication leaves the manual result truthful but keeps automatic mode unavailable.

### Operational lifecycle and semantic result stay separate

Run control records execution facts. Repair records product facts. A normally completed repair protocol, including a semantic `failed`, `blocked`, `escalated`, or `stalled` result, may close its durable operation as completed because the operation successfully produced its authoritative product result. CLI and adapter exit mapping still treats non-verified semantic outcomes as non-success. Infrastructure failure before a semantic result closes the durable operation as failed, interrupted, cancelled, or cleanup uncertain without manufacturing a repair verdict.

Cleanup uncertainty maps to semantic `escalated` when the current fenced owner or conservative recovery can safely publish that result. Pure cancellation does not become `failed` or `verified`; the durable cancellation state and repair-state cancellation facts remain authoritative, and a semantic result can remain absent. Recovery never changes an existing semantic result. Public status exposes both fields with their owners rather than combining them into one status string.

### Cleanup and recovery are protocol stages

Work cancellation propagates through runtime calls, approved commands, waits, and process groups. Cleanup uses a separate context derived without the work cancellation but bounded by the frozen cleanup deadline. It must prove child process termination, isolated workspace removal, apply compensation state, and lock disposition. Expiry records uncertainty; it does not convert a timeout into cleanup success.

Resume starts only from an immutable cycle boundary. A cycle with a committed apply but no conclusive reverification is not applied again. Recovery validates packet, confirmation, target, governed inputs, owner liveness, state pointers, apply journal, process facts, isolation path, counters, deadline, and locks. It may reuse valid immutable evidence, mark interruption, finish cleanup, or require human action. It cannot adopt a dead worker, reset limits, infer successful apply, infer successful cleanup, or infer verification.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Keep repair in `internal/sprint` with focused files | Preserves existing product ownership and keeps policy beside QA state. | The sprint package grows and needs disciplined file-level cohesion. | A new generic workflow or repair platform package would move product outcomes into infrastructure and create reverse knowledge. |
| Split durable acceptance from dispatch | Lets confirmation bind the accepted run while preserving confirmation-before-goroutine. | Changes the existing durable manager lifecycle and adds an accepted-but-unconfirmed recovery boundary. | Current immediate dispatch starts ownership control too early; confirmation before acceptance cannot bind durable identity. |
| One fixed authority chain | Makes every production mutation traceable and reviewable. | Fewer shortcuts and more persistence boundaries. | Adapter-specific or automatic-specific paths could bypass confirmation, fencing, cleanup, or outcome rules. |
| Immutable detail plus mutable current pointer | Crash leftovers cannot become current and evidence stays auditable. | More files and digest checks than one mutable JSON document. | One large mutable state file makes partial writes, resume boundaries, and evidence retention harder to prove. |
| Explicit flow-state schema increment | Preserves strict decoding and gives persisted state a defined migration. | Older binaries reject the new state after migration. | Silently accepting unknown fields weakens authority; dual formats create competing state. |
| Runtime edits only an isolated copy | Removes model write authority over production. | Copy, identity, patch derivation, and cleanup add time and disk use. | Direct execute-style target access cannot enforce the Sprint 38 trust boundary. |
| Product-derived and product-parsed patch | Scope and changed-byte facts come from filesystem differences, not model claims. | Strict text patch support rejects some otherwise valid binary or rename repairs. | Applying runtime-provided shell or Git commands expands authority and makes actual scope hard to prove. |
| Compensated per-file apply with an apply journal | Works on ordinary local filesystems without Git and supports failure injection. | Multi-file apply cannot be crash-atomic; uncertain crashes escalate. | Pretending several renames form one transaction would create false recovery confidence. Replacing the whole repository root is too disruptive. |
| Sequential reverification | Preserves widening order and first-failure meaning. | Slower than parallel checks. | Parallel gates could run wider work after a narrow failure and obscure target identity changes. |
| Same cycle implementation for both modes | Prevents automatic mode from drifting from manual safety rules. | The cycle result must carry enough facts for looping and stopping. | A separate automatic engine duplicates the highest-risk policy and apply code. |
| Separate operation lifecycle and repair outcome | Keeps cancellation, liveness, and product verification truthful. | Clients must display two related fields. | A combined status would let runtime completion imply verification or cancellation imply a semantic failure. |
| `terminalizing` state before lock release | Allows lock cleanup to be proved without admitting a competing mutation. | Adds a recovery boundary and acquisition recheck. | Publishing success before release can leave a verified result with an uncertain lock; releasing without a state barrier permits a race. |
| Fixed runtime boundary, no plugin registry | Keeps one auditable runtime and schema policy. | Third parties cannot add repair engines. | A plugin system introduces collision, version, timeout, and trust problems without a Sprint 38 use case. |

## Evidence

The report findings below are external evidence. The Sprint 38 conclusions are architectural inferences constrained by the governed requirements and current code; the studies do not prove UltraPlan-specific outcomes by themselves.

| Area | Report or implementation finding | Sprint 38 conclusion |
| --- | --- | --- |
| Module ownership | `studies/go-cli-study/reports/final/01-project-structure.md` finds thin entrypoints and one-way imports across mature Go CLIs. `projects/ultraplan-go/docs/ARCHITECTURE.md` already assigns repair semantics to `internal/sprint`. | Keep policy in `internal/sprint`, composition in `internal/app`, and adapters presentation-only. |
| Dependency seams | `studies/go-cli-study/reports/final/03-dependency-injection.md` favors explicit composition and volatile-boundary interfaces, while warning about globals and oversized dependency containers. | Inject runtime, process, clock, store hooks, and writer fencing at existing construction points. Keep pure repair decisions concrete. |
| Typed terminal facts | `studies/go-cli-study/reports/final/05-error-handling.md` shows typed errors and behavioral categories preserve recovery and exit information. | Use closed repair outcomes, typed blockers, and explicit recovery actions. Do not parse messages or runtime prose to decide outcomes. |
| Failure injection | `studies/go-cli-study/reports/final/06-io-abstraction.md` finds that injectable filesystem, process, and terminal boundaries make side effects testable. | Keep production apply and persistence behind narrow test seams or hooks, while avoiding a generic virtual filesystem. |
| Cancellation lifetime | `studies/go-cli-study/reports/final/07-state-context.md` supports root context propagation and a separate bounded cleanup context, citing restic's delayed cleanup context. | Cancel work first, then prove cleanup under its own finite context. Never use an unbounded background cleanup. |
| Concurrency ownership | `studies/go-cli-study/reports/final/08-concurrency.md` favors localized launch sites, explicit waits, cancellation, and bounded work. | The durable runner is the only repair launch site. Reverification stays sequential, and every worker or child has one owner and bounded join. |
| Canonical facts versus progress | `studies/go-cli-study/reports/final/10-logging-observability.md` favors structured correlation and separation of diagnostics from user output. | Correlate QA attempt, repair run, operation run, cycle, fence, and runtime IDs. Events remain sanitized observations; status reloads canonical records. |
| Test layering | `studies/go-cli-study/reports/final/11-testing-strategy.md` supports table tests, real command-path integration, fakes, and selective golden fixtures. | Use pure decision tables for scope, progress, and outcomes; real filesystem and process tests for apply, cleanup, and races; interface fixtures for public projections; real runtime only for proof gates. |
| Extension restraint | `studies/go-cli-study/reports/final/12-extensibility.md` documents subprocess isolation but also registry collisions, plugin timeout gaps, and version burden. | Reuse the fixed runtime boundary for isolated proposal generation. Do not add a repair plugin or operation registry. |
| Mutation trust | `studies/go-cli-study/reports/final/13-security.md` supports explicit permission gates, argument arrays, private temporary storage, schema validation, and redaction. | Treat confirmation, path scope, strict patch parsing, private records, explicit argv, and fail-closed unknown state as mandatory mutation controls. |
| Bounded work | `studies/go-cli-study/reports/final/14-performance.md` favors lazy initialization, bounded concurrency, streaming, incremental state, and finite retention. | Keep preparation runtime-free, enforce all automatic budgets before scheduling, persist counters incrementally, page public collections, and cap retained cycles and output. |
| Existing strict state | `internal/sprint/qa_state.go` already rejects symlinks and unknown fields, enforces private permissions and limits, verifies digests, checks writer fences before publication, writes detail before pointers, and rolls back canonical files. | Extend this pattern for repair records rather than introducing another store or database. |
| Existing flow projection | `internal/sprint/state.go` preserves review, smoke, and QA summaries across unrelated writes and strictly validates bounded pointers. | Add and preserve one repair summary with an explicit schema migration; keep details under `verification/`. |
| Existing isolation | `internal/platform/process/isolation.go` already provides bounded non-Git copies, link rejection, tree identity, changed-path comparison, and verified removal. | Reuse these mechanics for proposal work and actual-scope facts, but keep repair eligibility and outcomes in `internal/sprint`. |
| Existing durable ownership | `internal/app/durable_operations.go` accepts and claims before domain work, creates a run-control fence, routes cancellation, heartbeats ownership, and proposes one operational terminal result. | Reuse its ownership facts, but split acceptance from dispatch so repair confirmation precedes every goroutine and runtime. |
| Existing shared dispatch | `internal/app/operation_runner.go` is the common runtime-backed path and already injects QA writer fencing before calling the sprint service. | Add repair to this shared path. Do not let CLI, TUI, or web start a repair service independently. |
| Current implementation gap | `projects/ultraplan-go/sprints/38-bounded-repair/code-context.md` reports no repair service, packet schema, patch applicator, reverifier, outcome derivation, or manual-proof publisher, and notes that issue packet facts are spread across current QA records. | Build a focused repair protocol and explicit join validation. Do not pretend current execute or QA paths already authorize production repair. |

## Risks

### Multi-file apply can be interrupted

Ordinary filesystem renames do not make several file changes atomic. The apply journal and compensation path reduce damage but cannot prove rollback after process death or storage failure. Recovery must detect expected-preimage mismatch and escalate without another automatic write. Tests need interruption before each rename, during compensation, and before actual-scope publication.

### Packet derivation can join inconsistent authorities

Issue, theory, evidence, checks, review, smoke, and scope currently live in separate records. A superficially valid issue can point into stale or mismatched records if preparation checks only IDs. Packet construction must verify every digest, parent identity, attempt, target, policy, and expectation link as one snapshot. Missing exact reproducers or unambiguous production paths block preparation.

### Acceptance and confirmation can leave a stranded run

Durable acceptance must precede confirmation because acceptance identity is part of the digest. Persistence failure can leave an accepted but unconfirmed operation. Recovery must terminate that operation without dispatch and without manufacturing confirmation. This boundary needs direct restart tests.

### Terminalization introduces a narrow recovery state

The owner may die after releasing the mutation lock but before publishing cleanup and result. The `terminalizing` state blocks another mutation and gives recovery a clear boundary, but recovery must remain conservative. It may not publish a qualifying manual proof or infer verification even when earlier checks passed.

### Writer fence and file lock can disagree

The file lock proves local exclusion; run control proves durable owner generation. A live lock with a stale fence, or a live fence with an abandoned lock, must block mutation until reconciliation. Code review must reject any write path that checks only one of them or releases a lock by path without comparing ownership.

### Strict patch support may reject legitimate repairs

Rejecting binary patches, implicit renames, mode-only changes, and unsupported encodings narrows what repair can fix. That is intentional for Sprint 38. The service must return a concrete escalation reason rather than silently widening its parser or invoking Git.

### Reverification may lack immutable command coverage

Current approved QA checks may not provide an exact reproducer and every shard, theory, boundary, containing-suite, and delta-review descriptor required by the ladder. Repair must block at packet preparation when the catalog is incomplete. It must not let the runtime invent substitute commands.

### Semantic and operational terminals may drift

If the app maps every non-verified result to an operational failure, clients may treat a completed repair protocol as missing. If it maps runtime completion to verification, it creates false success. Contract tests must compare both authorities across normal completion, semantic failure, cancellation, interruption, persistence loss, cleanup uncertainty, late completion, and recovery.

### Manual proof can become too weak or too fragile

An underspecified fingerprint may admit automatic work after code, policy, or isolation changes. An overly broad machine fingerprint may invalidate proof for irrelevant host changes. The fingerprint inputs must be named, versioned product facts. Automatic status must report component-level mismatch rather than a generic stale flag.

### Detailed state can grow despite per-cycle limits

Proposal, command, and reverification evidence can accumulate across automatic cycles. Enforce per-record, per-cycle, per-run, and retained-cycle limits before writes. Pruning may remove old cycle bodies only after current pointers and result evidence no longer require them; summaries must disclose retention gaps.

### Cleanup timeout can be mistaken for cleanup proof

A bounded wait proves only that waiting ended. Process-tree absence, workspace removal, compensation state, and lock disposition each need affirmative facts. Any missing fact prevents a qualifying result and manual proof.

### The coordinator can become a large state machine

The workflow is necessarily stateful, but one giant method would hide authority checks and make failure boundaries hard to test. Keep one visible top-level protocol with named phase functions and pure helpers for packet validation, scope comparison, progress, stop decisions, and outcome derivation. Do not split those helpers into a general workflow framework or micro-packages.

### Decision

Proceed with a focused sprint-owned repair protocol, a strict extension of verification persistence, and a small durable-operation lifecycle refactor before mutable work. Manual and automatic modes share one cycle and one product-owned apply boundary. Production mutation remains isolated from the runtime, reverification remains ordered, cleanup remains part of truth, and run-control lifecycle never replaces the semantic repair result.
