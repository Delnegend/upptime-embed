package routes

import (
	"fmt"
	"io"
	"net/http"
)

func Graph(muxer *http.ServeMux) {
	muxer.HandleFunc("GET /api/graph/{username}/{reponame}/{slug}/{duration}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		reponame := r.PathValue("reponame")
		slug := r.PathValue("slug")
		duration := r.PathValue("duration")

		switch duration {
		case "day":
			duration = "-day"
		case "week":
			duration = "-week"
		case "month":
			duration = "-month"
		case "year":
			duration = "-year"
		case "all":
			duration = ""
		default:
			http.Error(w, "Invalid duration", http.StatusBadRequest)
			return
		}

		assetUrl := fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/refs/heads/master/graphs/%s/response-time%s.png",
			username,
			reponame,
			slug,
			duration,
		)

		resp, err := http.Get(assetUrl)
		if err != nil {
			http.Error(w, "Can't fetch graph", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	})
}
