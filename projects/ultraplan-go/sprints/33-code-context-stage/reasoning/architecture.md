> **Inputs Used:** `projects/ultraplan-go/sprints/33-code-context-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md`, `projects/ultraplan-go/sprints/33-code-context-stage/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`

# Architecture: Code-Context Stage Vertical Slice

This area decides where the new stage belongs, how it composes with the existing sprint flow, runtime, persistence, application, and web boundaries, and how legacy state remains usable. The scope is Sprint 33's vertical slice only; downstream prompt-prefix reuse, the manual skill, broad documentation, real-runtime dogfood, retrieval, caching, and new workflow frameworks remain outside this design.

## Area Decisions

### Owning module and dependency direction

`internal/sprint` owns the complete code-context product behavior: canonical stage identity and order, prerequisites, artifact mapping, target-resolution coordination, prompt construction, runtime request policy, structural validation, candidate promotion, flow transitions, and legacy-state compatibility. Add focused files to the existing package rather than creating a new package tree or generic stage framework.

The dependency direction remains:

```text
cmd/ultraplan -> internal/app -> internal/sprint
internal/web  -> internal/app -> internal/sprint
internal/sprint -> internal/project + internal/workspace + internal/platform/runtime + existing mechanical filesystem/config helpers
internal/platform/* -> no product modules
```

`internal/project` continues to own implementation-target configuration and resolution rules. `internal/workspace` continues to own workspace paths and embedded-default lookup. `internal/platform/runtime` receives only generic execution data such as prompt, working directory, model, variant, permissions, context, metadata, and expected output; it must not learn the meaning of `code-context`, sprint readiness, or Markdown validation. `internal/web` must not import `internal/sprint`, inspect repositories, parse CLI output, duplicate validators, or persist workflow truth.

### Stage and state model

Add `StageCodeContext` exactly once in the canonical ordered planning sequence:

```text
requirements -> code-context -> sprint-index
```

Every derived list, cumulative-flow calculation, status projection, artifact lookup, and prerequisite check must consume that canonical ordering rather than maintain a conflicting order. This is explicit fixed registration, not a runtime registry or plugin.

Readiness is strict: valid completed requirements make code-context ready; missing, invalid, failed, or incomplete requirements block it. Sprint-index becomes ready only after code-context execution has a successful terminal outcome and the authoritative `code-context.md` passes structural validation. Artifact presence or runtime exit status alone never establishes completion.

Persisted pre-code-context flow state is interpreted through deterministic compatibility logic. Loading preserves all known prior stage outcomes and derives the inserted stage's current projection without a hidden write. Read-only status remains read-only. The next explicit mutating state write emits the current canonical representation atomically. No parallel state file, migration command, or silent dual-write path is introduced.

### Stage service and execution flow

Implement one cohesive code-context stage service inside `internal/sprint`, using existing runtime and persistence seams. Its execution flow is:

1. Validate requirements and resolve project/sprint identity.
2. Resolve the implementation target through existing project/execute mechanisms.
3. Resolve the embedded-or-overridden prompt/template and effective stage model/variant with source metadata.
4. Build a deterministic prompt and generic runtime request with the caller-owned `context.Context` and read-only repository posture.
5. Generate into an isolated candidate location controlled by the sprint service, not directly over the authoritative artifact.
6. Require runtime success, candidate existence, and structural validation.
7. Atomically replace only the sprint-root `code-context.md` and atomically persist the truthful flow transition.
8. On failure or cancellation, preserve the previous valid artifact and record the failed, cancelled, interrupted, or cleanup-uncertain operation outcome.

Prompt preview and dry-run stop before runtime construction that has side effects, candidate creation, artifact replacement, or flow-state mutation. Runtime resolution may be lazy so status and preview remain cheap, but configuration selection and validation must still be deterministic and attributable before execution.

### Artifact and mutation boundary

The single authoritative output is `projects/<project>/sprints/<sprint>/code-context.md`. The runtime may read the resolved implementation repository, but neither the runtime nor UltraPlan may mutate implementation source, tests, Git state, governed inputs, or unrelated sprint artifacts during this stage.

Structural validation is centralized in `internal/sprint` and checks required sections, at least one exact language-tagged source excerpt, repository-relative contained paths, rationale, and well-formed optional ranges and symbols. It rejects placeholders, absolute or escaping paths, malformed ranges, missing fences, and absent or empty output. It does not claim that selected excerpts are semantically complete or that the pack replaces live repository inspection.

Use real filesystem behavior for containment, permissions, same-directory atomic rename, and rerun preservation tests. Do not introduce a broad virtual filesystem or repository interface solely for this stage. Reuse the existing generic runtime seam and any existing narrow atomic-write/target-resolution seams; keep internal parsing, validation, and orchestration concrete.

### Composition and configuration

Extend the existing application composition root with the minimum code-context dependencies. Do not add package-global runtime runners, test hooks, stage registries, or a second application container. Expensive runtime construction should remain behind the existing lazy boundary so prompt, validate, status, and dry-run do not pay runtime startup costs.

Model and variant resolution uses the existing fixed precedence and records whether a command override was explicitly supplied. Omitted flags must not mask stage-specific or global fallbacks. Merge first, validate the effective result, and expose only redacted source-aware projections.

### Cancellation, concurrency, and recovery

One caller-owned context reaches target resolution, runtime execution, event draining, waiting, validation, and normal persistence work. The shared operation layer owns cancellation; browser disconnect is only subscription loss. Server shutdown invokes the same canonical cancellation function with its existing reason and single-terminal-outcome arbitration.

The code-context stage itself remains sequential. There is no demonstrated need for repository-inspection fan-out, worker pools, or a stage-local event broker. Existing runtime and web-operation infrastructure may use goroutines, but every launch must have an owner, bounded buffers, explicit waiting, and cancellation. Bounded reconciliation may use the existing cleanup policy after work-context cancellation, but timeout must produce `interrupted` or `cleanup_uncertain`, never inferred success.

### Errors, observability, and verification

Reuse existing typed or stable error identities for prerequisite validation, configuration, conflict, runtime, missing output, invalid output, persistence, cancellation, interruption, and cleanup uncertainty. Preserve wrapped causes internally; render safe actionable details at CLI/web boundaries. Add no broad error hierarchy when an existing class plus wrapped detail is sufficient.

Emit existing structured operation/runtime metadata with project, sprint, stage, operation/run identity, attempt, runtime, model, variant, source, duration, validation outcome, cancellation reason, and terminal state where available. Do not emit full prompts, source excerpts, unsafe paths, secrets, raw provider payloads, or unbounded stderr.

Verification is layered: domain/order tests, compatibility fixtures, validator tables, fake-runtime service tests, real temporary repository/filesystem isolation tests, command/config/default tests, app/web contract tests, shutdown/cancellation tests, race tests, and repository-wide build/vet/test gates. Normal tests remain offline; real-runtime proof is Sprint 34 scope.

Final decision: proceed with a focused `internal/sprint` vertical slice and small additive wiring changes. No architecture-wide refactor is required, but any duplicated stage ordering or route-specific web behavior encountered during implementation must be consolidated into the existing canonical domain/app boundary before adding code-context branches.

## Trade-Offs

| Decision | Benefit | Cost / Limitation |
| --- | --- | --- |
| Keep all stage semantics in `internal/sprint` | Preserves one owner for state, validation, runtime coordination, and recovery | The sprint package grows and must stay split into focused files |
| Reuse generic runtime and app operation seams | Avoids duplicate process, cancellation, progress, and web workflow logic | Requires careful translation between product semantics and generic request/result types |
| Explicit fixed stage registration | Compile-time discoverability and deterministic order | Canonical projections must be tested against drift |
| Lazy runtime construction | Keeps status, prompt, validation, and dry-run cheap and non-mutating | Runtime/preflight failures occur later and need clear attribution |
| Isolated candidate then atomic promotion | Preserves the last valid artifact across failed reruns | Requires temporary-file cleanup and truthful separation of artifact validity from latest run outcome |
| Read-time compatibility without migration writes | Legacy workspaces remain usable and queries remain side-effect free | Compatibility interpretation remains part of the load path |
| Concrete internal helpers with narrow external seams | Keeps the flow traceable and avoids interface proliferation | Some tests use real temporary files rather than pure in-memory fakes |
| Sequential stage orchestration | Simplifies cancellation, mutation boundaries, and terminal-state reasoning | Does not optimize repository exploration latency without measured need |

Rejected alternatives:

- **New `internal/codecontext` product module:** rejected because code-context is sprint state and workflow behavior, not an independent product domain.
- **Generic stage framework or dynamic registry:** rejected because one fixed stage does not justify lifecycle, collision, or plugin complexity.
- **Runtime-owned validation or flow transitions:** rejected because the platform runtime must remain product-neutral and runtime success is not product success.
- **Web-owned code-context service or route workflow:** rejected because it would create a parallel application layer and durable-state ambiguity.
- **Direct generation over `code-context.md`:** rejected because failure or cancellation could corrupt the last valid artifact.
- **Read-time state migration:** rejected because status and other queries must not hide writes.
- **Parallel repository-inspection workers:** rejected because no measured latency or independent work units justify the race, leak, and cancellation burden.
- **Broad filesystem abstraction:** rejected because the critical guarantees depend on real local containment and rename semantics, while existing narrow seams already support deterministic tests.

## Evidence

### Report findings

- `studies/go-cli-study/reports/final/01-project-structure.md` finds that mature CLIs keep entrypoints thin and dependencies inward and acyclic (`restic/cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`; Helm's command/action split). This supports `internal/sprint` ownership behind app/interface adapters.
- `studies/go-cli-study/reports/final/02-command-architecture.md` finds that command factories and shared execution wrappers keep lifecycle behavior out of large handlers (`gh-cli/pkg/cmdutil/factory.go:16-43`, `rclone/cmd/cmd.go:240-340`). This supports additive command wiring and shared operation lifecycle reuse.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` finds that explicit manual composition, narrow volatile-boundary interfaces, and lazy factories dominate, while large containers and globals create coupling (`gh-cli/internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `restic/internal/restic/repository.go:18-66`). This supports extending existing composition without a DI framework or globals.
- `studies/go-cli-study/reports/final/04-configuration-management.md` finds that explicit precedence, changed-flag tracking, and validation after merge are required for truthful effective configuration (`go-task/internal/flags/flags.go:314-327`, `restic/internal/global/global.go:139,147`, `k9s/internal/config/k9s.go:423-451`). This governs model/variant resolution.
- `studies/go-cli-study/reports/final/05-error-handling.md` finds that wrapped causes plus stable programmatic classification support recovery and boundary-specific rendering (`helm/pkg/storage/driver/driver.go:27-48`, `gh-cli/internal/ghcmd/cmd.go:44-49,281-301`). This supports reusing stable failure classes rather than string matching.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` finds that external-system seams improve deterministic tests but broad filesystem interfaces impose maintenance cost (`gh-cli/pkg/iostreams/iostreams.go:551-568`, `restic/internal/fs/interface.go:10-31`). This supports narrow runtime/IO seams plus real filesystem contract tests.
- `studies/go-cli-study/reports/final/07-state-context.md` finds that caller-owned context must reach actual work and that cleanup may need a separately bounded context (`helm/pkg/cmd/install.go:333-347`, `restic/internal/restic/lock.go:290-305`). This supports canonical cancellation and explicit uncertainty after bounded reconciliation.
- `studies/go-cli-study/reports/final/08-concurrency.md` finds that high-quality concurrency localizes goroutine ownership, bounds work, and waits explicitly, while sequential execution avoids needless race/leak complexity (`restic/internal/repository/repository.go:567`, `opencode/cmd/root.go:261-279`, `urfave-cli/command_run.go:92`). This supports sequential stage orchestration.
- `studies/go-cli-study/reports/final/10-logging-observability.md` finds that structured fields and result/diagnostic separation preserve operator value and machine-readable output (`helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231`). This supports safe bounded metadata rather than raw runtime output.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` finds that command, fake, fixture, integration, and selective golden tests prove different contracts (`chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `restic/cmd/restic/integration_helpers_test.go:188-235`). This supports the layered verification plan.
- `studies/go-cli-study/reports/final/12-extensibility.md` finds that factories and additive options are lower-cost internal extension seams than registries/plugins for fixed workflows (`go-task/executor.go:20-24,91-122`, `rclone/fs/rc/registry.go:41-48`). This supports explicit canonical registration.
- `studies/go-cli-study/reports/final/13-security.md` finds that explicit trust boundaries, centralized validation, path controls, argument-safe execution, and redaction are required for untrusted inputs (`k9s/internal/config/json/validator.go:146`, `lazygit/cmd_obj_builder.go:38`, `restic/internal/options/secret_string.go:15-20`). This supports the read-only repository and single-artifact-write boundary.

### Sprint-specific inference

The handbook identifies the cross-cutting pressures but does not choose ownership. Applying the reports to the project architecture and sprint constraints yields the key conclusion: the existing module-driven design still fits, provided code-context is implemented as sprint-owned behavior with generic infrastructure beneath it and thin app/CLI/web projections above it. The requirements' explicit prohibition on new stage frameworks, repository indexes, web product services, and package-global runners confirms that the minimal direct vertical slice is the correct option.

## Risks

- **Stage-list drift:** duplicated order or artifact maps could place code-context differently across flow, status, and commands. Mitigation: canonical domain ordering plus exact-position tests across projections.
- **Product leakage into runtime:** adding sprint fields or validation to `platform/runtime` would reverse dependency direction. Mitigation: generic requests only and sprint-side translation tests.
- **Composition-root growth:** adding target, runtime, store, and operation collaborators can enlarge the app container. Mitigation: inject only existing cohesive services or lazy functions and keep stage helpers inside `internal/sprint`.
- **Artifact/state split-brain:** artifact promotion may succeed while state persistence fails, or cancellation may race with completion. Mitigation: preserve explicit transition ordering, atomic writes, single terminal arbitration, and recoverable diagnostics; never infer completion from file presence.
- **Legacy-state ambiguity:** deriving inserted-stage state incorrectly could block or skip work. Mitigation: fixtures for representative pre-stage states and proof that reads do not mutate them.
- **Repository mutation escape:** runtime permissions alone may not prove read-only behavior. Mitigation: permission policy, fixed expected output, prompt constraints, temporary-repository before/after comparisons, and Git-state checks in tests.
- **Over-abstraction:** similar planning stages may tempt a generic framework during implementation. Mitigation: reuse only proven mechanical helpers and keep code-context policy concrete.
- **Cleanup uncertainty:** cancellation during runtime or candidate promotion can leave temporary data or uncertain ownership. Mitigation: bounded cleanup, no authoritative promotion before validation, and explicit interrupted/cleanup-uncertain outcomes.
- **Open verification question:** implementation must confirm the existing atomic helper provides same-directory replace semantics and preserves the prior file on rename/write failure. If not, add the narrow mechanical guarantee without introducing a new persistence abstraction.
