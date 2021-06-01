# 结构

```
adapter 对外的接口 提供对外的api
domain 负责领域相关内容
database 具体的数据库实现文件
```

# 登录模块开发


# 数据库

```database
CREATE TABLE `UserInfo` (
	`useremail` VARCHAR(50) NOT NULL COMMENT '用户邮箱' COLLATE 'utf8_bin',
	`useraccount` INT(11) NOT NULL AUTO_INCREMENT COMMENT '用户账号',
	`username` VARCHAR(50) NULL DEFAULT 'xx1111' COMMENT '用户名' COLLATE 'utf8_bin',
	`signature` VARCHAR(50) NULL DEFAULT NULL COMMENT '用户签名' COLLATE 'utf8_bin',
	`avatar` VARCHAR(50) NULL DEFAULT NULL COMMENT '用户头像' COLLATE 'utf8_bin',
	`usersex` INT(11) NULL DEFAULT '1' COMMENT '用户性别1为男 0 为女',
	`userpassword` VARCHAR(50) NOT NULL DEFAULT '123456' COMMENT '用户密码' COLLATE 'utf8_bin',
	`userage` INT(11) NULL DEFAULT '18' COMMENT '用户年龄',
	PRIMARY KEY (`useraccount`, `useremail`) USING BTREE
)
COMMENT='用户信息表'
COLLATE='utf8_bin'
ENGINE=InnoDB
AUTO_INCREMENT=2
;


CREATE TABLE `UserFriend` (
	`launcher` INT(11) NOT NULL DEFAULT '0',
	`accepter` INT(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`launcher`, `accepter`) USING BTREE
)
COMMENT='用户-好友表'
COLLATE='latin1_swedish_ci'
ENGINE=InnoDB

```

# 要加入返回的http错误处理
