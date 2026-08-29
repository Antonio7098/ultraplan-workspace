# Source Analysis: temporal

## Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go 1.26 / Uber Fx DI, gRPC, Protobuf |
| Analyzed | 2026-08-28 |

## Summary

Temporal is not a user-facing agent harness with "tools" but a workflow orchestration server. Its extension model is **compile-time dependency injection** via `temporal.ServerOption` functional options wired through Uber Fx, not runtime plugin loading. Third parties extend the server by implementing Go interfaces (auth, persistence, archival, search-attribute mapping, dynamic config, TLS, metrics, gRPC interceptors) and injecting factories at `temporal.NewServer()` startup. Config-driven extension exists for SQL persistence (`pluginName`), archival `customStores`, and TLS/authorization YAML. There is no dynamic discovery, no sandboxing, and no hot reload — every extension runs in-process with full privilege and lifecycle is tied to Fx Start/Stop. Archival and search-attribute mapping are mature and well-documented; other points (custom datastore, metrics handler) are marked experimental and require source-available injection.

## Rating

**6 / 10 — Present but inconsistent, weakly documented or fragile**

Rationale: Explicit interfaces exist for ~10 extension families and are exercised in tests (`tests/testcore/onebox.go`, `common/archiver/provider/provider_test.go`), with operational guards (fail-fast validation, caching, fatal logs). However there is no runtime/dynamic loading, no isolation, no versioned/stable plugin API contract (several options marked `experimental`/`NOTE: may be changed`), and documentation is scattered (only `common/archiver/README.md` is comprehensive). A third party **can** add a new store/archiver/auth implementation without forking core, but only by compiling a custom binary — not by dropping a plugin at runtime.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Auth - Authorizer | `type Authorizer interface { Authorize(ctx context.Context, caller *Claims, target *CallTarget) (Result, error) }` + `GetAuthorizerFromConfig` switches on `config.Authorization.Authorizer` (`""`/`"default"`) and caller can inject via `WithAuthorizer` | `common/authorization/authorizer.go:53-56`, `common/authorization/authorizer.go:64-73`, `temporal/server_option.go:96-101` |
| Auth - ClaimMapper | `type ClaimMapper interface { GetClaims(authInfo *AuthInfo) (*Claims, error) }` + `WithClaimMapper(func(cfg *config.Config) ClaimMapper)` | `common/authorization/claim_mapper.go:28-31`, `temporal/server_option.go:110-114` |
| Auth - Audience | `type AudienceMapper` / `JWTAudienceMapper` + `WithAudienceGetter` | `common/authorization/audience_mapper.go:11-18`, `temporal/server_option.go:118-122` |
| Auth - TokenProvider | `WithTokenProvider(auth.TokenProvider)` validated against `global.authorization.remoteClusterAuth.require` and `remoteClusters` TLS | `temporal/server_option.go:202-207`, `temporal/fx.go:274-283` |
| TLS | `type TLSConfigProvider interface { GetInternodeServerConfig(...) ... GetRemoteClusterClientConfig(...) }` + `WithTLSConfigFactory` + `NewTLSConfigProviderFromConfig` defaulting to `localStoreCertProvider` | `common/rpc/encryption/tls_factory.go:17-24`, `temporal/server_option.go:103-108`, `common/rpc/encryption/tls_factory.go:60-76` |
| Persistence - custom datastore | `type AbstractDataStoreFactory interface { NewFactory(cfg CustomDatastoreConfig, ...) persistence.DataStoreFactory }` + `WithCustomDataStoreFactory` + `DataStoreFactoryProvider` branching on `Cassandra`/`SQL`/`CustomDataStoreConfig` | `common/persistence/client/abstract_data_store_factory.go:16-25`, `temporal/server_option.go:145-151`, `common/persistence/client/fx.go:187-209` |
| Persistence - visibility | `type VisibilityStoreFactory interface { NewVisibilityStore(cfg CustomDatastoreConfig, ...) (store.VisibilityStore, error) }` + `WithCustomVisibilityStoreFactory` + `newVisibilityStoreFromDataStoreConfig` fallback to custom factory | `common/persistence/visibility/factory.go:20-31`, `temporal/server_option.go:153-157`, `common/persistence/visibility/factory.go:284-298` |
| Archival - archivers | `type HistoryArchiver interface { Archive/Get/ValidateURI }` / `VisibilityArchiver` | `common/archiver/interface.go:44-58`, `common/archiver/interface.go:76-94` |
| Archival - factories | `CustomHistoryArchiverFactory` / `CustomVisibilityArchiverFactory` with `NewCustom*Params{Scheme,Configs,Logger}` + func adapters + provider caching + `ErrUnknownScheme` delegation | `common/archiver/provider/provider.go:54-63`, `common/archiver/provider/provider.go:65-67`, `common/archiver/provider/provider.go:123-184`, `common/archiver/provider/provider.go:186-246` |
| Archival - config | `CustomStores map[string]map[string]any yaml:"customStores"` for both history/visibility; `IndexName`+`Options` for custom datastores | `common/config/config.go:521`, `common/config/config.go:541`, `common/config/config.go:468-479` |
| Archival - injection | `WithCustomHistoryArchiverFactory` / `WithCustomVisibilityArchiverFactory` wired via `ServerOptionsProvider` → `ServiceProviderParamsCommon` → each service graph | `temporal/server_option.go:159-173`, `temporal/fx.go:110-111`, `temporal/fx.go:393-446` |
| Search attributes | `type Mapper interface { GetAlias/GetFieldName }` + `WithSearchAttributesMapper` | `common/searchattribute/mapper.go:20-23`, `temporal/server_option.go:183-188` |
| Dynamic config | `type Client interface { GetValue(key Key) []ConstrainedValue }` + optional `NotifyingClient` + `WithDynamicConfigClient` defaulting to `FileBasedClient` or `NoopClient` | `common/dynamicconfig/client.go:12-32`, `temporal/server_option.go:138-143`, `temporal/fx.go:207-221` |
| Metrics | `WithCustomMetricsHandler(metrics.Handler)` + `MetricsHandlerFromConfig` fallback | `temporal/server_option.go:209-214`, `temporal/fx.go:199-205` |
| gRPC interceptors | `WithChainedFrontendGrpcInterceptors(...grpc.UnaryServerInterceptor)` appended after internal interceptors; service-level `AdditionalInterceptors []grpc.UnaryServerInterceptor optional:"true"` | `temporal/server_option.go:190-200`, `service/fx.go:51-52`, `service/fx.go:144-179`, `temporal/fx.go:538-560` |
| RPC client factory | `WithClientFactoryProvider(client.FactoryProvider)` | `temporal/server_option.go:175-181`, `temporal/fx.go:193-196` |
| SQL plugin (internal) | `type Plugin interface { CreateDB(...)(GenericDB, error); GetVisibilityQueryConverter() }` — compile-time registration via `sqlplugin` packages (`mysql`, `postgresql`, `sqlite`) | `common/persistence/sql/sqlplugin/interfaces.go:32-36`, `common/persistence/sql/sqlplugin/mysql/plugin.go`, `common/persistence/sql/sqlplugin/postgresql/plugin.go` |
| Persistence resolver | `WithPersistenceServiceResolver(resolver.ServiceResolver)` | `temporal/server_option.go:124-129`, `temporal/fx.go:397-405` |
| Fx lifecycle | `ServerOptionsProvider` validates `remoteClusterAuth`+`TokenProvider`+TLS co-requirement; all extensions supplied via `fx.Supply`/`fx.Provide` with `fx.StartStopHook` lifecycle | `temporal/fx.go:274-283`, `temporal/fx.go:393-471`, `common/persistence/client/fx.go:66-85` |
| Tests exercising extension | `onebox.go` injects `WithAuthorizer(impl)`, `WithClaimMapper(...)`, `WithCustomDataStoreFactory(...)`, `WithSearchAttributesMapper(nil)`; `archivaltest` uses `WithCustomArchivers`; `lite_server.go` uses `WithAuthorizer`/`WithClaimMapper` | `tests/testcore/onebox.go:138-150`, `tests/archival_test.go:160`, `temporaltest/internal/lite_server.go:273-274`, `temporal/server_test.go:89-99` |
| Docs | `common/archiver/README.md` documents Option 2 custom factory end-to-end (interfaces, params, `customStores` YAML); `common/config/config.go:641-660` documents `Authorizer`/`ClaimMapper` string values; server_option.go godoc marks several factories `experimental` | `common/archiver/README.md:8-194`, `common/config/config.go:641-660`, `temporal/server_option.go:145-173` |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

Ten families, all requiring Go code at compile time:

- **AuthZ/AuthN**: `Authorizer`, `ClaimMapper`, `JWTAudienceMapper` (`authorization.Authorizer:64`, `authorization.ClaimMapper:28`, `authorization.JWTAudienceMapper`), plus `auth.TokenProvider` for outbound cross-cluster auth (`temporal/server_option.go:202`). Config YAML selects built-ins (`authorizer: default`/empty, `claimMapper: default`) (`common/config/config.go:648-651`, `temporal/fx.go:560-570`), but any implementation can be injected via `WithAuthorizer`/`WithClaimMapper`/`WithAudienceGetter`/`WithTokenProvider`.
- **Persistence**: `AbstractDataStoreFactory` for main store (`common/persistence/client/abstract_data_store_factory.go:16`) and `VisibilityStoreFactory` (`common/persistence/visibility/factory.go:20`). Selection via `persistence.datastores.<name>.customDatastore` (`common/config/config.go:468-479`).
- **Archival**: `HistoryArchiver`/`VisibilityArchiver` (`common/archiver/interface.go:44-94`) via `CustomHistoryArchiverFactory`/`CustomVisibilityArchiverFactory` (`common/archiver/provider/provider.go:54-63`). YAML `archival.*.provider.customStores.<scheme>` passes opaque `map[string]any` to factories.
- **Search attributes**: `searchattribute.Mapper` alias↔field mapping per namespace (`common/searchattribute/mapper.go:20`), via `WithSearchAttributesMapper`.
- **Dynamic config**: `dynamicconfig.Client` (`GetValue` + optional `NotifyingClient.Subscribe`) (`common/dynamicconfig/client.go:12`), via `WithDynamicConfigClient` (defaults to file-based or noop).
- **TLS**: `encryption.TLSConfigProvider` (`common/rpc/encryption/tls_factory.go:17`), via `WithTLSConfigFactory` (defaults to local-store file provider).
- **Metrics**: `metrics.Handler` via `WithCustomMetricsHandler`.
- **RPC**: `client.FactoryProvider` (history/matching/frontend client creation) via `WithClientFactoryProvider`; `resolver.ServiceResolver` via `WithPersistenceServiceResolver`; ElasticSearch HTTP client via `WithElasticsearchHttpClient`.
- **gRPC interceptors**: Frontend unary chain via `WithChainedFrontendGrpcInterceptors`; internal services via `AdditionalInterceptors` Fx optional param (`service/fx.go:51`).
- **SQL dialect plugins**: Internal `sqlplugin.Plugin` (`common/persistence/sql/sqlplugin/interfaces.go:32`) for MySQL/PostgreSQL/SQLite dialects — not a third-party plugin point per se but a separate compilation unit extending SQL support.

No extension point exists for workflow execution semantics, task queue matching, or history replication without forking core services.

**2. Can plugins be loaded at runtime?**

**No.** Evidence:

- All injection points are `temporal.ServerOption` functions applied before `fx.New` in `ServerOptionsProvider` (`temporal/fx.go:170-180`) and converted to Fx supplies (`temporal/fx.go:299-317`, `393-471`). There is no filesystem scan, registry, `plugin.Open`, WASM, or gRPC plugin discovery.
- Custom datastore requires `config.CustomDatastoreConfig.Name/IndexName/Options` to be present at startup; `DataStoreFactoryProvider` (`common/persistence/client/fx.go:187-209`) branches statically and `logger.Fatal` if missing.
- Archiver `customStores` (`common/config/config.go:521/541`) are read from YAML at startup and `GetHistoryArchiver` caches the first successful creation forever (`common/archiver/provider/provider.go:123-184`).
- Dynamic config is the only hot-reloadable piece, and even there the `Client` is fixed at startup; only its returned values change via `Subscribe` (`common/dynamicconfig/client.go:36-41`).
- Tests construct servers programmatically with `With*` options (`tests/testcore/onebox.go:135-181`), never by loading a `.so`.
- `WithCustomDataStoreFactory` godoc says `NOTE: this option is experimental and may be changed or removed` (`temporal/server_option.go:145-146`) — indicating an unstable compile-time hook, not a stable runtime plugin API.

To add a new extension, rebuild the `temporal` binary with your factory implementations imported.

**3. Are plugins isolated from each other?**

**No isolation.** Evidence:

- All factories/handlers run in the server process. `AbstractDataStoreFactory.NewFactory` receives the server logger/metrics/resolver and returns a `persistence.DataStoreFactory` that is wrapped directly (`common/persistence/client/fx.go:205`). Faulty factory can `logger.Fatal` and crash the entire server (`common/persistence/client/fx.go:206-208`, `common/persistence/visibility/factory.go:285-287`).
- `archiverProvider` shares a single `sync.RWMutex` cache for all schemes but no per-plugin classloader, resource quota, or error boundary beyond `ErrUnknownScheme` delegation (`common/archiver/provider/provider.go:69-85`, `136-146`).
- Interceptors are chained in-process: `grpc.ChainUnaryInterceptor(getUnaryInterceptors(...))` (`service/fx.go:154`). A panicking interceptor kills the gRPC handler; there is no per-interceptor recovery beyond the `ServiceErrorInterceptor`.
- Metrics handler, TLS provider, claim mapper, authorizer are global singletons supplied via `fx.Supply` (`temporal/fx.go:106-128`). One misbehaving implementation blocks all requests for that service.
- The only fault-isolation mechanism is `FaultInjection` (`common/config/config.go:287-357`) which is a test-oriented wrapper around persistence, not a plugin sandbox.

**4. Are extension points documented and stable?**

Partially. Archival is exemplary; others are uneven and explicitly unstable.

- **Well-documented**: `common/archiver/README.md:86-194` provides end-to-end guide for custom archivers (interfaces, params structs, factory func adapters, YAML snippets, FAQ on overriding built-ins). `common/searchattribute/mapper.go:18-19` has inline godoc on `WithSearchAttributesMapper`. Auth interfaces carry `// @@@SNIPSTART` doc tags (`common/authorization/authorizer.go:22`, `claim_mapper.go:28`) and config comments (`common/config/config.go:641-660`).
- **Weakly documented**: Persistence custom stores lack README; only `abstract_data_store_factory.go:14-15` comment says “can be used to implement custom datastore support outside of Temporal core.” `WithCustomDataStoreFactory`/`WithCustomVisibilityStoreFactory`/`WithCustom*ArchiverFactory` are annotated `experimental and may be changed or removed` (`temporal/server_option.go:146`, `159`, `168`). No versioning or SemVer contract for `Custom*Params`.
- **Tests as docs**: `tests/testcore/onebox.go:146-149` shows idiomatic auth injection using `temporalImpl`’s own `GetClaims`/`Authorize`; `tests/archival_test.go:160` shows `WithCustomArchivers`. But no example for `AbstractDataStoreFactory` outside `tests/testcore/test_cluster.go` fixtures.
- **Stability**: Godoc `With*` options are the de-facto stability surface (since `temporal.NewServer` is public). Archiver provider guarantees `ErrUnknownScheme` fallback precedence (`common/archiver/provider/provider.go:52-62`), but history indicates `customStores` replaced earlier per-provider config. `DefaultServices` and `Config.Validate()` enforce static membership (`temporal/fx.go:265-272`), not a pluggable service registry.

## Architectural Decisions

| Decision | Evidence | Tradeoff |
|----------|----------|----------|
| Compile-time functional options + Uber Fx over runtime plugin registry | `temporal/server_option.go:25-33` (interface + apply func), `temporal/fx.go:94-128` (ServerOptionsProvider struct with all extensions), `temporal/fx.go:393-471` (GetCommonServiceOptions converts to fx.Supply) | Type-safe, fast, no serialization; but requires rebuild and no hot reload. Favored for an infra server where correctness > extensibility agility. |
| Factory-per-scheme with `ErrUnknownScheme` delegation | `common/archiver/provider/provider.go:54-63`, `131-170` (custom first, then switch on `filestore`/`gcloud`/`s3store`) | Allows overriding built-ins without fork; opaque `Configs map[string]any` avoids schema coupling but sacrifices validation. |
| Opaque `map[string]any` for custom config | `common/config/config.go:478`, `common/archiver/provider/provider.go:40,48` (NewCustom*Params.Configs) | Maximum flexibility for third parties; but typos fail silently at factory parse time, no JSON schema. |
| Fx optional params for interceptors | `service/fx.go:51-52` (`optional:"true"`), `temporal/fx.go:538-560` (Decorate internal-frontend claim mapper) | Zero-cost when unused, avoids nil checks everywhere; but discovery is implicit — grep required to find injection points. |
| Fatally crash on invalid custom config vs return error | `common/persistence/client/fx.go:207` (`logger.Fatal`), `common/persistence/visibility/factory.go:286` | Fail-fast prevents split-brain; but violates graceful degradation for multi-tenant hosting. |
| Keep TLS/auth/metric as singletons | `temporal/fx.go:115-128`, `temporal/server_option.go:47-62` (single `authorizer`, `claimMapper` fields) | Simple wiring; precludes per-namespace/per-cluster plugin selection. |

## Notable Patterns

- **Functional options as plugin API**: `type ServerOption interface { apply(*serverOptions) }` + `applyFunc` (`temporal/server_option.go:25-33`) — standard Go pattern reused as extension surface. All ~12 `With*` constructors (`temporal/server_option.go:35-214`) follow same shape.
- **Factory function adapters**: `CustomHistoryArchiverFactoryFunc func(...)` implements `CustomHistoryArchiverFactory` via method (`common/archiver/provider/provider.go:65-92`) — lets users pass lambdas (`common/archiver/README.md:134-146`).
- **Double-checked locking cache**: `archiverProvider.GetHistoryArchiver` RLock → check → RUnlock → maybe create → Lock → re-check → insert (`common/archiver/provider/provider.go:123-184`) — safe lazy init without global lock.
- **Fx module composition**: `TopLevelModule` (`temporal/fx.go:131-154`) composes `dynamicconfig.Module`, `pprof.Module`, `chasm.Module`, etc., then `GetCommonServiceOptions` re-supplies dependencies per service graph (`temporal/fx.go:393-471`) — a workaround for cross-graph dependency propagation.
- **Test harness implements interfaces**: `tests/testcore/onebox.go:377-407` `temporalImpl.GetClaims`/`Authorize` flips based on `onGetClaims`/`onAuthorize` callbacks — shows how minimal an auth plugin can be.
- **Custom search attribute names via MapperProvider indirection**: `mapperProviderImpl` first checks `customMapper != nil` else falls back to `namespaceRegistry.GetCustomSearchAttributesMapper` with `backCompMapper` (`common/searchattribute/mapper.go:90-124`) — preserves legacy `Keyword01` behavior while allowing alias overrides.

## Tradeoffs

- **Compile-time safety vs runtime agility**: Strong typing catches wiring errors at `fx.New` (e.g., `temporal/fx.go:274-283` validates `TokenProvider` + TLS co-requirement). Cost: every new archiver/store requires a vendor fork or at least a `main.go` wrapper that imports both server and plugin.
- **In-process performance vs blast radius**: Direct interface calls avoid IPC overhead (archiver `Archive` called synchronously from history service). A buggy archiver blocks the archival queue processor for that shard.
- **Opaque config vs schema**: `CustomDatastoreConfig.Options map[string]any` (`common/config/config.go:478`) lets a plugin define arbitrary YAML, but server cannot validate it without invoking the factory. Compare to typed `Cassandra`/`SQL` structs (`common/config/config.go:359-465`) which are validated.
- **Global singletons vs per-tenant isolation**: One `Authorizer`/`ClaimMapper`/`MetricsHandler` for whole process simplifies authz reasoning (single `CallTarget` check in `authorization.Interceptor:common/authorization/interceptor.go`) but cannot offer namespace-isolated plugins.
- **Experimental tag protects velocity**: Marking `WithCustomDataStoreFactory` experimental (`temporal/server_option.go:146`) lets team evolve `AbstractDataStoreFactory` (recent addition of `serialization.Serializer` param) without SemVer break — at expense of third-party churn.

## Failure Modes / Edge Cases

- **Factory returns nil without error**: `archiverProvider.GetHistoryArchiver` treats `historyArchiver == nil` as “not handled” and falls through to built-in switch (`common/archiver/provider/provider.go:148-170`). If custom factory mistakenly returns `(nil, nil)` for its own scheme, server silently falls to `ErrUnknownScheme`/`ErrArchiverConfigNotFound` and archival fails without indicating the custom factory was consulted.
- **Config typo in `customStores.<scheme>`**: `customConfigs = provider.CustomStores[scheme]` may be `nil` (`common/archiver/provider/provider.go:133-135`); factory receives `nil` map and must handle it. No early validation in `Config.Validate()` for custom stores.
- **Panic in gRPC interceptor**: `WithChainedFrontendGrpcInterceptors` appends after internal interceptors (`temporal/fx.go:559`); `service/fx.go:160-179` builds chain with `ServiceErrorInterceptor` first. A panicking custom interceptor bypasses `RecoveryInterceptor` placement assumptions — depends on order.
- **Fault injection double-wrap**: `DataStoreFactoryProvider` layers `faultinjection.NewFaultInjectionDatastoreFactory` then `telemetry.NewTelemetryDataStoreFactory` (`common/persistence/client/fx.go:210-217`). Custom factories not aware of this order; if they internally wrap with own telemetry, double counting occurs.
- **TLS + TokenProvider mismatch**: Server boot fails with explicit error if `remoteClusterAuth.require && tokenProvider == nil` or `tokenProvider != nil && no remote TLS` (`temporal/fx.go:274-283`). Mis-ordered `With*` calls cannot fix config — must edit YAML or add `WithTLSConfigFactory`.
- **Visibility dual-manager inconsistency**: `NewManager` creates both primary and secondary managers from same `customVisibilityStoreFactory` (`common/persistence/visibility/factory.go:59-127`). If factory is stateful per index, sharing instance across both stores can cause data corruption; lifecycle not documented.
- **No per-plugin health signal**: `HealthSignalAggregator` (`common/persistence/client/fx.go:163-185`) aggregates persistence latency; a slow custom store degrades global metrics without attribution per plugin.

## Future Considerations

- Introduce a stable, versioned plugin SDK (separate Go module `go.temporal.io/server-plugin`) with SemVer on `AbstractDataStoreFactory` / `VisibilityStoreFactory` / `Custom*ArchiverFactory`, removing `experimental` guards and adding `go generate` mock parity.
- Add dynamic / out-of-process extension via gRPC sidecar (e.g., `Authorizer` as gRPC service) to enable policy updates without binary rebuild — leverage existing `GetClaims`/`Authorize` call sites (`common/authorization/interceptor.go:28-60`) as interception point.
- Add per-plugin isolation: run custom archivers/visibility stores behind circuit breaker (`common/circuitbreaker`) and bulkhead; surface per-plugin metrics tags (`metrics.VisibilityPluginNameTag` already exists: `common/persistence/visibility/factory.go:229`) consistently for auth/TLS.
- Validate `customStores` and `customDatastore.options` against JSON schema at `Config.Validate()` time via factory-provided `ValidateConfig(map[string]any) error` hook, failing fast with field-level errors instead of runtime `nil` map panics.
- Scope interceptors per-service/ per-namespace: expose stream interceptors for history/matching, not just frontend unary (`temporal/server_option.go:190`), and document ordering guarantees.
- Document all `With*` options in a single `docs/plugin-extensibility.md` (currently only archival has `common/archiver/README.md`), including matrix of stable vs experimental and minimal example `main.go` wrapper.

## Questions / Gaps

- No evidence of a plugin manifest, version negotiation, or capability advertisement — `GetAuthorizerFromConfig` only handles `""`/`"default"` (`common/authorization/authorizer.go:64-73`). How would a third party register a named built-in without code change? No registry found.
- No evidence of extension point for workflow/activity execution (e.g., custom workflow interceptors inside history) — Chasm components are internal (`chasm/`). Could an external CHASM library be injected? `chasm.Module` is Fx-provided but not exposed via `ServerOption`. Searched `chasm` for `FactoryProvider` — no hook found.
- Metrics `metrics.Handler` is replaceable, but tracing `[]otelsdktrace.SpanExporter` is only replaceable via `fx.Replace` in `temporal/fx.go:940-984`, not via `ServerOption` — inconsistency in observability extensibility.
- No integration test that verifies multiple custom archiver factories co-exist (e.g., `myscheme` + `filestore` override) — `common/archiver/provider/provider_test.go` only tests built-ins.
- `sqlplugin.Plugin` extension requires adding package under `common/persistence/sql/sqlplugin/*` and rebuilding — no `WithSqlPlugin` option — so adding a new SQL dialect still requires core modification. Is this intentional vs abstract datastore path?
- Does `WithPersistenceFactoryProvider` (`tests/testcore/onebox.go:167`) intended as public extension? It is not exported in `temporal/server_option.go` (only `WithCustomDataStoreFactory` is). Search found only test usage via `WithAdditionalStreamInterceptors` etc., not a stable hook.
- No UI/CLI plugin point — `tools/tdbg` and `temporal` CLI are not pluggable.

---

Generated by `21.01-plugin-and-extension-points` against `temporal`.
