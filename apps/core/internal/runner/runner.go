package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/runner/hook"
)

// Runner executes compiled SmokePlans. It is a pure library: no DB, no
// OpenAPI parser, no orchestration. All extension is via PreProcessor /
// PostProcessor hooks.
type Runner struct {
	opts Options
}

// New constructs a Runner with the provided Options.
// Use DefaultOptions() for the standard built-in hook set.
func New(opts Options) *Runner {
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Runner{opts: opts}
}

// Execute runs the plan and returns a structured RunResult.
func (r *Runner) Execute(ctx context.Context, plan *model.SmokePlan) *model.RunResult {
	start := time.Now()
	result := &model.RunResult{
		RunID:     plan.ID,
		Status:    "passed",
		StartedAt: start,
	}

	r.emit("run_started", map[string]string{"run_id": plan.ID})

	log.Info().Str("run_id", plan.ID).Str("env", plan.Environment.Name).Int("flows", len(plan.FlowPlans)).Int("suites", len(plan.SuitePlans)).Msg("smoke run started")

	vars := make(map[string]any)

	for _, fp := range plan.FlowPlans {
		fr := r.executeFlow(ctx, fp, plan, vars)
		result.Flows = append(result.Flows, fr)
		if fr.Status == "failed" || fr.Status == "error" {
			result.Status = "failed"
		}
	}

	for _, sp := range plan.SuitePlans {
		sr := r.executeSuite(ctx, sp, plan, vars)
		result.Suites = append(result.Suites, sr)
		if sr.Status == "failed" || sr.Status == "error" {
			result.Status = "failed"
		}
	}

	result.FinishedAt = time.Now()
	result.Duration = time.Since(start).Milliseconds()
	log.Info().Str("run_id", plan.ID).Str("status", result.Status).Int64("duration_ms", result.Duration).Msg("smoke run finished")
	r.emit("run_finished", result)
	return result
}

func (r *Runner) executeFlow(ctx context.Context, fp model.FlowPlan, plan *model.SmokePlan, vars map[string]any) model.FlowResult {
	fr := model.FlowResult{FlowID: fp.FlowID, Name: fp.Name, Status: "passed"}
	r.emit("flow.started", map[string]string{"flow_id": fp.FlowID, "name": fp.Name})

	for i := range fp.Steps {
		step := fp.Steps[i]
		r.emit("flow.step.started", map[string]string{"flow_id": fp.FlowID, "step": step.Name})
		sr := r.executeStep(ctx, &step, plan, vars)
		fr.Steps = append(fr.Steps, sr)
		r.emit("flow.step.finished", sr)
		if sr.Status == "failed" || sr.Status == "error" {
			fr.Status = "failed"
			break
		}
	}

	for i := range fp.Cleanup {
		step := fp.Cleanup[i]
		sr := r.executeStep(ctx, &step, plan, vars)
		fr.Cleanup = append(fr.Cleanup, sr)
	}

	r.emit("flow.finished", fr)
	return fr
}

func (r *Runner) executeSuite(ctx context.Context, sp model.SuitePlan, plan *model.SmokePlan, vars map[string]any) model.SuiteResult {
	sr := model.SuiteResult{SuiteID: sp.SuiteID, Name: sp.Name, Status: "passed"}
	r.emit("suite.started", map[string]string{"suite_id": sp.SuiteID, "name": sp.Name})

	for i := range sp.Cases {
		c := sp.Cases[i]
		r.emit("suite.case.started", map[string]string{"suite_id": sp.SuiteID, "operation_id": c.OperationID, "case_type": c.CaseType})
		stepResult := r.executeStep(ctx, &c.Step, plan, vars)
		cr := model.CaseResult{OperationID: c.OperationID, CaseType: c.CaseType, Step: stepResult}
		sr.Cases = append(sr.Cases, cr)
		r.emit("suite.case.result", cr)
		if stepResult.Status == "failed" || stepResult.Status == "error" {
			sr.Status = "failed"
		}
	}

	r.emit("suite.finished", sr)
	return sr
}

func (r *Runner) executeStep(ctx context.Context, step *model.PlannedStep, plan *model.SmokePlan, vars map[string]any) model.StepResult {
	start := time.Now()
	result := model.StepResult{Name: step.Name, Status: "passed"}

	url := plan.Environment.BaseURL + step.Path

	var bodyReader io.Reader
	if step.Body != nil {
		b, _ := json.Marshal(step.Body)
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, step.Method, url, bodyReader)
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		result.Duration = time.Since(start).Milliseconds()
		return result
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range step.Headers {
		req.Header.Set(k, v)
	}

	// Run pre-processors (interpolation, auth, env headers, custom).
	preCtx := &hook.RequestContext{Plan: plan, Step: step, Vars: vars, Request: req}
	for _, p := range r.opts.PreProcessors {
		if err := p.BeforeRequest(preCtx); err != nil {
			result.Status = "error"
			result.Error = err.Error()
			result.Duration = time.Since(start).Milliseconds()
			return result
		}
	}

	// Snapshot the request URL/method into the step result before sending.
	result.Request = model.RequestMeta{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: simpleHeaders(req.Header),
	}

	resp, err := r.opts.HTTPClient.Do(req)
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		result.Duration = time.Since(start).Milliseconds()
		log.Error().Str("step", step.Name).Str("method", step.Method).Str("url", req.URL.String()).Err(err).Msg("step request failed")
		return result
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	result.Response = model.ResponseMeta{
		Status:  resp.StatusCode,
		Headers: simpleHeaders(resp.Header),
	}
	var respJSON any
	if json.Unmarshal(respBody, &respJSON) == nil {
		result.Response.Body = respJSON
	}

	// Run post-processors (capture, trace, redact, assert, custom).
	postCtx := &hook.ResponseContext{
		Plan: plan, Step: step, Vars: vars,
		Request: req, Response: resp, Body: respBody,
	}
	for _, p := range r.opts.PostProcessors {
		if err := p.AfterResponse(postCtx, &result); err != nil {
			result.Status = "error"
			result.Error = err.Error()
			break
		}
	}

	result.Duration = time.Since(start).Milliseconds()
	lvl := log.Debug()
	if result.Status != "passed" {
		lvl = log.Warn()
	}
	lvl.Str("step", step.Name).Str("method", step.Method).Str("url", req.URL.String()).Int("status", resp.StatusCode).Str("result", result.Status).Int64("ms", result.Duration).Msg("step done")
	return result
}

func (r *Runner) emit(t string, data any) {
	if r.opts.EventEmitter == nil {
		return
	}
	r.opts.EventEmitter(Event{Type: t, Data: data})
}

func simpleHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}
