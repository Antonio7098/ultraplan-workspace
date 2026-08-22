# Source Analysis: opa

## 22.01 Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | opa (Open Policy Agent) |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go (module `github.com/open-policy-agent/opa`, go 1.25), single module plus one test-only side module (`e2e/go.mod`) |
| Analyzed | 2026-08-21 |

## Summary

OPA is a single Go module organized as a layered set of top-level packages: a parser/AST layer (`v1/ast`), an evaluation engine (`v1/topdown`), a query API (`v1/rego`), storage abstractions (`v1/storage` with `inmem` and `disk` backends), bundle handling (`v1/bundle`), a plugin framework (`v1/plugins` with bundle/logs/status/discovery/rest sub-plugins), an HTTP server (`v1/server`), a runtime assembly layer (`v1/runtime`), an embedding SDK (`v1/sdk`), the CLI (`cmd`), and shared internals (`internal/`, enforced private by the Go compiler).

Since v1.0, the codebase follows a strict two-tier versioning scheme: all implementation lives under `v1/`, while root-level packages are thin, deprecated compatibility shims consisting almost entirely of type aliases and delegating wrappers (e.g., `storage/storage.go:11-53`, `rego/rego.go:11-24`). Dependency direction is overwhelmingly root-shim → `v1`; only two edges run the other way across that boundary (`v1/capabilities → capabilities` for embedded data files, `v1/tester → cmd/formats`).

The evaluation core is genuinely separable from the operational shell: `go list -deps ./v1/rego ./v1/topdown` pulls in none of `plugins`, `server`, `sdk`, `runtime`, or `download`, so "use the tool system without pulling in the entire runtime" is satisfied at the library level. Optional engines (wasm, external OPA-SDK targets) are kept out of core via a registry populated by side-effect imports (`internal/rego/opa/engine.go:51-55`, `v1/features/wasm/wasm.go:16-20`). Weaknesses are concentrated at the middle of the stack: several first-level package pairs are mutually entangled at module granularity (server ↔ plugins, download ↔ plugins, loader ↔ bundle/util, cmd ↔ tester), `v1/util` is a grab-bag imported by everything, and no lint rule (e.g., depguard) or CI check enforces the dependency directions — they hold by convention.

## Rating

**7/10.** Clear, compiler-assisted model: `internal/` visibility is machine-enforced, the v0/v1 split is consistently annotated (34/34 shim packages carry identical deprecation doc comments; staticcheck configured to keep deprecation noise out of `v1/`, `.golangci.yaml:95-98`), and the core/runtime separation is real and verifiable via `go list -deps`. It falls short of 8–9 because boundary rules between mid-stack packages are implicit — there are five bidirectional first-level pairs, one library→CLI dependency leak (`v1/tester/reporter.go:17`), and zero automated dependency-direction checks (no depguard in `.golangci.yaml`, no `go list` boundary assertion in CI workflows or Makefile).

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.go:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Package structure | Single Go module `github.com/open-policy-agent/opa`; ~35 first-level packages mirrored under `v1/` | `go.mod:1-3`, directory listing of `v1/` |
| v0 shim pattern | Root `rego` re-exports only type aliases and delegating functions into `v1/rego` | `rego/rego.go:14-24`, `rego/rego.go:31-33`, `rego/rego.go:630-637` |
| v0 shim is pure delegation | Root `storage` contains only wrapper funcs calling `v1/storage`; no own logic | `storage/storage.go:9-53` |
| Deprecation annotation strategy | Identical `// Deprecated:` doc comment on every v0 shim package (34/34 checked: ast, bundle, compile, config, cover, debug, dependencies, download, format, hooks, ir, keys, loader, logging, metrics, plugins, profiler, refactor, rego, repl, resolver, runtime, schemas, sdk, server, storage, tester, topdown, tracing, types, util, version, capabilities, features) | `ast/doc.go:7-9`, `topdown/doc.go:13-15` |
| Lint enforces deprecation asymmetry | staticcheck deprecation warnings suppressed outside `v1/` so shims can reference v1 freely | `.golangci.yaml:94-98` |
| Stated design intent for v0/v1 | "Package v1 implements the v1 API… Most packages outside the v1 API are deprecated" ; docs describe v0-compatibility mode and discourage mixing v0+v1 imports | `v1/doc.go:7-11`; `docs/docs/v0-compatibility.md:147-151` |
| Internal API separation (compiler-enforced) | Shared `internal/` tree used widely by public v1 packages but invisible to external importers | `v1/ast/compile.go` (imports `internal/...`), `internal/` dir listing (38 entries incl. planner, wasm, providers/aws) |
| Internal subpackages inside modules | `ast/internal/{scanner,tokens}` and `storage/internal/{errors,ptr}` hide parser/store internals | `v1/ast/internal/`, `v1/storage/internal/` dirs; scanner consumed by `v1/ast/parser.go` |
| Dependency direction root → v1 | Entire `./v1/...` closure touches exactly two non-v1, non-internal OPA packages: `capabilities` (embedded JSON) and `cmd/formats` | `go list -deps ./v1/...` output; `v1/capabilities/capabilities.go:8`; `v1/tester/reporter.go:17` |
| Core usable without runtime shell | `go list -deps ./v1/rego ./v1/topdown` contains zero of plugins/server/sdk/runtime/download; closure sizes: ast=221, topdown=294, rego=346, sdk=502 total packages | `go list -deps` runs on this checkout |
| Layering: engine over AST/storage | `v1/topdown` imports only `v1/{ast,metrics,resolver,storage,tracing,types,util}` + internal helpers | `go list -f '{{.Imports}}' ./v1/topdown` |
| Query API composes engine + bundles + storage | `v1/rego` imports ast, bundle, ir, loader, storage/inmem, topdown, tracing (no server/plugins) | `go list -f '{{.Imports}}' ./v1/rego` |
| Server depends on plugin framework | `v1/server` imports `v1/plugins`, `v1/plugins/bundle`, `v1/plugins/server/decoding`, `v1/plugins/status` | `v1/server/server.go:40-42` |
| Reverse edge: plugins → server | Decision-logs plugin consumes `server.Info` and `server/types` | `v1/plugins/logs/plugin.go:31`, `v1/plugins/logs/plugin.go:716`, `v1/plugins/bundle/status.go:17` |
| Reverse edge: download → plugins | Downloader references `plugins.TriggerPeriodic` config types | `v1/download/download.go:27`, `v1/download/download.go:183` |
| Forward edge: plugins/bundle → download | Bundle plugin owns downloader lifecycle | `v1/plugins/bundle/config.go:17`, `v1/plugins/bundle/plugin.go:26` |
| Cycle broken via subpackage split | `v1/loader/loader.go:23` imports `v1/bundle`, while `v1/bundle` only imports leaf `v1/loader/filter` | `v1/loader/loader.go:23`, `v1/bundle/file.go:16`, `v1/bundle/filefs.go:10` |
| Grab-bag utility coupling | `v1/util` ↔ `v1/loader`: util imports `loader/extension` while loader imports util | `v1/util/json.go:17` vs `go list` edge loader→util |
| CLI ↔ tester entanglement | `v1/tester` (library) imports CLI flag helper `cmd/formats`; `cmd` imports `v1/tester` | `v1/tester/reporter.go:17`; `cmd/test.go` |
| Optional engine extension point | Registry with panic-on-duplicate `RegisterEngine`; wasm engine registered via blank/side-effect import in `features/wasm` | `internal/rego/opa/engine.go:49-60`, `v1/features/wasm/wasm.go:12-20` |
| Embedded capability data isolated | Historical capabilities JSON embedded as `embed.FS` in its own tiny package | `capabilities/capabilities.go:8-12` |
| Test-only side module | `e2e/` is a separate Go module with `replace` to parent, explicitly documented as never-to-be-imported | `e2e/go.mod:1-5`, `e2e/README.md:1-7` |
| E2E scenario organization | Per-scenario e2e packages (authz, http, logs, oci, tls, wasm…) live under `v1/test/e2e` rather than polluting product packages | `v1/test/e2e/` dir listing |
| Public API documentation | Docs name `v1/rego` and `v1/sdk` as the integration surfaces | `docs/docs/integration.md:29-36`, `docs/docs/integration.md:174-186` |
| In-v1 deprecation lifecycle | 91 `Deprecated:` markers inside `v1/` packages (symbols, not packages) | e.g., `v1/ast/builtins.go:332` |
| No automated boundary checks | No depguard/import-restriction linter configured; Makefile `check` target runs golangci-lint only | `.golangci.yaml:4-22`, `Makefile:161-173` |

## Answers to Dimension Questions

**1. Are modules cleanly separated?**
Largely yes at the extremes, muddier in the middle. The bottom layers are clean: `v1/ast` depends only on `v1/{types,util,capabilities,metrics}` plus internals, and `v1/topdown` adds only storage/resolver/tracing (see evidence table). The top layers compose downward cleanly (`v1/server/server.go:40-42` → plugins → rego → topdown → ast). But four mid-stack pairs are entangled at first-level granularity: server↔plugins (`v1/plugins/logs/plugin.go:31` vs `v1/server/server.go:40`), download↔plugins (`v1/download/download.go:27` vs `v1/plugins/bundle/config.go:17`), loader↔bundle/util (`v1/loader/loader.go:23`, `v1/util/json.go:17`), and cmd↔tester (`v1/tester/reporter.go:17`). These are cycles at module granularity, broken only because Go forbids package-level cycles — typically via leaf subpackages.

**2. Do dependencies flow in one direction?**
Within the v0→v1 axis, yes and it is nearly perfect: of every OPA package reachable from `./v1/...`, exactly two non-v1 edges exist (`v1/capabilities/capabilities.go:8`, which is a data-only `embed.FS` package at `capabilities/capabilities.go:8-12`, and `v1/tester/reporter.go:17`). Within the v1 layer graph, flow is mostly acyclic bottom-up (ast → topdown → rego → {plugins, sdk} → server → runtime), with the exceptions listed above. There is no enforcement mechanism proving this stays true: no depguard rules exist (`.golangci.yaml:4-22`) and CI workflows contain no dependency-graph assertions.

**3. Can modules be used independently?**
The headline question — "can you use the tool system without pulling in the entire runtime?" — is answered affirmatively by the build graph: `go list -deps ./v1/rego ./v1/topdown` returns zero occurrences of plugins, server, sdk, runtime, or download. A consumer of the policy engine needs only the ast/topdown/rego/storage subset (346 total packages including stdlib for rego, versus 502 for the full sdk). The SDK deliberately composes the full stack itself (`v1/sdk/opa.go:25-27`). Two caveats: (a) `v1/rego` transitively includes wasm-planning internals (`internal/planner`, `internal/compiler/wasm`) even if you never target wasm; (b) independence is a property of the import graph, not of declared interfaces — nothing prevents a future rego→plugins import except review.

**4. Are public APIs distinguished from internal ones?**
Yes, through three mechanisms. First, Go's `internal/` path visibility makes the 38-package `internal/` tree compiler-enforced private to the module (`internal/rego/opa/engine.go:51` is unreachable externally). Second, the v0 shims are uniformly marked `Deprecated` in doc comments (`ast/doc.go:7-9`), with linting configured to treat them as second-class (`.golangci.yaml:95-98`). Third, docs designate `v1/rego` and `v1/sdk` as supported integration surfaces (`docs/docs/integration.md:29-36`). What's missing is intra-module granularity: everything exported from any `v1/*` package is technically public to all sibling packages, and symbols like `plugins.Manager` end up consumed by lower-level packages such as the downloader (`v1/download/download.go:27,183`).

## Architectural Decisions

1. **Two-tier v0/v1 layout instead of a new module.** All logic moved to `v1/`, root packages became alias shims (`rego/rego.go:14-24`). This keeps a single `go.mod` (no multi-module sync burden) while giving v1 a fresh, syntax-default-changing API surface (`v1/doc.go:7-11`). Tradeoff: the shims must be maintained "for the lifetime of OPA v1.x" per `ast/doc.go:8`.

2. **Registry-based optional engines.** Rather than linking wasm into core, `internal/rego/opa` exposes `RegisterEngine`/`LookupEngine` (`internal/rego/opa/engine.go:51-60`), and `v1/features/wasm` self-registers on import (`v1/features/wasm/wasm.go:16-20`). This mirrors OPA's Rego-level design where builtins register into tables — extensibility without hard dependency.

3. **Cycle-breaking by subpackage extraction.** Where two first-level packages need each other, a leaf subpackage carries the shared piece: `loader/filter` for bundle↔loader (`v1/bundle/file.go:16`), `loader/extension` for util↔loader (`v1/util/json.go:17`). This works but signals the true module boundaries are finer than the directory names suggest.

4. **Test-dependency isolation in a separate Go module.** Heavy e2e deps (DB drivers) live in `e2e/go.mod` with a `replace` directive (`e2e/go.mod:1-5`) and a README forbidding imports (`e2e/README.md:5-7`) — protecting the main module's dep graph from bloat.

## Notable Patterns

- **Uniform shim generation**: all 34 root shim packages follow byte-for-byte identical deprecation wording, strongly suggesting tooling or a checklist drove the migration; staticcheck exclusions (`.golangci.yaml:95-98`) show the team tuned lint so shims don't trip their own deprecation warnings.
- **Data as a package**: versioned capabilities snapshots are embedded JSON in a dedicated package (`capabilities/capabilities.go:8-12`, 160+ `v*.json` files), keeping historical capability data out of `v1/capabilities` code.
- **Scenario-based e2e tree**: `v1/test/e2e/{authz,http,oci,tls,wasm,...}` gives cross-cutting behaviors their own packages rather than scattering integration tests inside product packages.
- **In-version deprecation**: deprecation is applied to individual exported symbols inside v1 too (91 sites, e.g., `v1/ast/builtins.go:332`), showing API hygiene continues post-1.0.

## Tradeoffs

- **Single module convenience vs. enforceable boundaries.** One `go.mod` keeps releases simple, but Go offers no way to say "v1/rego may not import v1/plugins"; only convention holds today.
- **Shim completeness vs. drift risk.** Because shims hand-forward every symbol (638 lines in `rego/rego.go` alone), every v1 API addition risks shim omission; there is no generated-code marker or test asserting shim parity.
- **Grab-bag `v1/util`.** Nearly everything imports it and it reaches into `loader/extension` (`v1/util/json.go:17`), making it a de-facto universal coupling point that caps how independent any package can be.

## Failure Modes / Edge Cases

- **Library depending on CLI code**: `v1/tester/reporter.go:17` importing `cmd/formats` means anyone embedding the tester library links in CLI flag machinery; if `cmd` ever grows heavier deps, library consumers pay. This is the single clearest layering violation found.
- **Bidirectional mid-stack pairs** (server↔plugins/logs, plugins/bundle↔download) make refactors brittle: changing `server.Info` ripples into decision-log plugins and vice versa (`v1/plugins/logs/plugin.go:716`).
- **Silent boundary erosion**: with no depguard/CI check, an accidental import (like the tester→cmd one) can land in review unnoticed; nothing fails mechanically.
- **Mixing v0 and v1 APIs in one binary** is explicitly unsupported and undefined-behavior territory (`docs/docs/v0-compatibility.md:147-151`), yet nothing in the build prevents it — only a doc warning.

## Future Considerations

- Add `depguard` (or a `go list -deps` assertion test, e.g., "no `cmd/` import outside `cmd/`") to make the observed layering mechanical rather than conventional — the tester→formats violation would have been caught by exactly this rule.
- Move `formats.Flag` constants out of `cmd/formats` into `v1/tester` or a neutral package to sever the last v1→root code dependency (the capabilities one is data-only and benign).
- Split `v1/util` along actual usage lines or fold hot paths into consumers, reducing the universal-import gravity well.
- Consider promoting the loader↔bundle seam (`loader/filter`, `loader/extension`) into a documented SPI so the extraction looks intentional.

## Questions / Gaps

- **No evidence found** of automated circular-dependency or layering tests: searches across `.github/workflows/*.yaml`, `Makefile:161-173`, `build/utils.sh`, and `.golangci.yaml` found golangci-lint and generic packaging scripts but no import-boundary assertions.
- **No evidence found** of generated shim-parity tests (a test asserting every exported v1 symbol has a v0 alias); the uniformity of shims suggests manual/tool-assisted authorship without a guard.
- Whether `v1/rego`'s transitive pull of wasm planning internals (`internal/planner`, `internal/compiler/wasm` via `go list -deps ./v1/rego`) is intentional payload or accepted cost could not be determined from code comments; the `Target("wasm")` option (`rego/rego.go`, `Target` delegation) implies it is by design.
- The `TODO: LINK TO V0 MIGRATION GUIDE` placeholder in `v1/doc.go:10` suggests the boundary documentation was incomplete at time of writing.

---

Generated by `22.01-package-and-module-boundaries` against `opa`.
