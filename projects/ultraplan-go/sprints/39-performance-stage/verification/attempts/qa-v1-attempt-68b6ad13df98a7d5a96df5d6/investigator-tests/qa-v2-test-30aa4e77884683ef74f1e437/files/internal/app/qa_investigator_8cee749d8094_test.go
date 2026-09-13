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

type qaInvestigator8ceeFailWriter struct{}

func (qaInvestigator8ceeFailWriter) Write([]byte) (int, error) {
	return 0, errors.New("injected stdout failure")
}

func TestQAInvestigator_8cee749d8094(t *testing.T) {
	const (
		matcher   = "ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_8cee749d8094"
		assertion = "For qa-v1-theory-3ba88b4e640cc758bd4724ff, assert the returned error for JSON and text modes when stdout fails, and separately assert whether a simultaneous operation error preserves both the operation and output failures."
	)

	rootPath := t.TempDir()
	writeCommandSprintProject(t, rootPath, "proj", "01-alpha")
	base := filepath.Join(rootPath, "projects", "proj")
	writeFixtureFileContent(t, base, commandProjectIndex(t), "project-index.md")
	writeFixtureFileContent(t, base, commandValidRequirements(), "sprints", "01-alpha", "requirements.md")
	writeFixtureFileContent(t, base, commandValidSprintIndex(), "sprints", "01-alpha", "sprint-index.md")
	root := workspace.Root{Path: rootPath}
	service := sprint.NewService(rootPath)

	run := func(stdout io.Writer, args ...string) error {
		return runSprintPerformance(dependencies{
			stdout:  stdout,
			stderr:  io.Discard,
			ctx:     context.Background(),
			workDir: rootPath,
			env:     map[string]string{},
		}, root, service, "proj", "01-alpha", args)
	}

	// Passing controls establish that the identical runtime-free result succeeds
	// and produces the expected representation with a healthy writer.
	var jsonControl bytes.Buffer
	if err := run(&jsonControl, "--dry-run", "--json"); err != nil {
		t.Fatalf("JSON passing control returned error: %v", err)
	}
	if got := jsonControl.String(); !strings.Contains(got, `"operation":"sprint.performance.dry-run"`) || !strings.Contains(got, `"status":"ok"`) {
		t.Fatalf("JSON passing control produced unexpected output: %q", got)
	}
	var textControl bytes.Buffer
	if err := run(&textControl, "--dry-run"); err != nil {
		t.Fatalf("text passing control returned error: %v", err)
	}
	if got := textControl.String(); !strings.Contains(got, "Performance dry-run: proj/01-alpha\n") || !strings.Contains(got, "Policy: disabled\n") {
		t.Fatalf("text passing control produced unexpected output: %q", got)
	}

	var defects []string
	if err := run(qaInvestigator8ceeFailWriter{}, "--dry-run", "--json"); err == nil || !strings.Contains(err.Error(), "injected stdout failure") {
		defects = append(defects, "JSON stdout failure was not returned")
	}
	if err := run(qaInvestigator8ceeFailWriter{}, "--dry-run"); err == nil || !strings.Contains(err.Error(), "injected stdout failure") {
		defects = append(defects, "text stdout failure was not returned")
	}

	// A missing project deterministically supplies an operation error. The
	// failing-writer result must retain both that error and the write failure.
	runMissing := func(stdout io.Writer) error {
		return runSprintPerformance(dependencies{
			stdout: stdout, stderr: io.Discard, ctx: context.Background(), workDir: rootPath, env: map[string]string{},
		}, root, service, "missing-project", "01-alpha", []string{"--dry-run", "--json"})
	}
	var operationOutput bytes.Buffer
	operationErr := runMissing(&operationOutput)
	if operationErr == nil || strings.Contains(operationErr.Error(), "injected stdout failure") {
		t.Fatalf("operation-error passing control did not isolate the operation failure: %v", operationErr)
	}
	combinedErr := runMissing(qaInvestigator8ceeFailWriter{})
	if combinedErr == nil || !strings.Contains(combinedErr.Error(), operationErr.Error()) || !strings.Contains(combinedErr.Error(), "injected stdout failure") {
		defects = append(defects, "simultaneous operation and stdout failures did not preserve both causes")
	}

	if len(defects) > 0 {
		t.Fatalf("%s assertion=%s observed=%s", matcher, assertion, strings.Join(defects, "; "))
	}
}
