package spec

import "strings"

func Classify(method, path, opID string) string {
	lower := strings.ToLower(opID + " " + path)
	switch method {
	case "GET":
		if isListLike(path, lower) {
			return "list"
		}
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	}
	return "unknown"
}

func isListLike(path, lower string) bool {
	if !strings.HasSuffix(path, "}") {
		return true
	}
	for _, kw := range []string{"list", "search", "query", "find"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func isDestructive(method, opID string) bool {
	if method == "DELETE" {
		return true
	}
	if method == "GET" {
		return false
	}
	destructiveKeywords := []string{"approve", "submit", "send", "publish", "pay", "cancel", "recalculate"}
	lower := strings.ToLower(opID)
	for _, kw := range destructiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	if method == "POST" || method == "PUT" || method == "PATCH" {
		return true
	}
	return false
}
