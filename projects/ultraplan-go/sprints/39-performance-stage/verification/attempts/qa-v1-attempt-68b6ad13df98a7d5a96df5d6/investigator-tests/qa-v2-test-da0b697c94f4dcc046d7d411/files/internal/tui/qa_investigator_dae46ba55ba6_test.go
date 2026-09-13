package tui

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/app"
	"github.com/Antonio7098/ultraplan-go/internal/sprint"
	"github.com/Antonio7098/ultraplan-go/internal/web"
)

type qaDAEPerformanceQueries struct {
	store     sprint.PerformanceStore
	attemptID string
	limits    []int
}

func (*qaDAEPerformanceQueries) Dashboard(context.Context) (app.WebDashboardResult, error) {
	return app.WebDashboardResult{}, nil
}
func (*qaDAEPerformanceQueries) Projects(context.Context) (app.WebProjectsResult, error) {
	return app.WebProjectsResult{}, nil
}
func (*qaDAEPerformanceQueries) Project(context.Context, string) (app.WebProjectResult, error) {
	return app.WebProjectResult{}, nil
}
func (*qaDAEPerformanceQueries) Sprint(context.Context, string, string) (app.WebSprintResult, error) {
	return app.WebSprintResult{}, nil
}
func (*qaDAEPerformanceQueries) Studies(context.Context) (app.WebStudiesResult, error) {
	return app.WebStudiesResult{}, nil
}
func (*qaDAEPerformanceQueries) Study(context.Context, string) (app.WebStudyResult, error) {
	return app.WebStudyResult{}, nil
}
func (*qaDAEPerformanceQueries) Validations(context.Context, string, string) (app.WebValidationResult, error) {
	return app.WebValidationResult{}, nil
}
func (*qaDAEPerformanceQueries) Artifact(context.Context, string) (app.WebArtifactPreview, error) {
	return app.WebArtifactPreview{}, nil
}
func (*qaDAEPerformanceQueries) Health(context.Context) (app.WebHealthResult, error) {
	return app.WebHealthResult{}, nil
}
func (q *qaDAEPerformanceQueries) PerformanceStatus(context.Context, app.PerformanceRequest) (app.PerformanceStatusResult, error) {
	return app.PerformanceStatusResult{PerformanceAttemptID: q.attemptID}, nil
}
func (q *qaDAEPerformanceQueries) PerformanceResult(context.Context, app.PerformanceRequest) (app.PerformanceStatusResult, error) {
	return app.PerformanceStatusResult{PerformanceAttemptID: q.attemptID}, nil
}
func (q *qaDAEPerformanceQueries) PerformanceEvidence(_ context.Context, req app.PerformanceRequest) (app.PerformanceEvidenceResult, error) {
	q.limits = append(q.limits, req.Limit)
	page, err := q.store.EvidencePage(q.attemptID, req.Cursor, req.Filter, req.Limit)
	if err != nil {
		return app.PerformanceEvidenceResult{}, err
	}
	return app.PerformanceEvidenceResult{SchemaVersion: 1, Records: page.Records, NextCursor: page.NextCursor, ReturnedCount: len(page.Records)}, nil
}

func TestQAInvestigator_dae46ba55ba6(t *testing.T) {
	root := t.TempDir()
	sp := sprint.Sprint{Project: "alpha", Slug: "39-performance", Path: filepath.Join(root, "projects", "alpha", "sprints", "39-performance")}
	if err := os.MkdirAll(sp.Path, 0o755); err != nil {
		t.Fatal(err)
	}
	attemptID := "performance-v1-attempt-0123456789abcdef01234567"
	store := sprint.NewPerformanceStore(root, sp).WithWriterFence(func(sprint.VerificationWriterToken) error { return nil })
	rel, err := sprint.PerformanceAttemptRelPath(sp, attemptID, "summary.json")
	if err != nil {
		t.Fatal(err)
	}
	token := sprint.VerificationWriterToken{RunID: "run", OperationalAttemptID: "attempt", FencingGeneration: 1}
	if _, err := store.PublishBytes(rel, "qa-investigator", []byte(`{"status":"retained"}`), true, token); err != nil {
		t.Fatal(err)
	}

	queries := &qaDAEPerformanceQueries{store: store, attemptID: attemptID}
	handler, err := web.NewHandler(web.HandlerOptions{Queries: queries, Authority: "127.0.0.1:8080", Diagnostics: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/projects/alpha/sprints/39-performance/performance/evidence", nil)
	request.Host = "127.0.0.1:8080"
	handler.ServeHTTP(response, request)

	if len(queries.limits) != 1 || queries.limits[0] != 0 {
		t.Fatalf("omitted request did not reach PerformanceEvidence and EvidencePage with Limit=0: limits=%v status=%d body=%s", queries.limits, response.Code, response.Body.String())
	}
	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "The request could not be completed.") {
		t.Fatalf("unexpected omitted-limit HTTP result: status=%d body=%s", response.Code, response.Body.String())
	}

	control, controlErr := queries.PerformanceEvidence(context.Background(), app.PerformanceRequest{Project: "alpha", Sprint: "39-performance", Limit: 1})
	if controlErr != nil || control.ReturnedCount != 1 || len(queries.limits) != 2 || queries.limits[1] != 1 {
		t.Fatalf("bounded-success control failed: result=%+v err=%v limits=%v", control, controlErr, queries.limits)
	}

	t.Log("assertion: For qa-v1-theory-f2dfc2c802889fddc3bbae5c, assert the omitted request's HTTP status and response body, record the exact Limit received by PerformanceEvidence, and compare it with a reachable bounded-success invocation.")
	t.Logf("omitted_status=%d omitted_limit=%d omitted_body=%s control_status=success control_limit=%d", response.Code, queries.limits[0], response.Body.String(), queries.limits[1])
	t.Fatal("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_dae46ba55ba6")
}
