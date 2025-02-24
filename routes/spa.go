package routes

import (
	"embed"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func Spa(muxer *http.ServeMux, frontend embed.FS) {
	muxer.HandleFunc("GET /{filepath...}", func(w http.ResponseWriter, r *http.Request) {
		pathStr := filepath.Clean(r.PathValue("filepath"))
		switch pathStr {
		case ".":
			pathStr = "index.html"
		case "200":
			pathStr = "200.html"
		case "404":
			pathStr = "404.html"
		}
		pathStr = filepath.Join("frontend/.output/public", pathStr)

		file, err := frontend.Open(strings.ReplaceAll(pathStr, "\\", "/"))
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		stats, err := file.Stat()
		if err != nil {
			http.Error(w, "Can't stat file", http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, pathStr, stats.ModTime(), file.(io.ReadSeeker))
	})
}
