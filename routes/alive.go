package routes

import "net/http"

func Alive(muxer *http.ServeMux) {
	muxer.HandleFunc("GET /api/alive", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
