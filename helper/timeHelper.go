package helper

import "time"

const layoutStr = "2006-01-02 15:04:05"

func ParseStringToTime(timeStr string) (timeRes time.Time) {
	timeRes, err := time.Parse(layoutStr, timeStr)
	if err != nil {
		timeRes, _ = time.Parse(layoutStr, "1000-01-01 00:00:00")
	}
	return
}
