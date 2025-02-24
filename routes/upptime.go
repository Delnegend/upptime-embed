package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"upptime-api/utils"
)

func Upptime(muxer *http.ServeMux) {
	muxer.HandleFunc("GET /api/{username}/{reponame}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		reponame := r.PathValue("reponame")
		contentURL := "https://raw.githubusercontent.com/" + username + "/" + reponame + "/refs/heads/master/README.md"

		resp, err := http.Get(contentURL)
		if err != nil {
			http.Error(w, "Can't fetch README.md", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Can't read README.md", http.StatusInternalServerError)
			return
		}
		bodyString := string(body)

		bodyLines := strings.Split(bodyString, "\n")
		if len(bodyLines) == 0 {
			http.Error(w, "Empty README.md", http.StatusBadRequest)
			return
		}

		a := strings.Index(bodyString, "<!--start: status pages-->")
		b := strings.Index(bodyString, "<!--end: status pages-->")
		details := utils.ExtractDetails(bodyString[a:b])

		respJson, err := json.Marshal(struct {
			Overall utils.Overall         `json:"Overall"`
			Details []utils.UptimeDetails `json:"Details"`
		}{
			Overall: utils.ExtractOverall(bodyString),
			Details: details,
		})
		if err != nil {
			http.Error(w, "Can't marshal json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	})
}
