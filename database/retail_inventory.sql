-- MySQL dump 10.13  Distrib 8.0.44, for Linux (x86_64)
--
-- Host: localhost    Database: retail_inventory
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
-- Table structure for table `Inventory_Logs`
--

DROP TABLE IF EXISTS `Inventory_Logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Inventory_Logs` (
  `log_id` binary(16) NOT NULL,
  `product_id` binary(16) NOT NULL,
  `user_id` binary(16) NOT NULL,
  `change_quantity` int NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`log_id`),
  KEY `idx_logs_product_id` (`product_id`),
  CONSTRAINT `fk_logs_product` FOREIGN KEY (`product_id`) REFERENCES `Product_Stocks` (`product_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Inventory_Logs`
--

LOCK TABLES `Inventory_Logs` WRITE;
/*!40000 ALTER TABLE `Inventory_Logs` DISABLE KEYS */;
INSERT INTO `Inventory_Logs` VALUES (_binary 'ö§°]˚Yß°\ˆ‚èùm†',_binary 'V>:µ\”\÷vLaÔπìΩ[',_binary 'V>:µ\”\÷vLaÔπìΩ[',100,'testing aja','2025-11-20 21:17:10'),(_binary 'ö§§&´É¡øG.ô.',_binary 'V>:µ\”\÷vLaÔπìΩ[',_binary 'V>:µ\”\÷vLaÔπìΩ[',-5,'Transaction: TX-12345','2025-11-20 21:20:13'),(_binary 'ö§∫\ [\›\√\Ì8£\…\Ïã',_binary 'ö§∫\ X.∫\◊Y\ZØ«í',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',100,'init stock from monolith','2025-11-20 21:44:56'),(_binary 'ö§ºBm\"†\:\ÿ\Ë\0?/L',_binary 'ö§ºBlçKáEÑW(4c',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',100,'init stock from monolith','2025-11-20 21:46:33'),(_binary 'ö§Ω0\·M<UåE\ƒ5˚s',_binary 'ö§ºBlçKáEÑW(4c',_binary 'öî»Ñ\‚	ø7ë£ ˝>',-1,'Transaction: 01KAJBTC70RT93JKPH7KAR92DP','2025-11-20 21:47:34'),(_binary 'ö§Ω0\·M<UåE\ƒH˛ºß',_binary 'ö§∫\ X.∫\◊Y\ZØ«í',_binary 'öî»Ñ\‚	ø7ë£ ˝>',-2,'Transaction: 01KAJBTC70RT93JKPH7KAR92DP','2025-11-20 21:47:34'),(_binary 'ö•\0ë\"`∞\Âf1´øs',_binary 'ö•\0é\Ôô}n¥æQπ',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',200,'init stock from monolith','2025-11-20 23:00:38'),(_binary 'ö•¡ô¥\≈\·œú\—\Î/ùj',_binary 'ö•\0é\Ôô}n¥æQπ',_binary 'öî»Ñ\‚	ø7ë£ ˝>',10,'Incoming Stock','2025-11-20 23:02:27'),(_binary 'ö•\Ú±\∆¯)\Ï\Úé≤\—	',_binary 'ö•\0é\Ôô}n¥æQπ',_binary 'öî»Ñ\‚	ø7ë£ ˝>',-1,'Transaction: 01KAJG5WNG9XSWBJ5B078X7TRK','2025-11-20 23:03:45'),(_binary 'ö•~#´\Z†{º\œ\Ùä',_binary 'ö•\0é\Ôô}n¥æQπ',_binary 'öî»Ñ\‚	ø7ë£ ˝>',10,'Incoming Stock','2025-11-20 23:21:49'),(_binary 'ö•\ \Ûü7ÁªÄ6PF',_binary 'ö•\0é\Ôô}n¥æQπ',_binary 'öî»Ñ\‚	ø7ë£ ˝>',10,'Incoming Stock','2025-11-20 23:22:09'),(_binary 'ö•i\‘UÜº’ó≠Ü=',_binary 'ö•i\‘Q3r$eè\›Ïìè§',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',200,'init stock from monolith','2025-11-21 00:56:08'),(_binary 'ö•k%\Îäƒ¶ë¶í)\'¶é',_binary 'ö•k%\È\¬\ÊQõ0EQ3Pq',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',250,'init stock from monolith','2025-11-21 00:57:34'),(_binary 'ö•k£sØ:ú#Ç˛',_binary 'ö•k£q^V-7óë0Ö',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',225,'init stock from monolith','2025-11-21 00:58:06'),(_binary 'ö•l)®A:ˇ≠ìB^V',_binary 'ö•l)ßåVÄS\ÌÖ{rú',_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',175,'init stock from monolith','2025-11-21 00:58:41'),(_binary 'ö•o[x\È-Ü.t1d9¥',_binary 'ö•l)ßåVÄS\ÌÖ{rú',_binary 'ö•[*}tˇb8˘πäò_\»',10,'Incoming Stock','2025-11-21 01:02:10'),(_binary 'ö•pêã:ƒÄ\·kêª\Â',_binary 'ö•i\‘Q3r$eè\›Ïìè§',_binary 'ö•[*}tˇb8˘πäò_\»',-3,'Transaction: 01KAJQ143GQSEJNDAM3XERK2QG','2025-11-21 01:03:29'),(_binary 'ö•pêã:ƒÄ\‚L\Ô…∫',_binary 'ö•k%\È\¬\ÊQõ0EQ3Pq',_binary 'ö•[*}tˇb8˘πäò_\»',-5,'Transaction: 01KAJQ143GQSEJNDAM3XERK2QG','2025-11-21 01:03:29'),(_binary 'ö•pêã:ƒÄ\‚œ∫\Ûö',_binary 'ö•k£q^V-7óë0Ö',_binary 'ö•[*}tˇb8˘πäò_\»',-2,'Transaction: 01KAJQ143GQSEJNDAM3XERK2QG','2025-11-21 01:03:29');
/*!40000 ALTER TABLE `Inventory_Logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Product_Stocks`
--

DROP TABLE IF EXISTS `Product_Stocks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Product_Stocks` (
  `product_id` binary(16) NOT NULL,
  `quantity` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Product_Stocks`
--

LOCK TABLES `Product_Stocks` WRITE;
/*!40000 ALTER TABLE `Product_Stocks` DISABLE KEYS */;
INSERT INTO `Product_Stocks` VALUES (_binary 'V>:µ\”\÷vLaÔπìΩ[',95),(_binary 'ö§∫\ X.∫\◊Y\ZØ«í',98),(_binary 'ö§ºBlçKáEÑW(4c',99),(_binary 'ö•\0é\Ôô}n¥æQπ',229),(_binary 'ö•i\‘Q3r$eè\›Ïìè§',197),(_binary 'ö•k%\È\¬\ÊQõ0EQ3Pq',245),(_binary 'ö•k£q^V-7óë0Ö',223),(_binary 'ö•l)ßåVÄS\ÌÖ{rú',185);
/*!40000 ALTER TABLE `Product_Stocks` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-11-21 21:21:46
