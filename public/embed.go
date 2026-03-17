package public

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/buffalo"
)

//go:embed assets
//go:embed favicon.svg
//go:embed index.html
//go:embed vite.svg
var files embed.FS

// FS returns the embedded filesystem for public assets
func FS() fs.FS {
	return buffalo.NewFS(files, "public")
}
