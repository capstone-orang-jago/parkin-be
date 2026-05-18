package helpers

import (
	"time"
)

const StandardLayout = "2006-01-02 15:04:05"

func GetLocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return time.Now().In(loc)
}

func FormatTime(t time.Time) string {
	return t.Format(StandardLayout)
}
