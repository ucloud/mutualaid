ALTER TABLE `mutualaid`.`aid`
    MODIFY COLUMN `group` smallint(3) NOT NULL COMMENT '求助人群: 1-重症患者,2-儿童婴儿,3-孕妇,40-老人,50-残障,60-外来务工人员,70-滞留人员,80-新冠阳性,90-医护工作者,100-街道社区,110-外籍人士,120 其他' AFTER `type`;
