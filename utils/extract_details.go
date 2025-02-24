package utils

import (
	"regexp"
	"strings"
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

	UptimeOverall string `json:"uptimeOverall,omitempty"`
	Uptime24h     string `json:"uptime24h,omitempty"`
	Uptime7d      string `json:"uptime7d,omitempty"`
	Uptime30d     string `json:"uptime30d,omitempty"`
	Uptime1y      string `json:"uptime1y,omitempty"`
}

var (
	rePublicUrl = regexp.MustCompile(`\| <img.+?src="(?<icon_url>.+?)".+?\[(?<title>.+?)\]\((?<url>.+?)\).+?(?<status>🟩|🟥|🟨).+?src="(?<graph_week>.+?)".+?time (?<resp_overall>.+?)".+?time (?<resp_24h>\d+?)".+?time (?<resp_7d>\d+?)".+?time (?<resp_30d>\d+?)".+?time (?<resp_1y>\d+?)".+?uptime (?<up_all>.+?%)".+?uptime (?<up_24h>.+?%)".+?uptime (?<up_7d>.+?%)".+?uptime (?<up_30d>.+?%)".+?uptime (?<up_1y>.+?%)"`)

	rePrivateUrl = regexp.MustCompile(`\| <img.+?src="(?<icon_url>.+?)".+?> (?<title>.+?) \|.+?(?<status>🟩|🟥|🟨).+?src="(?<graph_week>.+?)".+?time (?<resp_overall>.+?)".+?time (?<resp_24h>\d+?)".+?time (?<resp_7d>\d+?)".+?time (?<resp_30d>\d+?)".+?time (?<resp_1y>\d+?)".+?uptime (?<up_all>.+?%)".+?uptime (?<up_24h>.+?%)".+?uptime (?<up_7d>.+?%)".+?uptime (?<up_30d>.+?%)".+?uptime (?<up_1y>.+?%)"`)
)

func ExtractDetails(str string) []UptimeDetails {
	lines := strings.Split(str, "\n")
	details := []UptimeDetails{}

	for _, line := range lines {
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
			continue
		}

		details = append(details, UptimeDetails{
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
		})
	}

	return details
}
