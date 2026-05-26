package hook

import (
	"github.com/fanboykun/smokery/apps/core/internal/assertion"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// AssertionRunner runs all step assertions against the response.
// It populates step.Assertions and marks the step as failed if any assertion fails.
type AssertionRunner struct{}

func (h *AssertionRunner) AfterResponse(rctx *ResponseContext, step *model.StepResult) error {
	body := string(rctx.Body)
	for _, a := range rctx.Step.Assertions {
		ar := assertion.Run(a, rctx.Response.StatusCode, body)
		step.Assertions = append(step.Assertions, ar)
		if !ar.Passed {
			step.Status = "failed"
		}
	}
	return nil
}
