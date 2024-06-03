CREATE TABLE `users`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `name`        varchar(128) NOT NULL DEFAULT '' COMMENT 'User name',
    `college_name`        varchar(128) NOT NULL DEFAULT '' COMMENT 'College Name',
    `college_address`        varchar(128) NOT NULL DEFAULT '' COMMENT 'College Address',
    `emails`   Text NULL COMMENT 'User emails',
    `created_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User information create time',
    `updated_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User information update time',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'User information delete time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User information table'