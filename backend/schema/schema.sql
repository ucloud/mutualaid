use mutualaid;

CREATE TABLE `user` (
                        `id` bigint(19) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
                        `phone` varchar(60) NOT NULL COMMENT '用户手机号',
                        `name` varchar(60) NOT NULL COMMENT '姓名',
                        `icon` varchar(512) NOT NULL DEFAULT '' COMMENT '头像',
                        `openid` varchar(100) NOT NULL DEFAULT '' COMMENT '微信openid',
                        `mp_openid` varchar(100) NOT NULL DEFAULT '' COMMENT '微信公众号openid',
                        `unionid` varchar(100) NOT NULL DEFAULT '' COMMENT '微信unionid',
                        `addr` varchar(500) NOT NULL COMMENT '联系地址',
                        `community` varchar(100) NOT NULL DEFAULT '' COMMENT '社区名',
                        `status` smallint(3) NOT NULL DEFAULT 1 COMMENT '用户状态: 1-新建，2-正常，3-锁定',
                        `create_time` bigint(19) NOT NULL COMMENT '创建时间，用unix时间戳表示',
                        `update_time` bigint(19) NOT NULL COMMENT '更新时间，用unix时间戳表示',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `user_unionid_IDX` (`unionid`,`mp_openid`,`openid`) USING BTREE,
                        KEY `phone` (`phone`),
                        KEY `create_time` (`create_time`),
                        KEY `user_mp_openid_IDX` (`mp_openid`) USING BTREE,
                        KEY `openid` (`openid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12325526773432323 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE aid (
                     id                    bigint(19)    NOT NULL DEFAULT 0  COMMENT '求助ID，非自增',
                     user_id               bigint(19)    NOT NULL            COMMENT '用户ID',
                     type                  smallint(3)   NOT NULL            COMMENT '求助类型: 10-食品生活物资，15-就医，20-求药，25-防疫物资，30-隔离求助，35-心理援助，40-其他',
                     `group`               smallint(3)   NOT NULL            COMMENT '求助人群: 10-重症患者，15-儿童婴儿，20-孕妇，25-老人，30-残障，35-外来务工人员，40-滞留人员，45-新冠阳性，50-医护工作者，55-街道社区，60-外籍人士',
                     emergency_level       smallint(3)   NOT NULL            COMMENT '紧急程度: 1-威胁生命，2-威胁健康，3-处境困难，4-暂无危险',
                     status                smallint(3)   NOT NULL DEFAULT 10 COMMENT '求助状态: 10-已创建，15-已取消，20-已完成',
                     examine_status        smallint(3)   NOT NULL DEFAULT 10 COMMENT '审核状态: 10-待审核，15-审核不通过，20-审核通过',
                     finish_user_id        bigint(19)    NOT NULL DEFAULT 0  COMMENT '完成用户ID',
                     finish_time           int(11)       NOT NULL DEFAULT 0  COMMENT '完成时间，用unix时间戳表示',
                     message_count         int(11)       NOT NULL DEFAULT 0  COMMENT '留言数量，增加留言时更新',
                     content               varchar(400)  NOT NULL            COMMENT '描述',
                     longitude             DECIMAL(11,8) NOT NULL DEFAULT 0  COMMENT '坐标：经度',
                     latitude              DECIMAL(11,8) NOT NULL DEFAULT 0  COMMENT '坐标：纬度',
                     phone                 varchar(15)   NOT NULL            COMMENT '联系电话',
                     district              varchar(10)   NOT NULL DEFAULT '' COMMENT '区县，使用国家民政局标准',
                     address               varchar(100)  NOT NULL DEFAULT '' COMMENT '地址',
                     create_time           int(11)       NOT NULL            COMMENT '创建时间，用unix时间戳表示',
                     update_time           int(11)       NOT NULL            COMMENT '更新时间，用unix时间戳表示',
                     version               int(11)       NOT NULL DEFAULT 1  COMMENT '版本，用于乐观锁控制',
                     PRIMARY KEY (id),
                     KEY idx_user (user_id),
                     KEY idx_finish_user (finish_user_id),
                     KEY idx_create_time (create_time),
                     KEY idx_status (status)
) COMMENT '援助';

CREATE TABLE aid_messages (
                              id                    bigint(19)    NOT NULL            COMMENT '帮助ID',
                              aid_id                bigint(19)    NOT NULL            COMMENT '求助ID',
                              status                smallint(3)   NOT NULL DEFAULT 10 COMMENT '状态: 10-已创建，15-已取消，20-已完成',
                              examine_status        smallint(3)   NOT NULL DEFAULT 15 COMMENT '审核状态: 10-待审核，15-审核不通过，20-审核通过',
                              user_id               bigint(19)    NOT NULL            COMMENT '用户ID',
                              user_phone            varchar(15)   NOT NULL            COMMENT '联系电话',
                              content               varchar(500)  NOT NULL            COMMENT '帮助说明',
                              create_time           int(11)       NOT NULL            COMMENT '创建时间，用unix时间戳表示',
                              update_time           int(11)       NOT NULL            COMMENT '更新时间，用unix时间戳表示',
                              version               int(11)       NOT NULL DEFAULT 1  COMMENT '版本，用于乐观锁控制',
                              PRIMARY KEY (id),
                              KEY idx_aid (aid_id),
                              KEY idx_user (user_id)
) COMMENT '援助留言';


CREATE TABLE examine_user (
                              `id` INT  NOT NULL     COMMENT '用户ID,自增',
                              `name_cn` VARCHAR(45) NULL  COMMENT '中文名',
                              `user_name` VARCHAR(45) NULL COMMENT '用户名',
                              `password` VARCHAR(45) NULL COMMENT '密码加盐',
                              PRIMARY KEY (`id`),
                              UNIQUE INDEX `id_UNIQUE` (`id` ASC)  )
    ENGINE=InnoDB AUTO_INCREMENT=100032323 DEFAULT CHARSET=utf8mb4 COMMENT='审核员表';
