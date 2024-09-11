package utils

import "time"

func IsTimeIsLessThanGivenTime(timeIn24Hr int) bool {
	now := time.Now()
	hour := time.Date(now.Year(), now.Month(), now.Day(), timeIn24Hr, 0, 0, 0, now.Location())
	return now.Before(hour)
}

func GetDateDetails() (dayOfWeek int, dayOfMonth int, month int, year int) {
	now := time.Now()
	month = int(now.Month())
	dayOfMonth = now.Day()
	year = now.Year()
	// Convert Weekday to an integer (0 for Sunday, 1 for Monday, etc.)
	dayOfWeek = int(now.Weekday())
	return
}
