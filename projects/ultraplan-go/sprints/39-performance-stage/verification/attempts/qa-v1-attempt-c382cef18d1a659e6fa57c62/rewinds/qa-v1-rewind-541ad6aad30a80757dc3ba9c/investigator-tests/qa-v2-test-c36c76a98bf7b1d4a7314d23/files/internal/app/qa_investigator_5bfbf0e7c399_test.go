package app

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/platform/config"
	"github.com/Antonio7098/ultraplan-go/internal/workspace"
)

func TestQAInvestigator_5bfbf0e7c399(t *testing.T) {
	effective, err := config.Load(config.LoadOptions{Env: func(string) string { return "" }})
	if err != nil {
		t.Fatalf("control: production config load failed: %v", err)
	}
	if _, _, err := performanceSettings(effective); err != nil {
		t.Fatalf("control: valid effective settings failed: %v", err)
	}

	root := workspace.Root{Path: t.TempDir()}
	deps := dependencies{ctx: context.Background(), env: map[string]string{}}
	controlEmitted := false
	controlRunner := sharedOperationRunner(deps, root, effective, dashboardUseCases{root: root.Path})
	_, controlErr := controlRunner(context.Background(), OperationRequest{Kind: OperationKind("qa-investigator-control")}, func(OperationEvent) {
		controlEmitted = true
	})
	if controlErr == nil || !strings.Contains(controlErr.Error(), "unsupported runtime operation") || controlEmitted {
		t.Fatalf("control: runner dispatch result err=%v emitted=%t", controlErr, controlEmitted)
	}

	effective.Config.Performance.CommandTimeout = "not-a-duration"
	settingsErr := func() error {
		_, _, err := performanceSettings(effective)
		return err
	}()
	if settingsErr == nil || !strings.Contains(settingsErr.Error(), "performance.command_timeout") {
		t.Fatalf("fixture: exact effective value did not fail performanceSettings: %v", settingsErr)
	}

	emitted := false
	runner := sharedOperationRunner(deps, root, effective, dashboardUseCases{root: root.Path})
	_, runErr := runner(context.Background(), OperationRequest{Kind: OperationKind("qa-investigator-probe")}, func(OperationEvent) {
		emitted = true
	})
	if emitted {
		t.Fatal("operation emitted an event before configuration rejection")
	}
	entries, err := os.ReadDir(root.Path)
	if err != nil || len(entries) != 0 {
		t.Fatalf("operation side effect check: entries=%d err=%v", len(entries), err)
	}
	if runErr == nil || !strings.Contains(runErr.Error(), "performance.command_timeout") {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_5bfbf0e7c399 assertion: Show whether the runner still executes an operation or whether composition rejects the same configuration before runner invocation. runner returned %v", runErr)
	}
}
