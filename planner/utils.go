package planner

import "time"

const firstWeekDay = time.Saturday
const lastWeekDay = time.Friday

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetStartWeek(date time.Time) time.Time {
	date = Bod(date)
	for date.Weekday() != firstWeekDay {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func GetEndWeek(date time.Time) time.Time {
	date = Bod(date)
	for date.Weekday() != lastWeekDay {
		date = date.AddDate(0, 0, 1)
	}
	return date
}
