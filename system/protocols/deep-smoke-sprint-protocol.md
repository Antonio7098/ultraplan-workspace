# Deep Smoke Sprint Protocol

> Role: sprint-targeted real-system smoke orchestrator.
> Canonical output: `projects/<project>/sprints/<sprint>/smoke.md`.
> Detailed runs and issue evidence remain in the external smoke harness cataloged by `project-index.md`.

Deep smoke runs after sprint review. It is not a substitute for unit tests or `review.md`; it proves runtime-facing claims with executable evidence and links the relevant external harness run and issue records. Every normal smoke run first uses the configured smoke author model to inspect the sprint and target, create or update a durable sprint-specific suite in the external harness, and declare exactly which non-deterministic coverage that suite owns. UltraPlan then discovers and executes that authored suite.

Normal smoke execution must not change product source, product tests, governed sprint inputs, product documentation, or Git state. It may change only manifest-declared harness authoring paths, harness-declared run/evidence paths, sprint smoke flow state, and sprint-root `smoke.md`. UltraPlan records the author run and changed harness paths and rejects out-of-scope mutations.

## Order

```text
execute -> review -> smoke
```

Use smoke for real runtime/process behavior, CLI/TUI output, JSON, persistence, artifacts, locks, cancellation, timeout, auth/model/environment, security, redaction, containment, and mutation claims. Runtime fixes require passing smoke evidence.

## 1. Discover Context

Resolve the target, harness root, versioned harness manifest, `runs/`, and `issues/` from `project-index.md`. Read the governed sprint inputs, execute evidence, current `review.md`, and flow state. Review must be current and `pass`/`pass_with_findings` by default; failed/stale review requires an explicit diagnostic override.

## 2. Author The Sprint Suite

Run the configured smoke author against the governed requirements, reasoning, plan, execute evidence, review, product target, existing deterministic tests, and current harness. Build a coverage matrix for behaviors that normal unit and integration tests do not prove deterministically: real provider/model calls when the product path uses them, actual process and signal behavior, real filesystem/configuration boundaries, real network services, real browser engines, credentials, timing, cancellation, and platform behavior.

Create or update durable tests under the manifest-declared authoring paths. Do not merely rerun normal verification commands or label an empty/cached suite complete. Every complete sprint mapping must enumerate stable required coverage IDs, every selected suite must enumerate at least one test ID, and each coverage ID must be owned by an enumerated test. Provider substitution is forbidden when the product path being tested uses a real provider; offline product paths must not add a provider call just for appearance.

## 3. Read The Machine Harness Contract

The manifest declares schema/protocol version, executable/argument prefix, machine discovery/run commands, authoring paths, capabilities, evidence directories, and environment policy. Discovery supplies levels, suites, enumerated tests and coverage IDs, sprint mappings, prerequisites, evidence schema, and duration/cost class. README prose is documentation, never an executable command. Do not hardcode developer paths, models, or commands.

## 4. Select Scope

Prefer the authored sprint-specific suite, then directly relevant suites, then explicit tests followed by their containing required suite, and full harness only for broad closure. Record sufficiency. Missing, empty, or unassigned required coverage returns blocked and must be fixed in the harness before the sprint can pass smoke.

Show scope, prerequisites, runtime/model, duration/cost class, target/harness paths, allowed mutations, timeout, and evidence destination before confirmation.

## 5. Run Safely

Use explicit executable/argv, contained cwd, allowlisted environment, context timeout/cancellation, bounded output capture, descendant cleanup, and no shell interpolation. Validate run ID, argv/cwd, runtime/model, environment, time/duration, exit code, selected scope/counts, evidence paths/identity, issue IDs, and mutation summary. Harness `runs/` is the detailed source of truth; do not copy raw streams into the sprint.

## 6. Classify And Rerun

Classify product bug, harness bug, model variance, runtime behavior, environment/configuration, missing coverage/evidence, or weak assertion. Relevant issue files stay in harness `issues/` and are linked from `smoke.md`.

A resolution requires a fix outside normal smoke, a passing narrow rerun, a passing containing required suite, and resolved issue evidence. A narrow pass alone cannot establish the final required suite pass.

## 7. Write `smoke.md`

Atomically write sections for Smoke Context, Review Gate, Selected Scope And Rationale, Preconditions And Environment, Run Evidence, Findings, Open Issues, Resolved Issues, Mutation And Safety Check, and Verdict And Next Action.

Verdicts are `pass`, `pass_with_open_issues`, `fail`, `blocked`, or `not_applicable`. Missing prerequisites/coverage never pass; runtime-free scope is `not_applicable`. Failed/cancelled runs do not corrupt the last complete artifact. A successful explicit run may replace older `deep-smoke.md`/`smoke.md`.

## 8. TUI Parity

The TUI exposes the same typed smoke use case: review gate, scope, prerequisites, duration/cost class, confirmation, progress, cancellation, result, issues, evidence links, reruns, recovery, and `smoke.md` preview. It must not call CLI handlers, parse terminal output, or own harness execution.
