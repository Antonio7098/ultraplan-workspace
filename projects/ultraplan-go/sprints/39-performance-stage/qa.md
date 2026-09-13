# QA

Project: `ultraplan-go`
Sprint: `39-performance-stage`
Input fingerprint: `c0af02ba07d3924beeeebdecfecb684cde538f290672be7f05a5e423164f6d21`
Attempt: `qa-v1-attempt-68b6ad13df98a7d5a96df5d6`
Assessment: `pass_with_findings`

## Evidence

Operator retirement `qa-v2-retirement-8fa5f90b05bdeb7def7f0a0b`, request `qa-v2-request-13cc278803c32eaa2dbf06b2`, operator `codex`: Retained reproducer ends in an unconditional testing.T.Fatalf, so it cannot establish a repairable product defect..

- Retired `qa-v2-evidence-bd8fd861a8b2d54eb6480aba`, bundle `qa-v2-test-98a099d74bdc9e0e2061af79`: `unconditional_predicted_failure`. Original artifacts retained.

Operator retirement `qa-v2-retirement-3e83deaa03d33d25579dcb03`, request `qa-v2-request-b6170cb8d2c9b7efe9eb62a1`, operator `antonioborgerees`: Frozen test asserted acceptance while the promoted issue claimed missing rejection; the current implementation correctly rejected the zero-target passing record..

- Retired `qa-v2-evidence-267a2f024601a2e21585e870`, bundle `qa-v2-test-3bee0b7a9d3ae9994b9c28f4`: `unconditional_predicted_failure`. Original artifacts retained.
- Retired `qa-v2-evidence-50ac97b9100a4c15df0e01fc`, bundle `qa-v2-test-82b0c1d1cd30bf494317de77`: `unconditional_predicted_failure`. Original artifacts retained.
- Retired `qa-v2-evidence-59eccf8b0f94fef2945ac94c`, bundle `qa-v2-test-dc7bd57b7dd74e4b104afda0`: `unconditional_predicted_failure`. Original artifacts retained.
- Retired `qa-v2-evidence-97ab420cdfb8e8ffbd55d4ae`, bundle `qa-v2-test-ed61cdd33b508f5e21f741b0`: `unconditional_predicted_failure`. Original artifacts retained.
- Retired `qa-v2-evidence-9dd64f43be091f9c54b6fb1b`, bundle `qa-v2-test-0d0fa2ce66a25e252c2149e5`: `unconditional_predicted_failure`. Original artifacts retained.

Operator retirement `qa-v2-retirement-d957eeb92a8b85e58ded6895`, request `qa-v2-request-1cefed6a599adae46ba55ba6`, operator `antonioborgerees`: Frozen test expected the omitted evidence limit to remain zero and return HTTP 500, but current correct behavior applies limit 50 and returns HTTP 200; the test also ends in an unconditional predicted failure..

- Retired `qa-v2-evidence-67e0793d19fecc71197e3a0c`, bundle `qa-v2-test-da0b697c94f4dcc046d7d411`: `unconditional_predicted_failure`. Original artifacts retained.

Operator retirement `qa-v2-retirement-e38ffc5facfe83beb1dae23c`, request `qa-v2-request-b2ae680281cd729d3c90109f`, operator `antonioborgerees`: Governed repair repair-v1-run-ceb4a33e924baa707503fbe7 applied the fix and passed its exact reproducer and complete ladder; the pre-repair failing evidence must not keep the repaired issue in the retained queue..

- Retired `qa-v2-evidence-8f00144753ce57cb4d91beac`, bundle `qa-v2-test-54a2569f5e519e1993d77c95`: `operator_retirement`. Original artifacts retained.

Operator retirement `qa-v2-retirement-8eaa6ab3f18f8c78235b2125`, request `qa-v2-request-07d95cd845e0108c780842e4`, operator `antonioborgerees`: Frozen test runs a private two-iteration fake loop and deliberately calls its own propose method twice; it never exercises product cleanup behavior, so no governed product repair can satisfy it..

- Retired `qa-v2-evidence-372461ad324508ca6636dcb7`, bundle `qa-v2-test-03d557c2068f466cd567c3ea`: `operator_retirement`. Original artifacts retained.

Operator retirement `qa-v2-retirement-320113ce85f8a9ff9631a652`, request `qa-v2-request-ef633898930720f83a1a7704`, operator `antonioborgerees`: Governed repair repair-v1-run-d00241da2c8a025836383c73 already applied exact mode parsing and passed the immutable reproducer plus complete ladder; the retained pre-repair failure must not be selected again..

- Retired `qa-v2-evidence-408040e1837c91fedfe288e1`, bundle `qa-v2-test-313eb1ef83bfc4588a4583fd`: `operator_retirement`. Original artifacts retained.

Operator retirement `qa-v2-retirement-dbb55139da06e58f934fff2b`, request `qa-v2-request-96ab4d195642f558fdc3b798`, operator `antonioborgerees`: Governed repair repair-v1-run-c5e8e2a228dd383587fe8f7c applied strict JSON suffix rejection and passed its exact reproducer and complete ladder; retire the pre-repair failure from the retained queue..

- Retired `qa-v2-evidence-dbadb1804ad4ea0b7a534469`, bundle `qa-v2-test-d2b77ae8796965482151be3d`: `operator_retirement`. Original artifacts retained.

Operator retirement `qa-v2-retirement-10809cb75bc5062838a55028`, request `qa-v2-request-717ddd9a741f8cee749d8094`, operator `antonioborgerees`: Governed repair repair-v1-run-8a952c78b856d67f5490fedd fixed stdout failure propagation and passed its exact reproducers and complete ladder; retire this pre-repair evidence from reselection..

- Retired `qa-v2-evidence-03a7542f2f3f6cfb3f0c46f1`, bundle `qa-v2-test-30aa4e77884683ef74f1e437`: `operator_retirement`. Original artifacts retained.

Operator retirement `qa-v2-retirement-5a88167249fcf462e92ad018`, request `qa-v2-request-ba638561ab2878d7bc49d841`, operator `antonioborgerees`: Governed repair repair-v1-run-32cc7834f028722009a6b63a fixed report-only flow-state validity and passed its exact reproducer and complete ladder; retire the pre-repair failure from the retained queue..

- Retired `qa-v2-evidence-53b03de369a9bfd58bd53512`, bundle `qa-v2-test-69d2537e17f11e15ae7fd2e2`: `operator_retirement`. Original artifacts retained.

Accepted: `32`
Rejected: `2`
- `qa-v2-evidence-031915bea0e9fd348e0d2dda` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-03a7542f2f3f6cfb3f0c46f1` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-06c244731175c9105f0a0129` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-08bd0589eafbe38eae23ed4d` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-09721d5555a508f0830d8499` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-0e382bbab496e8efd1805043` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-0e4e3dc356711714be780b3d` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-1ba7d9f21d6b034d48920988` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-25d3e98f06890f16ff1e25d0` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-267a2f024601a2e21585e870` inconclusive, reason `failure_signature_mismatch`, contained `true`, cleanup `true`
- `qa-v2-evidence-29f3be0bbb947759a4129dfa` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-33e46626c6f545e2ddd8da96` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-372461ad324508ca6636dcb7` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-3dea7afa8584d44f529789a5` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-408040e1837c91fedfe288e1` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-4a56017ac8e5122050ab565f` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-4abc65866e42d7589a32b6a0` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-4f3c506b07d69fb9942415b8` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-50ac97b9100a4c15df0e01fc` inconclusive, reason `unrelated_compile_or_panic_failure`, contained `true`, cleanup `true`
- `qa-v2-evidence-53b03de369a9bfd58bd53512` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-59eccf8b0f94fef2945ac94c` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-5b45c8ec395971b8f54919cd` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-66e2bb4d026636f6ccd349d2` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-67e0793d19fecc71197e3a0c` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-6f0b3c6e84bd9765b5f6126d` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-8f00144753ce57cb4d91beac` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-97ab420cdfb8e8ffbd55d4ae` inconclusive, reason `failure_signature_mismatch`, contained `true`, cleanup `true`
- `qa-v2-evidence-9dd64f43be091f9c54b6fb1b` inconclusive, reason `unrelated_compile_or_panic_failure`, contained `true`, cleanup `true`
- `qa-v2-evidence-a2fdc6efba894a45bb4c2875` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-b4399611906c26c9c19a3bea` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-b56056ba7b5c86c6b74cf133` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-b5b5dbdca1e3c480d8a593f9` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-b6105f0613dd463baef8b97e` inconclusive, reason `renewal_review_inconclusive`, contained `true`, cleanup `true`
- `qa-v2-evidence-bd8fd861a8b2d54eb6480aba` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-be3cb76aefdeb3e0acf4dd5f` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-c055a1141b094f007c4d74fb` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-c0f42d428745dd46a6068f10` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-c4e1dc51960f4f3cf2df04ac` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-d1e543c41093b8fef36a873e` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-db7f9917320b9e4f2f6e86e1` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-dba00742a1372800e16303e9` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-dbadb1804ad4ea0b7a534469` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-de3d89be043e44015c79b173` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-f4ba74024edbf92397debb1c` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-f5bedcd64a3b884f7c029e3b` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-f97f0fa50ec5313302f68de1` inconclusive, reason `renewal_review_inconclusive`, contained `true`, cleanup `true`
- `qa-v2-evidence-fbc05e34b3dc98412929c5d2` pass, reason `check_passed`, contained `true`, cleanup `true`

## Rejected evidence
- `qa-v2-evidence-b6105f0613dd463baef8b97e` `evidence_blocked`: the evidence attempt did not complete
- `qa-v2-evidence-f97f0fa50ec5313302f68de1` `evidence_blocked`: the evidence attempt did not complete

## Issue candidates

Total: `0`
Promoted: `0`
Unpromoted: `0`

## Resolved claims

- `qa-v2-issue-01916e6c0e0566e8940471ab` Optimization continues after isolation cleanup failure: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-99d53cdcca1098534de53e2d/recovery-result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-99d53cdcca1098534de53e2d/qa-resolution/verification-5ba7d6b7921cf847e4424914499227b4fb2079abe260785f3d342fe29724177b.json`.

- `qa-v2-issue-377a3311d04fcdd86069f900` Performance CLI compatibility coverage is incomplete: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-1b4c5f6f05794a109ae3f83a/recovery-result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-1b4c5f6f05794a109ae3f83a/qa-resolution/verification-9181d88b2e6cfe6bdab248e122dfccd6df437f15f46cacd20350eff4c5c8d343.json`.

- `qa-v2-issue-7972b68261d6362477e5dfe2` Passing repair records need no target coverage: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-fa422483ec382a9eb3811722/result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-fa422483ec382a9eb3811722/qa-resolution/verification-13c52dfb7479531788c04effb7ffbaace01bebfeefaceb80133ce9b1717fcf11.json`.

- `qa-v2-issue-30084c02949cd30257cb8c94` Default performance evidence URL returns an error: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/qa-resolution/verification-f6abe5d7f003ab59d37f8e0e9ebc83177dd0e6437f3a1367f017e86e9672a5dc.json`.

- `qa-v2-issue-5dfa0c7375a6627e5fda7e93` Performance policy mode parsing is not exact: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/qa-resolution/verification-102cef10eecdb16ef6ac6a9a49e142b297ce13b339f4a0af2752b52f600472ba.json`.

- `qa-v2-issue-ba299cfce0c3e44dcf2e6426` Failed measurements evade the command budget: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/qa-resolution/verification-102cef10eecdb16ef6ac6a9a49e142b297ce13b339f4a0af2752b52f600472ba.json`.

- `qa-v2-issue-484ccba6a60d3497d4a35c1a` Performance JSON parser accepts malformed suffixes: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/qa-resolution/verification-102cef10eecdb16ef6ac6a9a49e142b297ce13b339f4a0af2752b52f600472ba.json`.

- `qa-v2-issue-cd244860592754728c9b5135` Report-only outcomes make flow state invalid: `resolved_by_verified_repair`; repair `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/result.json`; verification `projects/ultraplan-go/sprints/39-performance-stage/verification/attempts/qa-v1-attempt-68b6ad13df98a7d5a96df5d6/repairs/repair-v1-run-0b63617fc2789ddd585fe477/qa-resolution/verification-102cef10eecdb16ef6ac6a9a49e142b297ce13b339f4a0af2752b52f600472ba.json`.

## Promoted issues

None.

## Theory assertion coverage

| Theory | Accepted failing bundle / assertion |
| --- | --- |
| `qa-v1-theory-1007492d20bbcc01841fb291` | Missing evidence |
| `qa-v1-theory-35a6a9cbcf6592fb04bb7828` | qa-v2-test-aec431cf121861ff9963d55b / TestQAInvestigator_a02b8a23413a, qa-v2-test-b32c3bca434c97ee8b522d5a / TestQAInvestigator_dd232b1d46ba |
| `qa-v1-theory-3ba88b4e640cc758bd4724ff` | qa-v2-test-f50f6e6a22217003e70ab0c9 / TestQAInvestigator_e2eec52fa86d |
| `qa-v1-theory-44d9a1998804c860f04368c4` | qa-v2-test-16a0129229c8391e97b6af00 / TestQAInvestigator_d48db6847a74, qa-v2-test-f50f6e6a22217003e70ab0c9 / TestQAInvestigator_e2eec52fa86d |
| `qa-v1-theory-53f10905cf5035b40ecbb1e8` | Missing evidence |
| `qa-v1-theory-7ef7c4876857c8e4b85d4c0c` | qa-v2-test-731bcfbd077beb53b1c37659 / TestQAInvestigator_6b0e9329abe5, qa-v2-test-c418b7d27881f9b5b43f7e46 / TestQAInvestigator_d094a0cec36e |
| `qa-v1-theory-9a2d3dba12c4d353a09f5c53` | Missing evidence |
| `qa-v1-theory-bd5d5990c9d330ce736e467c` | Missing evidence |
| `qa-v1-theory-ea7c7adacc9ea256c6314ce8` | Missing evidence |
| `qa-v1-theory-f2dfc2c802889fddc3bbae5c` | Missing evidence |
| `qa-v1-theory-fa8487cda772b1748115a5c2` | Missing evidence |

## Next action

All retained QA obligations are resolved and verified on the current target.
