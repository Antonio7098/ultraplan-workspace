# Sprint Requirements: Validation Command, Diagnostics, and JSON Stability

> Project: `ultraplan-go`
> Sprint: `14-validation-and-diagnostics`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Make UltraPlan's study validation, health, and run status surfaces inspectable, automatable, and safe through a study validation command, stable JSON output shapes, and redacted diagnostics.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Study validation command | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Adds thin CLI wiring for `ultraplan study <study> validate`, including `--json` support and stable exit-code mapping. |
| Study validation service | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_command.go` | Implements study-owned validation orchestration for study structure, sources, dimensions, expected reports, final reports, summary, run state, and applicability-aware missing/inapplicable artifacts. |
| Validation result model | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Adds or extends typed validation summary/result structures with check names, severity, paths, status, guidance, safe diagnostics, and schema/version fields for JSON output. |
| Status JSON output | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Extends `ultraplan study <study> status` with `--json` using a documented stable machine-readable shape. |
| Health JSON output | `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go` | Extends `ultraplan health --json` to emit a stable, redacted machine-readable health result. |
| JSON rendering helpers | `/home/antonioborgerees/coding/ultraplan-go/internal/app/json_output.go` | Centralizes JSON envelope rendering for command name, workspace, status, timestamps, schema version, and result payload without ANSI text. |
| Diagnostics redaction updates | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/redaction.go` | Ensures all new validation, health, and status diagnostics use the existing redaction policy before text or JSON rendering. |
| Study validation tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_command_test.go` | Covers valid study validation, missing artifacts, invalid reports, malformed frontmatter, inapplicable Markdown pairs, run-state diagnostics, and secret-free diagnostics. |
| Command JSON tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_validate_commands_test.go` | Covers validate text/JSON output, exit codes, checkable failures, redaction, and no runtime invocation. |
| Status JSON tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` | Extends status tests to cover stable JSON output for counts, task states, lock diagnostics, run metadata summaries, and validation summaries. |
| Health JSON tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands_test.go` | Covers `ultraplan health --json` success/failure, runtime-health diagnostics, schema fields, and redaction. |
| Updated command help | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Documents `study <study> validate`, `study <study> validate --json`, and `study <study> status --json` in help text. |
| Sprint review evidence | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/review.md` | Records implementation evidence, deviations, tests, build status, and review conclusions after execution. |

## Acceptance Criteria

- [ ] `ultraplan study <study> validate` validates study structure, source discovery, dimension files, expected per-source reports, final reports, `summary.csv`, and `studies/<study>/.ultraplan/run-state.json` when present.
- [ ] Validation respects Markdown source applicability: inapplicable source/dimension pairs are reported as inapplicable or skipped, not as missing failures.
- [ ] Validation failures include check name, severity, artifact path when applicable, expected condition, observed safe detail, and corrective guidance.
- [ ] Validation returns exit code `0` when all required checks pass and exit code `5` for validation failures.
- [ ] Validation does not invoke agentwrap, OpenCode, network calls, source cloning, runtime execution, synthesis, repair, or prompt generation.
- [ ] `ultraplan study <study> validate --json` emits valid JSON with a stable envelope containing schema version, command, workspace, status, generated timestamp, and a result payload.
- [ ] `ultraplan health --json` emits valid JSON with stable health-check IDs, statuses, safe diagnostics, and runtime/config/workspace health summaries.
- [ ] `ultraplan study <study> status --json` emits valid JSON with stable run counts, task summaries, lock diagnostics, validation summaries, retry/fallback metadata where present, and unknown usage/cost represented as unknown rather than zero.
- [ ] JSON output contains no ANSI formatting and writes diagnostics as JSON fields rather than mixed human text.
- [ ] Text output remains deterministic, concise, and human-readable for validation, health, and status commands.
- [ ] All new JSON and text outputs redact API keys, secret config values, sensitive environment values, prompts, embedded Markdown document bodies, unsafe raw payload bytes, and full runtime stderr unless explicitly safe/truncated.
- [ ] Stable JSON schemas are versioned and covered by command-level tests that parse JSON and assert field names for success and failure cases.
- [ ] Health, status, and validation outputs include actionable next steps or guidance for failed checks without exposing secrets.
- [ ] Existing command behavior remains backward-compatible for text mode unless this sprint explicitly changes it and tests document the new behavior.
- [ ] The implementation stays module-owned: study validation behavior remains in `internal/study`, command parsing/rendering remains in `internal/app`, and platform runtime remains product-agnostic.
- [ ] The implementation introduces no global `internal/validation`, `internal/reports`, `internal/diagnostics`, `internal/status`, or workflow-engine package.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Implementing target, sprint planning, sprint execution, or sprint validation commands.
- Adding a browser UI, TUI dashboard, hosted service, metrics exporter, alerting integration, or local API server.
- Adding new runtime adapters or changing OpenCode process supervision outside the existing agentwrap-backed platform runtime boundary.
- Running analyses, synthesis, repairs, retries, fallback execution, or source cloning from the validation command.
- Implementing schema migrations beyond clear rejection or diagnostics for unsupported durable schema versions.
- Making real OpenCode/provider smoke tests mandatory for the default test suite.
- Changing report templates, prompt semantics, summary scoring rules, or code-reference extraction behavior except where required to expose validation diagnostics.
- Persisting unsafe raw runtime payload bytes or full native stderr in validation/status/health output.
- Introducing a general diagnostics framework unrelated to the three selected public surfaces: `validate`, `health --json`, and `status --json`.

## Constraints

- Study validation behavior must be owned by `internal/study`; `internal/app` may parse flags, call study services, render text/JSON, and map exit codes only.
- `internal/platform/runtime` must not import or know about studies, dimensions, sources, report validation, summaries, status, or validation-command semantics.
- JSON output must use a stable, versioned envelope and must not include ANSI escape sequences, progress text, or human-only formatting.
- Validation/status/health JSON must preserve unknown values as unknown/null/explicit known flags; unknown usage or cost must never be converted to zero.
- All user-visible diagnostics must pass through the established redaction policy before rendering.
- Paths in generated output should prefer workspace-relative paths; absolute paths may appear only when needed for local diagnostics and must not expose secrets.
- Durable state must be read strictly and rejected or diagnosed on unsupported schema versions; this sprint must not silently migrate state.
- Validation must be bounded to known workspace/study artifacts and must not recursively scan source repositories except for behavior already required by existing validators.
- Tests must be deterministic, offline, and fake-first; default verification must not require OpenCode, provider credentials, network access, or long sleeps.
- The sprint must carry forward prior review decisions: preserve cause chains, keep runtime success insufficient without artifact validation, keep status runtime-free, keep diagnostics redacted, and address the previously deferred stable JSON status surface.
- The implementation must not modify unrelated dirty files or revert user changes in the target repository.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 05: Markdown source discovery and applicability | Applicability-aware validation and status | Validation must distinguish missing from inapplicable Markdown document source/dimension pairs. |
| Sprint 06: Report validation and rating parsing | Reusable report checks | The validation command must reuse existing per-source/final report validators instead of duplicating report rules. |
| Sprint 08: Run state persistence and status | Run-state validation and status JSON | Status JSON and validation diagnostics depend on versioned `run-state.json`, task summaries, and resume validation behavior. |
| Sprint 09: Runtime integration and health | Health JSON and safe runtime diagnostics | Health output must use the existing agentwrap-backed runtime health boundary and redacted runtime diagnostics. |
| Sprint 10: Single analysis and synthesis | Execution validation semantics | Validation must preserve the rule that expected artifacts must exist and pass validation after runtime work. |
| Sprint 11: Run-all batch execution | Batch output validation and summary diagnostics | Validation must understand batch-produced per-source reports, synthesis outputs, and `summary.csv`. |
| Sprint 12: Durable run-loop | Run metadata, locks, retry/fallback/cancellation status | Status JSON must expose safe durable run-loop metadata and lock diagnostics without invoking runtime work. |
| Sprint 13: Summary generation and code-reference extraction | Summary and citation diagnostics | Validation should report summary presence/shape and may surface existing unresolved-reference diagnostics without changing extractor scope. |
| PRD/TRD JSON and diagnostics requirements | Stable automation surfaces | PRD requires machine-readable inspect/status/validate output; TRD requires stable JSON mode and safe diagnostics. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Study validation command exists and is scoped | Run `ultraplan study <study> validate --help`; inspect `internal/app` thin wiring and `internal/study` ownership. |
| Validation behavior is objective and actionable | Review `validation_command_test.go` fixtures for valid, missing, invalid, malformed, and inapplicable cases. |
| JSON schemas are stable and parseable | Command tests parse `--json` output and assert envelope fields, schema version, result fields, statuses, and absence of ANSI/text leakage. |
| Redaction is enforced | Tests seed secret-like config/env/error values and assert they are absent from validation, health, and status text/JSON output. |
| Runtime-free validation and status | Tests use fakes or nil runtime seams and assert validation/status do not invoke agentwrap/OpenCode/network/subprocess paths. |
| Health JSON uses runtime boundary correctly | Health command tests verify runtime health failures are reported through existing platform runtime abstractions with safe diagnostics. |
| Architecture boundaries remain intact | Review imports and file placement; confirm no product imports under `internal/platform/runtime` and no new global validation/diagnostics packages. |
| Durable state compatibility is safe | Tests cover missing, malformed, and unsupported `run-state.json` with clear diagnostics and no silent migration. |
| Prior follow-ups are addressed where in scope | Confirm stable `status --json`, redacted diagnostics, run metadata summaries, and validation diagnostics are present; record out-of-scope follow-ups explicitly. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
