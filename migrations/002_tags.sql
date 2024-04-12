-- +goose Up
CREATE TABLE tags
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL DEFAULT '',
    color VARCHAR(50) NOT NULL DEFAULT ''
);

CREATE TABLE meal_tags
(
    meal_id INTEGER REFERENCES meals,
    tag_id INTEGER REFERENCES tags,
    PRIMARY KEY (meal_id, tag_id)
);

-- +goose Down
DROP TABLE meal_tags;
DROP TABLE tags;