package tui

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/app"
	"github.com/Antonio7098/ultraplan-go/internal/sprint"
	"github.com/Antonio7098/ultraplan-go/internal/web"
)

type qaPerformanceEvidenceQueries struct {
	store     sprint.PerformanceStore
	attemptID string
	limits    []int
}

func (q *qaPerformanceEvidenceQueries) Dashboard(context.Context) (app.WebDashboardResult, error) {
	return app.WebDashboardResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Projects(context.Context) (app.WebProjectsResult, error) {
	return app.WebProjectsResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Project(context.Context, string) (app.WebProjectResult, error) {
	return app.WebProjectResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Sprint(context.Context, string, string) (app.WebSprintResult, error) {
	return app.WebSprintResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Studies(context.Context) (app.WebStudiesResult, error) {
	return app.WebStudiesResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Study(context.Context, string) (app.WebStudyResult, error) {
	return app.WebStudyResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Validations(context.Context, string, string) (app.WebValidationResult, error) {
	return app.WebValidationResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) Artifact(context.Context, string) (app.WebArtifactPreview, error) {
	return app.WebArtifactPreview{}, nil
}

func (q *qaPerformanceEvidenceQueries) Health(context.Context) (app.WebHealthResult, error) {
	return app.WebHealthResult{}, nil
}

func (q *qaPerformanceEvidenceQueries) PerformanceStatus(context.Context, app.PerformanceRequest) (app.PerformanceStatusResult, error) {
	return app.PerformanceStatusResult{PerformanceAttemptID: q.attemptID}, nil
}

func (q *qaPerformanceEvidenceQueries) PerformanceResult(context.Context, app.PerformanceRequest) (app.PerformanceStatusResult, error) {
	return app.PerformanceStatusResult{PerformanceAttemptID: q.attemptID}, nil
}

func (q *qaPerformanceEvidenceQueries) PerformanceEvidence(_ context.Context, req app.PerformanceRequest) (app.PerformanceEvidenceResult, error) {
	q.limits = append(q.limits, req.Limit)
	page, err := q.store.EvidencePage(q.attemptID, req.Cursor, req.Filter, req.Limit)
	if err != nil {
		return app.PerformanceEvidenceResult{}, err
	}
	return app.PerformanceEvidenceResult{SchemaVersion: 1, Records: page.Records, NextCursor: page.NextCursor, ReturnedCount: len(page.Records)}, nil
}

func TestQAInvestigator_2d93974a2b11(t *testing.T) {
	root := t.TempDir()
	sp := sprint.Sprint{Project: "alpha", Slug: "39-performance", Path: filepath.Join(root, "projects", "alpha", "sprints", "39-performance")}
	if err := os.MkdirAll(sp.Path, 0o755); err != nil {
		t.Fatal(err)
	}
	attemptID := "performance-v1-attempt-0123456789abcdef01234567"
	token := sprint.VerificationWriterToken{RunID: "run", OperationalAttemptID: "attempt", FencingGeneration: 1}
	store := sprint.NewPerformanceStore(root, sp).WithWriterFence(func(sprint.VerificationWriterToken) error { return nil })
	rel, err := sprint.PerformanceAttemptRelPath(sp, attemptID, "summary.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.PublishBytes(rel, "qa-investigator", []byte(`{"status":"retained"}`), true, token); err != nil {
		t.Fatal(err)
	}

	queries := &qaPerformanceEvidenceQueries{store: store, attemptID: attemptID}
	handler, err := web.NewHandler(web.HandlerOptions{Queries: queries, Authority: "127.0.0.1:8080", Diagnostics: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	request := func(path string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Host = "127.0.0.1:8080"
		handler.ServeHTTP(res, req)
		return res
	}

	control := request("/api/v1/projects/alpha/sprints/39-performance/performance/evidence?limit=1")
	if control.Code != http.StatusOK || len(queries.limits) != 1 || queries.limits[0] != 1 {
		t.Fatalf("passing control failed: status=%d limits=%v body=%s", control.Code, queries.limits, control.Body.String())
	}

	omitted := request("/api/v1/projects/alpha/sprints/39-performance/performance/evidence")
	if omitted.Code != http.StatusOK {
		t.Log("assertion: an omitted evidence limit must use a bounded product default and return HTTP 200")
		t.Logf("actual_status=%d forwarded_limit=%d response=%s", omitted.Code, queries.limits[len(queries.limits)-1], omitted.Body.String())
		t.Fatal("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_2d93974a2b11")
	}
}
