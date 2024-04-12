-- +goose Up
CREATE TABLE IF NOT EXISTS files
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         VARCHAR(255) NOT NULL DEFAULT '',
    content_type VARCHAR(50)  NOT NULL DEFAULT '',
    data         BLOB
);
CREATE TABLE IF NOT EXISTS meals
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          VARCHAR(255) DEFAULT '',
    description   TEXT         DEFAULT '',
    image_file_id INTEGER REFERENCES files
);
CREATE TABLE IF NOT EXISTS meals_of_day
(
    id      VARCHAR(10) PRIMARY KEY,
    date    TIMESTAMP NOT NULL,
    meal_id INTEGER REFERENCES meals
);

-- +goose Down
DROP TABLE meals_of_day;
DROP TABLE meals;
DROP TABLE files;