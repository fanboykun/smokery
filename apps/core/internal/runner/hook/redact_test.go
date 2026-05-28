package hook

import (
	"net/http"
	"testing"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

func TestRedactorDefaultHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	req.Header.Set("Cookie", "session=abc123")
	req.Header.Set("Content-Type", "application/json")

	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("Set-Cookie", "session=xyz")
	resp.Header.Set("X-Request-Id", "req-123")

	rctx := &ResponseContext{
		Request:  req,
		Response: resp,
		Step:     &model.PlannedStep{},
	}
	step := &model.StepResult{}

	h := &Redactor{}
	if err := h.AfterResponse(rctx, step); err != nil {
		t.Fatal(err)
	}

	// Request headers
	if step.Request.Headers["Authorization"] != "[REDACTED]" {
		t.Errorf("expected Authorization redacted, got %q", step.Request.Headers["Authorization"])
	}
	if step.Request.Headers["Cookie"] != "[REDACTED]" {
		t.Errorf("expected Cookie redacted, got %q", step.Request.Headers["Cookie"])
	}
	if step.Request.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type preserved, got %q", step.Request.Headers["Content-Type"])
	}

	// Response headers
	if step.Response.Headers["Set-Cookie"] != "[REDACTED]" {
		t.Errorf("expected Set-Cookie redacted, got %q", step.Response.Headers["Set-Cookie"])
	}
	if step.Response.Headers["X-Request-Id"] != "req-123" {
		t.Errorf("expected X-Request-Id preserved, got %q", step.Response.Headers["X-Request-Id"])
	}
}

func TestRedactorCustomSensitive(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Api-Key", "my-secret-key")
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Accept", "application/json")

	resp := &http.Response{Header: http.Header{}}

	rctx := &ResponseContext{
		Request:  req,
		Response: resp,
		Step:     &model.PlannedStep{},
	}
	step := &model.StepResult{}

	h := &Redactor{Sensitive: []string{"x-api-key", "authorization"}}
	if err := h.AfterResponse(rctx, step); err != nil {
		t.Fatal(err)
	}

	if step.Request.Headers["X-Api-Key"] != "[REDACTED]" {
		t.Errorf("expected X-Api-Key redacted, got %q", step.Request.Headers["X-Api-Key"])
	}
	if step.Request.Headers["Authorization"] != "[REDACTED]" {
		t.Errorf("expected Authorization redacted, got %q", step.Request.Headers["Authorization"])
	}
	if step.Request.Headers["Accept"] != "application/json" {
		t.Errorf("expected Accept preserved, got %q", step.Request.Headers["Accept"])
	}
}

func TestRedactorCaseInsensitive(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("AUTHORIZATION", "Bearer token")

	resp := &http.Response{Header: http.Header{}}

	rctx := &ResponseContext{
		Request:  req,
		Response: resp,
		Step:     &model.PlannedStep{},
	}
	step := &model.StepResult{}

	h := &Redactor{}
	_ = h.AfterResponse(rctx, step)

	// http.Header canonicalizes to "Authorization"
	if step.Request.Headers["Authorization"] != "[REDACTED]" {
		t.Errorf("expected case-insensitive redaction, got %q", step.Request.Headers["Authorization"])
	}
}
