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
  `user1` int NOT NULL,
  `user2` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chats`
--

LOCK TABLES `chats` WRITE;
/*!40000 ALTER TABLE `chats` DISABLE KEYS */;
INSERT INTO `chats` VALUES (12,21,25),(13,21,34),(14,21,35),(15,21,21),(16,35,34),(17,21,32),(18,36,21),(19,34,34);
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
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES (21,25,'accepted'),(21,32,'accepted'),(21,34,'accepted'),(21,35,'accepted'),(21,36,'accepted'),(21,47,'accepted'),(25,21,'accepted'),(32,21,'accepted'),(33,21,'accepted'),(34,21,'accepted'),(35,21,'accepted'),(40,21,'accepted'),(42,21,'invited');
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
  `to_id` int NOT NULL,
  `chats_id` int NOT NULL,
  PRIMARY KEY (`id`,`from_id`,`to_id`,`chats_id`),
  KEY `fk_messeges_users1_idx` (`from_id`),
  KEY `fk_messeges_chats1_idx` (`chats_id`),
  CONSTRAINT `fk_messeges_chats1` FOREIGN KEY (`chats_id`) REFERENCES `chats` (`id`),
  CONSTRAINT `fk_messeges_users1` FOREIGN KEY (`from_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=335 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messeges`
--

LOCK TABLES `messeges` WRITE;
/*!40000 ALTER TABLE `messeges` DISABLE KEYS */;
INSERT INTO `messeges` VALUES (118,'ты кто?','2026-03-24 10:34:35',21,25,12),(119,'Я Муртуз настоящий','2026-03-24 10:34:43',25,21,12),(121,'ты астоящий муртуз?!?!???!','2026-03-24 10:35:10',21,25,12),(124,'Привет! есть ошибка в слове \"сообщение\" ','2026-03-24 17:27:32',34,21,13),(125,'ббляяя','2026-03-24 17:27:52',21,34,13),(126,'ага','2026-03-24 17:27:57',34,21,13),(128,'но смотри, все работает!!','2026-03-24 17:28:06',34,21,13),(129,'Круто!','2026-03-24 17:28:09',34,21,13),(130,'ага, я проверял на разных браузерах','2026-03-24 17:28:28',21,34,13),(131,'но это было не то','2026-03-24 17:28:32',21,34,13),(132,'пунктик добавил, что сообщение с ошибкой','2026-03-24 17:29:32',21,34,13),(133,'поменяй фотку в профиле, отобразится сразу в чате интересно','2026-03-24 17:30:04',21,34,13),(134,'ну по идее должно','2026-03-24 17:30:20',21,34,13),(135,'теперь \"сообщение\"','2026-03-24 17:30:59',21,34,13),(136,'о, да','2026-03-24 17:32:14',34,21,13),(137,'ща фотку поставлю','2026-03-24 17:32:20',34,21,13),(138,'о ','2026-03-24 17:32:56',21,34,13),(139,'пост удалила','2026-03-24 17:32:59',21,34,13),(140,'а стоп','2026-03-24 17:33:06',21,34,13),(141,'у тебя не было поста вроде','2026-03-24 17:33:18',21,34,13),(143,'да','2026-03-24 17:33:32',34,21,13),(144,'о, все я вижу котика','2026-03-24 17:33:36',21,34,13),(145,'ща напишу постик, уже продумываю слова','2026-03-24 17:33:42',34,21,13),(146,'ты с телефона кстати?','2026-03-24 17:34:34',21,34,13),(147,'очень сложный пост','2026-03-24 17:36:14',21,34,13),(148,'нет, ноут','2026-03-24 17:38:26',34,21,13),(149,'нет, это грустный пост, а не сложный','2026-03-24 17:38:51',21,34,13),(150,'ахаххахахаххах','2026-03-24 17:38:58',34,21,13),(151,'такова моя реальность ','2026-03-24 17:39:14',34,21,13),(152,'а будет функция удалить сообщение?','2026-03-24 17:39:35',34,21,13),(153,'электроника момент?','2026-03-24 17:39:35',21,34,13),(154,'тоже начинаю делать','2026-03-24 17:39:46',21,34,13),(155,'all момент','2026-03-24 17:39:48',34,21,13),(156,'были бы реакции на посты...','2026-03-24 17:40:43',21,34,13),(157,'эх блин, да, уже открываю лк, чтобы хотя бы задание сегодня прочитать..','2026-03-24 17:40:47',34,21,13),(158,'поставил бы плюсик','2026-03-24 17:40:51',21,34,13),(159,'ахахахах','2026-03-24 17:40:56',34,21,13),(160,'спасибо тебе еще раз','2026-03-24 17:41:40',21,34,13),(161,'не буду задерживать тут','2026-03-24 17:41:47',21,34,13),(163,'??','2026-03-24 17:45:41',34,21,13),(164,'не поняла','2026-03-24 17:45:43',34,21,13),(165,'ну типа сразу с козырей ','2026-03-24 17:46:58',21,34,13),(168,'да я не специально, просто заметила','2026-03-24 17:48:51',34,21,13),(169,'да не за что','2026-03-24 17:48:54',34,21,13),(170,'привет','2026-03-25 21:23:44',21,35,14),(171,'привет','2026-03-25 21:24:01',35,21,14),(172,'заебал','2026-03-25 21:24:03',21,35,14),(173,'не хочешь познакомиться','2026-03-25 21:24:07',35,21,14),(174,'хочу срать','2026-03-25 21:24:10',21,35,14),(175,'у меня вилла есть','2026-03-25 21:24:10',35,21,14),(176,'и tample','2026-03-25 21:24:12',35,21,14),(177,'temple','2026-03-25 21:24:14',35,21,14),(178,'тебе 18 есть?','2026-03-25 21:24:22',35,21,14),(180,'тогда не интересует','2026-03-25 21:24:35',35,21,14),(181,'больше не пишите','2026-03-25 21:24:39',35,21,14),(183,'не понима','2026-03-25 21:24:48',35,21,14),(202,'вфвфв','2026-03-26 21:57:05',21,25,12),(203,'вфтвфв','2026-03-26 21:57:07',21,25,12),(205,'!!!!!!!!!!!!!!!!','2026-03-26 21:59:38',21,25,12),(206,'ksfnsfs','2026-03-26 22:01:19',25,21,12),(210,'!!!!!','2026-03-27 12:38:52',21,25,12),(222,'gjhghghg','2026-03-27 13:09:11',25,21,12),(223,'ghyigibdgibibhg','2026-03-27 13:09:12',25,21,12),(231,'аыаыыа','2026-03-27 13:15:40',21,25,12),(237,'аыаыа','2026-03-27 13:22:15',25,21,12),(239,'аыаыаыа','2026-03-27 13:22:30',25,21,12),(241,'сосал?','2026-04-04 13:42:06',36,21,18),(242,'привет','2026-04-05 12:52:20',21,34,13),(243,'смотри ','2026-04-05 12:52:23',21,34,13),(244,'отправь какое нибудь сообщение, а потом измени','2026-04-05 12:52:55',21,34,13),(245,'видишь : справа','2026-04-05 12:53:06',21,34,13),(246,'Здравствуй!!','2026-04-05 12:54:24',34,21,13),(247,'о, прикольно','2026-04-05 12:54:49',34,21,13),(248,'работает исправно, даже удаление','2026-04-05 12:54:55',34,21,13),(249,'измени на контекст какой нибудь, чтобы я понял, что оно изменено','2026-04-05 12:55:00',21,34,13),(251,'увидел','2026-04-05 12:55:38',21,34,13),(252,'бля, изменения вижу если только страницу обновить ','2026-04-05 12:55:57',21,34,13),(253,'дела...','2026-04-05 12:56:01',21,34,13),(254,'а у меня сразу','2026-04-05 12:56:05',34,21,13),(255,'смотри (изменил)','2026-04-05 12:56:14',21,34,13),(256,'попробуй тоже что-нибудь изменить','2026-04-05 12:56:20',34,21,13),(257,'видишь изменения ?','2026-04-05 12:56:32',21,34,13),(258,'неа','2026-04-05 12:56:38',34,21,13),(259,'обнови','2026-04-05 12:56:41',21,34,13),(260,'да, только после обновления','2026-04-05 12:56:55',34,21,13),(261,'удаление также скорее всего','2026-04-05 12:57:11',21,34,13),(262,'а еще не совсем удобно, что нужно самому прокручивать вниз, когда написал сообщение или когда тебе пришло','2026-04-05 12:57:29',34,21,13),(263,'согласен ','2026-04-05 12:57:38',21,34,13),(264,'еще нет смысла все сообщения сразу подгружать ','2026-04-05 12:57:51',21,34,13),(265,'нужно пагинацию (не уверен, то правильно написал это слово) сделать','2026-04-05 12:58:08',21,34,13),(266,'один хер, я не знаю что это','2026-04-05 12:58:38',34,21,13),(267,'при обнове сайта диалог отображается с самого начала, приходится вниз листать','2026-04-05 12:59:07',34,21,13),(268,'а если кнопочку сделать? или она сложнее, чем просто это исправить??','2026-04-05 12:59:23',34,21,13),(269,'знаешь на некоторых сайтах с товарами снизу есть переключение на другую страницу','2026-04-05 12:59:29',21,34,13),(270,'когда нажимаешь новый контент появляется ','2026-04-05 12:59:47',21,34,13),(271,'вот то же самое только не по кнопочкам, а при прокрутке страницы сделать ','2026-04-05 13:00:02',21,34,13),(272,'это и есть пагинация ','2026-04-05 13:00:06',21,34,13),(273,'ага) багов тут много','2026-04-05 13:00:15',21,34,13),(274,'оу, спасибо, буду знать','2026-04-05 13:00:26',34,21,13),(275,'кнопочку, которая будет вниз листать ?','2026-04-05 13:00:41',21,34,13),(276,'так ты только начал, это вообще круто, что тут уже столько всего работает','2026-04-05 13:01:01',34,21,13),(277,'да','2026-04-05 13:01:03',34,21,13),(278,'кста на посты теперь фотки можно загружать','2026-04-05 13:01:36',21,34,13),(279,'заметила, ща что-нибудь загружу','2026-04-05 13:01:49',34,21,13),(280,'нууу не только начал, первый коммит 16 февраля))','2026-04-05 13:02:14',21,34,13),(281,'ну это кошмар на самом деле, по учебе ничего не успеваю ','2026-04-05 13:03:03',21,34,13),(282,'как и все наверное','2026-04-05 13:03:11',21,34,13),(283,'вот бы учеба на паузу поставилась на месяц ','2026-04-05 13:04:10',21,34,13),(284,'всем бы хорошо стало','2026-04-05 13:04:17',21,34,13),(285,'еще совсем неудобно, постоянно на профиль друга заходить','2026-04-05 13:08:26',21,34,13),(286,'чтобы написать','2026-04-05 13:08:29',21,34,13),(287,'а еще я не уверен, что все фотки будут нормально отображаться ','2026-04-05 13:14:20',21,34,13),(288,'ура вербное воскресенье ','2026-04-05 13:17:07',21,34,13),(289,'\\^o^/','2026-04-05 13:17:25',21,34,13),(290,'да, начала завидовать школьникам, у которых ща каникулы между 3 и 4 четвертью.. долгов дохера и больше','2026-04-05 13:17:27',34,21,13),(291,'завтра шкильники идут на учебу','2026-04-05 13:18:02',21,34,13),(292,'да, с учетом того, что порой нужно писать нескольким людям одновременно','2026-04-05 13:18:04',34,21,13),(293,'а','2026-04-05 13:18:08',34,21,13),(294,'ой','2026-04-05 13:18:09',34,21,13),(295,'ну у них хотя бы неделька перерыва ','2026-04-05 13:18:22',34,21,13),(296,'ну да они неделю отдыхали','2026-04-05 13:18:25',21,34,13),(297,'кста добавил песни в свой плейлист','2026-04-05 13:19:27',21,34,13),(298,'не особо люблю русские песни, но эти понравились ','2026-04-05 13:19:42',21,34,13),(299,'похоже на lofi hip hop girl, слушаешь и кайфуешь','2026-04-05 13:20:26',21,34,13),(301,'интересно получиться ли этот проект сдать, как проект на летней практике ','2026-04-05 13:23:07',21,34,13),(302,'что то выдумать, из под пальца высосать и сдать ','2026-04-05 13:23:30',21,34,13),(303,'в смысле как проект??','2026-04-05 13:24:54',34,21,13),(304,'ну у нас же практика будет на кафедре этим летом','2026-04-05 13:25:31',21,34,13),(305,'а по поводу музыки сорян, что так поздно, я просто в вк не сижу','2026-04-05 13:25:37',34,21,13),(306,'так','2026-04-05 13:25:39',34,21,13),(307,'какой нибудь проект сделать','2026-04-05 13:26:09',21,34,13),(308,'типа бота в телеге ','2026-04-05 13:26:14',21,34,13),(309,'что то связанное с кафедрой ','2026-04-05 13:26:22',21,34,13),(310,'сказать мол, ну вот смотрите, это мессенджер для 14 кафедры, чтобы студенты могли общаться с преподавателями, телега же заблокированна, а студенты в макс переходить не хотят','2026-04-05 13:27:17',21,34,13),(311,'ну и все отдыхать всю практику ','2026-04-05 13:27:40',21,34,13),(312,'на чилл выйти','2026-04-05 13:27:44',21,34,13),(313,'ладно, мне пора идти, и тебя тут задерживать не буду, спасибо больше','2026-04-05 13:29:12',21,34,13),(314,'приятно, когда можешь показать кому то, над чем пахал десятки часов, и он оценить и на ошибки укажет ','2026-04-05 13:29:54',21,34,13),(315,'следующая просьба будет не скоро, я на каникулы)))','2026-04-05 13:30:24',21,34,13),(316,'Это можно использовать как \"Избранное\", но надо будет подумать. Плюс получится ли тут вести диалог?','2026-07-04 12:03:58',34,34,19),(331,'1','2026-07-04 12:08:13',34,34,19),(333,'3','2026-07-04 12:08:15',34,34,19),(334,'1','2026-07-04 12:10:02',34,21,13);
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `login_UNIQUE` (`login`)
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (21,'clesstirss@gmail.com','$2a$10$B5AgoXb.ajpgRfErzEAdMeJdywlrKWVs/xZ2BTTLWQHAXLsB5CfWC','Ильназ','Горе разработчик, горе разрабатываю\n','/static/uploads/avatars/1783126811_Снимок экрана 2026-07-01 151442.png','2026-03-12 09:28:45','Мужской','unknown'),(25,'clesstirss0@gmail.com','$2a$10$QKQJWvL03RV7hy.qJuf4jOzi24Wtql7Ip3uiy48pgRNao8RqqKLc6','Муртуз','Кучерявый и не секси\n','static/uploads/avatars/1773534869_Снимок экрана 2025-06-20 171955.png','2026-03-12 21:49:15','Гофер','unknow'),(32,'cle@gmail.com','$2a$10$F8ngN/7N37xnlOI04fFdbOXnseGZfJb.cy6zFQrbBEBmIjznp0sXm','Тестовый пользователь','Терпеть не могу гофера','static/uploads/avatars/1773952824_1000033619.jpg','2026-03-19 20:39:08','Мужской','unknow'),(34,'elina.marchenko.2017@mail.ru','$2a$10$LeUNNiAJEY0.l7HP4NFZF.h9yoHDBf7Xz2O0ES/8YkeNU572S6J.i','Эля','Радуюсь каждым мгновением!','static/uploads/avatars/1783122507_DSC06578.JPG','2026-03-20 14:36:10','Женский','unknow'),(35,'RaymondYoung1415@hotmail.com','$2a$10$6241qoZZUBiLRUfNRugU5.PowOPGZQNMQyYpTdHdhXKlkmkU6iGeq','Jefrey Epstein','Пользователь THE NOMAX','static/uploads/avatars/1774473930_Снимок экрана 2025-09-21 102833.png','2026-03-25 21:21:19','Гофер','unknow'),(36,'alxiksanov@gmail.com','$2a$10$CeaJYpPl9XmN1/bRPbV9De4UOuE7KuIJKe9FgjB2jqe90gCX8fyS2','goddamn Стейси ','Работайте братья\n','static/uploads/avatars/1775310017_IMG_20250911_210443_662.jpg','2026-04-04 13:33:44','Мужской','unknow'),(37,'popa@gmail.ru','$2a$10$mADhynM23e3al8s1/ymHO.yYZPR2QHcCMCNBtGKTLvKV9.RI6kxEC','Мур','Пользователь THE NOMAX','unknow','2026-04-04 18:51:21','Мужской','unknow'),(39,'clesgdgdgstirss@gmail.com','$2a$10$LItGPsuEo4WIQbro6SFZNuZV7a3hLS2ii1Qvv9cOJgBSoL.wzV2f6','gdgdgdgdg','Пользователь TheNomax','unknown','2026-07-01 21:18:53','Мужской','unknow'),(40,'clesstirss25@gmail.com','$2a$10$zi66OKssAGUus.0.qrXbwe7KCZnCJZX/CtI4BIg8pfglDtI09y3NS','Ильназ1fssf','Пользователь TheNomaxsfs','/static/uploads/avatars/1783247419_Снимок экрана 2025-06-20 171920.png','2026-07-01 21:19:11','Женский','unknow'),(41,'alisa.k2005@gmail.com','$2a$10$hGdvhGB677psHjPAj9B4/uXZ5vl0RI93rACDY/wYtMCbB3tdR5CNC','кузнецова алиса андреевна ','Пользователь THE NOMAX','unknow','2026-07-02 14:54:08','Женский','unknow'),(42,'alyssa.kk@icloud.com','$2a$10$XtqAI312lTkdjo.3mUikR.R33OhT3ablaUDzfJl2HNw4frjAXa.s2','Кузнецова Алиса Андреевна','Пользователь THE NOMAX','unknow','2026-07-02 14:56:13','Женский','unknow'),(47,'negr@nigger.nigga','$2a$10$2yBdTVIjZCRJ70q0o3EkOeeUau86b6EAHLdwWo1PlPiRu9CwOQ7vi','Негр','Пользователь THE NOMAX','unknow','2026-07-03 12:01:00','Гофер','unknow');
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
) ENGINE=InnoDB AUTO_INCREMENT=112 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `wall`
--

LOCK TABLES `wall` WRITE;
/*!40000 ALTER TABLE `wall` DISABLE KEYS */;
INSERT INTO `wall` VALUES (107,'','',34,'/static/img/1783120639_photo_2026-07-02_13-21-24.jpg','2026-07-03 23:17:19'),(111,'.','.',21,'/static/uploads/posts/1783124013_Снимок экрана 2025-10-12 145339.png','2026-07-04 00:13:33');
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

-- Dump completed on 2026-07-05 14:13:58
