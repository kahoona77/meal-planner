-- +goose Up
CREATE TABLE weekday_tags
(
    weekday INTEGER,
    tag_id  INTEGER REFERENCES tags,
    PRIMARY KEY (weekday, tag_id)
);


-- +goose Down
DROP TABLE weekday_tags;