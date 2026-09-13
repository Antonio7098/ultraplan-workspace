package sprint

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pprocess "github.com/Antonio7098/ultraplan-go/internal/platform/process"
	pruntime "github.com/Antonio7098/ultraplan-go/internal/platform/runtime"
)

func TestQAInvestigator_76299fb53a2e(t *testing.T) {
	const marker = "ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_76299fb53a2e"
	target := gitFixture(t)
	if err := os.WriteFile(filepath.Join(target, "work.go"), []byte("package performance\n\nconst work = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, target, "add", "work.go")
	runGitTest(t, target, "commit", "-m", "add production source")

	performanceTarget := PerformanceTarget{ID: "PERF-LATENCY", Scenario: "package operation", Metric: "latency", Comparator: PerformanceLessOrEqual, Value: "10", Unit: PerformanceNSPerOp, Gate: PerformanceRequired, Samples: 2, Basis: PerformanceAbsolute}
	descriptor := performanceDescriptor(PerformanceGoBenchmark, ".", "BenchmarkLatency", PerformanceNSPerOp, performanceTarget.ID)
	baseline := PerformanceMeasurement{Target: performanceTarget, Qualification: PerformanceQualification{Median: "20", Qualified: true}}
	limits := DefaultPerformanceLimits()
	limits.Warmups = 1
	limits.Commands = 20
	limits.OptimizationCycles = 8

	controlRunner := &qaPerformanceAccountingRunner{}
	control, err := MeasurePerformanceTarget(context.Background(), controlRunner, target, descriptor, performanceTarget, limits, nil, fingerprintOf("packet"), fingerprintOf("benchmark"), fingerprintOf("environment"))
	if err != nil || !control.Qualification.Qualified {
		t.Fatalf("passing control measurement = %#v, %v", control, err)
	}
	controlCounters := PerformanceCounters{Commands: limits.Warmups + performanceTarget.Samples, Warmups: limits.Warmups, Samples: performanceTarget.Samples}
	if controlRunner.calls != controlCounters.Commands {
		t.Fatalf("passing control calls = %d, accounted = %d", controlRunner.calls, controlCounters.Commands)
	}

	root := t.TempDir()
	sp := Sprint{Project: "proj", Slug: "39-performance", Path: filepath.Join(root, "projects", "proj", "sprints", "39-performance")}
	if err := os.MkdirAll(sp.Path, 0o755); err != nil {
		t.Fatal(err)
	}
	token := VerificationWriterToken{RunID: "run-performance", OperationalAttemptID: "attempt-performance", FencingGeneration: 1}
	attemptID := "performance-v1-attempt-0123456789abcdef01234567"
	ref := func(name string) *PerformanceArtifactRef {
		path, pathErr := PerformanceAttemptRelPath(sp, attemptID, name)
		if pathErr != nil {
			t.Fatal(pathErr)
		}
		return &PerformanceArtifactRef{Path: path, Digest: fingerprintOf(name)}
	}
	packetRef := ref("target-packet.json")
	state := PerformanceState{
		SchemaVersion: PerformanceSchemaVersion, Project: sp.Project, Sprint: sp.Slug, AttemptID: attemptID,
		Phase: PerformancePhaseOptimizing, Correlation: token, Deadline: time.Now().Add(time.Hour), Limits: limits,
		Packet: *packetRef, TargetPacketDigest: packetRef.Digest, Manifest: ref("benchmark-manifest.json"), Environment: ref("environment.json"), Baseline: ref("baseline.json"),
		CurrentImplementation: fingerprintOf("initial"), NextAction: "Optimize.", UpdatedAt: time.Now(),
	}
	store := NewPerformanceStore(root, sp).WithWriterFence(func(got VerificationWriterToken) error { return nil })
	runner := &qaPerformanceAccountingRunner{failEveryThirdMeasurement: true}
	service := NewService(root).WithRuntime(qaPerformanceAccountingRuntime{}).WithProcessRunner(runner)
	prepared := PerformancePrepareResult{Project: sp.Project, Sprint: sp.Slug, Packet: PerformanceTargetPacket{PacketDigest: packetRef.Digest, Targets: []PerformanceTarget{performanceTarget}}, Target: ExecuteTargetRef{Path: target}, Limits: limits}
	manifest := PerformanceBenchmarkManifest{Descriptors: []PerformanceDescriptor{descriptor}, BenchmarkDigest: fingerprintOf("benchmark")}

	result, err := service.optimizePerformance(context.Background(), prepared, manifest, []PerformanceMeasurement{baseline}, store, &state, token, nil)
	if err != nil {
		t.Fatalf("optimization reproduction returned unrelated error: %v", err)
	}
	if !result.Stalled || state.Counters.OptimizationCycles != limits.OptimizationCycles {
		t.Fatalf("optimization did not reach the bounded-cycle control: result=%#v counters=%#v", result, state.Counters)
	}
	if runner.calls <= limits.Commands || state.Counters.Commands >= runner.calls {
		return
	}
	assertion := "For qa-v1-theory-53f10905cf5035b40ecbb1e8, make candidate or regression measurement fail after one or more runner calls, record actual invocations and durable command counters, and assert whether a later cycle starts such that actual cumulative commands ex"
	t.Fatalf("%s assertion=%s actual_calls=%d durable_commands=%d command_limit=%d output_matcher=%s", marker, assertion, runner.calls, state.Counters.Commands, limits.Commands, marker)
}

type qaPerformanceAccountingRunner struct {
	calls                     int
	measurementCalls          int
	failEveryThirdMeasurement bool
}

func (r *qaPerformanceAccountingRunner) Run(_ context.Context, request pprocess.Request) (pprocess.Result, error) {
	r.calls++
	profile := false
	for _, arg := range request.Args {
		profile = profile || strings.HasPrefix(arg, "-cpuprofile=")
	}
	if profile {
		return pprocess.Result{Stdout: "profile", CleanupComplete: true}, nil
	}
	r.measurementCalls++
	if r.failEveryThirdMeasurement && r.measurementCalls%3 == 0 {
		return pprocess.Result{Stdout: "invalid benchmark output\n", CleanupComplete: true}, nil
	}
	return pprocess.Result{Stdout: "BenchmarkLatency-8 100 9 ns/op\n", CleanupComplete: true}, nil
}

type qaPerformanceAccountingRuntime struct{}

func (qaPerformanceAccountingRuntime) StartRun(_ context.Context, request pruntime.Request) (pruntime.Result, error) {
	path := filepath.Join(request.WorkDir, "work.go")
	if err := os.WriteFile(path, []byte("package performance\n\nconst work = 2\n"), 0o644); err != nil {
		return pruntime.Result{}, err
	}
	return pruntime.Result{Status: "success", RunID: "qa-performance-accounting"}, nil
}
