package utils

import "time"

const TimeFormat = time.RFC3339

func NowFormatted() string {
	return time.Now().Format(TimeFormat)
}
