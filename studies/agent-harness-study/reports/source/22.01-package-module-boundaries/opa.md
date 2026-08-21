# Source Analysis: opa

## Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | opa |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go (single Go module, module path `github.com/open-policy-agent/opa`, `go 1.25.0`) |
| Analyzed | 2026-08-21 |

## Summary

OPA is structured as a single Go module that exposes two parallel public trees: a deprecated top-level v0 API (`ast`, `rego`, `topdown`, `runtime`, `server`, `storage`, `plugins`, `bundle`, `loader`, `sdk`, `hooks`, `compile`, `format`, `version`, `capabilities`, `tracing`, `metrics`, `logging`, `keys`, `dependencies`, `features`, `profiler`, `refactor`, `debug`, `schemas`, `types`, `util`, `wasm`, `cover`, `download`, `repl`, `tester`, `test`, `resolver`, `ir`, `format`) and a `v1/` tree that owns the implementation (`v1/ast`, `v1/rego`, `v1/topdown`, `v1/runtime`, `v1/server`, `v1/storage`, `v1/plugins`, `v1/bundle`, `v1/loader`, `v1/sdk`, `v1/hooks`, `v1/compile`, `v1/format`, `v1/version`, `v1/capabilities`, `v1/tracing`, `v1/metrics`, `v1/logging`, `v1/keys`, `v1/dependencies`, `v1/features`, `v1/profiler`, `v1/refactor`, `v1/debug`, `v1/schemas`, `v1/types`, `v1/util`, `v1/wasm`, `v1/cover`, `v1/download`, `v1/repl`, `v1/resolver`, `v1/ir`).

The v0 root packages are pure re-export wrappers: each `.go` file at the root is dominated by `type X = v1.X` aliases and small factory functions that delegate to `v1.X`. For example `ast/term.go:11-60` is entirely `Location`, `Value`, `InterfaceToValue`, `ValueFromReader`, `As`, `Resolver`, `ValueResolver`, `UnknownValueErr`, `IsUnknownValueErr` aliases that point to `v1/ast`. The legacy tree is declared deprecated in each package's `doc.go` (e.g. `ast/doc.go:5`, `runtime/doc.go:7`, `server/doc.go:7`, `storage/doc.go:7`, `sdk/doc.go:5`).

The actual implementation lives under `v1/`. Public/private boundaries are enforced by an `internal/` tree (36 subpackages) plus extra `internal/` subtrees under `v1/ast/internal` (scanner, tokens), `v1/storage/internal` (errors, ptr), `v1/runtime/info`, and several `topdown/*` subpackages (`cache`, `print`, `copypropagation`, `builtins`, `lineage`). A separate `cmd/` subpackage builds the CLI via `main.go:11`. The CLI is a leaf: it imports `internal/`, `v1/`, and `cmd/internal/`, but nothing imports `cmd/`.

The dependency direction is largely one-way: `v1/ast` and `v1/storage` depend only on low-level utilities; `v1/topdown` adds `v1/ast`, `v1/storage`, metrics, tracing; `v1/rego` adds `v1/topdown` plus internal helpers; `v1/server` uses `v1/rego`, `v1/topdown`, `v1/ast`, `v1/storage`; `v1/sdk` and `v1/runtime` consume everything. No cycles were observed in the v1 layer.

The critical boundary question — *can you use the tool system without pulling in the entire runtime?* — is answered with a partial yes. Pure parsing/compilation (`v1/ast`) and direct evaluation (`v1/topdown`) can be used without the runtime or the HTTP server. The `v1/rego` package is a moderately heavy but self-contained evaluation API. However, the `v1/sdk` package imports `v1/server` (the 3,221-line HTTP server), so the SDK despite being a "library" entry point pulls in the entire HTTP governance surface.

## Rating

7 / 10 — Clear model, clean v0/v1 duality, explicit internal/ external markers, and a tightly enforced (lint-checked) boundary against the legacy root. Deductions for: the SDK dependency on the full HTTP server, the fact that the v0 root is a permanent compatibility shim that doubles the public surface, and the absence of automated boundary tests (no `go test` enforces "import nothing heavier than v1/ast").

## Evidence Collected

Every entry cites a file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Module identity | Single Go module, `go 1.25.0` | `go.mod:1-3` |
| v0/v1 split declared | `v1/doc.go` documents v1 as the new API, root as deprecated v0 | `v1/doc.go:5-8` |
| Legacy wrapper example | `ast` package is entirely type aliases to `v1/ast` | `ast/term.go:11-60` |
| Legacy wrapper example | `ast/parser.go` aliases `RegoV1`, `RegoV0`, `Parser`, `ParserOptions` | `ast/parser.go:13-46` |
| Legacy wrapper example | Default Rego version differs by surface (v0 retains v0 default) | `ast/parser.go:16-17` |
| Legacy wrapper example | `ast/compile.go` aliases `Compiler`, `CompilerStage`, `EvalMode*` | `ast/compile.go:14-58` |
| Legacy wrapper example | `hooks` package aliases `Hook`, `Hooks`, `ConfigHook` | `hooks/hooks.go:27-46` |
| Legacy wrapper example | `sdk/opa.go` aliases `OPA`, `Options`, `DecisionResult`, `PartialResult` | `sdk/opa.go:14-30` |
| Deprecation noted | `ast/doc.go` says deprecated for v0.x projects | `ast/doc.go:5-7` |
| Deprecation noted | `runtime/doc.go` says deprecated | `runtime/doc.go:7-9` |
| Deprecation noted | `server/doc.go`, `storage/doc.go` say deprecated | `server/doc.go:7-9`, `storage/doc.go:7-9` |
| Deprecation noted | `sdk/doc.go` says deprecated | `sdk/doc.go:5-7` |
| Deprecation noted | `hooks/doc.go` says deprecated | `hooks/doc.go:5-7` |
| Internal subtree | 36 internal subpackages under `internal/` | directory listing (e.g. `internal/cmd/`, `internal/wasm/`, `internal/rego/opa/`, `internal/planner/`, `internal/runtime/init/`) |
| v1/ast isolation | `v1/ast/term.go` only imports `v1/ast/json`, `v1/ast/location`, `v1/util`, stdlib | `v1/ast/term.go:7-25` |
| v1/ast isolation | `v1/ast/parser.go` imports only `v1/ast/internal/scanner`, `v1/ast/internal/tokens`, `v1/ast/json`, `v1/ast/location`, `v1/util`, `v1/types` | `v1/ast/parser.go:7-30` |
| v1/ast isolation | `v1/ast/compile.go` imports `internal/debug`, `internal/gojsonschema`, `v1/ast/location`, `v1/metrics`, `v1/types`, `v1/util` | `v1/ast/compile.go:14-23` |
| v1/ast isolation | `v1/ast/annotations.go` imports `internal/deepcopy`, `v1/ast/json`, `v1/util` | `v1/ast/annotations.go:14-20` |
| v1/topdown isolation | `v1/topdown/query.go` imports `v1/ast`, `v1/metrics`, `v1/resolver`, `v1/storage`, `v1/topdown/builtins`, `v1/topdown/cache`, `v1/topdown/copypropagation`, `v1/topdown/print`, `v1/tracing` — no server/sdk/runtime | `v1/topdown/query.go:10-19` |
| v1/storage isolation | `v1/storage/disk/txn.go` imports `v1/metrics`, `v1/storage`, `v1/storage/internal/errors`, `v1/storage/internal/ptr`, `v1/util` | `v1/storage/disk/txn.go:13-21` |
| v1/storage isolation | `v1/storage/inmem/inmem.go` imports `internal/merge`, `v1/ast`, `v1/storage`, `v1/storage/internal/errors`, `v1/util` | `v1/storage/inmem/inmem.go:23-31` |
| v1/storage interface | `Store` interface defines `NewTransaction`, `Read`, `Write`, `Commit`, `Abort`, `Truncate` | `v1/storage/interface.go:20-44` |
| v1/rego mediated by internal | `v1/rego/rego.go` imports `internal/bundle`, `internal/compiler/wasm`, `internal/future`, `internal/planner`, `internal/rego/opa`, `internal/wasm/encoding` | `v1/rego/rego.go:18-23` |
| v1/server composition | `v1/server/server.go` imports `v1/ast`, `v1/bundle`, `v1/config`, `v1/hooks`, `v1/logging`, `v1/metrics`, `v1/plugins`, `v1/plugins/bundle`, `v1/plugins/server/decoding`, `v1/plugins/server/encoding`, `v1/plugins/status`, `v1/rego`, `v1/storage`, `v1/topdown`, `v1/tracing`, `v1/util` — no v1/sdk, no v1/runtime | `v1/server/server.go:33-58` |
| v1/server composition | `v1/server/compile_handler.go` uses `internal/compile`, `v1/rego/compile`, `v1/server/failtracer` | `v1/server/compile_handler.go:16-27` |
| v1/server composition | `v1/server/authorizer/authorizer.go` uses `v1/rego`, `v1/server/identifier`, `v1/topdown/cache`, `v1/topdown/print` | `v1/server/authorizer/authorizer.go:14-22` |
| SDK pulls in server | `v1/sdk/opa.go` imports `v1/server` (HTTP server, 3,221 lines) plus server/types | `v1/sdk/opa.go:30-31` |
| SDK heavy imports | `v1/sdk/opa.go` import list: ast, bundle, hooks, logging, metrics, plugins, plugins/discovery, plugins/logs, rego, runtime/info, server, server/types, storage, topdown, topdown/builtins, topdown/cache, topdown/print, util, version + internal/ref, internal/uuid | `v1/sdk/opa.go:18-38` |
| v1/runtime composition | `v1/runtime/runtime.go` imports 11 internal/* packages and 16 v1/* packages | `v1/runtime/runtime.go:35-67` |
| Runtime uses init helper | `internal/runtime/init/init.go` exposes `InsertAndCompile`; consumed by plugin manager and runtime | `internal/runtime/init/init.go:5-22` |
| Extensibility via hooks | `v1/hooks/hooks.go` declares `Hook` (any), `Hooks`, `ConfigHook`, `ConfigDiscoveryHook`, `InterQueryCacheHook`, `InterQueryValueCacheHook`, `BundlePreActivateHook` | `v1/hooks/hooks.go:32-100` |
| Plugin manager | `v1/plugins/plugins.go` 1,426 lines implements `Factory`, `Manager`, plugin lifecycle | `v1/plugins/plugins.go:1-40` |
| Plugin subpackages | `v1/plugins/bundle`, `v1/plugins/discovery`, `v1/plugins/logs`, `v1/plugins/rest`, `v1/plugins/server/{decoding,encoding,metrics}`, `v1/plugins/status`, `v1/plugins/logger/file` | directory listing under `v1/plugins/` |
| REST plugin modular | `v1/plugins/rest/auth.go`, `auth_tls.go`, `aws.go`, `azure.go`, `gcp.go` | `v1/plugins/rest/` |
| Topdown subpackages | `v1/topdown/builtins`, `v1/topdown/cache`, `v1/topdown/copypropagation`, `v1/topdown/lineage`, `v1/topdown/print` | `v1/topdown/` |
| CLI is a leaf | `cmd/commands.go` builds the root command with cobra | `cmd/commands.go:7-50` |
| CLI imports CLI-internal | `cmd/eval.go` imports `cmd/internal/env`, `cmd/formats`, `internal/file/url`, `internal/presentation`, `v1/ast`, `v1/ast/location`, `v1/bundle`, `v1/compile`, `v1/cover`, `v1/loader`, `v1/metrics`, `v1/profiler`, `v1/rego`, `v1/runtime/info`, `v1/topdown`, `v1/topdown/lineage`, `v1/util` | `cmd/eval.go:20-37` |
| Lint rule for v0 deprecation | `staticcheck` exception lets the v0 root keep using deprecated symbols | `.golangci.yaml:95-98` |
| Lint exclusion for vendored code | `internal/gojsonschema` excluded from lint | `.golangci.yaml:103-105` |
| Test helpers for SDK | `v1/sdk/test/test.go` exports `MockBundle`, `MockOCIBundle`, `Ready`, `Server` | `v1/sdk/test/test.go:1-60` |
| Test helpers for util | `v1/util/test/` provides `tempfs`, `tempus`, `populate`, `benchmark`, `zeroreader`, `ci_skip` | `v1/util/test/` |
| No cross-source note | `main.go` only imports `github.com/open-policy-agent/opa/cmd` | `main.go:7-12` |

## Answers to Dimension Questions

1. **Are modules cleanly separated?**

   Yes, with one important caveat. The v1 layer is cleanly separated along a strict hierarchy: `v1/ast` (parser/compiler) and `v1/util`/`v1/types`/`v1/metrics` at the bottom; `v1/storage` (interface + `inmem` + `disk`) above; `v1/topdown` (evaluator) above that; `v1/rego` (evaluation façade) above `v1/topdown`; `v1/server` (HTTP) above `v1/rego`; and `v1/sdk`/`v1/runtime` at the top. Verifiable from import lists at `v1/ast/term.go:7-25`, `v1/storage/disk/txn.go:13-21`, `v1/topdown/query.go:10-19`, `v1/server/server.go:33-58`, `v1/sdk/opa.go:18-38`.

   The caveat is that `v1/sdk` imports `v1/server` (`v1/sdk/opa.go:30`), so the SDK heavyweight module is not really a "library" surface — it pulls in the full HTTP server. The `/workspace/studies/agent-harness-study/sources/opa/v1/sdk/opa.go:30-31` line is the canonical example of the boundary violation.

2. **Do dependencies flow in one direction?**

   Yes across v1. Observed direction: `v1/ast` → `v1/storage` → `v1/topdown` → `v1/rego` → `v1/server`. `v1/sdk` and `v1/runtime` consume all of them. The reverse direction is not present: `v1/ast` does not import `v1/topdown`, `v1/server`, `v1/sdk`, `v1/runtime`, `v1/rego`, `v1/plugins` (verified by grep across `v1/ast/*.go` excluding tests). `v1/storage` does not import `v1/topdown`, `v1/server`, `v1/sdk`, `v1/runtime` (verified for `v1/storage/inmem/inmem.go:23-31` and `v1/storage/disk/txn.go:13-21`). `v1/topdown` does not import `v1/server`, `v1/sdk`, `v1/runtime`, `v1/plugins` (verified for `v1/topdown/query.go:10-19`).

   The exception is the v0/v1 pair: the v0 root packages import `v1/` (e.g. `ast/term.go:11` imports `v1/ast`), but no `v1/` package imports the v0 root. This is an asymmetric edge, but it is one-way and confined to the entry-point packages.

   Internal subpackages import `v1/` heavily (e.g. `internal/planner/planner.go:14-18`, `internal/runtime/init/init.go:15-21`, `internal/pathwatcher/utils.go:14-18`), but nothing in `v1/` imports `internal/`. This is the canonical Go internal-package boundary, enforced by the compiler.

3. **Can modules be used independently?**

   - `v1/ast` (parsing, AST, compiler) — yes, standalone. Imports only `v1/ast/json`, `v1/ast/location`, `v1/ast/internal/*`, `v1/util`, `v1/types`, `v1/metrics`, `internal/deepcopy`, `internal/debug`, `internal/gojsonschema`, `internal/semver`. No server, no runtime. Evidence: `v1/ast/term.go:7-25`, `v1/ast/parser.go:7-30`, `v1/ast/compile.go:14-23`.
   - `v1/topdown` (evaluator) — yes, standalone if you supply a `storage.Store`. Evidence: `v1/topdown/query.go:10-19`.
   - `v1/storage/inmem` — yes, standalone. Evidence: `v1/storage/inmem/inmem.go:23-31`.
   - `v1/storage/disk` — yes, standalone. Evidence: `v1/storage/disk/txn.go:13-21`.
   - `v1/rego` — yes, but moderately heavy: pulls in `internal/bundle`, `internal/compiler/wasm`, `internal/future`, `internal/planner`, `internal/rego/opa`, `internal/wasm/encoding`, `v1/ast`, `v1/bundle`, `v1/ir`, `v1/loader`, `v1/loader/filter`, `v1/metrics`, `v1/resolver`, `v1/storage`, `v1/storage/inmem`, `v1/topdown`, `v1/topdown/builtins`, `v1/topdown/cache`, `v1/topdown/print`, `v1/tracing`, `v1/types`, `v1/util`. Evidence: `v1/rego/rego.go:18-39`.
   - `v1/sdk` — no, heavyweight. Brings `v1/server` (3,221 lines) and the plugin manager. Evidence: `v1/sdk/opa.go:30-31`.
   - `v1/runtime` — no, heaviest. Brings `internal/compiler`, `internal/config`, `internal/distributedtracing`, `internal/logging`, `internal/metricsexport`, `internal/pathwatcher`, `internal/prometheus`, `internal/ref`, `internal/runtime/init`, `internal/uuid`, `internal/versioncheck` plus the majority of `v1/`. Evidence: `v1/runtime/runtime.go:35-67`.

4. **Are public APIs distinguished from internal ones?**

   Yes, via three mechanisms:
   - Go's `internal/` directory rule (compiler-enforced): `internal/` (36 subpackages) and `v1/ast/internal/`, `v1/storage/internal/`, `v1/runtime/info/` (info is the metadata package for the runtime; not the runtime itself) cannot be imported by anything outside the parent tree.
   - Deprecation markers on the v0 root packages (`ast/doc.go:5-7`, `runtime/doc.go:7-9`, `server/doc.go:7-9`, `storage/doc.go:7-9`, `sdk/doc.go:5-7`, `hooks/doc.go:5-7`).
   - Documentation comments on the v1 packages (`v1/ast/doc.go:5-7`, `v1/topdown/doc.go:5-9`, `v1/util/doc.go:5`, `v1/hooks/hooks.go:11-26`).
   - Lint exception that suppresses `staticcheck` "deprecated" warnings in the v0 root — confirming the v0 root is *expected* to be the deprecated shim (`.golangci.yaml:95-98`).

   Additionally, the v0 root packages and the v1 packages are deliberately kept structurally identical so users can migrate by changing only the import path.

## Architectural Decisions

- **v0/v1 dual-tree as a compatibility affordance.** The repo keeps both `ast` and `v1/ast` wired up. The v0 root is a thin (one- or two-line) re-export layer for every exposed type, function, and constant. `ast/term.go:11-60` is the canonical example: `type Location = v1.Location`, `type Value = v1.Value`, `func InterfaceToValue(x any) (Value, error) { return v1.InterfaceToValue(x) }`. The decision is documented in `v1/doc.go:5-8` and reinforced by `sdk/doc.go:5-7` ("Deprecated: This package is intended for older projects transitioning from OPA v0.x").
- **Plugin manager as the composition root.** `v1/plugins/plugins.go` (1,426 lines) is treated as the canonical way to wire features. Each plugin (bundle, discovery, logs, status, server/decoding, server/encoding, server/metrics, rest, logger) lives in its own subpackage and registers a `plugins.Factory` (`v1/plugins/plugins.go:42-80`). This keeps cross-plugin coupling explicit and at one location.
- **Hooks as runtime extension points.** `v1/hooks/hooks.go:32-100` defines `Hook` as `any` and a set of marker interfaces (`ConfigHook`, `ConfigDiscoveryHook`, `InterQueryCacheHook`, `InterQueryValueCacheHook`, `BundlePreActivateHook`). The compiler/runtime checks at call time whether a hook implements the appropriate interface, keeping the surface additive.
- **Internal/ as the playground.** Engineering helpers that need to span the v1 boundary without becoming public API live in `internal/`. Examples: `internal/planner` (regulates the query planner that `v1/rego` consumes), `internal/runtime/init` (the bundle-and-compile bootstrap that `v1/runtime` and `v1/plugins` both use), `internal/ref` (tiny RFC 6901 / data-path parser shared by `v1/sdk` and `v1/plugins/bundle`).
- **Wasm is treated as an optional feature.** `cmd/features.go:10` does `import _ "github.com/open-policy-agent/opa/v1/features/wasm"` to keep the wasm runtime out of the default link. `internal/rego/opa/engine.go:13-27` raises a typed error when no engine is registered, with a helpful message pointing users to the wasm-enabled build.
- **Auth and identifier are factored out of `v1/server`.** `v1/server/authorizer/authorizer.go` and `v1/server/identifier/` (token, tls, certs) are separate subpackages. The decisions are localized and a custom auth scheme can be dropped in without modifying server.go.

## Notable Patterns

- **Bit-for-bit type-aliased re-export.** Every public symbol in a v0 root file is `type X = v1.X` or `var X = v1.X`. Verified by reading `ast/term.go:11-60`, `ast/parser.go:11-46`, `ast/compile.go:14-58`, `ast/builtins.go:13-60`, `hooks/hooks.go:27-46`, `sdk/opa.go:14-30`. This is a deterministic, machine-greppable pattern.
- **Lint-enforced deprecation surface.** `.golangci.yaml:95-98` explicitly allows the v0 root to use deprecated symbols without SA1019 noise. This is a deliberate policy: legacy public symbols sit at the v0 root, and the new symbols live under `v1/`.
- **Interface-as-package.** `v1/storage/interface.go` defines `Store`, `Transaction`, `Trigger`, `Policy`, `MakeDirer`, `NonEmptyer`, `Closer` (`v1/storage/interface.go:15-64`). Implementations live in `v1/storage/inmem/` and `v1/storage/disk/`. Mocking is in `internal/storage/mock/`. Tests in `v1/storage/disk/txn_test.go` and `v1/storage/inmem/inmem_test.go` validate both implementations against the same interface.
- **Compile-only vs eval split.** `v1/rego/compile/compile.go` separates compile from eval. The `v1/rego/rego.go` (`v1/rego/rego.go:18-39`) and `v1/rego/compile/compile.go:14-23` import lists show they share `v1/ast`, `v1/ir`, `v1/util`, `v1/metrics` but eval pulls in `v1/topdown` while compile does not.
- **Two test helper packages.** `v1/sdk/test/test.go` (Server, MockBundle, MockOCIBundle) and `v1/util/test/` (tempfs, tempus, populate, benchmark, zeroreader, ci_skip) keep test-time helpers out of the main packages. Evidence: `v1/sdk/test/test.go:1-60`, `v1/util/test/`.
- **Five topdown subpackages.** `v1/topdown/builtins`, `v1/topdown/cache`, `v1/topdown/copypropagation`, `v1/topdown/lineage`, `v1/topdown/print`. Each can be reasoned about in isolation. `v1/topdown/builtins/builtins.go:15-16` imports only `v1/ast` and `v1/util`.

## Tradeoffs

- **Two parallel public trees.** Maintaining the v0 root as a permanent compatibility shim doubles the public surface and forces every contributor to remember the deprecation status. The decision is paid for in lint exceptions (`.golangci.yaml:95-98`) and an explicit migration link in each v0 doc.go.
- **SDK is heavy.** Pulling `v1/server` into `v1/sdk` (`v1/sdk/opa.go:30`) means a program that wants only "evaluate a policy" through the SDK still pulls in the HTTP server, server types, authorizer, identifier, decoding/encoding plugins, and tracing. There is no lightweight "evaluation-only" SDK API; users wanting that must drop down to `v1/rego` or `v1/topdown` directly.
- **Internal/ sprawl.** 36 internal subpackages provide clear boundaries but also create a leaky abstraction: when something needs to live in `internal/` because it crosses the v1 boundary, the v1 package surfaces become exposed only via that internal package, e.g. `internal/runtime/init` is used by both `v1/runtime` and `v1/plugins`, and `internal/planner` is used only by `v1/rego`.
- **Vendored JSON-Schema code.** `internal/gojsonschema` is 38 `.go` files plus `LICENSE-APACHE-2.0.txt`. It is excluded from lint (`.golangci.yaml:103-105`) and used by `v1/ast/compile.go` for type checking. This is a tradeoff for control over the schema validator versus depending on `santhosh-tekuri/jsonschema` (which is also a direct dependency at `go.mod:28`).
- **Identification of v1 subtrees.** Subpackages like `v1/ast/internal/`, `v1/storage/internal/`, `v1/runtime/info/` are `internal/`-style but live under `v1/`. They can be imported by any other `v1/...` package but not by external consumers. The convention is consistent but not enforced by tooling.

## Failure Modes / Edge Cases

- **`v1/sdk` link bloat.** Any consumer of `v1/sdk` imports the entire HTTP server, including its auth, identifier, encoding, decoding, and metrics subpackages. No filesystem evidence of a build-tag-gated path to a lighter SDK.
- **Two `eval` implementations.** `v1/rego` (regal evaluation) and `v1/topdown` (raw evaluation) are both reachable. A user importing `v1/rego` pulls in `internal/compiler/wasm`, `internal/planner`, `internal/rego/opa`, `internal/wasm/encoding` (`v1/rego/rego.go:18-23`), which is a meaningful dependency footprint.
- **Wasm optionality requires a build tag.** No evidence in the repo of a `//go:build !wasm` or similar tag separating wasm-enabled builds. The wasm registry is empty unless the consumer imports `v1/features/wasm` (see `cmd/features.go:10`). Consumers who forget the import get the typed error `ErrEngineNotFound` at `internal/rego/opa/engine.go:13`.
- **Linter-only guarantee.** The boundary between v0 (deprecated) and v1 (current) is enforced by `staticcheck` text rules (`.golangci.yaml:95-102`) and `--max-same-issues: 0` (`.golangci.yaml:108`). There is no automated test that, e.g., asserts `v1/ast` never imports `v1/rego`. A regression test would lift the rating.
- **Internal package boundary drift.** `internal/runtime/init` is imported by both `v1/runtime` and `v1/plugins` and by `internal/pathwatcher`. If a new `v1/...` package tries to use it, nothing prevents that — the Go compiler will allow it as long as the importer is under `github.com/open-policy-agent/opa/`. Detection is by convention, not by tool.

## Future Considerations

- A lightweight "library-only" SDK entry point that does not import `v1/server`. Concretely, a `v1/sdk/eval` package that exposes `Decision` and `Partial` without `Configure`, `Start`, or `manager`. The current `v1/sdk/opa.go:191-196` builds a `plugins.Manager` for every Decision call; for a "pure eval" path this is unnecessary.
- Automated boundary tests. Ideally a `hack/verify_boundaries.sh` or a `go test` that calls `go list -deps` and asserts that `v1/ast`, `v1/storage`, `v1/topdown`, `v1/util` do not depend on `v1/server`, `v1/sdk`, `v1/runtime`, `v1/plugins`. The staticcheck text rules at `.golangci.yaml:95-98` prove the team thinks about this; a CI gate would tighten it.
- Trim v0 root. Roughly 16 v0 root packages are pure wrappers. Once OPA 1.x reaches EOL, deletion of these files is a single-PR change because the file contents are entirely aliases. Documenting the deletion policy in `CHANGELOG.md` would help.
- Modularize `v1/sdk`. Splitting into `v1/sdk/eval`, `v1/sdk/config`, `v1/sdk/server` would let users import only the slice they need. The current structure forces a single import graph.
- Build-tag-gated wasm. Move the wasm registry into a separate `v1/internal/wasm/` package so consumer binaries can opt out at compile time instead of at runtime via `internal/rego/opa/engine.go:13`.

## Questions / Gaps

- **Boundary tests.** No `go test` enforces "v1/ast does not import v1/topdown" or "v1/topdown does not import v1/server". A search for `go test` files containing such an assertion came up empty. The guarantee is enforced by code review and lint rules only.
- **Deprecated v0 root deletion timeline.** The deprecation message in `sdk/doc.go:5-7` and `ast/doc.go:5-7` says "will remain for the lifetime of OPA v1.x". The repo does not state when v1.x is EOL. `CHANGELOG.md` does not enumerate v0 root removal.
- **Internal package rule enforcement.** `internal/` is enforced by Go's compiler. The intra-tree placement of `v1/runtime/info/`, `v1/ast/internal/`, `v1/storage/internal/` is convention. No test asserts that these subpackages are not re-exported.
- **Build-tag asymmetry.** `cmd/features.go:10` does `import _ "github.com/open-policy-agent/opa/v1/features/wasm"`. There is no equivalent for the SDK or runtime. The implications for binary size vs SDK size are not documented in the README.
- **Two evaluation-entry points.** `v1/rego` and `v1/topdown` both evaluate Rego. The boundary between them is documented loosely: `v1/rego/doc.go:5-6` ("Package rego exposes high level APIs for evaluating Rego policies") vs `v1/topdown/doc.go:5-9` ("Package topdown provides low-level query evaluation support"). A newcomer does not learn which to choose from the comments.
- **Plugin lifecycle test isolation.** `v1/plugins/plugins.go` is 1,426 lines. The public surface (Factory, Manager, etc.) is exercised in `v1/plugins/plugins_test.go`. Whether a plugin can be used without the discovery plugin (and how) is not obvious from the public docs.

---

Generated by `22.01-package-and-module-boundaries` against `opa`.
