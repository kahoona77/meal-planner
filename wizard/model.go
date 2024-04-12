package wizard

import "time"

type Week struct {
	Start  time.Time
	End    time.Time
	Offset int
	Days   []*Day
}

type Day struct {
	Weekday int
	Date    time.Time
	Tags    []*WeekdayTag
}

type WeekdayTag struct {
	Weekday  int
	TagId    int64  `db:"tag_id"`
	TagName  string `db:"-"`
	TagColor string `db:"-"`
}
