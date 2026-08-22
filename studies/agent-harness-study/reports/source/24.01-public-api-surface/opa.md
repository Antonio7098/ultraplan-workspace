# Source Analysis: opa

## Public API Surface (Dimension 24.01)

### Source Info

| Field | Value |
|-------|-------|
| Name | opa (Open Policy Agent) |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go 1.25 (`go.mod:3`), Cobra CLI, net/http server, Wasm/plan compile targets |
| Analyzed | 2026-08-22 |

## Summary

OPA exposes a deliberately layered public API with four consumer-facing surfaces plus two extension surfaces:

1. **REST API** — the HTTP policy engine (`v1/server/server.go:895-918` registers the `/v0/data`, `/v1/data`, `/v1/policies`, `/v1/query`, `/v1/compile`, `/v1/config`, `/v1/status` route table), documented endpoint-by-endpoint in `docs/docs/rest-api.md` (~2,464 lines).
2. **Low-level Go evaluation API** — `v1/rego` ("Package rego exposes high level APIs for evaluating Rego policies", `v1/rego/rego.go:5`), built around `rego.New(...)` (`v1/rego/rego.go:1414`), `PrepareForEval` (`v1/rego/rego.go:1786`) and `PreparedEvalQuery.Eval` (`v1/rego/rego.go:559`). ~176 exported symbols.
3. **High-level embedding SDK** — `v1/sdk` ("a high-level API for embedding OPA inside of Go programs", `v1/sdk/opa.go:5`) with lifecycle-managed `sdk.OPA` (`v1/sdk/opa.go:43`), `New` (`v1/sdk/opa.go:66`), `Decision` (`v1/sdk/opa.go:300`), `Partial` (`v1/sdk/opa.go:438`), and `Stop` (`v1/sdk/opa.go:289`).
4. **CLI** — 16 subcommands registered in one place (`cmd/commands.go:35-50`: bench, build, capabilities, check, deps, eval, exec, fmt, inspect, oracle, parse, refactor, run, sign, test, version), executed by a minimal `main.go:14-30`.
5. **Extension APIs** — plugin `Factory`/`Plugin` interfaces (`v1/plugins/plugins.go:89-110`), `runtime.RegisterPlugin` (`v1/runtime/runtime.go:94-98`), custom storage backends (`v1/runtime/runtime.go:103-107`), and custom builtins (`rego.RegisterBuiltin1..Dyn`, `v1/rego/rego.go` equivalents at root shim `rego/rego.go:220-242`).

The distinguishing architectural move is the **v0/v1 dual-module split**: the single Go module `github.com/open-policy-agent/opa` (`go.mod:1`) keeps deprecated v0 packages at the repo root while all current code lives under `v1/`. Every v0 package carries an identical deprecation banner ("Deprecated: This package is intended for older projects transitioning from OPA v0.x…", e.g. `rego/doc.go:5-8`; 73 packages carry this marker), and each v0 symbol is a Go type alias or thin wrapper onto the `v1` implementation (e.g. `type OPA = v1.OPA` in `sdk/opa.go:13-18`, `type Rego = v1.Rego` in `rego/rego.go:193`), so both import paths share exactly one implementation. The `v1` module doc states the contract plainly: "The v1 API defaults to enforcing the v1 Rego syntax… Most packages outside the v1 API are deprecated" (`v1/doc.go:5-8`).

Stability is communicated through compiler-enforced boundaries (`internal/` packages at `internal/` and nested `v1/ast/internal/`, `v1/storage/internal/` — invisible to external importers per Go's internal-package rule), godoc `Deprecated:` markers inside v1 itself (e.g. `v1/ast/compile.go:592`, `v1/loader/loader.go:499-517`), explicit `EXPERIMENTAL:` labels (`v1/debug/debugger.go:6`, `v1/loader/extension/extension.go:15`, `v1/ast/interning.go:18`), and machine-readable builtin metadata generated into the repo (`main.go:32-36` generates `capabilities.json`, `builtin_metadata.json`, and `v1/ast/version_index.json`, the last embedded via `//go:embed` at `v1/ast/capabilities.go:39` to record which OPA version introduced each builtin).

## Rating

**8 / 10** — A clear, well-documented model with explicit interfaces and real operational safeguards. The v0/v1 alias split, uniform deprecation banners, lint-enforced deprecation hygiene (`.golangci.yaml:94-102`), extensive runnable examples, and machine-readable capability metadata make it easy for a new integration to use the stable API without reading internals. It falls short of 9-10 because compatibility is enforced socially rather than mechanically — there is no automated Go-API diff/gorelease check in CI (searched for `apidiff|gorelease|api-diff`; only `make check` → golangci-lint exists, `Makefile:161-165`) — and a few spots leak global mutable state (`cmd.RootCommand`, `cmd.UserAgent` in `cmd/commands.go:14-23`; `sdk.SetDefaultOptions` in `v1/sdk/options.go:23-33`).

## Evidence Collected

Every entry cites file paths with line numbers relative to the source root.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Module identity & retract | Module path pinned lower-case; pre-modules releases retracted so pkg.go.dev warns on wrong import path | `go.mod:1`, `go.mod:127-132` |
| v1 API contract | Package doc: v1 defaults to Rego v1 syntax; packages outside v1 are deprecated v0 API | `v1/doc.go:5-9` |
| v0 deprecation banners | Identical "Deprecated:" package comment across 73 v0 shim packages | `rego/doc.go:5-8`, `sdk/doc.go:5-8`, `topdown/doc.go:11` |
| v0 shims are aliases | `type OPA = v1.OPA`, `type Options = v1.Options`, etc.; `type Rego = v1.Rego` | `sdk/opa.go:13-18`, `rego/rego.go:28-217` |
| Layered integration story | Docs enumerate REST API, Go API, Wasm, IR, SDK as evaluation interfaces | `docs/docs/integration.md:23-36` |
| REST route table | All versioned endpoints registered on one mux; `/health` unversioned; method-not-allowed catch-alls | `v1/server/server.go:895-928` |
| REST reference docs | Per-endpoint documentation incl. Policy/Data/Query/Compile/Health/Config/Status APIs | `docs/docs/rest-api.md:48-2103` |
| Rego entry points | `New`, `PrepareForEval`, `(pq PreparedEvalQuery).Eval` | `v1/rego/rego.go:1414`, `v1/rego/rego.go:1786`, `v1/rego/rego.go:559` |
| Prepared-query encapsulation | `PreparedEvalQuery` holds pre-compiled state; Eval takes only `EvalOption`s | `v1/rego/rego.go:548-569` |
| SDK lifecycle | `sdk.New` blocks until ready unless `Options.Ready` channel supplied; atomic `Configure`; non-restartable `Stop` | `v1/sdk/options.go:54-57`, `v1/sdk/options.go:110-115`, `v1/sdk/opa.go:141-162`, `v1/sdk/opa.go:288-297` |
| SDK decision API | `Decision(ctx, DecisionOptions)` returns `DecisionResult{ID, Result, Provenance}`; options struct fully commented | `v1/sdk/opa.go:299-377` |
| CLI registration | Single `Command()` factory wires 16 subcommands; `RootCommand` kept for back-compat | `cmd/commands.go:13-52` |
| Binary entry point | `main.go` only executes `cmd.RootCommand` and maps `cmd.ExitError` to exit codes | `main.go:14-30` |
| Plugin extension contract | `Factory{Validate,New}` and `Plugin{Start,Stop,Reconfigure}` interfaces with step-by-step godoc | `v1/plugins/plugins.go:42-110` |
| Runtime extension hooks | `RegisterPlugin(name, factory)`; `RegisterStorageBackend(builder)` | `v1/runtime/runtime.go:90-107` |
| Custom builtins | `RegisterBuiltin1..4/Dyn` global registration + `Function1..Dyn` per-query options | `rego/rego.go:219-267` |
| Internal boundary | Root `internal/` tree (38 pkgs) + nested internals unreachable externally; 6 `v1/**/internal` packages | `internal/`, `v1/ast/internal/scanner`, `go list ./v1/...` |
| E2E isolation | Separate Go module for heavy test deps; README forbids importing it from OPA packages | `e2e/go.mod:1-6`, `e2e/README.md:1-6` |
| Deprecation hygiene inside v1 | `Deprecated:` markers with replacement guidance; lint config silences SA1019 in tests and v0-shim paths | `v1/ast/compile.go:592-600`, `v1/loader/loader.go:499-517`, `.golangci.yaml:94-102` |
| Experimental labels | `EXPERIMENTAL:` on debugger, loader extensions, AST interning; `UserAgent` marked experimental | `v1/debug/debugger.go:6`, `v1/loader/extension/extension.go:15-34`, `v1/ast/interning.go:18`, `cmd/commands.go:19-23` |
| Machine-readable API metadata | go:generate emits `capabilities.json`, `builtin_metadata.json`, builtin `version_index.json` (embedded) | `main.go:33-36`, `v1/ast/capabilities.go:39`, `v1/ast/version_index.json:1-60` |
| Godoc examples | 22 runnable `Example*` functions covering eval, partial eval, tracing, prepared queries, custom builtins | `v1/rego/example_test.go:30-985`, `v1/ast/example_test.go:13-55` |
| Runnable SDK example in docs | Complete program: mock bundle server, config, `sdk.New`, `Decision`, `Stop` | `docs/docs/integration.md:195-263` |
| Public test-support package | `sdk/test` provides `MockBundle`/test server used by the official example | `docs/docs/integration.md:204`, `v1/sdk/test/test.go:27` |
| v0/v1 migration guidance | Compatibility modes per surface; mixing v0+v1 packages declared unsupported anti-pattern | `docs/docs/v0-compatibility.md:32-58`, `docs/docs/v0-compatibility.md:147-151` |
| Extension docs | Custom builtins, plugins (with status), storage backends, runtime version override — all with Go examples | `docs/docs/extensions.md:10-575` |

## Answers to Dimension Questions

**1. What is the intended public API surface?**
Four consumer surfaces with an explicit division of labor, stated in the integration guide (`docs/docs/integration.md:23-36`): REST API for non-Go consumers, `v1/rego` for "only policy evaluation" embedding, `v1/sdk` when management features (bundle discovery, decision logs) are wanted, and the CLI for operator workflows (`cmd/commands.go:25-52`). Two extension surfaces round it out: plugins/builtins/storage backends (`v1/plugins/plugins.go:89-110`, `v1/runtime/runtime.go:94-107`) and the Wasm/plan compile target consumed via the documented IR format (`docs/docs/integration.md:33-34`, schema generated at `v1/ir/plan.schema.json` per `main.go:36`). Everything else is explicitly internal: `internal/` packages are compiler-invisible, and the e2e harness is quarantined in its own module with a warning not to import it (`e2e/README.md:3-6`).

**2. Is the stable API easy to distinguish from internal implementation details?**
Yes, unusually so. Three mechanisms work together: (a) the import-path rule — anything under `v1/` is current, root packages are deprecated shims carrying identical banners (`rego/doc.go:5-8`); (b) Go's `internal/` visibility rule hard-blocks the ~44 internal packages from external use; (c) godoc markers distinguish further tiers: `Deprecated:` with named replacements inside v1 (`v1/ast/policy.go:710` "Poor handling of ref rules. Use `(*Rule).Ref()` instead"), and `EXPERIMENTAL:` for actively-changing areas (`v1/debug/debugger.go:6`). The lint setup institutionalizes this: staticcheck SA1019 is suppressed in `_test.go` files and the v0-compat deprecation text is suppressed outside `v1/`, so the warnings stay actionable where they matter (`.golangci.yaml:94-102`). One wrinkle: `v1/topdown` is labeled "low-level query evaluation support" (`v1/topdown/doc.go:5-9`) yet is a mandatory dependency of the public `rego` API (tracers, caches, `print.Hook` all come from topdown — see imports at `v1/rego/rego.go:33-36`), so the effective stable boundary includes much of topdown whether intended or not.

**3. Does the API expose the right level of abstraction for agent harness users?**
Largely yes. The prepared-query model cleanly separates expensive compilation from hot-path evaluation (`v1/rego/rego.go:84-91`, `548-569`) — the right shape for per-request authorization. The SDK wraps plugin/discovery machinery behind three methods (`New`/`Decision`/`Stop`, `v1/sdk/opa.go:66,300,289`) while leaving deliberate escape hatches: `opa.Plugin("bundle")` for manual triggers (`docs/docs/integration.md:326-340`), `ManagerOpts []func(*plugins.Manager)` (`v1/sdk/options.go:90-93`), and `Hooks` (options.go:72-73, self-admittedly vague: "TODO(sr): find better words"). Abstraction leaks exist at the edges: `DecisionOptions.NDBCache` is typed `any` with a runtime type-assertion back to `builtins.NDBCache` inside the SDK (`v1/sdk/opa.go:310-316`), forcing users of that feature to import `v1/topdown/builtins` anyway; and `SetDefaultOptions` mutates package-global state consulted by every `sdk.New` call (`v1/sdk/options.go:23-33`), which is process-level hidden coupling rather than instance configuration.

**4. Are examples sufficient to use the API correctly without reading internals?**
Yes. Coverage is layered: 22 compilable godoc `Example*` functions in `v1/rego/example_test.go:30-985` (including error handling, tracing, transactions, custom builtins with caching and nondeterminism), 2 for the compiler in `v1/ast/example_test.go:13-55`; complete copy-paste programs in the docs for both SDK (`docs/docs/integration.md:195-263`) and rego-package usage (`docs/docs/integration.md:373-415`); a full extension tutorial with a working plugin ("Putting It Together", `docs/docs/extensions.md:220`); and the REST API documented per-endpoint with request/response samples (`docs/docs/rest-api.md:48-2103`). Notably, the docs' SDK example depends on the public `v1/sdk/test` mock-server package (`docs/docs/integration.md:204`), making the example actually runnable. Gaps: no godoc Example functions exist for `v1/sdk` itself (only prose docs), and correct concurrent use of `sdk.OPA` relies on comments like "This function is threadsafe" (`v1/sdk/opa.go:299,435`) rather than demonstrated patterns.

## Architectural Decisions

1. **One module, two eras: root = v0 shims, `v1/` = current API.** Rather than a separate `opa/v2` module, OPA kept a single module (`go.mod:1`) and relocated the living API under `v1/`, converting every legacy package into aliases over the new implementation (`rego/rego.go:193`, `sdk/opa.go:13-18`). This gives free cross-compatibility of types between v0 and v1 callers within the same dependency graph, at the cost of shipping the deprecated surface forever — the banner promises the shims "will remain for the lifetime of OPA v1.x" (`rego/doc.go:5`).
2. **Type aliases as the compatibility mechanism.** Because shims are `type X = v1.X`, values flow between old and new APIs without conversion; but the project must document mixing v0 and v1 *packages* in one program as unsupported due to differing default Rego versions (`docs/docs/v0-compatibility.md:147-151`) — a semantic, not syntactic, boundary.
3. **Functional-options everywhere, but split by phase.** Construction-time options are `func(r *Rego)` (`rego.New`, `v1/rego/rego.go:1414`), while evaluation-time options are `EvalOption`s restricted to what is safe to change per-call (`EvalContext`, `v1/rego/rego.go:93-130`); the doc comment states the invariant: "Any other options will need to be set on a new Rego object" (`rego/rego.go:38-39` in the v0 shim).
4. **Machine-readable contracts checked into the repo.** Builtin signatures (`builtin_metadata.json`, generated per `main.go:34`), capabilities per Rego version (`capabilities.json`, `main.go:33`), and per-builtin introduction versions (`v1/ast/version_index.json:1-20`, embedded at `v1/ast/capabilities.go:39`) turn language-surface stability into data that downstream tools (e.g., `opa capabilities`) can consume instead of prose.
5. **Extension via narrow interfaces registered globally.** Plugins implement a 5-method surface (`Factory.Validate/New` + `Plugin.Start/Stop/Reconfigure`, `v1/plugins/plugins.go:89-110`) and register against a name keyed into the YAML config (`v1/runtime/runtime.go:90-98`); the same pattern covers storage backends (`v1/runtime/runtime.go:100-107`). This keeps third-party code out of core while giving it first-class config, status, and discovery treatment.
6. **Test infrastructure as published API.** `v1/sdk/test` (bundle mock server) and `v1/sdk/test/testdata` ship as importable packages precisely so integrators can follow the documented workflow without hand-rolling servers (`v1/sdk/test/test.go:27`).

## Notable Patterns

- **Uniform deprecation scaffolding**: identical `doc.go` banners across 73 v0 packages, plus per-symbol `Deprecated:` notes that always name the replacement (`v1/ast/compile.go:592-600` "Use WithCapabilities instead"; `v1/loader/loader.go:499-517` "Use FileLoader.Filtered() instead").
- **Route-table-as-contract**: every HTTP endpoint is registered in one function with parallel method-not-allowed catch-alls so unsupported verbs fail loudly (`v1/server/server.go:895-928`), and both `v0/data` and `v1/data` coexist behind the same handler family — HTTP-level versioning mirroring the Go-level v0/v1 split.
- **Readiness as an API element**: the SDK makes initialization blocking-by-default with an opt-out channel (`Options.Ready`, `v1/sdk/options.go:54-57,110-115`) and documents plugin-status-driven readiness (`v1/sdk/opa.go:207-227`) — lifecycle is part of the contract, not an afterthought.
- **Docs and code share examples**: the website's snippets mirror the godoc examples' idioms (`PrepareForEval` then `Eval` with `EvalInput`), reducing drift between `docs/docs/integration.md:360-415` and `v1/rego/example_test.go:758`.
- **Lint-enforced API discipline**: `.golangci.yaml:94-102` encodes which deprecation warnings apply where, effectively making the deprecation policy executable.

## Tradeoffs

- **Alias-based compat vs. dead-weight surface**: keeping v0 wrappers for the lifetime of v1.x (`rego/doc.go:5`) doubles the number of import paths users must reason about and requires anti-pattern warnings (`docs/docs/v0-compatibility.md:147-151`) to prevent misuse.
- **Escape hatches vs. abstraction integrity**: exposing `ManagerOpts`, `Plugin()`, and raw `storage.Store` injection (`v1/sdk/options.go:68-70,90-93`) empowers integrators but means internal manager semantics can become de-facto public contract.
- **Global registration vs. testability/isolation**: `RegisterBuiltin1..4` and `runtime.RegisterPlugin` mutate process-global state (`rego/rego.go:219-242`; `v1/runtime/runtime.go:69-98` guarded by mutexes), which is simple but makes parallel testing with different registries impossible.
- **Richness vs. reviewability**: `v1/ast` alone exposes ~1,076 exported symbols (measured via `go doc -all`), `v1/rego` ~176 — excellent for tool authors, but every addition is a permanent semver obligation with no mechanical guard (no gorelease/apidiff in CI; verified by search and by `Makefile:161-173` where `check` is golangci-lint only).

## Failure Modes / Edge Cases

- **Mixing v0 and v1 packages** silently changes default Rego syntax expectations per object; the docs must warn users off it because nothing in the type system prevents it (`docs/docs/v0-compatibility.md:147-151`).
- **SDK default-options race window**: `sdk.New` reads `defaultOptions` under a mutex (`v1/sdk/opa.go:67-97`), but `SetDefaultOptions` taking effect mid-flight means behavior depends on call ordering across unrelated libraries linking the SDK (`v1/sdk/options.go:29-33`).
- **Blocking `sdk.New` footgun**: if callers pass no `Ready` channel, `New` blocks until plugins report ready; the docs example must remind users `Ready: make(chan struct{})` is "needed or else sdk.New will block" in the manual-trigger scenario (`docs/docs/integration.md:317`).
- **`NDBCache any` downcast**: passing a non-`builtins.NDBCache` value silently disables the cache rather than erroring (`v1/sdk/opa.go:310-316`).
- **Deprecated-but-present traps**: `Tracer`/`EvalTracer` remain callable in the shim (`rego/rego.go:73-75,452-457`) — safety nets rely on staticcheck adoption downstream, which OPA cannot enforce for consumers.
- **HTTP verb coverage**: correctness of the method-not-allowed catch-alls depends on hand-maintaining the list per resource; an undocumented resource would 404 rather than 405 (`v1/server/server.go:920-928`).

## Future Considerations

- Add a mechanical compatibility gate (e.g., gorelease/apidiff in CI) so breaking changes to `v1/**` exported symbols fail the build rather than being caught in review; today the only automated checks are linting (`Makefile:161-173`).
- Replace global-state knobs (`cmd.UserAgent`, `cmd.RootCommand`, `sdk.SetDefaultOptions`) with instance-scoped options; `cmd.UserAgent` already carries maintainer doubt ("consider this experimental… I have the hunch that we'll find a better way", `cmd/commands.go:19-23`).
- Tighten the `rego`↔`topdown` seam: either promote the tracer/cache/print types into a small stable surface package or mark `v1/topdown` clearly as supported-for-integrators to resolve the current half-internal status (`v1/topdown/doc.go:5-9` vs. its use in `v1/rego/rego.go:33-36` and `v1/sdk/opa.go:33-34`).
- Type `DecisionOptions.NDBCache` concretely (`builtins.NDBCache`) once import-cycle concerns allow, removing the `any`+downcast at `v1/sdk/opa.go:310-316`.
- Add godoc Examples for `v1/sdk` parity with `v1/rego`.

## Questions / Gaps

- No evidence found of an automated Go API compatibility check: searched for `apidiff`, `gorelease`, `api-diff`, and inspected `.github/workflows/pull-request.yaml` (lint + test jobs only) and the Makefile `check` target. Compatibility appears enforced by convention (semver, CHANGELOG.md, deprecation markers), not tooling.
- No formal statement found in-repo defining exactly which `v1/` packages are covered by the stability guarantee (e.g., is `v1/util` — a grab-bag of helpers, `v1/util/json.go:5` — considered public contract?). The `v1/doc.go:5-9` comment implies everything under `v1/` is the API, which would make even helper packages permanent obligations.
- Version-index metadata covers builtins only (`v1/ast/version_index.json`); no equivalent provenance tracking was found for Go API symbols themselves.
- The REST API's long-term versioning strategy (beyond coexisting `/v0` and `/v1` data routes, `v1/server/server.go:895-905`) is not addressed in the docs reviewed; `docs/docs/rest-api.md` documents both without a deprecation timeline.

---

Generated by Dimension 24.01 (Public API Surface) against `opa`.
