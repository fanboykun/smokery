package hook

import (
	"net/http"
	"strings"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// Default sensitive headers redacted in step request/response metadata.
var DefaultSensitiveHeaders = []string{
	"authorization", "cookie", "set-cookie",
}

// Redactor replaces sensitive header values with [REDACTED] in the step's
// request and response metadata.
type Redactor struct {
	Sensitive []string
}

func (h *Redactor) AfterResponse(rctx *ResponseContext, step *model.StepResult) error {
	sensitive := h.Sensitive
	if sensitive == nil {
		sensitive = DefaultSensitiveHeaders
	}
	step.Request.Headers = redactHeaders(rctx.Request.Header, sensitive)
	step.Response.Headers = redactHeaders(rctx.Response.Header, sensitive)
	return nil
}

func redactHeaders(h http.Header, sensitive []string) map[string]string {
	sset := make(map[string]bool, len(sensitive))
	for _, s := range sensitive {
		sset[strings.ToLower(s)] = true
	}
	out := make(map[string]string, len(h))
	for k, v := range h {
		if sset[strings.ToLower(k)] {
			out[k] = "[REDACTED]"
		} else if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}
