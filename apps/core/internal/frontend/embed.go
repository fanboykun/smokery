//go:build embed_frontend

package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFS embed.FS

// FS returns the embedded frontend static files, or nil if not embedded.
func FS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil
	}
	return sub
}
