# Sprint requirements: Evidence-producing QA and smoke integration

> Project: `ultraplan-go`
> Sprint: `37-evidence-qa-smoke`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint goal

Add isolated, evidence-producing QA with product-owned adjudication, canonical QA reporting, and smoke execution through QA while preserving target immutability and every existing smoke safety and compatibility guarantee.

## Required outputs

| Output | Path | Description |
| --- | --- | --- |
| Isolated investigation service | `../ultraplan-go/internal/sprint/qa_investigation.go` | Creates and validates one contained writable workspace per shard attempt, runs bounded evidence plans, captures results, and proves cleanup without writing to the target checkout. |
| Isolated investigation tests | `../ultraplan-go/internal/sprint/qa_investigation_test.go` | Covers copy and non-Git targets, identity and containment checks, symlink/path escape, denied target writes, bounded execution, cancellation, cleanup failure, and preserved evidence. |
| Generic isolation mechanics | `../ultraplan-go/internal/platform/process/isolation.go` | Provides reusable local copy, process-containment, cancellation, and cleanup mechanics without importing sprint QA types or deciding evidence outcomes. |
| Generic isolation tests | `../ultraplan-go/internal/platform/process/isolation_test.go` | Fault-injects creation, copy, command, cancellation, descendant cleanup, and removal failures against temporary workspaces. |
| Evidence and issue state extensions | `../ultraplan-go/internal/sprint/qa_state.go` | Extends versioned verification state with generated checks, evidence identities, cleanup facts, adjudication, bounded issues, regression candidates, freshness, and canonical QA pointers. |
| Evidence state tests | `../ultraplan-go/internal/sprint/qa_state_test.go` | Covers schema compatibility, strict validation, atomic publication, stale-writer fencing, digests, bounds, invalidation, and recovery for new evidence and issue records. |
| Global QA adjudicator | `../ultraplan-go/internal/sprint/qa_adjudication.go` | Grounds expectations, validates current and contained evidence, rejects flaky or invalid setups, groups root causes, promotes bounded issues, and classifies repair eligibility and regression candidates. |
| Adjudication tests | `../ultraplan-go/internal/sprint/qa_adjudication_test.go` | Covers promotion and rejection decisions, freshness, repeatability, setup validity, containment, evidence sufficiency, root-cause grouping, and model-output distrust. |
| QA orchestration and canonical summary | `../ultraplan-go/internal/sprint/qa.go` | Extends Sprint 36 QA orchestration with writable shard attempts, adjudication, deterministic assessment, canonical `qa.md`, cancellation, resume, and recovery. |
| QA orchestration tests | `../ultraplan-go/internal/sprint/qa_test.go` | Covers bounded writable execution, read-only compatibility, current-state publication, assessment rules, cancellation, resume, cleanup uncertainty, and no production repair. |
| Smoke QA executor adapter | `../ultraplan-go/internal/sprint/smoke.go` | Runs the existing smoke protocol as the `smoke` QA suite while preserving current discovery, containing-suite, execution, evidence, diagnostic, and compatibility behavior. |
| Smoke parity tests | `../ultraplan-go/internal/sprint/smoke_test.go` | Proves equivalent selection, invocation, cancellation, cleanup, evidence validation, verdict, `smoke.md`, and external-harness results for `smoke` and `qa --suite smoke`. |
| Shared QA use cases | `../ultraplan-go/internal/app/sprint_usecases.go` | Exposes typed writable-attempt, evidence, adjudication, issue, canonical QA, smoke-suite, cancellation, and recovery results to every local adapter. |
| Shared use-case tests | `../ultraplan-go/internal/app/sprint_usecases_test.go` | Proves adapter-independent results, durable run correlation, bounded projections, authorization-independent observation, and consistent next actions. |
| CLI and JSON commands | `../ultraplan-go/internal/app/sprint_commands.go` | Adds evidence-producing `qa`, `qa --suite smoke`, focused shard operation, status, cancellation, resume, and stable machine-readable output without adding a planning stage. |
| CLI and JSON tests | `../ultraplan-go/internal/app/sprint_commands_test.go` | Freezes help, flags, output schemas, exit behavior, compatibility aliases, blocked states, and app/CLI agreement. |
| TUI QA view | `../ultraplan-go/internal/tui/qa_view.go` | Presents shard evidence, adjudication, promoted issues, canonical assessment, smoke-suite status, blockers, cancellation, and recovery through app use cases. |
| TUI QA tests | `../ultraplan-go/internal/tui/qa_view_test.go` | Covers keyboard operation, bounded evidence rendering, hostile text, cancellation, dropped-delivery recovery, and parity fixtures. |
| Browser QA handlers | `../ultraplan-go/internal/web/qa_handlers.go` | Maps guarded versioned requests and read-only evidence queries to shared QA use cases and durable run progress without importing sprint internals. |
| Browser QA tests | `../ultraplan-go/internal/web/qa_handlers_test.go` | Covers routes, strict JSON, same-origin and CSRF checks, confirmation, hostile content, reconnect, restart recovery, cancellation, and parity. |
| Sprint QA presentation | `../ultraplan-go/internal/web/templates/sprint.html` | Shows current canonical QA status, evidence, adjudication, bounded issues, smoke compatibility, blockers, and next actions in a useful server-rendered view. |
| Browser QA enhancement | `../ultraplan-go/internal/web/static/js/operations.js` | Adds bounded progressive enhancement for QA evidence and cancellation over existing operation and durable-run APIs without client-side authority. |
| QA command documentation | `../ultraplan-go/docs/cli-reference.md` | Documents writable QA, `qa --suite smoke`, compatibility commands, JSON outcomes, blockers, cancellation, and recovery. |
| QA architecture documentation | `../ultraplan-go/docs/architecture.md` | Documents isolation, product/platform ownership, evidence authority, adjudication, state layout, issue bounds, and smoke integration. |
| QA user workflow | `../ultraplan-go/docs/user-guide.md` | Explains evidence plans, generated checks, adjudication, canonical QA, issues, smoke suites, and safe operator actions. |
| QA browser documentation | `../ultraplan-go/docs/local-web.md` | Documents guarded starts, evidence inspection, canonical status, durable progress, cancellation, reconnect, and recovery. |
| QA recovery documentation | `../ultraplan-go/docs/recovery.md` | Documents failed isolation, target drift, stale evidence, interrupted attempts, cleanup uncertainty, invalid state, and safe restart behavior. |
| QA JSON and state schemas | `../ultraplan-go/docs/phase3-json-schemas.md` | Documents additive evidence, adjudication, issue, assessment, and smoke-suite schemas plus compatibility and migration rules. |
| QA release checks | `../ultraplan-go/docs/release-checklist.md` | Adds isolation, non-mutation, evidence quality, issue audit, smoke parity, cross-surface, race, build, and gated dogfood checks. |
| Canonical QA report | `projects/<project>/sprints/<sprint>/qa.md` | Summarizes the current input fingerprint, coverage, evidence quality, adjudicated issues, assessment, smoke evidence, blockers, and next action. |
| Verification state root | `projects/<project>/sprints/<sprint>/verification/state.json` | Points to the current schema-versioned attempt and stores bounded canonical QA summary, freshness, assessment, artifact digests, and run correlation. |
| Attempt evidence records | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/evidence/<evidence-id>.json` | Stores evidence plan, commands, bounded results, target and workspace identities, generated-patch digest, repeatability, containment, and cleanup facts. |
| Generated verification patches | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/patches/<patch-id>.patch` | Preserves bounded generated tests, fixtures, probes, or smoke-scenario changes as evidence or regression candidates, never as production repair. |
| Adjudication outcome | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/adjudication.json` | Stores deterministic evidence-quality decisions, rejected evidence, root-cause groups, requested follow-up, assessment inputs, and promotion reasons. |
| Bounded issue records | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/issues.json` | Stores only adjudicated evidence-backed issues, repair eligibility, regression-candidate classification, and exact supporting evidence references. |
| Smoke compatibility report | `projects/<project>/sprints/<sprint>/smoke.md` | Retains the current linked smoke summary and canonical-versus-narrow evidence rules when smoke runs as a QA executor. |
| External smoke evidence | `../ultraplan-go-smoke/runs/<run-id>/` | Retains raw run JSON, stdout/stderr, per-test artifacts, identities, and cleanup evidence under the manifest-declared harness authority. |

## Acceptance criteria

- [ ] Sprint 36 has a current acceptable Conformance Review and required smoke evidence proving deterministic mapping, complete changed-path ownership, read-only investigation, cancellation, resume, fingerprint invalidation, synthesis, and target non-mutation before writable QA is enabled.
- [ ] Every writable shard attempt receives a distinct workspace. Before any write, product code proves source identity, target identity, root containment, path and symlink safety, writable-path policy, process containment, and cleanup capability; inability to prove any item produces `blocked` and starts no writable child work.
- [ ] The target checkout, production source, production tests, governed planning inputs, Sprint 36 historical theory records, smoke harness outside manifest-declared paths, and Git state remain unchanged by QA. Before-and-after identity checks detect drift and prevent evidence promotion.
- [ ] Investigators may create only targeted tests, fixtures, probes, smoke scenarios, and bounded experiments inside their assigned workspace. They cannot modify shared production code, repair the target, weaken existing tests or expectations, broaden shard scope, mutate Git, or promote their own claims.
- [ ] Each generated check has a frozen pre-execution plan with the theory and expectation references, confirmation, refutation, and inconclusive conditions, approved paths, explicit argv, bounded environment, timeout, output cap, and cleanup requirements.
- [ ] Commands execute without shell interpolation, use bounded environment forwarding and contained working directories, propagate context cancellation, terminate descendants, and record timeout, cancellation, exit, output truncation, and cleanup truthfully.
- [ ] Evidence records bind the current governed-input, implementation, map, shard, workspace, command, generated-patch, and external-evidence identities. Any stale, malformed, missing, mismatched, uncontained, or cleanup-uncertain identity prevents promotion and a passing assessment.
- [ ] Generated patches and bounded outputs survive workspace cleanup under versioned verification state. They remain evidence or regression candidates and are never applied to the target by this sprint.
- [ ] Writable concurrency is disabled or sequential until the selected isolation mechanism proves independent workspaces and process cleanup under fault and race tests. All concurrency, attempt, command, generated-check, output, storage, duration, retry, and follow-up limits are finite, validated, documented, and fixture-tested.
- [ ] Only the global adjudicator can promote an issue, mark repair eligibility, or classify a regression candidate. A failing command, generated check, investigator claim, or model response alone cannot do so.
- [ ] Adjudication validates expectation grounding, current implementation and input fingerprints, setup validity, containment, confirmation-condition fidelity, repeatability or deterministic sufficiency, flakiness, external evidence identity, severity, root-cause grouping, and evidence sufficiency before promotion.
- [ ] Adjudication retains confirmed, refuted, invalid, inconclusive, blocked, cross-shard, and not-applicable theories plus rejected evidence. Equivalent manifestations are grouped by root cause, and every promoted issue links exact current evidence and states why it is sufficient.
- [ ] Stale, malformed, flaky, ungrounded, diagnostic-only, narrow-only, failed-setup, missing-cleanup, or uncontained evidence cannot yield a passing canonical QA assessment. A narrow passing rerun cannot replace its required containing suite.
- [ ] Product code, not an investigator or summary model, derives the canonical QA assessment from the current Conformance Review, required QA evidence, adjudication, blockers, and containing-suite results. QA cannot overwrite or upgrade the independent Conformance Review verdict.
- [ ] `qa.md` is written atomically only from validated current state and includes the input fingerprint, map and shard coverage, evidence summary, rejected evidence, promoted issues, regression candidates, smoke evidence, assessment, blockers, and next action. A failed new run preserves the last complete report while current failure remains visible.
- [ ] Detailed maps, shards, theories, evidence, generated patches, adjudication, and issues remain schema-versioned under `verification/`. `flow-state.json` stores only bounded status, freshness, assessment, counts, next action, and contained pointers or digests.
- [ ] Unknown major state versions, unsupported migrations, invalid digests, path escapes, partial publication, stale writer tokens, or over-limit records fail closed with actionable recovery. Publication is atomic and dependency ordered, and prior valid state survives injected write failures.
- [ ] `ultraplan sprint <project> <sprint> qa --suite smoke` uses the existing manifest-driven smoke discovery, authoring, selection, containing-suite, invocation, and evidence-validation path rather than a second smoke implementation.
- [ ] `smoke` remains compatible with existing CLI, JSON, TUI, browser, flow-state, `smoke.md`, external run, and issue evidence behavior. Compatibility is not removed in this sprint, even after parity tests pass.
- [ ] Smoke-as-QA preserves manifest protocol validation, bounded authoring paths, explicit argv, environment allowlisting, contained cwd, timeout, cancellation, descendant cleanup, evidence schema and identity checks, diagnostic-only behavior, canonical-versus-narrow distinctions, and review gating.
- [ ] Raw smoke JSON, stdout/stderr, per-test artifacts, and harness issue files stay in manifest-declared external harness roots. UltraPlan stores only validated links, identities, bounded adjudication records, and the current human-readable `qa.md` and `smoke.md` summaries.
- [ ] CLI text/JSON, TUI, browser HTML/JSON, durable run inspection, `qa.md`, and verification state agree on current attempt, fingerprints, shard and evidence status, adjudication, promoted issues, assessment, smoke suite, blocker, cancellation, recovery, and next action for shared fixtures.
- [ ] Every runtime-backed writable QA or smoke-suite run uses Sprint 35 durable acceptance, fencing, liveness, progress, cancellation, redaction, replay, retention, reconciliation, and terminal arbitration. No adapter creates a second operation or event authority.
- [ ] Cancellation stops new scheduling, reaches active runtime and child processes, performs bounded cleanup, and preserves completed evidence. Disconnecting a browser or SSE client only stops observation; cleanup uncertainty is explicit and cannot pass.
- [ ] Browser starts require the existing confirmation, same-origin, CSRF, request-bound, and authorization checks. Server-rendered pages work without JavaScript, hostile evidence is escaped, and refresh, session rotation, reconnect, observer restart, and replay gaps recover from authoritative state.
- [ ] Deterministic offline tests cover all isolation, evidence, adjudication, state, smoke parity, cancellation, recovery, security, compatibility, and cross-surface cases. Race tests cover parallel attempts, stale writers, cancellation, cleanup, and terminal arbitration.
- [ ] A gated real-repository dogfood run produces contained generated evidence and at least one adjudication rejection or promotion audit, leaves the target identity unchanged, cleans its workspace, and runs smoke through QA with evidence identical in authority and guarantees to the compatibility command. Missing prerequisites produce `blocked` and do not satisfy the gate.
- [ ] `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` pass in `../ultraplan-go`.

## Non-goals

- Production repair, applying generated patches to the target, `repair --issue`, frozen repair packets, manual repair, automatic repair, repair cycles, progressive repair reverification, or convergence decisions; these belong to Sprint 38.
- Letting investigators, generated checks, runtime exit status, model prose, or external harness records directly promote issues or set the canonical QA assessment.
- Converting a regression candidate into a permanent product test or fixture as part of QA execution.
- Replacing, renaming, or removing `review`, `review.md`, `smoke`, `smoke.md`, current review/smoke verdict rules, or existing Phase 3 clients.
- Adding QA or repair to `PlanningStage`, changing the canonical planning artifact order, or merging Conformance Review, empirical QA, adjudication, smoke execution, and repair into one stage.
- A general-purpose issue tracker, assignment system, remote issue synchronization, mutable test-generation service, generic sandbox framework, workflow engine, scheduler, broker, daemon, plugin system, or remote worker protocol.
- Fixing Sprint 35 or Sprint 36 findings. Findings and failed checks are QA inputs; production changes require later governed repair scope.
- Global content identity, provenance, retrieval, RAG, embeddings, knowledge graphs, alternate authored-artifact persistence, hosted service, cloud authority, public exposure, multi-user collaboration, or remote execution.
- Automatic Git mutation, commits, resets, cleanup, worktree repair, or any assumption that the target is clean, committed, or Git-backed.

## Constraints

- `internal/sprint` owns QA investigation policy, evidence semantics, adjudication, issue promotion, canonical assessment, `qa.md`, verification state, and smoke compatibility. Generic platform packages may provide copy, sandbox, and process mechanics but must not import sprint QA types or decide product outcomes.
- CLI, TUI, and HTTP remain adapters over typed `internal/app` use cases. `internal/web` must not import `internal/sprint`, parse CLI output, invoke investigators or the smoke harness directly, or persist alternate QA truth.
- A writable investigator gets one fresh isolated workspace for one shard attempt. The shared target is always read-only to QA, and writable work fails closed when isolation, identity, containment, path safety, process cleanup, or cleanup completion is uncertain.
- Isolation must support local targets that are dirty, uncommitted, or not Git repositories. Git worktrees may be an optimization only when exact target identity and uncommitted state are represented safely.
- Agentwrap/OpenCode remains the agent runtime boundary, and the existing platform process seam remains the external-process boundary. Commands use explicit argv, contained working directories, bounded environment forwarding, finite timeouts, context cancellation, and descendant cleanup; no command is parsed from Markdown or model prose.
- All investigator permissions are default-deny. Allowed writes are limited to the assigned isolated workspace and declared verification-state outputs; smoke authoring and evidence writes remain limited to manifest-declared external harness roots.
- Detailed QA state is schema-versioned, private, path-contained, bounded, digest-linked, atomically published, and fail-closed on unknown major versions. Verification-scoped IDs retain explicit compatibility or migration behavior and do not claim to be global content identity.
- `flow-state.json` remains a bounded projection. Sprint 35 run control remains authoritative for operational acceptance, ownership, fencing, event order, liveness, cancellation, cleanup, and terminal observation; neither store absorbs the other.
- The external smoke harness remains authoritative for suites, raw runs, per-test artifacts, and harness issue files. QA integration cannot copy raw evidence into the project workspace, bypass manifest discovery, or treat a narrow or diagnostic run as canonical evidence.
- `review` remains read-only Conformance Review. QA may combine its current verdict as an assessment input but cannot mutate review evidence or change the review verdict.
- Product code validates structured investigator and adjudicator output and computes promotion and assessment decisions. Model output is untrusted input and cannot bypass deterministic rules.
- Normal tests are deterministic and offline with temporary workspaces, fake runtimes, fake process boundaries, and harness fixtures. Real runtime, process, browser, and smoke evidence is gated, bounded, redacted, and reported truthfully.
- Automatic Git mutation remains prohibited.

## Dependencies

| Prior sprint / output | Required for | Notes |
| --- | --- | --- |
| Sprint 36 implementation, current Conformance Review, and smoke evidence | Writable QA admission | Sprint 36 is planned but not yet executed or reviewed when these requirements were written. Sprint 37 must not enable writable work until deterministic mapping, full changed-path ownership, read-only investigation, cancellation, resume, invalidation, synthesis, and target non-mutation pass the roadmap gate. |
| Sprint 36 verification state under `projects/<project>/sprints/<sprint>/verification/` | Maps, shards, theories, attempts, synthesis, fingerprints, and bounded state | Extend the versioned state and preserve historical outcomes; do not rewrite read-only theory history to manufacture current evidence. |
| Sprint 35 review, `projects/ultraplan-go/sprints/35-durable-run-observability/review.md` | Operational risk inputs | Its `pass_with_findings` result remains evidence for QA mapping and adjudication, not implicit repair authorization. |
| Sprint 35 run-control implementation in `../ultraplan-go/internal/runcontrol/` | Durable execution, cancellation, replay, and recovery | Reuse the selected same-host authority and fencing model. Do not create a QA-specific operation registry or lifecycle store. |
| Existing review and smoke implementation in `../ultraplan-go/internal/sprint/` | Conformance input and smoke-as-QA parity | Preserve current review gating, harness authoring/discovery, containing-suite selection, invocation, evidence validation, flow state, and Markdown summaries. |
| External smoke harness manifest, `../ultraplan-go-smoke/ultraplan-smoke.json` | Suite discovery, bounded authoring, execution protocol, and evidence roots | The manifest, not prose, defines executable argv, protocol version, capabilities, allowed mutation roots, and durable evidence locations. |
| Current governed sprint inputs, execute evidence, implementation fingerprint, and Conformance Review | Grounded and fresh QA evidence | Missing or stale inputs block current adjudication and assessment; they are never inferred from old attempts. |
| Sprint 33 code-context and Sprint 34 shared-context behavior | Grounded investigator and adjudicator context | Reuse the exact validated context pack and current source references without adding retrieval, caching, or a parallel context manifest. |
| `../ultraplan-go/docs/plans/integrated-roadmap.md` | Phase 5 sequencing and Sprint 37 gate | Current implementation-repository plan was read directly as required by the project roadmap. |
| `../ultraplan-go/docs/plans/post-execution-qa-and-repair-loop.md` | Isolation, evidence, adjudication, artifact, and smoke migration design input | Current plan was read directly. This sprint adopts its evidence-producing and adjudication scope but defers production repair. |
| Project index, roadmap, PRD, TRD, and Architecture | Product scope, ownership, state authority, compatibility, and release criteria | These remain authoritative when implementation-repository plans offer optional mechanisms. |
| Agentwrap/OpenCode and existing platform process boundary | Bounded agent and child-process execution | UltraPlan supplies typed permissions, explicit evidence plans, validation, adjudication, and product-owned state without duplicating runtime supervision. |

## Review expectations

| What | How verified |
| --- | --- |
| Sprint 36 entry gate | Inspect current review and smoke artifacts plus deterministic map, changed-path coverage, read-only permission, resume, invalidation, synthesis, and target-identity evidence before testing writable admission. |
| Isolation and target immutability | Adversarial temporary-workspace and dogfood tests cover dirty and non-Git targets, copies, symlinks, path escapes, concurrent attempts, target drift, creation failure, cancellation, descendant processes, cleanup failure, and before/after identity. |
| Investigator boundaries | Runtime-request and process fixtures attempt production writes, governed-input changes, expectation weakening, shell indirection, environment escape, scope growth, Git commands, direct promotion, and repair; every attempt is denied or blocks the attempt. |
| Evidence integrity | Golden records trace each result to a frozen evidence plan, current fingerprints, generated patch, exact argv, bounded output, repeatability, workspace identity, containment, cleanup, and external hashes where applicable. |
| Adjudication quality | Table-driven fixtures cover grounded valid evidence and stale, malformed, flaky, subjective, unrelated-failure, diagnostic, narrow, uncontained, cleanup-uncertain, and contradictory cases; only valid current evidence can promote an issue. |
| Issue and regression-candidate bounds | State and audit tests prove global-only promotion, root-cause grouping, exact evidence links, separate severity and repair eligibility, bounded issue counts, retained rejected evidence, and no repair execution. |
| Canonical QA and state authority | Golden Markdown and JSON fixtures prove deterministic assessment, complete `qa.md`, schema compatibility, pointer-only flow state, dependency-ordered atomic writes, digest checks, stale-writer fencing, and prior-state preservation. |
| Smoke parity | Run identical fixture and gated harness selections through `smoke` and `qa --suite smoke`; compare protocol, selected containing suite/tests, argv, environment, timeout, cancellation, cleanup, run ID, evidence links, verdict, flow projection, and `smoke.md`. |
| Smoke safety and authority | Manifest, authoring-scope, path, process, evidence-schema, canonical-versus-narrow, diagnostic-only, review-gate, and external-root tests prove no guarantee or authority moved into QA. |
| Cross-surface agreement | Inspect one shared fixture through app DTOs, CLI text/JSON, TUI, browser HTML/JSON, durable run detail, `qa.md`, and verification state; canonical fingerprints, evidence, issues, assessment, blockers, and next actions must match. |
| Cancellation and recovery | Fault injection cancels queued and active writable and smoke attempts, drops observers, restarts the server, expires sessions, loses owners, corrupts state, and fails cleanup; no case infers success or loses completed valid evidence. |
| Browser security and accessibility | Handler/template tests cover confirmation, same-origin and CSRF policy, strict decoding, escaped hostile evidence, no-JavaScript snapshots, keyboard/focus behavior, reduced motion, mobile layout, bounded rendering, reconnect, and restart. |
| Documentation and compatibility | Review CLI help, JSON schemas, architecture, user workflow, browser behavior, recovery, release checklist, `review`/`review.md`, and `smoke`/`smoke.md` against executable fixtures and exact state paths. |
| Release and dogfood | Run offline, race, vet, build, and diff checks, then inspect the gated real-repository isolation record, adjudication audit, clean target identity, workspace cleanup, smoke parity, external evidence, redacted diagnostics, and durable run correlation. |
| Scope exclusions | Diff and dependency review confirm no product repair, generated-patch application, compatibility removal, planning-stage expansion, general issue tracker, content identity, retrieval, alternate product persistence, hosted operation, remote worker, or Git mutation. |
