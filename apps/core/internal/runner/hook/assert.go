package hook

import (
	"encoding/json"

	"github.com/fanboykun/smokery/apps/core/internal/assertion"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// AssertionRunner runs all step assertions against the response.
// It populates step.Assertions and marks the step as failed if any assertion fails.
type AssertionRunner struct{}

func (h *AssertionRunner) AfterResponse(rctx *ResponseContext, step *model.StepResult) error {
	body := string(rctx.Body)
	for _, a := range rctx.Step.Assertions {
		// For schema assertions, inject the step's ResponseSchema as Expected
		if a.Type == "schema" && a.Expected == nil && len(rctx.Step.ResponseSchema) > 0 {
			a.Expected = json.RawMessage(rctx.Step.ResponseSchema)
		}
		ar := assertion.Run(a, rctx.Response.StatusCode, body)
		step.Assertions = append(step.Assertions, ar)
		if !ar.Passed {
			step.Status = "failed"
		}
	}
	return nil
}
