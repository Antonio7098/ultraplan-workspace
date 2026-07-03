# Deep Smoke Sprint Protocol

> Role: sprint-targeted real-runtime smoke investigator.
> Input: one completed sprint slug, for example `14-validation-and-diagnostics` or `sprint-14`.
> Reads the project index to discover the smoke harness location.
> Evidence standard follows the smoke harness README.

This protocol deep-smokes a completed sprint against the project's real smoke harness. It is not a substitute for unit tests or sprint review. Its purpose is to prove runtime-facing claims with executable evidence and to file precise issue records when the real system behaves differently than expected.

This is an investigation and reporting protocol only. It must not change product source code, product tests, product documentation, or sprint artifacts except for the deep-smoke report it is asked to produce. The only repository it may change during execution is the smoke harness, and only when adding or correcting investigation coverage, run evidence, issue files, or resolved issue records.

## Purpose

Use this protocol after implementation and review when a sprint touches behavior that needs live end-to-end confidence:

- real runtime or subprocess behavior;
- CLI output, exit codes, JSON envelopes, or diagnostics;
- persisted state, report files, locks, or generated artifacts;
- cancellation, timeout, auth, model, or environment behavior;
- security, redaction, containment, or filesystem mutation claims.

Do not claim an issue is fixed from code inspection alone. Prove it with a passing smoke run that covers the affected path.

---

## 1. Discover Sprint And Smoke Context

### Required identifiers

```text
Project: <project>
Sprint: <sprint-slug-or-number>
Project index: .ultra/projects/<project>/project-index.md
Sprint directory: .ultra/projects/<project>/sprints/<sprint-slug>/
```

### Required inputs

Read the project index first. Use it to discover:

- target implementation directory;
- smoke harness path;
- smoke README path;
- relevant smoke command or suite naming convention;
- evidence directories for `runs/` and `issues/`.

Then read the sprint artifacts:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/sprint-index.md
.ultra/projects/<project>/sprints/<sprint-slug>/requirements.md
.ultra/projects/<project>/sprints/<sprint-slug>/technical-handbook.md
.ultra/projects/<project>/sprints/<sprint-slug>/reasoning.md
.ultra/projects/<project>/sprints/<sprint-slug>/plan.md
.ultra/projects/<project>/sprints/<sprint-slug>/review.md
```

If the input is a shorthand such as `sprint-14`, map it to the matching sprint directory before running anything. If no exact directory exists, list the project sprints and select the directory whose numeric prefix or slug clearly matches the input.

---

## 2. Read The Smoke Harness Contract

Open the smoke harness README discovered from the project index. Treat it as the source of truth for:

- entrypoint and environment variables;
- supported `--level` values and `--tests` filters;
- current baseline and known open issues;
- run evidence format;
- issue file workflow;
- resolution evidence rules.

For `ultraplan-go`, the verified entrypoint is currently:

```bash
OPENCODE_MODEL=minimax-coding-plan/MiniMax-M2.7 node --import tsx/esm src/cli.ts smoke --level <level> --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke
```

Do not assume this command for other projects. Read the project index and smoke README each time.

---

## 3. Select The Deep Smoke Scope

Choose the narrowest smoke scope that fully covers the sprint's delivered surface.

Selection order:

1. Prefer a sprint-specific suite, for example `--level sprint-14`, when the smoke README advertises one for the input sprint.
2. Use the suite names in the smoke README when they map directly to the sprint surface, for example `study-run-loop-deep`, `study-summary-deep`, or `code-extract-deep`.
3. Use `--tests <test-name>` for a known failing or recently fixed path, then rerun the containing `--level`.
4. Use `--level all` only when closing a broad issue set or when multiple suites are affected.

Record why the selected level or test set is sufficient. If the existing harness lacks coverage for the sprint, add or update smoke harness tests before claiming confidence. Do not add or change product tests as part of this protocol.

---

## 4. Run The Smoke

From the smoke harness directory:

1. Verify dependencies and entrypoint are available.
2. Run the selected smoke command.
3. Capture the run ID, exit code, command, cwd, runtime/model, stdout, stderr, and generated run files.
4. Open `runs/<run-id>-summary.md`.
5. Open `runs/<run-id>.json` when a failure needs full evidence.

For every run, use the harness `runs/` files as the source of truth. Do not rely only on terminal output.

---

## 5. Classify Failures

For each failing behavior, classify the cause before fixing anything:

- product bug;
- smoke harness bug;
- model variance;
- OpenCode/runtime behavior;
- environment/auth/model configuration;
- missing evidence or weak assertion.

Ask the concrete investigation questions from the smoke README, including whether the expected command was invoked, whether OpenCode emitted JSON, whether the model wrote the expected file, whether validation is too strict, whether stdout/stderr landed on the expected stream, whether model configuration reached OpenCode, whether child processes survived, whether secrets leaked, and whether filesystem mutation escaped the expected paths.

Do not collapse these into a generic runtime failure.

---

## 6. Investigate, Rerun, And Manage Issues

When a failure is real:

1. Determine whether the problem is in the product, harness, runtime, model, or environment.
2. If the smoke harness is wrong or incomplete, fix only the smoke harness.
3. If the product appears wrong, do not modify product source code. File or update an evidence-backed issue and describe the recommended product change.
4. Rerun the narrow failing test when a harness change or environment correction was made.
5. Rerun the relevant level when the selected smoke scope needs updated evidence.
6. For broad closure, rerun `--level all`.
7. File or update an evidence-backed issue under the harness `issues/` directory when behavior remains open.
8. Move an issue into `issues/resolved/` only after a passing run proves the affected path.

Product fixes are out of scope for this protocol. They may be recommended in the report, but must be performed through a separate implementation task.

An issue claim needs:

- failing command;
- cwd/workspace;
- runtime/model;
- stdout and stderr;
- exit code;
- run ID;
- expected behavior;
- actual behavior;
- reproduction notes;
- generated artifact path or hash when relevant.

A resolution claim needs:

- fix or harness change that resolved the observed behavior, noting whether the product fix happened outside this protocol;
- passing run ID;
- exact smoke command;
- evidence that the relevant test passed;
- no remaining top-level issue file for that behavior.

---

## 7. Produce The Deep Smoke Report

Write the result into the sprint's review or smoke evidence location used by the project. If no project-specific location exists, create:

```text
.ultra/projects/<project>/sprints/<sprint-slug>/deep-smoke.md
```

The report must include:

```markdown
# Deep Smoke: <sprint>

## Scope
- Sprint:
- Implementation directory:
- Smoke harness:
- Selected smoke level/tests:
- Reason scope is sufficient:

## Run Evidence
| Run ID | Command | Result | Summary |
|---|---|---|---|

## Findings
| ID | Severity | Status | Evidence | Required Action |
|---|---|---|---|---|

## Open Issues
- ...

## Resolved Issues
- ...

## Verdict
pass | pass_with_open_issues | fail
```

Use `pass` only when the selected smoke scope passed and no relevant top-level issue remains open. Use `pass_with_open_issues` when the sprint's primary smoke path passed but known or deferred evidence-backed issues remain. Use `fail` when a selected smoke path fails without an accepted deferral.
