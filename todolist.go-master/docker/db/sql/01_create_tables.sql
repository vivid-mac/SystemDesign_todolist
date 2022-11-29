-- Table for tasks
DROP TABLE IF EXISTS `tasks`;

CREATE TABLE `tasks` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `title` varchar(50) NOT NULL,
    `is_done` boolean NOT NULL DEFAULT b'0',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `explanation` varchar(256) NOT NULL DEFAULT "",
    `priority` varchar(256) NOT NULL DEFAULT "LOW",
    `deadline` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `tag` boolean NOT NULL DEFAULT b'0',
    `category` varchar(256) NOT NULL DEFAULT "",
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `users`;
 
CREATE TABLE `users` (
    `id`         bigint(20) NOT NULL AUTO_INCREMENT,
    `name`       varchar(50) NOT NULL UNIQUE,
    `password`   binary(32) NOT NULL,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `ownership`;
 
CREATE TABLE `ownership` (
    `user_id` bigint(20) NOT NULL,
    `task_id` bigint(20) NOT NULL,
    PRIMARY KEY (`user_id`, `task_id`)
) DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `groups`;

CREATE TABLE `groups` (
    `id`            bigint(20) NOT NULL AUTO_INCREMENT,
    `name`          varchar(50) NOT NULL UNIQUE,
    `created_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `belong`;

CREATE TABLE `belong` (
    `group_id`      bigint(20) NOT NULL,
    `user_id`       bigint(20) NOT NULL,
    PRIMARY KEY (`group_id`, `user_id`)
) DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `share`;

CREATE TABLE `share` (
    `group_id`      bigint(20) NOT NULL,
    `task_id`       bigint(20) NOT NULL,
    PRIMARY KEY (`group_id`, `task_id`)
) DEFAULT CHARSET=utf8mb4;
