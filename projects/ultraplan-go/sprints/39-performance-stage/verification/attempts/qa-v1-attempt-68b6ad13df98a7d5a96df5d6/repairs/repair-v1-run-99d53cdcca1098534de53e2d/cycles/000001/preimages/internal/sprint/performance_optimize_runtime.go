package sprint

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	pprocess "github.com/Antonio7098/ultraplan-go/internal/platform/process"
)

type performanceOptimizationResult struct {
	Changed             bool
	Stalled             bool
	FinalImplementation string
}

func (s Service) optimizePerformance(ctx context.Context, prepared PerformancePrepareResult, manifest PerformanceBenchmarkManifest, baseline []PerformanceMeasurement, store PerformanceStore, state *PerformanceState, token VerificationWriterToken, env []string) (performanceOptimizationResult, error) {
	identity, err := targetIdentity(prepared.Target.Path)
	if err != nil {
		return performanceOptimizationResult{}, err
	}
	out := performanceOptimizationResult{FinalImplementation: identity}
	if s.runtime == nil {
		return out, nil
	}
	descriptors := map[string]PerformanceDescriptor{}
	for _, descriptor := range manifest.Descriptors {
		for _, id := range descriptor.TargetIDs {
			descriptors[id] = descriptor
		}
	}
	current := make(map[string]PerformanceMeasurement, len(baseline))
	for _, measurement := range baseline {
		current[measurement.Target.ID] = measurement
	}
	baseRunner := s.processRunner
	if baseRunner == nil {
		baseRunner = pprocess.DirectRunner{}
	}
	for state.Counters.OptimizationCycles < state.Limits.OptimizationCycles {
		comparisons := make([]PerformanceComparison, 0, len(prepared.Packet.Targets))
		for _, target := range prepared.Packet.Targets {
			value := current[target.ID]
			comparison, compareErr := ComparePerformanceTarget(target, baselineMeasurement(baseline, target.ID).Qualification.Median, value.Qualification.Median, value.Qualification.Qualified, false)
			if compareErr != nil {
				return out, compareErr
			}
			comparisons = append(comparisons, comparison)
		}
		selected, ok := SelectPerformanceMiss(prepared.Packet.Targets, comparisons)
		if !ok {
			return out, nil
		}
		neededCommands := 1 + len(prepared.CorrectnessCommands) + state.Limits.Warmups + selected.Samples
		for _, result := range performanceAlreadyMet(prepared.Packet.Targets, comparisons, current) {
			neededCommands += state.Limits.Warmups + result.Target.Samples
		}
		// A cycle may start only when the complete final measurement and
		// correctness boundary still fits afterward.
		neededCommands += len(prepared.CorrectnessCommands)
		for _, target := range prepared.Packet.Targets {
			neededCommands += state.Limits.Warmups + target.Samples
		}
		if state.Counters.Commands+neededCommands > state.Limits.Commands {
			state.Counters.OptimizationCycles = state.Limits.OptimizationCycles
			state.NextAction = "Command budget exhausted; finalize without another optimization proposal."
			state.UpdatedAt = s.now().UTC()
			out.Stalled = true
			if _, err := store.PublishState(*state, token); err != nil {
				return out, err
			}
			return out, nil
		}
		cycle := state.Counters.OptimizationCycles + 1
		state.Counters.OptimizationCycles, state.Phase, state.NextAction, state.UpdatedAt = cycle, PerformancePhaseOptimizing, "Evaluate one isolated proposal for "+selected.ID+".", s.now().UTC()
		if _, err := store.PublishState(*state, token); err != nil {
			return out, err
		}
		isolation, err := createPerformanceIsolation(ctx, prepared.Target.Path, prepared.ProtectedRoots, state.Limits)
		if err != nil {
			return out, err
		}
		profileKind, err := PerformanceProfileForUnit(descriptors[selected.ID].RawUnit)
		if err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		profileOutput := filepath.Join(isolation.parent, fmt.Sprintf("profile-%06d.pprof", cycle))
		profileRequest, err := PerformanceProfileRequest(isolation.workspace.Path, profileOutput, descriptors[selected.ID], profileKind, state.Limits, env)
		if err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		state.Counters.Commands++
		profileResult, profileErr := isolation.runner(baseRunner).Run(ctx, profileRequest)
		if profileErr != nil || !profileResult.CleanupComplete || profileResult.StdoutTruncated || profileResult.StderrTruncated {
			cleanupErr := isolation.cleanup()
			return out, errors.Join(fmt.Errorf("performance profile command failed: %w", profileErr), cleanupErr)
		}
		profileBytes := []byte(profileResult.Stdout)
		if info, statErr := os.Lstat(profileOutput); statErr == nil {
			if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() > int64(state.Limits.ProfileBytes) {
				_ = isolation.cleanup()
				return out, fmt.Errorf("performance profile output is unsafe or exceeds its retained-byte limit")
			}
			if retained, readErr := os.ReadFile(profileOutput); readErr == nil {
				profileBytes = retained
			} else {
				_ = isolation.cleanup()
				return out, readErr
			}
		} else if !errors.Is(statErr, fs.ErrNotExist) {
			_ = isolation.cleanup()
			return out, statErr
		}
		if removeErr := os.Remove(profileOutput); removeErr != nil && !errors.Is(removeErr, fs.ErrNotExist) {
			_ = isolation.cleanup()
			return out, removeErr
		}
		if len(profileBytes) == 0 {
			profileBytes = []byte("profile command completed without retained output\n")
		}
		if len(profileBytes) > state.Limits.ProfileBytes {
			_ = isolation.cleanup()
			return out, fmt.Errorf("performance profile exceeds retained-byte limit")
		}
		profileDataPath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "profile.pprof")
		profileDataRef, err := store.PublishBytes(profileDataPath, "cycle-profile-data", profileBytes, true, token)
		if err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		profile := PerformanceProfile{SchemaVersion: PerformanceSchemaVersion, Cycle: cycle, TargetID: selected.ID, Kind: profileKind, Implementation: identity, Digest: profileDataRef.Digest, Bytes: len(profileBytes), Artifact: profileDataRef}
		profilePath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "profile.json")
		if _, err := store.PublishJSON(profilePath, "cycle-profile", profile, true, token); err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		state.Counters.RuntimeAttempts++
		if _, err := store.PublishState(*state, token); err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		prompt := performanceOptimizationPrompt(selected, profile)
		request := s.runtimeRequest(prompt, map[string]string{"project": prepared.Project, "sprint": prepared.Sprint, "stage": string(VerificationPhasePerformance), "operation": "performance-optimization", "cycle": fmt.Sprint(cycle), "target": selected.ID})
		if override, ok := s.verificationRuntime[VerificationPhasePerformance]; ok {
			request.Provider, request.Model = splitProviderModel(override.Model)
			request.Metadata["variant"] = override.Variant
		}
		request.WorkDir, request.Timeout, request.Sandbox, request.Permissions = isolation.workspace.Path, state.Limits.WallTime, "workspace_write", "restricted"
		request.Policy.Default = "deny"
		request.Policy.Tools = map[string]string{"read": "allow", "list": "allow", "search": "allow", "glob": "allow", "write": "allow", "edit": "allow", "patch": "allow", "bash": "deny", "shell": "deny"}
		runtimeResult, runErr := s.runtime.StartRun(ctx, request)
		if cleanupErr := s.deleteCompletedSessions(context.WithoutCancel(ctx), runtimeResult); cleanupErr != nil {
			runErr = errors.Join(runErr, cleanupErr)
		}
		if runErr != nil {
			_ = publishPerformanceRejection(store, state, token, cycle, selected.ID, fmt.Errorf("runtime proposal failed: %w", runErr))
			_ = isolation.cleanup()
			if ctx.Err() != nil {
				return out, ctx.Err()
			}
			if state.Counters.RuntimeAttempts >= state.Limits.OptimizationCycles*state.Limits.OptimizationAttempts {
				out.Stalled = true
				return out, nil
			}
			continue
		}
		changed, diffErr := pprocess.CompareTrees(context.WithoutCancel(ctx), prepared.Target.Path, isolation.workspace.Path, pprocess.IsolationLimits{MaxFiles: performanceIsolationFiles, MaxBytes: performanceIsolationBytes, MaxFileSize: performanceIsolationFileSize, Timeout: state.Limits.CommandTimeout})
		replacements, preimages, changedBytes, deriveErr := derivePerformanceProductionProposal(prepared.Target.Path, isolation.workspace.Path, changed, state.Limits)
		if diffErr != nil || deriveErr != nil {
			_ = publishPerformanceRejection(store, state, token, cycle, selected.ID, errors.Join(diffErr, deriveErr))
			_ = isolation.cleanup()
			continue
		}
		patch := []byte(renderPerformancePatch(prepared.Target.Path, isolation.workspace.Path, changed))
		if len(patch) == 0 || len(patch) > state.Limits.ProposalBytes {
			_ = publishPerformanceRejection(store, state, token, cycle, selected.ID, fmt.Errorf("proposal patch exceeds retained-byte limit"))
			_ = isolation.cleanup()
			continue
		}
		proposalPath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "proposal.patch")
		if _, err := store.PublishBytes(proposalPath, "cycle-proposal", patch, true, token); err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		hypothesis := PerformanceHypothesis{SchemaVersion: PerformanceSchemaVersion, Cycle: cycle, TargetID: selected.ID, Summary: boundRepairText(runtimeResult.TerminalOutput, 512), ProfileDigest: profile.Digest}
		if hypothesis.Summary == "" {
			hypothesis.Summary = "Runtime proposed one target-linked production change."
		}
		hypothesisPath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "hypothesis.json")
		if _, err := store.PublishJSON(hypothesisPath, "cycle-hypothesis", hypothesis, true, token); err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		scope := PerformanceProposalScope{SchemaVersion: PerformanceSchemaVersion, Cycle: cycle, TargetID: selected.ID, Before: identity, Paths: changed, ChangedBytes: changedBytes, Allowed: true}
		candidateIdentity, identityErr := targetIdentity(isolation.workspace.Path)
		if identityErr != nil {
			_ = isolation.cleanup()
			return out, identityErr
		}
		scope.Candidate = candidateIdentity
		scopePath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "scope.json")
		if _, err := store.PublishJSON(scopePath, "cycle-scope", scope, true, token); err != nil {
			_ = isolation.cleanup()
			return out, err
		}
		if err := runPerformanceCorrectness(ctx, isolation.runner(baseRunner), isolation.workspace.Path, prepared.CorrectnessCommands, state.Limits, env, &state.Counters); err != nil {
			_ = publishPerformanceRejection(store, state, token, cycle, selected.ID, fmt.Errorf("candidate correctness failed: %w", err))
			_ = isolation.cleanup()
			continue
		}
		candidateRunner := &performanceMeasurementAccountingRunner{delegate: isolation.runner(baseRunner), counters: &state.Counters, limits: state.Limits}
		candidate, err := MeasurePerformanceTarget(ctx, candidateRunner, isolation.workspace.Path, descriptors[selected.ID], selected, state.Limits, env, prepared.Packet.PacketDigest, manifest.BenchmarkDigest, currentPerformanceEnvironment(env).Identity)
		if _, publishErr := store.PublishState(*state, token); publishErr != nil {
			_ = isolation.cleanup()
			return out, publishErr
		}
		if err != nil {
			_ = publishPerformanceRejection(store, state, token, cycle, selected.ID, err)
			_ = isolation.cleanup()
			continue
		}
		alreadyMet := performanceAlreadyMet(prepared.Packet.Targets, comparisons, current)
		remeasurements := make(map[string]PerformanceMeasurement, len(alreadyMet))
		for i := range alreadyMet {
			target := alreadyMet[i].Target
			remeasurementRunner := &performanceMeasurementAccountingRunner{delegate: isolation.runner(baseRunner), counters: &state.Counters, limits: state.Limits}
			remeasured, measureErr := MeasurePerformanceTarget(ctx, remeasurementRunner, isolation.workspace.Path, descriptors[target.ID], target, state.Limits, env, prepared.Packet.PacketDigest, manifest.BenchmarkDigest, currentPerformanceEnvironment(env).Identity)
			if _, publishErr := store.PublishState(*state, token); publishErr != nil {
				_ = isolation.cleanup()
				return out, publishErr
			}
			if measureErr != nil {
				alreadyMet[i].Comparison.Outcome = PerformanceTargetInconclusive
				alreadyMet[i].FinalQualification = PerformanceQualification{Reason: measureErr.Error()}
				continue
			}
			remeasurements[target.ID] = remeasured
			comparison, compareErr := ComparePerformanceTarget(target, baselineMeasurement(baseline, target.ID).Qualification.Median, remeasured.Qualification.Median, remeasured.Qualification.Qualified, false)
			if compareErr != nil {
				alreadyMet[i].Comparison.Outcome = PerformanceTargetBlocked
				alreadyMet[i].FinalQualification = PerformanceQualification{Reason: compareErr.Error()}
				continue
			}
			alreadyMet[i].Comparison, alreadyMet[i].FinalQualification = comparison, remeasured.Qualification
		}
		if err := ValidatePerformanceRegressionSet(selected, current[selected.ID], candidate, alreadyMet); err != nil {
			_ = publishPerformanceRejection(store, state, token, cycle, selected.ID, err)
			_ = isolation.cleanup()
			continue
		}
		if cleanupErr := isolation.cleanup(); cleanupErr != nil {
			return out, cleanupErr
		}
		currentIdentity, err := targetIdentity(prepared.Target.Path)
		if err != nil || currentIdentity != identity {
			return out, fmt.Errorf("canonical source changed before performance apply")
		}
		operations := make([]RepairApplyOperation, 0, len(changed))
		for index, path := range changed {
			preimage, readErr := os.ReadFile(filepath.Join(prepared.Target.Path, filepath.FromSlash(path)))
			if readErr != nil || hashBytes(preimage) != preimages[path] {
				return out, errors.Join(fmt.Errorf("canonical preimage changed before performance apply"), readErr)
			}
			preimagePath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, fmt.Sprintf("preimage-%06d.bin", index+1))
			preimageRef, publishErr := store.PublishBytes(preimagePath, "cycle-apply-preimage", preimage, true, token)
			if publishErr != nil {
				return out, publishErr
			}
			operations = append(operations, RepairApplyOperation{Path: path, PreimageDigest: preimageRef.Digest, PreimagePath: preimageRef.Path, PostimageDigest: hashBytes(replacements[path])})
		}
		journal := PerformanceApplyJournal{SchemaVersion: PerformanceSchemaVersion, Cycle: cycle, State: "planned", Before: identity, Operations: operations}
		journalPath, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "apply-journal.json")
		journalRef, err := store.PublishJSON(journalPath, "cycle-apply-journal", journal, false, token)
		if err != nil {
			return out, err
		}
		state.ApplyJournal, state.CurrentImplementation, state.UpdatedAt = &journalRef, identity, s.now().UTC()
		if _, err := store.PublishState(*state, token); err != nil {
			return out, err
		}
		applied, _, err := applyRepairFiles(prepared.Target.Path, replacements, preimages, state.Limits.ChangedFiles, state.Limits.ChangedBytes)
		if err != nil {
			journal.State = "uncertain"
			journal.Operations = mergeRepairApplyOperations(applied, operations)
			if ref, publishErr := store.PublishJSON(journalPath, "cycle-apply-journal", journal, false, token); publishErr == nil {
				state.ApplyJournal = &ref
				_, _ = store.PublishState(*state, token)
			}
			return out, err
		}
		identity, err = targetIdentity(prepared.Target.Path)
		if err != nil {
			return out, err
		}
		journal.State, journal.After, journal.Operations = "applied", identity, mergeRepairApplyOperations(applied, operations)
		journalRef, err = store.PublishJSON(journalPath, "cycle-apply-journal", journal, false, token)
		if err != nil {
			return out, err
		}
		current[selected.ID], out.Changed, out.FinalImplementation = candidate, true, identity
		for id, measurement := range remeasurements {
			current[id] = measurement
		}
		state.ApplyJournal, state.CurrentImplementation, state.AcceptedMutations, state.UpdatedAt = &journalRef, identity, state.AcceptedMutations+1, s.now().UTC()
		if _, err := store.PublishState(*state, token); err != nil {
			return out, err
		}
		if err := invalidateAfterPerformanceMutation(s.root, store.sprint); err != nil {
			return out, err
		}
	}
	out.Stalled = true
	return out, nil
}

type performanceMeasurementAccountingRunner struct {
	delegate    pprocess.Runner
	counters    *PerformanceCounters
	limits      PerformanceLimits
	invocations int
}

func (r *performanceMeasurementAccountingRunner) Run(ctx context.Context, request pprocess.Request) (pprocess.Result, error) {
	if r.counters.Commands >= r.limits.Commands {
		return pprocess.Result{}, fmt.Errorf("performance command limit exhausted")
	}
	r.counters.Commands++
	if r.invocations < r.limits.Warmups {
		r.counters.Warmups++
	} else {
		r.counters.Samples++
	}
	r.invocations++
	return r.delegate.Run(ctx, request)
}

func performanceOptimizationPrompt(target PerformanceTarget, profile PerformanceProfile) string {
	return fmt.Sprintf("Optimize exactly target %s (%s, %s %s %s) using retained profile %s. Modify approved production source only. Do not change tests, benchmarks, fixtures, configuration, governed inputs, evidence, repository controls, or Git data. Do not run shell or Git commands. Product code will derive the actual patch, run frozen correctness first, remeasure, and decide promotion.\n", target.ID, target.Metric, target.Comparator, target.Value, target.Unit, profile.Digest)
}
func baselineMeasurement(values []PerformanceMeasurement, id string) PerformanceMeasurement {
	for _, value := range values {
		if value.Target.ID == id {
			return value
		}
	}
	return PerformanceMeasurement{}
}
func performanceAlreadyMet(targets []PerformanceTarget, comparisons []PerformanceComparison, current map[string]PerformanceMeasurement) []PerformanceTargetResult {
	byID := map[string]PerformanceComparison{}
	for _, c := range comparisons {
		byID[c.TargetID] = c
	}
	var out []PerformanceTargetResult
	for _, t := range targets {
		if byID[t.ID].Outcome == PerformanceTargetMet {
			out = append(out, PerformanceTargetResult{Target: t, FinalQualification: current[t.ID].Qualification, Comparison: byID[t.ID]})
		}
	}
	return out
}

func derivePerformanceProductionProposal(canonical, isolated string, paths []string, limits PerformanceLimits) (map[string][]byte, map[string]string, int64, error) {
	normalized, err := ValidatePerformanceOptimizationScope(paths, 1, limits)
	if err != nil && !strings.Contains(err.Error(), "byte limit") {
		return nil, nil, 0, err
	}
	replacements, preimages := map[string][]byte{}, map[string]string{}
	var changedBytes int64
	for _, rel := range normalized {
		beforePath, afterPath := filepath.Join(canonical, filepath.FromSlash(rel)), filepath.Join(isolated, filepath.FromSlash(rel))
		before, beforeErr := os.ReadFile(beforePath)
		after, afterErr := os.ReadFile(afterPath)
		if errors.Is(beforeErr, fs.ErrNotExist) {
			return nil, nil, 0, fmt.Errorf("optimization cannot create production file %q", rel)
		}
		if beforeErr != nil || afterErr != nil {
			return nil, nil, 0, errors.Join(beforeErr, afterErr)
		}
		if err := ensureRepairRegularPath(canonical, beforePath); err != nil {
			return nil, nil, 0, err
		}
		if err := ensureRepairRegularPath(isolated, afterPath); err != nil {
			return nil, nil, 0, err
		}
		changedBytes += int64(len(before) + len(after))
		if changedBytes > limits.ChangedBytes {
			return nil, nil, 0, fmt.Errorf("optimization changed bytes exceed limit")
		}
		replacements[rel] = after
		preimages[rel] = hashBytes(before)
	}
	if _, err := ValidatePerformanceOptimizationScope(normalized, changedBytes, limits); err != nil {
		return nil, nil, 0, err
	}
	return replacements, preimages, changedBytes, nil
}
func renderPerformancePatch(canonical, isolated string, paths []string) string {
	var b strings.Builder
	for _, rel := range paths {
		before, _ := os.ReadFile(filepath.Join(canonical, filepath.FromSlash(rel)))
		after, _ := os.ReadFile(filepath.Join(isolated, filepath.FromSlash(rel)))
		writeWholeFilePatch(&b, rel, before, after)
	}
	return b.String()
}
func publishPerformanceRejection(store PerformanceStore, state *PerformanceState, token VerificationWriterToken, cycle int, target string, cause error) error {
	record := map[string]any{"schema_version": PerformanceSchemaVersion, "cycle": cycle, "target_id": target, "accepted": false, "reason": displayPerformanceCause(cause)}
	path, _ := PerformanceCycleRelPath(store.sprint, state.AttemptID, cycle, "rejection.json")
	_, err := store.PublishJSON(path, "cycle-rejection", record, true, token)
	return err
}
func invalidateAfterPerformanceMutation(root string, sp Sprint) error {
	flow, err := LoadFlowState(root, sp)
	if err != nil {
		return err
	}
	if flow.Review != nil {
		flow.Review.Stale = true
	}
	if flow.Smoke != nil {
		flow.Smoke.Stale = true
	}
	if flow.QA != nil {
		flow.QA.Fresh = false
	}
	if flow.Repair != nil {
		flow.Repair.Fresh = false
	}
	return SaveFlowState(root, sp, flow)
}
