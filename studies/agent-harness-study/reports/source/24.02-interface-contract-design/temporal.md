# Source Analysis: temporal

## Interface Contract Design

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go (gRPC/protobuf services, Cassandra/SQL persistence, Elasticsearch visibility, fx DI) |
| Analyzed | 2026-08-22 |

## Summary

Temporal server defines its contracts at four layers, each with a different enforcement mechanism:

1. **Wire/schema layer** — protobuf definitions under `proto/internal/` generate the public and internal service APIs (`api/`), including typed error-detail payloads (`api/errordetails/v1/message.pb.go`).
2. **Service-client layer** — consumer-facing Go interfaces are the generated gRPC client interfaces (e.g., `workflowservice.WorkflowServiceClient`), wrapped by generated decorator stacks (timeout, metric, retryable) produced by a repo-owned code generator (`client/frontend/client.go:1-2`, `cmd/tools/genrpcwrappers/main.go`).
3. **Persistence layer** — a deliberate two-layer contract: low-level `*Store` interfaces implemented per datastore (Cassandra/SQL) and high-level `*Manager` interfaces that add serialization, then get wrapped in metric/rate-limit/retry decorators (`common/persistence/persistence_interface.go:23-28`, `common/persistence/client/factory.go:109-120`).
4. **Component framework (CHASM)** — a generic Go contract model where components, tasks, and libraries register into a validating registry; the `chasm.Engine` interface has two independent implementations (production `service/history/chasm_engine.go` and in-memory `chasm/chasmtest/test_engine.go`).

The distinguishing strength is that contracts encode *behavior*, not just signatures: error taxonomy carries retryability and redirect semantics (`common/persistence/error_type.go:5-22`, `common/util.go:265-295`), request structs carry lifecycle callbacks (`common/persistence/persistence_interface.go:213-217`), and registration-time validation rejects conflicting contracts before startup completes (`chasm/registry.go:247-277`). Substitutability is proven by conformance suites that run identical test suites against Cassandra, MySQL, PostgreSQL, and SQLite backends (`common/persistence/persistence-tests/persistence_test_base.go:174-183`).

## Rating

**8 / 10.** Contracts are explicit, layered, behavior-rich, validated at multiple points in time (compile, registration, schema), and exercised by cross-backend conformance suites plus a second full engine implementation used in tests. The score is held below 9 because several contracts are explicitly incomplete (`Speculative` transitions unimplemented in both engines, `chasm/engine.go:144-150`; `TransitionOption`s ignored by several engine methods, `chasm/engine.go:290-292`; `Workflow.LifecycleState` is a stub, `chasm/lib/workflow/workflow.go:53-61`), validation style is inconsistent (panics in options vs errors in registry, `chasm/registrable_component.go:96-101`), and some interfaces are very large with implicit invariants (`ExecutionStore`, `common/persistence/persistence_interface.go:116-167`).

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Central interface: persistence store contracts | `DataStoreFactory` vends `ShardStore`, `TaskStore`, `MetadataStore`, `ExecutionStore`, `Queue`, `QueueV2`, `ClusterMetadataStore`, `NexusEndpointStore`; header comment states intent to share serialization logic across SQL/Cassandra | common/persistence/persistence_interface.go:23-51 |
| Central interface: persistence manager contracts | `ExecutionManager` (30+ methods), `TaskManager`, `MetadataManager`, `ClusterMetadataManager`, `NexusEndpointManager`, `HistoryTaskQueueManager` | common/persistence/data_interfaces.go:1106-1270 |
| Lifecycle contract | `Closeable` interface with TODO "allow this method to return errors" | common/persistence/data_interfaces.go:1099-1103 |
| Behavior-bearing request contract | `InternalGetOrCreateShardRequest` carries a `CreateShardInfo` callback and a `LifecycleContext` documented as "cancelled when shard is unloaded" | common/persistence/persistence_interface.go:210-217 |
| Two-layer rationale | "Persistence interface is a lower layer of dataInterface… let different persistence implementations share common logic (serialization)" | common/persistence/persistence_interface.go:24-28 |
| Manager implements contract over store | `executionManagerImpl` implements `ExecutionManager` over `ExecutionStore` + `Serializer`, with static assertion `var _ ExecutionManager` | common/persistence/execution_manager.go:28-60 |
| Decorator stack composition | Factory wraps manager with rate-limiter then metrics decorators | common/persistence/client/factory.go:109-120 |
| Retry decorators with compile-time assertions | Seven retryable clients, each asserting `var _ XManager = (*xRetryablePersistenceClient)(nil)` | common/persistence/persistence_retryable_clients.go:11-61 |
| Error contract: conflict semantics | `IsConflictErr` groups three condition-failed error types | common/persistence/data_interfaces.go:1394-1402 |
| Error contract: outcome semantics | `OperationPossiblySucceeded` maps typed persistence errors to "write definitely not committed" vs possibly-committed | common/persistence/error_type.go:5-22 |
| Error contract: retryability | `IsPersistenceTransientError` (Unavailable/ResourceExhausted) and `IsServiceTransientError` (blacklist of non-retryable) | common/util.go:265-295 |
| Error contract: redirect semantics | `ShardOwnershipLost` carries `OwnerHost` and round-trips via gRPC status details (`errordetailsspb.ShardOwnershipLostFailure`) | common/serviceerror/shard_ownership_lost.go:11-57 |
| Consumer of error semantics | `BasicRedirector.redirectLoop` re-issues operation against `solErr.OwnerHost` | client/history/redirector.go:61-86 |
| Service client contract | `client.Bean` aggregates typed gRPC clients; remote clients lazily created and evicted on cluster-metadata change | client/client_bean.go:23-31, 96-113 |
| Generated decorator stack | `//go:generate go run ../../cmd/tools/genrpcwrappers` produces client, metric, and retryable wrappers implementing the generated gRPC interface | client/frontend/client.go:1-2; client/frontend/metric_client_gen.go:1-10 |
| Compile-time substitutability for client impl | `var _ workflowservice.WorkflowServiceClient = (*clientImpl)(nil)` | client/frontend/client.go:21 |
| CHASM component contract | `Component` (LifecycleState), `TerminableComponent` (Terminate with documented forced-termination triggers), `RootComponent` (ContextMetadata via gRPC trailers) | chasm/component.go:13-52 |
| Forward compatibility idiom | `mustEmbedUnimplementedComponent()` + `UnimplementedComponent` embed for adding methods without breaking implementations | chasm/component.go:17, 54-59 |
| Engine contract | `Engine` interface: StartExecution, UpdateWithStartExecution, UpdateComponent, ReadComponent, PollComponent, DeleteExecution, NotifyExecution | chasm/engine.go:16-59 |
| Semantic guarantee in engine contract | PollComponent documents monotonic-predicate requirement and `(nil, nil, nil)` long-poll-timeout convention | chasm/engine.go:368-382 |
| Context contract | `Context` (read) vs `MutableContext` (AddTask/SetRequestLinks/SetUserMetadata); `withValue` must return same concrete type; `RequestHeader` documented as empty outside gRPC-request contexts | chasm/context.go:17-69, 231-236 |
| Task handler contracts | `SideEffectTaskHandler` (runs outside lock, Go ctx, has Discard), `PureTaskHandler` (holds write lock, no I/O), `TaskValidator.Validate` with documented (true,nil)/(false,nil)/error semantics | chasm/task.go:27-73 |
| Base-struct embedding requirement | `SideEffectTaskHandlerBase`/`PureTaskHandlerBase` provide unexported marker methods and default Discard returning `ErrTaskDiscarded` | chasm/task_handler_base.go:5-18 |
| Library contract | `Library` interface with `mustEmbedUnimplementedLibrary()`; FQN = `libName.name` | chasm/library.go:11-23, 54-60 |
| Registration-time validation | Registry rejects duplicate FQNs, reserved IDs, ID collisions (hash of FQN), struct-kind violations, conflicting context-value keys | chasm/registry.go:234-332, 259-267 |
| Cross-field contract validation | Component with `Field[*Visibility]` must declare `WithBusinessIDAlias` or registration fails | chasm/registry.go:344-353 |
| Startup fail-fast | `chasm.Module` fx.Invoke registers CoreLibrary at boot; errors abort startup; production engine provided as `chasm.Engine` via fx | chasm/fx.go:5-11; service/history/chasm_engine.go:71-78 |
| Second independent engine implementation | `chasmtest.Engine` "implements [chasm.Engine]… matching the behavior of the production engine as closely as possible without persistence or shard logic", asserts `var _ chasm.Engine`, replicates reuse/conflict policies and `(nil, nil)` poll timeout | chasm/chasmtest/test_engine.go:26-42, 80, 188-224, 313-317 |
| Conformance suite across datastores | Shared `TestBase` with `NewTestBaseWithCassandra` / `NewTestBaseWithSQL` (mysql/postgresql/sqlite); same suites run per backend | common/persistence/persistence-tests/persistence_test_base.go:123-183 |
| Conformance suite instantiation | `persistencetests.NewTestBaseWithSQL(...)` in sqlite/mysql/postgresql tests; `NewTestBaseWithCassandra` in cassandra tests | common/persistence/tests/sqlite_test.go:555; common/persistence/tests/mysql_test.go:183; common/persistence/tests/cassandra_test.go:278 |
| SQL-dialect conformance layer | `common/persistence/sql/sqlplugin/tests/` runs per-table suites against each dialect plugin (mysql/postgresql/sqlite implement `sqlplugin.Plugin` interfaces) | common/persistence/sql/sqlplugin/tests/history_current_execution.go:1; common/persistence/sql/sqlplugin/interfaces.go:32-39 |
| Mock-based substitutability | 122 non-mock files carry `//go:generate mockgen`; service mocks generated under `api/historyservicemock/v1/` | api/historyservicemock/v1/service.pb.mock.go:1; grep "go:generate mockgen" = 122 files |
| Host application contract | `temporal.Server` interface (Start/Stop) + `Services` string list kept "to keep ServerOptions interface stable" | temporal/server.go:16-41 |
| Test server over public API | `temporaltest.TestServer` runs the real server on loopback for end-to-end tests, forcing `FailWorkflow` panic policy | temporaltest/server.go:19-63 |
| Archiver provider contract | `HistoryArchiver`/`VisibilityArchiver` document context-expiry, retry, and lossless behavior obligations per implementation | common/archiver/interface.go:44-94 |
| Task contract with capability interfaces | `Task` base interface plus optional `HasVersion`, `HasOutboundTaskGroup`, `HasDestination`, `HasArchetypeID` capability interfaces | service/history/tasks/task.go:14-52 |
| Request validation as contract enforcement | `RequestValidator` rejects incompatible reuse/conflict policy combinations with typed `serviceerror.InvalidArgument` sentinels | chasm/lib/workflow/validator.go:21-31, 111-124 |
| Wire-level error detail schemas | `api/errordetails/v1/message.pb.go` generated from `proto/internal/temporal/server/api/errordetails` carries structured error payloads across services | api/errordetails/v1/message.pb.go:1; proto/internal/ tree |

## Answers to Dimension Questions

**1. Are interfaces small, coherent, and owned by the consumer side?**
Mixed. The CHASM framework is the strongest example of consumer-side, capability-scoped design: `Component`, `TerminableComponent`, `RootComponent` form an inheritance ladder (chasm/component.go:13-52), and task handlers are split into `PureTaskHandler` vs `SideEffectTaskHandler` with a shared `TaskValidator` (chasm/task.go:30-73), with optional capability interfaces like `HasDestination` for tasks (service/history/tasks/task.go:44-47). The persistence layer is the counter-example: `ExecutionStore` (common/persistence/persistence_interface.go:116-167) and `ExecutionManager` (common/persistence/data_interfaces.go:1116-1173) each carry ~30 methods spanning workflow CRUD, task CRUD, DLQ management, and history-branch tree operations — several unrelated subdomains in one interface. Ownership is also producer-side in places: persistence interfaces live in the same package as the manager implementations that wrap them, and the `client.Bean` interface (client/client_bean.go:25-31) is defined next to its implementation.

**2. Do contracts specify behavior, not just method signatures?**
Yes, extensively — this is the codebase's defining trait. Examples: `CompleteTasksLessThan` documents the UnknownNumRowsAffected/limit contract (common/persistence/data_interfaces.go:1186-1195); `UpdateTaskQueueUserData` documents version-increment and error obligations (common/persistence/data_interfaces.go:1201-1208); `CreateQueue` must return a specific sentinel error (common/persistence/data_interfaces.go:1266-1267); `PollComponent` requires monotonic predicates and defines the timeout return convention (chasm/engine.go:368-382); `Validate` defines a three-way (true,nil)/(false,nil)/error protocol (chasm/task.go:68-72); `HistoryArchiver.Archive` specifies context-expiry and retry obligations (common/archiver/interface.go:45-52). Error *values* are part of the contract: `OperationPossiblySucceeded` (common/persistence/error_type.go:5-22) and `IsServiceTransientError` (common/util.go:276-295) give the error taxonomy operational meaning that every store/backend implementation must honor, and `ShardOwnershipLost.OwnerHost` is a machine-usable redirect hint transported through gRPC status details (common/serviceerror/shard_ownership_lost.go:40-48) and consumed by the redirector (client/history/redirector.go:79-85).

**3. Can providers, tools, stores, and runtimes be replaced safely?**
Largely yes. (a) Datastores: Cassandra and all SQL dialects implement the same `DataStoreFactory`/store contracts and are verified by the shared `persistencetests.TestBase` suites (common/persistence/persistence-tests/persistence_test_base.go:174-183) instantiated per backend (common/persistence/tests/cassandra_test.go:278, mysql_test.go:183, sqlite_test.go:555), plus a second dialect-level conformance layer in `sqlplugin/tests`. (b) Engines: the production `ChasmEngine` (service/history/chasm_engine.go:71-78) and `chasmtest.Engine` (chasm/chasmtest/test_engine.go:80) are two independent implementations of `chasm.Engine`, with the test engine deliberately replicating conflict/reuse-policy and poll-timeout semantics. (c) Cross-cutting concerns are substitutable decorators (retry, metrics, rate limit) inserted by the factory (common/persistence/client/factory.go:109-120), so backends never re-implement them. (d) Mocks are generated for 122 interfaces, making consumer-side substitution routine. Caveats: some contracts rely on undocumented invariants — e.g., `MetadataManager.WatchNamespaces` (common/persistence/data_interfaces.go:1228) has semantics discoverable only from implementations, and the `FairTaskManager` alias (common/persistence/data_interfaces.go:1213) has no SQL implementation yet, with the conformance suite's error check commented out (common/persistence/persistence-tests/persistence_test_base.go:240-242).

**4. Are compatibility failures caught early by tests or validation?**
Mostly yes, at four timescales. Compile time: static assertions such as `var _ ExecutionManager = (*executionRetryablePersistenceClient)(nil)` (common/persistence/persistence_retryable_clients.go:55-61) and `var _ workflowservice.WorkflowServiceClient = (*clientImpl)(nil)` (client/frontend/client.go:21). Schema time: protobuf codegen keeps wire contracts in sync (`make proto` / `make update-go-api`, AGENTS.md commands; generated mocks under `api/historyservicemock/v1/`). Registration time: `Registry.Register` fails fast on FQN/ID collisions, reserved archetype IDs, struct-kind violations, context-key conflicts (chasm/registry.go:247-277), and cross-field violations like missing business-ID aliases (chasm/registry.go:344-353); it is invoked during fx startup (chasm/fx.go:8-10). Test time: cross-backend conformance suites and the second engine implementation. Gaps: some option constructors panic instead of erroring (chasm/registrable_component.go:96-101), and several `TransitionOption`s are silently ignored by engine methods ("opts are currently ignored", chasm/engine.go:291, 334, 375), which is a contract that fails late or never.

## Architectural Decisions

1. **Two-layer persistence contract (Store vs Manager).** Stores are dumb, per-datastore CRUD with `Internal*` request/response types carrying serialized blobs; managers add serialization, stats, and paging (common/persistence/persistence_interface.go:24-28, common/persistence/execution_manager.go:28-60). This confines datastore-specific code to one layer and lets serialization evolve independently.

2. **Decorators over inheritance for cross-cutting concerns.** Rate limiting, metrics, retries, and health signals are separate wrapper types stacked by the factory (common/persistence/client/factory.go:109-120; common/persistence/persistence_retryable_clients.go:11-61). Any backend gets identical operational behavior for free — a substitutability guarantee by construction.

3. **Code-generated client wrappers.** A repo-owned generator (`cmd/tools/genrpcwrappers`, referenced at client/frontend/client.go:2) emits timeout, metric, and retryable wrappers for every service client, so adding an RPC updates all decorator layers mechanically — the contract surface cannot drift from the decorators.

4. **Registry with registration-time validation for the CHASM component model.** Components/tasks are registered by FQN with deterministic 32-bit IDs (`GenerateTypeID`, chasm/registrable_component.go:200-205); collisions, reserved names, and cross-field contract violations are rejected at startup (chasm/registry.go:234-353). This is schema-time validation for Go struct contracts.

5. **Unexported marker methods for embedding enforcement.** `mustEmbedUnimplementedComponent()` (chasm/component.go:17), `pureTaskHandler()`/`sideEffectTaskHandler()` (chasm/task.go:39,47), and base structs providing defaults (chasm/task_handler_base.go:5-18) reserve evolution room: new methods can get default implementations in the `Unimplemented*` embed without breaking external implementations — Go's answer to protobuf's `UnimplementedXServiceServer`.

6. **Behavior encoded in typed errors + gRPC status details.** Server-side error types implement `Status()` attaching proto error details, and a symmetric constructor rebuilds them from received statuses (common/serviceerror/shard_ownership_lost.go:35-56), enabling cross-process contracts like shard-owner redirection (client/history/redirector.go:72-86).

7. **Test doubles that are real implementations.** Rather than only gomock stubs, the repo ships a second full `chasm.Engine` (chasm/chasmtest/test_engine.go:26-42) and a real-server `temporaltest.TestServer` (temporaltest/server.go:19-33), so contract conformance is exercised against behavioral clones, not just call-recording mocks.

## Notable Patterns

- **Capability interfaces on top of a base contract**: `Task` plus `HasVersion`/`HasDestination`/`HasOutboundTaskGroup`/`HasArchetypeID` (service/history/tasks/task.go:35-52); processors type-assert to opt into features instead of bloating the base interface.
- **Generic adapter functions**: `StartExecution[C RootComponent, I any]`, `UpdateComponent[C any, R []byte | ComponentRef, ...]` wrap the non-generic `Engine` interface, giving callers type safety while the interface stays simple (chasm/engine.go:207-236, 296-331); panics inside user callbacks are captured into errors via `defer log.CapturePanic` (chasm/engine.go:218).
- **Context-keyed dependency injection inside the framework**: the engine is passed through `NewEngineContext` (chasm/engine.go:447-452) and installed by a gRPC interceptor (chasm/interceptors.go:21-33); per-library values are injected via `WithContextValues` with globally-unique-key enforcement (chasm/registrable_component.go:159-180; chasm/registry.go:259-267).
- **Immutable/mutable context split**: `Context` vs `MutableContext` (chasm/context.go:17-108) makes read-only transitions statically unable to emit tasks or metadata, and `ContextWithValue[C Context]` preserves the concrete context type (chasm/context.go:284-287).
- **Conformance-suite-per-backend harness**: one `TestBase`, N cluster adapters (Cassandra/MySQL/PostgreSQL/SQLite) (common/persistence/persistence-tests/persistence_test_base.go:123-191), with time-precision tolerance for cross-DB timestamp semantics (line 41-43).
- **Error-classification helpers as part of the public surface**: `IsConflictErr` (common/persistence/data_interfaces.go:1394-1402) and `OperationPossiblySucceeded` (common/persistence/error_type.go:5-22) let callers branch on semantics without knowing every backend error type.

## Tradeoffs

- **Large persistence interfaces vs. single-file discoverability.** `ExecutionManager`'s ~30 methods span workflow state, tasks, DLQ, and history trees (common/persistence/data_interfaces.go:1116-1173). Every new backend must implement everything at once; the payoff is that managers/decorators and the conformance suite can treat the whole domain uniformly.
- **Generated code vs. debuggability.** Client decorators and mocks are generated (client/frontend/client.go:1-2); correctness is high, but stack traces and reviews traverse generated files, and the generator itself becomes critical infrastructure.
- **Reflection-based registry vs. static typing.** CHASM validates Go struct contracts at runtime registration (chasm/registry.go:269-277) and erases types through `any` inside `RegistrableTask` (chasm/registrable_task.go:28-31), re-asserting `component.(C)` at the generic boundary (chasm/registrable_task.go:54-57). Flexibility for library authors costs a layer of unchecked assertions.
- **Behavior-rich doc comments vs. enforceability.** Contracts like "predicate must be monotonic" (chasm/engine.go:373-375) or archiver losslessness (common/archiver/interface.go:81-83) are documented but not machine-checked; only the conformance suites partially verify them.
- **`Close()` cannot fail.** The `Closeable` contract returns nothing, with an in-code TODO acknowledging the limitation (common/persistence/data_interfaces.go:1099-1103) — shutdown errors are unrepresentable in the current contract.

## Failure Modes / Edge Cases

- **Contract drift between engines.** `chasmtest.Engine` must manually mirror production semantics; the `Speculative` transition option is unimplemented in both and explicitly flagged (chasm/engine.go:144-150; chasm/chasmtest/test_engine.go:644-647), so a caller using it gets silently different behavior than documented.
- **Silently ignored options.** `UpdateComponent`, `ReadComponent`, and `PollComponent` ignore `TransitionOption`s ("opts are currently ignored", chasm/engine.go:291, 334, 375) — passing `WithBusinessIDPolicy` there is a no-op, a latent misuse trap.
- **Lifecycle stubs.** `Workflow.LifecycleState` always returns Running with NOTE comments explaining why (chasm/lib/workflow/workflow.go:53-61), and `Terminate` returns an internal error (chasm/lib/workflow/workflow.go:68-73) — implementers of `RootComponent` cannot rely on the interface's implied guarantees yet.
- **Hash-ID collision handling.** Component/task IDs are 32-bit FQN fingerprints; collisions are detected only at registration with an error (chasm/registry.go:255-257, 310-312), so a new library name can break startup — fail-fast, but a hard failure mode for third-party libraries.
- **Mixed validation styles.** `WithBusinessIDAlias`/`WithSearchAttributes` panic on duplicate aliases (chasm/registrable_component.go:94-101, 130-150) while `Registry.Register` returns errors (chasm/registry.go:69-102); option-time panics bypass the registry's structured error path.
- **Partial backend support masked in tests.** The FairTaskManager conformance error check is commented out with a TODO (common/persistence/persistence-tests/persistence_test_base.go:240-242), so the suite would not catch a missing SQL implementation.
- **Ambiguity in outcome semantics is explicit but default-permissive.** `OperationPossiblySucceeded` returns true for *unknown* error types (common/persistence/error_type.go:19-21), meaning new backend errors are conservatively treated as possibly-committed — safe for dedup, but it pushes complexity onto every writer to handle "maybe succeeded".

## Future Considerations

- Split `ExecutionStore`/`ExecutionManager` along subdomain lines (workflow state vs tasks vs DLQ vs history trees) or provide partial interfaces so new backends can land incrementally; the `InternalGetOrCreateShardRequest` callback pattern (common/persistence/persistence_interface.go:213-217) shows the team already pushes behavior into contracts where splitting is hard.
- Give `TransitionOption`s real effect or remove them from `ReadComponent`/`PollComponent` signatures to eliminate no-op parameters (chasm/engine.go:291, 334, 375).
- Replace option-time panics with registration-time errors so all CHASM contract violations surface through `Registry.Register` (chasm/registrable_component.go:96-101).
- Extend the conformance-suite pattern to the CHASM engine: `chasmtest.Engine` currently mirrors semantics by hand; shared behavioral tests (reuse/conflict policy tables, poll-timeout convention) could run against both engines.
- Allow `Closeable.Close()` to return errors once call sites can handle them (common/persistence/data_interfaces.go:1100).
- Stabilize `MetadataManager.WatchNamespaces` semantics (delivery guarantees, replay) in the interface doc (common/persistence/data_interfaces.go:1228) since it is part of a widely implemented contract.

## Questions / Gaps

- **No evidence found** of a formal, versioned compatibility policy for the internal persistence contracts (e.g., how `Internal*` request structs may evolve across server versions); the `Internal` prefix and two-layer split imply privacy, but no doc states the evolution rules.
- **No evidence found** of contract tests that verify `OperationPossiblySucceeded`/`IsConflictErr` classification for *every* backend implementation; classification correctness appears to rely on each backend translating errors into the shared types (e.g., `ConditionFailedError`), but a per-backend error-taxonomy conformance test was not located in `common/persistence/persistence-tests` or `common/persistence/tests`.
- `RegistrableTask.componentGoType` carries the comment "It is not clear how this one is used" (chasm/registrable_task.go:13), suggesting the task-registration contract itself has unresolved design debt; the impact on substitutability could not be fully determined from this source.
- The claimed equivalence between `chasmtest.Engine` and the production engine is asserted in comments (chasm/chasmtest/test_engine.go:29-30) but no shared conformance suite proving it was found in this source.

---

Generated by `24.02-interface-contract-design` against `temporal`.
