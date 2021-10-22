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
use linebot;
INSERT INTO stations (name, first_station, second_station) VALUES ("行き", "東青梅", "立川");