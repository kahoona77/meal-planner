-- +goose Up
-- Deleting a meal used to keep its tag assignments, and an old version stored some with meal_id 0.
-- The wizard picks meals by their tags, so these rows made it fail.
DELETE FROM meal_tags WHERE meal_id NOT IN (SELECT id FROM meals);

-- +goose Down
-- the deleted rows referenced meals that do not exist, nothing to restore
