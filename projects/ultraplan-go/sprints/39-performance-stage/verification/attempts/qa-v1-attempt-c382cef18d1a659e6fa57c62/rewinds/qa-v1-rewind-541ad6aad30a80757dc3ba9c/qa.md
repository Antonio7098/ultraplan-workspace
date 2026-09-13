# QA

Project: `ultraplan-go`
Sprint: `39-performance-stage`
Input fingerprint: `fa10644c48b9941884ca51d75057ecc62e7f85f161742d3097eb36bc270a7570`
Attempt: `qa-v1-attempt-c382cef18d1a659e6fa57c62`
Assessment: `blocked`

## Evidence

Accepted: `32`
Rejected: `0`
- `qa-v2-evidence-00dbf90eb8bd3fb7bafe24e5` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-14b29a0ed41b25ad9b6691c1` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-1e81f0cb77c6d0fd7ddc7a86` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-24ad942f17348380c0a64cfd` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-2e14e530b280e545e92c6989` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-32840b5cd1df0243c4bed24e` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-37ec38bb11681c3e3874fbcf` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-3bc66810797a2c15ecf0e6af` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-515a9349ab1e38d3bacffbb2` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-5ac4ce346fa6bace9cc10d80` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-61dc129e7128336484283c7d` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-68bb93e04be58e050319697b` fail, reason `predicted_failure_reproduced`, contained `true`, cleanup `true`
- `qa-v2-evidence-79a6d251bd02763160125d85` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-7ad1d2d79cdd153a166122e2` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-7f08761e25e1a2a48ec22c28` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-83daf2823f1110ab34e0f926` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-8b8e792871f357c04458bf7c` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-8cb5b9ced79072ab9e252b3d` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-96a8f6227ababbab0c097672` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-aa7a27c3807649fd206f2122` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-b1acb6b3c341536cea0525ea` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-b3f2138562e576213b297633` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-ba717168c668291afe6f1bd2` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-ba9bfa88fee4b87debe62179` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-cf61f3d66cc505154c711916` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-d6a19ad118982283d1201b7b` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-e3e09efe1bb94479a48975e7` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-e511b3db0edf4953ab5d3b01` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-eb38e2fa4d1a71fff4b1e062` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-ee8e63fdf4bdb98fb781baac` pass, reason `check_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-f500af63bda7e21bb3118a2d` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`
- `qa-v2-evidence-ffaefcf62364921459cad551` pass, reason `go_source_integrity_passed`, contained `true`, cleanup `true`

## Issue candidates

Total: `65`
Promoted: `1`
Unpromoted: `64`
- `qa-v1-arbiter-issue-00b760b8ff1978cdd8ae29b1` [high] Performance status surfaces omit durable run state at `internal/app/sprint_usecases.go; internal/web/handlers.go; internal/tui`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-08185f25ca285909f34b8a92` [high] Performance surfaces omit durable lifecycle state at `internal/app/sprint_usecases.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-0e251f225cb386ec9c98e70b` [medium] Policy validation mistakes pseudo target headings for declarations at `internal/sprint/index.go; internal/sprint/performance_targets.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-16d4c8ca368f724c50bd7b51` [high] Unavailable result loses error classification at `internal/app/sprint_usecases.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-1869d0041e03aadabc2cb167` [medium] Policy validation misclassifies fenced and commented headings at `internal/sprint/index.go; internal/sprint/performance_targets.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-1f3867631e4cd4140e2bbda0` [medium] Stale performance writer can publish flow summary at `internal/sprint/performance.go; internal/sprint/state.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-20ca7d37a2089084589974fa` [high] Performance surfaces omit durable run state at `internal/app/sprint_usecases.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-22c7f47af2883036849b3f4e` [medium] Pseudo target headings cause false policy mismatch at `internal/sprint/index.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-260bb2ee701726d0c2f255d8` [medium] Failure flow publication lacks final writer fence at `internal/sprint/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-2700697cfb03c0cc3b555aae` [medium] Evidence fallback violates the versioned response contract at `internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-2a8495df2dfd7fc7706d8e9f` [medium] Operation runner swallows performance configuration failure at `internal/app/operation_runner.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-2b28ac6f55450609429d374e` [medium] Descriptor markers need not immediately precede functions at `internal/sprint/performance_discovery.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-2c49684e0f83c9f59a9d540b` [high] Performance surfaces omit durable lifecycle state at `internal/app/sprint_usecases.go; internal/web/handlers.go; internal/web/templates/sprint.html; internal/tui`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-2f886a31bee5a57d3cf50632` [high] Missing performance result returns an unclassified error at `internal/app/sprint_usecases.go; internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-30111861767dfd2670c9551e` [medium] Descriptor discovery weakens marker adjacency at `internal/sprint/performance_discovery.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-3467039c2c57b94a7ee30779` [high] Performance policy is embedded in platform config at `internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-3b4b9f094c1e1ed89af4d76b` [high] Repair reverification accepts malformed target results at `internal/sprint/qa_repair_state.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-46d5f6069848caca0ffba83f` [medium] Pseudo performance headings trigger policy mismatch at `internal/sprint/index.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-4a1182f9829203d097b9a34c` [high] Documented performance maxima are rejected at `docs/configuration.md`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-510175ad99d31a759a8ea2bb` [high] Disabled projects expose performance mutations at `internal/tui/model.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-54f56fed54dd5675a39cb4d4` [high] TUI exposes mutations for disabled performance policy at `internal/tui/model.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-5755b409e74945ba6a47642d` [high] Repair reverification accepts malformed nested performance results at `internal/sprint/qa_repair_state.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-5d247538c0eefca99f87c4b0` [medium] Evidence fallback violates the versioned response contract at `internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-5ec354e4ce0c317e6d6c6bea` [high] Performance limit documentation contradicts validation at `docs/configuration.md; internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-5fc5c8443fd53a722e69d7b5` [high] Platform config contains bounded-context performance policy at `internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-60a41a36a418052a408fa47f` [high] Performance JSON parser accepts malformed trailing data at `internal/sprint/performance_measure.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-69f05dc7057e52426b65660a` [high] Disabled historical performance states expose mutations at `internal/tui/model.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-6d489eecd4f1d7901fe9455c` [high] Performance CLI loses failure exit categories at `internal/app/sprint_commands.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-71eebcecdaa868156132da3c` [high] Performance limit documentation contradicts validation at `docs/configuration.md; internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-766c1e0ae77b0e2872f8cdba` [high] Performance summaries omit durable lifecycle state at `internal/app/sprint_usecases.go; internal/web/handlers.go; internal/web/templates/sprint.html; internal/tui`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-807647a2078d976b603817fa` [medium] Evidence API has an incompatible fallback schema at `internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-83eae721403c69628743568d` [high] Platform config owns performance product policy at `internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-8459b233ed457ad16187ea73` [high] Unavailable result returns an unclassified error at `internal/app/sprint_usecases.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-84a0ee8828f3e415b400863e` [high] Performance JSON parser accepts malformed trailing bytes at `internal/sprint/performance_measure.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-84ca632f1f983370da8e43b4` [medium] Evidence fallback violates versioned response contract at `internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-86fc652cabc152633bbd1088` [medium] Pseudo performance headings trigger policy mismatch at `internal/sprint/index.go; internal/sprint/performance_targets.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-88cb9c677536a12726659e7a` [medium] Descriptor markers need not immediately precede functions at `internal/sprint/performance_discovery.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-89c0d8c263f701eda12ff78d` [high] Unavailable performance result is returned unclassified at `internal/app/sprint_usecases.go; internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-8c8c86a4dddf2af73083d609` [high] Performance JSON parser accepts malformed trailing data at `internal/sprint/performance_measure.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-8f73380166c68c44207b9363` [high] Performance CLI collapses operational exit classes at `internal/app/sprint_commands.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-9b7b566eda8a3201f07fc46b` [high] Performance product policy is embedded in platform config at `internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-9ef8e3508eb4d2defab3b64e` [high] Disabled performance policy still exposes TUI mutations at `internal/tui/model.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-a5bc1f6e7c25dffd9c5c6994` [medium] Operation runner swallows performance configuration errors at `internal/app/operation_runner.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-a85a5b53963e8ba4393f2e2f` [high] Performance summaries omit durable run state at `internal/app/sprint_usecases.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-b2a1e60b83cd15fc34b8b8fe` [medium] Flow-state failure publication lacks a final writer fence at `internal/sprint/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-b3df58b1b95e02beea4a1b89` [medium] Evidence compatibility response breaks schema at `internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-b87289beabb2a581ef90aa73` [medium] Descriptor discovery weakens immediate-precedence rule at `internal/sprint/performance_discovery.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-b93c47b4944ae4b5a829ec93` [high] Unavailable performance results bypass typed error mapping at `internal/app/sprint_usecases.go; internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-bbab961b7867b4b527d27409` [high] Performance CLI collapses operational exit classes at `internal/app/sprint_commands.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-bcc3a3ac83e9c7a85319226d` [medium] Operation runner discards fatal performance configuration errors at `internal/app/operation_runner.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-bf8fdf8d2ac392561f31448f` [medium] Failure summary publication is not writer-fenced at `internal/sprint/performance.go; internal/sprint/state.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-c573658a8955e11efd23b875` [high] Repair reverification omits nested performance validation at `internal/sprint/qa_repair_state.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-c6028d29cd7aa6fc06f11842` [high] Performance CLI loses category-specific exit behavior at `internal/app/sprint_commands.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-c9117dbdcdb2f825dfa8a906` [high] Performance JSON parser does not require EOF at `internal/sprint/performance_measure.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-ca217848a5c9009952979e82` [high] Performance CLI loses failure-specific exit behavior at `internal/app/sprint_commands.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-cbf5cb5e27d2a210dbb475b1` [high] Performance limit documentation contradicts validation at `docs/configuration.md`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-ce3a3a55d063c08f2364af7b` [high] Disabled projects expose historical performance mutations in TUI at `internal/tui/model.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-e9df984cd57a06dd6a2db153` [high] Performance limit documentation contradicts validation at `docs/configuration.md; internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-eadbd1b5233ad8107cd689c7` [high] Repair reverification omits nested performance validation at `internal/sprint/qa_repair_state.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-eb4a05bdecf39114ffae473a` [high] Documented performance maxima are rejected at `docs/configuration.md; internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-ee2a2ec706ae7d40fa7bd258` [medium] Evidence API fallback violates the versioned response schema at `internal/web/handlers.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-ef0827fa91c4b61ca8341289` [high] Platform config contains performance business policy at `internal/platform/config/performance.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-f5f6a18e980b4ba76a08ee8f` [medium] Discovery weakens marker adjacency rule at `internal/sprint/performance_discovery.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate
- `qa-v1-arbiter-issue-f7529d27f24b197684317376` [high] Unavailable performance results lose error classification at `internal/app/sprint_usecases.go`, outcome `unpromoted`, reason `promotion_evidence_missing`: no accepted failing evidence was linked to this candidate

## Promoted issues
- `qa-v2-issue-b6303030807f9608fa8ec3ad` [medium] Shared operation runner drops performance configuration failures at `internal/app/operation_runner.go`, evidence `qa-v2-evidence-515a9349ab1e38d3bacffbb2`, `qa-v2-evidence-68bb93e04be58e050319697b`, regression candidate `true`

## Next action

Collect promotion evidence for 64 unpromoted issue candidates.
