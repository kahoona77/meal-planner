package meals

import (
	"database/sql"
	"meal-planner/core"
)

type Meal struct {
	Id          int64         `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	ImageFileId sql.NullInt64 `db:"image_file_id"`
}

func InitDb(ctx *core.Ctx) {
	ctx.Db().MustExec(mealsSchema)
}

var mealsSchema = `
CREATE TABLE IF NOT EXISTS meals (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	name          VARCHAR(255)  DEFAULT '',
	description   TEXT          DEFAULT '',
	image_file_id INTEGER REFERENCES files
);`
