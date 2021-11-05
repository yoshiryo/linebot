DROP DATABASE IF EXISTS linebot;
CREATE DATABASE linebot;

DROP TABLE IF EXISTS linebot.stations;

CREATE TABLE linebot.stations
(
    id              INT                 NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(20)         NOT NULL,
    first_station   VARCHAR(50)         NOT NULL,
    second_station  VARCHAR(50)         NOT NULL                      
)DEFAULT CHARACTER SET=utf8;

CREATE TABLE linebot.mangas
(
    id              INT                 NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(20)         NOT NULL,
    update_at       VARCHAR(50)         NOT NULL                    
)DEFAULT CHARACTER SET=utf8;

use linebot;
INSERT INTO stations (name, first_station, second_station) VALUES ("行き", "東青梅", "立川");
INSERT INTO mangas (name, update_at) VALUES ("妖怪戦葬", "2021-11-04");
INSERT INTO mangas (name, update_at) VALUES ("ゲシュタルト", "2021-11-04");