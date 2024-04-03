-- +goose Up
CREATE TABLE weekday_tags
(
    weekday INTEGER PRIMARY KEY,
    tag_id  INTEGER REFERENCES tags
);


-- +goose Down
DROP TABLE weekday_tags;