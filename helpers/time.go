package helpers

import (
	"time"
)

const StandardLayout = "2006-01-02 15:04:05"

var witaZone = time.FixedZone("WITA", 8*60*60)

func GetLocalTime() time.Time {
	return time.Now().In(witaZone)
}

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(witaZone).Format(StandardLayout)
}
