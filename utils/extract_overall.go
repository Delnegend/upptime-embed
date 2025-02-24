package utils

import "regexp"

type Overall string

const (
	OverallAllSystemsOperational Overall = "all_good"
	OverallDegradedPerformance   Overall = "degraded"
	OverallCompleteOutage        Overall = "down"
	OverallPartialOutage         Overall = "partial"
	OverallUnknownn              Overall = "unknown"
)

var reOverall = regexp.MustCompile(`\<\!--live status--\>.+?\*\*(?<overall>.+?)\*`)

func ExtractOverall(str string) Overall {
	match := reOverall.FindStringSubmatch(str)
	if len(match) == 0 {
		return OverallUnknownn
	}
	switch match[1] {
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
