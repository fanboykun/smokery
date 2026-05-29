package hook

import "github.com/rs/zerolog/log"

// AuthInjector injects auth headers based on plan.Auth.
type AuthInjector struct{}

func (h *AuthInjector) BeforeRequest(rctx *RequestContext) error {
	auth := rctx.Plan.Auth
	if auth == nil {
		log.Debug().Msg("auth: no auth profile in plan")
		return nil
	}
	log.Debug().Str("type", auth.Type).Str("name", auth.Name).Interface("config_keys", configKeys(auth.Config)).Msg("auth: injecting")

	switch auth.Type {
	case "bearer":
		if token, ok := auth.Config["token"].(string); ok && token != "" {
			rctx.Request.Header.Set("Authorization", "Bearer "+token)
		} else {
			log.Warn().Msg("auth: bearer type but no 'token' key in config")
		}
	case "basic":
		user, _ := auth.Config["username"].(string)
		pass, _ := auth.Config["password"].(string)
		rctx.Request.SetBasicAuth(user, pass)
	case "api_key":
		header, _ := auth.Config["header"].(string)
		value, _ := auth.Config["value"].(string)
		// Also support "key" for backward compat
		if header == "" {
			header, _ = auth.Config["key"].(string)
		}
		if header != "" && value != "" {
			rctx.Request.Header.Set(header, value)
		} else {
			log.Warn().Str("header", header).Bool("has_value", value != "").Msg("auth: api_key type but missing header/value")
		}
	case "custom":
		// Custom: treat all config entries as header key-value pairs
		for k, v := range auth.Config {
			if s, ok := v.(string); ok {
				rctx.Request.Header.Set(k, s)
			}
		}
	}
	return nil
}

func configKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
