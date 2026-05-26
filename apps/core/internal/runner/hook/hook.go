// Package hook defines the runner's pluggable hook contracts and provides
// built-in pre/post processors. The runner depends on this package; this package
// must not import the runner.
package hook

import (
	"net/http"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// RequestContext is the input to a PreProcessor.
type RequestContext struct {
	Plan    *model.SmokePlan
	Step    *model.PlannedStep
	Vars    map[string]any
	Request *http.Request
}

// ResponseContext is the input to a PostProcessor.
type ResponseContext struct {
	Plan     *model.SmokePlan
	Step     *model.PlannedStep
	Vars     map[string]any
	Request  *http.Request
	Response *http.Response
	Body     []byte
}

// PreProcessor runs before each HTTP request is sent.
type PreProcessor interface {
	BeforeRequest(rctx *RequestContext) error
}

// PostProcessor runs after each HTTP response is received.
type PostProcessor interface {
	AfterResponse(rctx *ResponseContext, step *model.StepResult) error
}
