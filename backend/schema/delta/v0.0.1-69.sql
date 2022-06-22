ALTER TABLE mutualaid.`user` ADD mp_openid varchar(100) DEFAULT "" NOT NULL COMMENT '微信公众号openid';
CREATE INDEX user_mp_openid_IDX USING BTREE ON mutualaid.`user` (mp_openid);
ALTER TABLE mutualaid.`user` DROP INDEX openid;
CREATE INDEX openid USING BTREE ON mutualaid.`user` (openid);
CREATE UNIQUE INDEX uniwx USING BTREE ON mutualaid.`user` (unionid,mp_openid,openid);
