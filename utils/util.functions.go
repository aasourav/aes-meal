package utils

import "time"

func IsTimeIsLessThanGivenTime(timeIn24Hr int) bool {
	now := time.Now()
	hour := time.Date(now.Year(), now.Month(), now.Day(), timeIn24Hr, 0, 0, 0, now.Location())
	return now.Before(hour)
}
