package sprint

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pprocess "github.com/Antonio7098/ultraplan-go/internal/platform/process"
)

func TestQAInvestigator_b7efe9eb62a1(t *testing.T) {
	now := time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)
	gates := make([]RepairGateResult, 0, len(RepairGateOrder()))
	for _, gate := range RepairGateOrder() {
		gates = append(gates, RepairGateResult{Gate: gate, Status: RepairGatePassed})
	}
	valid := RepairReverification{
		SchemaVersion: QARepairSchemaVersion,
		RepairRunID:   "repair-v1-run-0123456789abcdef01234567",
		Cycle:         1,
		Gates:         gates,
		CompletedAt:   now,
		Performance: &RepairPerformanceReverification{
			SchemaVersion:          PerformanceSchemaVersion,
			AttemptID:              "performance-v1-attempt-0123456789abcdef01234567",
			PriorResultDigest:      fingerprintOf("prior-result"),
			RepairedImplementation: fingerprintOf("repaired-implementation"),
			Outcome:                PerformancePassed,
			CleanupComplete:        true,
			CompletedAt:            now,
		},
	}
	control := valid
	control.Performance = &RepairPerformanceReverification{
		SchemaVersion:          PerformanceSchemaVersion,
		AttemptID:              valid.Performance.AttemptID,
		PriorResultDigest:      valid.Performance.PriorResultDigest,
		RepairedImplementation: valid.Performance.RepairedImplementation,
		Outcome:                PerformancePassed,
		CleanupComplete:        false,
		CompletedAt:            now,
	}
	if err := ValidateRepairReverification(control); err == nil {
		t.Fatal("control: passing performance reverification without cleanup proof was accepted")
	}
	missingResultsAccepted := ValidateRepairReverification(valid) == nil

	root := workspaceFixture(t)
	sp := sprintFixture(t, root, "proj", "39-performance")
	target := gitFixture(t)
	benchmark := `package performance

// ultraplan:performance-v1 go-benchmark-v1 ns/op PERF-A
func BenchmarkA() {}

// ultraplan:performance-v1 go-benchmark-v1 ns/op PERF-B
func BenchmarkB() {}
`
	if err := os.WriteFile(filepath.Join(target, "performance_test.go"), []byte(benchmark), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, target, "add", "performance_test.go")
	runGitTest(t, target, "commit", "-m", "add performance benchmarks")
	writeFileContent(t, filepath.Join(root, "projects", "proj"), projectIndexForTarget(target)+"\n## Performance Policy\n\n- **Mode:** enabled\n", "project-index.md")
	writeFileContent(t, sp.Path, `# Requirements

## Performance Targets

| ID | Scenario | Metric | Comparator | Value | Unit | Gate | Samples | Basis |
| --- | --- | --- | --- | ---: | --- | --- | ---: | --- |
| PERF-A | operation A | latency change | <= | 10 | % | required | 5 | baseline |
| PERF-B | operation B | latency change | <= | 10 | % | required | 5 | baseline |
`, "requirements.md")
	writeFileContent(t, sp.Path, "# Execute\n\nComplete.\n", "execute.md")
	task := ExecuteTaskRecord{ID: "task-performance", Identity: ExecuteTaskIdentity{Name: "Performance implementation", PlanLine: 1}, Status: ExecuteTaskComplete, CreatedAt: now, UpdatedAt: now, CompletedAt: &now}
	execute := NewExecuteRunState(sp, ExecuteTargetRef{Path: target, Source: "project-index.md"}, ArtifactRelPath(sp, StagePlan), fingerprintOf("plan"), []ExecuteTaskRecord{task}, now)
	if err := SaveExecuteRunState(root, sp, execute); err != nil {
		t.Fatal(err)
	}
	if err := SaveFlowState(root, sp, NewFlowState(sp, completeStates(sp), now)); err != nil {
		t.Fatal(err)
	}
	token := VerificationWriterToken{RunID: "run-performance", OperationalAttemptID: "attempt-performance", FencingGeneration: 1}
	runner := qaInvestigatorPerformanceRunner{}
	service := NewService(root).WithClock(func() time.Time { return now }).WithProcessRunner(runner).WithVerificationWriterFence(func(got VerificationWriterToken) error {
		if got != token {
			return context.Canceled
		}
		return nil
	})
	prepared, err := service.PreparePerformance("proj", "39", true)
	if err != nil || !prepared.Ready {
		t.Fatalf("prepare performance: %#v, %v", prepared, err)
	}
	manifest := PerformanceBenchmarkManifest{
		SchemaVersion: PerformanceSchemaVersion,
		PacketDigest:  prepared.Packet.PacketDigest,
		RunnerVersion: PerformanceRunnerVersion,
		ParserVersion: PerformanceParserVersion,
		Descriptors: []PerformanceDescriptor{
			performanceDescriptor(PerformanceGoBenchmark, ".", "BenchmarkA", PerformanceNSPerOp, "PERF-A"),
			performanceDescriptor(PerformanceGoBenchmark, ".", "BenchmarkB", PerformanceNSPerOp, "PERF-B"),
		},
	}
	if err := ValidatePerformanceManifest(prepared.Packet, &manifest); err != nil {
		t.Fatal(err)
	}
	if _, err := service.RunPerformance(context.Background(), "proj", "39", PerformanceRunRequest{WriterToken: token, Manifest: manifest, Environment: allowPerformanceEnvironment(os.Environ())}); err != nil {
		t.Fatal(err)
	}
	ordered := service.runRepairPerformanceReverification(context.Background(), sp, target)
	if ordered.Outcome != PerformancePassed || len(ordered.TargetResults) != 2 || ordered.TargetResults[0].BaselineQualification.Median != "100" || ordered.TargetResults[1].BaselineQualification.Median != "200" {
		t.Fatalf("control: correctly ordered baselines were not associated by target: %#v", ordered)
	}
	state, err := LoadPerformanceState(root, sp)
	if err != nil {
		t.Fatal(err)
	}
	store := NewPerformanceStore(root, sp).WithWriterFence(service.verificationWriterFence)
	var baseline []PerformanceMeasurement
	if err := store.ReadJSON(*state.Baseline, &baseline); err != nil {
		t.Fatal(err)
	}
	baseline[0], baseline[1] = baseline[1], baseline[0]
	reorderedRef, err := store.PublishJSON(state.Baseline.Path, "baseline", baseline, false, token)
	if err != nil {
		t.Fatal(err)
	}
	state.Baseline = &reorderedRef
	if _, err := store.PublishState(state, token); err != nil {
		t.Fatal(err)
	}
	reordered := service.runRepairPerformanceReverification(context.Background(), sp, target)
	reorderedMisassociated := reordered.Outcome == PerformanceTargetMiss && len(reordered.TargetResults) == 2 &&
		reordered.TargetResults[0].Target.ID == "PERF-A" && reordered.TargetResults[0].BaselineQualification.Median == "200" &&
		reordered.TargetResults[1].Target.ID == "PERF-B" && reordered.TargetResults[1].BaselineQualification.Median == "100"
	if missingResultsAccepted && reorderedMisassociated {
		fmt.Println("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_b7efe9eb62a1")
		t.Fatalf("Assertions for qa-v1-theory-ea7c7adacc9ea256c6314ce8 must show whether a passing performance reverification with zero or incomplete TargetResults is accepted. Separate assertions for qa-v1-theory-1007492d20bbcc01841fb291 must use at least two distinguishab")
	}
	if !missingResultsAccepted {
		t.Fatal("passing performance reverification with zero TargetResults was rejected")
	}
	t.Fatalf("reordered baselines were not observably misassociated: %#v", reordered)
}

type qaInvestigatorPerformanceRunner struct{}

func (qaInvestigatorPerformanceRunner) Run(_ context.Context, request pprocess.Request) (pprocess.Result, error) {
	args := strings.Join(request.Args, " ")
	stdout := ""
	if strings.Contains(args, "-bench=^BenchmarkA$") {
		stdout = "BenchmarkA-8 100 100 ns/op\n"
	}
	if strings.Contains(args, "-bench=^BenchmarkB$") {
		stdout = "BenchmarkB-8 100 200 ns/op\n"
	}
	return pprocess.Result{Stdout: stdout, CleanupComplete: true}, nil
}
