package utils

import "time"

func ItTimeIsInRange(beforeTime int, afterTime int) bool {
	now := time.Now()
	before := time.Date(now.Year(), now.Month(), now.Day(), beforeTime, 0, 0, 0, now.Location())
	after := time.Date(now.Year(), now.Month(), now.Day(), beforeTime, 0, 0, 0, now.Location())
	return now.Before(before) && now.After(after)
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
