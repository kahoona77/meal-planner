package wizard

import (
	"github.com/jmoiron/sqlx"
	"meal-planner/core"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(c core.Context) *Repository {
	return &Repository{db: c.Db()}
}

func (r *Repository) GetWeekdayTags() ([]*WeekdayTag, error) {
	var tags []*WeekdayTag
	err := r.db.Select(&tags, "SELECT * FROM weekday_tags")
	return tags, err
}

const weekdayTagsInsert = `INSERT INTO weekday_tags (weekday, tag_id) VALUES (:weekday, :tag_id)`

func (r *Repository) SetWeekdayTags(weekday int, tags []*WeekdayTag) error {
	// delete old tags
	_, err := r.db.Exec("DELETE FROM weekday_tags WHERE weekday=$1", weekday)
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		// insert new tags
		_, err = r.db.NamedExec(weekdayTagsInsert, tags)
		return err
	}

	return nil
}
