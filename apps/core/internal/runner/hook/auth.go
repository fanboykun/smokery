package hook

// AuthInjector injects auth headers based on plan.Auth.
type AuthInjector struct{}

func (h *AuthInjector) BeforeRequest(rctx *RequestContext) error {
	auth := rctx.Plan.Auth
	if auth == nil {
		return nil
	}
	switch auth.Type {
	case "bearer":
		if token, ok := auth.Config["token"].(string); ok {
			rctx.Request.Header.Set("Authorization", "Bearer "+token)
		}
	case "basic":
		user, _ := auth.Config["username"].(string)
		pass, _ := auth.Config["password"].(string)
		rctx.Request.SetBasicAuth(user, pass)
	case "api_key":
		key, _ := auth.Config["key"].(string)
		value, _ := auth.Config["value"].(string)
		if key != "" {
			rctx.Request.Header.Set(key, value)
		}
	}
	return nil
}
