package hook

import "github.com/fanboykun/smokery/apps/core/internal/model"

// TraceExtractor pulls X-Trace-Id and X-Request-Id from the response.
type TraceExtractor struct {
	TraceHeader   string // default "X-Trace-Id"
	RequestHeader string // default "X-Request-Id"
}

func (h *TraceExtractor) AfterResponse(rctx *ResponseContext, step *model.StepResult) error {
	traceHdr := h.TraceHeader
	if traceHdr == "" {
		traceHdr = "X-Trace-Id"
	}
	reqHdr := h.RequestHeader
	if reqHdr == "" {
		reqHdr = "X-Request-Id"
	}
	step.Response.TraceID = rctx.Response.Header.Get(traceHdr)
	step.Response.RequestID = rctx.Response.Header.Get(reqHdr)
	return nil
}
