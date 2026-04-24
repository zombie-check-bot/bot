-- +goose Up
-- +goose StatementBegin
CREATE TABLE `notifications` (
    `id` SERIAL PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `type` VARCHAR(32) NOT NULL,
    `channel` VARCHAR(32) NOT NULL,
    `recipient` VARCHAR(255) NOT NULL,
    `sent_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX `idx_notifications_user_type_sent_at` ON `notifications` (`user_id`, `type`, `sent_at`);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `notifications`;
-- +goose StatementEnd