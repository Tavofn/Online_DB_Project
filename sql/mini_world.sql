-- MySQL Script generated by MySQL Workbench
-- Fri Mar  3 01:36:03 2023
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';


CREATE SCHEMA IF NOT EXISTS `3380-project`;
USE `3380-project` ;

DROP TABLE IF EXISTS `3380-project`.`user` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`user` (
  `UserID` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  `date_registered` DATETIME NOT NULL,
  `name_of_user` VARCHAR(45) NOT NULL,
  `access_level` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`UserID`),
  UNIQUE INDEX `UserID_UNIQUE` (`UserID` ASC)
  );



DROP TABLE IF EXISTS `3380-project`.`album` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`album` (
  `AlbumID` INT NOT NULL,
  `release_date` VARCHAR(45) NOT NULL,
  `album_title` VARCHAR(45) NOT NULL,
  `time` INT NOT NULL,
  `average_rating` FLOAT NOT NULL,
  `UserID` INT NOT NULL,
  PRIMARY KEY (`AlbumID`),
  INDEX `FK_UserID_idx` (`UserID` ASC),
  CONSTRAINT `FK_UserID`
    FOREIGN KEY (`UserID`)
    REFERENCES `3380-project`.`user` (`UserID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);



DROP TABLE IF EXISTS `3380-project`.`artist` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`artist` (
  `ArtistID` INT NOT NULL,
  `artist_name` VARCHAR(45) NOT NULL,
  `genre` VARCHAR(45) NOT NULL,
  `average_rating` FLOAT NOT NULL,
  PRIMARY KEY (`ArtistID`),
  UNIQUE INDEX `ArtistID_UNIQUE` (`ArtistID` ASC) );



DROP TABLE IF EXISTS `3380-project`.`genre` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`genre` (
  `GenreID` INT NOT NULL,
  `genre_name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`GenreID`));



DROP TABLE IF EXISTS `3380-project`.`album_genre_artist` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`album_genre_artist` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `AlbumID` INT NOT NULL,
  `GenreID` INT NOT NULL,
  `ArtistID` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `FK_AlbumID_idx` (`AlbumID` ASC),
  INDEX `FK_GenreID_idx` (`GenreID` ASC),
  INDEX `FK_ArtistID_idx` (`ArtistID` ASC),
  CONSTRAINT `FK_AlbumID_Genre`
    FOREIGN KEY (`AlbumID`)
    REFERENCES `3380-project`.`album` (`AlbumID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_ArtistID_Genre`
    FOREIGN KEY (`ArtistID`)
    REFERENCES `3380-project`.`artist` (`ArtistID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_GenreID`
    FOREIGN KEY (`GenreID`)
    REFERENCES `3380-project`.`genre` (`GenreID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);


DROP TABLE IF EXISTS `3380-project`.`song` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`song` (
  `songID` INT NOT NULL,
  `release_date` DATETIME NOT NULL,
  `album_title` VARCHAR(45) NOT NULL,
  `time` INT NOT NULL,
  `average_rating` FLOAT NOT NULL,
  `mp3_file` VARCHAR(45) NOT NULL,
  `UserID` INT NOT NULL,
  PRIMARY KEY (`songID`),
  INDEX `UserID_idx` (`UserID` ASC),
  CONSTRAINT `UserID`
    FOREIGN KEY (`UserID`)
    REFERENCES `3380-project`.`user` (`UserID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);



DROP TABLE IF EXISTS `3380-project`.`artist_work` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`artist_work` (
  `ArtistID` INT NOT NULL,
  `songID` INT NOT NULL,
  `AlbumID` INT NOT NULL,
  PRIMARY KEY (`ArtistID`, `songID`, `AlbumID`),
  INDEX `FK_ArtistID_idx` (`ArtistID` ASC),
  INDEX `FK_SongID_idx` (`songID` ASC),
  INDEX `FK_AlbumID_idx` (`AlbumID` ASC),
  CONSTRAINT `FK_AlbumID_Album`
    FOREIGN KEY (`AlbumID`)
    REFERENCES `3380-project`.`album` (`AlbumID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_ArtistID`
    FOREIGN KEY (`ArtistID`)
    REFERENCES `3380-project`.`artist` (`ArtistID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_SongID_Songs`
    FOREIGN KEY (`songID`)
    REFERENCES `3380-project`.`song` (`songID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);

-- on line 141 and CONTRAINT, not too sure why userID was initially null, changed to not null
DROP TABLE IF EXISTS `3380-project`.`comment` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`comment` (
  `CommentID` INT NOT NULL,
  `date_written` VARCHAR(45) NOT NULL,
  `text` VARCHAR(45) NOT NULL,
  `UserID` INT NOT NULL,
  -- INT NULL DEFAULT NULL 
  PRIMARY KEY (`CommentID`),
  INDEX `FK_Comment_idx` (`UserID` ASC),
  CONSTRAINT `FK_Comment`
    FOREIGN KEY (`UserID`)
    REFERENCES `3380-project`.`user` (`UserID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);



DROP TABLE IF EXISTS `3380-project`.`playlist` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`playlist` (
  `Playlist_ID` INT NOT NULL,
  `date_created` VARCHAR(45) NOT NULL,
  `time` INT NOT NULL,
  `playlist_name` VARCHAR(45) NOT NULL,
  `UserID` INT NULL DEFAULT NULL,
  PRIMARY KEY (`Playlist_ID`),
  INDEX `FK_Playlist_idx` (`UserID` ASC),
  CONSTRAINT `FK_Playlist`
    FOREIGN KEY (`UserID`)
    REFERENCES `3380-project`.`user` (`UserID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);


DROP TABLE IF EXISTS `3380-project`.`playlist_song` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`playlist_song` (
  `playlist_id` INT NOT NULL,
  `song_id` INT NOT NULL,
  `album_id` INT NOT NULL,
  PRIMARY KEY (`playlist_id`, `song_id`, `album_id`),
  INDEX `FK_PlaylistContent_Song_idx` (`song_id` ASC),
  INDEX `FK_PlaylistContent_Album_idx` (`album_id` ASC),
  CONSTRAINT `FK_PlaylistContent_Album`
    FOREIGN KEY (`album_id`)
    REFERENCES `3380-project`.`album` (`AlbumID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_PlaylistContent_Playlist`
    FOREIGN KEY (`playlist_id`)
    REFERENCES `3380-project`.`playlist` (`Playlist_ID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_PlaylistContent_Song`
    FOREIGN KEY (`song_id`)
    REFERENCES `3380-project`.`song` (`songID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);


DROP TABLE IF EXISTS `3380-project`.`producer` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`producer` (
  `ProducerID` INT NOT NULL,
  `company` VARCHAR(45) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`ProducerID`));



DROP TABLE IF EXISTS `3380-project`.`produces` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`produces` (
  `producer_workID` INT NOT NULL,
  `ProducerID` INT NOT NULL,
  `SongID` INT NULL DEFAULT NULL,
  `AlbumID` INT NULL DEFAULT NULL,
  PRIMARY KEY (`producer_workID`),
  INDEX `FK_ProducerID_idx` (`ProducerID` ASC),
  INDEX `FK_SongID_idx` (`SongID` ASC),
  INDEX `FK_AlbumID_idx` (`AlbumID` ASC),
  CONSTRAINT `FK_AlbumID`
    FOREIGN KEY (`AlbumID`)
    REFERENCES `3380-project`.`album` (`AlbumID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_ProducerID`
    FOREIGN KEY (`ProducerID`)
    REFERENCES `3380-project`.`producer` (`ProducerID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_SongID`
    FOREIGN KEY (`SongID`)
    REFERENCES `3380-project`.`song` (`songID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);


DROP TABLE IF EXISTS `3380-project`.`rating` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`rating` (
  `ratingID` INT NOT NULL,
  `score` FLOAT NOT NULL,
  `created_date` VARCHAR(45) NOT NULL,
  `UserID` INT NULL DEFAULT NULL,
  PRIMARY KEY (`ratingID`),
  INDEX `FK_UserID_Rating_idx` (`UserID` ASC),
  CONSTRAINT `FK_UserID_Rating`
    FOREIGN KEY (`UserID`)
    REFERENCES `3380-project`.`user` (`UserID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);




DROP TABLE IF EXISTS `3380-project`.`song_comment_rating` ;
CREATE TABLE IF NOT EXISTS `3380-project`.`song_comment_rating` (
  `songID` INT NOT NULL,
  `CommentID` INT NOT NULL,
  `ratingID` INT NOT NULL,
  PRIMARY KEY (`songID`, `CommentID`, `ratingID`),
  INDEX `FK_song_comment_rating_songID_idx` (`songID` ASC),
  INDEX `FK_song_comment_rating_commentID_idx` (`CommentID` ASC),
  INDEX `FK_song_comment_rating_ratingID_idx` (`ratingID` ASC),
  CONSTRAINT `FK_song_comment_rating_commentID`
    FOREIGN KEY (`CommentID`)
    REFERENCES `3380-project`.`comment` (`CommentID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_song_comment_rating_ratingID`
    FOREIGN KEY (`ratingID`)
    REFERENCES `3380-project`.`rating` (`ratingID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `FK_song_comment_rating_songID`
    FOREIGN KEY (`songID`)
    REFERENCES `3380-project`.`song` (`songID`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);



-- SET SQL_MODE=@OLD_SQL_MODE;
-- SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
-- SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
