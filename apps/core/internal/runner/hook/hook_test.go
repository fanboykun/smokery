package hook

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

func TestAuthInjectorBearer(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rctx := &RequestContext{
		Plan:    &model.SmokePlan{Auth: &model.AuthProfile{Type: "bearer", Config: map[string]any{"token": "abc123"}}},
		Step:    &model.PlannedStep{},
		Request: req,
	}
	h := &AuthInjector{}
	if err := h.BeforeRequest(rctx); err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer abc123" {
		t.Errorf("expected Bearer abc123, got %q", got)
	}
}

func TestAuthInjectorBasic(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rctx := &RequestContext{
		Plan:    &model.SmokePlan{Auth: &model.AuthProfile{Type: "basic", Config: map[string]any{"username": "user", "password": "pass"}}},
		Step:    &model.PlannedStep{},
		Request: req,
	}
	h := &AuthInjector{}
	_ = h.BeforeRequest(rctx)
	u, p, ok := req.BasicAuth()
	if !ok || u != "user" || p != "pass" {
		t.Errorf("expected basic auth user/pass, got %s/%s ok=%v", u, p, ok)
	}
}

func TestAuthInjectorAPIKey(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rctx := &RequestContext{
		Plan:    &model.SmokePlan{Auth: &model.AuthProfile{Type: "api_key", Config: map[string]any{"key": "X-API-Key", "value": "secret"}}},
		Step:    &model.PlannedStep{},
		Request: req,
	}
	h := &AuthInjector{}
	_ = h.BeforeRequest(rctx)
	if got := req.Header.Get("X-API-Key"); got != "secret" {
		t.Errorf("expected secret, got %q", got)
	}
}

func TestAuthInjectorNilAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rctx := &RequestContext{Plan: &model.SmokePlan{}, Step: &model.PlannedStep{}, Request: req}
	h := &AuthInjector{}
	if err := h.BeforeRequest(rctx); err != nil {
		t.Fatal(err)
	}
}

func TestVariableInterpolatorPath(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/users/{id}", nil)
	rctx := &RequestContext{
		Plan:    &model.SmokePlan{},
		Step:    &model.PlannedStep{Params: map[string]any{"id": "42"}},
		Vars:    map[string]any{},
		Request: req,
	}
	h := &VariableInterpolator{}
	_ = h.BeforeRequest(rctx)
	if req.URL.Path != "/users/42" {
		t.Errorf("expected /users/42, got %s", req.URL.Path)
	}
}

func TestVariableInterpolatorBody(t *testing.T) {
	body := `{"user_id":"{{captured_id}}"}`
	req, _ := http.NewRequest("POST", "http://example.com", io.NopCloser(bytes.NewReader([]byte(body))))
	rctx := &RequestContext{
		Plan:    &model.SmokePlan{},
		Step:    &model.PlannedStep{},
		Vars:    map[string]any{"captured_id": "abc-123"},
		Request: req,
	}
	h := &VariableInterpolator{}
	_ = h.BeforeRequest(rctx)
	result, _ := io.ReadAll(req.Body)
	if string(result) != `{"user_id":"abc-123"}` {
		t.Errorf("unexpected body: %s", result)
	}
}

func TestCaptureFromBody(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	rctx := &ResponseContext{
		Step: &model.PlannedStep{Captures: []model.Capture{
			{Name: "user_id", Source: "body", Path: "data.id"},
		}},
		Vars:     map[string]any{},
		Response: resp,
		Body:     []byte(`{"data":{"id":"xyz-789","name":"test"}}`),
	}
	step := &model.StepResult{}
	h := &Capture{}
	_ = h.AfterResponse(rctx, step)
	if step.Captures["user_id"] != "xyz-789" {
		t.Errorf("expected xyz-789, got %v", step.Captures["user_id"])
	}
	if rctx.Vars["user_id"] != "xyz-789" {
		t.Error("expected var to be set")
	}
}

func TestCaptureFromHeader(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("X-Request-Id", "req-abc")
	rctx := &ResponseContext{
		Step: &model.PlannedStep{Captures: []model.Capture{
			{Name: "req_id", Source: "header", Path: "X-Request-Id"},
		}},
		Vars:     map[string]any{},
		Response: resp,
		Body:     []byte(`{}`),
	}
	step := &model.StepResult{}
	h := &Capture{}
	_ = h.AfterResponse(rctx, step)
	if step.Captures["req_id"] != "req-abc" {
		t.Errorf("expected req-abc, got %v", step.Captures["req_id"])
	}
}

func TestTraceExtractor(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("X-Trace-Id", "trace-123")
	resp.Header.Set("X-Request-Id", "req-456")
	rctx := &ResponseContext{
		Step:     &model.PlannedStep{},
		Response: resp,
		Body:     []byte(`{}`),
	}
	step := &model.StepResult{}
	h := &TraceExtractor{}
	_ = h.AfterResponse(rctx, step)
	if step.Response.TraceID != "trace-123" {
		t.Errorf("expected trace-123, got %q", step.Response.TraceID)
	}
	if step.Response.RequestID != "req-456" {
		t.Errorf("expected req-456, got %q", step.Response.RequestID)
	}
}

func TestTraceExtractorCustomHeaders(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("Traceparent", "00-abc-def-01")
	rctx := &ResponseContext{Step: &model.PlannedStep{}, Response: resp, Body: []byte(`{}`)}
	step := &model.StepResult{}
	h := &TraceExtractor{TraceHeader: "Traceparent"}
	_ = h.AfterResponse(rctx, step)
	if step.Response.TraceID != "00-abc-def-01" {
		t.Errorf("expected traceparent value, got %q", step.Response.TraceID)
	}
}

func TestEnvironmentHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Existing", "keep")
	rctx := &RequestContext{
		Plan: &model.SmokePlan{Environment: model.Environment{
			Headers: map[string]string{"X-Source": "smokery", "Existing": "overwrite"},
		}},
		Step:    &model.PlannedStep{},
		Request: req,
	}
	h := &EnvironmentHeaders{}
	_ = h.BeforeRequest(rctx)
	if req.Header.Get("X-Source") != "smokery" {
		t.Errorf("expected X-Source=smokery, got %q", req.Header.Get("X-Source"))
	}
	// Should NOT overwrite existing header
	if req.Header.Get("Existing") != "keep" {
		t.Errorf("expected Existing=keep (not overwritten), got %q", req.Header.Get("Existing"))
	}
}
