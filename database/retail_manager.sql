-- MySQL dump 10.13  Distrib 8.0.44, for Linux (x86_64)
--
-- Host: localhost    Database: retail_manager
-- ------------------------------------------------------
-- Server version	8.0.44-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Categories`
--

DROP TABLE IF EXISTS `Categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Categories` (
  `category_id` binary(16) NOT NULL,
  `category_name` varchar(100) NOT NULL,
  PRIMARY KEY (`category_id`),
  UNIQUE KEY `category_name` (`category_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Categories`
--

LOCK TABLES `Categories` WRITE;
/*!40000 ALTER TABLE `Categories` DISABLE KEYS */;
INSERT INTO `Categories` VALUES (_binary 'ö•e÷Ñ\”\Ì>\“{»ùC\»\','Herbs');
/*!40000 ALTER TABLE `Categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Inventory_Log`
--

DROP TABLE IF EXISTS `Inventory_Log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Inventory_Log` (
  `log_id` binary(16) NOT NULL,
  `product_id` binary(16) NOT NULL,
  `user_id` binary(16) NOT NULL COMMENT 'ID of the user who performed the adjustment',
  `change_quantity` int NOT NULL COMMENT 'Positive (e.g., stock received), Negative (e.g., damaged/lost)',
  `reason` varchar(255) DEFAULT NULL COMMENT 'e.g., Incoming Stock, Stock Take, Damaged Return',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`log_id`),
  KEY `product_id` (`product_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `Inventory_Log_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `Products` (`product_id`),
  CONSTRAINT `Inventory_Log_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `Users` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Inventory_Log`
--

LOCK TABLES `Inventory_Log` WRITE;
/*!40000 ALTER TABLE `Inventory_Log` DISABLE KEYS */;
/*!40000 ALTER TABLE `Inventory_Log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Products`
--

DROP TABLE IF EXISTS `Products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Products` (
  `product_id` binary(16) NOT NULL,
  `product_name` varchar(150) NOT NULL,
  `purchase_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `selling_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `stock_quantity` int NOT NULL DEFAULT '0',
  `category_id` binary(16) NOT NULL,
  `supplier_id` binary(16) NOT NULL,
  PRIMARY KEY (`product_id`),
  KEY `category_id` (`category_id`),
  KEY `supplier_id` (`supplier_id`),
  CONSTRAINT `Products_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `Categories` (`category_id`) ON DELETE RESTRICT,
  CONSTRAINT `Products_ibfk_2` FOREIGN KEY (`supplier_id`) REFERENCES `Suppliers` (`supplier_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Products`
--

LOCK TABLES `Products` WRITE;
/*!40000 ALTER TABLE `Products` DISABLE KEYS */;
INSERT INTO `Products` VALUES (_binary 'ö•i\‘Q3r$eè\›Ïìè§','Wedang Uwuh - Original',8250.75,12000.00,200,_binary 'ö•e÷Ñ\”\Ì>\“{»ùC\»\',_binary 'ö•hH˘/˚ª±˚›É-\À'),(_binary 'ö•k%\È\¬\ÊQõ0EQ3Pq','Wedang Uwuh - Jahe Merah',10750.75,15000.00,250,_binary 'ö•e÷Ñ\”\Ì>\“{»ùC\»\',_binary 'ö•hH˘/˚ª±˚›É-\À'),(_binary 'ö•k£q^V-7óë0Ö','Wedang Uwuh - Bunga Telang',9250.25,12000.00,225,_binary 'ö•e÷Ñ\”\Ì>\“{»ùC\»\',_binary 'ö•hH˘/˚ª±˚›É-\À'),(_binary 'ö•l)ßåVÄS\ÌÖ{rú','Wedang Uwuh - Lemon',14500.50,18000.00,175,_binary 'ö•e÷Ñ\”\Ì>\“{»ùC\»\',_binary 'ö•hH˘/˚ª±˚›É-\À');
/*!40000 ALTER TABLE `Products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Roles`
--

DROP TABLE IF EXISTS `Roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Roles` (
  `role_id` int NOT NULL AUTO_INCREMENT,
  `role_name` varchar(50) NOT NULL COMMENT 'example: admin, cashier',
  PRIMARY KEY (`role_id`),
  UNIQUE KEY `role_name` (`role_name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Roles`
--

LOCK TABLES `Roles` WRITE;
/*!40000 ALTER TABLE `Roles` DISABLE KEYS */;
INSERT INTO `Roles` VALUES (1,'admin'),(2,'cashier');
/*!40000 ALTER TABLE `Roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Suppliers`
--

DROP TABLE IF EXISTS `Suppliers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Suppliers` (
  `supplier_id` binary(16) NOT NULL,
  `supplier_name` varchar(150) NOT NULL,
  `phone_number` varchar(20) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`supplier_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Suppliers`
--

LOCK TABLES `Suppliers` WRITE;
/*!40000 ALTER TABLE `Suppliers` DISABLE KEYS */;
INSERT INTO `Suppliers` VALUES (_binary 'ö•f¥Dõ˝F\ﬁ\›e¶\Ãr','PT. Sekar Jaya','+62-555-555-555','sekarjaya@gmail.com'),(_binary 'ö•hH˘/˚ª±˚›É-\À','PT. Rempahkarta','+62-123-456-789','rempahkarta@gmail.com');
/*!40000 ALTER TABLE `Suppliers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Transaction_Details`
--

DROP TABLE IF EXISTS `Transaction_Details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Transaction_Details` (
  `detail_id` binary(16) NOT NULL,
  `transaction_id` binary(16) NOT NULL,
  `product_id` binary(16) NOT NULL,
  `quantity` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`detail_id`),
  KEY `transaction_id` (`transaction_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `Transaction_Details_ibfk_1` FOREIGN KEY (`transaction_id`) REFERENCES `Transactions` (`transaction_id`) ON DELETE CASCADE,
  CONSTRAINT `Transaction_Details_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `Products` (`product_id`) ON DELETE RESTRICT,
  CONSTRAINT `Transaction_Details_chk_1` CHECK ((`quantity` > 0))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Transaction_Details`
--

LOCK TABLES `Transaction_Details` WRITE;
/*!40000 ALTER TABLE `Transaction_Details` DISABLE KEYS */;
INSERT INTO `Transaction_Details` VALUES (_binary 'ö•pêpæ]*µT∑d≤7',_binary 'ö•pêpæ]*µT]âä\',_binary 'ö•i\‘Q3r$eè\›Ïìè§',3,12000.00),(_binary 'ö•pêpæ]*µT\”\Ûæl',_binary 'ö•pêpæ]*µT]âä\',_binary 'ö•k%\È\¬\ÊQõ0EQ3Pq',5,15000.00),(_binary 'ö•pêpæ]*µT\Í≥\–\«',_binary 'ö•pêpæ]*µT]âä\',_binary 'ö•k£q^V-7óë0Ö',2,12000.00);
/*!40000 ALTER TABLE `Transaction_Details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Transactions`
--

DROP TABLE IF EXISTS `Transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Transactions` (
  `transaction_id` binary(16) NOT NULL,
  `transaction_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` binary(16) NOT NULL,
  PRIMARY KEY (`transaction_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `Transactions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `Users` (`user_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Transactions`
--

LOCK TABLES `Transactions` WRITE;
/*!40000 ALTER TABLE `Transactions` DISABLE KEYS */;
INSERT INTO `Transactions` VALUES (_binary 'ö•pêpæ]*µT]âä\','2025-11-21 08:03:29',_binary 'ö•[*}tˇb8˘πäò_\»');
/*!40000 ALTER TABLE `Transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `User_Roles`
--

DROP TABLE IF EXISTS `User_Roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `User_Roles` (
  `user_id` binary(16) NOT NULL,
  `role_id` int NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `User_Roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `Users` (`user_id`) ON DELETE CASCADE,
  CONSTRAINT `User_Roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `Roles` (`role_id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User_Roles`
--

LOCK TABLES `User_Roles` WRITE;
/*!40000 ALTER TABLE `User_Roles` DISABLE KEYS */;
INSERT INTO `User_Roles` VALUES (_binary 'ö•[*}tˇb8˘πäò_\»',1),(_binary 'ö•bé\∆yyÅZå˚0G',1),(_binary 'ö•c®¡ƒï\Á˙\È\Ë\Ÿ',2),(_binary 'ö•c5\ˆÑi	~ßißC\'',2),(_binary 'ö•cz\⁄“å\ZÑø\Œ]\ÓΩ',2),(_binary 'ö¶¢-w∑´ØYUª£ô\€@',2),(_binary 'ö¶ØôπEôXπM∂˛\"\“\‡',2),(_binary 'ö¶±æcû\\<\ÿ›ôZD',2),(_binary 'ö¶≥tPi†\Î\Z4–ùz∑á',2);
/*!40000 ALTER TABLE `User_Roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Users`
--

DROP TABLE IF EXISTS `Users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Users` (
  `user_id` binary(16) NOT NULL,
  `username` varchar(100) NOT NULL,
  `hashed_password` varchar(255) NOT NULL COMMENT 'password must be hashed',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Users`
--

LOCK TABLES `Users` WRITE;
/*!40000 ALTER TABLE `Users` DISABLE KEYS */;
INSERT INTO `Users` VALUES (_binary 'ö•[*}tˇb8˘πäò_\»','admin_toko','$2a$10$eCN9a5MVh5S/s/HkKNc9eu8xmwUUgNxVGfg6Njz/mXABGuf2C8QEG'),(_binary 'ö•bé\∆yyÅZå˚0G','admin_pison','$2a$10$z1GsbvQeleIGkvK3scd.UuMdCTVv7Q2SDTtyfljeSKOfIRzBpAIz6'),(_binary 'ö•c®¡ƒï\Á˙\È\Ë\Ÿ','cashier_asa','$2a$10$rhBLZ2lk2zAFWY8TdybvM.x4TK/toaJimgv9ZujkwdFBnI2EQUFXC'),(_binary 'ö•c5\ˆÑi	~ßißC\'','cashier_beni','$2a$10$o3s5krS060QeaKSCTsuSy.jSCSzYOJFJgbQUtvrvhrLUW7iVpj06e'),(_binary 'ö•cz\⁄“å\ZÑø\Œ]\ÓΩ','cashier_candra','$2a$10$u8QvbmIkCbejZUmllJoJQOtEBS5tOOB51sam.BQd/wa9HBXGxnBBC'),(_binary 'ö¶¢-w∑´ØYUª£ô\€@','cashier_deni','$2a$10$ge83RB7XzhFVjO9L6bVfKuyh7lyE1SOCbiKaovjO9AoyXk1eGL35K'),(_binary 'ö¶ØôπEôXπM∂˛\"\“\‡','cashier_eko','$2a$10$6XPJ5huXerVqwViE/ZiTC.1arNGB6P35mdsDYSwpyGusenyEsvhE.'),(_binary 'ö¶±æcû\\<\ÿ›ôZD','cashier_ferry','$2a$10$o5cvYrMkMknSMa5UijClUudecLgrRzmjTqZ0f9UrHkzuN3mOMyeJW'),(_binary 'ö¶≥tPi†\Î\Z4–ùz∑á','cashier_gerrard','$2a$10$KWEsX.AgKsaOtwoqckqXquNd.cpRMeEnLCMWzIg4orCp20mudY3rK');
/*!40000 ALTER TABLE `Users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-11-21 21:22:02
