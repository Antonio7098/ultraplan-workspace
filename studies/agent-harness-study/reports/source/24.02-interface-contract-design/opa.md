# Source Analysis: opa

## 24.02 Interface Contract Design

### Source Info

| Field | Value |
|-------|-------|
| Name | opa (Open Policy Agent) |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go (module `github.com/open-policy-agent/opa`), embedded WASM, JSON schemas, Rego policy language |
| Analyzed | 2026-08-22 |

## Summary

OPA defines its cross-boundary contracts almost exclusively as small Go interfaces with explicit behavioral documentation, plus a set of generated, schema-checked artifacts (`capabilities.json`, `builtin_metadata.json`, `v1/ir/plan.schema.json`) that pin down the machine-facing surface. The dominant seams are: the storage backend contract (`sources/opa/v1/storage/interface.go:20-44`), the built-in function extension contract (`sources/opa/v1/topdown/builtins.go:63-68`), the plugin lifecycle contract (`sources/opa/v1/plugins/plugins.go:106-110`), the query tracer contract (`sources/opa/v1/topdown/trace.go:183-187`), and the host-application embedding API in `rego`/`sdk` (functional-options based, e.g. `sources/opa/v1/rego/rego.go:199`). Contracts are validated at three levels: compile time via Go types; "schema time" via generated JSON schemas and CI-enforced regeneration of capability artifacts (`sources/opa/main.go:33-36`, `sources/opa/.github/workflows/pull-request.yaml:100-121`); and runtime via arity/type checks against declared function signatures (`sources/opa/v1/topdown/eval.go:1017-1018`, `sources/opa/v1/topdown/eval.go:2079`). Evolution is handled with an unusually disciplined deprecation pattern: old interfaces are kept and adapted to new ones by wrappers (`sources/opa/v1/topdown/trace.go:195-217`), unsupported capabilities are surfaced through dedicated sentinel errors rather than silent no-ops (`sources/opa/v1/storage/errors.go:31-42`), and the whole v0 surface is aliased onto `v1` packages via type aliases (`sources/opa/storage/storage.go:10-30`). The main weaknesses are stringly-typed storage error codes compared by type assertion without unwrap support (`sources/opa/v1/storage/errors.go:58-96`), a wide grab-bag `BuiltinContext` struct (`sources/opa/v1/topdown/builtins.go:37-61`), and the absence of a shared conformance suite that exercises all `Store` implementations against one scenario matrix.

## Rating

**8 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale: contracts are small, named, and documented at the definition site; behavior (not just signatures) is encoded in doc comments and enforced at runtime (arity checks, error classification, sentinel "not supported" errors); compatibility across versions is actively engineered (deprecated interfaces wrapped, generated artifacts CI-checked for drift). It falls short of 9–10 because several predicates rely on exact type assertions rather than `errors.As` unwrapping (`sources/opa/v1/storage/errors.go:58-96`), `Manager.AuthPlugin` performs an unchecked interface type assertion that can panic on misregistration (`sources/opa/v1/plugins/plugins.go:740-749`), and there is no shared store-conformance test suite — each `Store` implementation is tested only by its own package-local tests.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Storage transaction contract | `Transaction` is a minimal ID-bearing handle; `Store` composes `Trigger` + `Policy` + txn/read/write/commit/truncate/abort; comment states Commit error implies automatic abort | `sources/opa/v1/storage/interface.go:15-44` |
| Optional storage capabilities | `MakeDirer`, `NonEmptyer`, `Closer` are separate optional interfaces a Store may realize | `sources/opa/v1/storage/interface.go:46-63` |
| Unsupported-capability shims | `WritesNotSupported`, `PolicyNotSupported`, `TriggersNotSupported` embeddable structs return typed sentinel errors instead of panicking | `sources/opa/v1/storage/interface.go:142-148`, `160-180`, `238-245` |
| Trigger semantics | `TriggerConfig.OnCommit` invoked after successful commit before other clients see changes; `TriggerEvent` carries typed `PolicyEvent`/`DataEvent` change records | `sources/opa/v1/storage/interface.go:220-230`, `182-201` |
| Storage error taxonomy | Eight documented string error codes; `Error{Code, Message}` type; predicate helpers `IsNotFound`, `IsWriteConflictError`, `IsInvalidPatch`, `IsInvalidTransaction`; deprecated `IsIndexingNotSupported` stub returns false | `sources/opa/v1/storage/errors.go:11-42`, `45-55`, `58-96` |
| Built-in function contract | `BuiltinFunc` continuation-passing signature `(bctx BuiltinContext, operands []*ast.Term, iter func(*ast.Term) error) error`; registration mutates global map | `sources/opa/v1/topdown/builtins.go:63-68`, `90-93`, `127` |
| Builtin context propagation | `BuiltinContext` carries request context, cancel signal, seed, wall clock, metrics, caches, print hook, round-tripper customizer, caller metadata | `sources/opa/v1/topdown/builtins.go:35-61` |
| Builtin error classification | `handleBuiltinErr` distinguishes `BuiltinEmpty`, `*Error`/`Halt` (pass-through), `builtins.ErrOperand` (→ `TypeErr`), default → `BuiltinErr`, wrapping with location | `sources/opa/v1/topdown/builtins.go:182-203` |
| Deprecated builtin style retained | `FunctionalBuiltin1..4` marked deprecated with adapters still registered | `sources/opa/v1/topdown/builtins.go:22-33`, `95-113` |
| Tracer contract evolution | `Tracer` deprecated; replacement `QueryTracer{Enabled, TraceEvent(Event), Config() TraceConfig}`; `legacyTracer` adapts old→new preserving plug-local-vars behavior | `sources/opa/v1/topdown/trace.go:172-217` |
| Plugin lifecycle contract | `Plugin{Start(ctx) error, Stop(ctx), Reconfigure(ctx, config any)}`; optional `Triggerable`; `LoggerPlugin` extends `Plugin` with `slog.Handler` | `sources/opa/v1/plugins/plugins.go:106-123` |
| Plugin state reporting | Typed `State` enum (`NOT_READY`, `OK`, …) plugins must report through the manager's status channel | `sources/opa/v1/plugins/plugins.go:125-140` |
| Manager lifecycle semantics | `Register` enqueues status init event; `Start` fail-fast on first plugin start error; `Stop` documents graceful-shutdown deadline and non-reentrancy | `sources/opa/v1/plugins/plugins.go:700-714`, `870-918`, `920-926` |
| Auth plugin behavioral contract | Doc comment mandates single-call `NewClient` vs per-request `Prepare` split ("MUST" language) | `sources/opa/v1/plugins/rest/rest.go:42-51` |
| Embedding API options pattern | `EvalOption func(*EvalContext)` functional options (`EvalInput`, `EvalTransaction`, `EvalQueryTracer`, …) over an unexported `EvalContext` with exported accessors | `sources/opa/v1/rego/rego.go:95-199`, `202-274` |
| Prepared-query contract | `PreparedEvalQuery.Eval` documents txn reuse rules (new txn opened unless `EvalTransaction` given) and option-override precedence | `sources/opa/v1/rego/rego.go:550-569` |
| Aggregate error contract | `Errors []error` with singular/plural formatting; `IsPartialEvaluationNotEffectiveErr` sentinel check | `sources/opa/v1/rego/rego.go:592-619` |
| External resolver contract | One-method `Resolver` interface with explicit `Input{Ref, Input, Metrics}` / `Result{Value}` value types | `sources/opa/v1/resolver/interface.go:14-29` |
| Print hook contract | Single-method `Hook.Print(Context, string) error` injected via `BuiltinContext` | `sources/opa/v1/topdown/print/print.go:19-21` |
| Build target enumeration | `TargetWasm`/`TargetPlan` constants plus `Targets` slice; `init()` rejects unknown targets via `slices.Contains(Targets, c.target)` | `sources/opa/v1/compile/compile.go:44-54`, `420-425` |
| IR plan schema generation | `genplanschema` reflects `v1/ir/ir.go` into `plan.schema.json` ($id `https://openpolicyagent.org/schemas/ir/v1/plan.schema.json`), wired by `go:generate` | `sources/opa/internal/cmd/genplanschema/main.go:1-3`, `33-54`; `sources/opa/main.go:36`; artifact `sources/opa/v1/ir/plan.schema.json` |
| Capabilities as versioned surface | `capabilities.json` and `builtin_metadata.json` generated from code via `go:generate`; CI regenerates and uploads `capabilities.json` so drift breaks the build | `sources/opa/main.go:33-34`; `sources/opa/.github/workflows/pull-request.yaml:100-121` |
| Builtin declaration schema | `ast.Builtin{Name, Description, Decl *types.Function, Infix, Relation, Deprecated, ...}` is the serialized contract unit inside capabilities | `sources/opa/v1/ast/builtins.go:3594-3606` |
| Declaration-level contract tests | `TestBuiltinDeclRoundtrip`, `TestAllBuiltinsHaveDescribedArguments` verify every builtin declares its argument structure | `sources/opa/v1/ast/builtins_test.go:16`, `33` |
| Capability ordering/versioning tests | `TestCapabilitiesAddBuiltinSorted`, `TestCapabilitiesMinimumCompatibleVersion` keep the published capability list deterministic and version-compatible | `sources/opa/v1/ast/capabilities_test.go:185`, `219` |
| Runtime arity validation | Evaluator reads declared arity from `TypeEnv` (`e.compiler.TypeEnv.GetByRef(ref).(*types.Function).Arity()`) and from `bi.Decl.Arity()` when dispatching builtins | `sources/opa/v1/topdown/eval.go:1017-1018`, `1054-1058`, `2079-2090` |
| Runtime error routing | Errors from builtin iterator callbacks wrapped in `Halt` to avoid polluting `builtinErrors`; other builtin errors appended to query error set | `sources/opa/v1/topdown/eval.go:2151-2169` |
| v0 → v1 alias shim | Root `storage` package re-exports v1 identifiers as type aliases (`type Store = v1.Store`, etc.) so both import paths satisfy the same interface | `sources/opa/storage/storage.go:10-30` |
| Host SDK entry contract | `sdk.OPA` instance constructed via `New(ctx, Options)`; `Options` struct in same package forms the embedding config surface | `sources/opa/sdk/opaq.go`; `sources/opa/v1/sdk/opa.go:43-66`; `sources/opa/v1/sdk/options.go:36` |

## Answers to Dimension Questions

**1. Are interfaces small, coherent, and owned by the consumer side?**
Mostly small and coherent: `Resolver` is one method (`sources/opa/v1/resolver/interface.go:15-17`), `print.Hook` one method (`sources/opa/v1/topdown/print/print.go:19-21`), `Plugin` three lifecycle methods (`sources/opa/v1/plugins/plugins.go:106-110`). They are not strictly consumer-owned: `storage.Store` (11 methods, `sources/opa/v1/storage/interface.go:20-44`) and `topdown.BuiltinFunc` are defined by the producing package, which both sides import — a shared-kernel model rather than Go-idiomatic consumer-side definition. OPA compensates by splitting optional behavior into separate narrow interfaces (`MakeDirer`, `NonEmptyer`, `Closer`, `Triggerable`, `LoggerPlugin`) instead of growing the core interface. The `Rego` embedding API avoids interface bloat entirely with functional options over an unexported context struct (`sources/opa/v1/rego/rego.go:95-132`, `199`). Counter-example of width: `BuiltinContext` has grown to ~25 fields including two deprecated ones (`sources/opa/v1/topdown/builtins.go:37-61`).

**2. Do contracts specify behavior, not just method signatures?**
Yes, in three ways. First, normative doc comments: `Commit` must auto-abort on error (`sources/opa/v1/storage/interface.go:33-35`), auth plugin `NewClient` is once-per-client while `Prepare` is per-request (`sources/opa/v1/plugins/rest/rest.go:44-50`), `OnCommit` fires after commit but before other clients observe changes (`sources/opa/v1/storage/interface.go:226-229`). Second, structural scaffolding encodes behavior: embeddable `WritesNotSupported`/`PolicyNotSupported`/`TriggersNotSupported` structs make partial-capability stores return correct sentinel errors (`sources/opa/v1/storage/interface.go:142-148`, `160-180`, `238-245`). Third, the evaluator enforces semantics at runtime: builtin results are routed differently depending on whether the declaration has a result (`e.bi.Decl.Result() == nil` branches, `sources/opa/v1/topdown/eval.go:2134-2141`), and errors are classified into `Halt` vs collected builtin errors vs type errors (`sources/opa/v1/topdown/builtins.go:182-203`, `sources/opa/v1/topdown/eval.go:2162-2169`).

**3. Can providers, tools, stores, and runtimes be replaced safely?**
Largely yes. A second `storage.Store` implementation exists in-tree (`v1/storage/inmem`, `v1/storage/disk`) behind the same interface, and the optional-capability interfaces plus sentinel errors let an implementation declare what it lacks. Custom builtins, tracers, resolvers, auth plugins, and loggers all plug in via narrow interfaces; the legacy `Tracer` keeps pre-existing implementations working through `WrapLegacyTracer` (`sources/opa/v1/topdown/trace.go:213-217`). Caveats: (a) `Manager.AuthPlugin` does an unchecked assertion `plugin.(rest.HTTPAuthPlugin)` that will panic if a plugin registered under an auth name doesn't implement the interface (`sources/opa/v1/plugins/plugins.go:740-749`) — substitutability there rests on naming convention, not validation; (b) `RegisterBuiltinFunc` writes a process-global map with no collision detection or unregister hook (`sources/opa/v1/topdown/builtins.go:90-93`, `127`), so two independent providers claiming the same builtin name cannot coexist safely; (c) no shared conformance suite proves third-party `Store` implementations behave identically to in-tree ones (see Gaps).

**4. Are compatibility failures caught early by tests or validation?**
Yes for the declarative/machine-readable surface, partially for Go interfaces. Generated artifacts (`capabilities.json`, `builtin_metadata.json`, `v1/ir/plan.schema.json`) are produced by `go:generate` from code (`sources/opa/main.go:33-36`), and CI regenerates them so drift between registered builtins and the published capability files fails the build (`sources/opa/.github/workflows/pull-request.yaml:100-121`). Declaration-shape invariants are unit-tested (`TestAllBuiltinsHaveDescribedArguments`, `sources/opa/v1/ast/builtins_test.go:33`; roundtrip test at `:16`; capability sort order and minimum-version logic at `sources/opa/v1/ast/capabilities_test.go:185`, `219`). Unknown build targets are rejected early in `Compiler.init` (`sources/opa/v1/compile/compile.go:425`). For Go interface evolution, the safety net is compiler-checked aliases plus wrapper adapters rather than dedicated conformance tests; runtime misuse (wrong arity) surfaces only during evaluation (`sources/opa/v1/topdown/eval.go:2079-2090`).

## Architectural Decisions

- **Shared-kernel interfaces instead of consumer-side ports.** Core seams (`storage.Store`, `topdown.BuiltinFunc`, `plugins.Plugin`) live in the engine's own packages; every extension point imports them. This trades strict dependency inversion for discoverability and a stable, centrally documented vocabulary.
- **Continuation-passing style for builtins.** `BuiltinFunc(bctx, operands, iter)` lets a builtin yield multiple results (relations) and integrates with evaluation control flow; the older functional style remains only through deprecated wrappers (`sources/opa/v1/topdown/builtins.go:63-68`, `139-180`).
- **Capability negotiation via optional interfaces + sentinel errors.** Stores advertise partial support by embedding `*NotSupported` helper structs whose methods return coded errors (`sources/opa/v1/storage/interface.go:142-180`, `sources/opa/v1/storage/errors.go:31-42`) — callers probe with type checks/predicates instead of relying on hidden defaults.
- **Code-generated, CI-pinned machine contracts.** The externally observable builtin surface and the planner IR are derived from Go source by generators (`sources/opa/internal/cmd/genplanschema/main.go:33-54`; `go:generate` wiring at `sources/opa/main.go:33-36`) and checked in, making drift visible in review and CI.
- **Compatibility by adapter layer.** Deprecated interfaces stay compilable and semantically preserved via wrappers (`legacyTracer` forcing `PlugLocalVars: true` for old tracers, `sources/opa/v1/topdown/trace.go:203-211`; root-package type aliases onto v1, `sources/opa/storage/storage.go:10-30`).
- **Functional options for the embedding API.** `PreparedEvalQuery.Eval(ctx, opts...)` separates expensive preparation from cheap per-evaluation configuration, with documented transaction ownership rules (`sources/opa/v1/rego/rego.go:554-569`).

## Notable Patterns

- **Embeddable default-behavior structs**: `WritesNotSupported`, `PolicyNotSupported`, `TriggersNotSupported` give partial implementations correct-by-construction error behavior (`sources/opa/v1/storage/interface.go:144-148`).
- **Adapter/wrapper for interface migration**: `legacyTracer` + `WrapLegacyTracer` translate the deprecated `Tracer` to `QueryTracer` without changing call sites (`sources/opa/v1/topdown/trace.go:194-217`).
- **Sentinel-error predicates**: `IsNotFound`, `IsWriteConflictError`, `IsInvalidPatch`, `IsInvalidTransaction` form the error-matching API (`sources/opa/v1/storage/errors.go:57-96`).
- **Typed change events**: `TriggerEvent`/`PolicyEvent`/`DataEvent` give triggers structured, inspectable deltas instead of raw diffs (`sources/opa/v1/storage/interface.go:182-201`).
- **Normative MUST-language in doc comments** for cross-boundary timing guarantees (`sources/opa/v1/plugins/rest/rest.go:44-51`).
- **Versioned capability sets**: `Capabilities.MinimumCompatibleVersion` maintained by tests (`sources/opa/v1/ast/capabilities_test.go:219`) allows consumers to gate features by OPA version.
- **Single-method extension points** (`resolver.Resolver`, `print.Hook`) keep third-party obligations minimal (`sources/opa/v1/resolver/interface.go:15-17`; `sources/opa/v1/topdown/print/print.go:19-21`).

## Tradeoffs

- **Global registration state vs pluggability**: `builtinFunctions[name] = ...` makes custom builtins trivially easy but process-global — no namespacing, no removal, duplicate-name last-writer-wins (`sources/opa/v1/topdown/builtins.go:91-93`, `127`). Two independent implementations of the same builtin name cannot coexist.
- **Stringly-typed storage error codes vs rich errors**: codes are stable, serializable strings (`sources/opa/v1/storage/errors.go:11-42`), but predicates use direct `err.(*Error)` assertions, so a wrapped or re-boxed storage error fails `IsNotFound` silently (`sources/opa/v1/storage/errors.go:58-63`) — fragile under middleware-style error decoration elsewhere in the codebase.
- **Wide `BuiltinContext` vs convenience**: exposing ~25 fields (including deprecated `Tracers`) gives builtin authors everything but couples them to evaluator internals; additions are backward-compatible yet grow the frozen-in-capabilities surface (`sources/opa/v1/topdown/builtins.go:37-61`).
- **Doc-comment behavior specs vs enforcement**: lifecycle rules like "cannot call Stop twice or it will hang" (`sources/opa/v1/plugins/plugins.go:925`) or "Commit must auto-abort" (`sources/opa/v1/storage/interface.go:33-35`) rely on implementer discipline; only some are backed by tests.
- **Generated-artifact pinning vs maintenance cost**: checked-in `capabilities.json` (~94 KB) and `builtin_metadata.json` (~548 KB) guarantee reviewability but add generate-and-commit ceremony to every builtin change (`sources/opa/main.go:33-34`).

## Failure Modes / Edge Cases

- **Unchecked type assertion panic**: registering any plugin under an auth-plugin name without implementing `rest.HTTPAuthPlugin` panics later at lookup (`plugin.(rest.HTTPAuthPlugin)`, `sources/opa/v1/plugins/plugins.go:745`).
- **Silent mismatch on wrapped errors**: storage helpers returning wrapped errors break `IsNotFound`-style checks because they assert the concrete `*storage.Error` type only (`sources/opa/v1/storage/errors.go:59-62`).
- **Deprecated-but-live paths can diverge**: `FunctionalBuiltin4` treats `BuiltinEmpty` specially while `BuiltinFunc` uses nil-return convention; mixing conventions changes whether "undefined" vs error results (`sources/opa/v1/topdown/builtins.go:169-179`, `120-125`).
- **Manager.Stop reentrancy hang**: calling `Stop` twice hangs by documented limitation (`sources/opa/v1/plugins/plugins.go:925`), an edge case delegated to callers.
- **Arity edge handling in evaluator**: dropping the trailing term when operand count exceeds declared arity (`sources/opa/v1/topdown/eval.go:2083-2090`) is subtle; a builtin whose declared arity disagrees with its implementation would misroute operands rather than fail loudly.
- **Partial-evaluation ineffectiveness**: signaled through a package-private sentinel reachable only via `IsPartialEvaluationNotEffectiveErr` on an `Errors` aggregate (`sources/opa/v1/rego/rego.go:609-619`) — callers must know to use the predicate rather than compare errors.

## Future Considerations

- Introduce `errors.Is/As` support (or `Unwrap`) for `storage.Error` so predicates survive error wrapping (`sources/opa/v1/storage/errors.go:44-96`).
- Validate at `Manager.Register` time (not lookup time) that plugins registered as auth plugins implement `rest.HTTPAuthPlugin` (`sources/opa/v1/plugins/plugins.go:702-714`, `740-749`).
- Provide namespaced or instance-scoped builtin registration with duplicate detection to replace the global map (`sources/opa/v1/topdown/builtins.go:127`).
- Extract a shared storage conformance suite (scenario matrix run against `inmem`, `disk`, and third-party stores) to make substitutability provable; current tests are package-local (`sources/opa/v1/storage/storage_test.go`, `sources/opa/v1/storage/inmem/*_test.go`).
- Continue shrinking `BuiltinContext` toward grouped sub-contexts as deprecated fields (`Tracers`) retire (`sources/opa/v1/topdown/builtins.go:49-50`).

## Questions / Gaps

- **No cross-store conformance evidence found.** Searches for shared test harnesses (`storagetest`, `conformance`) across `v1/storage/**` returned no common suite; whether out-of-tree `Store` implementations (e.g., the documented disk store or community stores) satisfy identical semantics is unverified within this source snapshot.
- **Contract coverage of the REST server surface was not deeply audited.** `/v1/*` endpoint schemas exist under `sources/opa/v1/server/`, but this analysis focused on Go-level and generated-artifact contracts; a full API-contract study would need to inspect `server/types` and OpenAPI tooling separately.
- **Wasm SDK boundary**: how `wasm` runtime inputs/outputs conform to the plan IR schema (`sources/opa/v1/ir/plan.schema.json`) at runtime (validation vs trust) was not traced end-to-end; the schema is generated from `v1/ir/ir.go` (`sources/opa/internal/cmd/genplanschema/main.go:54`) but runtime enforcement points were not located.
- Whether `RequestMetadata`/`ResponseMetadata` propagation through `BuiltinContext` (`sources/opa/v1/topdown/builtins.go:57-58`) has a written stability guarantee could not be confirmed from code alone — comments say "for use by wrapping projects" without a versioning promise.

---

Generated by `24.02-interface-contract-design` against `opa`.
