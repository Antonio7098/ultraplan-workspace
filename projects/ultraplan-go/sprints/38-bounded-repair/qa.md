# QA

Project: `ultraplan-go`
Sprint: `38-bounded-repair`
Input fingerprint: `7777777777777777777777777777777777777777777777777777777777777777`
Attempt: `qa-v1-attempt-c496c705da66f4b04235f119`
Assessment: `pass_with_findings`

## Evidence

Accepted: `1`
Rejected: `0`
- `qa-v2-evidence-85ab3151ab9b9aea2d3fce18` fail, reason `compile_failure`, contained `true`, cleanup `true`

## Promoted issues
- `qa-v2-issue-f1cc3019f6400a0cfa998cb1` [medium] Sentinel compile failure at `internal/sprint/repair_dogfood_sentinel.go`, evidence `qa-v2-evidence-85ab3151ab9b9aea2d3fce18`, regression candidate `true`

## Smoke evidence

Verdict: `pass`
Run: `dogfood-containing-smoke`

## Next action

Repair the promoted sentinel issue.
