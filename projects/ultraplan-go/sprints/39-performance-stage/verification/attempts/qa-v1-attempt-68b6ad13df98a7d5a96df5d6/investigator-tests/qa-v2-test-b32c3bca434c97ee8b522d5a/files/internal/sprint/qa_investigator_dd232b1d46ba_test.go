package sprint

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pprocess "github.com/Antonio7098/ultraplan-go/internal/platform/process"
	pruntime "github.com/Antonio7098/ultraplan-go/internal/platform/runtime"
)

func TestQAInvestigator_dd232b1d46ba(t *testing.T) {
	const marker = "ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_dd232b1d46ba"
	run := func(t *testing.T, failCleanup bool) (int, performanceOptimizationResult, error) {
		t.Helper()
		target := gitFixture(t)
		if err := os.WriteFile(filepath.Join(target, "work.go"), []byte("package performance\n\nconst work = 1\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		runGitTest(t, target, "add", "work.go")
		runGitTest(t, target, "commit", "-m", "add production source")

		performanceTarget := PerformanceTarget{ID: "PERF-CLEANUP", Scenario: "package operation", Metric: "latency", Comparator: PerformanceLessOrEqual, Value: "10", Unit: PerformanceNSPerOp, Gate: PerformanceRequired, Samples: 5, Basis: PerformanceAbsolute}
		descriptor := performanceDescriptor(PerformanceGoBenchmark, ".", "BenchmarkLatency", PerformanceNSPerOp, performanceTarget.ID)
		baseline := PerformanceMeasurement{Target: performanceTarget, Qualification: PerformanceQualification{Median: "20", Qualified: true}}
		limits := DefaultPerformanceLimits()
		limits.Warmups = 1
		limits.Commands = 40
		limits.OptimizationCycles = 2
		limits.OptimizationAttempts = 1

		root := t.TempDir()
		sp := Sprint{Project: "proj", Slug: "39-performance", Path: filepath.Join(root, "projects", "proj", "sprints", "39-performance")}
		if err := os.MkdirAll(sp.Path, 0o755); err != nil {
			t.Fatal(err)
		}
		token := VerificationWriterToken{RunID: "run-performance", OperationalAttemptID: "attempt-performance", FencingGeneration: 1}
		attemptID := "performance-v1-attempt-0123456789abcdef01234567"
		ref := func(name string) *PerformanceArtifactRef {
			path, err := PerformanceAttemptRelPath(sp, attemptID, name)
			if err != nil {
				t.Fatal(err)
			}
			return &PerformanceArtifactRef{Path: path, Digest: fingerprintOf(name)}
		}
		packetRef := ref("target-packet.json")
		state := PerformanceState{
			SchemaVersion: PerformanceSchemaVersion, Project: sp.Project, Sprint: sp.Slug, AttemptID: attemptID,
			Phase: PerformancePhaseOptimizing, Fresh: true, Correlation: token, Deadline: time.Now().Add(time.Hour), Limits: limits,
			Packet: *packetRef, TargetPacketDigest: packetRef.Digest, Manifest: ref("benchmark-manifest.json"), Environment: ref("environment.json"), Baseline: ref("baseline.json"),
			CurrentImplementation: fingerprintOf("initial"), NextAction: "Optimize.", UpdatedAt: time.Now(),
		}
		store := NewPerformanceStore(root, sp).WithWriterFence(func(VerificationWriterToken) error { return nil })
		runtime := &qaCleanupFailureRuntimeDD232{failCleanup: failCleanup}
		service := NewService(root).WithRuntime(runtime).WithProcessRunner(qaCleanupProfileRunnerDD232{})
		prepared := PerformancePrepareResult{Project: sp.Project, Sprint: sp.Slug, Packet: PerformanceTargetPacket{PacketDigest: packetRef.Digest, Targets: []PerformanceTarget{performanceTarget}}, Target: ExecuteTargetRef{Path: target}, Limits: limits}
		manifest := PerformanceBenchmarkManifest{Descriptors: []PerformanceDescriptor{descriptor}, BenchmarkDigest: fingerprintOf("benchmark")}

		result, err := service.optimizePerformance(context.Background(), prepared, manifest, []PerformanceMeasurement{baseline}, store, &state, token, nil)
		for _, parent := range runtime.parents {
			_ = os.RemoveAll(parent)
		}
		return runtime.calls, result, err
	}

	controlCalls, controlResult, controlErr := run(t, false)
	if controlErr != nil || !controlResult.Stalled || controlCalls != 2 {
		t.Fatalf("successful-cleanup control failed: calls=%d result=%#v err=%v", controlCalls, controlResult, controlErr)
	}

	calls, result, err := run(t, true)
	cleanupUncertain := err != nil && strings.Contains(err.Error(), "performance cleanup is uncertain")
	if calls > 1 || !cleanupUncertain {
		assertion := "Induce an actual optimization rejection followed by cleanup failure, then assert whether later work starts and whether returned or persisted state is cleanup_uncertain for qa-v1-theory-35a6a9cbcf6592fb04bb7828."
		t.Fatalf("%s assertion=%s later_work_started=%t runtime_calls=%d stalled=%t cleanup_uncertain=%t err=%v", marker, assertion, calls > 1, calls, result.Stalled, cleanupUncertain, err)
	}
}

type qaCleanupFailureRuntimeDD232 struct {
	calls       int
	failCleanup bool
	parents     []string
}

func (r *qaCleanupFailureRuntimeDD232) StartRun(_ context.Context, request pruntime.Request) (pruntime.Result, error) {
	r.calls++
	if r.failCleanup && r.calls == 1 {
		parent := filepath.Dir(request.WorkDir)
		r.parents = append(r.parents, parent)
		if err := os.WriteFile(filepath.Join(parent, "cleanup-blocker"), []byte("retain parent"), 0o600); err != nil {
			return pruntime.Result{}, err
		}
	}
	return pruntime.Result{}, errors.New("proposal rejected")
}

type qaCleanupProfileRunnerDD232 struct{}

func (qaCleanupProfileRunnerDD232) Run(_ context.Context, _ pprocess.Request) (pprocess.Result, error) {
	return pprocess.Result{Stdout: "profile", CleanupComplete: true}, nil
}
