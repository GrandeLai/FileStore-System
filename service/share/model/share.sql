CREATE TABLE `share` (
   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
   `user_id` bigint unsigned NOT NULL DEFAULT '0',
   `store_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '公共池中的唯一标识',
   `expired_time` int(11) NOT NULL DEFAULT '0' COMMENT '失效时间，单位秒, 【0-永不失效】',
   `share_url` varchar(32) NOT NULL DEFAULT '' COMMENT '分享链接后缀',
   `extraction_code` varchar(10) NOT NULL DEFAULT '' COMMENT '提取码',
   `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_user_id` (`user_id`),
   KEY `idx_store_id` (`store_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;