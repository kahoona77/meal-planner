package views

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/a-h/templ"
)

var weekdayNames = [...]string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"}

func weekdayName(weekday time.Weekday) string {
	return weekdayNames[weekday]
}

const dateFormat = "02.01.2006"

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
