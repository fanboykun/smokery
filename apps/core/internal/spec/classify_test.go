package spec

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		method, path, opID, want string
	}{
		{"GET", "/users", "listUsers", "list"},
		{"GET", "/users/{id}", "getUser", "read"},
		{"POST", "/users", "createUser", "create"},
		{"PUT", "/users/{id}", "updateUser", "update"},
		{"DELETE", "/users/{id}", "deleteUser", "delete"},
		{"GET", "/search", "searchItems", "list"},
	}
	for _, tt := range tests {
		got := Classify(tt.method, tt.path, tt.opID)
		if got != tt.want {
			t.Errorf("Classify(%s, %s, %s) = %s, want %s", tt.method, tt.path, tt.opID, got, tt.want)
		}
	}
}

func TestIsDestructive(t *testing.T) {
	tests := []struct {
		method, opID string
		want         bool
	}{
		{"DELETE", "deleteUser", true},
		{"GET", "getUser", false},
		{"POST", "createUser", true},
		{"POST", "approveOrder", true},
		{"PUT", "updateUser", true},
	}
	for _, tt := range tests {
		got := isDestructive(tt.method, tt.opID)
		if got != tt.want {
			t.Errorf("isDestructive(%s, %s) = %v, want %v", tt.method, tt.opID, got, tt.want)
		}
	}
}
