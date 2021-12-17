use bridgx;
--
-- Table structure for table `account`
--

DROP TABLE IF EXISTS `account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account`
(
    `id`             bigint(20) NOT NULL AUTO_INCREMENT,
    `account_name`   varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `account_key`    varchar(256) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `encrypted_account_secret` varchar(1024) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `salt` varchar(256) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `org_id`         bigint(20) DEFAULT '0',
    `provider`       varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `create_at`      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `create_by`      varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `update_by`      varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
    `deleted_at`     timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY              `ak_index` (`account_key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `b_security_group`
--

DROP TABLE IF EXISTS `b_security_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `b_security_group`
(
    `id`                  bigint(20) NOT NULL AUTO_INCREMENT,
    `vpc_id`              varchar(255) NOT NULL,
    `security_group_id`   varchar(255) NOT NULL DEFAULT '',
    `name`                varchar(255) NOT NULL DEFAULT '',
    `security_group_type` varchar(255) NOT NULL DEFAULT '',
    `is_del`              tinyint(3) NOT NULL DEFAULT '0',
    `create_at`           datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    `update_at`           datetime     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_vpc_group_name` (`vpc_id`,`security_group_id`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='安全组表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `b_security_group_rule`
--

DROP TABLE IF EXISTS `b_security_group_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `b_security_group_rule`
(
    `id`                bigint(20) NOT NULL AUTO_INCREMENT,
    `vpc_id`            varchar(255) NOT NULL,
    `security_group_id` varchar(255) NOT NULL DEFAULT '',
    `port_range`        varchar(255) NOT NULL DEFAULT '',
    `protocol`          varchar(255) NOT NULL DEFAULT '',
    `direction`         varchar(10)  NOT NULL,
    `other_group_id`    varchar(255) NOT NULL,
    `cidr_ip`           varchar(255) NOT NULL,
    `prefix_list_id`    varchar(255) NOT NULL,
    `is_del`            tinyint(3) NOT NULL DEFAULT '0',
    `create_at`         datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    `update_at`         datetime     NOT NULL,
    PRIMARY KEY (`id`),
    KEY                 `idx_vpc_group_protocol` (`vpc_id`,`security_group_id`,`protocol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全组规则表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `b_switch`
--

DROP TABLE IF EXISTS `b_switch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `b_switch`
(
    `id`                         bigint(20) NOT NULL AUTO_INCREMENT,
    `vpc_id`                     varchar(255) NOT NULL,
    `switch_id`                  varchar(255) NOT NULL DEFAULT '',
    `zone_id`                    varchar(255) NOT NULL DEFAULT '',
    `name`                       varchar(255) NOT NULL DEFAULT '',
    `cidr_block`                 varchar(255) NOT NULL DEFAULT '',
    `gateway_ip`                 varchar(128) NOT NULL DEFAULT '',
    `v_status`                   varchar(64)  NOT NULL,
    `is_default`                 tinyint(3) NOT NULL DEFAULT '0' COMMENT '0 非默认 1 默认',
    `available_ip_address_count` int(100) NOT NULL,
    `is_del`                     tinyint(3) NOT NULL DEFAULT '0',
    `create_at`                  datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    `update_at`                  datetime     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_vpc_switch` (`vpc_id`,`switch_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='交换机表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `b_vpc`
--

DROP TABLE IF EXISTS `b_vpc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `b_vpc`
(
    `id`         bigint(20) NOT NULL AUTO_INCREMENT,
    `ak`         varchar(255) NOT NULL,
    `vpc_id`     varchar(255) NOT NULL,
    `region_id`  varchar(255) NOT NULL DEFAULT '',
    `name`       varchar(255) NOT NULL,
    `cidr_block` varchar(255) NOT NULL DEFAULT '',
    `switch_ids` varchar(255) NOT NULL DEFAULT '[]',
    `provider`   varchar(255) NOT NULL,
    `v_status`   varchar(20)  NOT NULL DEFAULT 'Pending',
    `is_del`     tinyint(3) NOT NULL DEFAULT '0',
    `create_at`  datetime     NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    `update_at`  datetime     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_ak_region_vpc_id` (`ak`,`region_id`,`vpc_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='vpc 表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cluster`
--

DROP TABLE IF EXISTS `cluster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cluster`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `cluster_name`    varchar(64) COLLATE utf8mb4_bin  NOT NULL,
    `cluster_desc`    varchar(128) COLLATE utf8mb4_bin NOT NULL,
    `expect_count`    int(7) NOT NULL DEFAULT '0',
    `status`          varchar(32) COLLATE utf8mb4_bin           DEFAULT NULL,
    `region_id`       varchar(64) COLLATE utf8mb4_bin           DEFAULT NULL,
    `zone_id`         varchar(64) COLLATE utf8mb4_bin           DEFAULT NULL,
    `instance_type`   varchar(32) COLLATE utf8mb4_bin           DEFAULT NULL,
    `image`           varchar(512) COLLATE utf8mb4_bin          DEFAULT NULL,
    `provider`        varchar(64) COLLATE utf8mb4_bin           DEFAULT NULL,
    `password`        varchar(128) COLLATE utf8mb4_bin          DEFAULT NULL,
    `account_key`     varchar(128) COLLATE utf8mb4_bin          DEFAULT NULL,
    `image_config`    varchar(256) COLLATE utf8mb4_bin          DEFAULT '',
    `network_config`  varchar(4096) COLLATE utf8mb4_bin         DEFAULT NULL,
    `storage_config`  varchar(4096) COLLATE utf8mb4_bin         DEFAULT NULL,
    `charge_config`   varchar(4096) COLLATE utf8mb4_bin         DEFAULT NULL,
    `create_at`       timestamp                        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`       timestamp                        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `create_by`       varchar(32) COLLATE utf8mb4_bin           DEFAULT '',
    `update_by`       varchar(32) COLLATE utf8mb4_bin           DEFAULT '',
    `deleted_at`      timestamp NULL DEFAULT NULL,
    `delete_uniq_key` bigint(20) DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `cluster_cluster_name_uindex` (`cluster_name`, `delete_uniq_key`),
    KEY `cluster_account_key_index` (`account_key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cluster_tag`
--

DROP TABLE IF EXISTS `cluster_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cluster_tag`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `cluster_name` varchar(32) COLLATE utf8mb4_bin DEFAULT NULL,
    `tag_key`      varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
    `tag_value`    varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
    `create_at`    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `update_at`    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `cluster_tags_cluster_name_tag_key_uindex` (`cluster_name`,`tag_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `instance`
--

DROP TABLE IF EXISTS `instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `instance`
(
    `id`             bigint(20) NOT NULL AUTO_INCREMENT,
    `cluster_name`   varchar(64)          DEFAULT NULL,
    `task_id`        bigint(20) NOT NULL DEFAULT '-1',
    `shrink_task_id` bigint(20) NOT NULL DEFAULT '-1',
    `instance_id`    varchar(255)         DEFAULT NULL,
    `status`         varchar(32) NOT NULL DEFAULT 'UNDEFINED',
    `ip_inner`       varchar(255)         DEFAULT NULL,
    `ip_outer`       varchar(255)         DEFAULT NULL,
    `create_at`      timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `delete_at`      timestamp NULL DEFAULT NULL,
    `running_at`     timestamp NULL DEFAULT NULL,
    `charge_type`    varchar(32) collate utf8mb4_bin NOT NULL DEFAULT 'PostPaid',
    `expire_at`      timestamp NULL DEFAULT NULL comment 'PrePaid instance expire time',
    PRIMARY KEY (`id`),
    KEY              `idx_ip_inner` (`ip_inner`),
    KEY              `instance_cluster_name_status_index` (`cluster_name`,`status`),
    KEY              `instance_shrink_task_id_index` (`shrink_task_id`),
    KEY              `instance_task_id_index` (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `instance_type`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `instance_type`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT,
    `provider`  varchar(20)  NOT NULL,
    `region_id` varchar(200) NOT NULL,
    `zone_id`   varchar(200) NOT NULL,
    `type_name` varchar(50)  NOT NULL,
    `family`    varchar(50)           DEFAULT NULL,
    `core`      int(11) NOT NULL COMMENT '核心数,单位 核',
    `memory`    int(11) NOT NULL COMMENT '内存,单位 G',
    `i_status`  tinyint(4) NOT NULL COMMENT '0 待激活 1 已激活 2 已过期',
    `create_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `update_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `org`
--

DROP TABLE IF EXISTS `org`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `org`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT,
    `org_name`  varchar(128) COLLATE utf8mb4_bin NOT NULL,
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY         `org_name_index` (`org_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `task`
--

DROP TABLE IF EXISTS `task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `task`
(
    `id`             bigint(20) NOT NULL,
    `task_name`      varchar(64) COLLATE utf8mb4_bin NOT NULL,
    `status`         varchar(32) COLLATE utf8mb4_bin DEFAULT NULL,
    `task_action`    varchar(32) COLLATE utf8mb4_bin DEFAULT NULL,
    `task_filter`    varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
    `task_info`      text COLLATE utf8mb4_bin,
    `task_result`    text COLLATE utf8mb4_bin,
    `err_msg`        text COLLATE utf8mb4_bin,
    `support_cancel` tinyint(1) DEFAULT NULL,
    `finish_time`    timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `create_at`      timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `update_at`      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY              `task_task_filter_index` (`task_filter`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT,
    `username`    varchar(128) COLLATE utf8mb4_bin NOT NULL,
    `password`    varchar(128) COLLATE utf8mb4_bin NOT NULL,
    `user_type`   tinyint(1) NOT NULL DEFAULT '0',
    `user_status` varchar(16) COLLATE utf8mb4_bin  NOT NULL DEFAULT 'enable',
    `org_id`      bigint(20) NOT NULL DEFAULT '0',
    `create_at`   timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `create_by`   varchar(128) COLLATE utf8mb4_bin NOT NULL,
    `update_at`   timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

drop table if exists `order_202101`;
create table `order_202101`
(
    `id`               bigint(20) not null auto_increment,
    `account_name`     varchar(64) collate utf8mb4_bin   not null default '',
    `order_id`         varchar(32) collate utf8mb4_bin   not null,
    `order_time`       timestamp                         not null,
    `product`          varchar(32) collate utf8mb4_bin   not null comment 'esc',
    `quantity`         int                               not null default '1',
    `usage_start_time` timestamp                         not null DEFAULT current_timestamp,
    `usage_end_time`   timestamp                         not null DEFAULT current_timestamp,
    `provider`         varchar(64) collate utf8mb4_bin   not null default 'AlibabaCloud',
    `region_id`        varchar(64) collate utf8mb4_bin   not null default '',
    `charge_type`      varchar(32) collate utf8mb4_bin   not null default 'PostPaid',
    `pay_status`       tinyint                           not null default '1' comment '1:已支付，2：未支付，3：取消',
    `currency`         varchar(8) collate utf8mb4_bin    not null default 'CNY',
    `cost`             float                             not null default '0',
    `extend`           varchar(4096) collate utf8mb4_bin not null default '{"main_order_id":"","order_type":"new"}' comment '拓展字段，json格式',
    `create_at`        timestamp                         not null DEFAULT current_timestamp,
    `update_at`        timestamp                         not null DEFAULT current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key uniq_account_order_id(`account_name`,`order_id`),
    index              idx_account_order_time(`account_name`,`order_time`),
    index              idx_cost(`cost`),
    index              idx_update_at(`update_at`)
)engine=innodb default charset=utf8mb4 collate=utf8mb4_bin;


-- gf_cloud tables
drop table if exists kubernetes_infos;
create table kubernetes_infos (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL ,
  `region` varchar(255) NOT NULL ,
  `cloud_type` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `config` TEXT,
  `install_step` varchar(255) ,
  `message` varchar (255) ,
  `bridgx_cluster_name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `created_user` varchar(255),
  `created_time` int(11),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_cluster_name` (`name`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin;

drop table if exists instance_groups;
CREATE TABLE `instance_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `kubernetes_id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `image` varchar(255) NOT NULL,
  `cpu` varchar(255) NOT NULL,
  `memory` varchar(255) NOT NULL,
  `disk` varchar(255) NOT NULL,
  `instance_count` int(11) NOT NULL DEFAULT '0',
  `created_user` varchar(255) NOT NULL,
  `created_user_id` int(11) NOT NULL DEFAULT '0',
  `ssh_pwd` VARCHAR(500) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_kubernetes_id_name` (`kubernetes_id`,`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

drop table if exists instance_forms;
CREATE TABLE `instance_forms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `execute_status` varchar(100) NOT NULL,
  `cluster_name` varchar(255) DEFAULT NULL,
  `instance_group` varchar(255) NOT NULL,
  `cpu` varchar(255) NOT NULL,
  `memory` varchar(255) NOT NULL,
  `disk` varchar(255) NOT NULL,
  `updated_instance_count` int(11) NOT NULL DEFAULT '0',
  `host_time` int(11) NOT NULL,
  `opt_type` varchar(50) DEFAULT NULL,
  `created_user_id` int(11) NOT NULL DEFAULT '0',
  `created_user_name` varchar(255) NOT NULL,
  `created_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

drop table if exists kubernetes_install_steps;
CREATE TABLE `kubernetes_install_steps` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `kubernetes_id` int(11) NOT NULL,
  `host_ip` varchar(255) ,
  `operation` varchar(255) ,
  `message` varchar(255) ,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `operation_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `operation` varchar(256) NOT NULL DEFAULT '' COMMENT '操作',
    `object_name` varchar(128) NOT NULL DEFAULT '' COMMENT '备操作的对象名,一般为表名',
    `operator` bigint(20) NOT NULL DEFAULT '0' COMMENT '操作人的 id',
    `diff` varchar(4096) DEFAULT NULL,
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='操作日志'


-- init super admin info
INSERT INTO `user`
VALUES (1, 'root', '87d9bb400c0634691f0e3baaf1e2fd0d', 1, 'enable', 1, '2021-11-09 12:29:44', '',
        '2021-11-09 12:29:44');

-- init org info
INSERT INTO `org`
VALUES (1, '星汉未来', '2021-11-16 03:44:33', '2021-11-16 03:44:33');
