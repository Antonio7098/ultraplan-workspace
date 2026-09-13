package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/sprint"
	"github.com/Antonio7098/ultraplan-go/internal/workspace"
)

type qaInvestigatorFailWriter struct{}

func (qaInvestigatorFailWriter) Write([]byte) (int, error) {
	return 0, errors.New("injected stdout failure")
}

func TestQAInvestigator_e2eec52fa86d(t *testing.T) {
	const marker = "ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_e2eec52fa86d"

	prepare := PerformancePrepareResult{
		SchemaVersion: 1,
		Project:       "proj",
		Sprint:        "01-alpha",
		PolicyMode:    "disabled",
		DryRun:        true,
		NextAction:    "Enable Performance Policy in project-index.md to use this optional phase.",
	}

	// Passing control: the text projection is written completely to a healthy writer.
	var text bytes.Buffer
	renderSprintPerformance(dependencies{stdout: &text}, "dry-run", prepare)
	if got := text.String(); !strings.Contains(got, "Performance dry-run: proj/01-alpha\n") || !strings.Contains(got, "Policy: disabled\n") {
		t.Fatalf("passing text control produced an unstable projection: %q", got)
	}

	rootPath := t.TempDir()
	writeCommandSprintProject(t, rootPath, "proj", "01-alpha")
	base := filepath.Join(rootPath, "projects", "proj")
	writeFixtureFileContent(t, base, commandProjectIndex(t), "project-index.md")
	writeFixtureFileContent(t, base, commandValidRequirements(), "sprints", "01-alpha", "requirements.md")
	writeFixtureFileContent(t, base, commandValidSprintIndex(), "sprints", "01-alpha", "sprint-index.md")
	service := sprint.NewService(rootPath)
	root := workspace.Root{Path: rootPath}

	// Passing control: runtime-free JSON dry-run succeeds with a healthy writer.
	var jsonOutput bytes.Buffer
	goodDeps := dependencies{stdout: &jsonOutput, stderr: io.Discard, ctx: context.Background(), workDir: rootPath, env: map[string]string{}}
	if err := runSprintPerformance(goodDeps, root, service, "proj", "01-alpha", []string{"--dry-run", "--json"}); err != nil {
		t.Fatalf("passing JSON control returned error: %v", err)
	}
	for _, field := range []string{`"schema_version":1`, `"operation":"sprint.performance.dry-run"`, `"status":"ok"`, `"result":`} {
		if !strings.Contains(jsonOutput.String(), field) {
			t.Fatalf("passing JSON control missing %s: %s", field, jsonOutput.String())
		}
	}

	badDeps := goodDeps
	badDeps.stdout = qaInvestigatorFailWriter{}
	jsonErr := runSprintPerformance(badDeps, root, service, "proj", "01-alpha", []string{"--dry-run", "--json"})

	// Text has no error-returning surface, while JSON incorrectly reports success.
	renderSprintPerformance(badDeps, "dry-run", prepare)
	if jsonErr == nil {
		t.Fatalf("%s assertion=Assertions for qa-v1-theory-3ba88b4e640cc758bd4724ff must prove JSON Encode and text rendering write failures affect the returned error. Separate assertions for qa-v1-theory-44d9a1998804c860f04368c4 must verify stable JSON envelope fields, stdout/stderr se", marker)
	}
}
