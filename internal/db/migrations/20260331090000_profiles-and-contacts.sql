-- +goose Up
-- +goose StatementBegin
CREATE TABLE `profiles` (
    `user_id` CHAR(36) PRIMARY KEY,
    `username` VARCHAR(255) NULL,
    `display_name` VARCHAR(255) NULL,
    `locale` VARCHAR(16) NOT NULL DEFAULT 'en',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `fk_profiles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `contacts` (
    `id` CHAR(36) PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `name` VARCHAR(255) NULL,
    `type` VARCHAR(32) NOT NULL,
    `value` VARCHAR(255) NOT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `fk_trusted_contacts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `uk_trusted_contacts_owner_type_value` UNIQUE (`user_id`, `type`, `value`)
);
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `contacts`;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS `profiles`;
-- +goose StatementEnd