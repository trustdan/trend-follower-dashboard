// Package webui exposes an fs.FS with the compiled Svelte app.
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// Sub returns an fs.FS rooted at the "dist" folder.
func Sub() (fs.FS, error) {
	return fs.Sub(dist, "dist")
}
