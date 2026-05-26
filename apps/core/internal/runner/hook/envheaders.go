package hook

// EnvironmentHeaders applies plan.Environment.Headers as defaults
// (only if the request doesn't already have that header set).
type EnvironmentHeaders struct{}

func (h *EnvironmentHeaders) BeforeRequest(rctx *RequestContext) error {
	for k, v := range rctx.Plan.Environment.Headers {
		if rctx.Request.Header.Get(k) == "" {
			rctx.Request.Header.Set(k, v)
		}
	}
	return nil
}
