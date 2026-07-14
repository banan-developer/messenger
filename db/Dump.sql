-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: messanger
-- ------------------------------------------------------
-- Server version	8.0.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `chats`
--

DROP TABLE IF EXISTS `chats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chats` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `is_group` tinyint(1) DEFAULT '0',
  `avatar_url` varchar(255) DEFAULT NULL,
  `created_by` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chats`
--

LOCK TABLES `chats` WRITE;
/*!40000 ALTER TABLE `chats` DISABLE KEYS */;
INSERT INTO `chats` VALUES (12,NULL,0,NULL,0),(13,NULL,0,NULL,0),(14,NULL,0,NULL,0),(15,NULL,0,NULL,0),(16,NULL,0,NULL,0),(17,NULL,0,NULL,0),(18,NULL,0,NULL,0),(19,NULL,0,NULL,0),(20,NULL,0,NULL,0),(21,NULL,0,NULL,0),(22,NULL,0,NULL,0),(23,NULL,0,NULL,0),(24,NULL,0,NULL,0),(26,'1442',1,'',21);
/*!40000 ALTER TABLE `chats` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friends`
--

DROP TABLE IF EXISTS `friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friends` (
  `friend_id` int NOT NULL AUTO_INCREMENT,
  `users_id` int NOT NULL,
  `status` varchar(45) NOT NULL,
  PRIMARY KEY (`friend_id`,`users_id`),
  KEY `fk_friends_users1_idx` (`users_id`),
  CONSTRAINT `fk_friends_users1` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES (21,32,'accepted'),(21,34,'accepted'),(21,35,'accepted'),(21,36,'accepted'),(21,47,'accepted'),(32,21,'accepted'),(33,21,'accepted'),(34,21,'accepted'),(35,21,'accepted'),(40,21,'accepted'),(42,21,'invited'),(48,21,'accepted');
/*!40000 ALTER TABLE `friends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `img`
--

DROP TABLE IF EXISTS `img`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `img` (
  `id_img` int NOT NULL AUTO_INCREMENT,
  `img_src` varchar(255) NOT NULL,
  `users_id` int NOT NULL,
  PRIMARY KEY (`id_img`,`users_id`),
  KEY `fk_img_users1_idx` (`users_id`),
  CONSTRAINT `fk_img_users1` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `img`
--

LOCK TABLES `img` WRITE;
/*!40000 ALTER TABLE `img` DISABLE KEYS */;
/*!40000 ALTER TABLE `img` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messeges`
--

DROP TABLE IF EXISTS `messeges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messeges` (
  `id` int NOT NULL AUTO_INCREMENT,
  `text` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `from_id` int NOT NULL,
  `to_id` int NOT NULL DEFAULT '0',
  `chats_id` int NOT NULL,
  `attachment_url` varchar(255) DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_messeges_users1_idx` (`from_id`),
  KEY `fk_messeges_chats1_idx` (`chats_id`),
  KEY `idx_chats_id` (`chats_id`),
  KEY `idx_from_id` (`from_id`),
  CONSTRAINT `fk_messeges_chats1` FOREIGN KEY (`chats_id`) REFERENCES `chats` (`id`),
  CONSTRAINT `fk_messeges_users1` FOREIGN KEY (`from_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=362 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messeges`
--

LOCK TABLES `messeges` WRITE;
/*!40000 ALTER TABLE `messeges` DISABLE KEYS */;
INSERT INTO `messeges` VALUES (118,'ты кто?','2026-03-24 10:34:35',21,25,12,NULL,0),(119,'Я Муртуз настоящий','2026-03-24 10:34:43',25,21,12,NULL,0),(121,'ты астоящий муртуз?!?!???!','2026-03-24 10:35:10',21,25,12,NULL,0),(124,'Привет! есть ошибка в слове \"сообщение\" ','2026-03-24 17:27:32',34,21,13,NULL,0),(125,'ббляяя','2026-03-24 17:27:52',21,34,13,NULL,0),(126,'ага','2026-03-24 17:27:57',34,21,13,NULL,0),(128,'но смотри, все работает!!','2026-03-24 17:28:06',34,21,13,NULL,0),(129,'Круто!','2026-03-24 17:28:09',34,21,13,NULL,0),(130,'ага, я проверял на разных браузерах','2026-03-24 17:28:28',21,34,13,NULL,0),(131,'но это было не то','2026-03-24 17:28:32',21,34,13,NULL,0),(132,'пунктик добавил, что сообщение с ошибкой','2026-03-24 17:29:32',21,34,13,NULL,0),(133,'поменяй фотку в профиле, отобразится сразу в чате интересно','2026-03-24 17:30:04',21,34,13,NULL,0),(134,'ну по идее должно','2026-03-24 17:30:20',21,34,13,NULL,0),(135,'теперь \"сообщение\"','2026-03-24 17:30:59',21,34,13,NULL,0),(136,'о, да','2026-03-24 17:32:14',34,21,13,NULL,0),(137,'ща фотку поставлю','2026-03-24 17:32:20',34,21,13,NULL,0),(138,'о ','2026-03-24 17:32:56',21,34,13,NULL,0),(139,'пост удалила','2026-03-24 17:32:59',21,34,13,NULL,0),(140,'а стоп','2026-03-24 17:33:06',21,34,13,NULL,0),(141,'у тебя не было поста вроде','2026-03-24 17:33:18',21,34,13,NULL,0),(143,'да','2026-03-24 17:33:32',34,21,13,NULL,0),(144,'о, все я вижу котика','2026-03-24 17:33:36',21,34,13,NULL,0),(145,'ща напишу постик, уже продумываю слова','2026-03-24 17:33:42',34,21,13,NULL,0),(146,'ты с телефона кстати?','2026-03-24 17:34:34',21,34,13,NULL,0),(147,'очень сложный пост','2026-03-24 17:36:14',21,34,13,NULL,0),(148,'нет, ноут','2026-03-24 17:38:26',34,21,13,NULL,0),(149,'нет, это грустный пост, а не сложный','2026-03-24 17:38:51',21,34,13,NULL,0),(150,'ахаххахахаххах','2026-03-24 17:38:58',34,21,13,NULL,0),(151,'такова моя реальность ','2026-03-24 17:39:14',34,21,13,NULL,0),(152,'а будет функция удалить сообщение?','2026-03-24 17:39:35',34,21,13,NULL,0),(153,'электроника момент?','2026-03-24 17:39:35',21,34,13,NULL,0),(154,'тоже начинаю делать','2026-03-24 17:39:46',21,34,13,NULL,0),(155,'all момент','2026-03-24 17:39:48',34,21,13,NULL,0),(156,'были бы реакции на посты...','2026-03-24 17:40:43',21,34,13,NULL,0),(157,'эх блин, да, уже открываю лк, чтобы хотя бы задание сегодня прочитать..','2026-03-24 17:40:47',34,21,13,NULL,0),(158,'поставил бы плюсик','2026-03-24 17:40:51',21,34,13,NULL,0),(159,'ахахахах','2026-03-24 17:40:56',34,21,13,NULL,0),(160,'спасибо тебе еще раз','2026-03-24 17:41:40',21,34,13,NULL,0),(161,'не буду задерживать тут','2026-03-24 17:41:47',21,34,13,NULL,0),(163,'??','2026-03-24 17:45:41',34,21,13,NULL,0),(164,'не поняла','2026-03-24 17:45:43',34,21,13,NULL,0),(165,'ну типа сразу с козырей ','2026-03-24 17:46:58',21,34,13,NULL,0),(168,'да я не специально, просто заметила','2026-03-24 17:48:51',34,21,13,NULL,0),(169,'да не за что','2026-03-24 17:48:54',34,21,13,NULL,0),(170,'привет','2026-03-25 21:23:44',21,35,14,NULL,0),(171,'привет','2026-03-25 21:24:01',35,21,14,NULL,0),(172,'заебал','2026-03-25 21:24:03',21,35,14,NULL,0),(173,'не хочешь познакомиться','2026-03-25 21:24:07',35,21,14,NULL,0),(174,'хочу срать','2026-03-25 21:24:10',21,35,14,NULL,0),(175,'у меня вилла есть','2026-03-25 21:24:10',35,21,14,NULL,0),(176,'и tample','2026-03-25 21:24:12',35,21,14,NULL,0),(177,'temple','2026-03-25 21:24:14',35,21,14,NULL,0),(178,'тебе 18 есть?','2026-03-25 21:24:22',35,21,14,NULL,0),(180,'тогда не интересует','2026-03-25 21:24:35',35,21,14,NULL,0),(181,'больше не пишите','2026-03-25 21:24:39',35,21,14,NULL,0),(183,'не понима','2026-03-25 21:24:48',35,21,14,NULL,0),(202,'вфвфв','2026-03-26 21:57:05',21,25,12,NULL,0),(203,'вфтвфв','2026-03-26 21:57:07',21,25,12,NULL,0),(205,'!!!!!!!!!!!!!!!!','2026-03-26 21:59:38',21,25,12,NULL,0),(206,'ksfnsfs','2026-03-26 22:01:19',25,21,12,NULL,0),(210,'!!!!!','2026-03-27 12:38:52',21,25,12,NULL,0),(222,'gjhghghg','2026-03-27 13:09:11',25,21,12,NULL,0),(223,'ghyigibdgibibhg','2026-03-27 13:09:12',25,21,12,NULL,0),(231,'аыаыыа','2026-03-27 13:15:40',21,25,12,NULL,0),(237,'аыаыа','2026-03-27 13:22:15',25,21,12,NULL,0),(239,'аыаыаыа','2026-03-27 13:22:30',25,21,12,NULL,0),(241,'сосал?','2026-04-04 13:42:06',36,21,18,NULL,0),(242,'привет','2026-04-05 12:52:20',21,34,13,NULL,0),(243,'смотри ','2026-04-05 12:52:23',21,34,13,NULL,0),(244,'отправь какое нибудь сообщение, а потом измени','2026-04-05 12:52:55',21,34,13,NULL,0),(245,'видишь : справа','2026-04-05 12:53:06',21,34,13,NULL,0),(246,'Здравствуй!!','2026-04-05 12:54:24',34,21,13,NULL,0),(247,'о, прикольно','2026-04-05 12:54:49',34,21,13,NULL,0),(248,'работает исправно, даже удаление','2026-04-05 12:54:55',34,21,13,NULL,0),(249,'измени на контекст какой нибудь, чтобы я понял, что оно изменено','2026-04-05 12:55:00',21,34,13,NULL,0),(251,'увидел','2026-04-05 12:55:38',21,34,13,NULL,0),(252,'бля, изменения вижу если только страницу обновить ','2026-04-05 12:55:57',21,34,13,NULL,0),(253,'дела...','2026-04-05 12:56:01',21,34,13,NULL,0),(254,'а у меня сразу','2026-04-05 12:56:05',34,21,13,NULL,0),(255,'смотри (изменил)','2026-04-05 12:56:14',21,34,13,NULL,0),(256,'попробуй тоже что-нибудь изменить','2026-04-05 12:56:20',34,21,13,NULL,0),(257,'видишь изменения ?','2026-04-05 12:56:32',21,34,13,NULL,0),(258,'неа','2026-04-05 12:56:38',34,21,13,NULL,0),(259,'обнови','2026-04-05 12:56:41',21,34,13,NULL,0),(260,'да, только после обновления','2026-04-05 12:56:55',34,21,13,NULL,0),(261,'удаление также скорее всего','2026-04-05 12:57:11',21,34,13,NULL,0),(262,'а еще не совсем удобно, что нужно самому прокручивать вниз, когда написал сообщение или когда тебе пришло','2026-04-05 12:57:29',34,21,13,NULL,0),(263,'согласен ','2026-04-05 12:57:38',21,34,13,NULL,0),(264,'еще нет смысла все сообщения сразу подгружать ','2026-04-05 12:57:51',21,34,13,NULL,0),(265,'нужно пагинацию (не уверен, то правильно написал это слово) сделать','2026-04-05 12:58:08',21,34,13,NULL,0),(266,'один хер, я не знаю что это','2026-04-05 12:58:38',34,21,13,NULL,0),(267,'при обнове сайта диалог отображается с самого начала, приходится вниз листать','2026-04-05 12:59:07',34,21,13,NULL,0),(268,'а если кнопочку сделать? или она сложнее, чем просто это исправить??','2026-04-05 12:59:23',34,21,13,NULL,0),(269,'знаешь на некоторых сайтах с товарами снизу есть переключение на другую страницу','2026-04-05 12:59:29',21,34,13,NULL,0),(270,'когда нажимаешь новый контент появляется ','2026-04-05 12:59:47',21,34,13,NULL,0),(271,'вот то же самое только не по кнопочкам, а при прокрутке страницы сделать ','2026-04-05 13:00:02',21,34,13,NULL,0),(272,'это и есть пагинация ','2026-04-05 13:00:06',21,34,13,NULL,0),(273,'ага) багов тут много','2026-04-05 13:00:15',21,34,13,NULL,0),(274,'оу, спасибо, буду знать','2026-04-05 13:00:26',34,21,13,NULL,0),(275,'кнопочку, которая будет вниз листать ?','2026-04-05 13:00:41',21,34,13,NULL,0),(276,'так ты только начал, это вообще круто, что тут уже столько всего работает','2026-04-05 13:01:01',34,21,13,NULL,0),(277,'да','2026-04-05 13:01:03',34,21,13,NULL,0),(278,'кста на посты теперь фотки можно загружать','2026-04-05 13:01:36',21,34,13,NULL,0),(279,'заметила, ща что-нибудь загружу','2026-04-05 13:01:49',34,21,13,NULL,0),(280,'нууу не только начал, первый коммит 16 февраля))','2026-04-05 13:02:14',21,34,13,NULL,0),(281,'ну это кошмар на самом деле, по учебе ничего не успеваю ','2026-04-05 13:03:03',21,34,13,NULL,0),(282,'как и все наверное','2026-04-05 13:03:11',21,34,13,NULL,0),(283,'вот бы учеба на паузу поставилась на месяц ','2026-04-05 13:04:10',21,34,13,NULL,0),(284,'всем бы хорошо стало','2026-04-05 13:04:17',21,34,13,NULL,0),(285,'еще совсем неудобно, постоянно на профиль друга заходить','2026-04-05 13:08:26',21,34,13,NULL,0),(286,'чтобы написать','2026-04-05 13:08:29',21,34,13,NULL,0),(287,'а еще я не уверен, что все фотки будут нормально отображаться ','2026-04-05 13:14:20',21,34,13,NULL,0),(288,'ура вербное воскресенье ','2026-04-05 13:17:07',21,34,13,NULL,0),(289,'\\^o^/','2026-04-05 13:17:25',21,34,13,NULL,0),(290,'да, начала завидовать школьникам, у которых ща каникулы между 3 и 4 четвертью.. долгов дохера и больше','2026-04-05 13:17:27',34,21,13,NULL,0),(291,'завтра шкильники идут на учебу','2026-04-05 13:18:02',21,34,13,NULL,0),(292,'да, с учетом того, что порой нужно писать нескольким людям одновременно','2026-04-05 13:18:04',34,21,13,NULL,0),(293,'а','2026-04-05 13:18:08',34,21,13,NULL,0),(294,'ой','2026-04-05 13:18:09',34,21,13,NULL,0),(295,'ну у них хотя бы неделька перерыва ','2026-04-05 13:18:22',34,21,13,NULL,0),(296,'ну да они неделю отдыхали','2026-04-05 13:18:25',21,34,13,NULL,0),(297,'кста добавил песни в свой плейлист','2026-04-05 13:19:27',21,34,13,NULL,0),(298,'не особо люблю русские песни, но эти понравились ','2026-04-05 13:19:42',21,34,13,NULL,0),(299,'похоже на lofi hip hop girl, слушаешь и кайфуешь','2026-04-05 13:20:26',21,34,13,NULL,0),(301,'интересно получиться ли этот проект сдать, как проект на летней практике ','2026-04-05 13:23:07',21,34,13,NULL,0),(302,'что то выдумать, из под пальца высосать и сдать ','2026-04-05 13:23:30',21,34,13,NULL,0),(303,'в смысле как проект??','2026-04-05 13:24:54',34,21,13,NULL,0),(304,'ну у нас же практика будет на кафедре этим летом','2026-04-05 13:25:31',21,34,13,NULL,0),(305,'а по поводу музыки сорян, что так поздно, я просто в вк не сижу','2026-04-05 13:25:37',34,21,13,NULL,0),(306,'так','2026-04-05 13:25:39',34,21,13,NULL,0),(307,'какой нибудь проект сделать','2026-04-05 13:26:09',21,34,13,NULL,0),(308,'типа бота в телеге ','2026-04-05 13:26:14',21,34,13,NULL,0),(309,'что то связанное с кафедрой ','2026-04-05 13:26:22',21,34,13,NULL,0),(310,'сказать мол, ну вот смотрите, это мессенджер для 14 кафедры, чтобы студенты могли общаться с преподавателями, телега же заблокированна, а студенты в макс переходить не хотят','2026-04-05 13:27:17',21,34,13,NULL,0),(311,'ну и все отдыхать всю практику ','2026-04-05 13:27:40',21,34,13,NULL,0),(312,'на чилл выйти','2026-04-05 13:27:44',21,34,13,NULL,0),(313,'ладно, мне пора идти, и тебя тут задерживать не буду, спасибо больше','2026-04-05 13:29:12',21,34,13,NULL,0),(314,'приятно, когда можешь показать кому то, над чем пахал десятки часов, и он оценить и на ошибки укажет ','2026-04-05 13:29:54',21,34,13,NULL,0),(315,'следующая просьба будет не скоро, я на каникулы)))','2026-04-05 13:30:24',21,34,13,NULL,0),(316,'Это можно использовать как \"Избранное\", но надо будет подумать. Плюс получится ли тут вести диалог?','2026-07-04 12:03:58',34,34,19,NULL,0),(331,'1','2026-07-04 12:08:13',34,34,19,NULL,0),(333,'3','2026-07-04 12:08:15',34,34,19,NULL,0),(334,'1','2026-07-04 12:10:02',34,21,13,NULL,0),(355,'dadada','2026-07-12 22:21:40',34,0,20,'/static/uploads/images/1783894900594868600.png',0),(356,'fsfsfsf','2026-07-12 22:46:58',21,0,26,'',0),(357,'ddd','2026-07-12 22:47:45',21,0,26,'',0),(358,'d','2026-07-12 22:47:48',21,0,26,'',0),(359,'выфвфвфв','2026-07-12 22:50:45',34,0,26,'',0),(360,'','2026-07-12 22:50:54',21,0,26,'/static/uploads/images/1783896654168136100.png',0),(361,'','2026-07-13 16:01:35',21,0,20,'/static/uploads/images/1783958495088806700.png',0);
/*!40000 ALTER TABLE `messeges` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `login` varchar(45) NOT NULL,
  `password` varchar(255) NOT NULL,
  `name` varchar(120) NOT NULL,
  `about` varchar(255) DEFAULT NULL,
  `avatar_url` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `sex` varchar(45) NOT NULL,
  `avatar_img` varchar(45) DEFAULT NULL,
  `role` varchar(45) DEFAULT NULL,
  `group_name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `login_UNIQUE` (`login`)
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (21,'clesstirss@gmail.com','$2a$10$B5AgoXb.ajpgRfErzEAdMeJdywlrKWVs/xZ2BTTLWQHAXLsB5CfWC','Ильназ','Горе разработчик, горе разрабатываю\n','/static/uploads/avatars/1783126811_Снимок экрана 2026-07-01 151442.png','2026-03-12 09:28:45','Мужской','unknown','Студент','1442'),(25,'clesstirss0@gmail.com','$2a$10$QKQJWvL03RV7hy.qJuf4jOzi24Wtql7Ip3uiy48pgRNao8RqqKLc6','Муртуз','Кучерявый и не секси\n','static/uploads/avatars/1773534869_Снимок экрана 2025-06-20 171955.png','2026-03-12 21:49:15','Гофер','unknow','Студент',NULL),(32,'cle@gmail.com','$2a$10$F8ngN/7N37xnlOI04fFdbOXnseGZfJb.cy6zFQrbBEBmIjznp0sXm','Тестовый пользователь','Терпеть не могу гофера','static/uploads/avatars/1773952824_1000033619.jpg','2026-03-19 20:39:08','Мужской','unknow','Преподаватель',NULL),(34,'elina.marchenko.2017@mail.ru','$2a$10$LeUNNiAJEY0.l7HP4NFZF.h9yoHDBf7Xz2O0ES/8YkeNU572S6J.i','Эля','Радуюсь каждым мгновением!','/static/uploads/avatars/1783460476_Снимок экрана 2026-03-10 114551.png','2026-03-20 14:36:10','Женский','unknow','Студент','1442'),(35,'RaymondYoung1415@hotmail.com','$2a$10$6241qoZZUBiLRUfNRugU5.PowOPGZQNMQyYpTdHdhXKlkmkU6iGeq','Jefrey Epstein','Пользователь THE NOMAX','static/uploads/avatars/1774473930_Снимок экрана 2025-09-21 102833.png','2026-03-25 21:21:19','Гофер','unknow','Студент',NULL),(36,'alxiksanov@gmail.com','$2a$10$CeaJYpPl9XmN1/bRPbV9De4UOuE7KuIJKe9FgjB2jqe90gCX8fyS2','goddamn Стейси ','Работайте братья\n','static/uploads/avatars/1775310017_IMG_20250911_210443_662.jpg','2026-04-04 13:33:44','Мужской','unknow','Студент',NULL),(37,'popa@gmail.ru','$2a$10$mADhynM23e3al8s1/ymHO.yYZPR2QHcCMCNBtGKTLvKV9.RI6kxEC','Мур','Пользователь THE NOMAX','unknow','2026-04-04 18:51:21','Мужской','unknow','Студент',NULL),(39,'clesgdgdgstirss@gmail.com','$2a$10$LItGPsuEo4WIQbro6SFZNuZV7a3hLS2ii1Qvv9cOJgBSoL.wzV2f6','gdgdgdgdg','Пользователь TheNomax','unknown','2026-07-01 21:18:53','Мужской','unknow','Студент',NULL),(40,'clesstirss25@gmail.com','$2a$10$zi66OKssAGUus.0.qrXbwe7KCZnCJZX/CtI4BIg8pfglDtI09y3NS','Ильназ1fssf','Пользователь TheNomaxsfs','/static/uploads/avatars/1783247419_Снимок экрана 2025-06-20 171920.png','2026-07-01 21:19:11','Женский','unknow','Студент',NULL),(41,'alisa.k2005@gmail.com','$2a$10$hGdvhGB677psHjPAj9B4/uXZ5vl0RI93rACDY/wYtMCbB3tdR5CNC','кузнецова алиса андреевна ','Пользователь THE NOMAX','unknow','2026-07-02 14:54:08','Женский','unknow','Студент',NULL),(42,'alyssa.kk@icloud.com','$2a$10$XtqAI312lTkdjo.3mUikR.R33OhT3ablaUDzfJl2HNw4frjAXa.s2','Кузнецова Алиса Андреевна','Пользователь THE NOMAX','unknow','2026-07-02 14:56:13','Женский','unknow','Студент',NULL),(47,'negr@nigger.nigga','$2a$10$2yBdTVIjZCRJ70q0o3EkOeeUau86b6EAHLdwWo1PlPiRu9CwOQ7vi','Негр','Пользователь THE NOMAX','unknow','2026-07-03 12:01:00','Гофер','unknow','Студент',NULL),(48,'Fclesstirss@gmail.com','$2a$10$OG/qATA3RtXFuoYHMJPuIezY0GqPA7Mn1xS8lpiBxF.cq2xz5GytG','Ахун','Пользователь TheNomax ','/static/uploads/avatars/1783289417_Снимок экрана 2025-06-20 171920.png','2026-07-05 21:32:01','Мужской','unknown','Студент',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users_has_chats`
--

DROP TABLE IF EXISTS `users_has_chats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users_has_chats` (
  `users_id` int NOT NULL,
  `chats_id` int NOT NULL,
  PRIMARY KEY (`users_id`,`chats_id`),
  KEY `fk_users_has_chats_chats1_idx` (`chats_id`),
  KEY `fk_users_has_chats_users_idx` (`users_id`),
  CONSTRAINT `fk_users_has_chats_chats1` FOREIGN KEY (`chats_id`) REFERENCES `chats` (`id`),
  CONSTRAINT `fk_users_has_chats_users` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users_has_chats`
--

LOCK TABLES `users_has_chats` WRITE;
/*!40000 ALTER TABLE `users_has_chats` DISABLE KEYS */;
INSERT INTO `users_has_chats` VALUES (21,20),(34,20),(21,21),(40,21),(21,22),(40,22),(21,23),(32,23),(21,24),(32,24),(21,26),(32,26),(34,26),(40,26);
/*!40000 ALTER TABLE `users_has_chats` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `wall`
--

DROP TABLE IF EXISTS `wall`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `wall` (
  `idwall` int NOT NULL AUTO_INCREMENT,
  `title` varchar(45) NOT NULL,
  `text` text NOT NULL,
  `users_id` int NOT NULL,
  `img_scr` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`idwall`,`users_id`),
  KEY `fk_wall_users1_idx` (`users_id`),
  CONSTRAINT `fk_wall_users1` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=118 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `wall`
--

LOCK TABLES `wall` WRITE;
/*!40000 ALTER TABLE `wall` DISABLE KEYS */;
INSERT INTO `wall` VALUES (113,'Отвечаю','.',48,'/static/uploads/posts/1783289409_1.jpg','2026-07-05 22:10:09'),(115,'Стена','стена\r\n',21,'/static/uploads/posts/1783628949_аыаыаыа.jpg','2026-07-09 20:29:09'),(116,'Важная новость','Если пишу, что я против путина, значит меня взломали',21,'/static/uploads/posts/1783629034_1.jpg','2026-07-09 20:30:34');
/*!40000 ALTER TABLE `wall` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'messanger'
--

--
-- Dumping routines for database 'messanger'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-07-14  5:12:47
