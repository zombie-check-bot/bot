-- +goose Up
-- +goose StatementBegin
CREATE TABLE `activity` (
    `_id` SERIAL PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `unq_activity_log` UNIQUE (`user_id`, `created_at`),
    CONSTRAINT `fk_activity_log_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX `idx_activity_created_at_user_id` ON `activity` (`created_at`, `user_id`);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP INDEX `idx_activity_created_at_user_id` ON `activity`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS `activity`;
-- +goose StatementEnd