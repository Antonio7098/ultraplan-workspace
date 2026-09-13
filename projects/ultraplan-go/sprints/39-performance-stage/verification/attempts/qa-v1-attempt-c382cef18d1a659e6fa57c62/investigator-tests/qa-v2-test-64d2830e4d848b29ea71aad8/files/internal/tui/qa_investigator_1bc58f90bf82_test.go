package tui

import (
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/app"
)

func TestQAInvestigator_1bc58f90bf82(t *testing.T) {
	data := fixtureDashboard()
	summary := &data.Sprints[0]
	model := Model{Data: data, Routes: []Route{{Kind: RouteSprint, Project: summary.Project, Sprint: summary.Slug}}}

	performanceKinds := func() map[app.OperationKind]bool {
		kinds := make(map[app.OperationKind]bool)
		for _, item := range model.navItems() {
			if item.Operation != nil && strings.HasPrefix(string(item.Operation.Kind), "performance-") {
				kinds[item.Operation.Kind] = true
			}
		}
		return kinds
	}
	mutating := []app.OperationKind{
		app.OperationPerformanceStart,
		app.OperationPerformanceResume,
		app.OperationPerformanceCancel,
		app.OperationPerformanceRecover,
	}

	// Control: disabled policy without history exposes no performance actions.
	summary.Performance = app.PerformanceStatusResult{PolicyMode: "disabled", Enabled: false}
	model.Data = data
	if kinds := performanceKinds(); len(kinds) != 0 {
		t.Fatalf("disabled-without-attempt control exposed performance operations: %v", kinds)
	}

	// Control: an enabled current attempt exposes all expected mutation actions.
	summary.Performance = app.PerformanceStatusResult{PolicyMode: "enabled", Enabled: true, PerformanceAttemptID: "performance-v1-attempt-current"}
	model.Data = data
	currentKinds := performanceKinds()
	for _, kind := range mutating {
		if !currentKinds[kind] {
			t.Fatalf("enabled-current-attempt control missing %s", kind)
		}
	}

	// Historical data remains inspectable, but disabled policy must hide mutations.
	summary.Performance = app.PerformanceStatusResult{PolicyMode: "disabled", Enabled: false, Phase: "terminal", Fresh: false, Outcome: "passed", PerformanceAttemptID: "performance-v1-attempt-historical"}
	model.Data = data
	historicalKinds := performanceKinds()
	var exposed []string
	for _, kind := range mutating {
		if historicalKinds[kind] {
			exposed = append(exposed, string(kind))
		}
	}
	if len(exposed) != 0 {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_1bc58f90bf82 assertion=Assert start, resume, cancel, and recover operation kinds are present. exposed_disabled_historical_mutations=%s", strings.Join(exposed, ","))
	}
}
