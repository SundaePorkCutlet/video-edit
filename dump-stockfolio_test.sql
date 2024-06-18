/*!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.4.2-MariaDB, for osx10.19 (arm64)
--
-- Host: localhost    Database: stockfolio_test
-- ------------------------------------------------------
-- Server version	11.2.2-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `concat_history`
--

DROP TABLE IF EXISTS `concat_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `concat_history` (
  `uuid` char(40) NOT NULL,
  `concat_video_uuid_list` text DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `concat_history`
--

LOCK TABLES `concat_history` WRITE;
/*!40000 ALTER TABLE `concat_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `concat_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `encode_history`
--

DROP TABLE IF EXISTS `encode_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `encode_history` (
  `uuid` char(40) NOT NULL,
  `origin_video_uuid` char(40) DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `encode_history`
--

LOCK TABLES `encode_history` WRITE;
/*!40000 ALTER TABLE `encode_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `encode_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `trim_history`
--

DROP TABLE IF EXISTS `trim_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trim_history` (
  `uuid` char(40) NOT NULL,
  `origin_video_uuid` char(40) DEFAULT NULL,
  `start_time` int(4) DEFAULT NULL,
  `end_time` int(4) DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `trim_history`
--

LOCK TABLES `trim_history` WRITE;
/*!40000 ALTER TABLE `trim_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `trim_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `video`
--

DROP TABLE IF EXISTS `video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `video` (
  `uuid` char(40) NOT NULL,
  `path` text DEFAULT NULL,
  `video_name` varchar(200) DEFAULT NULL,
  `extension` varchar(10) DEFAULT NULL,
  `upload_time` char(14) DEFAULT '00000000000000',
  `is_trimed` int(1) DEFAULT 0,
  `trim_time` char(14) DEFAULT '00000000000000',
  `is_concated` int(1) DEFAULT 0,
  `concat_time` char(14) DEFAULT '00000000000000',
  `is_encoded` int(1) DEFAULT 0,
  `encode_time` char(14) DEFAULT '00000000000000',
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `video`
--

LOCK TABLES `video` WRITE;
/*!40000 ALTER TABLE `video` DISABLE KEYS */;
/*!40000 ALTER TABLE `video` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'stockfolio_test'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2024-06-18 16:56:11
