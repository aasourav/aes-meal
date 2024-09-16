package utils

import (
	"time"
)

func ItTimeIsInRange(afterTime int, beforeTime int) bool {
	now := time.Now()
	// before := time.Date(now.Year(), now.Month(), now.Day(), beforeTime, 0, 0, 0, now.Location())
	// after := time.Date(now.Year(), now.Month(), now.Day(), beforeTime, 0, 0, 0, now.Location())
	curHour := now.Local().Hour()
	// fmt.Println("HI temp:", now.Before(before), now.After(after))
	// return now.Before(before) && now.After(after)
	return curHour < beforeTime && curHour > afterTime

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
