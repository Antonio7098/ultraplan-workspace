# Sprint Plan: Documentation and Packaging

> Project: `ultraplan-go`
> Sprint: `15-documentation-and-packaging`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/15-documentation-and-packaging/requirements.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/reasoning.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/sprint-index.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`; area reasoning omitted because `projects/ultraplan-go/sprints/15-documentation-and-packaging/reasoning/` is not present.

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** none; the selected `reasoning/` directory is absent, and `reasoning.md` records architecture reasoning directly from `ARCHITECTURE.md`, `sprint-index.md`, and `technical-handbook.md`.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-06-13
- **Completion Date:** 2026-06-13

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Ship study-side documentation set only | `reasoning.md#decision-1-ship-study-side-documentation-set-only` | Write the user guide, CLI reference, and README update around current study-side workflows only; explicitly exclude target/sprint workflows, SaaS, browser UI, multi-user collaboration, automatic Git mutation, and direct OpenCode ownership claims. |
| Make configuration and runtime boundary documentation authoritative | `reasoning.md#decision-2-make-configuration-and-runtime-boundary-documentation-authoritative` | Write `docs/configuration.md` with supported fields, precedence, redaction, runtime/model/retry/fallback settings, schema-version rejection, and agentwrap/OpenCode mapping. |
| Publish recovery runbook for operational failure modes | `reasoning.md#decision-3-publish-recovery-runbook-for-operational-failure-modes` | Write `docs/recovery.md` for validation failures, missing artifacts, stale tasks, locks, `--force-unlock`, cancellation, retry/fallback metadata, partial completion, and loud atomic write failures. |
| Package four local release binaries from the single CLI entrypoint | `reasoning.md#decision-4-package-four-local-release-binaries-from-the-single-cli-entrypoint` | Build `./cmd/ultraplan` for Linux and macOS on amd64 and arm64 into the required `dist/` filenames; do not publish, sign, notarize, tag, or upload. |
| Generate exact checksums and redacted smoke evidence | `reasoning.md#decision-5-generate-exact-checksums-and-redacted-smoke-evidence` | Generate `dist/checksums.txt` with exactly four SHA-256 entries and create `dist/smoke-evidence.md` with offline gates, build targets, checksum command, redaction notes, and gated OpenCode pass/fail/skip status. |
| Use a release checklist as the publishing gate | `reasoning.md#decision-6-use-a-release-checklist-as-the-publishing-gate` | Write `docs/release-checklist.md` with test, race, build, packaging, checksum, smoke, local `replace`, supported-platform, prompt/version metadata, runtime redaction, stable JSON, recovery, and security review gates. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| REQ-1, Documentation, Workflows | `docs/user-guide.md` covers workspace setup, study init/listing, Markdown source applicability, prompt preview, run, synthesize, run-all, run-loop/status, validate, summary, and code extraction. | Document exists; manual review against implemented help and study-side scope. |
| REQ-2, Configuration, Security, LLM Runtime | `docs/configuration.md` documents supported config fields, precedence, redaction, runtime/model settings, retry/fallback settings, and schema-version rejection. | Document exists; review against `ultraplan.yml`, config/help output, and TRD config requirements. |
| REQ-3, Errors, Workflows, Persistence And Migrations | `docs/recovery.md` gives safe operator recovery procedures without inventing new behavior. | Document exists; review for `--force-unlock`, stale state, validation, cancellation, retry/fallback, missing output, and atomic write guidance. |
| REQ-4, LLM Runtime, Testing | `docs/opencode-smoke.md` documents gated agentwrap/OpenCode smoke prerequisites, commands, expected artifacts, cleanup, and skip path. | Document exists; default tests remain offline; smoke evidence records pass/fail/skip. |
| REQ-5, Testing, Security, CLI Surface | `docs/release-checklist.md` covers release gates and audits, including local `agentwrap` replace and stable JSON/recovery/redaction review. | Document exists; checklist includes all required gates. |
| REQ-6, CLI Surface, Documentation | `docs/cli-reference.md` documents public commands, output modes, and Sprint 14 stable JSON surfaces only. | Document exists; compare with `ultraplan --help`, command help, and tests. |
| REQ-7, Documentation | `README.md` includes install, quickstart, supported scope, documentation index, and explicit deferred target/sprint notice. | README review confirms links and exclusions. |
| REQ-8 to REQ-11, Architecture, Testing | Four release binaries are built from `./cmd/ultraplan` with matching `GOOS`/`GOARCH`. | Required files exist under `dist/`; commands recorded in `dist/smoke-evidence.md`. |
| REQ-12, Observability, Security | `dist/checksums.txt` has exactly one SHA-256 entry for each required binary and no unrelated entries. | Checksum generation command and file review. |
| REQ-13, Testing, Observability, Security | `dist/smoke-evidence.md` records exact offline commands, results, build targets, checksum command, and gated OpenCode result or skip reason. | Evidence file exists and passes redaction review. |
| AC-1 to AC-17 | Acceptance criteria in `requirements.md` define release readiness. | Verification commands, artifact inspection, and docs/security/scope reviews. |
| Architecture | Preserve single CLI entrypoint and module ownership; docs must not create public Go API or runtime ownership claims. | Build from `./cmd/ultraplan`; docs review against `ARCHITECTURE.md`. |
| Security | No secrets, provider tokens, full prompts, full reports, full sensitive env dumps, or raw unsafe payloads in docs/evidence. | Manual redaction review of docs and `dist/smoke-evidence.md`. |

## Tasks

- [x] **Task 1: Inspect Implemented Public Surface**
  > Executes: Decision 1, Decision 2, Decision 3, Decision 6, REQ-1 to REQ-7, AC-8 to AC-17
  - [x] Run or inspect `ultraplan --help` and command-specific help for the current public study-side command surface.
  - [x] Inspect existing config examples/tests to confirm supported `ultraplan.yml` fields, precedence behavior, redaction behavior, runtime/model/retry/fallback settings, and schema-version rejection diagnostics.
  - [x] Inspect Sprint 14 implementation/tests for stable JSON surfaces and identify any debug or undocumented JSON that must not be labeled stable.
  - [x] Inspect current README and docs to preserve accurate existing wording and avoid contradicting current behavior.
  - [x] Record any insufficient evidence as an open question in the relevant document rather than guessing.

- [x] **Task 2: Write User-Facing Study Documentation**
  > Executes: Decision 1, REQ-1, REQ-6, REQ-7, AC-8, AC-12, AC-13, AC-14
  - [x] Create `docs/user-guide.md` with an end-to-end study-side workflow from workspace setup through study init, listing, Markdown applicability, prompt preview, single run, synthesize, run-all, run-loop/status, validate, summary, and code extraction.
  - [x] Create `docs/cli-reference.md` with current public commands, flags/output modes, exit/diagnostic expectations where implemented, and stable Sprint 14 JSON surfaces only.
  - [x] Update `README.md` with install/build basics, quickstart, supported scope, documentation index, and explicit deferred target/sprint, hosted SaaS, browser UI, multi-user collaboration, and automatic Git mutation notices.
  - [x] Review examples for consistency with current help output and avoid unsupported provider-agnostic real-runtime guarantees.

- [x] **Task 3: Write Operator And Runtime Documentation**
  > Executes: Decision 2, Decision 3, Decision 5, Decision 6, REQ-2 to REQ-5, AC-7, AC-9 to AC-11, AC-15 to AC-17
  - [x] Create `docs/configuration.md` with supported fields, precedence order of built-in defaults, workspace config, environment variables, and command flags.
  - [x] Document redaction, runtime/model settings, agentwrap/OpenCode mapping, retry/fallback settings, unsupported fields, and schema-version rejection/diagnosis behavior.
  - [x] Create `docs/recovery.md` covering validation repair workflow, missing artifacts, stale running tasks, lock diagnostics, `--force-unlock` semantics and risk, cancellation recovery, partial completion, retry/fallback metadata, and atomic write failure behavior.
  - [x] Create `docs/opencode-smoke.md` with prerequisites, required environment/config, gated commands, expected artifacts, cleanup notes, and how to keep real OpenCode smoke outside normal tests.
  - [x] Create `docs/release-checklist.md` with offline gates, race tests, build, packaging, checksums, smoke evidence, local `replace github.com/Antonio7098/agentwrap => ../agentwrap` audit, platform notes, prompt/version metadata review, runtime metadata redaction, stable JSON review, recovery review, and security review.

- [x] **Task 4: Run Offline Release Gates**
  > Executes: Decision 4, Decision 5, REQ-8 to REQ-13, AC-1 to AC-6
  - [x] From `/home/antonioborgerees/coding/ultraplan-go`, run `go test ./...` and record exact command, working directory, result, and any relevant safe output summary.
  - [x] From `/home/antonioborgerees/coding/ultraplan-go`, run `go test -race ./...` and record exact command, working directory, result, and any relevant safe output summary.
  - [x] From `/home/antonioborgerees/coding/ultraplan-go`, run `go build ./cmd/ultraplan` and record the result.
  - [x] Treat failures as release blockers unless explicitly triaged outside this sprint plan.

- [x] **Task 5: Build Local Release Artifacts**
  > Executes: Decision 4, REQ-8 to REQ-11, AC-3, AC-4
  - [x] Ensure `/home/antonioborgerees/coding/ultraplan-go/dist/` exists for release outputs.
  - [x] Build `GOOS=linux GOARCH=amd64 go build -o dist/ultraplan-linux-amd64 ./cmd/ultraplan`.
  - [x] Build `GOOS=linux GOARCH=arm64 go build -o dist/ultraplan-linux-arm64 ./cmd/ultraplan`.
  - [x] Build `GOOS=darwin GOARCH=amd64 go build -o dist/ultraplan-darwin-amd64 ./cmd/ultraplan`.
  - [x] Build `GOOS=darwin GOARCH=arm64 go build -o dist/ultraplan-darwin-arm64 ./cmd/ultraplan`.
  - [x] Inspect `dist/` to confirm all four required binaries exist and no naming mismatch occurred.

- [x] **Task 6: Generate Checksums And Smoke Evidence**
  > Executes: Decision 5, REQ-12, REQ-13, AC-5, AC-6, AC-7
  - [x] Generate `dist/checksums.txt` with exactly the four required binary filenames and SHA-256 hashes.
  - [x] Verify `dist/checksums.txt` contains no entries for unrelated files such as `smoke-evidence.md`.
  - [x] Create `dist/smoke-evidence.md` with exact offline verification commands/results, build target commands, checksum generation command, working directory, timestamp/date, and redaction notes.
  - [x] Run the gated OpenCode smoke procedure only when prerequisites are available; otherwise record an explicit skip reason naming the missing OpenCode/provider/network/config prerequisite without dumping secrets or full environment.

- [x] **Task 7: Final Scope, Security, And Review Preparation**
  > Executes: Decision 1 to Decision 6, AC-12 to AC-17
  - [x] Review README and docs for absent target/sprint command claims, hosted service claims, browser UI claims, multi-user claims, direct OpenCode supervision claims, unsupported provider-agnostic real-runtime guarantees, and automatic Git mutation claims.
  - [x] Review docs and `dist/smoke-evidence.md` for secrets, provider tokens, full prompts, full report bodies, full sensitive environment variables, and raw unsafe runtime payloads.
  - [x] Review `docs/cli-reference.md` against current help and Sprint 14 stable JSON tests/schemas.
  - [x] Review `docs/release-checklist.md` for dependency provenance, local `agentwrap` replace disposition, platform limitations, prompt/version metadata, runtime metadata redaction, stable JSON, and recovery gates.
  - [x] Ensure `review.md` can run the Architecture Review and Sprint Review protocols without guessing intent.

## Evidence Checklist

- [x] `go test ./...` passes without OpenCode, provider credentials, network access, or real subprocess smoke fixtures.
- [x] `go test -race ./...` passes.
- [x] `go build ./cmd/ultraplan` passes.
- [x] Four required binaries exist under `dist/` and match their `GOOS`/`GOARCH` filenames.
- [x] `dist/checksums.txt` contains exactly four SHA-256 entries for required binaries.
- [x] `dist/smoke-evidence.md` records exact offline commands/results, build targets, checksum command, and gated OpenCode pass/fail/skip status.
- [x] Required docs and README update exist at the paths named in `requirements.md`.
- [x] Docs cover only current study-side capabilities and explicitly state deferred scope.
- [x] Stable JSON documentation is limited to compatibility-sensitive Sprint 14 surfaces.
- [x] Release checklist includes local `agentwrap` replace audit and metadata/redaction/recovery/stable JSON reviews.
- [x] Redaction/security review finds no secrets, full prompts, full report bodies, full sensitive env dumps, provider tokens, or raw unsafe runtime payloads.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Architecture Review and Sprint Review inputs are ready.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline tests | `go test ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go` without OpenCode, credentials, network, or real subprocess smoke fixtures. |
| Race tests | `go test -race ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| CLI build | `go build ./cmd/ultraplan` | Passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| Linux amd64 package | `GOOS=linux GOARCH=amd64 go build -o dist/ultraplan-linux-amd64 ./cmd/ultraplan` | Produces `dist/ultraplan-linux-amd64`. |
| Linux arm64 package | `GOOS=linux GOARCH=arm64 go build -o dist/ultraplan-linux-arm64 ./cmd/ultraplan` | Produces `dist/ultraplan-linux-arm64`. |
| macOS amd64 package | `GOOS=darwin GOARCH=amd64 go build -o dist/ultraplan-darwin-amd64 ./cmd/ultraplan` | Produces `dist/ultraplan-darwin-amd64`. |
| macOS arm64 package | `GOOS=darwin GOARCH=arm64 go build -o dist/ultraplan-darwin-arm64 ./cmd/ultraplan` | Produces `dist/ultraplan-darwin-arm64`. |
| Checksums | `sha256sum dist/ultraplan-linux-amd64 dist/ultraplan-linux-arm64 dist/ultraplan-darwin-amd64 dist/ultraplan-darwin-arm64 > dist/checksums.txt` | `dist/checksums.txt` has exactly four entries and no unrelated files. |
| Help comparison | `./ultraplan --help` and command-specific help, or equivalent built binary help commands | `docs/cli-reference.md`, `docs/user-guide.md`, and README examples match implemented public commands. |
| Gated OpenCode smoke | Commands from `docs/opencode-smoke.md` | Records pass/fail when prerequisites exist, or explicit skip reason when unavailable. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Current CLI help or stable JSON behavior may differ from assumptions. | `reasoning.md` assumptions | Inspect help/tests before final docs; document open questions rather than guessing. | closed |
| Real OpenCode/provider environment may be unavailable. | `reasoning.md` risks | Keep smoke gated and record explicit missing prerequisite skip reason in `dist/smoke-evidence.md`. | mitigated: existing smoke harness evidence reviewed; no new full smoke run after packaging-only dependency update |
| Local `agentwrap` replace directive may remain. | `reasoning.md` risks and `technical-handbook.md` warning | Include release checklist audit and disposition before publishing artifacts. | closed: `agentwrap` fix pushed at `bc655a2`; UltraPlan uses non-local pseudo-version `v0.0.0-20260613122459-bc655a256a5f` |
| Documentation may accidentally claim deferred target/sprint workflows. | PRD/TRD and `reasoning.md` risks | Review README/docs against deferred scope and excluded context. | closed |
| Smoke evidence could leak secrets or unsafe payloads. | Security contract and `reasoning.md` risks | Record commands/results/summaries only; perform redaction review. | closed |
| Manual packaging commands may target the wrong platform or directory. | `reasoning.md` risks | Record exact working directory, commands, `GOOS`/`GOARCH`, and inspect artifact names. | closed |
| Race tests may expose existing implementation defects. | `reasoning.md` risks | Treat failures as release blockers or require explicit triage before release. | closed |
| macOS artifacts are cross-compiled but unsigned and unnotarized. | Sprint non-goals and `reasoning.md` assumptions | State local artifact scope and no signing/notarization in README/checklist. | mitigated |
| Darwin cross-compilation fails in pinned `github.com/Antonio7098/agentwrap/opencode`. | Execution evidence | Fixed and pushed `agentwrap` build tags, upgraded UltraPlan to the pushed pseudo-version, rebuilt Darwin artifacts, and generated four-entry checksums. | closed |

## Review Inputs

Review should use:

- `requirements.md`
- `sprint-index.md`
- `technical-handbook.md`
- `reasoning.md`
- this `plan.md`
- `docs/PRD.md`
- `docs/TRD.md`
- `docs/ARCHITECTURE.md`
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- `dist/smoke-evidence.md`
- `dist/checksums.txt`
- required release binaries under `dist/`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-12 / planning | Created sprint plan from `reasoning.md` without implementing code. | Plan carries decisions, evidence, risks, assumptions, commands, and review inputs forward. |
| 2026-06-13 / implementation | Added README release update and docs: `docs/user-guide.md`, `docs/cli-reference.md`, `docs/configuration.md`, `docs/recovery.md`, `docs/opencode-smoke.md`, `docs/release-checklist.md`. | Docs reviewed against current help/config/status/code surfaces and study-side scope. |
| 2026-06-13 / verification | Ran `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. | All three commands passed. |
| 2026-06-13 / packaging | Built Linux amd64 and Linux arm64 binaries. Darwin amd64 and Darwin arm64 builds failed in pinned `github.com/Antonio7098/agentwrap/opencode` because `syscall.SysProcAttr.Pdeathsig` is Linux-specific. | Release blocked; `dist/smoke-evidence.md` records failure. `dist/checksums.txt` intentionally not generated because exactly four binaries are required. |
| 2026-06-13 / dependency fix | Patched `/home/antonioborgerees/coding/agentwrap` to split Linux-only `Pdeathsig` process setup from Unix process-group setup, committed `bc655a2`, and pushed `origin/main`. | `agentwrap` native tests passed; Darwin opencode compile checks passed via `go test -c`. |
| 2026-06-13 / final packaging | Upgraded UltraPlan to non-local `github.com/Antonio7098/agentwrap v0.0.0-20260613122459-bc655a256a5f`, removed local replace, re-ran tests/race/native build/all four package builds/checksum generation, and reviewed existing smoke harness evidence. | Four binaries exist under `dist/`; `dist/checksums.txt` has exactly four entries; smoke evidence references `/home/antonioborgerees/coding/ultraplan-go-smoke` latest full and Sprint 14 runs. |

## Completion Criteria

- [x] All documentation outputs from `requirements.md` exist and are reviewed against current behavior.
- [x] README links every new document and states deferred target/sprint and non-release publication scope.
- [x] All four release binaries exist under `dist/` with correct target names.
- [x] `dist/checksums.txt` contains exactly the four required binary hashes.
- [x] `dist/smoke-evidence.md` records offline gates, packaging, checksums, and gated OpenCode pass/fail/skip status with redaction notes.
- [x] Verification commands were run or failures/deferrals are explicitly documented.
- [x] Evidence satisfies `reasoning.md` Expected Evidence and `requirements.md` acceptance criteria.
- [x] Architecture Review and Sprint Review can evaluate conformance without guessing intent.
