package main

import (
	"embed"
	"log/slog"
	"net/http"
	"upptime-api/routes"
)

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
