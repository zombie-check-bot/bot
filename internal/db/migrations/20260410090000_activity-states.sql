-- +goose Up
-- +goose StatementBegin
CREATE TABLE `activity_states` (
    `user_id` CHAR(36) PRIMARY KEY,
    `last_alive` DATETIME NOT NULL,
    `check_interval_days` INT NOT NULL DEFAULT 7,
    `timeout_days` INT NOT NULL DEFAULT 14,
    `reminders_json` JSON NOT NULL,
    `is_notified` BOOLEAN NOT NULL DEFAULT FALSE,
    `notified_at` DATETIME NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `fk_activity_states_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `notification_logs` (
    `id` CHAR(36) PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `type` VARCHAR(32) NOT NULL,
    `recipient` VARCHAR(255) NOT NULL,
    `sent_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `fk_notification_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `notification_logs`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS `activity_states`;
-- +goose StatementEnd
