package sprint

import (
	"context"
	"fmt"
	"os"
	"strings"

	pprocess "github.com/Antonio7098/ultraplan-go/internal/platform/process"
)

// runRepairPerformanceReverification measures the frozen target set against
// the repaired implementation after the functional repair ladder has passed.
// It never updates or reuses the old performance result as current evidence.
func (s Service) runRepairPerformanceReverification(ctx context.Context, sp Sprint, target string) RepairPerformanceReverification {
	record := RepairPerformanceReverification{SchemaVersion: PerformanceSchemaVersion, CleanupComplete: true, Outcome: PerformanceBlocked}
	blocked := func(cause error, cleanup bool) RepairPerformanceReverification {
		record.Blocker = displayPerformanceCause(cause)
		record.CleanupComplete = cleanup
		if !cleanup {
			record.Outcome = PerformanceCleanupUncertain
		}
		record.CompletedAt = s.now().UTC()
		return record
	}
	state, err := LoadPerformanceState(s.root, sp)
	if err != nil || state.Result == nil || state.Manifest == nil || state.Environment == nil || state.Baseline == nil || state.Phase != PerformancePhaseTerminal || !state.Fresh {
		if err == nil {
			err = fmt.Errorf("performance state is incomplete or stale")
		}
		return blocked(fmt.Errorf("fresh frozen performance authority is unavailable: %w", err), true)
	}
	record.AttemptID, record.PriorResultDigest = state.AttemptID, state.Result.Digest
	store := NewPerformanceStore(s.root, sp)
	if _, err := store.readResult(*state.Result); err != nil {
		return blocked(fmt.Errorf("read prior performance result: %w", err), true)
	}
	var packet PerformanceTargetPacket
	var manifest PerformanceBenchmarkManifest
	var environment PerformanceEnvironment
	var baseline []PerformanceMeasurement
	if err := store.ReadJSON(state.Packet, &packet); err != nil {
		return blocked(err, true)
	}
	if err := store.ReadJSON(*state.Manifest, &manifest); err != nil {
		return blocked(err, true)
	}
	if err := store.ReadJSON(*state.Environment, &environment); err != nil {
		return blocked(err, true)
	}
	if err := store.ReadJSON(*state.Baseline, &baseline); err != nil {
		return blocked(err, true)
	}
	if len(baseline) != len(packet.Targets) || environment.Identity != currentPerformanceEnvironment(os.Environ()).Identity {
		return blocked(fmt.Errorf("frozen performance environment or baseline coverage changed"), true)
	}
	currentManifest, err := DiscoverPerformanceManifest(target, packet)
	if err != nil || currentManifest.BenchmarkDigest != manifest.BenchmarkDigest {
		return blocked(fmt.Errorf("frozen performance benchmark changed: %w", err), true)
	}
	repairedIdentity, err := targetIdentity(target)
	if err != nil {
		return blocked(err, true)
	}
	record.RepairedImplementation = repairedIdentity
	isolation, err := createPerformanceIsolation(ctx, target, []string{sp.Path}, state.Limits)
	if err != nil {
		return blocked(err, true)
	}
	baseRunner := s.processRunner
	if baseRunner == nil {
		baseRunner = pprocess.DirectRunner{}
	}
	descriptors := make(map[string]PerformanceDescriptor, len(packet.Targets))
	for _, descriptor := range manifest.Descriptors {
		for _, targetID := range descriptor.TargetIDs {
			descriptors[targetID] = descriptor
		}
	}
	outcomes := make([]PerformanceTargetOutcome, 0, len(packet.Targets))
	for index, performanceTarget := range packet.Targets {
		measurement, measureErr := MeasurePerformanceTarget(ctx, isolation.runner(baseRunner), isolation.workspace.Path, descriptors[performanceTarget.ID], performanceTarget, state.Limits, allowPerformanceEnvironment(os.Environ()), packet.PacketDigest, manifest.BenchmarkDigest, environment.Identity)
		if measureErr != nil {
			cleanupErr := isolation.cleanup()
			return blocked(fmt.Errorf("repair performance target %s: %w", performanceTarget.ID, measureErr), cleanupErr == nil)
		}
		comparison, compareErr := ComparePerformanceTarget(performanceTarget, baseline[index].Qualification.Median, measurement.Qualification.Median, measurement.Qualification.Qualified, false)
		if compareErr != nil {
			cleanupErr := isolation.cleanup()
			return blocked(compareErr, cleanupErr == nil)
		}
		record.TargetResults = append(record.TargetResults, PerformanceTargetResult{Target: performanceTarget, BaselineQualification: baseline[index].Qualification, FinalQualification: measurement.Qualification, Comparison: comparison})
		outcomes = append(outcomes, comparison.Outcome)
	}
	if err := isolation.assertUnchanged(ctx); err != nil {
		cleanupErr := isolation.cleanup()
		return blocked(err, cleanupErr == nil)
	}
	if err := isolation.cleanup(); err != nil {
		return blocked(err, false)
	}
	finalIdentity, err := targetIdentity(target)
	if err != nil || finalIdentity != repairedIdentity {
		return blocked(fmt.Errorf("repaired implementation changed during performance reverification: %w", err), true)
	}
	record.Outcome = DerivePerformanceRunOutcome(outcomes, false, false, false, false)
	if record.Outcome != PerformancePassed && record.Outcome != PerformancePassedWithReports {
		var failed []string
		for _, result := range record.TargetResults {
			if result.Target.Gate == PerformanceRequired && result.Comparison.Outcome != PerformanceTargetMet {
				failed = append(failed, result.Target.ID)
			}
		}
		record.Blocker = "repaired implementation did not preserve required performance targets"
		if len(failed) > 0 {
			record.Blocker += ": " + strings.Join(failed, ", ")
		}
	}
	record.CompletedAt = s.now().UTC()
	return record
}
