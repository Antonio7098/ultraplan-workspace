# Source Analysis: opa

## Dimension 13.01: Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | opa (Open Policy Agent) |
| Path | `studies/agent-harness-study/sources/opa` |
| Language / Stack | Go (policy engine: Rego parser/compiler, evaluator, HTTP server, plugin framework, WASM SDK) |
| Analyzed | 2026-08-21 |

## Summary

OPA classifies errors with a **layer-local, code-prefixed string taxonomy** rather than a single global error enum. Each subsystem defines its own `Error` struct carrying a `Code string` plus message/location, with code constants prefixed by subsystem: `rego_*` for parse/compile (`v1/ast/errors.go:48-66`), `eval_*` for runtime evaluation (`v1/topdown/errors.go:35-60`), `storage_*` for the store (`v1/storage/errors.go:11-42`), `bundle_error` for bundle download/activation (`v1/plugins/bundle/status.go:20-22`), and SDK-level codes like `opa_undefined_error` (`v1/sdk/opa.go:538-541`) and the WASM SDK's closed set (`internal/wasm/sdk/opa/errors/errors.go:12-30`). A second, coarser taxonomy exists at the REST API boundary (`internal_error`, `evaluation_error`, `unauthorized`, `invalid_parameter`, `invalid_operation`, `resource_not_found`, `resource_conflict`, `undefined_document`, `v1/server/types/types.go:119-128`).

Classification is genuinely load-bearing: a type/code switch in `writer.ErrorAuto` maps error kinds to HTTP statuses and API codes (`v1/server/writer/writer.go:27-42`), `handleBuiltinErr` decides whether a builtin failure becomes `eval_type_error`, `eval_builtin_error`, or nothing (`v1/topdown/builtins.go:182-203`), and `Status.SetError` uses `errors.As` to split bundle failures into AST vs HTTP vs generic for the status API (`v1/plugins/bundle/status.go:65-98`). A special `Halt`/`HaltError` type marks errors that must stop evaluation entirely (`v1/topdown/errors.go:14-24`, `v1/rego/errors.go:5-18`).

The taxonomy is stage-based (parsing → compilation → evaluation) and documented as such in a dedicated errors guide with per-error pages (`docs/docs/errors/index.md:14-32,47-48`). Note: OPA is a policy engine, not an LLM agent harness — the dimension's "model/provider" categories have no direct counterpart; the nearest mappings are "user/authoring" (parse/compile), "tool" (builtins), "infrastructure" (storage, bundle download), "policy decision" (authorization rejections), and "context/timeout" (cancellation errors).

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale:
- **Clear model**: every layer has a typed error with a string `Code`, and codes follow a discoverable `subsystem_noun_error` convention (`v1/ast/errors.go:50-65`, `v1/topdown/errors.go:38-59`, `v1/storage/errors.go:13-41`).
- **Used for handling**: routing decisions demonstrably key off type/code (`v1/server/writer/writer.go:27-42`, `v1/topdown/builtins.go:182-203`, `v1/tester/runner.go:1146`, `v1/debug/debugger.go:339`).
- **Tested**: `errors.Is`/`As` semantics, `Halt` wrapping, and cancellation matching are covered (`v1/topdown/errors_test.go:12-100`; `v1/plugins/bundle/errors_test.go:42-83`).
- **Documented**: a user-facing errors guide organizes errors by stage with fix guidance (`docs/docs/errors/index.md`).
- **Not 8-9** because: the taxonomy is fragmented across six-plus struct definitions with no shared interface or registry; several components emit unclassified `fmt.Errorf` strings (e.g., decision-log mask rules, `v1/plugins/logs/mask.go:69,110,161`); retry decisions in the bundle plugin are count-based, not error-code-based (`v1/plugins/bundle/plugin.go:408`); and one routing quirk maps write conflicts to HTTP 404 (`v1/server/writer/writer.go:31-32`).

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| AST error type + codes | `Error{Code, Message, Location, Details}`; codes `rego_parse_error`, `rego_compile_error`, `rego_type_error`, `rego_unsafe_var_error`, `rego_recursion_error`, `rego_format_error` | `v1/ast/errors.go:82-87`, `v1/ast/errors.go:48-66` |
| AST code matcher | `IsError(code, err)` matches by code string; re-exported in v0 shim | `v1/ast/errors.go:69-74`, `ast/errors.go:34` |
| Eval error type + codes | `Error{Code, Message, Location}`; `eval_internal_error`, `eval_cancel_error`, `eval_conflict_error`, `eval_type_error`, `eval_builtin_error`, `eval_with_merge_error` | `v1/topdown/errors.go:28-33`, `v1/topdown/errors.go:35-60` |
| Halt (stop) signal | `Halt` type: "policy evaluation should stop immediately", with `Unwrap` | `v1/topdown/errors.go:14-24` |
| Code-based matching | `IsCancel` matches `CancelErr`; `Error.Is` enables `errors.Is` with code/message/location matching | `v1/topdown/errors.go:69-82` |
| Storage codes | `storage_internal_error`, `storage_not_found_error`, `storage_write_conflict_error`, `storage_invalid_patch_error`, `storage_invalid_txn_error`, plus not-supported variants | `v1/storage/errors.go:11-42` |
| Storage classifiers | `IsNotFound`, `IsWriteConflictError`, `IsInvalidPatch`, `IsInvalidTransaction` helpers | `v1/storage/errors.go:58-90` |
| REST API error codes | `internal_error`, `evaluation_error`, `unauthorized`, `invalid_parameter`, `invalid_operation`, `resource_not_found`, `resource_conflict`, `undefined_document` | `v1/server/types/types.go:119-128` |
| API error envelope | `ErrorV1{Code, Message, Errors []error}` with `WithASTErrors` to embed detailed AST errors | `v1/server/types/types.go:131-163` |
| User-input error marker | `BadRequestErr` string type + `IsBadRequest` for caller-caused failures | `v1/server/types/types.go:614-637` |
| HTTP dispatch by error kind | `ErrorAuto` type-switch: BadRequest→400/invalid_parameter, write conflict→404/resource_conflict, topdown error→500/internal_error, invalid patch→400, not found→404/resource_not_found, default→500 | `v1/server/writer/writer.go:27-42` |
| Builtin error classification | `handleBuiltinErr`: `BuiltinEmpty`→nil, `*Error`/`Halt`→passthrough, `ErrOperand`→`eval_type_error`, default→`eval_builtin_error` | `v1/topdown/builtins.go:182-203` |
| Stop-vs-collect routing in eval loop | `Halt` unwraps and propagates; other builtin errors are appended to `builtinErrors` and evaluation continues | `v1/topdown/eval.go:2151-2168` |
| Timeout/cancel classification | `http.send` errors: client timeout→generic message; context deadline→`Halt{Error{Code: eval_cancel_error}}` | `v1/topdown/http.go:270-285` |
| Rego-level Halt | `HaltError` + `NewHaltError` for custom builtin implementations; wrapped into `topdown.Halt{topdown.Error{BuiltinErr}}` | `v1/rego/errors.go:5-18`, `v1/rego/rego.go:3027-3053` |
| Bundle error classification | `Error{BundleName, Code, HTTPCode, Message, Err}`; `NewBundleError` extracts `download.HTTPError` status via `errors.As` | `v1/plugins/bundle/errors.go:25-53` |
| Bundle status classification | `Status.SetError` switch on `errors.As`: `ast.Errors`→compile message, `download.HTTPError`→`http_code`, default→generic; code always `bundle_error` | `v1/plugins/bundle/status.go:65-98` |
| Provider HTTP error type | `download.HTTPError{StatusCode}` | `v1/download/download.go:441-447` |
| Decision-log upload error type | `HTTPError{StatusCode}` with distinct message; `SetError` mirrors bundle pattern | `v1/plugins/logs/status/status.go:57-62`, `v1/plugins/logs/status/status.go:32-49` |
| Cancellation routing | test runner stops on `topdown.IsCancel(err)` unless deadline exceeded | `v1/tester/runner.go:1146-1148` |
| Debugger routing by code | debugger checks `topdownErr.Code == topdown.CancelErr` | `v1/debug/debugger.go:339` |
| SDK undefined classification | `opa_undefined_error` + `IsUndefinedErr` | `v1/sdk/opa.go:529-554` |
| WASM SDK closed code set | six codes; `New` panics on unknown code — stricter than the Go layers | `internal/wasm/sdk/opa/errors/errors.go:12-46` |
| Tests for taxonomy semantics | `TestErrorWrapping`: `IsError`, `IsCancel`, `Halt` wrapping, `errors.Is` with code/message/location | `v1/topdown/errors_test.go:12-100` |
| Tests for bundle error unwrap | `errors.As` through `Errors` list and `Error.Err` chain, HTTPCode propagation | `v1/plugins/bundle/errors_test.go:42-91` |
| Documentation of taxonomy | errors guide: stage table (parsing/compilation/evaluation), per-error pages, undefined-by-default semantics, `--strict-builtin-errors` | `docs/docs/errors/index.md:14-32`, `docs/docs/errors/index.md:47-48`, `docs/docs/errors/index.md:127-148` |
| API code documentation | REST API examples show `"code": "invalid_parameter"` / `internal_error` responses | `docs/docs/rest-api.md:2212`, `docs/docs/ocp/api-reference.md:790` |

## Answers to Dimension Questions

**1. Are errors classified by source?**
Yes, by *subsystem and pipeline stage* rather than by a declared source-category enum. Each layer owns a struct with a prefixed string code: `rego_*` for authoring-time errors raised by the parser/compiler (`v1/ast/errors.go:48-66`, created throughout `v1/ast/parser_ext.go:671-760` and `v1/ast/check.go:301-473`), `eval_*` for runtime (`v1/topdown/errors.go:35-60`), `storage_*` (`v1/storage/errors.go:11-42`), and provider-facing `bundle_error` carrying the upstream HTTP status (`v1/plugins/bundle/errors.go:25-45`). Mapping to this dimension's categories: user→`rego_parse_error`/`invalid_parameter` (`v1/server/writer/writer.go:29-30`); validation→`rego_type_error`/`rego_unsafe_var_error`; tool→`eval_builtin_error` (`v1/topdown/builtins.go:196-201`); policy decision→`unauthorized` (`v1/server/authorizer/authorizer.go:155-164`); infrastructure→`storage_*` and `download.HTTPError` (`v1/download/download.go:441-443`); context/timeout→`eval_cancel_error` (`v1/topdown/http.go:276-283`). No "model" or "LLM provider" category exists — out of scope for this codebase.

**2. Is the taxonomy used for handling?**
Yes, in several concrete dispatch points. `writer.ErrorAuto` selects the HTTP status and API code purely from error type/code (`v1/server/writer/writer.go:27-42`); `handleBuiltinErr` chooses between ignore/passthrough/type/builtin classification (`v1/topdown/builtins.go:182-203`); the eval loop uses the `Halt` type to decide stop-immediately vs collect-and-continue (`v1/topdown/eval.go:2162-2168`); `Status.SetError` routes to different status payloads via `errors.As` (`v1/plugins/bundle/status.go:70-97`); the test runner and debugger branch on `topdown.IsCancel` / `CancelErr` (`v1/tester/runner.go:1146`, `v1/debug/debugger.go:339`). However, *retry* decisions are not taxonomy-driven: bundle activation retries up to `maxActivationRetry` (10) attempts regardless of error cause (`v1/plugins/bundle/plugin.go:34-44`, `v1/plugins/bundle/plugin.go:408-434`), and download retries are schedule-based rather than status-code-based.

**3. Are error categories documented?**
Yes for user-facing errors: `docs/docs/errors/index.md:14-32` maintains a stage/category/message table with a page per documented error explaining cause and fix, and frames the three stages explicitly (`docs/docs/errors/index.md:47-48`). It also documents the key semantic tradeoff that builtin failures default to *undefined* rather than errors, with `--strict-builtin-errors` / `--show-builtin-errors` opt-outs (`docs/docs/errors/index.md:127-148`). REST API error codes appear in API examples (`docs/docs/rest-api.md:2212`). Go-level codes are documented only as doc comments on the constants (`v1/topdown/errors.go:37-59`); there is no single reference enumerating all codes across layers.

**4. Can new error types be added without breaking existing handling?**
Largely yes. Codes are plain strings, not iota enums, so new codes are additive; every routing switch has a safe fallback (`ErrorAuto` default→500/internal_error at `v1/server/writer/writer.go:39-40`; `handleBuiltinErr` default→`eval_builtin_error` at `v1/topdown/builtins.go:195-201`; `Status.SetError` default branch at `v1/plugins/bundle/status.go:92-96`). `topdown.Error.Is` matches on code only when set (`v1/topdown/errors.go:74-82`), so unknown codes still flow through wrapping. The exception is the WASM SDK, whose `New` constructor panics on codes outside its fixed set of six (`internal/wasm/sdk/opa/errors/errors.go:39-46`) — a deliberately closed taxonomy. Caveat: because `ErrorAuto` dispatches on type before code, a *new Go error type* will silently fall to the 500 default unless the switch is extended; the taxonomy is extensible for codes but not automatically for new types.

**Guiding question — can you tell from the error type whether to retry, escalate, or stop?**
Partially. *Stop* is explicit: `topdown.Halt` and `rego.HaltError` exist precisely to abort evaluation (`v1/topdown/errors.go:14-18`, `v1/rego/errors.go:3-7`), and `eval_cancel_error` signals cancellation/timeout (`v1/topdown/http.go:276-283`). *Escalate* is expressible at the API boundary (400 user error vs 500 internal via `ErrorAuto`). *Retry* is not derivable from the taxonomy — nothing in the code distinguishes retryable from non-retryable bundle/storage failures; retries are count- or timer-driven (`v1/plugins/bundle/plugin.go:408`), and the status API merely exposes `http_code` for external consumers to decide (`v1/plugins/bundle/status.go:86-90`).

## Architectural Decisions

1. **Layer-local error structs with string codes, not a central error hierarchy.** `ast.Error` (`v1/ast/errors.go:82`), `topdown.Error` (`v1/topdown/errors.go:28`), `storage.Error` (`v1/storage/errors.go:45`), `sdk.Error` (`v1/sdk/opa.go:529`), and WASM `errors.Error` (`internal/wasm/sdk/opa/errors/errors.go:33`) are structurally similar but unrelated types. This keeps packages dependency-free but means cross-layer handling must type-switch per layer.
2. **Two-level taxonomy: detailed internal codes, coarse API codes.** Internal `rego_*`/`eval_*`/`storage_*` codes are embedded as an `Errors []error` list inside the coarse API `ErrorV1` via `WithASTErrors` (`v1/server/types/types.go:156-163`), preserving detail for clients while keeping the wire contract stable.
3. **A dedicated control-flow error type (`Halt`) separate from the taxonomy.** Stop-semantics are encoded in the type, not the code, so any code can be fatal when wrapped in `Halt` (`v1/topdown/errors.go:14-24`, `v1/topdown/eval.go:2163-2165`).
4. **Errors-as-values with `errors.Is/As` support.** `topdown.Error.Is` implements partial matching on code/message/location (`v1/topdown/errors.go:74-82`), and bundle errors implement `Unwrap` chains so `errors.As` can recover `download.HTTPError` (`v1/plugins/bundle/errors.go:51-53`, tested in `v1/plugins/bundle/errors_test.go:57-71`).
5. **Undefined-by-default for builtin failures.** Builtin errors are collected rather than propagated unless strict mode is enabled (`v1/topdown/eval.go:2166-2167`; `v1/server/types/types.go:609-611` for the `strict-builtin-errors` parameter; documented at `docs/docs/errors/index.md:127-148`).

## Notable Patterns

- **Classification helper functions per package**: `IsNotFound`/`IsWriteConflictError`/... (`v1/storage/errors.go:58-90`), `IsCancel`/`IsError` (`v1/topdown/errors.go:63-71`), `IsBadRequest` (`v1/server/types/types.go:633-637`), `IsUndefinedErr` (`v1/sdk/opa.go:550-554`) — callers never string-compare codes directly.
- **Constructor functions that pin codes**: `functionConflictErr`, `completeDocConflictErr`, `objectDocKeyConflictErr`, `mergeConflictErr`, `internalErr` centralize code assignment (`v1/topdown/errors.go:117-163`).
- **`errors.As`-based status enrichment**: both bundle (`v1/plugins/bundle/status.go:70-97`) and decision-log (`v1/plugins/logs/status/status.go:32-49`) status objects classify errors into AST vs HTTP vs generic for observability.
- **Location-carrying errors**: both AST and eval errors attach `Location` for precise reporting (`v1/ast/errors.go:85`, `v1/topdown/errors.go:31`), and `Errors.Sort` orders by location (`v1/ast/errors.go:38-46`).
- **Wrapping to preserve cause**: `topdown.Error.Wrap/Unwrap` (`v1/topdown/errors.go:108-115`) and `rego.HaltError` wrapping into `topdown.Halt` (`v1/rego/rego.go:3036-3041`).

## Tradeoffs

- **Fragmentation vs decoupling**: per-package structs avoid import cycles but duplicate shape and prevent a single `errors.As` target for "any OPA error"; the API layer must know each concrete type (`v1/server/writer/writer.go:28-41`).
- **String codes vs enum**: strings are wire-friendly and additive, but typos compile fine; only the WASM SDK enforces a closed set (via panic at construction, `internal/wasm/sdk/opa/errors/errors.go:43-44`).
- **Undefined-by-default vs fail-fast**: swallowing builtin errors keeps policies resilient to flaky side effects (e.g., `http.send`), but hides real failures unless operators know to enable strict mode (`docs/docs/errors/index.md:145-148`).
- **Stage-based docs vs source-based codes**: the errors guide is organized by stage (`docs/docs/errors/index.md:47-48`), which matches the authoring workflow but does not document infrastructure/provider codes (`storage_*`, `bundle_error`) in the same place.

## Failure Modes / Edge Cases

- **HTTP status/code mismatch**: `ErrorAuto` maps `storage.IsWriteConflictError` to HTTP 404 while labeling it `resource_conflict` (`v1/server/writer/writer.go:31-32`) — a 404 status for a conflict condition can confuse generic HTTP clients.
- **Retry blindness**: bundle activation retries a fixed 10 times regardless of whether the failure is transient or permanent (`v1/plugins/bundle/plugin.go:34-44,408-434`); the taxonomy cannot express "do not retry".
- **Timeout message laundering**: `http.send` client timeouts are deliberately re-worded into a generic "request timed out" message, losing the underlying error detail unless the context was cancelled (`v1/topdown/http.go:271-275`).
- **Silent error collection**: non-Halt builtin errors are appended to `builtinErrors` and evaluation proceeds (`v1/topdown/eval.go:2166-2167`); callers that never inspect the collection see undefined results with no signal.
- **New-type blind spot**: an error type not anticipated by `ErrorAuto`'s switch is reported as `internal_error`/500 even if it is a user fault, until the switch is updated (`v1/server/writer/writer.go:39-40`).
- **Unclassified subsystems**: decision-log mask rule failures use bare `fmt.Errorf` values with no codes (`v1/plugins/logs/mask.go:69,110,161,199`), making programmatic classification impossible there; they are funneled through an `OnRuleError` callback instead (`v1/plugins/logs/mask.go:41,406`).

## Future Considerations

- Introduce a shared interface (e.g., `Code() string`) implemented by all layer error structs so API/status layers can classify without per-type switches, while keeping the per-layer structs.
- Add a retryability dimension (e.g., a `Retryable bool` or `retryable_*` code family) so the bundle/download plugins can make count-free retry decisions; today only `download.HTTPError.StatusCode` gives a hint (`v1/download/download.go:441-443`).
- Publish a machine-readable registry of all codes (the pieces exist as constants in `v1/ast/errors.go:48-66`, `v1/topdown/errors.go:35-60`, `v1/storage/errors.go:11-42`) and cross-link it from `docs/docs/errors/index.md`.
- Reconcile the write-conflict HTTP status (404 vs 409) in `v1/server/writer/writer.go:31-32`.

## Questions / Gaps

- No evidence found of a documented, exhaustive list of REST API `code` values (e.g., in `docs/docs/rest-api.md`); codes appear only in examples (`docs/docs/rest-api.md:2212`). Searched `docs/` for `invalid_parameter`, `evaluation_error`, `error codes`.
- `ast.IsError(code, err)` (`v1/ast/errors.go:69-74`) has no production callers inside this source (grep across `*.go` found only the v0 shim re-export at `ast/errors.go:34`); whether external consumers rely on it could not be determined from this source alone.
- No evidence found of timeout-specific error *codes* distinct from `eval_cancel_error`; timeouts are folded into cancellation (`v1/topdown/http_slow_test.go:70,94` expects `CancelErr` with "timed out" message) or into generic builtin errors.
- Whether `eval_internal_error` is ever surfaced distinctly to clients (vs the coarse `internal_error`) could not be traced end-to-end within the search boundary; `ErrorAuto` wraps any `topdown.Error` as `CodeInternal` (`v1/server/writer/writer.go:33-34`).

---

Generated by `13.01-error-taxonomy` against `opa`.
