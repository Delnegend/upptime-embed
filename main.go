package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"time"
	"upptime-api/routes"

	"github.com/lmittmann/tint"
)

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC1123Z,
		}),
	))
}

//go:embed all:frontend/.output/public
var frontend embed.FS

func main() {
	muxer := http.NewServeMux()

	routes.Alive(muxer)
	routes.Upptime(muxer)
	routes.Graph(muxer)
	routes.Spa(muxer, frontend)

	slog.Info("Listening on port 3001")
	http.ListenAndServe(":3001", muxer)
}
