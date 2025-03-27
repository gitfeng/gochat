use gochat;

CREATE TABLE customers (
    `id` BIGINT UNSIGNED AUTO_INCREMENT COMMENT '客户唯一标识',
    `customer_name` VARCHAR(128) NOT NULL COMMENT '客户名称（中文支持）',
    `password` VARCHAR(255) NOT NULL COMMENT 'BCrypt加密密码',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
    PRIMARY KEY (`id`),    
    UNIQUE INDEX idx_customer_name (customer_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户信息表';

CREATE TABLE messages (
    `id` BIGINT UNSIGNED AUTO_INCREMENT COMMENT '消息唯一ID',
    `customer_id` BIGINT UNSIGNED NOT NULL COMMENT '关联客户ID',
    `message` VARCHAR(1024) NOT NULL COMMENT '消息内容（支持中文）',
    `sender` VARCHAR(32) NOT NULL COMMENT '发送者标识',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
    `message_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:普通消息,1:feedback引导消息',
    PRIMARY KEY (`id`),
    INDEX idx_customer_at (customer_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户消息记录表';

CREATE TABLE `feedback` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '反馈唯一ID',
  `customer_id` bigint unsigned NOT NULL COMMENT '关联客户ID',
  `score` tinyint unsigned NOT NULL COMMENT '用户评分 (0-10)',
  `comment` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL COMMENT '反馈内容（支持中文）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '反馈创建时间',
  `sentiment` tinyint NOT NULL DEFAULT '0' COMMENT '0:neutral,1:positive,-1:negative',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_customer_feedback` (`customer_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户反馈记录表';

insert into customers (customer_name, password) values ('admin', '$2a$10$3Jj2V5s933h86X46z1z5Y.5z1z1z1z1z1z1z1z1z1z1z1z1z1z1z1z1z1z');