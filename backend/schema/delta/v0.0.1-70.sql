-- 用户表使用原有的status处理封禁状态
-- AId考虑有还原与恢复功能，增加字段来区分
-- Aid 默认状态为10，待审核默认不展示
ALTER TABLE
    `mutualaid`.`aid`
ADD
    COLUMN `examine_status` SMALLINT NOT NULL DEFAULT '10'
AFTER
    `status`;

-- AidMessage默认状态为20，审核通过
ALTER TABLE
    `mutualaid`.`aid_messages`
ADD
    COLUMN `examine_status` SMALLINT NOT NULL DEFAULT '20'
AFTER
    `status`;

-- 创建审核用户表
CREATE TABLE `mutualaid`.`examine_user` (
    `id` INT COMMENT '用户ID,自增',
    `name_cn` VARCHAR(45) NULL COMMENT '中文名',
    `user_name` VARCHAR(45) NULL COMMENT '用户名',
    `password` VARCHAR(45) NULL COMMENT '密码加盐',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) 
) ENGINE = InnoDB AUTO_INCREMENT = 10000032323 DEFAULT CHARSET = utf8mb4 COMMENT = '用户表';


-- 插入默认审核员
INSERT INTO
    `mutualaid`.`examine_user` (`id`, `name_cn`, `user_name`, `password`)
VALUES
    (
        '11991',
        '王珏',
        'ashly.wang',
        '65fe8599c5958b2ded217b845c686f63'
    );