CREATE TABLE `user` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户密码',
    `email` varchar(64) DEFAULT '' COMMENT '邮箱',
    `phone` varchar(128) DEFAULT '' COMMENT '手机号',
    `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
    `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机号是否已验证',
    `now_volume`   int(11) NOT NULL DEFAULT '0' COMMENT '当前存储容量',
    `total_volume` int(11) NOT NULL DEFAULT '1000' COMMENT '最大存储容量',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username_unique` (`user_name`),
    UNIQUE KEY `idx_email_unique` (`email`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;