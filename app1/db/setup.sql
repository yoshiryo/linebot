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
    updateAt        VARCHAR(50)         NOT NULL                    
)DEFAULT CHARACTER SET=utf8;

use linebot;
