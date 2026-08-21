# Source Analysis: opa

## 13.01 Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | opa |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go (with auxiliary Rego policy language and WASM SDK) |
| Analyzed | 2026-08-21 |

## Summary

OPA organizes its error taxonomy by **subsystem** rather than by an external "model/provider/tool/validation/policy/context/user/infrastructure/timeout" axis. Each subsystem declares its own typed `Error` struct with a closed list of `Code` constants, an `Is*` predicate family, and an `Is(target error) bool` method so callers can match on `errors.Is`. The taxonomy is split across `ast` (parse/compile/format), `topdown` (evaluation), `storage` (in-memory/disk/backend), `sdk` (Go SDK wrapper), `internal/wasm/sdk/opa/errors` (WASM host ABI, closed code set), `server/types` (REST envelope), and per-plugin status types (`plugins/bundle/status.go`, `plugins/logs/status/status.go`). Errors are routed at the HTTP boundary in `v1/server/writer/writer.go:27-42`, where `ErrorAuto` dispatches on `types.IsBadRequest`, `storage.IsWriteConflictError`, `topdown.IsError`, `storage.IsInvalidPatch`, and `storage.IsNotFound` to map onto the right `CodeInternal`/`CodeInvalidParameter`/`CodeResourceConflict`/`CodeResourceNotFound` HTTP response. The taxonomy is documented in `docs/docs/errors/index.md` with a per-stage table (parsing / compilation / evaluation) and individual pages per error message.

The taxonomy is mature, exercised by tests, and stable across the v0/v1 module split, but it is **not unified**: there is no single "OPA Error" base class, and the `code` of a `topdown.Error` (`eval_builtin_error`) tells you nothing about whether the failure was from the document store, the network, the policy, or an internal panic. To make that determination, callers must chain `errors.As`/`errors.Is` on the typed wrapper. New error codes can be added safely because `Is()` matching uses empty-wildcard semantics, but a caller using a `switch` on `err.Code` (as `writer.ErrorAuto` does internally) will silently fall through to the default branch and may mis-route.

## Rating

**7 / 10** — Clear closed-code taxonomy per subsystem, used for routing at the HTTP boundary, documented in `docs/docs/errors/index.md`, exercised by tests (`v1/topdown/errors_test.go`, `v1/storage/errors_test.go`, `v1/plugins/bundle/errors_test.go`, `internal/wasm/sdk/opa/errors/errors.go`), and conforms to Go `errors.Is`/`errors.As` conventions. Deductions: the taxonomy is not unified across subsystems (no common base, no common `Code` enum), the code namespaces overlap (`internal_error` exists in both `sdk/opa.go` and `server/types/types.go` and `internal/wasm/sdk/opa/errors/errors.go`), some error types are unclassified (`v1/loader/errors.go:15`, `v1/rego/rego.go:593`), and the code labels do not map cleanly to the "model/provider/tool/validation/policy/context/user/infrastructure/timeout" categories called out in the dimension prompt — they map to subsystem boundaries (AST/topdown/storage/sdk/server) instead.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| AST typed error + code constants | `Error` struct with `Code`, `Message`, `Location`, `Details`; codes `ParseErr`="rego_parse_error", `CompileErr`="rego_compile_error", `TypeErr`="rego_type_error", `UnsafeVarErr`="rego_unsafe_var_error", `RecursionErr`="rego_recursion_error", `FormatErr`="rego_format_error" | `v1/ast/errors.go:48-87` |
| AST `IsError` predicate | `IsError(code, err)` predicate on `*ast.Error` | `v1/ast/errors.go:68-74` |
| AST `Errors` aggregator | `type Errors []*Error` with `Sort()`, `Error()` for plural errors | `v1/ast/errors.go:14-46` |
| Topdown typed error + code constants | `Error` struct with `Code`, `Message`, `Location`, `err`; codes `InternalErr`="eval_internal_error", `CancelErr`="eval_cancel_error", `ConflictErr`="eval_conflict_error", `TypeErr`="eval_type_error", `BuiltinErr`="eval_builtin_error", `WithMergeErr`="eval_with_merge_error" | `v1/topdown/errors.go:26-60` |
| Topdown `IsError` / `IsCancel` predicates | `IsError(err)` via `errors.As`, `IsCancel(err)` via `errors.Is` against `&Error{Code: CancelErr}` | `v1/topdown/errors.go:62-71` |
| Topdown `Error.Is` matching | `Error.Is(target error) bool` with empty-wildcard semantics on `Code`, `Message`, `Location` so `errors.Is(err, &topdown.Error{Code: x})` works | `v1/topdown/errors.go:73-82` |
| Topdown `Halt` control-flow error | `Halt{Err error}` used by built-ins to abort evaluation; `Unwrap()` | `v1/topdown/errors.go:14-24` |
| Storage typed error + code constants | `Error` struct with `Code`, `Message`; codes `InternalErr`="storage_internal_error", `NotFoundErr`="storage_not_found_error", `WriteConflictErr`="storage_write_conflict_error", `InvalidPatchErr`="storage_invalid_patch_error", `InvalidTransactionErr`="storage_invalid_txn_error", `TriggersNotSupportedErr`="storage_triggers_not_supported_error", `WritesNotSupportedErr`="storage_writes_not_supported_error", `PolicyNotSupportedErr`="storage_policy_not_supported_error" | `v1/storage/errors.go:11-48` |
| Storage `Is*` predicates | `IsNotFound`, `IsWriteConflictError`, `IsInvalidPatch`, `IsInvalidTransaction`, plus deprecated `IsIndexingNotSupported` | `v1/storage/errors.go:57-96` |
| SDK Go error | `Error` struct with `Code`, `Message`; constant `UndefinedErr`="opa_undefined_error"; `undefinedDecisionErr` constructor | `v1/sdk/opa.go:528-548` |
| SDK `IsUndefinedErr` predicate | `IsUndefinedErr(err)` typed assertion | `v1/sdk/opa.go:550-554` |
| SDK propagates `IsNotFound` | `bundles()` swallows `storage.IsNotFound` so missing data is not a hard error | `v1/sdk/opa.go:748-762` |
| WASM host error (closed code set) | `Error` struct; constants `InvalidConfigErr`="invalid_config", `InvalidPolicyOrDataErr`="invalid_policy_or_data", `InvalidBundleErr`="invalid_bundle", `NotReadyErr`="not_ready", `InternalErr`="internal_error", `CancelledErr`="cancelled" | `internal/wasm/sdk/opa/errors/errors.go:12-30` |
| WASM `New()` panics on unknown code | `New(code, msg)` enforces a closed code set by panicking with "unknown error code: <code>" | `internal/wasm/sdk/opa/errors/errors.go:39-46` |
| WASM `Is` matching | `Error.Is(target error) bool` with empty-wildcard semantics on `Code` and `Message` | `internal/wasm/sdk/opa/errors/errors.go:62-77` |
| WASM `IsCancel` predicate | `IsCancel(err)` via `errorHasCode(err, CancelledErr)` | `internal/wasm/sdk/opa/errors/errors.go:57-60` |
| Server REST error envelope | `CodeInternal`="internal_error", `CodeEvaluation`="evaluation_error", `CodeUnauthorized`="unauthorized", `CodeInvalidParameter`="invalid_parameter", `CodeInvalidOperation`="invalid_operation", `CodeResourceNotFound`="resource_not_found", `CodeResourceConflict`="resource_conflict", `CodeUndefinedDocument`="undefined_document" | `v1/server/types/types.go:118-128` |
| Server `ErrorV1` response + `BadRequestErr` | `ErrorV1{Code, Message, Errors}` JSON envelope; `BadRequestErr` string sentinel for client-input errors | `v1/server/types/types.go:130-169`, `v1/server/types/types.go:614-637` |
| Server `IsBadRequest` predicate | `IsBadRequest(err)` type-asserts on `BadRequestErr` | `v1/server/types/types.go:633-637` |
| Server error routing (HTTP) | `ErrorAuto` dispatches `types.IsBadRequest` → 400, `storage.IsWriteConflictError` → 404, `topdown.IsError` → 500, `storage.IsInvalidPatch` → 400, `storage.IsNotFound` → 404, default → 500 | `v1/server/writer/writer.go:27-42` |
| Server error routing (compile handler) | Uses `types.NewErrorV1(types.CodeInvalidParameter, …).WithASTErrors(err)` to surface AST error lists | `v1/server/compile_handler.go:282`, `v1/server/compile_handler.go:312`, `v1/server/compile_handler.go:485-487` |
| Server error routing (authorizer) | `MsgUnauthorizedUndefinedError`, `MsgUnauthorizedError`, `MsgUndefinedError` constants drive the 401/500 responses | `v1/server/types/types.go:172-185`, `v1/server/authorizer/authorizer.go:136-164` |
| Rego package `Errors` aggregator (untyped) | `type Errors []error` with `Error()`; no `Is*` predicate for individual codes | `v1/rego/rego.go:592-619` |
| Rego `HaltError` | `HaltError{err}` for custom built-ins to abort evaluation; `errors.As(&e)` used in `finishFunction` | `v1/rego/errors.go:5-18`, `v1/rego/rego.go:3027-3042` |
| Rego `IsPartialEvaluationNotEffectiveErr` | `IsPartialEvaluationNotEffectiveErr` checks the typed `Errors[0] == errPartialEvaluationNotEffective` sentinel | `v1/rego/rego.go:609-619` |
| Loader `Errors` aggregator (untyped) | `type Errors []error` with `add()` that flattens `ast.Errors`; no per-code predicate | `v1/loader/errors.go:14-56` |
| Compiler typed errors | `invalidEntrypointErr`, `undefinedEntrypointErr` per-instance errors; no `Is*` predicate | `v1/compile/compile.go:882-897` |
| Bundle plugin status error aggregation | `SetError` branches on `ast.Errors` (compile error), `download.HTTPError` (HTTP error), and generic fallback, all normalized to `code="bundle_error"` | `v1/plugins/bundle/status.go:65-98` |
| Bundle plugin typed error | `Error{BundleName, Code, HTTPCode, Message, Err}`; `Errors []Error` with `Unwrap() []error`; `NewBundleError` classifies by `errors.As` on `download.HTTPError` | `v1/plugins/bundle/errors.go:25-53` |
| Bundle error tests | `TestErrors`, `TestUnwrap`, `TestUnwrap`, `TestHTTPErrorWrapping`, `TestASTErrorsWrapping`, `TestGenericErrorWrapping` verify classification and `errors.As` behavior | `v1/plugins/bundle/errors_test.go:11-144` |
| Decision log status error aggregation | `SetError` branches on `HTTPError` (logs-specific) and generic fallback, all normalized to `code="decision_log_error"` | `v1/plugins/logs/status/status.go:32-51` |
| Decision log `HTTPError` | `HTTPError{StatusCode}` distinct from `download.HTTPError` | `v1/plugins/logs/status/status.go:57-63` |
| Coverage typed error | `CoverageThresholdError{Coverage, Threshold, Report}` with formatted Line-by-line renderer | `v1/cover/threshold_error.go:16-46` |
| Download HTTP error | `HTTPError{StatusCode}`; carried in `bundle.Error` and `decision_log.Status` dispatch | `v1/download/download.go:441-447` |
| Tester cancel handling | `IsCancel(err)` from both topdown and WASM decides whether to stop the test run; `IsCancel && ctx.Err() == context.DeadlineExceeded` distinguished from user cancel | `v1/tester/runner.go:999`, `v1/tester/runner.go:1146` |
| Tester `BufferEmpty` / `UploadCancelled` | Two private control-flow error types raised by the buffer and matched with `errors.Is` | `v1/plugins/logs/plugin.go:965-981`, `v1/plugins/logs/plugin.go:693, 1009` |
| Auth provider typed errors | `gcpMetadataError`, `azureManagedIdentitiesError`, `awsCredentialCheckError` per cloud; `Unwrap` on the GCP variant | `v1/plugins/rest/gcp.go:30-41`, `v1/plugins/rest/azure.go:37-46`, `v1/plugins/rest/auth.go:739-770` |
| Oracle typed error | `Error{Code}` with `ErrNoMatchFound` and `ErrNoDefinitionFound` value sentinels | `v1/ast/oracle/oracle.go:11-44` |
| REPL typed error | `Error{Code, Message}` with `BadArgsErr`="bad arguments" | `v1/repl/errors.go:9-30` |
| Bundle verify sentinel errors | `errUnauthorized` (`errors.New("401 Unauthorized")`) wrapped with `%w` and matched with `errors.Is` in the test harness | `v1/download/testharness.go:27-155`, `v1/download/testharness.go:391` |
| Term `UnknownValueErr` | `UnknownValueErr` sentinel for "value not resolvable"; `IsUnknownValueErr` predicate | `v1/ast/term.go:167-179` |
| Topdown error test (wrapping, matching) | `TestErrorWrapping` covers `IsError`, `IsCancel`, `Halt`, `errors.Is` on `Code`, `Message`, `Location`, and wrapped builtin errors | `v1/topdown/errors_test.go:12-117` |
| Storage error test | `TestIsNotFound` asserts that `IsNotFound` distinguishes `NotFoundErr` from `InternalErr` | `v1/storage/errors_test.go:9-27` |
| WASM `New`-panic test | `New` panics on unknown codes — checked by behavior of the closed-enum pattern | `internal/wasm/sdk/opa/errors/errors.go:39-46` |
| Documentation: errors index | Per-stage table (parsing / compilation / evaluation) of error codes with links to detail pages | `docs/docs/errors/index.md:14-32` |
| Documentation: error categories by stage | Parsing, compilation, evaluation stages documented with examples and `--strict-builtin-errors` flag | `docs/docs/errors/index.md:46-148` |
| Documentation: contrib adding built-in | Built-in functions tested with `want_error_code: eval_builtin_error` | `docs/docs/contrib-adding-builtin-functions.md:135` |
| Documentation: build-in error code enums | `topdown.BuiltinErr` constant cited in custom-built-in guide | `v1/topdown/errors.go:53-56` |
| Test that `errors.As` recovers an `ast.Errors` from `rego` | `if !errors.As(err, &errs)` pattern in `rego_external_source_test.go` | `v1/rego/rego_external_source_test.go:320-321` |
| Configuration error limit | `Plugins.MaxErrors` / `runtime.Params.ErrorLimit` caps the number of plugin errors before abort | `v1/plugins/plugins.go:390`, `v1/runtime/runtime.go:179-181`, `v1/runtime/runtime.go:517` |
| Decision log `MaskError` | `errMaskInvalidObject` sentinel for invalid mask rule input | `v1/plugins/logs/mask.go:29` |

## Answers to Dimension Questions

1. **Are errors classified by source?**
   - Yes, but by **subsystem** (AST compile vs. topdown evaluation vs. storage backend vs. SDK vs. server HTTP) rather than by the "model/provider/tool/validation/policy/context/user/infrastructure/timeout" axis named in the dimension. Each subsystem declares its own `Error` struct with a closed `Code` constant set (`v1/ast/errors.go:48-66`, `v1/topdown/errors.go:35-60`, `v1/storage/errors.go:11-42`, `internal/wasm/sdk/opa/errors/errors.go:12-30`). Cross-subsystem classification (e.g. "is this a network error or a policy error?") has to be reconstructed by the caller using `errors.As` on `download.HTTPError`, `topdown.Error`, `ast.Errors`, `storage.Error`, and `bundle.Error` wrappers together.

2. **Is the taxonomy used for handling?**
   - Yes. The flagship example is `v1/server/writer/writer.go:27-42` (`ErrorAuto`) which dispatches `types.IsBadRequest` → 400, `storage.IsWriteConflictError` → 404, `topdown.IsError` → 500, `storage.IsInvalidPatch` → 400, `storage.IsNotFound` → 404, default → 500. Bundle and decision-log status objects (`v1/plugins/bundle/status.go:65-98`, `v1/plugins/logs/status/status.go:32-51`) similarly use `errors.As` to differentiate AST/HTTP/generic errors. The tester uses `topdown.IsCancel` / `wasm_errors.IsCancel` to decide whether to stop the run (`v1/tester/runner.go:999`). The rego package detects `HaltError` via `errors.As` to convert a custom-built-in abort into a `topdown.Halt` (`v1/rego/rego.go:3027-3042`).

3. **Are error categories documented?**
   - Yes. `docs/docs/errors/index.md:14-32` provides a table that maps every documented error code to its stage (parsing / compilation / evaluation). Each entry has a dedicated page (`docs/docs/errors/rego-parse-error/*.md`, `docs/docs/errors/rego-type-error/*.md`, `docs/docs/errors/eval-conflict-error/*.md`, etc.) explaining the message, the cause, and how to fix it. The Go-side constants are documented with the constants themselves (`v1/ast/errors.go:48-66`, `v1/topdown/errors.go:35-60`, `v1/storage/errors.go:11-42`).

4. **Can new error types be added without breaking existing handling?**
   - Yes in most packages, because each `Error` implements an `Is(target error) bool` method with empty-wildcard semantics (`v1/topdown/errors.go:73-82`, `internal/wasm/sdk/opa/errors/errors.go:62-70`) and the call sites use `errors.Is`/`errors.As` rather than concrete `==` comparisons. Adding a new `Code` constant to `v1/ast/errors.go` or `v1/topdown/errors.go` is non-breaking for downstream callers. **Exception**: the WASM SDK enforces a closed code set in `New()` and panics on unknown codes (`internal/wasm/sdk/opa/errors/errors.go:39-46`), and the server writer's `switch` falls through to `default` (`v1/server/writer/writer.go:39-41`) so a new error type that is not handled gets a 500 instead of a more specific code. The `server/types` REST envelope (`Code*` constants) is itself an open extension surface — adding a new code does not break the JSON decoder.

## Architectural Decisions

- **Per-subsystem error types, not one global hierarchy.** Each subsystem owns its `Error` struct and its `Code` constants. There is no `opa.Error` interface or shared base class. This keeps packages independent but means callers must know which subsystem a failure came from to inspect `Code`.
- **Use Go `errors.Is` / `errors.As` consistently.** Every typed error implements `Unwrap()` and most implement `Is(target error) bool` with empty-wildcard semantics, so callers can pattern-match on `*Error{Code: x}` without a custom predicate. New errors are wired in via `errors.As` (e.g. `v1/server/writer/writer.go:33`, `v1/plugins/bundle/status.go:77-86`).
- **Closed-code WASM ABI.** The WASM host-facing error type (`internal/wasm/sdk/opa/errors/errors.go`) deliberately locks the code set (`InvalidConfigErr`, `InvalidPolicyOrDataErr`, `InvalidBundleErr`, `NotReadyErr`, `InternalErr`, `CancelledErr`) because the host ABI is a stable contract. Unknown codes panic in `New()` so an embedding SDK cannot silently emit a code the host cannot handle.
- **HTTP envelope is decoupled from internal codes.** `v1/server/types/types.go:118-128` defines its own `CodeInternal`/`CodeInvalidParameter`/`CodeResourceConflict`/etc. REST envelope codes. The server writer translates internal codes into envelope codes via a `switch` over `Is*` predicates (`v1/server/writer/writer.go:27-42`). This is intentional: the REST API is a public contract and the internal codes change; the translator is the only place that has to be updated.
- **Status types use a single `Code` per category.** `plugins/bundle/status.go:21` and `plugins/logs/status/status.go:19` each define a single `errCode` constant ("bundle_error", "decision_log_error") and stuff the underlying error/message into `HTTPCode`/`Message`/`Errors` fields. Callers see one code per category and pull the underlying details from the structured fields.
- **Control-flow errors are wrapped, not exposed.** `topdown.Halt{}` (`v1/topdown/errors.go:14-24`) and `rego.HaltError{}` (`v1/rego/errors.go:5-18`) are dedicated value types that satisfy `error` so the evaluator can use sentinel-style `errors.As` matching. Plugins follow the same pattern (`bufferEmpty{}`, `uploadCancelled{}` in `v1/plugins/logs/plugin.go:965-981`).

## Notable Patterns

- **Empty-wildcard `Is()` matching.** Both `v1/topdown/errors.go:73-82` and `internal/wasm/sdk/opa/errors/errors.go:62-70` define `Is(target error) bool` so that `errors.Is(err, &topdown.Error{Code: topdown.BuiltinErr})` matches any error whose `Code=="eval_builtin_error"` regardless of message and location. This is the cleanest example of using Go's `errors.Is` to maintain a closed-code taxonomy without coupling predicate helpers to specific messages.
- **Aggregator types consistent with stdlib `errors.Join`.** `v1/ast/errors.go:14-46` (`Errors []*Error`), `v1/loader/errors.go:14-39` (`Errors []error`), `v1/rego/rego.go:592-619` (`Errors []error`), `v1/plugins/bundle/errors.go:11-53` (`Errors []Error` with `Unwrap() []error`) all implement `Error() string` and either `Sort()` or `Unwrap()`. The bundle `Errors` type explicitly supports `errors.Join` semantics (`v1/plugins/bundle/errors.go:13-19`).
- **HTTP-error sentinel + `errors.As` fan-out.** `v1/download/download.go:441-447` defines `HTTPError{StatusCode}` once, and the same type is consumed by `v1/plugins/bundle/errors.go:33-44` (via `errors.As`), `v1/plugins/bundle/status.go:86-90`, and `v1/server/writer/writer.go:31-32` (via `storage.IsWriteConflictError`, which is a different but parallel pattern). Callers recover the HTTP status by `errors.As` rather than string matching.
- **SDK `IsNotFound` swallow in callers.** `v1/sdk/opa.go:751`, `v1/bundle/store.go:88`, `v1/bundle/store.go:173`, `v1/bundle/store.go:617`, `v1/server/server.go:1056-1068/1986/2233-2239/2520-2529`, `v1/storage/inmem/txn.go:355`, `v1/topdown/eval.go:1933-1941` all use `!storage.IsNotFound(err)` as the canonical "this is a real error" check. The pattern is consistent across the codebase.
- **Built-in emit `eval_builtin_error` and `Halt`.** Custom built-ins raise `HaltError` to stop evaluation; `v1/rego/rego.go:3027-3042` (`finishFunction`) catches it with `errors.As` and re-emits a `topdown.Error{Code: BuiltinErr}` wrapped in `topdown.Halt{Err: ...}`. This shows the intended escape hatch: built-ins stay typed as Go errors, but the evaluator gets a uniform `topdown.Error` code.

## Tradeoffs

- **Multiple parallel code namespaces.** The same string `"internal_error"` appears in `v1/sdk/opa.go`-adjacent `UndefinedErr`-style codes, `v1/server/types/types.go:120` (`CodeInternal`), and `internal/wasm/sdk/opa/errors/errors.go:26`. The dimension's "model vs. provider vs. tool" distinction is lost: a caller looking at a server response sees `internal_error` and has no way to tell whether the cause was a topdown eval panic, a storage transaction failure, or a SDK wiring problem.
- **No central error type or interface.** A function that may return any kind of OPA error has to enumerate `*ast.Error`, `*topdown.Error`, `*storage.Error`, `*sdk.Error`, `os.Error`, `bundle.Error`, `loader.Errors`, `rego.Errors`, etc. There is no `errors.As(err, &opaErr)` shim.
- **Aggregator types are not usable with `errors.Is` for a sub-code.** `v1/loader/errors.go:14-39` and `v1/rego/rego.go:592-619` provide `Errors []error` but no `IsError(code, err)` predicate analogous to `ast.IsError`. A caller wanting to know "is this loader error a parse error?" has to flatten the list and check each element.
- **HTTP writer's `switch` falls through.** `v1/server/writer/writer.go:27-42` enumerates known error types in a `switch`; a new error type added in a subsystem that the writer does not yet know about will be mapped to `default → 500 internal_error`. This is a load-bearing failure mode: silently mis-classifying a `BadRequestErr` as 500 is a regression.
- **Closed WASM code set panics on unknown values.** `internal/wasm/sdk/opa/errors/errors.go:39-46` is the only place in the codebase that panics on construction of an error. This is appropriate for a host ABI but it means a downstream embedder that adds a new code will crash the process at first use.
- **Status single-code-per-category hides the structured detail.** `v1/plugins/bundle/status.go:21` always sets `code="bundle_error"` even when the underlying cause is an AST compile error (which has its own "rego_compile_error" code) or an HTTP 503. The wire payload still carries the original code in `Errors []error`, but a client that only checks `code == "bundle_error"` cannot distinguish the failure class.

## Failure Modes / Edge Cases

- **Cross-subsystem error lost during wrapping.** `v1/rego/rego.go:3027-3042` re-wraps a `HaltError` into a `topdown.Error`, dropping the original error type from the `errors.As` chain unless the caller has kept a reference. The `Wrap` flow is correct (`topdown.Error.Wrap` at `v1/topdown/errors.go:108-115`), but with `*HaltError` interfaces this is a common source of "IsCancel said false" bugs.
- **Predicate vs. message divergence.** `topdown.IsError(err)` (`v1/topdown/errors.go:62-66`) only checks the type, not the code. A caller wanting "is this a cancellation?" must use `topdown.IsCancel` (`v1/topdown/errors.go:68-71`) which goes via `errors.Is`, not `IsError`. The two helpers are not interchangeable.
- **Writer fall-through mis-classification.** `v1/server/writer/writer.go:39-41` defaults to 500 internal_error. If a new typed error is added without updating the switch, the request is silently mis-classified. The pattern is tested in `v1/server/server_test.go` but not exhaustively.
- **WASM `New` panic.** `internal/wasm/sdk/opa/errors/errors.go:39-46` will panic if a host calls `New("unknown_code", "")`. This is intended for ABI stability but it also means a custom built-in that constructs an error from a user-supplied code cannot safely use `errors.New`.
- **Aggregators' `Error()` strings are not unique.** `v1/loader/errors.go:17-29` and `v1/rego/rego.go:595-607` both produce strings like `"1 error occurred: …"`. A caller that does `err.Error() == "1 error occurred: …"` to detect these types will mis-match across loaders.
- **Coverage error is in its own place.** `v1/cover/threshold_error.go:16-46` defines `CoverageThresholdError` with no `Is*` predicate and no `Code` constant. Test runners that want to detect "coverage not met" have to type-assert directly.

## Future Considerations

- **Unify the parallel code namespaces.** `internal_error`, `cancel_error`, `not_found_error` recur in `ast`, `topdown`, `storage`, `sdk`, `wasm/sdk/opa/errors`. A shared `opaerrors` package with one `Code` enum and one `Error{Op, Code, ...}` struct would let callers pattern-match across subsystems. The cost is a breaking change to every existing public error type.
- **Surface the source category in the error.** Adding a `Category` field (e.g. `Model`, `Provider`, `Tool`, `Validation`, `Policy`, `Context`, `User`, `Infrastructure`, `Timeout`) on a shared `Error` struct would let `ErrorAuto` route by category instead of by type. The per-subsystem `Code` can remain as a finer-grained label.
- **Promote `errors.Is`-compatible predicates for aggregators.** `v1/loader/errors.go` and `v1/rego/rego.go` should expose `IsCode(code string) bool` analogous to `ast.IsError(code, err)` so callers don't have to flatten.
- **Make the server writer's `switch` exhaustive.** Convert `v1/server/writer/writer.go:27-42` to a `map[errorType]httpStatus` table or a typed `RoutingRule` registry so adding a new error type warns at build time if the writer doesn't know about it.
- **Add retry hints to the taxonomy.** `Retryable interface { RetryAfter() time.Duration }` consumed in `v1/plugins/bundle/plugin.go` and `v1/plugins/logs/plugin.go` would let transient infrastructure errors (`HTTPError` on 5xx) be retried automatically while permanent ones (404, 401) short-circuit.
- **Replace `panic` with `sentinel` in WASM SDK.** `internal/wasm/sdk/opa/errors/errors.go:39-46` should return `ErrUnknownCode` from `New` instead of panicking, then have the host surface it. This is friendlier to extensions.
- **Add an observability hook.** `Error.TraceID() string` (or `SpanCtx` field) populated at the point of failure would let traces correlate the error code with the policy/built-in line that produced it. Today the correlation is reconstructed manually from `*ast.Location`.

## Questions / Gaps

- **How is "retry" determined?** OPA does not appear to have a `Retryable` interface in the studied sources. Bundle plugin retries activation (`v1/plugins/bundle/plugin.go:44` `maxActivationRetry=10`) but uses a fixed cap, not a per-error-type retry hint. No clear evidence of a "retry this error" classification — dimension question "can you tell from the error type whether to retry, escalate, or stop?" is only partially answered: `topdown.IsCancel` ⇒ stop, `storage.IsNotFound` ⇒ not-found-escalate, but other error codes do not carry retry semantics.
- **Is there a "policy" vs. "validation" vs. "context" split?** The evidence suggests that the natural OPA axis is subsystem (parse / compile / evaluate / store / serve) rather than the "model/provider/tool/validation/policy/context/user/infrastructure/timeout" axis the dimension prompt lists. The studied code does not name any of those categories explicitly.
- **Are infra/timeout errors typed?** `v1/util/wait.go:27` returns `errors.New("timeout")` as a bare sentinel. `v1/plugins/rest/auth.go:781` checks `errors.Is(err, context.Canceled)` and `errors.Is(err, context.DeadlineExceeded)`. There is no central `opa.TimeoutError` or `opa.ContextError` type. No clear evidence of a typed "context lost" / "deadline exceeded" classification outside of `context.Context` itself.
- **How do custom built-in authors learn the code contract?** `docs/docs/contrib-adding-builtin-functions.md:135` shows `want_error_code: eval_builtin_error` in tests but does not explain the closed-code WASM contract. Embedders who want to raise a custom code will hit the `New()` panic in `internal/wasm/sdk/opa/errors/errors.go:39-46` with no documented extension path.
- **Where is the single source of truth for the error code list?** The codes are scattered across `v1/ast/errors.go:48-66`, `v1/topdown/errors.go:35-60`, `v1/storage/errors.go:11-42`, `internal/wasm/sdk/opa/errors/errors.go:12-30`, `v1/sdk/opa.go:538-548`, `v1/server/types/types.go:118-128`, `v1/plugins/bundle/status.go:21` (`errCode`), `v1/plugins/logs/status/status.go:19` (`errCode`). There is no `errors.go` aggregating these. `docs/docs/errors/index.md` is the only place that lists them together, and it only covers AST/topdown codes (no storage, SDK, server, or plugin codes).
- **Is the taxonomy durable under v0→v1 module split?** The v1 paths (`v1/ast`, `v1/topdown`, `v1/storage`, `v1/sdk`) are the canonical homes; the unversioned mirrors (`ast`, `topdown`, `storage`, `sdk`) thin-import from v1. No clear evidence of a separate, parallel v0 code set — the v1 codes are the codes.

---

Generated by `13.01-error-taxonomy` against `opa`.
