package wizard

import (
	"meal-planner/planner"
	"time"
)

type Wizard interface {
	Generate(wizardWeek Week) (*planner.Week, error)
}

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
