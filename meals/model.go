package meals

import "meal-planner/core"

type Meal struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func InitDb(ctx *core.Ctx) {
	ctx.Db().MustExec(mealsSchema)
}

var mealsSchema = `
CREATE TABLE IF NOT EXISTS meals (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	name        VARCHAR(255)  DEFAULT '',
	description TEXT          DEFAULT ''
);`
