package views

import (
	"database/sql"
	"fmt"
	"math"
	"meal-planner/planner"
	"regexp"
	"strings"
	"time"

	"github.com/a-h/templ"
)

var weekdayNames = [...]string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"}

func weekdayName(weekday time.Weekday) string {
	return weekdayNames[weekday]
}

// weekdayShort returns the abbreviation of the weekday, e.g. "Mo".
func weekdayShort(weekday time.Weekday) string {
	return weekdayNames[weekday][:2]
}

const dateFormat = "02.01.2006"

// weekOffsetURL returns the URL of the week with the given offset from the current one.
// The current week is "./" (the start page) and not "", which is no URL for layout.
func weekOffsetURL(offset int) string {
	if offset == 0 {
		return "./"
	}
	return fmt.Sprintf("offset/%d", offset)
}

// weekURL returns the URL of the week containing date.
func weekURL(date time.Time) string {
	weeks := planner.GetStartWeek(date).Sub(planner.GetStartWeek(time.Now())).Hours() / (24 * 7)
	return weekOffsetURL(int(math.Round(weeks)))
}

// selectBackURL is where the back arrow of the meal selection leads: the meal of the day,
// or the week if nothing is planned (empty days link to the selection directly).
func selectBackURL(mod *planner.MealOfTheDay) string {
	if !mod.MealId.Valid {
		return weekURL(mod.Date)
	}
	return "planner/" + mod.Id
}

// isToday reports whether t is on the current day.
func isToday(t time.Time) bool {
	return t.Format(time.DateOnly) == time.Now().Format(time.DateOnly)
}

func isoWeek(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

// navSection returns the href of the navigation entry a page belongs to ("", "meals"
// or "tags"), path is relative to the base path as returned by web.Path.
func navSection(path string) string {
	section, _, _ := strings.Cut(path, "/")
	switch section {
	case "meals", "tags":
		return section
	default:
		return ""
	}
}

// fileURL returns the (relative) URL of an uploaded file, or "" if there is none.
func fileURL(id sql.NullInt64) string {
	if !id.Valid {
		return ""
	}
	return fmt.Sprintf("files/%d", id.Int64)
}

var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{3,8}$`)

// colorStyle returns "property: color" for an inline style. Tag colors come from the
// database, so only hex colors are let through.
func colorStyle(property string, color string) templ.SafeCSS {
	if !hexColor.MatchString(color) {
		return ""
	}
	return templ.SafeCSS(property + ": " + color)
}
