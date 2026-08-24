# Architecture reasoning: evidence-producing QA and smoke integration

> **Inputs Used:** `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/code-context.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `../ultraplan-go/internal/platform/process/process.go`, `../ultraplan-go/internal/platform/process/process_unix.go`, `../ultraplan-go/internal/platform/process/process_other.go`, `../ultraplan-go/internal/platform/runtime/runtime.go`, `../ultraplan-go/internal/sprint/service.go`, `../ultraplan-go/internal/sprint/state.go`, `../ultraplan-go/internal/sprint/smoke.go`, `../ultraplan-go/internal/sprint/smoke_types.go`, `../ultraplan-go/internal/app/durable_operations.go`

> **Selected Area Template:** Architecture, from `system/reasoning/architecture_reasoning_template.md`.

This area decides ownership, state authority, isolation, identity, publication, adjudication, cancellation, and smoke delegation. The aim is contained empirical QA, not a generic sandbox or repair system.

## Area Decisions

### Key conclusion

The current module architecture still fits, but Sprint 37 needs a small refactor at two existing boundaries before feature work starts. `internal/platform/process` must report stronger containment and cleanup facts than `DirectRunner` currently proves, and the smoke entry point must expose one internal execution path that both command spellings can call without taking the sprint mutation lease twice. No new top-level product module, workflow engine, registry, daemon, or plugin system is justified.

The feature remains one sprint-owned verification workflow composed from four existing authorities:

| Authority | Owns | Must not own |
| --- | --- | --- |
| `internal/sprint` | QA admission, evidence plans, verification identities, adjudication, issue promotion, assessment, `qa.md`, and smoke projection | OS process mechanics or durable operation arbitration |
| `internal/platform/process` | Private workspace creation, bounded copy, filesystem/process containment capability, explicit argv execution, cancellation, descendant cleanup, and removal facts | Shards, theories, evidence sufficiency, issue severity, or verdicts |
| `internal/runcontrol` through `internal/app` | Durable acceptance, owner claim, run ID, fencing, liveness, replay, cancellation request, and terminal arbitration | QA freshness, evidence validity, or canonical assessment |
| Existing smoke code and external harness | Manifest discovery, authoring scope, selection, containing-suite rules, invocation, raw evidence, harness issues, smoke verdict, and `smoke.md` | General QA issue promotion or production repair |

`internal/web` continues to depend only on `internal/app`. CLI and TUI also call the same app use cases. Adapters may render product decisions, but they cannot derive them.

### Entry gate and workflow

Writable admission is evaluated on every run, not once at installation or by sprint number. It requires current Sprint 36 state, a current acceptable Conformance Review, required smoke evidence, complete changed-path ownership, read-only theory synthesis, cancellation and resume evidence, fingerprint invalidation, and target non-mutation evidence. Missing or stale proof returns `blocked` before workspace creation.

The product workflow is:

1. `internal/app` durably accepts and claims the operation, then passes the fenced run context to `internal/sprint`.
2. `internal/sprint` acquires the existing per-sprint mutation lease and validates the Sprint 36 gate and all QA limits.
3. Product code freezes governed-input, implementation, map, shard, theory, and review identities.
4. A read-only investigator returns structured generated files and a frozen evidence plan. It cannot edit files, invoke a shell, promote an issue, or choose an assessment.
5. `internal/sprint` validates the plan and asks `internal/platform/process` for one fresh private workspace.
6. Platform code copies the current target, verifies the copy, and proves the selected filesystem and process containment capabilities.
7. Product code rechecks target identity immediately before the first generated file is materialized. It writes only validated new check, fixture, or probe paths into the copy.
8. Platform code executes exact argv under the frozen limits. The copied source is read-only to the child; only attempt-local cache and temporary directories are writable.
9. Product code captures bounded results, computes the generated patch, preserves patch and evidence records, rechecks target identity, and requests bounded workspace cleanup.
10. The global adjudicator validates the preserved records. Product code alone groups root causes, promotes issues, classifies regression candidates, and computes assessment.
11. Publication commits an immutable attempt generation, then updates the bounded current-state and `qa.md` projections under fencing.
12. Run control records the terminal execution outcome independently from the QA assessment.

Cancellation stops new scheduling, cancels the active runtime or process tree, preserves completed evidence, and enters a separate bounded cleanup context. A cancelled operation can retain a previous canonical assessment but cannot create a new passing one.

### Isolation mechanism

The first writable implementation uses a full local copy, not a Git worktree. A copy represents dirty, uncommitted, and non-Git targets without changing Git or assuming a commit is the source identity. The copy excludes Git administrative data. Product code separately records and rechecks Git control-state identity when the target is Git-backed.

`internal/platform/process/isolation.go` should expose a narrow concrete service, not a generic sandbox framework. Its input contains a source root, private workspace parent, copy limits, read-only execution root, exact writable cache/temp paths, command request, and cleanup deadline. Its output contains paths and mechanical facts such as bytes copied, manifest digest, containment backend, process identity, cancellation, truncation, descendant cleanup, and removal result. It has no sprint imports and no verdict field.

The copy algorithm uses `Lstat` and a deterministic relative-path walk. It rejects devices, sockets, FIFOs, path escapes, cycles, and symlinks whose resolved target leaves the source root. It copies regular-file contents rather than preserving hard links. It preserves only the executable permission bit needed by the target and creates the workspace parent with private permissions. Every destination parent is checked without following an unvalidated link.

The existing process-group implementation is necessary but insufficient for writable admission. On Linux, the initial containment backend uses a private mount, PID, and network namespace through a validated local sandbox launcher such as `bwrap`. System binaries and libraries are mounted read-only, the copied source is mounted read-only, the generated cache/temp roots are mounted writable, capabilities are dropped, network is disabled, and the launcher dies with its parent. `DirectRunner` still owns explicit argv, bounded streams, timeout, signals, and wait behavior. The isolation layer verifies backend availability and cleanup capability before product code writes generated files.

On any platform where the configured backend cannot prove read-only source mounts, writable-path confinement, process-group ownership, descendant termination, and workspace removal, writable QA returns `blocked`. There is no best-effort writable mode. macOS and other platforms may add a focused backend later, but an unproven backend cannot claim parity.

### Product-mediated generated checks

The investigator runtime remains read-only and default-deny. It returns schema-validated generated file records plus the evidence plan. Product code, not the model runtime, materializes those records after all pre-write checks pass. This keeps the agent runtime behind Agentwrap while avoiding direct model-controlled writes to the target or arbitrary host paths.

Generated paths must be absent from the copied baseline or must already have been created by the same attempt. An attempt cannot modify copied production files, existing tests, governed inputs, Git metadata, or files created by another attempt. The initial implementation accepts new Go test files, bounded fixtures, and data-only probes. It rejects generated shell scripts, executable command text, and changes to existing expectations. Product code creates the patch from the validated additions and stores it as evidence; no code path applies that patch to the target.

Each plan freezes:

- expectation and theory references;
- confirmation, refutation, and inconclusive conditions;
- exact new paths and generated bytes or their digests;
- executable and argv as separate fields;
- working directory relative to the copy;
- environment names and product-supplied values;
- command timeout and stdout/stderr caps;
- required structured result or artifact;
- cleanup requirements and repeatability class.

The command allowlist is product configuration, not model output. Shell interpreters, Git mutation, package-install commands, network tools, and commands parsed from Markdown are denied. Runtime caches, `HOME`, and `TMPDIR` are redirected under the attempt workspace so a normal test command does not write elsewhere.

### Verification-scoped identity

Sprint 37 uses scoped SHA-256 identities and does not introduce the later global content identity contract.

| Identity | Inputs |
| --- | --- |
| Governed input | Exact required sprint artifacts, selected contracts/protocols, and current Conformance Review fingerprint |
| Implementation | Canonical target root identity plus a deterministic full tree manifest of relative path, entry kind, mode, link target, size, and regular-file digest |
| Git control state | Resolved Git administration path and digest of its current state when present; absent is valid for non-Git targets |
| Map and shard | Sprint 36 schema version, map fingerprint, shard definition, and theory set |
| Workspace | Attempt ID, canonical workspace root, copied manifest digest, containment backend, and creation facts |
| Evidence plan | All frozen plan fields in canonical order |
| Command | Executable identity, argv, relative cwd, allowlisted environment digest, timeout, and output bounds |
| Generated patch | Exact stored patch bytes |
| External smoke | Existing validated harness, protocol, run, containing suite, evidence path, size, and digest fields |

Attempt IDs derive from the durable run ID plus a bounded local attempt sequence. Evidence, patch, command, and issue IDs are deterministic within that attempt. They are stable enough for references and migrations but are not advertised as workspace-global identifiers.

The target manifest is captured before copy, checked after copy, checked immediately before the first generated write, and checked after command completion. A mismatch stops promotion and marks all evidence from the drifting interval stale. Cleanup still runs. Concurrent user edits are not reverted; they simply prevent the attempt from becoming current.

### Evidence validity and adjudication

Evidence has a local record validator and a relational attempt validator. Local validation checks schema, bounds, enums, exact argv, output facts, patch digest, containment facts, and cleanup fields. Relational validation checks current fingerprints, map/shard membership, workspace identity, command-to-plan linkage, external evidence identity, target before/after equality, and publication generation.

The adjudicator may use an Agentwrap model for bounded structured analysis, but its output is advisory input. Deterministic product code owns the final decision. A command failure, investigator claim, model classification, or harness issue file cannot promote an issue by itself.

One execution is sufficient only when the evidence plan declares a deterministic assertion, setup completed, the decisive structured result is complete, and no contradictory current evidence exists. Timing-sensitive, external, model-dependent, or otherwise variable evidence requires two matching executions in distinct workspaces. Any contradiction or variance marks the evidence flaky and blocks promotion pending bounded follow-up.

Output truncation is always recorded. Truncated human diagnostics may accompany otherwise complete structured evidence, but retained text beyond the cap is never assumed. If confirmation or refutation depends on truncated bytes, the result is inconclusive and cannot support promotion or a passing assessment.

Root-cause grouping uses expectation ID, primary component, normalized failure signature, affected shard set, and exact evidence links. A model may suggest a group, but product code recomputes and validates the key. Equivalent manifestations become one bounded issue. Every promoted issue states severity, why the evidence is sufficient, repair eligibility, regression-candidate status, and exact current evidence IDs.

### Assessment precedence

Execution status and assessment remain separate. The assessment values are `incomplete`, `blocked`, `fail`, `pass_with_findings`, and `pass`.

1. `incomplete` applies before a current admissible attempt has all mandatory evidence, including after cancellation.
2. `blocked` wins when the Sprint 36 gate, current Conformance Review, target identity, containment, cleanup, required evidence, or required containing suite cannot be proven.
3. `fail` applies when valid current evidence promotes a blocker or high-severity issue, or the required containing smoke suite fails.
4. `pass_with_findings` applies only when every mandatory check is current and sufficient, cleanup is certain, containing suites are canonical, and promoted issues are medium, low, or informational.
5. `pass` requires a current acceptable Conformance Review, complete current QA coverage, no promoted issue, no unresolved required rejection, certain cleanup, and passing canonical containing suites.

QA reads the Conformance Review verdict but never rewrites or upgrades it. Diagnostic-only or narrow smoke evidence cannot satisfy a containing-suite requirement.

### State model and publication

Detailed QA state lives under `verification/`; `flow-state.json` receives only a bounded projection. The verification root starts at schema version 2, treating Sprint 36 fixtures as version 1. Migration from version 1 is explicit and additive: it preserves maps, shards, theories, attempts, and synthesis, then adds empty evidence/adjudication/issue fields and marks the assessment incomplete until current Sprint 37 evidence exists. Unknown versions, unknown fields, invalid digests, escaping pointers, over-limit records, and unsupported migrations fail closed.

Each attempt is immutable after terminal publication. The writer uses this order:

1. Write patches and evidence records into a private staging directory.
2. Validate every record and digest from disk.
3. Write adjudication, issues, attempt metadata, and the attempt-local rendered QA report.
4. Sync files and directories, then atomically rename the complete attempt directory into `verification/attempts/<attempt-id>`.
5. Prepare a root `qa.md` projection and `verification/state.json` generation that carry the same attempt ID and report digest.
6. Under the sprint mutation lease and stale-writer token, commit the generation with a small publication journal and retained prior pair.
7. Update the bounded `flow-state.json` projection only after the verification generation is committed.

Readers validate the state pointer, attempt completion marker, and root report digest as one generation. If a crash or injected failure leaves a mismatched pair, recovery ignores the incomplete generation and restores or serves the retained prior complete pair. The failed current attempt remains visible through its immutable attempt record and durable run. This is needed because separate filesystem files cannot be replaced atomically as one operation.

The stale-writer token contains the expected verification generation, current attempt ID, durable run ID, run-control fencing generation, and sprint mutation lease identity. Verification state checks the token but does not renew leases, arbitrate terminal outcomes, or become a second run-control store.

### Smoke delegation

`qa --suite smoke` maps directly to the existing smoke operation. It does not enter writable investigation or parse the manifest itself. The two command spellings share the same durable operation kind, mutation lease, `Service.RunSmoke` orchestration, discovery, authoring, selection, argv, environment, timeout, cancellation, cleanup, evidence validation, verdict, external run ID, flow projection, and `smoke.md` publication.

To avoid nested locking, refactor smoke into one private execution method that assumes the caller already owns the mutation context. The public `RunSmoke` and the QA suite adapter both acquire or receive that same context and call the private method. Adapter-source metadata may explain which spelling started the run, but it cannot alter defaults or results.

After smoke validates its own result, QA stores only validated links, identity fields, bounded adjudication facts, and canonical-versus-diagnostic status. Raw JSON, stdout, stderr, per-test artifacts, and harness issue files remain in manifest-declared external roots.

### Durable operation and adapter boundaries

Every runtime-backed QA or smoke run follows `Accept -> Claim -> Running` in `internal/app/durable_operations.go` before sprint work starts. Product progress is appended through the existing fenced event path. Cancellation is requested against the durable run, acknowledged by the current owner, propagated through the operation context, and finally arbitrated by run control.

QA state records run correlation and product outcome only. It does not duplicate leases, heartbeats, cancellation requests, replay cursors, or terminal arbitration. A persistence failure in the durable event path cancels active work and prevents optimistic success.

The app layer exposes bounded summary and immutable attempt-detail DTOs. CLI text, JSON, TUI, server-rendered HTML, and progressive enhancement all consume those DTOs. No adapter reads verification JSON directly, invokes the process runner, calls the smoke harness, or decides freshness, evidence sufficiency, severity, or next action.

### Configuration and resource limits

One `QASettings` value is resolved with existing precedence, validated after merge, and injected into `sprint.Service`. Request fields may narrow shard scope but cannot widen these limits.

| Limit | Sprint 37 default |
| --- | --- |
| Writable concurrency | 1 |
| Shards in one run | 32 |
| Attempts per shard | 2 |
| Generated checks per attempt | 4 |
| Commands per attempt | 8 |
| Evidence records per attempt | 32 |
| Adjudication follow-ups per run | 8 |
| Promoted issues per attempt | 64 |
| Investigator structured output | 256 KiB |
| One generated file | 256 KiB |
| Generated patch total per attempt | 2 MiB |
| Command stdout | 1 MiB |
| Command stderr | 512 KiB |
| Argv | 128 elements and 16 KiB encoded |
| Forwarded environment | 32 names and 32 KiB encoded |
| One command | 5 minutes |
| One attempt | 20 minutes |
| One QA run | 2 hours |
| Cleanup grace | 10 seconds |
| Command retries | 1 |
| Copy | 500,000 entries and 8 GiB |
| Retained state per run | 256 MiB |
| Retained verification state per sprint | 1 GiB |

The total run duration and retained-state caps override multiplied per-item limits. Hitting any cap stops new scheduling, preserves completed evidence, and produces a typed blocked or incomplete outcome. Existing smoke settings remain authoritative for the smoke suite; QA does not replace them with these defaults.

### Error, observability, and test boundaries

Expected product errors are typed and include blocked admission, unsupported isolation, target drift, invalid plan, denied path, stale evidence, invalid external identity, cleanup uncertainty, state version, digest mismatch, stale writer, cancellation, and publication recovery. Adapters map types to stable machine codes and separate recovery hints from debug detail.

Safe events include run ID, attempt ID, shard ID, evidence ID, phase, containment backend, command digest, duration, truncation, cleanup certainty, assessment, and error code. They omit generated source, raw model output, full argv values when sensitive, environment values, absolute target paths, and external raw evidence. Progress remains observation; durable run and verification state decide outcomes.

Tests follow the ownership boundary:

- `internal/platform/process` fault-injects creation, copy, manifest, command, cancellation, descendant, and removal behavior without sprint types.
- `internal/sprint` table-tests admission, identity, plan policy, relational evidence validation, repeatability, adjudication, assessment, migration, publication, recovery, and smoke projection.
- `internal/app` tests durable acceptance, typed DTO bounds, cancellation, run correlation, and adapter-independent next actions.
- CLI, TUI, and web tests use the same fixtures and freeze only documented text/JSON/HTML contracts.
- Race tests cover concurrent starts, stale writers, cancellation versus completion, cleanup, and terminal arbitration.
- Gated dogfood proves a real contained generated check, one adjudication rejection or promotion audit, unchanged target identity, workspace cleanup, and smoke command parity.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Product-mediated generated files | The model never receives arbitrary host write authority, and existing source/tests cannot be weakened. | Generated checks must fit a structured file schema and cannot use an interactive editor. | Direct agent editing was rejected because permission metadata alone is weaker than product validation. |
| Full local copy | Represents dirty and non-Git targets and leaves the target untouched. | Copy time and disk can be high for large repositories. | Git worktrees were rejected as the primary mechanism because they omit dirty state and exclude non-Git targets. |
| OS-enforced read-only source execution | A test binary cannot write the target or copied production files even if source code is hostile. | Writable QA is unavailable where the backend cannot prove confinement. | Cwd checks and prompt instructions were rejected because a child can use absolute paths. |
| Sequential writable attempts | Simplifies identity checks, storage budgets, cleanup, and stale-writer ordering. | A many-shard run can take longer. | Parallel writable attempts were rejected until race and fault evidence proves independent workspaces and descendants. |
| Scoped digests rather than global IDs | Gives exact lineage without pulling post-Sprint-39 content identity into scope. | IDs are not portable outside one verification model. | A global provenance or content-addressed artifact service was rejected as premature. |
| Deterministic product adjudication with optional model advice | Keeps promotion and assessment testable while still allowing bounded semantic analysis. | Product rules and normalization need substantial fixtures. | Model-owned promotion was rejected because model prose and investigator claims are untrusted. |
| Immutable attempts plus bounded projections | Preserves audit evidence and recovery while keeping status views small. | Publication and migration logic are more involved. | Putting all detail in `flow-state.json` was rejected because it would become large and mix authorities. |
| Publication journal for the state/report pair | Preserves the prior complete pair across multi-file failure. | Recovery has another small state machine to test. | Independent atomic renames were rejected because they can expose mismatched generations after a crash. |
| Direct smoke delegation | Preserves all current protocol and compatibility guarantees. | QA must accept smoke-specific result asymmetry. | A QA-native smoke implementation or generic executor registry was rejected because either would duplicate authority. |
| Hard aggregate limits | Bounds cost, memory, disk, and rendered data even when per-item maxima multiply. | Large investigations may stop incomplete and require a narrower rerun. | Independent limits without a total budget were rejected because their product can still be unsafe. |

The main accepted trade-off is availability for proof. Writable QA blocks on unsupported containment rather than running with weaker guarantees. That is slower to roll out, but it is the only honest fit for target immutability and evidence promotion.

## Evidence

The report findings support the boundaries and mechanisms above; the sprint-specific choices remain inferences from those findings and the governed requirements.

- `studies/go-cli-study/reports/final/01-project-structure.md` finds that mature CLIs keep entry points thin and dependencies one-way, with Restic's domain interface/implementation split as a concrete example. This supports `web -> app -> sprint -> platform`, with policy in `internal/sprint` and mechanics in `internal/platform/process`. The decision to keep one sprint package with focused files, rather than adding subpackages or a generic QA module, is the project-specific inference.
- `studies/go-cli-study/reports/final/02-command-architecture.md` documents factory-built thin commands and Helm's command-to-action delegation. This supports one app operation and compatibility aliases that delegate instead of embedding QA rules in command handlers.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` finds that manual composition roots and focused interfaces improve test isolation, while global mutable state and context service lookup hide dependencies. This supports injected isolation, process, runtime, clock, and store seams. Context carries cancellation and correlation only.
- `studies/go-cli-study/reports/final/04-configuration-management.md` supports explicit precedence and validation after all sources merge. The fixed `QASettings` table and aggregate budget are the sprint-specific application of that finding.
- `studies/go-cli-study/reports/final/05-error-handling.md` supports typed errors, preserved chains, stable exit mapping, and separate user/operational detail. The QA error taxonomy therefore represents target drift, invalid evidence, and cleanup uncertainty as machine decisions rather than message parsing.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` shows that injectable filesystem, backend, process, and output seams make fault paths testable. Restic's functional-field backend mocks and Gh CLI's test IO constructor directly support the proposed fault-injected isolation and adapter tests.
- `studies/go-cli-study/reports/final/07-state-context.md` finds that one root context should reach every I/O layer and cites Restic's separate bounded cleanup context. That supports cancellation propagation plus cleanup that continues long enough to report certainty without changing the cancelled work outcome.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports localized launch sites, bounded fan-out, explicit waits, and timeout-backed shutdown. It also documents fire-and-forget and unbounded-wait failures. The decision to serialize writable attempts is a conservative project inference until isolation passes race tests.
- `studies/go-cli-study/reports/final/09-terminal-ux.md` finds that long work needs interruptible progress and non-TTY fallback. It also shows why presentation cannot be the operation record. QA progress is therefore bounded observation over durable state.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports structured fields and strict stdout/stderr separation. The chosen fields correlate run, attempt, shard, command, evidence, and cleanup without retaining secrets or raw evidence.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports table-driven logic tests, command scenarios, functional fakes, and golden files for compatibility output. The split test plan follows those findings while avoiding snapshots of incidental internals.
- `studies/go-cli-study/reports/final/12-extensibility.md` shows that versioned adapters can preserve an existing subprocess authority. Its Helm examples support the narrow smoke adapter. The same report's registry collision and plugin complexity findings argue against a QA executor registry.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit argv, private temporary directories, default-deny permission decisions, path canonicalization, and secret redaction. The stronger read-only mount and product-mediated write decisions are required inferences because lexical checks alone cannot stop hostile child code.
- `studies/go-cli-study/reports/final/14-performance.md` supports chunked copy, bounded buffers, finite concurrency, and measurement before pooling. The design bounds copy, command output, patches, state, duration, and fan-out before considering low-level optimization.

The live repository confirms the need for the small refactors. `../ultraplan-go/internal/platform/process/process.go` already records explicit argv, bounded output, cancellation, and cleanup fields, but `process_unix.go` declares cleanup complete after process-group signalling without proving a filesystem write boundary or surviving descendants. `process_other.go` cannot prove descendant cleanup after its grace period. `../ultraplan-go/internal/sprint/smoke.go` already centralizes manifest-driven smoke behavior, while `Service.RunSmoke` owns the mutation lease and current state publication. `../ultraplan-go/internal/app/durable_operations.go` already performs durable acceptance, owner claim, event fencing, cancellation routing, and terminal proposal. These are extension points to strengthen or delegate through, not behaviors to recreate.

The project documents add the authority rules that the comparative studies do not address. `projects/ultraplan-go/docs/ARCHITECTURE.md` assigns QA semantics to `internal/sprint`, process mechanics to `internal/platform/process`, adapter composition to `internal/app`, and run lifecycle to `internal/runcontrol`. `projects/ultraplan-go/docs/PRD.md` and `projects/ultraplan-go/docs/TRD.md` require Sprint 36 gating, canonical `qa.md`, versioned detailed state, smoke compatibility, and no repair or Git mutation in Sprint 37.

## Risks

- The live implementation has no Sprint 36 or Sprint 37 `qa*.go` files. Writable admission must stay blocked until current Sprint 36 artifacts and executable validation exist; chronology or fixture-only proof is not enough.
- A full manifest of a large target is expensive. Incremental hashing may be added only if it produces the same identity and cannot miss an entry; exclusions based on convenience would weaken non-mutation proof.
- Git control-state hashing can race with unrelated user Git activity. The correct outcome is target drift and no promotion, not retrying silently or modifying the user's repository.
- The initial Linux containment backend adds a local prerequisite. Missing namespace support, sandbox executable, or descendant proof makes writable QA unavailable. Documentation and status must name the failed capability.
- Read-only source mounts can break tests that write beside source files. Evidence plans must redirect caches and temporary output. Tests that require broader writes are invalid for Sprint 37 rather than grounds to relax containment.
- Product-mediated generated files reduce the attack area but do not make generated content safe. Compilers and test binaries still process hostile bytes, so command execution needs the OS boundary and output/time caps.
- Symlink replacement between validation and copy is a time-of-check/time-of-use risk. The copy implementation must use descriptor-relative operations or revalidate each opened entry; path-string checks alone are insufficient.
- A process-group kill does not prove that every descendant stopped. The containment backend must pair namespace ownership with bounded wait and an explicit post-kill probe. Uncertainty maps to cleanup-uncertain, never success.
- Multi-file publication cannot rely on rename folklore. Recovery tests must inject failure and process death at every journal, report, state, pointer, and directory-sync boundary.
- Schema version 2 assumes Sprint 36 state version 1. If Sprint 36 lands with a different persisted contract, implementation must revise this migration before writing any Sprint 37 state rather than accepting ambiguous records.
- The fixed defaults may be too small for a real repository or too large when multiplied. Tests must verify the aggregate cap, and dogfood must measure copy time, retained bytes, truncation, and blocked rates before defaults are raised.
- Two repeat runs can still share an environmental cause. Distinct workspaces reduce contamination but do not prove independence from the host toolchain. Adjudication must retain setup identity and contradictory evidence.
- Lower-severity promoted issues can yield `pass_with_findings` only when all required evidence and containing suites are valid. A count-only implementation would create false passes.
- Smoke parity can drift through a default, lock, operation kind, or result mapping even when both spellings call similar code. Parity tests must compare invocation, environment, run ID, evidence links, cleanup, verdict, flow state, and exact `smoke.md` digest.
- A failed rerun leaves a current failed attempt beside an older canonical report. Every projection must label both attempt IDs and must not present the retained report as evidence that the new run succeeded.
- Durable run success and QA pass are different facts. Run control may record that orchestration completed while QA assessment is fail or blocked; adapters must not collapse these fields.
- Hostile evidence may contain invalid UTF-8, terminal controls, Markdown, HTML, URLs, and very long lines. State validation must normalize or reject invalid encodings, and each adapter must escape its own output medium.

The architecture decision is to proceed after the two boundary refactors. Complexity is high but required by the proof obligations. The design deliberately removes three tempting sources of accidental complexity: parallel writable attempts, a second smoke implementation, and model-owned promotion.
