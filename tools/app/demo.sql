-- MySQL dump 10.13  Distrib 5.7.16, for Win64 (x86_64)
--
-- Host: 192.168.99.100    Database: kolide
-- ------------------------------------------------------
-- Server version	5.7.16

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `app_configs`
--

DROP TABLE IF EXISTS `app_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_configs` (
  `id` int(10) unsigned NOT NULL DEFAULT '1',
  `org_name` varchar(255) NOT NULL DEFAULT '',
  `org_logo_url` varchar(255) NOT NULL DEFAULT '',
  `kolide_server_url` varchar(255) NOT NULL DEFAULT '',
  `smtp_configured` tinyint(1) NOT NULL DEFAULT '0',
  `smtp_sender_address` varchar(255) NOT NULL DEFAULT '',
  `smtp_server` varchar(255) NOT NULL DEFAULT '',
  `smtp_port` int(10) unsigned NOT NULL DEFAULT '587',
  `smtp_authentication_type` int(11) NOT NULL DEFAULT '0',
  `smtp_enable_ssl_tls` tinyint(1) NOT NULL DEFAULT '1',
  `smtp_authentication_method` int(11) NOT NULL DEFAULT '0',
  `smtp_domain` varchar(255) NOT NULL DEFAULT '',
  `smtp_user_name` varchar(255) NOT NULL DEFAULT '',
  `smtp_password` varchar(255) NOT NULL DEFAULT '',
  `smtp_verify_ssl_certs` tinyint(1) NOT NULL DEFAULT '1',
  `smtp_enable_start_tls` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_configs`
--

LOCK TABLES `app_configs` WRITE;
/*!40000 ALTER TABLE `app_configs` DISABLE KEYS */;
INSERT INTO `app_configs` VALUES (1,'Kolide','https://www.kolide.co/assets/kolide-nav-logo.svg','https://demo.kolide.kolide.net',0,'','',587,0,1,0,'','','',1,1);
/*!40000 ALTER TABLE `app_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `decorators`
--

DROP TABLE IF EXISTS `decorators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `decorators` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `query` text NOT NULL,
  `type` int(10) unsigned NOT NULL,
  `interval` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `decorators`
--

LOCK TABLES `decorators` WRITE;
/*!40000 ALTER TABLE `decorators` DISABLE KEYS */;
/*!40000 ALTER TABLE `decorators` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `distributed_query_campaign_targets`
--

DROP TABLE IF EXISTS `distributed_query_campaign_targets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `distributed_query_campaign_targets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `type` int(11) DEFAULT NULL,
  `distributed_query_campaign_id` int(10) unsigned DEFAULT NULL,
  `target_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `distributed_query_campaign_targets`
--

LOCK TABLES `distributed_query_campaign_targets` WRITE;
/*!40000 ALTER TABLE `distributed_query_campaign_targets` DISABLE KEYS */;
/*!40000 ALTER TABLE `distributed_query_campaign_targets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `distributed_query_campaigns`
--

DROP TABLE IF EXISTS `distributed_query_campaigns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `distributed_query_campaigns` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `query_id` int(10) unsigned DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `distributed_query_campaigns`
--

LOCK TABLES `distributed_query_campaigns` WRITE;
/*!40000 ALTER TABLE `distributed_query_campaigns` DISABLE KEYS */;
/*!40000 ALTER TABLE `distributed_query_campaigns` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `distributed_query_executions`
--

DROP TABLE IF EXISTS `distributed_query_executions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `distributed_query_executions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` int(10) unsigned DEFAULT NULL,
  `distributed_query_campaign_id` int(10) unsigned DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `error` varchar(1024) DEFAULT NULL,
  `execution_duration` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_dqe_unique_host_dqc_id` (`host_id`,`distributed_query_campaign_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `distributed_query_executions`
--

LOCK TABLES `distributed_query_executions` WRITE;
/*!40000 ALTER TABLE `distributed_query_executions` DISABLE KEYS */;
/*!40000 ALTER TABLE `distributed_query_executions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `file_integrity_monitoring_files`
--

DROP TABLE IF EXISTS `file_integrity_monitoring_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `file_integrity_monitoring_files` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `file` varchar(255) NOT NULL DEFAULT '',
  `file_integrity_monitoring_id` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fim_unique_file_name` (`file`) USING BTREE,
  KEY `fk_file_integrity_monitoring` (`file_integrity_monitoring_id`),
  CONSTRAINT `fk_file_integrity_monitoring` FOREIGN KEY (`file_integrity_monitoring_id`) REFERENCES `file_integrity_monitorings` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `file_integrity_monitoring_files`
--

LOCK TABLES `file_integrity_monitoring_files` WRITE;
/*!40000 ALTER TABLE `file_integrity_monitoring_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `file_integrity_monitoring_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `file_integrity_monitorings`
--

DROP TABLE IF EXISTS `file_integrity_monitorings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `file_integrity_monitorings` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `section_name` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_section_name` (`section_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `file_integrity_monitorings`
--

LOCK TABLES `file_integrity_monitorings` WRITE;
/*!40000 ALTER TABLE `file_integrity_monitorings` DISABLE KEYS */;
/*!40000 ALTER TABLE `file_integrity_monitorings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hosts`
--

DROP TABLE IF EXISTS `hosts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hosts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `osquery_host_id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `detail_update_time` timestamp NULL DEFAULT NULL,
  `node_key` varchar(255) DEFAULT NULL,
  `host_name` varchar(255) NOT NULL DEFAULT '',
  `uuid` varchar(255) NOT NULL DEFAULT '',
  `platform` varchar(255) NOT NULL DEFAULT '',
  `osquery_version` varchar(255) NOT NULL DEFAULT '',
  `os_version` varchar(255) NOT NULL DEFAULT '',
  `build` varchar(255) NOT NULL DEFAULT '',
  `platform_like` varchar(255) NOT NULL DEFAULT '',
  `code_name` varchar(255) NOT NULL DEFAULT '',
  `uptime` bigint(20) NOT NULL DEFAULT '0',
  `physical_memory` bigint(20) NOT NULL DEFAULT '0',
  `cpu_type` varchar(255) NOT NULL DEFAULT '',
  `cpu_subtype` varchar(255) NOT NULL DEFAULT '',
  `cpu_brand` varchar(255) NOT NULL DEFAULT '',
  `cpu_physical_cores` int(11) NOT NULL DEFAULT '0',
  `cpu_logical_cores` int(11) NOT NULL DEFAULT '0',
  `hardware_vendor` varchar(255) NOT NULL DEFAULT '',
  `hardware_model` varchar(255) NOT NULL DEFAULT '',
  `hardware_version` varchar(255) NOT NULL DEFAULT '',
  `hardware_serial` varchar(255) NOT NULL DEFAULT '',
  `computer_name` varchar(255) NOT NULL DEFAULT '',
  `primary_ip_id` int(10) unsigned DEFAULT NULL,
  `seen_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_osquery_host_id` (`osquery_host_id`),
  UNIQUE KEY `idx_host_unique_nodekey` (`node_key`),
  FULLTEXT KEY `hosts_search` (`host_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hosts`
--

LOCK TABLES `hosts` WRITE;
/*!40000 ALTER TABLE `hosts` DISABLE KEYS */;
/*!40000 ALTER TABLE `hosts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invites`
--

DROP TABLE IF EXISTS `invites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invites` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `invited_by` int(10) unsigned NOT NULL,
  `email` varchar(255) NOT NULL,
  `admin` tinyint(1) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `position` varchar(255) DEFAULT NULL,
  `token` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_invite_unique_email` (`email`),
  UNIQUE KEY `idx_invite_unique_key` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invites`
--

LOCK TABLES `invites` WRITE;
/*!40000 ALTER TABLE `invites` DISABLE KEYS */;
/*!40000 ALTER TABLE `invites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `label_query_executions`
--

DROP TABLE IF EXISTS `label_query_executions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `label_query_executions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `matches` tinyint(1) NOT NULL DEFAULT '0',
  `label_id` int(10) unsigned DEFAULT NULL,
  `host_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_lqe_label_host` (`label_id`,`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `label_query_executions`
--

LOCK TABLES `label_query_executions` WRITE;
/*!40000 ALTER TABLE `label_query_executions` DISABLE KEYS */;
/*!40000 ALTER TABLE `label_query_executions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `labels`
--

DROP TABLE IF EXISTS `labels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `labels` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `query` text NOT NULL,
  `platform` varchar(255) DEFAULT NULL,
  `label_type` int(10) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_label_unique_name` (`name`),
  FULLTEXT KEY `labels_search` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `labels`
--

LOCK TABLES `labels` WRITE;
/*!40000 ALTER TABLE `labels` DISABLE KEYS */;
INSERT INTO `labels` VALUES (1,'2017-01-18 21:41:16','2017-01-18 21:41:16',NULL,0,'All Hosts','','select 1;','',1),(2,'2017-01-18 21:41:16','2017-01-18 21:41:16',NULL,0,'Mac OS X','','select 1 from osquery_info where build_platform = \'darwin\';','darwin',1),(3,'2017-01-18 21:41:16','2017-01-18 21:41:16',NULL,0,'Ubuntu Linux','','select 1 from osquery_info where build_platform = \'ubuntu\';','ubuntu',1),(4,'2017-01-18 21:41:16','2017-01-18 21:41:16',NULL,0,'CentOS Linux','','select 1 from osquery_info where build_platform = \'centos\';','centos',1),(5,'2017-01-18 21:41:16','2017-01-18 21:41:16',NULL,0,'MS Windows','','select 1 from osquery_info where build_platform = \'windows\';','windows',1),(6,'2017-01-19 01:22:08','2017-01-19 01:23:40',NULL,0,'macOS - update needed','The macOS hosts which have not yet updated to macOS Sierra.','select * from os_version where version != \'10.12\';','',0),(7,'2017-01-19 01:25:13','2017-01-19 01:25:13',NULL,0,'Windows- update needed','Windows hosts which have not installed Windows 10.','select * from os_version where major != \'10\';','',0);
/*!40000 ALTER TABLE `labels` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `migration_status_data`
--

DROP TABLE IF EXISTS `migration_status_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migration_status_data` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `version_id` bigint(20) NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migration_status_data`
--

LOCK TABLES `migration_status_data` WRITE;
/*!40000 ALTER TABLE `migration_status_data` DISABLE KEYS */;
INSERT INTO `migration_status_data` VALUES (1,0,1,'2017-01-18 21:41:16'),(2,20161223115449,1,'2017-01-18 21:41:16'),(3,20161229171615,1,'2017-01-18 21:41:16');
/*!40000 ALTER TABLE `migration_status_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `migration_status_tables`
--

DROP TABLE IF EXISTS `migration_status_tables`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migration_status_tables` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `version_id` bigint(20) NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migration_status_tables`
--

LOCK TABLES `migration_status_tables` WRITE;
/*!40000 ALTER TABLE `migration_status_tables` DISABLE KEYS */;
INSERT INTO `migration_status_tables` VALUES (1,0,1,'2017-01-20 08:04:28'),(2,20161118193812,1,'2017-01-20 08:04:28'),(3,20161118211713,1,'2017-01-20 08:04:28'),(4,20161118212436,1,'2017-01-20 08:04:28'),(5,20161118212515,1,'2017-01-20 08:04:28'),(6,20161118212528,1,'2017-01-20 08:04:28'),(7,20161118212538,1,'2017-01-20 08:04:28'),(8,20161118212549,1,'2017-01-20 08:04:28'),(9,20161118212557,1,'2017-01-20 08:04:28'),(10,20161118212604,1,'2017-01-20 08:04:28'),(11,20161118212613,1,'2017-01-20 08:04:28'),(12,20161118212621,1,'2017-01-20 08:04:28'),(13,20161118212630,1,'2017-01-20 08:04:28'),(14,20161118212641,1,'2017-01-20 08:04:28'),(15,20161118212649,1,'2017-01-20 08:04:28'),(16,20161118212656,1,'2017-01-20 08:04:28'),(17,20161118212758,1,'2017-01-20 08:04:28'),(18,20161128234849,1,'2017-01-20 08:04:28'),(19,20161230162221,1,'2017-01-20 08:04:28'),(20,20170104113816,1,'2017-01-20 08:04:28'),(21,20170105151732,1,'2017-01-20 08:04:28'),(22,20170108191242,1,'2017-01-20 08:04:28'),(23,20170109094020,1,'2017-01-20 08:04:28'),(24,20170109130438,1,'2017-01-20 08:04:28'),(25,20170110202752,1,'2017-01-20 08:04:28'),(26,20170111133013,1,'2017-01-20 08:04:28'),(27,20170117025759,1,'2017-01-20 08:04:28');
/*!40000 ALTER TABLE `migration_status_tables` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `network_interfaces`
--

DROP TABLE IF EXISTS `network_interfaces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `network_interfaces` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` int(10) unsigned NOT NULL,
  `mac` varchar(255) NOT NULL DEFAULT '',
  `ip_address` varchar(255) NOT NULL DEFAULT '',
  `broadcast` varchar(255) NOT NULL DEFAULT '',
  `ibytes` bigint(20) NOT NULL DEFAULT '0',
  `interface` varchar(255) NOT NULL DEFAULT '',
  `ipackets` bigint(20) NOT NULL DEFAULT '0',
  `last_change` bigint(20) NOT NULL DEFAULT '0',
  `mask` varchar(255) NOT NULL DEFAULT '',
  `metric` int(11) NOT NULL DEFAULT '0',
  `mtu` int(11) NOT NULL DEFAULT '0',
  `obytes` bigint(20) NOT NULL DEFAULT '0',
  `ierrors` bigint(20) NOT NULL DEFAULT '0',
  `oerrors` bigint(20) NOT NULL DEFAULT '0',
  `opackets` bigint(20) NOT NULL DEFAULT '0',
  `point_to_point` varchar(255) NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_network_interfaces_unique_ip_host_intf` (`ip_address`,`host_id`,`interface`),
  KEY `idx_network_interfaces_hosts_fk` (`host_id`),
  FULLTEXT KEY `ip_address_search` (`ip_address`),
  CONSTRAINT `network_interfaces_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `network_interfaces`
--

LOCK TABLES `network_interfaces` WRITE;
/*!40000 ALTER TABLE `network_interfaces` DISABLE KEYS */;
/*!40000 ALTER TABLE `network_interfaces` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `options`
--

DROP TABLE IF EXISTS `options`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `options` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `type` int(10) unsigned NOT NULL,
  `value` varchar(255) NOT NULL,
  `read_only` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_option_unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `options`
--

LOCK TABLES `options` WRITE;
/*!40000 ALTER TABLE `options` DISABLE KEYS */;
INSERT INTO `options` VALUES (1,'disable_distributed',2,'false',1),(2,'distributed_plugin',0,'\"tls\"',1),(3,'distributed_tls_read_endpoint',0,'\"/api/v1/osquery/distributed/read\"',1),(4,'distributed_tls_write_endpoint',0,'\"/api/v1/osquery/distributed/write\"',1),(5,'pack_delimiter',0,'\"/\"',1),(6,'aws_access_key_id',0,'null',0),(7,'aws_firehose_period',1,'null',0),(8,'aws_firehose_stream',0,'null',0),(9,'aws_kinesis_period',1,'null',0),(10,'aws_kinesis_random_partition_key',2,'null',0),(11,'aws_kinesis_stream',0,'null',0),(12,'aws_profile_name',0,'null',0),(13,'aws_region',0,'null',0),(14,'aws_secret_access_key',0,'null',0),(15,'aws_sts_arn_role',0,'null',0),(16,'aws_sts_region',0,'null',0),(17,'aws_sts_session_name',0,'null',0),(18,'aws_sts_timeout',1,'null',0),(19,'buffered_log_max',1,'null',0),(20,'decorations_top_level',2,'null',0),(21,'disable_caching',2,'null',0),(22,'disable_database',2,'null',0),(23,'disable_decorators',2,'null',0),(24,'disable_events',2,'null',0),(25,'disable_kernel',2,'null',0),(26,'disable_logging',2,'null',0),(27,'disable_tables',0,'null',0),(28,'distributed_interval',1,'10',0),(29,'distributed_tls_max_attempts',1,'3',0),(30,'enable_foreign',2,'null',0),(31,'enable_monitor',2,'null',0),(32,'ephemeral',2,'null',0),(33,'events_expiry',1,'null',0),(34,'events_max',1,'null',0),(35,'events_optimize',2,'null',0),(36,'host_identifier',0,'null',0),(37,'logger_event_type',2,'null',0),(38,'logger_mode',0,'null',0),(39,'logger_path',0,'null',0),(40,'logger_plugin',0,'\"tls\"',0),(41,'logger_secondary_status_only',2,'null',0),(42,'logger_syslog_facility',1,'null',0),(43,'logger_tls_compress',2,'null',0),(44,'logger_tls_endpoint',0,'\"/api/v1/osquery/log\"',0),(45,'logger_tls_max',1,'null',0),(46,'logger_tls_period',1,'10',0),(47,'pack_refresh_interval',1,'null',0),(48,'read_max',1,'null',0),(49,'read_user_max',1,'null',0),(50,'schedule_default_interval',1,'null',0),(51,'schedule_splay_percent',1,'null',0),(52,'schedule_timeout',1,'null',0),(53,'utc',2,'null',0),(54,'value_max',1,'null',0),(55,'verbose',2,'null',0),(56,'worker_threads',1,'null',0);
/*!40000 ALTER TABLE `options` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pack_targets`
--

DROP TABLE IF EXISTS `pack_targets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pack_targets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pack_id` int(10) unsigned DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `target_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `constraint_pack_target_unique` (`pack_id`,`target_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pack_targets`
--

LOCK TABLES `pack_targets` WRITE;
/*!40000 ALTER TABLE `pack_targets` DISABLE KEYS */;
INSERT INTO `pack_targets` VALUES (1,1,0,1),(2,2,0,1),(3,3,0,1),(4,4,0,1),(5,5,0,1),(6,6,0,1),(7,7,0,1),(8,8,0,3),(9,8,0,4);
/*!40000 ALTER TABLE `pack_targets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `packs`
--

DROP TABLE IF EXISTS `packs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `packs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `disabled` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `platform` varchar(255) DEFAULT NULL,
  `created_by` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_pack_unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `packs`
--

LOCK TABLES `packs` WRITE;
/*!40000 ALTER TABLE `packs` DISABLE KEYS */;
INSERT INTO `packs` VALUES (1,'2017-01-19 01:07:01','2017-01-19 01:09:03',NULL,0,0,'Intrusion Detection','A collection of queries that detect indicators of initial compromise via various tactics, techniques, and procedures.','',1),(2,'2017-01-19 01:08:08','2017-01-19 01:08:08',NULL,0,0,'Osquery Monitoring','Osquery exposes several tables which allow you to query the internal operations of the osqueryd process itself. This pack contains queries that allow us to maintain insight into the health and performance of the osquery fleet.','',1),(3,'2017-01-19 01:10:38','2017-01-19 01:10:38',NULL,0,0,'Asset Management','A collection of queries that tracks the company\'s assets, installed applications, etc.','',1),(4,'2017-01-19 01:12:28','2017-01-19 01:12:28',NULL,0,0,'Hardware Monitoring','A collection of queries which monitor the changes that occur in the lower-level, hardware configurations of assets. ','',1),(5,'2017-01-19 01:13:51','2017-01-19 01:13:51',NULL,0,0,'Incident Response','While responding to an incident, it\'s often useful to have a collection of certain historical data to be able to piece together the incident timeline. This pack is a collection of queries which are useful to have during the incident response process.','',1),(6,'2017-01-19 01:14:56','2017-01-19 01:14:56',NULL,0,0,'Compliance','In order to maintain compliance, we have to ensure that we are tracking certain events and operations that occur throughout our fleet.','',1),(7,'2017-01-19 01:16:51','2017-01-19 01:16:51',NULL,0,0,'Vulnerability Management','In order to ensure that our assets are not running vulnerable versions of key software, we deploy queries within this pack to track important application values and versions.','',1),(8,'2017-01-20 08:16:52','2017-01-20 08:16:52',NULL,0,0,'Systems Monitoring','Queries which track the health, stability, and performance of a system from an operations perspective.','',1);
/*!40000 ALTER TABLE `packs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `password_reset_requests`
--

DROP TABLE IF EXISTS `password_reset_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `password_reset_requests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `expires_at` timestamp NOT NULL DEFAULT '1970-01-01 00:00:01',
  `user_id` int(10) unsigned NOT NULL,
  `token` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `password_reset_requests`
--

LOCK TABLES `password_reset_requests` WRITE;
/*!40000 ALTER TABLE `password_reset_requests` DISABLE KEYS */;
/*!40000 ALTER TABLE `password_reset_requests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `queries`
--

DROP TABLE IF EXISTS `queries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `queries` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `saved` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `query` text NOT NULL,
  `author_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_query_unique_name` (`name`),
  UNIQUE KEY `constraint_query_name_unique` (`name`),
  KEY `author_id` (`author_id`),
  CONSTRAINT `queries_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `queries`
--

LOCK TABLES `queries` WRITE;
/*!40000 ALTER TABLE `queries` DISABLE KEYS */;
INSERT INTO `queries` VALUES (1,'2017-01-20 00:53:10','2017-01-20 00:53:10',NULL,0,1,'Osquery Events','Information about osquery\'s event publishers and subscribers, which are the implementation components of event-based tables.','select * from osquery_events;',1),(2,'2017-01-20 00:53:46','2017-01-20 00:53:46',NULL,0,1,'Osquery Extensions','A list of active osquery extensions.','select * from osquery_extensions;',1),(3,'2017-01-20 00:54:27','2017-01-20 00:54:27',NULL,0,1,'Osquery Flags','The values of configurable flags which modify osquery\'s behavior. ','select * from osquery_flags;',1),(4,'2017-01-20 00:55:04','2017-01-20 00:55:04',NULL,0,1,'Osquery General Info','Top-level information about the running osquery instance.','select * from osquery_info;',1),(5,'2017-01-20 00:55:43','2017-01-20 00:55:43',NULL,0,1,'Osquery Packs','Information about the current query packs that are loaded in osquery.','select * from osquery_packs;',1),(6,'2017-01-20 00:56:18','2017-01-20 00:56:18',NULL,0,1,'Osquery Registry','Information about the active items/plugins in the osquery application registry.','select * from osquery_registry;',1),(7,'2017-01-20 00:56:42','2017-01-20 00:56:42',NULL,0,1,'Osquery Schedule','Information about the current queries that are scheduled in osquery.','select * from osquery_schedule;',1),(8,'2017-01-20 00:59:50','2017-01-20 00:59:50',NULL,0,1,'Hosts File','A line-parsed readout of the /etc/hosts file.','select * from etc_hosts;',1),(9,'2017-01-20 01:00:12','2017-01-20 01:00:12',NULL,0,1,'Protocols File','A line-parsed readout of the /etc/protocols file.','select * from etc_protocols;',1),(10,'2017-01-20 01:00:30','2017-01-20 01:00:30',NULL,0,1,'Services File','A line-parsed readout of the /etc/services file.','select * from etc_services;',1),(11,'2017-01-20 01:01:14','2017-01-20 01:01:14',NULL,0,1,'OS Version Info','Information about the operating system name and version.','select * from os_version;',1),(12,'2017-01-20 01:01:50','2017-01-20 01:01:50',NULL,0,1,'System Info','Interesting system information about a host.','select * from system_info;',1),(13,'2017-01-20 01:03:42','2017-01-20 01:03:42',NULL,0,1,'Users','Information about the users on a system and their groups.','SELECT * FROM users u JOIN groups g WHERE u.gid = g.gid;',1),(14,'2017-01-20 01:04:23','2017-01-20 01:04:23',NULL,0,1,'Windows Services','All installed Windows services and relevant data.','select * from services;',1),(15,'2017-01-20 01:04:59','2017-01-20 01:04:59',NULL,0,1,'Windows Registry','All of the Windows registry hives.','select * from registry;',1),(16,'2017-01-20 01:05:32','2017-01-20 01:05:32',NULL,0,1,'Windows Drivers','Lists all installed and loaded Windows Drivers and their relevant data.','select * from drivers;',1),(17,'2017-01-20 01:05:56','2017-01-20 01:05:56',NULL,0,1,'Windows Patches','Lists all the patches applied. Note: This does not include patches applied via MSI or downloaded from Windows Update (e.g. Service Packs).','select * from patches;',1),(18,'2017-01-20 01:12:11','2017-01-20 01:12:11',NULL,0,1,'Windows Application Compatibility Shims','Application Compatibility shims are a way to persist malware. This table presents information about the Application Compatibility Shims from the registry in a nice format.','select * from appcompat_shims;',1),(19,'2017-01-20 01:13:34','2017-01-20 01:13:34',NULL,0,1,'Kernel Info','Basic information about the active kernel.','select * from kernel_info;',1);
/*!40000 ALTER TABLE `queries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scheduled_queries`
--

DROP TABLE IF EXISTS `scheduled_queries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scheduled_queries` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `pack_id` int(10) unsigned DEFAULT NULL,
  `query_id` int(10) unsigned DEFAULT NULL,
  `interval` int(10) unsigned DEFAULT NULL,
  `snapshot` tinyint(1) DEFAULT NULL,
  `removed` tinyint(1) DEFAULT NULL,
  `platform` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `shard` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scheduled_queries`
--

LOCK TABLES `scheduled_queries` WRITE;
/*!40000 ALTER TABLE `scheduled_queries` DISABLE KEYS */;
INSERT INTO `scheduled_queries` VALUES (1,'2017-01-20 00:57:12','2017-01-20 00:57:12',NULL,0,2,1,3600,1,0,'','',NULL),(2,'2017-01-20 00:57:27','2017-01-20 00:57:27',NULL,0,2,2,3600,1,0,'','',NULL),(3,'2017-01-20 00:57:40','2017-01-20 00:57:40',NULL,0,2,3,3600,1,0,NULL,NULL,NULL),(4,'2017-01-20 00:57:48','2017-01-20 00:57:48',NULL,0,2,4,3600,1,0,NULL,NULL,NULL),(5,'2017-01-20 00:57:56','2017-01-20 00:57:56',NULL,0,2,5,3600,1,0,NULL,NULL,NULL),(6,'2017-01-20 00:58:05','2017-01-20 00:58:05',NULL,0,2,6,3600,1,0,NULL,NULL,NULL),(7,'2017-01-20 00:58:12','2017-01-20 00:58:12',NULL,0,2,7,3600,1,0,NULL,NULL,NULL);
/*!40000 ALTER TABLE `scheduled_queries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sessions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `accessed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` int(10) unsigned NOT NULL,
  `key` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_session_unique_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sessions`
--

LOCK TABLES `sessions` WRITE;
/*!40000 ALTER TABLE `sessions` DISABLE KEYS */;
INSERT INTO `sessions` VALUES (1,'2017-01-20 08:09:01','2017-01-20 08:17:09',1,'qRDbkVCGURs3Auh+3RN5SZF1umFouMQIU7LXT6mzLge04jMRT8Z+FcIfrKYyU28X7G5RkhJH3T9ee9Uby2TFQQ==');
/*!40000 ALTER TABLE `sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT '0',
  `username` varchar(255) NOT NULL,
  `password` varbinary(255) NOT NULL,
  `salt` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL DEFAULT '',
  `email` varchar(255) NOT NULL,
  `admin` tinyint(1) NOT NULL DEFAULT '0',
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `admin_forced_password_reset` tinyint(1) NOT NULL DEFAULT '0',
  `gravatar_url` varchar(255) NOT NULL DEFAULT '',
  `position` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_unique_username` (`username`),
  UNIQUE KEY `idx_user_unique_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2017-01-18 21:43:48','2017-01-18 21:44:42',NULL,0,'administrator','$2a$12$KPbYHDTqvraN72M9csSRP.SENlgc5Q10zzMH2Wlr5JCEXHwv0P0AS','iqZ/6SSoXgWezAlM7HXJ7Vph96CsYDxb','John Doe','demo-admin@kolide.co',1,1,0,'','Security Engineer');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yara_file_paths`
--

DROP TABLE IF EXISTS `yara_file_paths`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yara_file_paths` (
  `file_integrity_monitoring_id` int(11) NOT NULL DEFAULT '0',
  `yara_signature_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`file_integrity_monitoring_id`,`yara_signature_id`),
  KEY `fk_yara_signature_id` (`yara_signature_id`),
  CONSTRAINT `fk_file_integrity_monitoring_id` FOREIGN KEY (`file_integrity_monitoring_id`) REFERENCES `file_integrity_monitorings` (`id`),
  CONSTRAINT `fk_yara_signature_id` FOREIGN KEY (`yara_signature_id`) REFERENCES `yara_signatures` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yara_file_paths`
--

LOCK TABLES `yara_file_paths` WRITE;
/*!40000 ALTER TABLE `yara_file_paths` DISABLE KEYS */;
/*!40000 ALTER TABLE `yara_file_paths` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yara_signature_paths`
--

DROP TABLE IF EXISTS `yara_signature_paths`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yara_signature_paths` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_path` varchar(255) NOT NULL DEFAULT '',
  `yara_signature_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_yara_signature` (`yara_signature_id`),
  CONSTRAINT `fk_yara_signature` FOREIGN KEY (`yara_signature_id`) REFERENCES `yara_signatures` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yara_signature_paths`
--

LOCK TABLES `yara_signature_paths` WRITE;
/*!40000 ALTER TABLE `yara_signature_paths` DISABLE KEYS */;
/*!40000 ALTER TABLE `yara_signature_paths` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yara_signatures`
--

DROP TABLE IF EXISTS `yara_signatures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yara_signatures` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `signature_name` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_yara_signatures_unique_name` (`signature_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yara_signatures`
--

LOCK TABLES `yara_signatures` WRITE;
/*!40000 ALTER TABLE `yara_signatures` DISABLE KEYS */;
/*!40000 ALTER TABLE `yara_signatures` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-01-20  1:17:26
