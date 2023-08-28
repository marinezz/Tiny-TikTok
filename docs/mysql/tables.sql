/*
 用户表
 */
CREATE TABLE `user` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_name` varchar(256) DEFAULT NULL,
    `pass_word` varchar(256) NOT NULL,
    `avatar` varchar(256) DEFAULT NULL,
    `background_image` varchar(256) DEFAULT NULL,
    `signature` varchar(256) DEFAULT '该用户还没有简介',
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=812575311663105 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


/*
 视频表
 */
CREATE TABLE `video` (
     `id` bigint NOT NULL AUTO_INCREMENT,
     `auth_id` bigint DEFAULT NULL,
     `title` varchar(256) DEFAULT NULL,
     `cover_url` varchar(256) DEFAULT NULL,
     `play_url` varchar(256) DEFAULT NULL,
     `favorite_count` bigint DEFAULT '0',
     `comment_count` bigint DEFAULT '0',
     `creat_at` datetime DEFAULT NULL,
     PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2276964627783681 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


 /*
  消息表
  */
CREATE TABLE `message` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `user_id` bigint DEFAULT NULL,
   `to_user_id` bigint DEFAULT NULL,
   `message` varchar(256) DEFAULT NULL,
   `created_at` varchar(256) DEFAULT NULL,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


/*
 评论表
 */
CREATE TABLE `comment` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `user_id` bigint DEFAULT NULL,
   `video_id` bigint DEFAULT NULL,
   `creat_at` datetime DEFAULT NULL,
   `comment_status` tinyint(1) DEFAULT NULL,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


/*
 点赞表
 */
CREATE TABLE `favorite` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint DEFAULT NULL,
    `video_id` bigint DEFAULT NULL,
    `is_favorite` tinyint(1) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2277278458191873 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


/*
 关注表
 */
CREATE TABLE `follow` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `to_user_id` bigint DEFAULT NULL,
  `is_follow` int DEFAULT (2),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7638289369411586 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
