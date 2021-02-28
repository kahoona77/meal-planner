package planner

import "time"

const firstWeekDay = time.Saturday
const lastWeekDay = time.Friday

func GetStartWeek(date time.Time) time.Time { //get monday 00:00:00
	for date.Weekday() != firstWeekDay {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func GetEndWeek(date time.Time) time.Time { //get monday 00:00:00
	for date.Weekday() != lastWeekDay {
		date = date.AddDate(0, 0, 1)
	}
	return date
}

func IsToday(date time.Time) bool { //get monday 00:00:00
	now := time.Now()
	return now.Year() == date.Year() && now.Month() == date.Month() && now.Day() == date.Day()
}
