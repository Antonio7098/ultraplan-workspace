package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestQAInvestigator_6b0e9329abe5(t *testing.T) {
	const theory = "qa-v1-theory-7ef7c4876857c8e4b85d4c0c"

	// Non-performance compatibility control.
	stdout, stderr, status := runForTest([]string{"sprint", "--help"})
	if status != ExitOK || stderr != "" || !strings.Contains(stdout, "ultraplan sprint") {
		t.Fatalf("%s non-performance control: status=%d stdout=%q stderr=%q", theory, status, stdout, stderr)
	}

	workspace := initializedWorkspace(t)
	writeCommandSprintProject(t, workspace, "proj", "01-alpha")
	projectBase := workspace + "/projects/proj"
	sprintBase := projectBase + "/sprints/01-alpha"
	writeFixtureFileContent(t, projectBase, commandProjectIndex(t), "project-index.md")
	writeFixtureFileContent(t, sprintBase, commandValidSprintIndex(), "sprint-index.md")
	writeFixtureFileContent(t, sprintBase, commandValidRequirements(), "requirements.md")

	// Runtime-free success JSON contract.
	stdout, stderr, status = runForTest([]string{"--workspace", workspace, "sprint", "proj", "01", "performance", "--dry-run", "--json"})
	if status != ExitOK || stderr != "" {
		t.Fatalf("%s dry-run separation: status=%d stdout=%q stderr=%q", theory, status, stdout, stderr)
	}
	var envelope struct {
		SchemaVersion int             `json:"schema_version"`
		Operation     string          `json:"operation"`
		Status        string          `json:"status"`
		Result        json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil || envelope.SchemaVersion != 1 || envelope.Operation != "sprint.performance.dry-run" || envelope.Status != "ok" || len(envelope.Result) == 0 {
		t.Fatalf("%s dry-run JSON envelope=%+v err=%v stdout=%q", theory, envelope, err, stdout)
	}

	// Runtime-free text contract.
	stdout, stderr, status = runForTest([]string{"--workspace", workspace, "sprint", "proj", "01", "performance", "status"})
	if status != ExitOK || stderr != "" || !strings.Contains(stdout, "Performance status: proj/01-alpha") || !strings.Contains(stdout, "Policy: disabled") {
		t.Fatalf("%s status text: status=%d stdout=%q stderr=%q", theory, status, stdout, stderr)
	}

	// Operational failures retain a JSON error envelope on stdout and a nonzero exit.
	stdout, _, status = runForTest([]string{"--workspace", workspace, "sprint", "missing", "01", "performance", "status", "--json"})
	var failure map[string]json.RawMessage
	if status != ExitValidation || json.Unmarshal([]byte(stdout), &failure) != nil || failure["error"] == nil || failure["result"] == nil {
		t.Fatalf("%s operational failure contract: status=%d stdout=%q", theory, status, stdout)
	}

	// Product-owned compatibility coverage must lock every semantic exit branch.
	files := []string{"performance_test.go", "sprint_commands_test.go"}
	var productTests strings.Builder
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s read compatibility tests: %v", theory, err)
		}
		productTests.Write(data)
	}
	required := []string{"PerformanceTargetMiss", "PerformanceCancelled", "ExitCancel", "ExitPartial"}
	missing := make([]string, 0, len(required))
	for _, token := range required {
		if !strings.Contains(productTests.String(), token) {
			missing = append(missing, token)
		}
	}
	if len(missing) != 0 {
		const assertion = "Exercise representative performance commands and assert stable JSON envelope fields, text output, stdout/stderr separation, and success, target-miss, cancellation, partial, and operational-failure exit mappings. Name qa-v1-theory-7ef7c4876857c8e4b85d4c0c i"
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_6b0e9329abe5 assertion=%s theory=%s missing=%s", assertion, theory, fmt.Sprint(missing))
	}
}
