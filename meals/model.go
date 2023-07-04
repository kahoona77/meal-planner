package meals

import (
	"database/sql"
)

type Meal struct {
	Id          int64         `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	ImageFileId sql.NullInt64 `db:"image_file_id"`
}

type Tag struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	Color string `db:"color"`
}
