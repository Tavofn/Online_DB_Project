-- MySQL dump 10.13  Distrib 8.0.32, for Win64 (x86_64)
--
-- Host: team3-music-database-2023.mysql.database.azure.com    Database: 3380-project
-- ------------------------------------------------------
-- Server version	5.7.40-log

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
-- Table structure for table `album`
--

DROP TABLE IF EXISTS `album`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `album` (
  `AlbumID` int(11) NOT NULL,
  `release_date` varchar(45) NOT NULL,
  `album_title` varchar(45) NOT NULL,
  `time` int(11) NOT NULL,
  `average_rating` float NOT NULL,
  `UserID` int(11) NOT NULL,
  PRIMARY KEY (`AlbumID`),
  KEY `FK_UserID_idx` (`UserID`),
  CONSTRAINT `FK_UserID` FOREIGN KEY (`UserID`) REFERENCES `user` (`UserID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `album`
--

LOCK TABLES `album` WRITE;
/*!40000 ALTER TABLE `album` DISABLE KEYS */;
/*!40000 ALTER TABLE `album` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `album_genre_artist`
--

DROP TABLE IF EXISTS `album_genre_artist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `album_genre_artist` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `AlbumID` int(11) NOT NULL,
  `GenreID` int(11) NOT NULL,
  `ArtistID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_AlbumID_idx` (`AlbumID`),
  KEY `FK_GenreID_idx` (`GenreID`),
  KEY `FK_ArtistID_idx` (`ArtistID`),
  CONSTRAINT `FK_AlbumID_Genre` FOREIGN KEY (`AlbumID`) REFERENCES `album` (`AlbumID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_ArtistID_Genre` FOREIGN KEY (`ArtistID`) REFERENCES `artist` (`ArtistID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_GenreID` FOREIGN KEY (`GenreID`) REFERENCES `genre` (`GenreID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `album_genre_artist`
--

LOCK TABLES `album_genre_artist` WRITE;
/*!40000 ALTER TABLE `album_genre_artist` DISABLE KEYS */;
/*!40000 ALTER TABLE `album_genre_artist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `artist`
--

DROP TABLE IF EXISTS `artist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `artist` (
  `ArtistID` int(11) NOT NULL,
  `artist_name` varchar(45) NOT NULL,
  `genre` varchar(45) NOT NULL,
  `average_rating` float NOT NULL,
  PRIMARY KEY (`ArtistID`),
  UNIQUE KEY `ArtistID_UNIQUE` (`ArtistID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `artist`
--

LOCK TABLES `artist` WRITE;
/*!40000 ALTER TABLE `artist` DISABLE KEYS */;
/*!40000 ALTER TABLE `artist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `artist_work`
--

DROP TABLE IF EXISTS `artist_work`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `artist_work` (
  `ArtistID` int(11) NOT NULL,
  `songID` int(11) NOT NULL,
  `AlbumID` int(11) NOT NULL,
  PRIMARY KEY (`ArtistID`,`songID`,`AlbumID`),
  KEY `FK_ArtistID_idx` (`ArtistID`),
  KEY `FK_SongID_idx` (`songID`),
  KEY `FK_AlbumID_idx` (`AlbumID`),
  CONSTRAINT `FK_AlbumID_Album` FOREIGN KEY (`AlbumID`) REFERENCES `album` (`AlbumID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_ArtistID` FOREIGN KEY (`ArtistID`) REFERENCES `artist` (`ArtistID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_SongID_Songs` FOREIGN KEY (`songID`) REFERENCES `song` (`songID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `artist_work`
--

LOCK TABLES `artist_work` WRITE;
/*!40000 ALTER TABLE `artist_work` DISABLE KEYS */;
/*!40000 ALTER TABLE `artist_work` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comment`
--

DROP TABLE IF EXISTS `comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comment` (
  `CommentID` int(11) NOT NULL,
  `date_written` varchar(45) NOT NULL,
  `text` varchar(45) NOT NULL,
  `UserID` int(11) NOT NULL,
  PRIMARY KEY (`CommentID`),
  KEY `FK_Comment_idx` (`UserID`),
  CONSTRAINT `FK_Comment` FOREIGN KEY (`UserID`) REFERENCES `user` (`UserID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comment`
--

LOCK TABLES `comment` WRITE;
/*!40000 ALTER TABLE `comment` DISABLE KEYS */;
/*!40000 ALTER TABLE `comment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `genre`
--

DROP TABLE IF EXISTS `genre`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `genre` (
  `GenreID` int(11) NOT NULL,
  `genre_name` varchar(45) NOT NULL,
  PRIMARY KEY (`GenreID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `genre`
--

LOCK TABLES `genre` WRITE;
/*!40000 ALTER TABLE `genre` DISABLE KEYS */;
/*!40000 ALTER TABLE `genre` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `playlist`
--

DROP TABLE IF EXISTS `playlist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `playlist` (
  `Playlist_ID` int(11) NOT NULL,
  `date_created` varchar(45) NOT NULL,
  `time` int(11) NOT NULL,
  `playlist_name` varchar(45) NOT NULL,
  `UserID` int(11) DEFAULT NULL,
  PRIMARY KEY (`Playlist_ID`),
  KEY `FK_Playlist_idx` (`UserID`),
  CONSTRAINT `FK_Playlist` FOREIGN KEY (`UserID`) REFERENCES `user` (`UserID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `playlist`
--

LOCK TABLES `playlist` WRITE;
/*!40000 ALTER TABLE `playlist` DISABLE KEYS */;
/*!40000 ALTER TABLE `playlist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `playlist_song`
--

DROP TABLE IF EXISTS `playlist_song`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `playlist_song` (
  `playlist_id` int(11) NOT NULL,
  `song_id` int(11) NOT NULL,
  `album_id` int(11) NOT NULL,
  PRIMARY KEY (`playlist_id`,`song_id`,`album_id`),
  KEY `FK_PlaylistContent_Song_idx` (`song_id`),
  KEY `FK_PlaylistContent_Album_idx` (`album_id`),
  CONSTRAINT `FK_PlaylistContent_Album` FOREIGN KEY (`album_id`) REFERENCES `album` (`AlbumID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_PlaylistContent_Playlist` FOREIGN KEY (`playlist_id`) REFERENCES `playlist` (`Playlist_ID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_PlaylistContent_Song` FOREIGN KEY (`song_id`) REFERENCES `song` (`songID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `playlist_song`
--

LOCK TABLES `playlist_song` WRITE;
/*!40000 ALTER TABLE `playlist_song` DISABLE KEYS */;
/*!40000 ALTER TABLE `playlist_song` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `producer`
--

DROP TABLE IF EXISTS `producer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `producer` (
  `ProducerID` int(11) NOT NULL,
  `company` varchar(45) NOT NULL,
  `name` varchar(45) NOT NULL,
  PRIMARY KEY (`ProducerID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `producer`
--

LOCK TABLES `producer` WRITE;
/*!40000 ALTER TABLE `producer` DISABLE KEYS */;
/*!40000 ALTER TABLE `producer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `produces`
--

DROP TABLE IF EXISTS `produces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `produces` (
  `producer_workID` int(11) NOT NULL,
  `ProducerID` int(11) NOT NULL,
  `SongID` int(11) DEFAULT NULL,
  `AlbumID` int(11) DEFAULT NULL,
  PRIMARY KEY (`producer_workID`),
  KEY `FK_ProducerID_idx` (`ProducerID`),
  KEY `FK_SongID_idx` (`SongID`),
  KEY `FK_AlbumID_idx` (`AlbumID`),
  CONSTRAINT `FK_AlbumID` FOREIGN KEY (`AlbumID`) REFERENCES `album` (`AlbumID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_ProducerID` FOREIGN KEY (`ProducerID`) REFERENCES `producer` (`ProducerID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_SongID` FOREIGN KEY (`SongID`) REFERENCES `song` (`songID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `produces`
--

LOCK TABLES `produces` WRITE;
/*!40000 ALTER TABLE `produces` DISABLE KEYS */;
/*!40000 ALTER TABLE `produces` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rating`
--

DROP TABLE IF EXISTS `rating`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rating` (
  `ratingID` int(11) NOT NULL,
  `score` float NOT NULL,
  `created_date` varchar(45) NOT NULL,
  `UserID` int(11) DEFAULT NULL,
  PRIMARY KEY (`ratingID`),
  KEY `FK_UserID_Rating_idx` (`UserID`),
  CONSTRAINT `FK_UserID_Rating` FOREIGN KEY (`UserID`) REFERENCES `user` (`UserID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rating`
--

LOCK TABLES `rating` WRITE;
/*!40000 ALTER TABLE `rating` DISABLE KEYS */;
/*!40000 ALTER TABLE `rating` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `song`
--

DROP TABLE IF EXISTS `song`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `song` (
  `songID` int(11) NOT NULL,
  `release_date` datetime NOT NULL,
  `album_title` varchar(45) NOT NULL,
  `time` int(11) NOT NULL,
  `average_rating` float NOT NULL,
  `mp3_file` varchar(45) NOT NULL,
  `UserID` int(11) NOT NULL,
  PRIMARY KEY (`songID`),
  KEY `UserID_idx` (`UserID`),
  CONSTRAINT `UserID` FOREIGN KEY (`UserID`) REFERENCES `user` (`UserID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `song`
--

LOCK TABLES `song` WRITE;
/*!40000 ALTER TABLE `song` DISABLE KEYS */;
/*!40000 ALTER TABLE `song` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `song_comment_rating`
--

DROP TABLE IF EXISTS `song_comment_rating`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `song_comment_rating` (
  `songID` int(11) NOT NULL,
  `CommentID` int(11) NOT NULL,
  `ratingID` int(11) NOT NULL,
  PRIMARY KEY (`songID`,`CommentID`,`ratingID`),
  KEY `FK_song_comment_rating_songID_idx` (`songID`),
  KEY `FK_song_comment_rating_commentID_idx` (`CommentID`),
  KEY `FK_song_comment_rating_ratingID_idx` (`ratingID`),
  CONSTRAINT `FK_song_comment_rating_commentID` FOREIGN KEY (`CommentID`) REFERENCES `comment` (`CommentID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_song_comment_rating_ratingID` FOREIGN KEY (`ratingID`) REFERENCES `rating` (`ratingID`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_song_comment_rating_songID` FOREIGN KEY (`songID`) REFERENCES `song` (`songID`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `song_comment_rating`
--

LOCK TABLES `song_comment_rating` WRITE;
/*!40000 ALTER TABLE `song_comment_rating` DISABLE KEYS */;
/*!40000 ALTER TABLE `song_comment_rating` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `UserID` int(11) NOT NULL,
  `username` varchar(45) NOT NULL,
  `password` varchar(45) NOT NULL,
  `date_registered` datetime NOT NULL,
  `name_of_user` varchar(45) NOT NULL,
  `access_level` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`UserID`),
  UNIQUE KEY `UserID_UNIQUE` (`UserID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-03-06 18:03:27
