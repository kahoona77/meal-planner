package web

var weekdays = map[int]string{
	0: "Sonntag",
	1: "Montag",
	2: "Dienstag",
	3: "Mittwoch",
	4: "Donnerstag",
	5: "Freitag",
	6: "Samstag",
}

func formatWeekday(weekday int) string {
	return weekdays[weekday]
}
