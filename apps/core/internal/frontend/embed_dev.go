//go:build !embed_frontend

package frontend

import "io/fs"

// FS returns nil when the frontend is not embedded (dev mode).
func FS() fs.FS {
	return nil
}
