package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"github.com/fanboykun/smokery/apps/api/internal/assertion"
	"github.com/fanboykun/smokery/apps/api/internal/model"
)

type Runner struct {
	Client *http.Client
}

func New() *Runner {
	return &Runner{Client: &http.Client{Timeout: 30 * time.Second}}
}

func (r *Runner) Execute(ctx context.Context, plan *model.SmokePlan) *model.RunResult {
	start := time.Now()
	result := &model.RunResult{
		RunID:     plan.ID,
		Status:    "passed",
		StartedAt: start,
	}

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
	return result
}

func (r *Runner) executeFlow(ctx context.Context, fp model.FlowPlan, plan *model.SmokePlan, vars map[string]any) model.FlowResult {
	fr := model.FlowResult{FlowID: fp.FlowID, Name: fp.Name, Status: "passed"}

	for _, step := range fp.Steps {
		sr := r.executeStep(ctx, step, plan, vars)
		fr.Steps = append(fr.Steps, sr)
		if sr.Status == "failed" || sr.Status == "error" {
			fr.Status = "failed"
			break
		}
	}

	for _, step := range fp.Cleanup {
		sr := r.executeStep(ctx, step, plan, vars)
		fr.Cleanup = append(fr.Cleanup, sr)
	}

	return fr
}

func (r *Runner) executeSuite(ctx context.Context, sp model.SuitePlan, plan *model.SmokePlan, vars map[string]any) model.SuiteResult {
	sr := model.SuiteResult{SuiteID: sp.SuiteID, Name: sp.Name, Status: "passed"}

	for _, c := range sp.Cases {
		stepResult := r.executeStep(ctx, c.Step, plan, vars)
		sr.Cases = append(sr.Cases, model.CaseResult{
			OperationID: c.OperationID, CaseType: c.CaseType, Step: stepResult,
		})
		if stepResult.Status == "failed" || stepResult.Status == "error" {
			sr.Status = "failed"
		}
	}

	return sr
}

func (r *Runner) executeStep(ctx context.Context, step model.PlannedStep, plan *model.SmokePlan, vars map[string]any) model.StepResult {
	start := time.Now()
	result := model.StepResult{Name: step.Name, Status: "passed"}

	url := plan.Environment.BaseURL + interpolatePath(step.Path, step.Params, vars)

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
		req.Header.Set(k, interpolateString(v, vars))
	}
	for k, v := range plan.Environment.Headers {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}

	if plan.Auth != nil {
		injectAuth(req, plan.Auth)
	}

	result.Request = model.RequestMeta{Method: step.Method, URL: url, Headers: redactHeaders(req.Header)}

	resp, err := r.Client.Do(req)
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		result.Duration = time.Since(start).Milliseconds()
		return result
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	respBodyStr := string(respBody)

	result.Response = model.ResponseMeta{
		Status:    resp.StatusCode,
		Headers:   redactHeaders(resp.Header),
		TraceID:   resp.Header.Get("X-Trace-Id"),
		RequestID: resp.Header.Get("X-Request-Id"),
	}

	var respJSON any
	if json.Unmarshal(respBody, &respJSON) == nil {
		result.Response.Body = respJSON
	}

	if len(step.Captures) > 0 {
		result.Captures = make(map[string]any)
		for _, cap := range step.Captures {
			val := captureValue(cap, respBodyStr, resp.Header)
			if val != nil {
				result.Captures[cap.Name] = val
				vars[cap.Name] = val
			}
		}
	}

	for _, a := range step.Assertions {
		ar := runAssertion(a, resp.StatusCode, respBodyStr)
		result.Assertions = append(result.Assertions, ar)
		if !ar.Passed {
			result.Status = "failed"
		}
	}

	result.Duration = time.Since(start).Milliseconds()
	return result
}

func runAssertion(a model.Assertion, statusCode int, body string) model.AssertionResult {
	return assertion.Run(a, statusCode, body)
}

func injectAuth(req *http.Request, auth *model.AuthProfile) {
	switch auth.Type {
	case "bearer":
		if token, ok := auth.Config["token"].(string); ok {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	case "basic":
		user, _ := auth.Config["username"].(string)
		pass, _ := auth.Config["password"].(string)
		req.SetBasicAuth(user, pass)
	case "api_key":
		key, _ := auth.Config["key"].(string)
		value, _ := auth.Config["value"].(string)
		req.Header.Set(key, value)
	}
}

func captureValue(cap model.Capture, body string, headers http.Header) any {
	switch cap.Source {
	case "body":
		r := gjson.Get(body, cap.Path)
		if r.Exists() {
			return r.Value()
		}
	case "header":
		if v := headers.Get(cap.Path); v != "" {
			return v
		}
	}
	return nil
}

func interpolatePath(path string, params map[string]any, vars map[string]any) string {
	result := path
	for k, v := range params {
		result = strings.ReplaceAll(result, "{"+k+"}", fmt.Sprintf("%v", v))
	}
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{"+k+"}", fmt.Sprintf("%v", v))
	}
	return result
}

func interpolateString(s string, vars map[string]any) string {
	result := s
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", fmt.Sprintf("%v", v))
	}
	return result
}

var sensitiveHeaders = map[string]bool{
	"authorization": true, "cookie": true, "set-cookie": true,
}

func redactHeaders(h http.Header) map[string]string {
	result := make(map[string]string)
	for k, v := range h {
		if sensitiveHeaders[strings.ToLower(k)] {
			result[k] = "[REDACTED]"
		} else if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}
