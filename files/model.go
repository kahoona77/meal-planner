package files

import "meal-planner/core"

type File struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	ContentType string `db:"content_type"`
	Data        []byte `db:"data"`
}

func InitDb(ctx *core.Ctx) {
	ctx.Db().MustExec(filesSchema)
}

var filesSchema = `
CREATE TABLE IF NOT EXISTS files (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	name         VARCHAR(255) NOT NULL DEFAULT '',
	content_type VARCHAR(50) NOT NULL DEFAULT '',
	data         BLOB
);`
