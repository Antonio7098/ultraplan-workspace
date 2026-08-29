# Source Analysis: opa

## 21.01 Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | opa |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go (go.mod:1.25), Rego |
| Analyzed | 2026-08-28 |

## Summary

OPA exposes a coherent but compile-time-bound extension model centered on `plugins.Factory`/`plugins.Plugin` registered via `runtime.RegisterPlugin`, plus orthogonal hooks for config mutation, custom built-ins, custom storage, loader file extensions, external rule sources, and REST auth plugins. The plugin manager (`v1/plugins/plugins.go`) provides explicit lifecycle (`Start`/`Stop`/`Reconfigure`), status reporting (`StateOK/StateErr/StateWarn`), and discovery-driven reconfiguration. There is no dynamic `.so` or WASM plugin loading; all extensions must be linked at build time via Go imports. Isolation is logical (per-plugin status, mutex-protected manager) not sandboxed. Documentation in `docs/docs/extensions.md` is thorough for built-ins and plugins but marks loader extensions as EXPERIMENTAL; API stability is enforced via `v1/` import path and tests.

## Rating

**7 / 10 — Clear model with tests, explicit interfaces, and operational safeguards, but not durable under untrusted third-party isolation or true runtime dynamic loading.**

Rationale: Factory/Plugin interfaces are explicit and versioned (`v1/plugins`), manager lifecycle covers init/start/stop/reconfigure with graceful shutdown, status observability, discovery can hot-swap config and start new plugins without restart. Built-in and storage extension points are tested. Gaps that prevent 9-10: no process/memory sandbox, built-in registry is global non-thread-safe (`ast.RegisterBuiltin`), plugin loading requires recompilation (no runtime binary plug-in), loader extension explicitly experimental, and inter-plugin failure isolation relies on caller handling rather than supervisor.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interfaces - Factory | `Factory` interface with `Validate(*Manager, []byte)` and `New(*Manager, any) Plugin` | `v1/plugins/plugins.go:89-92` |
| Extension interfaces - Plugin | `Plugin` interface `Start(ctx) error`, `Stop(ctx)`, `Reconfigure(ctx, any)` | `v1/plugins/plugins.go:106-110` |
| Extension interfaces - Triggerable | `Triggerable{Trigger(ctx) error}` for manual trigger plugins | `v1/plugins/plugins.go:112-114` |
| Extension interfaces - LoggerPlugin | `LoggerPlugin` extends Plugin with `Logger() slog.Handler` | `v1/plugins/plugins.go:118-123` |
| Extension interfaces - Status/State | `StateNotReady/StateOK/StateErr/StateWarn`, `Status`, `StatusListener` | `v1/plugins/plugins.go:127-188` |
| Plugin loader - global factory registry | `registeredPlugins map[string]plugins.Factory` + mutex | `v1/runtime/runtime.go:69-71` |
| Plugin loader - RegisterPlugin | `func RegisterPlugin(name string, factory plugins.Factory)` idempotent global registration | `v1/runtime/runtime.go:94-98` |
| Plugin loader - Manager.Register | `Manager.Register(name, plugin)` appends to slice, signals `statusInitPlugin` | `v1/plugins/plugins.go:702-714` |
| Plugin lifecycle - Manager.Start | iterates `toStart` and calls `plugin.Start(ctx)`; recompiles external sources after start | `v1/plugins/plugins.go:870-918` |
| Plugin lifecycle - Manager.Stop | iterates `toStop`, calls `Stop(ctx)` with graceful shutdown period, drains status channel | `v1/plugins/plugins.go:926-964` |
| Plugin lifecycle - Reconfigure | Discovery validates via Factory then calls `Plugin.Reconfigure(ctx, config)` | `v1/plugins/plugins.go:105-110` comment + `v1/plugins/discovery/discovery.go:740-749` |
| Plugin lifecycle - discovery getCustomPlugins | factories validated then `manager.Plugin(name)` lookup vs `factory.New` + `manager.Register` | `v1/plugins/discovery/discovery.go:740-750` |
| Plugin lifecycle - bundle/decision_logs/status wiring | `getBundlePlugin`, `getDecisionLogsPlugin`, `getStatusPlugin` create-or-reconfigure | `v1/plugins/discovery/discovery.go:706-738` |
| Plugin isolation - status goroutine | dedicated `pluginStatusLoop` with copy-on-write `copyPluginStatus`, final status snapshot | `v1/plugins/plugins.go:1384-1422` |
| Plugin isolation - no sandbox | Manager embeds shared `Store`, `Config`, `compiler`, `services`, `Logger` with `sync.Mutex` only | `v1/plugins/plugins.go:192-245` |
| Hooks - generic Hook type | `type Hook any` with runtime type-assertion dispatch | `v1/hooks/hooks.go:32` |
| Hooks - collection | `Hooks{m map[Hook]struct{}}` with `Append`, `Each`, `Validate` | `v1/hooks/hooks.go:35-53` |
| Hooks - ConfigHook | `OnConfig(context.Context, *config.Config) (*config.Config, error)` | `v1/hooks/hooks.go:65-72` |
| Hooks - ConfigDiscoveryHook | `OnConfigDiscovery(context.Context, *config.Config)` | `v1/hooks/hooks.go:76-78` |
| Hooks - BundlePreActivateHook | `OnBundlePreActivate(ctx, bundleName, manifest) error` enables `RegisterExternalSource` during compile | `v1/hooks/hooks.go:96-98` |
| Hooks - validation | `Hooks.Validate()` enforces known hook types, rejects unknown | `v1/hooks/hooks.go:100-113` |
| Builtin extension - global registry | `var Builtins []*Builtin`, `BuiltinMap`, `RegisterBuiltin(*Builtin)` appends, not thread-safe | `v1/ast/builtins.go:14-24` |
| Builtin extension - rego package wrappers | `RegisterBuiltin1..4/Dyn` call `ast.RegisterBuiltin` + `topdown.RegisterBuiltinFunc` | `v1/rego/rego.go:746-811` |
| Builtin extension - per-query Functions | `rego.Function1..Dyn` options + `WithBuiltinFuncs(map[string]*topdown.Builtin)` | `rego/rego.go:244-266` + `v1/rego/rego.go:636` |
| Storage extension - backend builder | `type StorageBackendBuilder func(ctx, Logger, Registerer, []byte, id) (Store,error)` + `RegisterStorageBackend` | `v1/runtime/runtime.go:84-107` |
| Storage extension - Store interface | `storage.Store` used throughout Manager; `storage.Closer` for cleanup | `v1/plugins/plugins.go:193` + `v1/runtime/runtime.go:104` comment |
| Loader extension - experimental handler | `type Handler func([]byte, any) error`, `RegisterExtension(name,Handler)`, `FindExtension(ext)` with sync.Mutex | `v1/loader/extension/extension.go:14-39` |
| External rule source - compiler hook | `ExternalRuleSource{Refs()[]Ref, Init(Ref) (ExternalRuleIndex,error)}` | `v1/ast/external_source.go:13-21` |
| External rule source - index | `ExternalRuleIndex{Opts, Lookup(ctx,...)([]*Rule,ExternalRuleIndex,error)}` | `v1/ast/external_source.go:24-41` |
| External rule source - Manager integration | `Manager.RegisterExternalSource(pkgRef, source)` + `GetExternalSources()`; used in `Manager.Start` recompile | `v1/plugins/plugins.go:845-867` + `v1/plugins/plugins.go:900-915` |
| External rule source - compiler wiring | `Compiler.WithExternalSource(Ref,Source)` and `externalSources *util.HasherMap` iterated during compile | `v1/ast/compile.go:132` + `v1/ast/compile.go:1063` |
| REST auth plugin | `type HTTPAuthPlugin interface{NewClient(Config)(*http.Client,error); Prepare(*http.Request) error}` | `v1/plugins/rest/auth.go:42-51` docs + `v1/plugins/rest/rest.go:42-51` |
| Server extensibility - ExtraRoute/Middleware | `Manager.ExtraRoute(path,name,HandlerFunc)`, `ExtraMiddleware(...func(http.Handler)http.Handler)` called from plugin init | `v1/plugins/plugins.go:782-804` |
| Config - plugins map | `config.Plugins map[string]json.RawMessage` validated via factories in `getPluginSet` | `v1/plugins/discovery/discovery.go:619-636` |
| Tests - cache triggers | `TestManagerCacheTriggers`, `TestManagerNDCacheTriggers` | `v1/plugins/plugins_test.go:26-100` |
| Tests - status listeners | `TestManagerPluginStatusListener` with 2 listeners, unregister semantics | `v1/plugins/plugins_test.go:102-167` |
| Tests - stop safety | `TestUpdatePluginStatusAfterStop`, `TestPluginStatusAfterStop`, `TestRegisterAfterStop` verify non-hang after `Stop` | `v1/plugins/plugins_test.go:464-570` |
| Tests - external sources | `TestExternalSourceIntegration` 3 sub-tests: wired after start, stop cleanup, no recompilation when none | `v1/plugins/plugins_test.go:661-769` |
| Docs - plugins | Custom Plugins section: Factory/Plugin, `runtime.RegisterPlugin`, decision logger example | `docs/docs/extensions.md:190-349` |
| Docs - custom builtins | `Custom Built-in Functions in Go` + appendix with `rego.RegisterBuiltin2` / `rego.Function2` | `docs/docs/extensions.md:10-188` + `509-574` |
| Docs - storage backend | Custom Storage Backends with `RegisterStorageBackend` example | `docs/docs/extensions.md:369-402` |
| Docs - plugin config example | `type PluginFactory struct{}`, `Validate` + `New`, `runtime.RegisterPlugin` | `docs/docs/configuration.md:730-788` |
| Compatibility - init.go default plugin | `filelogger` registered in `init()` as sole default: `registeredPlugins = map[string]plugins.Factory{filelogger.Name: ...}` | `v1/runtime/runtime.go:1144-1148` |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

- **Runtime plugins** – any behavior via `plugins.Factory`/`plugins.Plugin` registered with `runtime.RegisterPlugin` (`v1/runtime/runtime.go:94-98`, `v1/plugins/plugins.go:89-110`). In-tree examples: bundle (`v1/plugins/bundle`), decision logs (`v1/plugins/logs`), status (`v1/plugins/status`), discovery (`v1/plugins/discovery`), REST clients (`v1/plugins/rest`). A custom `decision_logs.plugin` implementation is the canonical example (`docs/docs/extensions.md:222-280`).
- **Custom built-ins** – global (`ast.RegisterBuiltin` `v1/ast/builtins.go:22`, `rego.RegisterBuiltin1..Dyn` `v1/rego/rego.go:746-810`) and per-query (`rego.Function1..Dyn` `rego/rego.go:244-266`). Namespacing via `.` is encouraged.
- **Custom storage** – `storage.Store` via `runtime.RegisterStorageBackend` (`v1/runtime/runtime.go:100-107`) with optional `storage.Closer`.
- **Loader file extensions** – `loader/extension.RegisterExtension(".ext", Handler)` (`v1/loader/extension/extension.go:23`) – explicitly EXPERIMENTAL (`v1/loader/extension/extension.go:15-22`).
- **External rule sources** – `ast.ExternalRuleSource` / `ExternalRuleIndex` registered via `Manager.RegisterExternalSource` (`v1/plugins/plugins.go:848`) and wired into `ast.Compiler.WithExternalSource` (`v1/ast/compile.go:1063`).
- **REST auth** – `rest.HTTPAuthPlugin` with `NewClient`/`Prepare` (`v1/plugins/rest/rest.go:42-51`), configured via `services.*.credentials.plugin`.
- **HTTP surface** – `Manager.ExtraRoute` / `ExtraMiddleware` / `ExtraAuthorizerRoute` (`v1/plugins/plugins.go:788-815`) plus `Params.Hooks` generic hooks (`v1/hooks/hooks.go:32-98`).
- **Hooks** – `ConfigHook`, `ConfigDiscoveryHook`, `BundlePreActivateHook`, `InterQueryCacheHook`s (`v1/hooks/hooks.go:65-98`) attached via `Params.Hooks` (`v1/runtime/runtime.go:289-290`).

What is *not* directly pluggable: eval engine stages (except via `ExternalSourceOptions.VisibleRefs`/`SkippedStages` `v1/ast/external_source.go:50-71`), capabilities, or OPA’s own bundle activation without defining an activator plugin.

**2. Can plugins be loaded at runtime?**

**No – compile-time linking; runtime reconfiguration but not runtime binary loading.**

- Registration is via Go init: `runtime.RegisterPlugin` writes to a process-global map protected by `registeredPluginsMux` (`v1/runtime/runtime.go:94-98`). The discovery plugin holds that map and merges it into its factories (`v1/plugins/discovery/discovery.go:76-83`). Adding a new plugin type requires importing its package and rebuilding the `opa` binary (`docs/docs/extensions.md:314-331` shows `go build -o opa++`).
- There is no `plugin.Open` / `go plugin` dynamic `.so` loader, no WASM sidecar, and no config-only reference to an unregistered name – `getPluginSet` errors `plugin %q not registered` (`v1/plugins/discovery/discovery.go:622`).
- What *is* dynamic: once a factory is compiled in, discovery can instantiate/reconfigure/start new instances without restart. `Discovery.oneShot` → `reconfigure` → `getPluginSet` → `getCustomPlugins` validates via `Factory.Validate`, then `Start` or `Reconfigure` (`v1/plugins/discovery/discovery.go:351-437`, `740-750`). REST auth plugins similarly resolve via `Manager.AuthPlugin` at client creation time. Built-ins registered via `RegisterBuiltin*` are process-global and cannot be hot-swapped after `Manager.Start`.

**3. Are plugins isolated from each other?**

**Logically isolated, not sandboxed – shared process, cooperative.**

- **State isolation:** each plugin has independent `Status` (`StateNotReady/OK/Err/Warn`) and listeners observe a copy (`copyPluginStatus` `v1/plugins/plugins.go:1075-1088`). `pluginStatusLoop` serializes updates via a goroutine and `statusCh` (`v1/plugins/plugins.go:1384-1422`). After `Manager.Stop`, `finalPluginStatus` snapshot + closed `pluginStatusDoneCh` prevent hangs (`v1/plugins/plugins_test.go:464-523` verifies non-hang).
- **Resource isolation:** no. All plugins share `Manager.Store`, `Manager.Config`, `Manager.compiler`/`wasmResolvers`, `Manager.services`, `Logger`, router (`v1/plugins/plugins.go:192-245`). Hostile or buggy plugins can block `Manager.Start` (sequential `Start` calls `v1/plugins/plugins.go:892-896` returns on first error), mutate shared store, or panic the process. Built-in registration mutates global slices/maps unsafely after init (`v1/ast/builtins.go:18-24` comment: “not thread-safe”).
- **Failure propagation:** `Manager.Start` fails fast on any plugin `Start` error; `Reconfigure` is best-effort per plugin. No supervisor restart; `TriggerMode` (`periodic`/`manual`) only governs bundle/logs polling (`v1/plugins/plugins.go:152-164`). The fix for pre-`Start` external source recompilation (`v1/plugins/plugins.go:900-915`) shows tight coupling.

**4. Are extension points documented and stable?**

- **Documented:** Yes for primary paths. `docs/docs/extensions.md` covers custom built-ins (declaration, memoize, nondeterministic, context handling – lines `10-188`, appendix `431-574`), custom plugins (Factory/Plugin, status reporting, full decision logger example – `190-349`), and custom storage (`369-402`). `docs/docs/configuration.md:730-788` shows `PluginFactory` Validate/New skeleton and `runtime.RegisterPlugin`. `plugins/plugins.go:26-87` doc-comments describe two-step Validate→New and start/reconfigure lifecycle.
- **Stability:** Moderate–strong via `v1/` import path. Core plugin interfaces live in `v1/plugins/plugins.go` and are re-exported by `plugins/plugins.go` (`26-73`) – a classic compat shim. Hooks are versioned `v1/hooks`. Storage interface stability is Go-interface-based. Weak spots: loader `extension` is explicitly EXPERIMENTAL with disclaimer “may go away” (`v1/loader/extension/extension.go:15-22` and `CHANGELOG.md:3080`), built-in registration is global and documented as “only during initialization” (`v1/ast/builtins.go:18-21`), and `ExternalRuleSource` is new (2026 copyright) with less external docs.
- **Tests as spec:** Plugin manager status, cache triggers, stop safety, and external sources are covered (`v1/plugins/plugins_test.go:26-769`), but there is no fuzz or chaos test for plugin panics/isolation, and server `ExtraRoute` has no contract test visible in this study.

## Architectural Decisions

- **Global factory registry + discovery-driven activation** (`v1/runtime/runtime.go:69-98`, `v1/plugins/discovery/discovery.go:606-750`). Keeps core `runtime` free of plugin imports, delegates validation/instantiation to discovery’s `getPluginSet`. Tradeoff: simple static linking but prevents untrusted dynamic plugins.
- **Manager as God object** (`v1/plugins/plugins.go:192-245`). Exposes store, compiler, services, keys, logger, router, tracer, hooks through one type so plugins need only one dependency. Tradeoff: ergonomics vs. least privilege; plugin can touch anything.
- **Two-phase plugin construction (`Validate` then `New`)** (`v1/plugins/plugins.go:42-50`). Allows discovery to reject bad config without leaking partially-constructed plugin. Mirrors bundle plugin pattern.
- **Status as observable** (`State*` + `UpdatePluginStatus` + `RegisterPluginStatusListener` + `pluginStatusLoop` – `v1/plugins/plugins.go:129-188`, `1384-1422`). Gives server `ready` semantics (`runtime.waitPluginsReady` `v1/runtime/runtime.go:1057-1075`) and external health checks.
- **Hooks as `any` with interface assertion** (`v1/hooks/hooks.go:16-32`, `v1/plugins/plugins.go:573-581`). Allows additive hook types without breaking `Hook` interface; `Validate` enumerates allowed types (`v1/hooks/hooks.go:100-113`).
- **ExternalRuleSource as lazy rule-tree delegate** (`v1/ast/external_source.go:13-41`, `v1/ast/compile.go:1063`, `v1/plugins/plugins.go:848-915`). Lets plugins supply rules without writing to store, with per-evaluation `ExternalRuleIndex` state and optional `Close`. Tradeoff: intrusive compiler coupling for power.

## Notable Patterns

- **Compat shim layer:** top-level `plugins/` and `hooks/` re-export `v1/...` via type aliases (`plugins/plugins.go:73`, `hooks/hooks.go:27`). Keeps old import paths working while encouraging `v1/` migration.
- **Functional options for Manager:** `Info`, `InitBundles`, `Logger`, `WithRouter`, `WithHooks`, etc. (`v1/plugins/plugins.go:367-535`) – consistent with `Params` opts.
- **Copy-on-read config:** `Manager.GetConfig().Clone()` (`v1/plugins/plugins.go:693-698`) and `Reconfigure`’s deep copy + bootstrap-label preservation (`v1/plugins/plugins.go:980-1032`) prevent data races during discovery updates.
- **Graceful shutdown orchestration:** `Manager.Stop` honors `GracefulShutdownPeriod`, calls plugin `Stop(ctx)` with timeout context, closes store if `storage.Closer`, drains status goroutine (`v1/plugins/plugins.go:926-964`, `v1/runtime/runtime.go:1016-1055`).

## Tradeoffs

- **Extensibility vs. safety:** powerful plugin hook surface (store, compiler, HTTP routes) enables Envoy/Kubernetes integrations but gives plugins full blast radius; no capability object limits what a plugin can do.
- **Build-time vs. runtime extensibility:** static linking ensures type safety and Go toolchain optimizations, but operators must fork/build `opa++` for custom logic – undesirable for SaaS/multi-tenant deployments.
- **Global built-in map:** zero-cost dispatch for eval hot path, but `RegisterBuiltin` is not safe after start and creates ordering hazards; per-query `WithBuiltinFuncs` mitigation is newer and not universal.
- **Discovery centralization:** single `discovery` plugin owns all lifecycle transitions (`getPluginSet` validates bundles, decision logs, status, plus custom factories). Simplifies ordering but makes discovery a single point of failure – `reconfigure` error aborts entire update (`v1/plugins/discovery/discovery.go:368-371`).
- **EXPERIMENTAL loader extension:** enables YAML/JSON-extension handlers without forking, but lack of stability promise discourages production adoption.

## Failure Modes / Edge Cases

- **Unregistered plugin blocks bundle activation.** `getPluginSet` returns `plugin %q not registered` (`v1/plugins/discovery/discovery.go:622`) which aborts `processBundle` → `MaxActivationRetry` loop in `loadAndActivateBundleFromDisk` (`v1/plugins/discovery/discovery.go:269-301`) but not in normal download path where error surfaces as discovery error status.
- **Start ordering hazard.** Plugins that call `RegisterExternalSource` must do so before `Manager.Start` recompiles (`v1/plugins/plugins.go:900-908`); late registration after start has no recompile path except next bundle activation.
- **Stop-after-stop hang prevention is tested** (`v1/plugins/plugins_test.go:524-570` non-hang assertions) by closing `pluginStatusDoneCh` and returning snapshot; missing this would deadlock `UpdatePluginStatus`.
- **Compiler cache poisoning.** `SetCompilerOnContext`/`SetWasmResolversOnContext` optimization (`v1/plugins/plugins.go:294-321`) skips recompilation if bundle plugin sets context; mis-set compiler bypasses validation and can install stale policy.
- **Auth plugin panic.** `Manager.AuthPlugin` does unchecked `.(rest.HTTPAuthPlugin)` type assertion (`v1/plugins/plugins.go:739-749`) – non-auth plugins registered under a service name cause panic.
- **ExtraRoute double-register panic.** `ExtraRoute` panics on duplicate path (`v1/plugins/plugins.go:789-791`); two plugins racing to register same path crash the server.
- **Discovery cannot update its own service/key** – explicit error “updates to the discovery service are not allowed” (`v1/plugins/discovery/discovery.go:502`) and similar for keys (`v1/plugins/discovery/discovery.go:516`) – operator must restart to rotate discovery credentials.
- **External source overlay visibility misconfiguration.** `ExternalSourceOptions.VisibleRefs=nil` isolates, but `[]Ref{MustParseRef("data")}` exposes entire rule tree – overly permissive config can cause compile-time leakage or performance blow-up.

## Future Considerations

- Add capability-scoped Manager view (read-only store txn, no `ExtraRoute` after server start) to limit plugin blast radius; enforce via interface rather than convention.
- Replace global `ast.Builtins`/`BuiltinMap` mutation with per-`Compiler` registration or `Capabilities` injection; deprecate `RegisterBuiltin` after start path.
- Stabilize `loader/extension` or remove – either promote with semver guarantee or document replacement (external source + bundle plugin).
- Provide dynamic backend for WASM ABI without recompilation (e.g., gRPC plugin host) if multi-tenant isolation becomes a goal; current model is unsuitable for untrusted third parties.
- Harden discovery with per-plugin validation isolation so one bad custom plugin config does not block bundle/status/logs updates.
- Document hook ordering guarantees or provide ordered `Hooks` variant; currently map-range execution is unspecified (`v1/hooks/hooks.go:23-28` note).
- Add supervisor for plugin panics (recover in `Manager.Start` loop and mark `StateErr`) instead of crashing runtime.

## Questions / Gaps

- No evidence of plugin dependency ordering or topological sort – discovery starts bundles/status/logs before custom plugins (`v1/plugins/discovery/discovery.go:669-702` vs `740-750`); what if a custom plugin provides a service needed by bundle downloading?
- No observability for extension points beyond status: no metrics for custom built-in invocation count/latency; `LoggerPlugin` is the only performance-adjacent hook.
- No documentation or test for `storage.Store` custom backend lifecycle under concurrent `Manager.Reconfigure` – search of `v1/plugins/plugins_test.go` and `v1/runtime/runtime_test.go` found no custom store test except `TestManager*` inmem variants.
- Source isolation prevented inspecting sibling `sources/` for comparative plugin models, so cannot assess whether OPA is more/less extensible than agent harnesses in this study cohort – flagged as out-of-scope per hard rule 1.

---

Generated by `Dimension 21.01: Plugin and Extension Points` against `opa`.
