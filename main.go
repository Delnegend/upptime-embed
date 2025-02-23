package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
)

type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusDegrade Status = "degrade"
	StatusUnknown Status = "unknown"
)

type UptimeDetails struct {
	IconUrl string `json:"iconUrl,omitempty"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Url     string `json:"url,omitempty"`
	Status  Status `json:"status"`

	ResponseOverall string `json:"responseOverall,omitempty"`
	Response24h     string `json:"response24h,omitempty"`
	Response7d      string `json:"response7d,omitempty"`
	Response30d     string `json:"response30d,omitempty"`
	Response1y      string `json:"response1y,omitempty"`

	// 0-100
	UptimeOverall string `json:"uptimeOverall,omitempty"`
	Uptime24h     string `json:"uptime24h,omitempty"`
	Uptime7d      string `json:"uptime7d,omitempty"`
	Uptime30d     string `json:"uptime30d,omitempty"`
	Uptime1y      string `json:"uptime1y,omitempty"`
}

var (
	rePublicUrl = regexp.MustCompile(`\| <img.+?src="(?<icon_url>.+?)".+?\[(?<title>.+?)\]\((?<url>.+?)\).+?(?<status>🟩|🟥|🟨).+?src="(?<graph_week>.+?)".+?time (?<resp_overall>.+?)".+?time (?<resp_24h>\d+?)".+?time (?<resp_7d>\d+?)".+?time (?<resp_30d>\d+?)".+?time (?<resp_1y>\d+?)".+?uptime (?<up_all>.+?%)".+?uptime (?<up_24h>.+?%)".+?uptime (?<up_7d>.+?%)".+?uptime (?<up_30d>.+?%)".+?uptime (?<up_1y>.+?%)"`)

	rePrivateUrl = regexp.MustCompile(`\| <img.+?src="(?<icon_url>.+?)".+?> (?<title>.+?) \|.+?(?<status>🟩|🟥|🟨).+?src="(?<graph_week>.+?)".+?time (?<resp_overall>.+?)".+?time (?<resp_24h>\d+?)".+?time (?<resp_7d>\d+?)".+?time (?<resp_30d>\d+?)".+?time (?<resp_1y>\d+?)".+?uptime (?<up_all>.+?%)".+?uptime (?<up_24h>.+?%)".+?uptime (?<up_7d>.+?%)".+?uptime (?<up_30d>.+?%)".+?uptime (?<up_1y>.+?%)"`)

	seats = make(chan struct{}, 20)
)

func extractDetails(str string) []UptimeDetails {
	lines := strings.Split(str, "\n")
	details := map[int]UptimeDetails{}
	var wg sync.WaitGroup

	for i, line := range lines {
		wg.Add(1)
		seats <- struct{}{}
		go func(i int, line string) {
			defer wg.Done()
			defer func() { <-seats }()

			var match []string
			var get func(label string) string
			if m := rePublicUrl.FindStringSubmatch(line); m != nil {
				match, get = m, func(label string) string {
					if i := rePublicUrl.SubexpIndex(label); i != -1 {
						return match[i]
					}
					return ""
				}
			} else if m := rePrivateUrl.FindStringSubmatch(line); m != nil {
				match, get = m, func(label string) string {
					if i := rePrivateUrl.SubexpIndex(label); i != -1 {
						return match[i]
					}
					return ""
				}
			} else {
				return
			}

			details[i] = UptimeDetails{
				IconUrl: get("icon_url"),
				Title:   get("title"),
				Slug:    strings.Split(get("graph_week"), "/")[2],
				Url:     get("url"),
				Status: func() Status {
					switch get("status") {
					case "🟩":
						return StatusUp
					case "🟥":
						return StatusDown
					case "🟨":
						return StatusDegrade
					default:
						return StatusUnknown
					}
				}(),

				ResponseOverall: get("resp_overall"),
				Response24h:     get("resp_24h"),
				Response7d:      get("resp_7d"),
				Response30d:     get("resp_30d"),
				Response1y:      get("resp_1y"),

				UptimeOverall: get("up_all"),
				Uptime24h:     get("up_24h"),
				Uptime7d:      get("up_7d"),
				Uptime30d:     get("up_30d"),
				Uptime1y:      get("up_1y"),
			}
		}(i, line)
	}
	wg.Wait()

	result := make([]UptimeDetails, 0, len(details))
	for _, k := range func() []int {
		keys := make([]int, 0, len(details))
		for k := range details {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		return keys
	}() {
		result = append(result, details[k])
	}

	return result
}

type Overall string

const (
	OverallAllSystemsOperational Overall = "all_good"
	OverallDegradedPerformance   Overall = "degraded"
	OverallCompleteOutage        Overall = "down"
	OverallPartialOutage         Overall = "partial"
	OverallUnknownn              Overall = "unknown"
)

func extractOverall(str string) Overall {
	scoped := str[strings.Index(str, "<!--live status--> ")+len("<!--live status--> "):]
	cleaned := strings.Trim(scoped, `*`)
	cleaned = strings.TrimSpace(cleaned)
	switch cleaned {
	case "🟩 All systems operational":
		return OverallAllSystemsOperational
	case "🟨 Degraded performance":
		return OverallDegradedPerformance
	case "🟥 Complete outage":
		return OverallCompleteOutage
	case "🟧 Partial outage":
		return OverallPartialOutage
	default:
		return OverallUnknownn
	}
}

func main() {
	http.HandleFunc("GET /api/alive", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("GET /api/{username}/{reponame}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		reponame := r.PathValue("reponame")
		contentURL := "https://raw.githubusercontent.com/" + username + "/" + reponame + "/refs/heads/master/README.md"

		resp, err := http.Get(contentURL)
		if err != nil {
			w.Write([]byte("ERROR"))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't read response body: " + err.Error()))
			return
		}

		bodyString := string(body)
		firstLine := bodyString[:strings.Index(bodyString, "\n")]

		a := strings.Index(bodyString, "<!--start: status pages-->")
		b := strings.Index(bodyString, "<!--end: status pages-->")
		details := extractDetails(bodyString[a:b])

		respJson, err := json.Marshal(struct {
			Overall Overall         `json:"overall"`
			Details []UptimeDetails `json:"details"`
		}{
			Overall: extractOverall(firstLine),
			Details: details,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't marshal json: " + err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	})

	http.HandleFunc("GET /api/graph/{username}/{reponame}/{slug}/{duration}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		reponame := r.PathValue("reponame")
		slug := r.PathValue("slug")
		duration := r.PathValue("duration")

		for _, d := range []string{"day", "week", "month", "year"} {
			if d == duration {
				break
			}
			http.Error(w, "Invalid duration", http.StatusBadRequest)
		}

		assetUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/heads/master/graphs/%s/response-time-%s.png", username, reponame, slug, duration)

		resp, err := http.Get(assetUrl)
		if err != nil {
			w.Write([]byte("ERROR"))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte("ERROR"))
			return
		}

		w.Write(body)
	})

	http.ListenAndServe(":8080", nil)
}
