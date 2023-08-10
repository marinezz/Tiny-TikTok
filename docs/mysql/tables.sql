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



 /*
  消息表
  */


/*
 评论表
 */




/*
 点赞表
 */



/*
 关注表
 */