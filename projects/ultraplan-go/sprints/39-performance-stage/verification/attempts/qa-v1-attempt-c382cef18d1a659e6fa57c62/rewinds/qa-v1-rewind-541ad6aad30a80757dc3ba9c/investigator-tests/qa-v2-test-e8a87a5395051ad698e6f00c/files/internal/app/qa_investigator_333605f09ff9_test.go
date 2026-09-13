package app

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/platform/config"
	"github.com/Antonio7098/ultraplan-go/internal/workspace"
)

func TestQAInvestigator_333605f09ff9(t *testing.T) {
	effective, err := config.Load(config.LoadOptions{Env: func(string) string { return "" }})
	if err != nil {
		t.Fatalf("control: load default effective configuration: %v", err)
	}
	if _, _, err := performanceSettings(effective); err != nil {
		t.Fatalf("control: default effective performance settings: %v", err)
	}

	effective.Config.Performance.CommandTimeout = "not-a-duration"
	if _, _, err := performanceSettings(effective); err == nil {
		t.Fatal("fixture: invalid performance command timeout was accepted")
	}

	root := workspace.Root{Path: t.TempDir()}
	deps := dependencies{stdout: io.Discard, stderr: io.Discard, ctx: context.Background(), env: map[string]string{}}
	runner := sharedOperationRunner(deps, root, effective, dashboardUseCases{root: root.Path})
	emitted := false
	_, runErr := runner(context.Background(), OperationRequest{Kind: OperationSmokeStart, Project: "missing-project", Sprint: "missing-sprint"}, func(OperationEvent) {
		emitted = true
	})
	if emitted {
		t.Fatal("operation emitted progress before invalid configuration was rejected")
	}
	entries, readErr := os.ReadDir(root.Path)
	if readErr != nil {
		t.Fatalf("inspect operation side effects: %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("operation wrote %d filesystem entries before invalid configuration was rejected", len(entries))
	}
	if runErr == nil || !strings.Contains(runErr.Error(), "performance.config") {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_333605f09ff9 assertion: shared operation runner must reject invalid performance.command_timeout before operation dispatch; got error %v", runErr)
	}
}
